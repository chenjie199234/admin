package initinternal

import (
	"context"
	"encoding/base64"
	"errors"
	"sync"
	"time"

	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"

	"github.com/chenjie199234/Corelib/cerror"
	"github.com/chenjie199234/Corelib/crpc"
	"github.com/chenjie199234/Corelib/discover"
	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/metadata"
	"github.com/chenjie199234/Corelib/util/common"
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
}
type app struct {
	sync.RWMutex
	summary      *model.AppSummary
	notices      map[chan *struct{}]*struct{}
	di           discover.DI
	diactive     int64
	client       *crpc.CrpcClient
	clientactive int64
	delstatus    bool
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
	for _, v := range s.apps {
		app := v
		app.Lock()
		app.delstatus = true
		for notice := range app.notices {
			delete(app.notices, notice)
			close(notice)
		}
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
			//discover should always stop after client
			if v.client != nil {
				if now.UnixNano()-v.clientactive >= time.Minute.Nanoseconds() {
					//client timeout
					if now.UnixNano()-v.diactive >= time.Minute.Nanoseconds() {
						//client and di both timeout
						di := v.di
						client := v.client
						v.di = nil
						v.diactive = 0
						v.client = nil
						v.clientactive = 0
						go func() {
							client.Close(false)
							di.Stop()
						}()
					} else {
						//only client timeout
						client := v.client
						v.client = nil
						v.clientactive = 0
						go client.Close(false)
					}

				}
			} else if v.di != nil && now.UnixNano()-v.diactive >= time.Minute.Nanoseconds() {
				di := v.di
				v.di = nil
				v.diactive = 0
				go di.Stop()
			}
			v.Unlock()
		}
		s.lker.RUnlock()
	}
}
func (s *InternalSdk) mongoGetAllApp() error {
	filter := bson.M{"key": "", "index": 0}
	var cursor *mongo.Cursor
	primarylocalOPTS := options.Collection().SetReadPreference(readpref.Primary()).SetReadConcern(readconcern.Local())
	cursor, e := s.db.Database("app").Collection("config", primarylocalOPTS).Find(context.Background(), filter)
	if e != nil {
		log.Error(nil, "[InitWatch] get all app config failed", log.CError(e))
		return e
	}
	defer cursor.Close(context.Background())
	apps := make([]*model.AppSummary, 0, cursor.RemainingBatchLength())
	if e := cursor.All(context.Background(), &apps); e != nil {
		log.Error(nil, "[InitWatch] get all app config failed", log.CError(e))
		return e
	}
	for _, v := range apps {
		if e := s.decodeProxyPath(v); e != nil {
			return e
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
	watchfilter := mongo.Pipeline{bson.D{primitive.E{Key: "$match", Value: bson.M{"ns.db": "app", "ns.coll": "config"}}}}
	var stream *mongo.ChangeStream
	for {
		for stream == nil {
			//connect
			var e error
			opts := options.ChangeStream().SetFullDocument(options.UpdateLookup).SetStartAtOperationTime(s.start)
			if stream, e = s.db.Watch(context.Background(), watchfilter, opts); e != nil {
				log.Error(nil, "[InitWatch] get stream failed", log.CError(e))
				stream = nil
				time.Sleep(time.Millisecond * 100)
				continue
			}
		}
		for stream.Next(context.Background()) {
			s.start.T, s.start.I = stream.Current.Lookup("clusterTime").Timestamp()
			s.start.I++
			switch stream.Current.Lookup("operationType").StringValue() {
			case "drop":
				//drop collection
				log.Error(nil, "[InitWatch] all configs deleted")
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
					log.Error(nil, "[InitWatch] document format wrong",
						log.String("project_id", projectid),
						log.String("group", gname),
						log.String("app", aname),
						log.CError(e))
					continue
				}
				//decode proxy path
				if e := s.decodeProxyPath(summary); e != nil {
					log.Error(nil, "[InitWatch] db data broken",
						log.String("project_id", projectid),
						log.String("group", gname),
						log.String("app", aname),
						log.CError(e))
					continue
				}
				log.Debug(nil, "[InitWatch] updated",
					log.String("project_id", summary.ProjectID),
					log.String("group", summary.Group),
					log.String("app", summary.App),
					log.Any("keys", summary.Keys))
				s.lker.Lock()
				if exist, ok := s.apps[summary.ID.Hex()]; !ok {
					tmp := &app{
						summary: summary,
						notices: make(map[chan *struct{}]*struct{}, 10),
					}
					s.apps[summary.ID.Hex()] = tmp
					s.appsNameIndex[summary.ProjectName+"-"+summary.Group+"."+summary.App] = tmp
					s.appsIDIndex[summary.ProjectID+"-"+summary.Group+"."+summary.App] = tmp
				} else {
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
						exist.summary.KubernetesNs != summary.KubernetesNs ||
						exist.summary.KubernetesLS != summary.KubernetesLS ||
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
						log.Debug(nil, "[InitWatch] discover changed",
							log.String("project_id", projectid),
							log.String("group", gname),
							log.String("app", aname))
						//discover should always stop after client
						if exist.client != nil {
							di := exist.di
							client := exist.client
							exist.di = nil
							exist.diactive = 0
							exist.client = nil
							exist.clientactive = 0
							go func() {
								client.Close(false)
								di.Stop()
							}()
						} else if exist.di != nil {
							di := exist.di
							exist.di = nil
							exist.diactive = 0
							go di.Stop()
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
					//this is not the summary
					s.lker.Unlock()
					break
				}
				delete(s.apps, objid)
				delete(s.appsNameIndex, exist.summary.ProjectName+"-"+exist.summary.Group+"."+exist.summary.App)
				delete(s.appsIDIndex, exist.summary.ProjectID+"-"+exist.summary.Group+"."+exist.summary.App)
				exist.Lock()
				s.lker.Unlock()
				log.Debug(nil, "[InitWatch] deleted",
					log.String("project_id", exist.summary.ProjectID),
					log.String("group", exist.summary.Group),
					log.String("app", exist.summary.App))
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
			log.Error(nil, "[InitWatch] stream disconnected", log.CError(stream.Err()))
		}
		stream.Close(context.Background())
		stream = nil
	}
}
func (s *InternalSdk) decodeProxyPath(app *model.AppSummary) error {
	tmp := make(map[string]*model.ProxyPath)
	for path, info := range app.Paths {
		realpath, e := base64.StdEncoding.DecodeString(path)
		if e != nil {
			log.Error(nil, "[InitWatch] app's proxy path's base64 format wrong",
				log.String("project_id", app.ProjectID),
				log.String("group", app.Group),
				log.String("app", app.App),
				log.String("base64_path", path),
				log.CError(e))
			return e
		}
		tmp[common.Byte2str(realpath)] = info
	}
	app.Paths = tmp
	return nil
}

// if you don't need the notice,remember to call the cancel
// also remember to check the status of the notice,if notice is closed,means this app deleted
func (s *InternalSdk) GetNoticeByProjectID(pid, g, a string) (notice <-chan *struct{}, cancel func(), e error) {
	fullname := pid + "-" + g + "." + a
	s.lker.Lock()
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
// also remember to check the status of the notice,if notice is closed,means this app deleted
func (s *InternalSdk) GetNoticeByProjectName(pname, g, a string) (notice <-chan *struct{}, cancel func(), e error) {
	fullname := pname + "-" + g + "." + a
	s.lker.Lock()
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

func (s *InternalSdk) GetAppIPsByProjectID(ctx context.Context, pid, g, a string) ([]string, error) {
	fullname := pid + "-" + g + "." + a
	s.lker.RLock()
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
		switch app.summary.DiscoverMode {
		case "kubernetes":
			app.di, e = discover.NewKubernetesDiscover(app.summary.ProjectName, g, a, app.summary.KubernetesNs, app.summary.KubernetesLS, 9000, 10000, 8000)
		case "dns":
			interval := time.Duration(app.summary.DnsInterval) * time.Second
			app.di, e = discover.NewDNSDiscover(app.summary.ProjectName, g, a, app.summary.DnsHost, interval, 9000, 10000, 8000)
		case "static":
			app.di, e = discover.NewStaticDiscover(app.summary.ProjectName, g, a, app.summary.StaticAddrs, 9000, 10000, 8000)
		default:
			e = errors.New("unknown discover mode")
		}
		if e != nil {
			app.Unlock()
			return nil, e
		}
	}
	//copy the di pointer
	di := app.di
	app.Unlock()
	ch, cancel := di.GetNotice()
	defer cancel()
	select {
	case _, ok := <-ch:
		if !ok {
			return nil, cerror.ErrDiscoverStopped
		}
		tmp, _, _ := di.GetAddrs(discover.NotNeed)
		addrs := make([]string, 0, len(tmp))
		for addr := range tmp {
			addrs = append(addrs, addr)
		}
		return addrs, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
func (s *InternalSdk) GetAppIPsByProjectName(ctx context.Context, pname, g, a string) ([]string, error) {
	fullname := pname + "-" + g + "." + a
	s.lker.RLock()
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
		switch app.summary.DiscoverMode {
		case "kubernetes":
			app.di, e = discover.NewKubernetesDiscover(app.summary.ProjectName, g, a, app.summary.KubernetesNs, app.summary.KubernetesLS, 9000, 10000, 8000)
		case "dns":
			interval := time.Duration(app.summary.DnsInterval) * time.Second
			app.di, e = discover.NewDNSDiscover(app.summary.ProjectName, g, a, app.summary.DnsHost, interval, 9000, 10000, 8000)
		case "static":
			app.di, e = discover.NewStaticDiscover(app.summary.ProjectName, g, a, app.summary.StaticAddrs, 9000, 10000, 8000)
		default:
			e = errors.New("unknown discover mode")
		}
		if e != nil {
			app.Unlock()
			return nil, e
		}
	}
	//copy the di pointer
	di := app.di
	app.Unlock()
	ch, cancel := di.GetNotice()
	defer cancel()
	select {
	case _, ok := <-ch:
		if !ok {
			return nil, cerror.ErrDiscoverStopped
		}
		tmp, _, _ := di.GetAddrs(discover.NotNeed)
		addrs := make([]string, 0, len(tmp))
		for addr := range tmp {
			addrs = append(addrs, addr)
		}
		return addrs, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

type PermissionCheckHandler func(ctx context.Context, nodeid string, read, write, admin bool) error

func (s *InternalSdk) CallByPrjoectID(ctx context.Context, pid, g, a string, path string, reqdata []byte, pcheck PermissionCheckHandler) ([]byte, error) {
	fullname := pid + "-" + g + "." + a
	s.lker.RLock()
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
	pathinfo, ok := app.summary.Paths[path]
	if !ok {
		app.Unlock()
		return nil, ecode.ErrProxyPathNotExist
	}
	if pcheck != nil {
		if e := pcheck(ctx, pathinfo.PermissionNodeID, pathinfo.PermissionRead, pathinfo.PermissionWrite, pathinfo.PermissionAdmin); e != nil {
			app.Unlock()
			return nil, e
		}
	}
	if app.client == nil {
		var e error
		if app.di == nil {
			switch app.summary.DiscoverMode {
			case "kubernetes":
				app.di, e = discover.NewKubernetesDiscover(app.summary.ProjectName, g, a, app.summary.KubernetesNs, app.summary.KubernetesLS, 9000, 10000, 8000)
			case "dns":
				interval := time.Duration(app.summary.DnsInterval) * time.Second
				app.di, e = discover.NewDNSDiscover(app.summary.ProjectName, g, a, app.summary.DnsHost, interval, 9000, 10000, 8000)
			case "static":
				app.di, e = discover.NewStaticDiscover(app.summary.ProjectName, g, a, app.summary.StaticAddrs, 9000, 10000, 8000)
			default:
				e = errors.New("unknown discover mode")
			}
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
	now := time.Now().UnixNano()
	app.diactive = now
	app.clientactive = now
	//copy the client pointer
	client := app.client
	app.Unlock()
	return client.Call(ctx, path, reqdata, metadata.GetMetadata(ctx))
}

func (s *InternalSdk) CallByPrjoectName(ctx context.Context, pname, g, a string, path string, reqdata []byte, pcheck PermissionCheckHandler) ([]byte, error) {
	fullname := pname + "-" + g + "." + a
	s.lker.RLock()
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
	pathinfo, ok := app.summary.Paths[path]
	if !ok {
		app.Unlock()
		return nil, ecode.ErrProxyPathNotExist
	}
	if pcheck != nil {
		if e := pcheck(ctx, pathinfo.PermissionNodeID, pathinfo.PermissionRead, pathinfo.PermissionWrite, pathinfo.PermissionAdmin); e != nil {
			app.Unlock()
			return nil, e
		}
	}
	if app.client == nil {
		var e error
		if app.di == nil {
			switch app.summary.DiscoverMode {
			case "kubernetes":
				app.di, e = discover.NewKubernetesDiscover(app.summary.ProjectName, g, a, app.summary.KubernetesNs, app.summary.KubernetesLS, 9000, 10000, 8000)
			case "dns":
				interval := time.Duration(app.summary.DnsInterval) * time.Second
				app.di, e = discover.NewDNSDiscover(app.summary.ProjectName, g, a, app.summary.DnsHost, interval, 9000, 10000, 8000)
			case "static":
				app.di, e = discover.NewStaticDiscover(app.summary.ProjectName, g, a, app.summary.StaticAddrs, 9000, 10000, 8000)
			default:
				e = errors.New("unknown discover mode")
			}
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
	now := time.Now().UnixNano()
	app.diactive = now
	app.clientactive = now
	//copy the client pointer
	client := app.client
	app.Unlock()
	return client.Call(ctx, path, reqdata, metadata.GetMetadata(ctx))
}
