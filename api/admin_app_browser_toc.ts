// Code generated by protoc-gen-browser. DO NOT EDIT.
// version:
// 	protoc-gen-browser v0.0.97<br />
// 	protoc             v4.25.1<br />
// source: api/admin_app.proto<br />

export interface LogicError{
	code: number;
	msg: string;
}

export class DelAppReq{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	g_name: string = ''
	a_name: string = ''
	secret: string = ''
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.g_name){
			tmp["g_name"]=this.g_name
		}
		if(this.a_name){
			tmp["a_name"]=this.a_name
		}
		if(this.secret){
			tmp["secret"]=this.secret
		}
		return tmp
	}
}
export class DelAppResp{
	fromOBJ(obj:Object){
	}
}
export class DelKeyReq{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	g_name: string = ''
	a_name: string = ''
	key: string = ''//can't contain '.' in key
	secret: string = ''
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.g_name){
			tmp["g_name"]=this.g_name
		}
		if(this.a_name){
			tmp["a_name"]=this.a_name
		}
		if(this.key){
			tmp["key"]=this.key
		}
		if(this.secret){
			tmp["secret"]=this.secret
		}
		return tmp
	}
}
export class DelKeyResp{
	fromOBJ(obj:Object){
	}
}
export class DelProxyReq{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	g_name: string = ''
	a_name: string = ''
	path: string = ''
	secret: string = ''
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.g_name){
			tmp["g_name"]=this.g_name
		}
		if(this.a_name){
			tmp["a_name"]=this.a_name
		}
		if(this.path){
			tmp["path"]=this.path
		}
		if(this.secret){
			tmp["secret"]=this.secret
		}
		return tmp
	}
}
export class DelProxyResp{
	fromOBJ(obj:Object){
	}
}
export class GetAppReq{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	g_name: string = ''
	a_name: string = ''
	secret: string = ''
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.g_name){
			tmp["g_name"]=this.g_name
		}
		if(this.a_name){
			tmp["a_name"]=this.a_name
		}
		if(this.secret){
			tmp["secret"]=this.secret
		}
		return tmp
	}
}
export class GetAppResp{
	discover_mode: string = ''//can be one of "kubernetes" / "dns" / "static"
	kubernetes_namespace: string = ''//when discover_mode == "kubernetes"
	kubernetes_labelselector: string = ''//when discover_mode == "kubernetes"
	kubernetes_fieldselector: string = ''//when discover_mode == "kubernetes"
	dns_host: string = ''//when discover_mode == "dns"
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	dns_interval: number = 0//when discover_mode == "dns",unit seconds
	static_addrs: Array<string>|null = null//when discover_mode == "static"
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	crpc_port: number = 0
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	cgrpc_port: number = 0
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	web_port: number = 0
	keys: Map<string,KeyConfigInfo|null>|null = null
	paths: Map<string,ProxyPathInfo|null>|null = null
	fromOBJ(obj:Object){
		if(obj["discover_mode"]){
			this.discover_mode=obj["discover_mode"]
		}
		if(obj["kubernetes_namespace"]){
			this.kubernetes_namespace=obj["kubernetes_namespace"]
		}
		if(obj["kubernetes_labelselector"]){
			this.kubernetes_labelselector=obj["kubernetes_labelselector"]
		}
		if(obj["kubernetes_fieldselector"]){
			this.kubernetes_fieldselector=obj["kubernetes_fieldselector"]
		}
		if(obj["dns_host"]){
			this.dns_host=obj["dns_host"]
		}
		if(obj["dns_interval"]){
			this.dns_interval=obj["dns_interval"]
		}
		if(obj["static_addrs"] && obj["static_addrs"].length>0){
			this.static_addrs=obj["static_addrs"]
		}
		if(obj["crpc_port"]){
			this.crpc_port=obj["crpc_port"]
		}
		if(obj["cgrpc_port"]){
			this.cgrpc_port=obj["cgrpc_port"]
		}
		if(obj["web_port"]){
			this.web_port=obj["web_port"]
		}
		if(obj["keys"] && Object.keys(obj["keys"]).length>0){
			this.keys=new Map<string,KeyConfigInfo|null>()
			for(let key of Object.keys(obj["keys"])){
				if(obj["keys"][key]){
					let tmp = new KeyConfigInfo()
					tmp.fromOBJ(obj["keys"][key])
					this.keys.set(key,tmp)
				}else{
					this.keys.set(key,null)
				}
			}
		}
		if(obj["paths"] && Object.keys(obj["paths"]).length>0){
			this.paths=new Map<string,ProxyPathInfo|null>()
			for(let key of Object.keys(obj["paths"])){
				if(obj["paths"][key]){
					let tmp = new ProxyPathInfo()
					tmp.fromOBJ(obj["paths"][key])
					this.paths.set(key,tmp)
				}else{
					this.paths.set(key,null)
				}
			}
		}
	}
}
export class GetInstanceInfoReq{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	g_name: string = ''
	a_name: string = ''
	secret: string = ''
	addr: string = ''
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.g_name){
			tmp["g_name"]=this.g_name
		}
		if(this.a_name){
			tmp["a_name"]=this.a_name
		}
		if(this.secret){
			tmp["secret"]=this.secret
		}
		if(this.addr){
			tmp["addr"]=this.addr
		}
		return tmp
	}
}
export class GetInstanceInfoResp{
	info: InstanceInfo|null = null
	fromOBJ(obj:Object){
		if(obj["info"]){
			this.info=new InstanceInfo()
			this.info.fromOBJ(obj["info"])
		}
	}
}
export class GetInstancesReq{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	g_name: string = ''
	a_name: string = ''
	secret: string = ''
	with_info: boolean = false
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.g_name){
			tmp["g_name"]=this.g_name
		}
		if(this.a_name){
			tmp["a_name"]=this.a_name
		}
		if(this.secret){
			tmp["secret"]=this.secret
		}
		if(this.with_info){
			tmp["with_info"]=this.with_info
		}
		return tmp
	}
}
export class GetInstancesResp{
	instances: Map<string,InstanceInfo|null>|null = null//key addr,value info,if with_info is false,value is empty
	fromOBJ(obj:Object){
		if(obj["instances"] && Object.keys(obj["instances"]).length>0){
			this.instances=new Map<string,InstanceInfo|null>()
			for(let key of Object.keys(obj["instances"])){
				if(obj["instances"][key]){
					let tmp = new InstanceInfo()
					tmp.fromOBJ(obj["instances"][key])
					this.instances.set(key,tmp)
				}else{
					this.instances.set(key,null)
				}
			}
		}
	}
}
export class GetKeyConfigReq{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	g_name: string = ''
	a_name: string = ''
	key: string = ''//can't contain '.' in key
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	index: number = 0//0 means return current active config,config's index start from 1
	secret: string = ''
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.g_name){
			tmp["g_name"]=this.g_name
		}
		if(this.a_name){
			tmp["a_name"]=this.a_name
		}
		if(this.key){
			tmp["key"]=this.key
		}
		if(this.index){
			tmp["index"]=this.index
		}
		if(this.secret){
			tmp["secret"]=this.secret
		}
		return tmp
	}
}
export class GetKeyConfigResp{
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	cur_index: number = 0//current active config index,0 means not exist
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	max_index: number = 0//current max config index,0 means not exist
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	cur_version: number = 0//current active config version,config's version start from 1
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	this_index: number = 0//the config data below belong's to which index
	value: string = ''
	value_type: string = ''
	fromOBJ(obj:Object){
		if(obj["cur_index"]){
			this.cur_index=obj["cur_index"]
		}
		if(obj["max_index"]){
			this.max_index=obj["max_index"]
		}
		if(obj["cur_version"]){
			this.cur_version=obj["cur_version"]
		}
		if(obj["this_index"]){
			this.this_index=obj["this_index"]
		}
		if(obj["value"]){
			this.value=obj["value"]
		}
		if(obj["value_type"]){
			this.value_type=obj["value_type"]
		}
	}
}
export class InstanceInfo{
	name: string = ''
	//Warning!!!Type is uint64,be careful of sign(+)
	total_mem: bigint = BigInt(0)
	//Warning!!!Type is uint64,be careful of sign(+)
	cur_mem_usage: bigint = BigInt(0)
	cpu_num: number = 0
	cur_cpu_usage: number = 0
	fromOBJ(obj:Object){
		if(obj["name"]){
			this.name=obj["name"]
		}
		if(obj["total_mem"]){
			this.total_mem=BigInt(obj["total_mem"])
		}
		if(obj["cur_mem_usage"]){
			this.cur_mem_usage=BigInt(obj["cur_mem_usage"])
		}
		if(obj["cpu_num"]){
			this.cpu_num=obj["cpu_num"]
		}
		if(obj["cur_cpu_usage"]){
			this.cur_cpu_usage=obj["cur_cpu_usage"]
		}
	}
}
export class KeyConfigInfo{
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	cur_index: number = 0
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	max_index: number = 0
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	cur_version: number = 0
	cur_value: string = ''
	cur_value_type: string = ''
	fromOBJ(obj:Object){
		if(obj["cur_index"]){
			this.cur_index=obj["cur_index"]
		}
		if(obj["max_index"]){
			this.max_index=obj["max_index"]
		}
		if(obj["cur_version"]){
			this.cur_version=obj["cur_version"]
		}
		if(obj["cur_value"]){
			this.cur_value=obj["cur_value"]
		}
		if(obj["cur_value_type"]){
			this.cur_value_type=obj["cur_value_type"]
		}
	}
}
export class ProxyPathInfo{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	node_id: Array<number>|null = null
	read: boolean = false//need read permission on this node
	write: boolean = false//need write permission on this node
	admin: boolean = false//need admin permission on this node
	fromOBJ(obj:Object){
		if(obj["node_id"] && obj["node_id"].length>0){
			this.node_id=obj["node_id"]
		}
		if(obj["read"]){
			this.read=obj["read"]
		}
		if(obj["write"]){
			this.write=obj["write"]
		}
		if(obj["admin"]){
			this.admin=obj["admin"]
		}
	}
}
export class ProxyReq{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	g_name: string = ''
	a_name: string = ''
	path: string = ''
	data: string = ''
	force_addr: string = ''
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.g_name){
			tmp["g_name"]=this.g_name
		}
		if(this.a_name){
			tmp["a_name"]=this.a_name
		}
		if(this.path){
			tmp["path"]=this.path
		}
		if(this.data){
			tmp["data"]=this.data
		}
		if(this.force_addr){
			tmp["force_addr"]=this.force_addr
		}
		return tmp
	}
}
export class ProxyResp{
	data: string = ''
	fromOBJ(obj:Object){
		if(obj["data"]){
			this.data=obj["data"]
		}
	}
}
export class RollbackReq{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	g_name: string = ''
	a_name: string = ''
	key: string = ''//can't contain '.' in key
	secret: string = ''
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	index: number = 0
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.g_name){
			tmp["g_name"]=this.g_name
		}
		if(this.a_name){
			tmp["a_name"]=this.a_name
		}
		if(this.key){
			tmp["key"]=this.key
		}
		if(this.secret){
			tmp["secret"]=this.secret
		}
		if(this.index){
			tmp["index"]=this.index
		}
		return tmp
	}
}
export class RollbackResp{
	fromOBJ(obj:Object){
	}
}
export class SetAppReq{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	g_name: string = ''
	a_name: string = ''
	secret: string = ''
	discover_mode: string = ''
	kubernetes_namespace: string = ''//when discover_mode == "kubernetes"
	kubernetes_labelselector: string = ''//when discover_mode == "kubernetes"
	kubernetes_fieldselector: string = ''//when discover_mode == "kubernetes"
	dns_host: string = ''//when discover_mode == "dns"
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	dns_interval: number = 0//when discover_mode == "dns",unit seconds
	static_addrs: Array<string>|null = null//when discover_mode == "static"
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	crpc_port: number = 0
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	cgrpc_port: number = 0
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	web_port: number = 0
	new_app: boolean = false//true: create a new app. false: update the already exist app
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.g_name){
			tmp["g_name"]=this.g_name
		}
		if(this.a_name){
			tmp["a_name"]=this.a_name
		}
		if(this.secret){
			tmp["secret"]=this.secret
		}
		if(this.discover_mode){
			tmp["discover_mode"]=this.discover_mode
		}
		if(this.kubernetes_namespace){
			tmp["kubernetes_namespace"]=this.kubernetes_namespace
		}
		if(this.kubernetes_labelselector){
			tmp["kubernetes_labelselector"]=this.kubernetes_labelselector
		}
		if(this.kubernetes_fieldselector){
			tmp["kubernetes_fieldselector"]=this.kubernetes_fieldselector
		}
		if(this.dns_host){
			tmp["dns_host"]=this.dns_host
		}
		if(this.dns_interval){
			tmp["dns_interval"]=this.dns_interval
		}
		if(this.static_addrs && this.static_addrs.length>0){
			tmp["static_addrs"]=this.static_addrs
		}
		if(this.crpc_port){
			tmp["crpc_port"]=this.crpc_port
		}
		if(this.cgrpc_port){
			tmp["cgrpc_port"]=this.cgrpc_port
		}
		if(this.web_port){
			tmp["web_port"]=this.web_port
		}
		if(this.new_app){
			tmp["new_app"]=this.new_app
		}
		return tmp
	}
}
export class SetAppResp{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	node_id: Array<number>|null = null
	fromOBJ(obj:Object){
		if(obj["node_id"] && obj["node_id"].length>0){
			this.node_id=obj["node_id"]
		}
	}
}
export class SetKeyConfigReq{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	g_name: string = ''
	a_name: string = ''
	key: string = ''//can't contain '.' in key
	value: string = ''
	value_type: string = ''
	secret: string = ''
	new_key: boolean = false//true: create a new key config. false: update the already exist key config
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.g_name){
			tmp["g_name"]=this.g_name
		}
		if(this.a_name){
			tmp["a_name"]=this.a_name
		}
		if(this.key){
			tmp["key"]=this.key
		}
		if(this.value){
			tmp["value"]=this.value
		}
		if(this.value_type){
			tmp["value_type"]=this.value_type
		}
		if(this.secret){
			tmp["secret"]=this.secret
		}
		if(this.new_key){
			tmp["new_key"]=this.new_key
		}
		return tmp
	}
}
export class SetKeyConfigResp{
	fromOBJ(obj:Object){
	}
}
export class SetProxyReq{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	g_name: string = ''
	a_name: string = ''
	path: string = ''
	read: boolean = false//need read permission on this node
	write: boolean = false//need write permission on this node
	admin: boolean = false//need admin permission on this node
	secret: string = ''
	new_path: boolean = false//true: create a new proxy path config. false: update the already exist proxy path setting
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.g_name){
			tmp["g_name"]=this.g_name
		}
		if(this.a_name){
			tmp["a_name"]=this.a_name
		}
		if(this.path){
			tmp["path"]=this.path
		}
		if(this.read){
			tmp["read"]=this.read
		}
		if(this.write){
			tmp["write"]=this.write
		}
		if(this.admin){
			tmp["admin"]=this.admin
		}
		if(this.secret){
			tmp["secret"]=this.secret
		}
		if(this.new_path){
			tmp["new_path"]=this.new_path
		}
		return tmp
	}
}
export class SetProxyResp{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	node_id: Array<number>|null = null
	fromOBJ(obj:Object){
		if(obj["node_id"] && obj["node_id"].length>0){
			this.node_id=obj["node_id"]
		}
	}
}
export class UpdateAppSecretReq{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	g_name: string = ''
	a_name: string = ''
	old_secret: string = ''
	new_secret: string = ''
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.g_name){
			tmp["g_name"]=this.g_name
		}
		if(this.a_name){
			tmp["a_name"]=this.a_name
		}
		if(this.old_secret){
			tmp["old_secret"]=this.old_secret
		}
		if(this.new_secret){
			tmp["new_secret"]=this.new_secret
		}
		return tmp
	}
}
export class UpdateAppSecretResp{
	fromOBJ(obj:Object){
	}
}
export class WatchConfigReq{
	project_name: string = ''
	g_name: string = ''
	a_name: string = ''
	//map's key is config's keyname,map's value is config's cur version
	//if cur version == 0 means return current active config
	//if all cur version is the newest,the request will block until a new version come
	//if some keys' version is the newest,and some keys' version is old,then the keys with old version will return newest version and datas,the newest's keys will only return version
	//Warning!!!map's value's type is uint32,be careful of sign(+) and overflow
	keys: Map<string,number>|null = null//can't contain '.' in key
	toJSON(){
		let tmp = {}
		if(this.project_name){
			tmp["project_name"]=this.project_name
		}
		if(this.g_name){
			tmp["g_name"]=this.g_name
		}
		if(this.a_name){
			tmp["a_name"]=this.a_name
		}
		if(this.keys && this.keys.size>0){
			tmp["keys"]={}
			for(let [k,v] of this.keys){
				tmp["keys"][k]=v
			}
		}
		return tmp
	}
}
export class WatchConfigResp{
	datas: Map<string,WatchData|null>|null = null
	fromOBJ(obj:Object){
		if(obj["datas"] && Object.keys(obj["datas"]).length>0){
			this.datas=new Map<string,WatchData|null>()
			for(let key of Object.keys(obj["datas"])){
				if(obj["datas"][key]){
					let tmp = new WatchData()
					tmp.fromOBJ(obj["datas"][key])
					this.datas.set(key,tmp)
				}else{
					this.datas.set(key,null)
				}
			}
		}
	}
}
export class WatchData{
	key: string = ''
	value: string = ''
	value_type: string = ''
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	version: number = 0
	fromOBJ(obj:Object){
		if(obj["key"]){
			this.key=obj["key"]
		}
		if(obj["value"]){
			this.value=obj["value"]
		}
		if(obj["value_type"]){
			this.value_type=obj["value_type"]
		}
		if(obj["version"]){
			this.version=obj["version"]
		}
	}
}
export class WatchDiscoverReq{
	project_name: string = ''
	g_name: string = ''
	a_name: string = ''
	cur_discover_mode: string = ''
	cur_dns_host: string = ''//when discover_mode == "dns"
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	cur_dns_interval: number = 0//when cur_discover_mode == "dns",unit seconds
	cur_static_addrs: Array<string>|null = null//when cur_discover_mode == "static"
	cur_kubernetes_namespace: string = ''//when cur_discover_mode == "kubernetes"
	cur_kubernetes_labelselector: string = ''//when cur_discover_mode == "kubernetes"
	cur_kubernetes_fieldselector: string = ''//when cur_discover_mode == "kubernetes"
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	cur_crpc_port: number = 0
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	cur_cgrpc_port: number = 0
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	cur_web_port: number = 0
	toJSON(){
		let tmp = {}
		if(this.project_name){
			tmp["project_name"]=this.project_name
		}
		if(this.g_name){
			tmp["g_name"]=this.g_name
		}
		if(this.a_name){
			tmp["a_name"]=this.a_name
		}
		if(this.cur_discover_mode){
			tmp["cur_discover_mode"]=this.cur_discover_mode
		}
		if(this.cur_dns_host){
			tmp["cur_dns_host"]=this.cur_dns_host
		}
		if(this.cur_dns_interval){
			tmp["cur_dns_interval"]=this.cur_dns_interval
		}
		if(this.cur_static_addrs && this.cur_static_addrs.length>0){
			tmp["cur_static_addrs"]=this.cur_static_addrs
		}
		if(this.cur_kubernetes_namespace){
			tmp["cur_kubernetes_namespace"]=this.cur_kubernetes_namespace
		}
		if(this.cur_kubernetes_labelselector){
			tmp["cur_kubernetes_labelselector"]=this.cur_kubernetes_labelselector
		}
		if(this.cur_kubernetes_fieldselector){
			tmp["cur_kubernetes_fieldselector"]=this.cur_kubernetes_fieldselector
		}
		if(this.cur_crpc_port){
			tmp["cur_crpc_port"]=this.cur_crpc_port
		}
		if(this.cur_cgrpc_port){
			tmp["cur_cgrpc_port"]=this.cur_cgrpc_port
		}
		if(this.cur_web_port){
			tmp["cur_web_port"]=this.cur_web_port
		}
		return tmp
	}
}
export class WatchDiscoverResp{
	discover_mode: string = ''
	dns_host: string = ''//when discover_mode == "dns"
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	dns_interval: number = 0//when discover_mode == "dns"
	static_addrs: Array<string>|null = null//when discover_mode == "static"
	kubernetes_namespace: string = ''//when discover_mode == "kubernetes"
	kubernetes_labelselector: string = ''//when discover_mode == "kubernetes"
	kubernetes_fieldselector: string = ''//when discover_mode == "kubernetes"
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	crpc_port: number = 0
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	cgrpc_port: number = 0
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	web_port: number = 0
	fromOBJ(obj:Object){
		if(obj["discover_mode"]){
			this.discover_mode=obj["discover_mode"]
		}
		if(obj["dns_host"]){
			this.dns_host=obj["dns_host"]
		}
		if(obj["dns_interval"]){
			this.dns_interval=obj["dns_interval"]
		}
		if(obj["static_addrs"] && obj["static_addrs"].length>0){
			this.static_addrs=obj["static_addrs"]
		}
		if(obj["kubernetes_namespace"]){
			this.kubernetes_namespace=obj["kubernetes_namespace"]
		}
		if(obj["kubernetes_labelselector"]){
			this.kubernetes_labelselector=obj["kubernetes_labelselector"]
		}
		if(obj["kubernetes_fieldselector"]){
			this.kubernetes_fieldselector=obj["kubernetes_fieldselector"]
		}
		if(obj["crpc_port"]){
			this.crpc_port=obj["crpc_port"]
		}
		if(obj["cgrpc_port"]){
			this.cgrpc_port=obj["cgrpc_port"]
		}
		if(obj["web_port"]){
			this.web_port=obj["web_port"]
		}
	}
}
//timeout's unit is millisecond,it will be used when > 0
function call(timeout: number,url: string,opts: Object,error: (arg: LogicError)=>void,success: (arg: Object)=>void){
	let tid: number|null = null
	if(timeout>0){
		const c = new AbortController()
		opts["signal"] = c.signal
		tid = setTimeout(()=>{c.abort()},timeout)
	}
	let ok=false
	fetch(url,opts)
	.then(r=>{
		ok=r.ok
		if(r.ok){
			return r.json()
		}
		return r.text()
	})
	.then(d=>{
		if(!ok){
			throw d
		}
		success(d)
	})
	.catch(e=>{
		if(e instanceof Error){
			error({code:-1,msg:e.message})
		}else if(e.length>0 && e[0]=='{' && e[e.length-1]=='}'){
			error(JSON.parse(e))
		}else{
			error({code:-1,msg:e})
		}
	})
	.finally(()=>{
		if(tid){
			clearTimeout(tid)
		}
	})
}
const _WebPathAppGetApp: string ="/admin.app/get_app";
const _WebPathAppSetApp: string ="/admin.app/set_app";
const _WebPathAppDelApp: string ="/admin.app/del_app";
const _WebPathAppUpdateAppSecret: string ="/admin.app/update_app_secret";
const _WebPathAppDelKey: string ="/admin.app/del_key";
const _WebPathAppGetKeyConfig: string ="/admin.app/get_key_config";
const _WebPathAppSetKeyConfig: string ="/admin.app/set_key_config";
const _WebPathAppRollback: string ="/admin.app/rollback";
const _WebPathAppWatchConfig: string ="/admin.app/watch_config";
const _WebPathAppWatchDiscover: string ="/admin.app/watch_discover";
const _WebPathAppGetInstances: string ="/admin.app/get_instances";
const _WebPathAppGetInstanceInfo: string ="/admin.app/get_instance_info";
const _WebPathAppSetProxy: string ="/admin.app/set_proxy";
const _WebPathAppDelProxy: string ="/admin.app/del_proxy";
const _WebPathAppProxy: string ="/admin.app/proxy";
//ToC means this is for users
export class AppBrowserClientToC {
	constructor(host: string){
		if(!host || host.length==0){
			throw "AppBrowserClientToC's host missing"
		}
		this.host=host
	}
	//timeout's unit is millisecond,it will be used when > 0
	get_app(header: Object,req: GetAppReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: GetAppResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathAppGetApp,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new GetAppResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	set_app(header: Object,req: SetAppReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: SetAppResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathAppSetApp,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new SetAppResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	del_app(header: Object,req: DelAppReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: DelAppResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathAppDelApp,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new DelAppResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	update_app_secret(header: Object,req: UpdateAppSecretReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: UpdateAppSecretResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathAppUpdateAppSecret,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new UpdateAppSecretResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	del_key(header: Object,req: DelKeyReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: DelKeyResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathAppDelKey,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new DelKeyResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	get_key_config(header: Object,req: GetKeyConfigReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: GetKeyConfigResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathAppGetKeyConfig,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new GetKeyConfigResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	set_key_config(header: Object,req: SetKeyConfigReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: SetKeyConfigResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathAppSetKeyConfig,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new SetKeyConfigResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	rollback(header: Object,req: RollbackReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: RollbackResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathAppRollback,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new RollbackResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	watch_config(header: Object,req: WatchConfigReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: WatchConfigResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathAppWatchConfig,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new WatchConfigResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	watch_discover(header: Object,req: WatchDiscoverReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: WatchDiscoverResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathAppWatchDiscover,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new WatchDiscoverResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	get_instances(header: Object,req: GetInstancesReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: GetInstancesResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathAppGetInstances,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new GetInstancesResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	get_instance_info(header: Object,req: GetInstanceInfoReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: GetInstanceInfoResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathAppGetInstanceInfo,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new GetInstanceInfoResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	set_proxy(header: Object,req: SetProxyReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: SetProxyResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathAppSetProxy,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new SetProxyResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	del_proxy(header: Object,req: DelProxyReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: DelProxyResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathAppDelProxy,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new DelProxyResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	proxy(header: Object,req: ProxyReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: ProxyResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathAppProxy,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new ProxyResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	private host: string
}
