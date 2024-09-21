package initinternal

import (
	"context"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"
	"google.golang.org/protobuf/proto"

	"github.com/chenjie199234/Corelib/cerror"
	"github.com/chenjie199234/Corelib/crpc"
	"github.com/chenjie199234/Corelib/discover"
	"github.com/chenjie199234/Corelib/trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type InternalSdk struct {
	secret string
	db     *mongo.Client
	start  *primitive.Timestamp

	lker          *sync.RWMutex
	apps          map[string]*app //key: projectid-group.app
	appsNameIndex map[string]*app //key: projectname-group.app
	appsIDIndex   map[string]*app //key: ObjectID.Hex()
	closing       bool
	stopwatch     context.CancelFunc
}
type app struct {
	sync.RWMutex
	summary   *model.AppSummary
	notices   map[chan *struct{}]*struct{}
	di        discover.DI
	client    *crpc.CrpcClient
	active    int64
	delstatus bool
}

func InitWatch(secret string, db *mongo.Client) (*InternalSdk, error) {
	sdk := &InternalSdk{
		secret: secret,
		db:     db,
		start:  &primitive.Timestamp{T: uint32(time.Now().Unix() - 1), I: 0},

		lker:          &sync.RWMutex{},
		apps:          make(map[string]*app, 10), //mongodb _id为key
		appsNameIndex: make(map[string]*app, 10), //ProjectName-Group.App为key
		appsIDIndex:   make(map[string]*app, 10), //ProjectId-Group.App为key
	}
	if e := sdk.mongoGetAllApp(); e != nil {
		return nil, e
	}
	go sdk.watch()
	go sdk.timeout()
	return sdk, nil
}
func (s *InternalSdk) Stop() {
	s.lker.Lock()
	defer s.lker.Unlock()
	s.closing = true
	if s.stopwatch != nil {
		s.stopwatch()
	}
	for _, v := range s.apps {
		app := v
		app.Lock()
		app.delstatus = true
		for notice := range app.notices {
			delete(app.notices, notice)
			close(notice)
		}
		//discover should always stop after client
		if app.client != nil {
			go func() {
				app.client.Close(false)
				app.di.Stop()
			}()
		} else if app.di != nil {
			go app.di.Stop()
		}
		app.Unlock()
	}
	s.apps = make(map[string]*app, 10)
	s.appsNameIndex = make(map[string]*app, 10)
	s.appsIDIndex = make(map[string]*app, 10)
}
func (s *InternalSdk) timeout() {
	tker := time.NewTicker(time.Minute)
	for {
		now := <-tker.C
		s.lker.RLock()
		for _, v := range s.apps {
			v.Lock()
			if now.UnixNano()-v.active >= time.Minute.Nanoseconds() {
				if v.client != nil {
					oldclient := v.client
					olddi := v.di
					v.client = nil
					v.di = nil
					v.active = 0
					go func() {
						oldclient.Close(false)
						olddi.Stop()
					}()
				} else if v.di != nil {
					olddi := v.di
					v.di = nil
					v.active = 0
					go olddi.Stop()
				}
			}
			v.Unlock()
		}
		s.lker.RUnlock()
	}
}
func createdi(summary *model.AppSummary) (discover.DI, error) {
	switch summary.DiscoverMode {
	case "kubernetes":
		if summary.KubernetesNs == "" || summary.KubernetesLS == "" {
			slog.ErrorContext(nil, "[InitWatch] discover info broken",
				slog.String("project_id", summary.ProjectID),
				slog.String("project_name", summary.ProjectName),
				slog.String("group", summary.Group),
				slog.String("app", summary.App))
			return nil, ecode.ErrDBDataBroken
		} else {
			di, e := discover.NewKubernetesDiscover(
				summary.ProjectName,
				summary.Group,
				summary.App,
				summary.KubernetesNs,
				summary.KubernetesFS,
				summary.KubernetesLS,
				int(summary.CrpcPort),
				int(summary.CGrpcPort),
				int(summary.WebPort))
			if e != nil {
				slog.ErrorContext(nil, "[InitWatch] new discover failed",
					slog.String("project_id", summary.ProjectID),
					slog.String("project_name", summary.ProjectName),
					slog.String("group", summary.Group),
					slog.String("app", summary.App),
					slog.String("error", e.Error()))
				return nil, ecode.ErrDBDataBroken
			}
			return di, nil
		}
	case "dns":
		if summary.DnsHost == "" || summary.DnsInterval <= 0 {
			slog.ErrorContext(nil, "[InitWatch] discover info broken",
				slog.String("project_id", summary.ProjectID),
				slog.String("project_name", summary.ProjectName),
				slog.String("group", summary.Group),
				slog.String("app", summary.App))
			return nil, ecode.ErrDBDataBroken
		} else {
			interval := time.Duration(summary.DnsInterval) * time.Second
			di, e := discover.NewDNSDiscover(
				summary.ProjectName,
				summary.Group,
				summary.App,
				summary.DnsHost,
				interval,
				int(summary.CrpcPort),
				int(summary.CGrpcPort),
				int(summary.WebPort))
			if e != nil {
				slog.ErrorContext(nil, "[InitWatch] new discover failed",
					slog.String("project_id", summary.ProjectID),
					slog.String("project_name", summary.ProjectName),
					slog.String("group", summary.Group),
					slog.String("app", summary.App),
					slog.String("error", e.Error()))
				return nil, ecode.ErrDBDataBroken
			}
			return di, nil
		}
	case "static":
		if len(summary.StaticAddrs) == 0 {
			slog.ErrorContext(nil, "[InitWatch] discover info broken",
				slog.String("project_id", summary.ProjectID),
				slog.String("project_name", summary.ProjectName),
				slog.String("group", summary.Group),
				slog.String("app", summary.App))
			return nil, ecode.ErrDBDataBroken
		} else {
			di, e := discover.NewStaticDiscover(
				summary.ProjectName,
				summary.Group,
				summary.App,
				summary.StaticAddrs,
				int(summary.CrpcPort),
				int(summary.CGrpcPort),
				int(summary.WebPort))
			if e != nil {
				slog.ErrorContext(nil, "[InitWatch] new discover failed",
					slog.String("project_id", summary.ProjectID),
					slog.String("project_name", summary.ProjectName),
					slog.String("group", summary.Group),
					slog.String("app", summary.App),
					slog.String("error", e.Error()))
				return nil, ecode.ErrDBDataBroken
			}
			return di, nil
		}
	default:
		return nil, ecode.ErrDBDataBroken
	}
}
func (s *InternalSdk) mongoGetAllApp() error {
	filter := bson.M{"key": "", "index": 0}
	var cursor *mongo.Cursor
	primarylocalOPTS := options.Collection().SetReadPreference(readpref.Primary()).SetReadConcern(readconcern.Local())
	cursor, e := s.db.Database("app").Collection("config", primarylocalOPTS).Find(context.Background(), filter)
	if e != nil {
		slog.ErrorContext(nil, "[InitWatch] get all app config failed", slog.String("error", e.Error()))
		return e
	}
	defer cursor.Close(context.Background())
	apps := make([]*model.AppSummary, 0, cursor.RemainingBatchLength())
	if e := cursor.All(context.Background(), &apps); e != nil {
		slog.ErrorContext(nil, "[InitWatch] get all app config failed", slog.String("error", e.Error()))
		return e
	}
	for _, v := range apps {
		undup := make(map[string]*struct{}, len(v.StaticAddrs))
		for _, v := range v.StaticAddrs {
			undup[v] = nil
		}
		v.StaticAddrs = make([]string, 0, len(undup))
		for k := range undup {
			v.StaticAddrs = append(v.StaticAddrs, k)
		}
		tmp := &app{
			summary: v,
			notices: make(map[chan *struct{}]*struct{}, 10),
		}
		s.apps[v.ID.Hex()] = tmp
		s.appsNameIndex[v.ProjectName+"-"+v.Group+"."+v.App] = tmp
		s.appsIDIndex[v.ProjectID+"-"+v.Group+"."+v.App] = tmp
	}
	return nil
}
func (s *InternalSdk) watch() {
	var stream *mongo.ChangeStream
	ctx, span := trace.NewSpan(context.Background(), "InitWatch", trace.Client, nil)
	defer span.Finish(nil)
	ctx, s.stopwatch = context.WithCancel(ctx)
	for {
		for stream == nil {
			if s.closing {
				return
			}
			//connect
			var e error
			opts := options.ChangeStream().SetFullDocument(options.UpdateLookup).SetStartAtOperationTime(s.start)
			if stream, e = s.db.Database("app").Collection("config").Watch(ctx, mongo.Pipeline{}, opts); e != nil {
				slog.ErrorContext(nil, "[InitWatch] get stream failed", slog.String("error", e.Error()))
				stream = nil
				time.Sleep(time.Millisecond * 100)
				continue
			}
		}
		for stream.Next(ctx) {
			s.start.T, s.start.I = stream.Current.Lookup("clusterTime").Timestamp()
			s.start.I++
			switch stream.Current.Lookup("operationType").StringValue() {
			case "drop":
				//drop collection
				slog.ErrorContext(nil, "[InitWatch] all configs deleted")
				s.lker.Lock()
				for _, v := range s.apps {
					app := v
					app.Lock()
					app.delstatus = true
					for notice := range app.notices {
						delete(app.notices, notice)
						close(notice)
					}
					//discover should always stop after client
					if app.client != nil {
						go func() {
							app.client.Close(false)
							app.di.Stop()
						}()
					} else if app.di != nil {
						go app.di.Stop()
					}
					app.Unlock()
				}
				s.apps = make(map[string]*app, 10)
				s.appsNameIndex = make(map[string]*app, 10)
				s.appsIDIndex = make(map[string]*app, 10)
				s.lker.Unlock()
			case "insert":
				//insert document
				fallthrough
			case "update":
				//update document
				projectid, pok := stream.Current.Lookup("fullDocument").Document().Lookup("project_id").StringValueOK()
				gname, gok := stream.Current.Lookup("fullDocument").Document().Lookup("group").StringValueOK()
				aname, aok := stream.Current.Lookup("fullDocument").Document().Lookup("app").StringValueOK()
				key, kok := stream.Current.Lookup("fullDocument").Document().Lookup("key").StringValueOK()
				index, iok := stream.Current.Lookup("fullDocument").Document().Lookup("index").AsInt64OK()
				if !pok || !gok || !aok || !kok || !iok {
					//unknown doc
					continue
				}
				if key != "" || index != 0 {
					//this is not the app summary
					continue
				}
				//this is the app summary
				summary := &model.AppSummary{}
				if e := stream.Current.Lookup("fullDocument").Unmarshal(summary); e != nil {
					slog.ErrorContext(nil, "[InitWatch] document format wrong",
						slog.String("project_id", projectid),
						slog.String("group", gname),
						slog.String("app", aname),
						slog.String("error", e.Error()))
					continue
				}
				slog.DebugContext(nil, "[InitWatch] updated",
					slog.String("project_id", summary.ProjectID),
					slog.String("group", summary.Group),
					slog.String("app", summary.App),
					slog.Any("keys", summary.Keys))
				s.lker.Lock()
				if exist, ok := s.apps[summary.ID.Hex()]; !ok {
					//this is a new app
					tmp := &app{
						summary: summary,
						notices: make(map[chan *struct{}]*struct{}, 10),
					}
					s.apps[summary.ID.Hex()] = tmp
					s.appsNameIndex[summary.ProjectName+"-"+summary.Group+"."+summary.App] = tmp
					s.appsIDIndex[summary.ProjectID+"-"+summary.Group+"."+summary.App] = tmp
				} else {
					//this is an old app
					exist.Lock()
					if exist.summary.ProjectName != summary.ProjectName {
						//ProjectName changed
						delete(s.appsNameIndex, exist.summary.ProjectName+"-"+exist.summary.Group+"."+exist.summary.App)
						s.appsNameIndex[summary.ProjectName+"-"+summary.Group+"."+summary.App] = exist
					}
					if exist.summary.ProjectID != summary.ProjectID {
						//ProjectID changed
						delete(s.appsIDIndex, exist.summary.ProjectID+"-"+exist.summary.Group+"."+exist.summary.App)
						s.appsIDIndex[summary.ProjectID+"-"+summary.Group+"."+summary.App] = exist
					}
					undup := make(map[string]*struct{}, len(summary.StaticAddrs))
					for _, v := range summary.StaticAddrs {
						undup[v] = nil
					}
					summary.StaticAddrs = make([]string, 0, len(undup))
					for k := range undup {
						summary.StaticAddrs = append(summary.StaticAddrs, k)
					}
					discoverchanged := exist.summary.DiscoverMode != summary.DiscoverMode ||
						exist.summary.CrpcPort != summary.CrpcPort ||
						exist.summary.CGrpcPort != summary.CGrpcPort ||
						exist.summary.WebPort != summary.WebPort ||
						exist.summary.KubernetesNs != summary.KubernetesNs ||
						exist.summary.KubernetesLS != summary.KubernetesLS ||
						exist.summary.KubernetesFS != summary.KubernetesFS ||
						exist.summary.DnsHost != summary.DnsHost ||
						exist.summary.DnsInterval != summary.DnsInterval ||
						len(exist.summary.StaticAddrs) != len(summary.StaticAddrs)
					if !discoverchanged {
						for _, a := range exist.summary.StaticAddrs {
							find := false
							for _, b := range summary.StaticAddrs {
								if a == b {
									find = true
									break
								}
							}
							if !find {
								discoverchanged = true
								break
							}
						}
					}
					if discoverchanged {
						slog.DebugContext(nil, "[InitWatch] discover changed",
							slog.String("project_id", projectid),
							slog.String("group", gname),
							slog.String("app", aname))
						//discover should always stop after client
						if exist.client != nil {
							oldclient := exist.client
							olddi := exist.di
							exist.client = nil
							exist.di = nil
							exist.active = 0
							go func() {
								oldclient.Close(false)
								olddi.Stop()
							}()
						} else if exist.di != nil {
							olddi := exist.di
							exist.di = nil
							exist.active = 0
							go olddi.Stop()
						}
					}
					exist.summary = summary
					for notice := range exist.notices {
						select {
						case notice <- nil:
						default:
						}
					}
					exist.Unlock()
				}
				s.lker.Unlock()
			case "delete":
				//delete document
				objid := stream.Current.Lookup("documentKey").Document().Lookup("_id").ObjectID().Hex()
				s.lker.Lock()
				exist, ok := s.apps[objid]
				if !ok {
					//this is not the app summary
					s.lker.Unlock()
					break
				}
				//this is the app summary
				delete(s.apps, objid)
				delete(s.appsNameIndex, exist.summary.ProjectName+"-"+exist.summary.Group+"."+exist.summary.App)
				delete(s.appsIDIndex, exist.summary.ProjectID+"-"+exist.summary.Group+"."+exist.summary.App)
				exist.Lock()
				s.lker.Unlock()
				slog.DebugContext(nil, "[InitWatch] deleted",
					slog.String("project_id", exist.summary.ProjectID),
					slog.String("group", exist.summary.Group),
					slog.String("app", exist.summary.App))
				exist.delstatus = true
				for notice := range exist.notices {
					delete(exist.notices, notice)
					close(notice)
				}
				//discover should always stop after client
				if exist.client != nil {
					go func() {
						exist.client.Close(false)
						exist.di.Stop()
					}()
				} else if exist.di != nil {
					go exist.di.Stop()
				}
				exist.Unlock()
			}
		}
		if stream.Err() != nil {
			slog.ErrorContext(nil, "[InitWatch] stream disconnected", slog.String("error", stream.Err().Error()))
		}
		stream.Close(nil)
		stream = nil
	}
}

// if you don't need the notice,remember to call the cancel
func (s *InternalSdk) GetNoticeByProjectID(pid, g, a string) (notice <-chan *struct{}, cancel func(), e error) {
	fullname := pid + "-" + g + "." + a
	s.lker.Lock()
	if s.closing {
		s.lker.Unlock()
		return nil, nil, ecode.ErrServerClosing
	}
	app, ok := s.appsIDIndex[fullname]
	if !ok {
		s.lker.Unlock()
		return nil, nil, ecode.ErrAppNotExist
	}
	app.Lock()
	s.lker.Unlock()
	defer app.Unlock()
	ch := make(chan *struct{}, 1)
	ch <- nil
	app.notices[ch] = nil
	return ch, func() {
		app.Lock()
		defer app.Unlock()
		if _, ok := app.notices[ch]; ok {
			delete(app.notices, ch)
			close(ch)
		}
	}, nil
}

// if you don't need the notice,remember to call the cancel
func (s *InternalSdk) GetNoticeByProjectName(pname, g, a string) (notice <-chan *struct{}, cancel func(), e error) {
	fullname := pname + "-" + g + "." + a
	s.lker.Lock()
	if s.closing {
		s.lker.Unlock()
		return nil, nil, ecode.ErrServerClosing
	}
	app, ok := s.appsNameIndex[fullname]
	if !ok {
		s.lker.Unlock()
		return nil, nil, ecode.ErrAppNotExist
	}
	app.Lock()
	s.lker.Unlock()
	defer app.Unlock()
	ch := make(chan *struct{}, 1)
	ch <- nil
	app.notices[ch] = nil
	return ch, func() {
		app.Lock()
		defer app.Unlock()
		if _, ok := app.notices[ch]; ok {
			delete(app.notices, ch)
			close(ch)
		}
	}, nil
}
func (s *InternalSdk) GetAppConfigByProjectID(pid, g, a string) (*model.AppSummary, error) {
	fullname := pid + "-" + g + "." + a
	s.lker.RLock()
	if s.closing {
		s.lker.RUnlock()
		return nil, ecode.ErrServerClosing
	}
	app, ok := s.appsIDIndex[fullname]
	if !ok {
		s.lker.RUnlock()
		return nil, ecode.ErrAppNotExist
	}
	app.RLock()
	s.lker.RUnlock()
	defer app.RUnlock()
	if app.delstatus {
		return nil, ecode.ErrAppNotExist
	}
	return app.summary, nil
}
func (s *InternalSdk) GetAppConfigByProjectName(pname, g, a string) (*model.AppSummary, error) {
	fullname := pname + "-" + g + "." + a
	s.lker.RLock()
	if s.closing {
		s.lker.RUnlock()
		return nil, ecode.ErrServerClosing
	}
	app, ok := s.appsNameIndex[fullname]
	if !ok {
		s.lker.RUnlock()
		return nil, ecode.ErrAppNotExist
	}
	app.RLock()
	s.lker.RUnlock()
	defer app.RUnlock()
	if app.delstatus {
		return nil, ecode.ErrAppNotExist
	}
	return app.summary, nil
}

func (s *InternalSdk) GetAppAddrsByProjectID(ctx context.Context, pid, g, a string) ([]string, error) {
	fullname := pid + "-" + g + "." + a
	s.lker.RLock()
	if s.closing {
		s.lker.RUnlock()
		return nil, ecode.ErrServerClosing
	}
	app, ok := s.appsIDIndex[fullname]
	if !ok {
		s.lker.RUnlock()
		return nil, ecode.ErrAppNotExist
	}
	app.Lock()
	s.lker.RUnlock()
	if app.delstatus {
		app.Unlock()
		return nil, ecode.ErrAppNotExist
	}
	if app.di == nil {
		var e error
		app.di, e = createdi(app.summary)
		if e != nil {
			app.Unlock()
			return nil, e
		}
	}
	app.active = time.Now().UnixNano()
	//copy the di pointer
	di := app.di
	app.Unlock()
	ch, cancel := di.GetNotice()
	defer cancel()
	select {
	case <-ch:
	case <-ctx.Done():
		return nil, cerror.Convert(ctx.Err())
	}
	tmp, _, _ := di.GetAddrs(discover.NotNeed)
	addrs := make([]string, 0, len(tmp))
	for addr := range tmp {
		addrs = append(addrs, addr)
	}
	return addrs, nil
}

func (s *InternalSdk) GetAppAddrsByProjectName(ctx context.Context, pname, g, a string) ([]string, error) {
	fullname := pname + "-" + g + "." + a
	s.lker.RLock()
	if s.closing {
		s.lker.RUnlock()
		return nil, ecode.ErrServerClosing
	}
	app, ok := s.appsNameIndex[fullname]
	if !ok {
		s.lker.RUnlock()
		return nil, ecode.ErrAppNotExist
	}
	app.Lock()
	s.lker.RUnlock()
	if app.delstatus {
		app.Unlock()
		return nil, ecode.ErrAppNotExist
	}
	if app.di == nil {
		var e error
		app.di, e = createdi(app.summary)
		if e != nil {
			app.Unlock()
			return nil, e
		}
	}
	app.active = time.Now().UnixNano()
	//copy the di pointer
	di := app.di
	app.Unlock()
	ch, cancel := di.GetNotice()
	defer cancel()
	select {
	case <-ch:
	case <-ctx.Done():
		return nil, cerror.Convert(ctx.Err())
	}
	tmp, _, _ := di.GetAddrs(discover.NotNeed)
	addrs := make([]string, 0, len(tmp))
	for addr := range tmp {
		addrs = append(addrs, addr)
	}
	return addrs, nil
}

func (s *InternalSdk) PingByPrjoectID(ctx context.Context, pid, g, a string, forceaddr string) (*api.Pingresp, error) {
	fullname := pid + "-" + g + "." + a
	s.lker.RLock()
	if s.closing {
		s.lker.RUnlock()
		return nil, ecode.ErrServerClosing
	}
	app, ok := s.appsIDIndex[fullname]
	if !ok {
		s.lker.RUnlock()
		return nil, ecode.ErrAppNotExist
	}
	app.Lock()
	s.lker.RUnlock()
	if app.delstatus {
		app.Unlock()
		return nil, ecode.ErrAppNotExist
	}
	if app.client == nil {
		var e error
		if app.di == nil {
			app.di, e = createdi(app.summary)
			if e != nil {
				app.Unlock()
				return nil, e
			}
		}
		app.client, e = crpc.NewCrpcClient(nil, app.di, model.Project, model.Group, model.Name, app.summary.ProjectName, app.summary.Group, app.summary.App, nil)
		if e != nil {
			app.Unlock()
			return nil, e
		}
	}
	app.active = time.Now().UnixNano()
	if forceaddr != "" && app.summary.CrpcPort != 0 {
		if strings.Contains(forceaddr, ":") {
			//ipv6
			forceaddr = "[" + forceaddr + "]:" + strconv.FormatUint(uint64(app.summary.CrpcPort), 10)
		} else {
			//ipv4 or host
			forceaddr = forceaddr + ":" + strconv.FormatUint(uint64(app.summary.CrpcPort), 10)
		}
	}
	//copy the client pointer
	client := app.client
	app.Unlock()
	in, _ := proto.Marshal(&api.Pingreq{Timestamp: time.Now().UnixNano()})
	var resp *api.Pingresp
	if e := client.Call(crpc.WithForceAddr(ctx, forceaddr), "/"+a+".status/ping", in, func(cctx *crpc.CallContext) error {
		out, e := cctx.Recv()
		if e != nil {
			return e
		}
		resp = &api.Pingresp{}
		if e := proto.Unmarshal(out, resp); e != nil {
			return ecode.ErrResp
		}
		return nil
	}); e != nil {
		return nil, e
	}
	return resp, nil
}

func (s *InternalSdk) PingByPrjoectName(ctx context.Context, pname, g, a string, forceaddr string) (*api.Pingresp, error) {
	fullname := pname + "-" + g + "." + a
	s.lker.RLock()
	if s.closing {
		s.lker.RUnlock()
		return nil, ecode.ErrServerClosing
	}
	app, ok := s.appsNameIndex[fullname]
	if !ok {
		s.lker.RUnlock()
		return nil, ecode.ErrAppNotExist
	}
	app.Lock()
	s.lker.RUnlock()
	if app.delstatus {
		app.Unlock()
		return nil, ecode.ErrAppNotExist
	}
	if app.client == nil {
		var e error
		if app.di == nil {
			app.di, e = createdi(app.summary)
			if e != nil {
				app.Unlock()
				return nil, e
			}
		}
		app.client, e = crpc.NewCrpcClient(nil, app.di, model.Project, model.Group, model.Name, app.summary.ProjectName, app.summary.Group, app.summary.App, nil)
		if e != nil {
			app.Unlock()
			return nil, e
		}
	}
	app.active = time.Now().UnixNano()
	if forceaddr != "" && app.summary.CrpcPort != 0 {
		if strings.Contains(forceaddr, ":") {
			//ipv6
			forceaddr = "[" + forceaddr + "]:" + strconv.FormatUint(uint64(app.summary.CrpcPort), 10)
		} else {
			//ipv4 or host
			forceaddr = forceaddr + ":" + strconv.FormatUint(uint64(app.summary.CrpcPort), 10)
		}
	}
	//copy the client pointer
	client := app.client
	app.Unlock()
	in, _ := proto.Marshal(&api.Pingreq{Timestamp: time.Now().UnixNano()})
	var resp *api.Pingresp
	if e := client.Call(crpc.WithForceAddr(ctx, forceaddr), "/"+a+".status/ping", in, func(cctx *crpc.CallContext) error {
		out, e := cctx.Recv()
		if e != nil {
			return e
		}
		resp = &api.Pingresp{}
		if e := proto.Unmarshal(out, resp); e != nil {
			return ecode.ErrResp
		}
		return nil
	}); e != nil {
		return nil, e
	}
	return resp, nil
}
