<script setup lang="ts">
import { ref,computed } from 'vue'
import * as appAPI from './api/admin_app_browser_toc'
import * as permissionAPI from './api/admin_permission_browser_toc'
import * as state from './state'
import * as client from './client'

const all=computed(()=>{
	let tmp: {[k:string]: {[k:string]:permissionAPI.NodeInfo}} = {}
	if(!state.page.node!.children){
		return tmp
	}
	for(let n of state.page.node!.children){
		if(!n){
			continue
		}
		let pieces:string[]=n.node_name.split(".")
		if(pieces.length!=2){
			state.set_alert("error",-1,"app's permission node's nodename format wrong,should be 'group.app',nodeid:"+n.node_id!.toString())
			continue
		}
		if(tmp[pieces[0]]){
			tmp[pieces[0]][pieces[1]] = n
		}else{
			tmp[pieces[0]]={}
			tmp[pieces[0]][pieces[1]] = n
		}
	}
	return tmp
})

const curg=ref<string>("")
const cura=ref<string>("")
const secret=ref<string>("")
const t_secret=ref<boolean>(false)
function selfapp():boolean{
	let nodeid=all.value[curg.value][cura.value].node_id
	return nodeid![1]==1&&nodeid![2]==2&&nodeid![3]==1
}
function canwrite():boolean{
	return all.value[curg.value][cura.value].canwrite||all.value[curg.value][cura.value].admin
}
function mustadmin():boolean{
	return all.value[curg.value][cura.value].admin
}

const config_proxy_instance=ref<string>("")

const discovermode=ref<string>("")
const kubernetesns=ref<string>("")
const kubernetesls=ref<string>("")
const kubernetesfs=ref<string>("")
const dnshost=ref<string>("")
const dnsinterval=ref<number>(0)
const staticaddrs=ref<string[]>([])
const crpc_port=ref<number>(0)
const cgrpc_port=ref<number>(0)
const web_port=ref<number>(0)

const keys=ref<Map<string,appAPI.KeyConfigInfo>>(new Map())
const t_keys_hover=ref<boolean>(false)
const keyhover=ref<string>("")

const proxys=ref<Map<string,appAPI.ProxyPathInfo>>(new Map())
const t_proxys_hover=ref<boolean>(false)
const proxyhover=ref<string>("")

const instances=ref<Map<string,appAPI.InstanceInfo|null>>(new Map())
const t_instances_hover=ref<boolean>(false)

const get_app_status=ref<boolean>(false)

function get_app(){
	if(curg.value==""||cura.value==""){
		keys.value=new Map()
		proxys.value=new Map()
		config_proxy_instance.value=""
		discovermode.value=""
		kubernetesns.value=""
		kubernetesls.value=""
		kubernetesfs.value=""
		dnshost.value=""
		dnsinterval.value=0
		staticaddrs.value=[]
		state.set_alert("error",-2,"Group and App must be selected!")
		return
	}
	if(!all.value[curg.value][cura.value].node_id||all.value[curg.value][cura.value].node_id!.length!=4){
		keys.value=new Map()
		proxys.value=new Map()
		config_proxy_instance.value=""
		discovermode.value=""
		kubernetesns.value=""
		kubernetesls.value=""
		kubernetesfs.value=""
		dnshost.value=""
		dnsinterval.value=0
		staticaddrs.value=[]
		state.set_alert("error",-2,"Missing node_id on Group:"+curg.value+" App:"+cura.value)
		return
	}
	if(!state.set_load()){
		return
	}
	let req=new appAPI.GetAppReq()
	req.project_id=state.project.info!.project_id
	req.g_name=curg.value
	req.a_name=cura.value
	req.secret=secret.value
	client.appClient.get_app({"Token":state.user.token},req,client.timeout,(e: appAPI.LogicError)=>{
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
	},(resp: appAPI.GetAppResp)=>{
		if(resp.keys){
			keys.value=new Map()
			let tmp = [...resp.keys.entries()].sort()
			for(let i=0;i<tmp.length;i++){
				if(tmp[i][1]){
					keys.value.set(tmp[i][0],tmp[i][1]!)
				}
			}
		}else{
			keys.value = new Map()
		}
		if(resp.paths){
			proxys.value = new Map()
			let tmp = [...resp.paths.entries()].sort()
			for(let i=0;i<tmp.length;i++){
				if(tmp[i][1]){
					proxys.value.set(tmp[i][0],tmp[i][1]!)
				}
			}
		}else{
			proxys.value = new Map()
		}
		discovermode.value=resp.discover_mode
		kubernetesns.value=resp.kubernetes_namespace
		kubernetesls.value=resp.kubernetes_labelselector
		kubernetesfs.value=resp.kubernetes_fieldselector
		dnshost.value=resp.dns_host
		dnsinterval.value=resp.dns_interval
		if(resp.static_addrs){
			staticaddrs.value=resp.static_addrs
		}else{
			staticaddrs.value=[]
		}
		crpc_port.value=resp.crpc_port
		cgrpc_port.value=resp.cgrpc_port
		web_port.value=resp.web_port
		get_app_status.value=true
		state.clear_load()
	})
}
function get_instances(withinfo: boolean){
	if(!state.set_load()){
		return
	}
	let req=new appAPI.GetInstancesReq()
	req.project_id=state.project.info!.project_id
	req.g_name=curg.value
	req.a_name=cura.value
	req.secret=secret.value
	req.with_info=withinfo
	client.appClient.get_instances({"Token":state.user.token},req,client.timeout,(e: appAPI.LogicError)=>{
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
	},(resp: appAPI.GetInstancesResp)=>{
		if(resp.instances){
			instances.value=new Map()
			let tmp = [...resp.instances.entries()].sort()
			for(let i=0;i<tmp.length;i++){
				if(tmp[i][1]){
					instances.value.set(tmp[i][0],tmp[i][1]!)
				}else{
					instances.value.set(tmp[i][0],null)
				}
			}
		}else{
			instances.value=new Map()
		}
		state.clear_load()
	})
}
function get_instance(addr: string){
	if(!state.set_load()){
		return
	}
	let req=new appAPI.GetInstanceInfoReq()
	req.project_id=state.project.info!.project_id
	req.g_name=curg.value
	req.a_name=cura.value
	req.secret=secret.value
	req.addr=addr
	client.appClient.get_instance_info({"Token":state.user.token},req,client.timeout,(e: appAPI.LogicError)=>{
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
	},(resp: appAPI.GetInstanceInfoResp)=>{
		if(resp.info){
			if(!instances.value){
				instances.value=new Map()
			}
			instances.value.set(addr,resp.info)
		}
		state.clear_load()
	})
}

const ing=ref<boolean>(false)
const optype=ref<string>("")

//add app
const new_g=ref<string>("")
const new_a=ref<string>("")
const new_secret=ref<string>("")
const new_discovermode=ref<string>("")
const new_kubernetesns=ref<string>("")
const new_kubernetesls=ref<string>("")
const new_kubernetesfs=ref<string>("")
const new_dnshost=ref<string>("")
const new_dnsinterval=ref<number>(0)
const new_staticaddrs=ref<string[]>([])
const new_staticaddr=ref<string>("")
const new_crpc_port=ref<number>(0)
const new_cgrpc_port=ref<number>(0)
const new_web_port=ref<number>(0)
function reset_new_app(){
	new_kubernetesns.value=''
	new_kubernetesls.value=''
	new_kubernetesfs.value=''
	new_dnshost.value=''
	new_dnsinterval.value=0
	new_staticaddrs.value=[]
	new_staticaddr.value=''
	new_crpc_port.value=0
	new_cgrpc_port.value=0
	new_web_port.value=0
}
function add_app_able() :boolean{
	if(new_g.value==''){
		return false
	}
	if(new_a.value==''){
		return false
	}
	if(new_discovermode.value==''){
		return false
	}
	if(new_discovermode.value=='kubernetes'&&(new_kubernetesns.value==''||(new_kubernetesls.value==''&&new_kubernetesfs.value==''))){
		return false
	}
	if(new_discovermode.value=='dns'&&new_dnshost.value==''){
		return false
	}
	if(new_discovermode.value=='static'&&new_staticaddrs.value.length==0){
		return false
	}
	if(new_crpc_port.value==0){
		return false
	}
	if(new_cgrpc_port.value==0){
		return false
	}
	if(new_web_port.value==0){
		return false
	}
	return true
}

//update app
const update_new_secret=ref<string>("")
const update_new_discovermode=ref<string>("")
const update_new_kubernetesns=ref<string>("")
const update_new_kubernetesls=ref<string>("")
const update_new_kubernetesfs=ref<string>("")
const update_new_dnshost=ref<string>("")
const update_new_dnsinterval=ref<number>(0)
const update_new_staticaddrs=ref<string[]>([])
const update_new_staticaddr=ref<string>("")
const update_new_crpc_port=ref<number>(0)
const update_new_cgrpc_port=ref<number>(0)
const update_new_web_port=ref<number>(0)
function update_secret_able():boolean{
	return update_new_secret.value != secret.value
}
function reset_update_discover(dname:string){
	if(dname!=''){
		update_new_discovermode.value=dname
	}
	update_new_kubernetesns.value=''
	update_new_kubernetesls.value=''
	update_new_kubernetesfs.value=''
	update_new_dnshost.value=''
	update_new_dnsinterval.value=0
	update_new_staticaddrs.value=[]
	update_new_staticaddr.value=''
	update_new_crpc_port.value=crpc_port.value
	update_new_cgrpc_port.value=cgrpc_port.value
	update_new_web_port.value=web_port.value
	if(update_new_discovermode.value=="kubernetes"){
		update_new_kubernetesns.value=kubernetesns.value
		update_new_kubernetesls.value=kubernetesls.value
		update_new_kubernetesfs.value=kubernetesfs.value
	}else if(update_new_discovermode.value=="dns"){
		update_new_dnshost.value=dnshost.value
		update_new_dnsinterval.value=dnsinterval.value
	}else if(update_new_discovermode.value=="static"){
		update_new_staticaddrs.value=[...staticaddrs.value]
	}
}
function update_discover_able():boolean{
	if(update_new_discovermode.value==''){
		return false
	}
	if(update_new_crpc_port.value==0){
		return false
	}
	if(update_new_cgrpc_port.value==0){
		return false
	}
	if(update_new_web_port.value==0){
		return false
	}
	let sameport = update_new_crpc_port.value==crpc_port.value&&update_new_cgrpc_port.value==cgrpc_port.value&&update_new_web_port.value==web_port.value
	if(update_new_discovermode.value=='kubernetes'){
		if(update_new_kubernetesns.value==''||(update_new_kubernetesls.value==''&&update_new_kubernetesfs.value=='')){
			return false
		}
		let samekubernetes = update_new_kubernetesns.value==kubernetesns.value&&
			update_new_kubernetesls.value==kubernetesls.value&&
			update_new_kubernetesfs.value==kubernetesfs.value
		return !sameport || !samekubernetes
	}
	if(update_new_discovermode.value=='dns'){
		if(update_new_dnshost.value==''){
			return false
		}
		let samedns = update_new_dnshost.value==dnshost.value&&update_new_dnsinterval.value==dnsinterval.value
		return !sameport || !samedns
	}
	if(update_new_discovermode.value=='static'){
		if(update_new_staticaddrs.value.length==0){
			return false
		}
		let sameaddr = update_new_staticaddrs.value.length==staticaddrs.value.length
		if(sameaddr){
			sameaddr = staticaddrs.value.every(function(v,i){return v==update_new_staticaddrs.value[i]})
		}
		console.log(sameport)
		console.log(sameaddr)
		return !sameport || !sameaddr
	}
	return true
}

//add key config
//const config_value_types=ref<string[]>(["json","raw","yaml","toml"])
const config_value_type=ref<string>("json")
const config_value=ref<string>("{\n}")
const config_key=ref<string>("")
function reset_add_key(){
	config_key.value=""
	config_value.value="{\n}"
	config_value_type.value="json"
}
function add_key_able():boolean{
	if(config_key.value==''){
		return false
	}
	if(keys.value.has(config_key.value)){
		return false
	}
	if(!is_json_obj(config_value.value)){
		return false
	}
	return true
}

const cur_key=ref<string>("")

//rollback
const rollback_key_index=ref<number>(0)
const rollback_key_value=ref<string>("")
const rollback_key_value_type=ref<string>("")

//edit
const edit_key_value=ref<string>("")
const edit_key_value_type=ref<string>("")
function edit_commit_able():boolean{
	if(edit_key_value_type.value=="json"){
		if(!is_json_obj(edit_key_value.value)){
			return false
		}
		return JSON.stringify(JSON.parse(edit_key_value.value),null,4)!=JSON.stringify(JSON.parse(keys.value.get(cur_key.value)!.cur_value),null,4)
	}
	return edit_key_value.value!=keys.value.get(cur_key.value)!.cur_value
}

//add proxy
const new_proxy_path=ref<string>("")
const new_proxy_permission_read=ref<boolean>(false)
const new_proxy_permission_write=ref<boolean>(false)
const new_proxy_permission_admin=ref<boolean>(false)
function reset_add_proxy(){
	new_proxy_path.value=''
	new_proxy_permission_read.value=false
	new_proxy_permission_write.value=false
	new_proxy_permission_admin.value=false
}
function add_proxy_able():boolean{
	return new_proxy_path.value!=""&&!proxys.value.has(new_proxy_path.value)
}

const cur_proxy=ref<string>("")

//update proxy
const update_proxy_permission_read=ref<boolean>(false)
const update_proxy_permission_write=ref<boolean>(false)
const update_proxy_permission_admin=ref<boolean>(false)

function update_proxy_able():boolean{
	if(proxys.value.get(cur_proxy.value)!.read!=update_proxy_permission_read.value){
		return true
	}
	if(proxys.value.get(cur_proxy.value)!.write!=update_proxy_permission_write.value){
		return true
	}
	if(proxys.value.get(cur_proxy.value)!.admin!=update_proxy_permission_admin.value){
		return true
	}
	return false
}

//proxy
const reqdata=ref<string>("")
const respdata=ref<string>("")
const forceaddr=ref<string>("")

function app_op(){
	if(!state.set_load()){
		return
	}
	switch(optype.value){
		case 'del_app':{
			let req=new appAPI.DelAppReq()
			req.project_id=state.project.info!.project_id
			req.g_name=curg.value
			req.a_name=cura.value
			req.secret=secret.value
			client.appClient.del_app({"Token":state.user.token},req,client.timeout,(e: appAPI.LogicError)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(_resp: appAPI.DelAppResp)=>{
				let node = all.value[curg.value][cura.value]
				if(state.page.node!.children){
					for(let i=0;i<state.page.node!.children!.length;i++){
						if(node == state.page.node!.children![i]){
							state.page.node!.children!.splice(i,1)
							break
						}
					}
				}
				curg.value=""
				cura.value=""
				secret.value=""
				discovermode.value=""
				keys.value=new Map()
				proxys.value=new Map()
				get_app_status.value=false
				config_proxy_instance.value=""
				ing.value=false
				state.clear_load()
			})
			break
		}
		case 'add_app':{
			let req=new appAPI.SetAppReq()
			req.project_id=state.project.info!.project_id
			req.g_name=new_g.value
			req.a_name=new_a.value
			req.secret=new_secret.value
			req.discover_mode=new_discovermode.value
			req.kubernetes_namespace=new_kubernetesns.value
			req.kubernetes_labelselector=new_kubernetesls.value
			req.kubernetes_fieldselector=new_kubernetesfs.value
			req.dns_host=new_dnshost.value
			req.dns_interval=new_dnsinterval.value
			req.static_addrs=new_staticaddrs.value
			req.crpc_port=new_crpc_port.value
			req.cgrpc_port=new_cgrpc_port.value
			req.web_port=new_web_port.value
			req.new_app=true
			client.appClient.set_app({"Token":state.user.token},req,client.timeout,(e: appAPI.LogicError)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: appAPI.SetAppResp)=>{
				if(all.value[new_g.value] && all.value[new_g.value][new_a.value]){
					return
				}
				let tmp = new permissionAPI.NodeInfo()
				tmp.node_id=resp.node_id
				tmp.node_name=new_g.value+"."+new_a.value
				tmp.node_data=""
				tmp.canread=true
				tmp.canwrite=true
				tmp.admin=true
				tmp.children=[]
				if(!state.page.node!.children){
					state.page.node!.children=[]
				}
				state.page.node!.children!.push(tmp)
				new_g.value=""
				new_a.value=""
				new_secret.value=""
				new_discovermode.value=""
				ing.value=false
				state.clear_load()
			})
			break
		}
		case 'update_discover':{
			let req=new appAPI.SetAppReq()
			req.project_id=state.project.info!.project_id
			req.g_name=curg.value
			req.a_name=cura.value
			req.secret=secret.value
			req.discover_mode=update_new_discovermode.value
			req.kubernetes_namespace=update_new_kubernetesns.value
			req.kubernetes_labelselector=update_new_kubernetesls.value
			req.kubernetes_fieldselector=update_new_kubernetesfs.value
			req.dns_host=update_new_dnshost.value
			req.dns_interval=update_new_dnsinterval.value
			req.static_addrs=update_new_staticaddrs.value
			req.crpc_port=update_new_crpc_port.value
			req.cgrpc_port=update_new_cgrpc_port.value
			req.web_port=update_new_web_port.value
			req.new_app=false
			client.appClient.set_app({"Token":state.user.token},req,client.timeout,(e: appAPI.LogicError)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(_: appAPI.SetAppResp)=>{
				discovermode.value=update_new_discovermode.value
				kubernetesns.value=update_new_kubernetesns.value
				kubernetesls.value=update_new_kubernetesls.value
				dnshost.value=update_new_dnshost.value
				dnsinterval.value=update_new_dnsinterval.value
				staticaddrs.value=update_new_staticaddrs.value
				update_new_discovermode.value=""
				crpc_port.value=update_new_crpc_port.value,
				cgrpc_port.value=update_new_cgrpc_port.value,
				web_port.value=update_new_web_port.value,
				ing.value=false
				state.clear_load()
			})
			break
		}
		case 'update_secret':{
			let req=new appAPI.UpdateAppSecretReq()
			req.project_id=state.project.info!.project_id
			req.g_name=curg.value
			req.a_name=cura.value
			req.old_secret=secret.value
			req.new_secret=update_new_secret.value
			client.appClient.update_app_secret({"Token":state.user.token},req,client.timeout,(e: appAPI.LogicError)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(_resp: appAPI.UpdateAppSecretResp)=>{
				secret.value=update_new_secret.value
				update_new_secret.value=""
				ing.value=false
				state.clear_load()
			})
			break
		}
		case 'get_rollback_key':{
			let req=new appAPI.GetKeyConfigReq()
			req.project_id=state.project.info!.project_id
			req.g_name=curg.value
			req.a_name=cura.value
			req.secret=secret.value
			req.key=cur_key.value
			req.index=rollback_key_index.value
			client.appClient.get_key_config({"Token":state.user.token},req,client.timeout,(e: appAPI.LogicError)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: appAPI.GetKeyConfigResp)=>{
				if(resp.value){
					rollback_key_value.value=resp.value
					rollback_key_value_type.value=resp.value_type
				}else{
					rollback_key_value.value="{}"
					rollback_key_value_type.value="json"
				}
				state.clear_load()
			})
			break
		}
		case 'add_key':{
			let req=new appAPI.SetKeyConfigReq()
			req.project_id=state.project.info!.project_id
			req.g_name=curg.value
			req.a_name=cura.value
			req.secret=secret.value
			req.key=config_key.value
			req.value=config_value.value
			req.value_type=config_value_type.value
			req.new_key=true
			client.appClient.set_key_config({"Token":state.user.token},req,client.timeout,(e: appAPI.LogicError)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(_resp: appAPI.SetKeyConfigResp)=>{
				let tmp = new appAPI.KeyConfigInfo()
				tmp.cur_index=1
				tmp.max_index=1
				tmp.cur_version=1
				tmp.cur_value=config_value.value
				tmp.cur_value_type=config_value_type.value
				keys.value.set(config_key.value,tmp)
				config_key.value = ''
				config_value.value = '{\n}'
				config_value_type.value = 'json'
				ing.value=false
				state.clear_load()
			})
			break
		}
		case 'update_key':{
			let req=new appAPI.SetKeyConfigReq()
			req.project_id=state.project.info!.project_id
			req.g_name=curg.value
			req.a_name=cura.value
			req.secret=secret.value
			req.key=cur_key.value
			req.value=edit_key_value.value
			req.value_type=edit_key_value_type.value
			req.new_key=false
			client.appClient.set_key_config({"Token":state.user.token},req,client.timeout,(e: appAPI.LogicError)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(_resp: appAPI.SetKeyConfigResp)=>{
				edit_key_value_type.value=''
				ing.value=false
				state.clear_load()
				get_app()
			})
			break
		}
		case 'rollback_key':{
			let req=new appAPI.RollbackReq()
			req.project_id=state.project.info!.project_id
			req.g_name=curg.value
			req.a_name=cura.value
			req.secret=secret.value
			req.key=cur_key.value
			req.index=rollback_key_index.value
			client.appClient.rollback({"Token":state.user.token},req,client.timeout,(e: appAPI.LogicError)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(_resp: appAPI.RollbackResp)=>{
				rollback_key_index.value=0
				ing.value=false
				state.clear_load()
				get_app()
			})
			break
		}
		case 'del_key':{
			let req=new appAPI.DelKeyReq()
			req.project_id=state.project.info!.project_id
			req.g_name=curg.value
			req.a_name=cura.value
			req.secret=secret.value
			req.key=cur_key.value
			client.appClient.del_key({"Token":state.user.token},req,client.timeout,(e: appAPI.LogicError)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(_resp: appAPI.DelKeyResp)=>{
				keys.value.delete(cur_key.value)
				cur_key.value=""
				ing.value=false
				state.clear_load()
			})
			break
		}
		case 'add_proxy':{
			let req=new appAPI.SetProxyReq()
			req.project_id=state.project.info!.project_id
			req.g_name=curg.value
			req.a_name=cura.value
			req.secret=secret.value
			req.path=new_proxy_path.value
			req.read=new_proxy_permission_read.value
			req.write=new_proxy_permission_write.value
			req.admin=new_proxy_permission_admin.value
			req.new_path=true
			client.appClient.set_proxy({"Token":state.user.token},req,client.timeout,(e: appAPI.LogicError)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: appAPI.SetProxyResp)=>{
				if(new_proxy_path.value[0]!='/'){
					new_proxy_path.value="/"+new_proxy_path.value
				}
				let tmp=new appAPI.ProxyPathInfo()
				tmp.node_id=resp.node_id
				tmp.read=new_proxy_permission_read.value
				tmp.write=new_proxy_permission_write.value
				tmp.admin=new_proxy_permission_admin.value
				proxys.value.set(new_proxy_path.value,tmp)
				new_proxy_path.value=''
				new_proxy_permission_read.value=false
				new_proxy_permission_write.value=false
				new_proxy_permission_admin.value=false
				ing.value=false
				state.clear_load()
			})
			break
		}
		case 'update_proxy':{
			let req=new appAPI.SetProxyReq()
			req.project_id=state.project.info!.project_id
			req.g_name=curg.value
			req.a_name=cura.value
			req.secret=secret.value
			req.path=cur_proxy.value
			req.read=update_proxy_permission_read.value
			req.write=update_proxy_permission_write.value
			req.admin=update_proxy_permission_admin.value
			req.new_path=false
			client.appClient.set_proxy({"Token":state.user.token},req,client.timeout,(e: appAPI.LogicError)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(_resp: appAPI.SetProxyResp)=>{
				proxys.value.get(cur_proxy.value)!.read=update_proxy_permission_read.value
				proxys.value.get(cur_proxy.value)!.write=update_proxy_permission_write.value
				proxys.value.get(cur_proxy.value)!.admin=update_proxy_permission_admin.value
				ing.value=false
				state.clear_load()
			})
			break
		}
		case 'del_proxy':{
			let req=new appAPI.DelProxyReq()
			req.project_id=state.project.info!.project_id
			req.g_name=curg.value
			req.a_name=cura.value
			req.secret=secret.value
			req.path=cur_proxy.value
			client.appClient.del_proxy({"Token":state.user.token},req,client.timeout,(e: appAPI.LogicError)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(_resp: appAPI.DelProxyResp)=>{
				proxys.value.delete(cur_proxy.value)
				cur_proxy.value=''
				ing.value=false
				state.clear_load()
			})
			break
		}
		case 'proxy':{
			let req=new appAPI.ProxyReq()
			req.project_id=state.project.info!.project_id
			req.g_name=curg.value
			req.a_name=cura.value
			req.path=cur_proxy.value
			req.data=JSON.stringify(JSON.parse(reqdata.value))
			req.force_addr=forceaddr.value
			client.appClient.proxy({"Token":state.user.token},req,client.timeout,(e: appAPI.LogicError)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(r: appAPI.ProxyResp)=>{
				if(r.data==""){
					respdata.value="{}"
				}else{
					respdata.value=JSON.stringify(JSON.parse(r.data),null,4)
				}
				ing.value=false
				state.clear_load()
			})
			break
		}
		default:{
			state.clear_load()
			state.set_alert("error",-2,"unknown operation")
		}
	}
}
function is_json_obj(str :string):boolean{
	if(str.length<2){
		return false
	}
	if(str[0]!='{'&&str[str.length-1]!='}'){
		return false
	}
	try{
		JSON.parse(str)
	}catch(e){
		return false
	}
	return true
}
</script>
<template>
	<VaModal v-model="ing" :mobileFullscreen="false" hideDefaultActions noDismiss blur :overlay="false" maxWidth="800px" maxHeight="600px" @beforeOpen="(el)=>{el.querySelector('.va-modal__dialog').style.width='auto'}">
		<template #default>
			<div v-if="optype=='del_app'" style="display:flex;flex-direction:column">
				<VaCard style="min-width:350px;witdh:auto;text-align:center" color="primary" gradient>
					<VaCardContent style="font-size:20px">
						<p><b>Delete app: {{ cura }} in group: {{ curg }}</b></p>
						<p><b>Please confirm</b></p>
					</VaCardContent>
				</VaCard>
				<div style="display:flex;justify-content:center">
					<VaButton style="width:80px;margin:10px 10px 0 0" @click="app_op" gradient>Del</VaButton>
					<VaButton style="width:80px;margin:10px 0 0 10px" @click="ing=false" gradient>Cancel</VaButton>
				</div>
			</div>
			<div v-else-if="optype=='add_app'" style="display:flex;flex-direction:column">
				<VaCard style="min-width:350px;witdh:auto;text-align:center" color="primary" gradient>
					<VaCardContent style="font-size:20px"><b>Add App</b></VaCardContent>
				</VaCard>
				<VaInput type="text" label="Group*" style="margin-top:10px" v-model.trim="new_g" />
				<VaInput type="text" label="App*" style="margin-top:10px" v-model.trim="new_a" />
				<VaInput type="text" label="Secret" style="margin-top:10px" v-model.trim="new_secret" :max-length="31" />
				<VaRadio style="display:flex;justify-content:center;margin-top:10px"
					v-model="new_discovermode"
					:options='["dns","kubernetes","static"]'
					@update:modelValue="reset_new_app()" />
				<VaInput v-if="new_discovermode=='kubernetes'" type="text" label="Kubernetes Namesapce*" style="margin-top:10px" v-model.trim="new_kubernetesns" />
				<b v-if="new_discovermode=='kubernetes'" style="color:var(--va-primary);margin-top:10px;font-size:12px">KUBERNETES SELECTOR*</b>
				<div v-if="new_discovermode=='kubernetes'" style="display:flex;justify-content:space-between">
					<VaInput type="text" label="Label Selector" style="flex:1;margin-right:5px" v-model.trim="new_kubernetesls" />
					<VaInput type="text" label="Field Selector" style="flex:1;margin-left:5px" v-model.trim="new_kubernetesfs" />
				</div>
				<VaInput v-if="new_discovermode=='dns'" type="text" label="Dns Host*" style="margin-top:10px" v-model.trim="new_dnshost" />
				<VaInput v-if="new_discovermode=='dns'" type="number" label="Dns Interval(seconds)*" style="margin-top:10px" v-model.number="new_dnsinterval" />
				<b v-if="new_discovermode=='static'&&new_staticaddrs.length>0" style="color:var(--va-primary);margin-top:10px;font-size:12px">CURRENT ADDRS*</b>
				<div v-if="new_discovermode=='static'" v-for="(addr,index) in new_staticaddrs" style="display:flex;align-items:end;margin-top:4px">
					<VaInput v-if="addr!=''" style="flex:1;margin-right:5px" readonly type="text" v-model="new_staticaddrs[index]" />
					<VaButton v-if="addr!=''" @click="new_staticaddrs.splice(index,1)" gradient>X</VaButton>
				</div>
				<div v-if="new_discovermode=='static'" style="display:flex;align-items:end;margin-top:10px">
					<VaInput style="flex:1;margin-right:5px" type="text" label="New Addr" v-model.trim="new_staticaddr" />
					<VaButton
						:disabled="new_staticaddr==''||new_staticaddrs.includes(new_staticaddr)" 
						@click="new_staticaddrs.push(new_staticaddr);new_staticaddr=''"
						gradient
					>+</VaButton>
				</div>
				<div style="display:flex;justify-content:space-around;margin-top:10px">
					<VaInput v-if="new_discovermode!=''" type="number" label="Crpc Port*" style="width:100px" v-model.number="new_crpc_port"/>
					<VaInput v-if="new_discovermode!=''" type="number" label="CGrpc Port*" style="width:100px" v-model.number="new_cgrpc_port"/>
					<VaInput v-if="new_discovermode!=''" type="number" label="Web Port*" style="width:100px" v-model.number="new_web_port"/>
				</div>
				<div style="display:flex;justify-content:center">
					<VaButton @click="app_op" style="margin:10px 10px 0 0" :disabled="!add_app_able()" gradient>Add</VaButton>
					<VaButton @click="new_g='';new_a='';new_secret='';new_discovermode='';ing=false" style="margin:10px 0 0 10px" gradient>Cancel</VaButton>
				</div>
			</div>
			<div v-else-if="optype=='update_secret'" style="display:flex;flex-direction:column">
				<VaCard style="min-width:350px;witdh:auto;text-align:center" color="primary" gradient>
					<VaCardContent style="font-size:20px"><b>Update Secret</b></VaCardContent>
				</VaCard>
				<VaInput :type="t_secret?'text':'password'" label="New Secret" style="margin-top:10px" v-model.trim="update_new_secret" :max-length="31">
					<template #appendInner>
						<VaIcon :name="t_secret?'◎':'◉'" size="small" color="var(--va-primary)" @click="t_secret=!t_secret" />
					</template>
				</VaInput>
				<div style="display:flex;justify-content:center">
					<VaButton @click="app_op" style="margin:10px 10px 0 0" :disabled="!update_secret_able()" gradient>Update</VaButton>
					<VaButton @click="update_new_secret='';ing=false" style="margin:10px 0 0 10px" gradient>Cancel</VaButton>
				</div>
			</div>
			<div v-else-if="optype=='update_discover'" style="display:flex;flex-direction:column">
				<VaCard style="min-width:350px;witdh:auto;text-align:center" color="primary" gradient>
					<VaCardContent style="font-size:20px"><b>Update Discover</b></VaCardContent>
				</VaCard>
				<VaRadio style="display:flex;justify-content:center;margin-top:10px"
					v-model="update_new_discovermode"
					:options='["dns","kubernetes","static"]'
					@update:modelValue="reset_update_discover('')"/>
				<VaInput v-if="update_new_discovermode=='kubernetes'" type="text" label="Kubernetes Namesapce*" style="margin-top:10px" v-model.trim="update_new_kubernetesns" />
				<b v-if="update_new_discovermode=='kubernetes'" style="color:var(--va-primary);margin-top:10px;font-size:12px">KUBERNETES SELECTOR*</b>
				<div v-if="update_new_discovermode=='kubernetes'" style="display:flex;justify-content:space-between">
					<VaInput type="text" label="Label Selector" style="flex:1;margin-right:5px" v-model.trim="update_new_kubernetesls" />
					<VaInput type="text" label="Field Selector" style="flex:1;margin-left:5px" v-model.trim="update_new_kubernetesfs" />
				</div>
				<VaInput v-if="update_new_discovermode=='dns'" type="text" label="Dns Host*" style="margin-top:10px" v-model.trim="update_new_dnshost" />
				<VaInput v-if="update_new_discovermode=='dns'" type="number" label="Dns Interval(seconds)*" style="margin-top:10px" v-model.number="update_new_dnsinterval" />
				<b v-if="update_new_discovermode=='static'&&update_new_staticaddrs.length>0" style="color:var(--va-primary);margin-top:10px;font-size:12px">CURRENT ADDRS*</b>
				<div v-if="update_new_discovermode=='static'" v-for="(addr,index) in update_new_staticaddrs" style="display:flex;align-items:end;margin-top:4px">
					<VaInput v-if="addr!=''" style="flex:1;margin-right:5px" readonly type="text" v-model="update_new_staticaddrs[index]" />
					<VaButton v-if="addr!=''" @click="update_new_staticaddrs.splice(index,1)" gradient>X</VaButton>
				</div>
				<div v-if="update_new_discovermode=='static'" style="display:flex;align-items:end;margin-top:10px">
					<VaInput style="flex:1;margin-right:5px" type="text" label="New Addr" v-model.trim="update_new_staticaddr" />
					<VaButton
						:disabled="update_new_staticaddr==''||update_new_staticaddrs.includes(update_new_staticaddr)" 
						@click="update_new_staticaddrs.push(update_new_staticaddr);update_new_staticaddr=''"
						gradient
					>+</VaButton>
				</div>
				<div style="display:flex;justify-content:space-around;margin-top:10px">
					<VaInput v-if="update_new_discovermode!=''" type="number" label="Crpc Port*" style="width:100px" v-model.number="update_new_crpc_port"/>
					<VaInput v-if="update_new_discovermode!=''" type="number" label="CGrpc Port*" style="width:100px" v-model.number="update_new_cgrpc_port"/>
					<VaInput v-if="update_new_discovermode!=''" type="number" label="Web Port*" style="width:100px" v-model.number="update_new_web_port"/>
				</div>
				<div style="display:flex;justify-content:center">
					<VaButton @click="app_op" style="margin:10px 10px 0 0" :disabled="!update_discover_able()" gradient>Update</VaButton>
					<VaButton @click="update_new_discovermode='';ing=false" style="margin:10px 0 0 10px" gradient>Cancel</VaButton>
				</div>
			</div>
			<div v-else-if="optype=='del_key'" style="display:flex;flex-direction:column">
				<VaCard  style="min-width:350px;width:auto;text-align:center" color="primary" gradient >
					<VaCardContent style="font-size:20px">
						<p><b>Delete key config: {{ cur_key }}</b></p>
						<p><b>Please confirm</b></p>
					</VaCardContent>
				</VaCard>
				<div style="display:flex;justify-content:center">
					<VaButton style="width:80px;margin:10px 10px 0 0" @click="app_op" gradient>Del</VaButton>
					<VaButton style="width:80px;margin:10px 0 0 10px" @click="ing=false" gradient>Cancel</VaButton>
				</div>
			</div>
			<div v-else-if="optype=='add_key'" style="display:flex;flex-direction:column">
				<VaCard style="min-width:350px;witdh:auto;text-align:center" color="primary" gradient>
					<VaCardContent style="font-size:20px"><b>Add Key Config</b></VaCardContent>
				</VaCard>
				<VaInput type="text" label="Key_Name*" style="margin-top:10px" v-model.trim="config_key" />
				<VaRadio
					style="margin-top:10px;display:flex;justify-content:space-evenly;align-items:center"
					:options='["json","raw","yaml","toml"]'
					v-model="config_value_type"
					disabled />
				<b style="font-size:13px;color:var(--va-primary)">CONTENT</b>
				<textarea
					style="border:1px solid var(--va-background-element);border-radius:5px;margin-top:10px;height:300px;resize:none"
					v-model.trim="config_value" />
				<div style="display:flex;justify-content:center">
					<VaButton style="width:80px;margin:10px 10px 0 0" @click="app_op" :disabled="!add_key_able()" gradient >Add</VaButton>
					<VaButton style="width:80px;margin:10px 0 0 10px" @click="reset_add_key();ing=false" gradient>Cancel</VaButton>
				</div>
			</div>
			<div v-else-if="optype=='update_key'" style="display:flex;flex-direction:column">
				<VaCard style="min-width:350px;witdh:auto;text-align:center" color="primary" gradient>
					<VaCardContent style="font-size:20px">
						<p><b>Update key config: {{ cur_key }}</b></p>
						<p><b>Please confirm</b></p>
					</VaCardContent>
				</VaCard>
				<div style="display:flex;justify-content:center">
					<VaButton style="width:80px;margin:10px 10px 0 0" @click="app_op" gradient>Update</VaButton>
					<VaButton style="width:80px;margin:10px 0 0 10px" @click="ing=false" gradient>Cancel</VaButton>
				</div>
			</div>
			<div v-else-if="optype=='rollback_key'" style="display:flex;flex-direction:column">
				<VaCard style="min-width:350px;witdh:auto;text-align:center" color="primary" gradient>
					<VaCardContent style="font-size:20px">
						<p><b>Rollback key config: {{ cur_key }} to config id: {{ rollback_key_index }}</b></p>
						<p><b>Please confirm</b></p>
					</VaCardContent>
				</VaCard>
				<div style="display:flex;justify-content:center">
					<VaButton style="width:80px;margin:10px 10px 0 0" @click="app_op" gradient>Rollback</VaButton>
					<VaButton style="width:80px;margin:10px 0 0 10px" @click="ing=false" gradient>Cancel</VaButton>
				</div>
			</div>
			<div v-else-if="optype=='add_proxy'" style="display:flex;flex-direction:column">
				<VaCard style="min-width:350px;witdh:auto;text-align:center" color="primary" gradient>
					<VaCardContent style="font-size:20px"><b>Add Proxy Path</b></VaCardContent>
				</VaCard>
				<VaInput label="Path" style="margin-top:10px" v-model.trim="new_proxy_path"/>
				<div style="display:flex;justify-content:space-around;margin-top:10px">
					<VaSwitch
						v-model="new_proxy_permission_read"
						true-inner-label="Read"
						false-inner-label="Read"
						@update:modelValue="()=>{
							if(!new_proxy_permission_read){
								new_proxy_permission_write=false
								new_proxy_permission_admin=false
							}
						}"
					/>
					<VaSwitch
						v-model="new_proxy_permission_write"
						true-inner-label="Write"
						false-inner-label="Write"
						@update:modelValue="()=>{
							if(!new_proxy_permission_write){
								new_proxy_permission_admin=false
							}else{
								new_proxy_permission_read=true
							}
						}"
					/>
					<VaSwitch
						v-model="new_proxy_permission_admin"
						true-inner-label="Admin"
						false-inner-label="Admin"
						@update:modelValue="()=>{
							if(new_proxy_permission_admin){
								new_proxy_permission_read=true
								new_proxy_permission_write=true
							}
						}"
					/>
				</div>
				<div style="display:flex;justify-content:center">
					<VaButton style="width:80px;margin:10px 10px 0 0" @click="app_op" :disabled="!add_proxy_able()" gradient>Add</VaButton>
					<VaButton style="width:80px;margin:10px 0 0 10px" @click="reset_add_proxy();ing=false" gradient>Cancel</VaButton>
				</div>
			</div>
			<div v-else-if="optype=='update_proxy'" style="display:flex;flex-direction:column">
				<VaCard style="min-width:350px;witdh:auto;text-align:center" color="primary" gradient>
					<VaCardContent style="font-size:20px">
						<p><b>Update proxy path: {{ cur_proxy }}</b></p>
						<p><b>Please confirm</b></p>
					</VaCardContent>
				</VaCard>
				<div style="display:flex;justify-content:center">
					<VaButton style="width:80px;margin:10px 10px 0 0" @click="app_op" gradient>Update</VaButton>
					<VaButton style="width:80px;margin:10px 0 0 10px" @click="ing=false" gradient>Cancel</VaButton>
				</div>
			</div>
			<div v-else-if="optype=='del_proxy'" style="display:flex;flex-direction:column">
				<VaCard style="min-width:350px;witdh:auto;text-align:center" color="primary" gradient>
					<VaCardContent style="font-size:20px">
						<p><b>Delete proxy path: {{ cur_proxy }}</b></p>
						<p><b>Please confirm</b></p>
					</VaCardContent>
				</VaCard>
				<div style="display:flex;justify-content:center">
					<VaButton style="width:80px;margin:10px 10px 0 0" @click="app_op" gradient>Del</VaButton>
					<VaButton style="width:80px;margin:10px 0 0 10px" @click="ing=false" gradient>Cancel</VaButton>
				</div>
			</div>
			<div v-else-if="optype=='proxy'" style="display:flex;flex-direction:column">
				<VaCard style="min-width:350px;witdh:auto;text-align:center" color="primary" gradient>
					<VaCardContent style="font-size:20px">
						<p><b>Call path: {{ cur_proxy }}</b></p>
						<p><b>Please confirm</b></p>
					</VaCardContent>
				</VaCard>
				<div style="display:flex;justify-content:center">
					<VaButton style="width:80px;margin:10px 10px 0 0" @click="app_op" gradient>Proxy</VaButton>
					<VaButton style="width:80px;margin:10px 0 0 10px" @click="ing=false" gradient>Cancel</VaButton>
				</div>
			</div>
		</template>
	</VaModal>
	<div style="flex:1;display:flex;flex-direction:column;margin:1px;overflow-y:auto">
		<div v-if="!get_app_status" style="align-self:center;display:flex;margin:1px 0">
			<VaSelect
				v-model="curg"
				:options="Object.keys(all)"
				noOptionsText="No Groups"
				placeholder="Group*"
				dropdownIcon=""
				style="width:150px;margin-right:1px"
				outline
				trigger="hover"
				:hoverOverTimeout="0"
				:hoverOutTimeout="100"
			>
				<template #option='{option,selectOption}'>
					<VaHover
						stateful
						@click="()=>{
							if(curg!=option){
								cura=''
								secret=''
								keys=new Map()
								proxys=new Map()
								//instances=[]
								get_app_status=false
								config_proxy_instance=''
							}
							selectOption(option)
						}">
						<template #default='{hover}'>
							<div
								style="padding:10px;cursor:pointer"
								:style="{'background-color':hover?'var(--va-background-border)':'',color:curg==option?'green':'black'}"
							>
								{{option}}
							</div>
						</template>
					</VaHover>
				</template>
			</VaSelect>
			<VaSelect
				v-model="cura"
				:options="curg==''?[]:Object.keys(all[curg])"
				noOptionsText="No Apps"
				placeholder="App*"
				dropdownIcon=""
				style="width:150px;margin:0 1px"
				outline
				trigger="hover"
				:hoverOverTimeout="0"
				:hoverOutTimeout="100"
			>
				<template #option='{option,selectOption}'>
					<VaHover
						stateful
						@click="()=>{
							if(cura!=option){
								secret=''
								keys=new Map()
								proxys=new Map()
								//instances=[]
								get_app_status=false
								config_proxy_instance=''
							}
							selectOption(option)
						}"
					>
						<template #default='{hover}'>
							<div
								style="padding:10px;cursor:pointer"
								:style="{'background-color':hover?'var(--va-background-border)':'',color:cura==option?'green':'black'}"
							>
								{{option}}
							</div>
						</template>
					</VaHover>
				</template>
			</VaSelect>
			<VaInput :type="t_secret?'text':'password'" v-model.trim="secret" outline placeholder="Secret" :max-length="31" style="width:250px;margin:0 1px">
				<template #appendInner>
					<VaIcon :name="t_secret?'◎':'◉'" size="small" color="var(--va-primary)" @click="t_secret=!t_secret" />
				</template>
			</VaInput>
			<VaButton style="margin:0 5px" :disabled="curg==''||cura==''" @click="get_app" gradient>Search</VaButton>
			<VaButton v-if="state.page.node!.admin" style="margin:0 5px" @click="reset_new_app();optype='add_app';ing=true" gradient>Add</VaButton>
		</div>
		<div v-else style="align-self:center;display:flex;align-items:center;margin:1px 0">
			<VaInput readonly v-model="curg" style="width:150px;margin-right:1px"/>
			<VaInput readonly v-model="cura" style="width:150px;margin-left:1px"/>
			<VaButton v-if="mustadmin()" style="margin:0 5px" @click="optype='update_secret';ing=true" gradient>UpdateSecret</VaButton>
			<VaButton v-if="mustadmin()" style="margin:0 5px" @click="reset_update_discover(discovermode);optype='update_discover';ing=true" gradient>UpdateDiscover</VaButton>
			<VaButton v-if="state.page.node!.admin&&!selfapp()" style="margin:0 5px" @click="optype='del_app';ing=true" gradient>Delete</VaButton>
			<VaButton style="margin:0 5px" @click="config_proxy_instance='';get_app_status=false" gradient>Back</VaButton>
		</div>
		<!-- configs -->
		<div
			v-if="get_app_status&&(config_proxy_instance=='config'||config_proxy_instance=='')"
			style="display:flex;align-items:center;margin:1px 0;cursor:pointer"
			:style="{'background-color':t_keys_hover?'var(--va-shadow)':'var(--va-background-element)'}"
			@click="()=>{
				if(config_proxy_instance==''){
					config_proxy_instance='config'
					cur_key=''
				}else{
					config_proxy_instance=''
				}
			}"
			@mouseover="t_keys_hover=true"
			@mouseout="t_keys_hover=false"
		>
			<span style="width:40px;padding:12px 20px;color:var(--va-primary)">{{ config_proxy_instance=='config'?'-':'+' }}</span>
			<span style="flex:1;padding:12px;color:var(--va-primary)">Configs</span>
			<VaButton
				v-if="canwrite()"
				style="height:30px;margin-right:20px"
				size="small"
				gradient
				@mouseover.stop=""
				@mouseout.stop=""
				@click.stop="optype='add_key';ing=true"
			>
				ADD
			</VaButton>
		</div>
		<!-- keys -->
		<div v-if="get_app_status&&config_proxy_instance=='config'&&keys.size>0" style="overflow-y:auto;flex:1;display:flex;flex-direction:column">
			<template v-for="key of keys.keys()">
				<div
					v-if="cur_key==''||cur_key==key"
					style="cursor:pointer;display:flex;align-items:center;margin:1px 10px"
					:style="{'background-color':keyhover==key?'var(--va-shadow)':'var(--va-background-element)'}"
					@click="()=>{
						if(cur_key==''){
							cur_key=key
							rollback_key_index=0
							edit_key_value_type=''
						}else{
							cur_key=''
						}
					}"
					@mouseover="keyhover=key"
					@mouseout="keyhover=''"
				>
					<span style="width:35px;padding:12px;color:var(--va-primary)"> {{ cur_key!=''&&cur_key==key?'-':'+' }} </span>
					<span style="padding:12px;color:var(--va-primary)">{{key}}</span>
				</div>
				<div v-if="cur_key==key" style="flex:1;display:flex;margin:1px 20px;overflow-y:auto">
					<div style="flex:1;display:flex;flex-direction:column">
						<textarea
							style="border:1px solid var(--va-background-element);border-radius:5px;flex:1;overflow-y:auto;resize:none"
							readonly>{{JSON.stringify(JSON.parse(keys.get(key)!.cur_value),null,4)}}</textarea>
						<div style="align-self:center;display:flex;align-items:center">
							<b style="color:var(--va-primary);margin:2px 10px">Current Version:  {{ keys.get(key)!.cur_version}}</b>
							<b style="color:var(--va-primary);margin:2px 10px">Current ID:  {{ keys.get(key)!.cur_index }}</b>
							<VaDropdown
								:disabled="rollback_key_index!=0||edit_key_value_type!=''"
								trigger="hover"
								placement="top"
								:hoverOverTimeout="0"
								:hoverOutTimeout="100">
								<template #anchor>
									<VaButton style="width:60px;height:30px;margin:2px 10px" size="small" gradient>History</VaButton>
								</template>
								<VaDropdownContent>
									<div style="max-height:300px;overflow-y:auto;display:flex;flex-direction:column">
										<VaButton
											v-for="index of keys.get(key)!.max_index"
											size="small"
											gradient
											style="height:24px;width:42px;padding:5px 0;margin:2px;cursor:pointer"
											:disabled="keys.get(key)!.cur_index==keys.get(key)!.max_index-index+1"
											@click="rollback_key_index=keys.get(key)!.max_index-index+1;optype='get_rollback_key';app_op()"
										>
											{{keys.get(key)!.max_index-index+1}}
										</VaButton>
									</div>
								</VaDropdownContent>
							</VaDropdown>
							<VaButton
								v-if="canwrite()"
								size="small"
								gradient
								style="width:60px;height:30px;margin:2px 10px"
								:disabled="rollback_key_index!=0||edit_key_value_type!=''"
								@click="()=>{
									if(keys.get(key)!.cur_value){
										if(keys.get(key)!.cur_value_type=='json'){
											edit_key_value=JSON.stringify(JSON.parse(keys.get(key)!.cur_value),null,4)
										}else{
											edit_key_value=keys.get(key)!.cur_value
										}
										edit_key_value_type=keys.get(key)!.cur_value_type
									}else{
										edit_key_value='{\n}'
										edit_key_value_type='json'
									}
								}"
							>
								Edit
							</VaButton>
							<VaButton
								v-if="canwrite()&&(!selfapp()||(key!='AppConfig'&&key!='SourceConfig'))"
								size="small"
								gradient
								style="width:60px;height:30px;margin:2px 10px"
								@click.stop="optype='del_key';ing=true"
							>
								Del
							</VaButton>
						</div>
					</div>
					<VaDivider v-if="rollback_key_index!=0||edit_key_value_type!=''" vertical style="margin:0 4px" />
					<div v-if="rollback_key_index!=0||edit_key_value_type!=''" style="flex:1;display:flex;flex-direction:column">
						<textarea 
							v-if="rollback_key_index!=0"
							style="border:1px solid var(--va-background-element);border-radius:5px;flex:1;overflow-y:auto;resize:none"
							readonly>{{rollback_key_value_type=='json'?JSON.stringify(JSON.parse(rollback_key_value),null,4):rollback_key_value}}</textarea>
						<textarea
							v-if="edit_key_value_type!=''"
							style="border:1px solid var(--va-background-element);border-radius:5px;flex:1;overflow-y:auto;resize:none"
							v-model.trim="edit_key_value" />
						<div style="display:flex;align-items:center">
							<VaRadio
								v-if="rollback_key_index!=0"
								:options='["json","raw","yaml","toml"]'
								v-model.trim="rollback_key_value_type"
								style="margin:4px"
								disabled
							/>
							<VaRadio
								v-if="edit_key_value_type!=''"
								:options='["json","raw","yaml","toml"]'
								v-model.trim="edit_key_value_type"
								style="margin:4px"
								disabled
							/>
							<span style="flex:1"></span>
							<VaButton
								v-if="rollback_key_index!=0&&canwrite()"
								size="small"
								gradient
								style="width:60px;height:30px;margin:2px"
								@click="optype='rollback_key';ing=true"
							>
								Rollback
							</VaButton>
							<VaButton
								v-if="edit_key_value_type!=''&&canwrite()"
								:disabled="!edit_commit_able()"
								size="small"
								gradient
								style="width:60px;height:30px;margin:2px"
								@click="optype='update_key';ing=true"
							>
								Update
							</VaButton>
							<VaButton
								size="small"
								gradient
								style="width:60px;height:30px;margin:2px"
								@click="rollback_key_index=0;edit_key_value_type=''"
							>
								Cancel
							</VaButton>
						</div>
					</div>
				</div>
			</template>
		</div>
		<div v-if="get_app_status&&config_proxy_instance=='config'&&keys.size==0">
			<p style="margin:1px 10px;padding:12px;background-color:var(--va-background-element);color:var(--va-primary)">No Config Keys</p>
		</div>
		<!-- proxys -->
		<div
			v-if="get_app_status&&(config_proxy_instance=='proxy'||config_proxy_instance=='')" 
			style="display:flex;align-items:center;margin:1px 0;cursor:pointer"
			:style="{'background-color':t_proxys_hover?'var(--va-shadow)':'var(--va-background-element)'}"
			@click="()=>{
				if(config_proxy_instance==''){
					config_proxy_instance='proxy'
					cur_proxy=''
				}else{
					config_proxy_instance=''
				}
			}"
			@mouseover="t_proxys_hover=true"
			@mouseout="t_proxys_hover=false"
		>
			<span style="width:40px;padding:12px 20px;color:var(--va-primary)">{{ config_proxy_instance=='proxy'?'-':'+' }}</span>
			<span style="flex:1;padding:12px;color:var(--va-primary)">Proxys</span>
			<VaButton
				v-if="canwrite()"
				style="height:30px;margin-right:20px"
				size="small"
				gradient
				@mouseover.stop=""
				@mouseout.stop=""
				@click.stop="reset_add_proxy();optype='add_proxy';ing=true"
			>
				ADD
			</VaButton>
		</div>
		<!-- paths -->
		<div v-if="get_app_status&&config_proxy_instance=='proxy'&&proxys.size>0" style="overflow-y:auto;flex:1;display:flex;flex-direction:column">
			<template v-for="proxy of proxys.keys()">
				<div
					v-if="cur_proxy==''||cur_proxy==proxy"
					style="cursor:pointer;display:flex;align-items:center;margin:1px 10px"
					:style="{'background-color':proxyhover==proxy?'var(--va-shadow)':'var(--va-background-element)'}"
					@mouseover="proxyhover=proxy"
					@mouseout="proxyhover=''"
					@click="()=>{
						if(cur_proxy==''){
							cur_proxy=proxy
							reqdata='{\n}'
							respdata=''
							forceaddr=''
							update_proxy_permission_read=proxys.get(proxy)!.read
							update_proxy_permission_write=proxys.get(proxy)!.write
							update_proxy_permission_admin=proxys.get(proxy)!.admin
						}else{
							cur_proxy=''
						}
					}"
				>
					<span style="width:35px;padding:12px;color:var(--va-primary)"> {{ cur_proxy!=''&&cur_proxy==proxy?'-':'+' }} </span>
					<span style="padding:12px;color:var(--va-primary)">{{proxy}}</span>
				</div>
				<div v-if="cur_proxy==proxy" style="flex:1;display:flex;margin:1px 20px;overflow-y:auto">
					<div style="flex:1;display:flex;flex-direction:column">
						<textarea
							style="border:1px solid var(--va-background-element);border-radius:5px;flex:1;overflow-y:auto;resize:none"
							v-model.trim="reqdata"
							:readonly="respdata!=''" />
						<div style="display:flex;align-items:center">
							<VaButton
								style="width:60px;height:30px;margin:2px 0"
								size="small"
								gradient
								@click="optype='proxy';ing=true"
								:disabled="respdata!=''||!is_json_obj(reqdata)"
							>
								Proxy
							</VaButton>
							<VaSelect
								:disabled="respdata!=''"
								v-model="forceaddr"
								:options="[...instances.keys()]"
								noOptionsText="No Servers"
								style="margin:2px 10px"
								placeholder="Specific Server(random when empty)"
								dropdownIcon=""
								outline
								@click="get_instances(false)">
								<template #option="{option,selectOption}">
									<VaHover stateful @click="()=>{
										if(forceaddr==option){
											selectOption(option)
											forceaddr=''
										}else{
											selectOption(option)
											forceaddr=option
										}
									}">
										<template #default="{hover}">
											<div
												style="padding:10px;cursor:pointer"
												:style="{'background-color':hover?'var(--va-background-border)':'',color:forceaddr==option?'green':'black'}">
												{{option}}{{forceaddr==option?"    (click to cancel)":""}}
											</div>
										</template>
									</VaHover>
								</template>
							</VaSelect>
							<VaSwitch
								:disabled="respdata!=''||!canwrite()"
								style="margin:2px 10px"
								v-model="update_proxy_permission_read"
								true-inner-label="Read"
								false-inner-label="Read"
								@update:model-value="()=>{
									if(!update_proxy_permission_read){
										update_proxy_permission_write=false
										update_proxy_permission_admin=false
									}
								}"
							/>
							<VaSwitch
								:disabled="respdata!=''||!canwrite()"
								style="margin:2px 10px"
								v-model="update_proxy_permission_write"
								true-inner-label="Write"
								false-inner-label="Write"
								@update:model-value="()=>{
									if(!update_proxy_permission_write){
										update_proxy_permission_admin=false
									}else{
										update_proxy_permission_read=true
									}
								}"
							/>
							<VaSwitch
								:disabled="respdata!=''||!canwrite()"
								style="margin:2px 10px"
								v-model="update_proxy_permission_admin"
								true-inner-label="Admin"
								false-inner-label="Admin"
								@update:model-value="()=>{
									if(update_proxy_permission_admin){
										update_proxy_permission_read=true
										update_proxy_permission_write=true
									}
								}"
							/>
							<VaButton
								v-if="canwrite()"
								style="width:60px;height:30px;margin:2px 10px"
								size="small"
								gradient
								:disabled="respdata!=''||!update_proxy_able()"
								@click="optype='update_proxy';ing=true"
							>
								Update
							</VaButton>
							<VaButton
								v-if="canwrite()"
								style="width:60px;height:30px;margin:2px 10px"
								size="small"
								gradient
								:disabled="respdata!=''"
								@click.stop="optype='del_proxy';ing=true"
							>
								DEL
							</VaButton>
						</div>
					</div>
					<VaDivider v-if="respdata!=''" vertical style="margin:0 4px" />
					<div v-if="respdata!=''" style="flex:1;display:flex;flex-direction:column;overflow-y:auto">
						<textarea
							style="border:1px solid var(--va-background-element);border-radius:5px;flex:1;overflow-y:auto;resize:none"
							v-model="respdata"
							readonly />
						<VaButton style="align-self:center;width:60px;height:30px;margin:2px" size="small" gradient @click="respdata=''">OK</VaButton>
					</div>
				</div>
			</template>
		</div>
		<div v-if="get_app_status&&config_proxy_instance=='proxy'&&proxys.size==0" style="width:100%">
			<p style="margin:1px 10px;padding:12px;background-color:var(--va-background-element);color:var(--va-primary)">No Proxy Paths</p>
		</div>
		<!-- instances -->
		<div
			v-if="get_app_status&&(config_proxy_instance=='instance'||config_proxy_instance=='')"
			style="display:flex;align-items:center;margin:1px 0;cursor:pointer"
			:style="{'background-color':t_instances_hover?'var(--va-shadow)':'var(--va-background-element)'}"
			@click="()=>{
				if(config_proxy_instance==''){
					instances=new Map()
					config_proxy_instance='instance'
					get_instances(true)
				}else{
					config_proxy_instance=''
				}
			}"
			@mouseover="t_instances_hover=true"
			@mouseout="t_instances_hover=false"
		>
			<span style="width:40px;padding:12px 20px;color:var(--va-primary)">{{ config_proxy_instance=='instance'?'-':'+' }}</span>
			<span style="flex:1;padding:12px;color:var(--va-primary)">Instances</span>
			<VaButton
				v-if="config_proxy_instance=='instance'"
				style="width:60px;height:30px;margin-right:20px"
				size="small"
				gradient
				@mouseover.stop=""
				@mouseout.stop=""
				@click.stop="get_instances(true)"
			>
				refresh
			</VaButton>
		</div>
		<div v-if="get_app_status&&config_proxy_instance=='instance'" style="overflow-y:auto;flex:1;display:flex;flex-wrap:wrap">
			<div v-if="instances.size==0"
				style="width:300px;height:150px;border:1px solid var(--va-primary);border-radius:5px;margin:5px;display:flex;justify-content:center;align-items:center">
				No Instances
			</div>
			<div v-for="instanceaddr of instances.keys()"
				style="position:relative;width:300px;height:150px;margin:5px;border:1px solid var(--va-primary);border-radius:5px">
				<VaButton style="position:absolute;right:1px;top:1px" size="small" gradient @click="get_instance(instanceaddr)">refresh</VaButton>
				<div style="width:100%;height:100%;display:flex;flex-direction:column;justify-content:space-around">
					<div style="margin:1px;display:flex">
						<span style="width:90px;margin-left:10px">Addr</span>
						<VaDivider vertical />
						<span>{{instanceaddr}}</span>
					</div>
					<div v-if="!instances.get(instanceaddr)||instances.get(instanceaddr)!.cpu_num==0" style="margin:1px;display:flex">
						<span style="width:90px;margin-left:10px">SysInfo</span>
						<VaDivider vertical />
						<span>get failed</span>
					</div>
					<div v-if="instances.get(instanceaddr)&&instances.get(instanceaddr)!.cpu_num!=0" style="margin:1px;display:flex">
						<span style="width:90px;margin-left:10px">CpuNum</span>
						<VaDivider vertical />
						<span>{{instances.get(instanceaddr)!.cpu_num}}</span>
					</div>
					<div v-if="instances.get(instanceaddr)&&instances.get(instanceaddr)!.cpu_num!=0" style="margin:1px;display:flex">
						<span style="width:90px;margin-left:10px">CpuUsage</span>
						<VaDivider vertical />
						<span>{{(instances.get(instanceaddr)!.cur_cpu_usage*100).toFixed(2)}}%</span>
					</div>
					<div v-if="instances.get(instanceaddr)&&instances.get(instanceaddr)!.cpu_num!=0" style="margin:1px;display:flex">
						<span style="width:90px;margin-left:10px">MemTotal</span>
						<VaDivider vertical />
						<span>{{(Number(instances.get(instanceaddr)!.total_mem)/1024/1024).toFixed(2)}}MB</span>
					</div>
					<div v-if="instances.get(instanceaddr)&&instances.get(instanceaddr)!.cpu_num!=0" style="margin:1px;display:flex">
						<span style="width:90px;margin-left:10px">MemUsage</span>
						<VaDivider vertical />
						<span>{{(Number(instances.get(instanceaddr)!.cur_mem_usage)/Number(instances.get(instanceaddr)!.total_mem)*100).toFixed(2)}}%</span>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>
