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
const dnshost=ref<string>("")
const dnsinterval=ref<number>(0)
const staticaddrs=ref<string[]>([])

const keys=ref<Map<string,appAPI.KeyConfigInfo>>(new Map())
const t_keys_hover=ref<boolean>(false)
const keyhover=ref<string>("")

const proxys=ref<Map<string,appAPI.ProxyPathInfo>>(new Map())
const t_proxys_hover=ref<boolean>(false)
const proxyhover=ref<string>("")

//const instances=ref<appAPI.InstanceInfo[]>([])
//const t_instances_hover=ref<boolean>(false)

const get_app_status=ref<boolean>(false)

function get_app(){
	if(curg.value==""||cura.value==""){
		keys.value=new Map()
		proxys.value=new Map()
		config_proxy_instance.value=""
		discovermode.value=""
		kubernetesns.value=""
		kubernetesls.value=""
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
		dnshost.value=""
		dnsinterval.value=0
		staticaddrs.value=[]
		state.set_alert("error",-2,"Missing node_id on Group:"+curg.value+" App:"+cura.value)
		return
	}
	if(!state.set_load()){
		return
	}
	let req = {
		project_id:state.project.info!.project_id,
		g_name:curg.value,
		a_name:cura.value,
		secret:secret.value,
	}
	client.appClient.get_app({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
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
	      	dnshost.value=resp.dns_host
	      	dnsinterval.value=resp.dns_interval
		if(resp.static_addrs){
			staticaddrs.value=resp.static_addrs
		}else{
			staticaddrs.value=[]
		}
		get_app_status.value=true
		state.clear_load()
	})
}
function get_instances(){
	if(!get_app_status.value){
		instances.value=[]
		return
	}
	if(!state.set_load()){
		return
	}
	let req = {
		project_id:state.project.info.project_id,
		g_name:curg.value,
		a_name:cura.value,
		secret:secret.value,
	}
	client.appClient.get_instances({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
	},(resp: appAPI.GetAppInstancesResp)=>{
		if(resp.instances){
			let tmp:appAPI.InstanceInfo[]=[]
			for(let i=0;i<resp.instances.length;i++){
				if(resp.instances[i]){
					tmp.push(resp.instances[i]!)
				}
			}
			instances.value=tmp
		}else{
			instances.value=[]
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
const new_dnshost=ref<string>("")
const new_dnsinterval=ref<number>(0)
const new_staticaddrs=ref<string[]>([])
const new_staticaddr=ref<string>("")
function reset_new_app(){
	new_kubernetesns.value=''
	new_kubernetesls.value=''
	new_dnshost.value=''
	new_dnsinterval.value=0
	new_staticaddrs.value=[]
	new_staticaddr.value=''
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
	if(new_discovermode.value=='kubernetes'&&(new_kubernetesns.value==''||new_kubernetesls.value=='')){
		return false
	}
	if(new_discovermode.value=='dns'&&new_dnshost.value==''){
		return false
	}
	if(new_discovermode.value=='static'&&new_staticaddrs.value.length==0){
		return false
	}
	return true
}

//update app
const update_new_secret=ref<string>("")
const update_new_discovermode=ref<string>("")
const update_new_kubernetesns=ref<string>("")
const update_new_kubernetesls=ref<string>("")
const update_new_dnshost=ref<string>("")
const update_new_dnsinterval=ref<number>(0)
const update_new_staticaddrs=ref<string[]>([])
const update_new_staticaddr=ref<string>("")
function update_secret_able():boolean{
	return update_new_secret.value != secret.value
}
function reset_update_discover(dname:string){
	if(dname!=''){
		update_new_discovermode.value=dname
	}
	update_new_kubernetesns.value=''
	update_new_kubernetesls.value=''
	update_new_dnshost.value=''
	update_new_dnsinterval.value=0
	update_new_staticaddrs.value=[]
	update_new_staticaddr.value=''
	if(update_new_discovermode.value=="kubernetes"){
		update_new_kubernetesns.value=kubernetesns.value
		update_new_kubernetesls.value=kubernetesls.value
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
	if(update_new_discovermode.value=='kubernetes'){
		if(update_new_kubernetesns.value==''||update_new_kubernetesls.value==''){
			return false
		}
		if(update_new_kubernetesns.value==kubernetesns.value&&update_new_kubernetesls.value==kubernetesls.value){
			return false
		}
	}
	if(update_new_discovermode.value=='dns'){
		if(update_new_dnshost.value==''){
			return false
		}
		if(update_new_dnshost.value==dnshost.value&&update_new_dnsinterval.value==dnsinterval.value){
			return false
		}
	}
	if(update_new_discovermode.value=='static'){
		if(update_new_staticaddrs.value.length==0){
			return false
		}
		if(update_new_staticaddrs.value.length==staticaddrs.value.length){
			if(staticaddrs.value.every(function(v,i){return v==update_new_staticaddrs.value[i]})){
				return false
			}
		}
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

function app_op(){
	if(!state.set_load()){
		return
	}
	switch(optype.value){
		case 'del_app':{
			let req = {
				project_id:state.project.info!.project_id,
				g_name:curg.value,
				a_name:cura.value,
				secret:secret.value,
			}
			client.appClient.del_app({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
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
			let req = {
				project_id:state.project.info!.project_id,
				g_name:new_g.value,
				a_name:new_a.value,
				secret:new_secret.value,
				discover_mode:new_discovermode.value,
				kubernetes_namespace:new_kubernetesns.value,
	      			kubernetes_labelselector:new_kubernetesls.value,
	      			dns_host:new_dnshost.value,
	      			dns_interval:new_dnsinterval.value,
	      			static_addrs:new_staticaddrs.value,
	      			new_app:true,
			}
			client.appClient.set_app({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: appAPI.SetAppResp)=>{
				if(all.value[new_g.value] && all.value[new_g.value][new_a.value]){
					return
				}
				if(state.page.node!.children){
					state.page.node!.children!.push({
						node_id:resp.node_id,
						node_name:new_g.value+"."+new_a.value,
						node_data:"",
						canread:true,
						canwrite:true,
						admin:true,
						children:[],
					})
				}else{
					state.page.node!.children=[{
						node_id:resp.node_id,
						node_name:new_g.value+"."+new_a.value,
						node_data:"",
						canread:true,
						canwrite:true,
						admin:true,
						children:[],
					}]
				}
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
			let req = {
				project_id:state.project.info!.project_id,
				g_name:curg.value,
				a_name:cura.value,
				secret:secret.value,
				discover_mode:update_new_discovermode.value,
				kubernetes_namespace:update_new_kubernetesns.value,
	      			kubernetes_labelselector:update_new_kubernetesls.value,
	      			dns_host:update_new_dnshost.value,
	      			dns_interval:update_new_dnsinterval.value,
	      			static_addrs:update_new_staticaddrs.value,
	      			new_app:false,
			}
			client.appClient.set_app({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
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
				ing.value=false
				state.clear_load()
			})
			break
		}
		case 'update_secret':{
			let req = {
				project_id:state.project.info!.project_id,
				g_name:curg.value,
				a_name:cura.value,
				old_secret:secret.value,
				new_secret:update_new_secret.value
			}
			client.appClient.update_app_secret({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
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
			let req = {
				project_id:state.project.info!.project_id,
				g_name:curg.value,
				a_name:cura.value,
				secret:secret.value,
				key:cur_key.value,
				index:rollback_key_index.value,
			}
			client.appClient.get_key_config({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
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
			let req = {
				project_id:state.project.info!.project_id,
				g_name:curg.value,
				a_name:cura.value,
				secret:secret.value,
				key:config_key.value,
				value:config_value.value,
				value_type:config_value_type.value,
				new_key:true,
			}
			client.appClient.set_key_config({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(_resp: appAPI.SetKeyConfigResp)=>{
				keys.value.set(config_key.value,{
					cur_index:1,
					max_index:1,
					cur_version:1,
					cur_value:config_value.value,
					cur_value_type:config_value_type.value,
				})
				config_key.value = ''
				config_value.value = '{\n}'
				config_value_type.value = 'json'
				ing.value=false
				state.clear_load()
			})
			break
		}
		case 'update_key':{
			let req = {
				project_id:state.project.info!.project_id,
				g_name:curg.value,
				a_name:cura.value,
				secret:secret.value,
				key:cur_key.value,
				value:edit_key_value.value,
				value_type:edit_key_value_type.value,
				new_key:false,
			}
			client.appClient.set_key_config({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
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
			let req = {
				project_id:state.project.info!.project_id,
				g_name:curg.value,
				a_name:cura.value,
				secret:secret.value,
				key:cur_key.value,
				index:rollback_key_index.value,
			}
			client.appClient.rollback({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
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
			let req = {
				project_id:state.project.info!.project_id,
				g_name:curg.value,
				a_name:cura.value,
				secret:secret.value,
				key:cur_key.value,
			}
			client.appClient.del_key({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
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
			let req = {
				project_id:state.project.info!.project_id,
				g_name:curg.value,
				a_name:cura.value,
				secret:secret.value,
				path:new_proxy_path.value,
				read:new_proxy_permission_read.value,
				write:new_proxy_permission_write.value,
				admin:new_proxy_permission_admin.value,
				new_path:true,
			}
			client.appClient.set_proxy({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: appAPI.SetProxyResp)=>{
				if(new_proxy_path.value[0]!='/'){
					new_proxy_path.value="/"+new_proxy_path.value
				}
				proxys.value.set(new_proxy_path.value,{
					node_id:resp.node_id,
					read:new_proxy_permission_read.value,
					write:new_proxy_permission_write.value,
					admin:new_proxy_permission_admin.value,
				})
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
			let req = {
				project_id:state.project.info!.project_id,
				g_name:curg.value,
				a_name:cura.value,
				secret:secret.value,
				path:cur_proxy.value,
				read:update_proxy_permission_read.value,
				write:update_proxy_permission_write.value,
				admin:update_proxy_permission_admin.value,
				new_path:false,
			}
			client.appClient.set_proxy({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
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
			let req = {
				project_id:state.project.info!.project_id,
				g_name:curg.value,
				a_name:cura.value,
				secret:secret.value,
				path:cur_proxy.value,
			}
			client.appClient.del_proxy({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
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
			let data=JSON.stringify(JSON.parse(reqdata.value))
			let req = {
				project_id:state.project.info!.project_id,
				g_name:curg.value,
				a_name:cura.value,
				path:cur_proxy.value,
				data:data,
			}
			client.appClient.proxy({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
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
	<va-modal v-model="ing" attach-element="#app" max-width="800px" max-height="600px" hide-default-actions no-dismiss overlay-opacity="0.2" z-index="999">
		<template #default>
			<div v-if="optype=='del_app'" style="display:flex;flex-direction:column">
				<va-card style="min-width:350px;witdh:auto;text-align:center" color="primary" gradient>
					<va-card-content style="font-size:20px">
						<p><b>Delete app: {{ cura }} in group: {{ curg }}</b></p>
						<p><b>Please confirm</b></p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:10px 10px 0 0" @click="app_op" gradient>Del</va-button>
					<va-button style="width:80px;margin:10px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='add_app'" style="display:flex;flex-direction:column">
				<va-card style="min-width:350px;witdh:auto;text-align:center" color="primary" gradient>
					<va-card-content style="font-size:20px"><b>Add App</b></va-card-content>
				</va-card>
				<va-input type="text" label="Group*" style="margin-top:10px" v-model.trim="new_g" />
				<va-input type="text" label="App*" style="margin-top:10px" v-model.trim="new_a" />
				<va-input type="text" label="Secret" style="margin-top:10px" v-model.trim="new_secret" :max-length="31" />
				<va-radio style="display:flex;justify-content:center;margin-top:10px"
					v-model="new_discovermode"
					:options='["dns","kubernetes","static"]'
					@update:modelValue="reset_new_app()" />
				<va-input v-if="new_discovermode=='kubernetes'" type="text" label="Kubernetes Namesapce*" style="margin-top:10px" v-model.trim="new_kubernetesns" />
				<va-input v-if="new_discovermode=='kubernetes'" type="text" label="Kubernetes Label Selector*" style="margin-top:10px" v-model.trim="new_kubernetesls" />
				<va-input v-if="new_discovermode=='dns'" type="text" label="Dns Host*" style="margin-top:10px" v-model.trim="new_dnshost" />
				<va-input v-if="new_discovermode=='dns'" type="text" label="Dns Interval(seconds)*" style="margin-top:10px" v-model.number="new_dnsinterval" />
				<div v-if="new_discovermode=='static'" v-for="(addr,index) in new_staticaddrs" style="display:flex;align-items:end;margin-top:10px">
					<va-input v-if="addr!=''" style="flex:1;margin-right:5px" readonly type="text" v-model="new_staticaddrs[index]" />
					<va-button v-if="addr!=''" @click="new_staticaddrs.splice(index,1)" gradient>X</va-button>
				</div>
				<div v-if="new_discovermode=='static'" style="display:flex;align-items:end;margin-top:10px">
					<va-input style="flex:1;margin-right:5px" type="text" label="New Addr" v-model.trim="new_staticaddr" />
					<va-button
						:disabled="new_staticaddr==''||new_staticaddrs.includes(new_staticaddr)" 
						@click="new_staticaddrs.push(new_staticaddr);new_staticaddr=''"
						gradient
					>+</va-button>
				</div>
				<div style="display:flex;justify-content:center">
					<va-button @click="app_op" style="margin:10px 10px 0 0" :disabled="!add_app_able()" gradient>Add</va-button>
					<va-button @click="new_g='';new_a='';new_secret='';new_discovermode='';ing=false" style="margin:10px 0 0 10px" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='update_secret'" style="display:flex;flex-direction:column">
				<va-card style="min-width:350px;witdh:auto;text-align:center" color="primary" gradient>
					<va-card-content style="font-size:20px"><b>Update Secret</b></va-card-content>
				</va-card>
				<va-input type="text" label="New Secret" style="margin-top:10px" v-model.trim="update_new_secret" />
				<div style="display:flex;justify-content:center">
					<va-button @click="app_op" style="margin:10px 10px 0 0" :disabled="!update_secret_able()" gradient>Update</va-button>
					<va-button @click="update_new_secret='';ing=false" style="margin:10px 0 0 10px" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='update_discover'" style="display:flex;flex-direction:column">
				<va-card style="min-width:350px;witdh:auto;text-align:center" color="primary" gradient>
					<va-card-content style="font-size:20px"><b>Update Discover</b></va-card-content>
				</va-card>
				<va-radio style="display:flex;justify-content:center;margin-top:10px"
					v-model="update_new_discovermode"
					:options='["dns","kubernetes","static"]'
					@update:modelValue="reset_update_discover('')"/>
				<va-input v-if="update_new_discovermode=='kubernetes'" type="text" label="Kubernetes Namesapce*" style="margin-top:10px" v-model.trim="update_new_kubernetesns" />
				<va-input v-if="update_new_discovermode=='kubernetes'" type="text" label="Kubernetes Label Selector*" style="margin-top:10px" v-model.trim="update_new_kubernetesls" />
				<va-input v-if="update_new_discovermode=='dns'" type="text" label="Dns Host*" style="margin-top:10px" v-model.trim="update_new_dnshost" />
				<va-input v-if="update_new_discovermode=='dns'" type="text" label="Dns Interval(seconds)*" style="margin-top:10px" v-model.number="update_new_dnsinterval" />
				<div v-if="update_new_discovermode=='static'" v-for="(addr,index) in update_new_staticaddrs" style="display:flex;align-items:end;margin-top:10px">
					<va-input v-if="addr!=''" style="flex:1;margin-right:5px" readonly type="text" v-model="update_new_staticaddrs[index]" />
					<va-button v-if="addr!=''" @click="update_new_staticaddrs.splice(index,1)" gradient>X</va-button>
				</div>
				<div v-if="update_new_discovermode=='static'" style="display:flex;align-items:end;margin-top:10px">
					<va-input style="flex:1;margin-right:5px" type="text" label="New Addr" v-model.trim="update_new_staticaddr" />
					<va-button
						:disabled="update_new_staticaddr==''||update_new_staticaddrs.includes(update_new_staticaddr)" 
						@click="update_new_staticaddrs.push(update_new_staticaddr);update_new_staticaddr=''"
						gradient
					>+</va-button>
				</div>
				<div style="display:flex;justify-content:center">
					<va-button @click="app_op" style="margin:10px 10px 0 0" :disabled="!update_discover_able()" gradient>Update</va-button>
					<va-button @click="update_new_discovermode='';ing=false" style="margin:10px 0 0 10px" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='del_key'" style="display:flex;flex-direction:column">
				<va-card  style="min-width:350px;width:auto;text-align:center" color="primary" gradient >
					<va-card-content style="font-size:20px">
						<p><b>Delete key config: {{ cur_key }}</b></p>
						<p><b>Please confirm</b></p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:10px 10px 0 0" @click="app_op" gradient>Del</va-button>
					<va-button style="width:80px;margin:10px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='add_key'" style="display:flex;flex-direction:column">
				<va-card style="min-width:350px;witdh:auto;text-align:center" color="primary" gradient>
					<va-card-content style="font-size:20px"><b>Add Key Config</b></va-card-content>
				</va-card>
				<va-input type="text" label="Key_Name*" style="margin-top:10px;width:600px" v-model.trim="config_key" />
				<va-radio
					style="margin-top:10px;display:flex;justify-content:space-evenly;align-items:center"
					:options='["json","raw","yaml","toml"]'
					v-model="config_value_type"
					disabled />
				<va-textarea
					style="margin-top:10px;width:600px"
					label="Content"
					v-model.trim="config_value"
					:minRows="15"
					:maxRows="15"
					autosize
					:resize="false"
					:rules="[(v)=>is_json_obj(v)]"/>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:10px 10px 0 0" @click="app_op" :disabled="!add_key_able()" gradient >Add</va-button>
					<va-button style="width:80px;margin:10px 0 0 10px" @click="reset_add_key();ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='update_key'" style="display:flex;flex-direction:column">
				<va-card style="min-width:350px;witdh:auto;text-align:center" color="primary" gradient>
					<va-card-content style="font-size:20px">
						<p><b>Update key config: {{ cur_key }}</b></p>
						<p><b>Please confirm</b></p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:10px 10px 0 0" @click="app_op" gradient>Update</va-button>
					<va-button style="width:80px;margin:10px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='rollback_key'" style="display:flex;flex-direction:column">
				<va-card style="min-width:350px;witdh:auto;text-align:center" color="primary" gradient>
					<va-card-content style="font-size:20px">
						<p><b>Rollback key config: {{ cur_key }} to config id: {{ rollback_key_index }}</b></p>
						<p><b>Please confirm</b></p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:10px 10px 0 0" @click="app_op" gradient>Rollback</va-button>
					<va-button style="width:80px;margin:10px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='add_proxy'" style="display:flex;flex-direction:column">
				<va-card style="min-width:350px;witdh:auto;text-align:center" color="primary" gradient>
					<va-card-content style="font-size:20px"><b>Add Proxy Path</b></va-card-content>
				</va-card>
				<va-input label="Path" style="width:500px;margin-top:10px" v-model.trim="new_proxy_path"/>
				<div style="display:flex;justify-content:space-around;margin-top:10px">
					<va-switch
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
					<va-switch
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
					<va-switch
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
					<va-button style="width:80px;margin:10px 10px 0 0" @click="app_op" :disabled="!add_proxy_able()" gradient>Add</va-button>
					<va-button style="width:80px;margin:10px 0 0 10px" @click="reset_add_proxy();ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='update_proxy'" style="display:flex;flex-direction:column">
				<va-card style="min-width:350px;witdh:auto;text-align:center" color="primary" gradient>
					<va-card-content style="font-size:20px">
						<p><b>Update proxy path: {{ cur_proxy }}</b></p>
						<p><b>Please confirm</b></p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:10px 10px 0 0" @click="app_op" gradient>Update</va-button>
					<va-button style="width:80px;margin:10px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='del_proxy'" style="display:flex;flex-direction:column">
				<va-card style="min-width:350px;witdh:auto;text-align:center" color="primary" gradient>
					<va-card-content style="font-size:20px">
						<p><b>Delete proxy path: {{ cur_proxy }}</b></p>
						<p><b>Please confirm</b></p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:10px 10px 0 0" @click="app_op" gradient>Del</va-button>
					<va-button style="width:80px;margin:10px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='proxy'" style="display:flex;flex-direction:column">
				<va-card style="min-width:350px;witdh:auto;text-align:center" color="primary" gradient>
					<va-card-content style="font-size:20px">
						<p><b>Call path: {{ cur_proxy }}</b></p>
						<p><b>Please confirm</b></p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:10px 10px 0 0" @click="app_op" gradient>Proxy</va-button>
					<va-button style="width:80px;margin:10px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
		</template>
	</va-modal>
	<div style="flex:1;display:flex;flex-direction:column;align-items:center;margin:1px;overflow-y:auto">
		<div v-if="!get_app_status" style="display:flex;margin:1px 0">
			<va-select
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
					<va-hover
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
					</va-hover>
				</template>
			</va-select>
			<va-select
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
					<va-hover
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
					</va-hover>
				</template>
			</va-select>
			<va-input :type="t_secret?'text':'password'" v-model.trim="secret" outline placeholder="Secret" :max-length="31" style="width:250px;margin:0 1px">
				<template #appendInner>
					<va-icon :name="t_secret?'◎':'◉'" size="small" color="var(--va-primary)" @click="t_secret=!t_secret" />
				</template>
			</va-input>
			<va-button style="margin:0 5px" :disabled="curg==''||cura==''" @click="get_app" gradient>Search</va-button>
			<va-button v-if="state.page.node!.admin" style="margin:0 5px" @click="reset_new_app();optype='add_app';ing=true" gradient>Add</va-button>
		</div>
		<div v-else style="display:flex;margin:1px 0;align-self:center">
			<va-button v-if="mustadmin()" style="margin:0 5px" @click="optype='update_secret';ing=true" gradient>UpdateSecret</va-button>
			<va-button v-if="mustadmin()" style="margin:0 5px" @click="reset_update_discover(discovermode);optype='update_discover';ing=true" gradient>UpdateDiscover</va-button>
			<va-button v-if="state.page.node!.admin&&!selfapp()" style="margin:0 5px" @click="optype='del_app';ing=true" gradient>Delete</va-button>
			<va-button style="margin:0 5px" @click="config_proxy_instance='';get_app_status=false" gradient>Back</va-button>
		</div>
		<!-- configs -->
		<div
			v-if="get_app_status&&(config_proxy_instance=='config'||config_proxy_instance=='')"
			style="width:100%;display:flex;align-items:center;margin:1px 0;cursor:pointer"
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
			<va-button
				v-if="canwrite()"
				style="height:30px;margin-right:20px"
				size="small"
				gradient
				@mouseover.stop=""
				@mouseout.stop=""
				@click.stop="optype='add_key';ing=true"
			>
				ADD
			</va-button>
		</div>
		<!-- keys -->
		<div v-if="get_app_status&&config_proxy_instance=='config'&&keys.size>0" style="width:100%;overflow-y:auto;flex:1;display:flex;flex-direction:column">
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
					<div style="flex:1;display:flex;flex-direction:column;overflow-y:auto">
						<va-textarea
							:modelValue="JSON.stringify(JSON.parse(keys.get(key)!.cur_value),null,4)"
							style="flex:1"
							readonly
							autosize
							:resize='false'/>
						<div style="align-self:center;display:flex;align-items:center">
							<b style="color:var(--va-primary);margin:2px 10px">Current Config ID:  {{ keys.get(key)!.cur_index }}</b>
							<va-dropdown
								:disabled="rollback_key_index!=0||edit_key_value_type!=''"
								trigger="hover"
								placement="top"
								:hoverOverTimeout="0"
								:hoverOutTimeout="100">
								<template #anchor>
									<va-button style="width:60px;height:30px;margin:2px 10px" size="small" gradient>History</va-button>
								</template>
								<va-dropdown-content>
									<div style="max-height:300px;overflow-y:auto;display:flex;flex-direction:column">
										<va-button
											v-for="index of keys.get(key)!.max_index"
											size="small"
											gradient
											style="height:24px;width:42px;padding:5px 0;margin:2px;cursor:pointer"
											:disabled="keys.get(key)!.cur_index==keys.get(key)!.max_index-index+1"
											@click="rollback_key_index=keys.get(key)!.max_index-index+1;optype='get_rollback_key';app_op()"
										>
											{{keys.get(key)!.max_index-index+1}}
										</va-button>
									</div>
								</va-dropdown-content>
							</va-dropdown>
							<va-button
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
							</va-button>
							<va-button
								v-if="canwrite()&&(!selfapp()||(key!='AppConfig'&&key!='SourceConfig'))"
								size="small"
								gradient
								style="width:60px;height:30px;margin:2px 10px"
								@click.stop="optype='del_key';ing=true"
							>
								Del
							</va-button>
						</div>
					</div>
					<va-divider v-if="rollback_key_index!=0||edit_key_value_type!=''" vertical style="margin:0 4px" />
					<div v-if="rollback_key_index!=0||edit_key_value_type!=''" style="flex:1;display:flex;flex-direction:column">
						<va-textarea
							v-if="rollback_key_index!=0"
							style="flex:1"
							:modelValue="rollback_key_value_type=='json'?JSON.stringify(JSON.parse(rollback_key_value),null,4):rollback_key_value"
							readonly
							autosize
							:resize="false"/>
						<va-textarea
							v-if="edit_key_value_type!=''"
							style="flex:1"
							v-model.trim="edit_key_value"
							autosize
							:resize="false">
						</va-textarea>
						<div style="display:flex;align-items:center">
							<va-radio
								v-if="rollback_key_index!=0"
								:options='["json","raw","yaml","toml"]'
								v-model.trim="rollback_key_value_type"
								style="margin:4px"
								disabled
							/>
							<va-radio
								v-if="edit_key_value_type!=''"
								:options='["json","raw","yaml","toml"]'
								v-model.trim="edit_key_value_type"
								style="margin:4px"
								disabled
							/>
							<span style="flex:1"></span>
							<va-button
								v-if="rollback_key_index!=0&&canwrite()"
								size="small"
								gradient
								style="width:60px;height:30px;margin:2px"
								@click="optype='rollback_key';ing=true"
							>
								Rollback
							</va-button>
							<va-button
								v-if="edit_key_value_type!=''&&canwrite()"
								:disabled="!edit_commit_able()"
								size="small"
								gradient
								style="width:60px;height:30px;margin:2px"
								@click="optype='update_key';ing=true"
							>
								Update
							</va-button>
							<va-button
								size="small"
								gradient
								style="width:60px;height:30px;margin:2px"
								@click="rollback_key_index=0;edit_key_value_type=''"
							>
								Cancel
							</va-button>
						</div>
					</div>
				</div>
			</template>
		</div>
		<div style="width:100%" v-if="get_app_status&&config_proxy_instance=='config'&&keys.size==0">
			<p style="margin:1px 10px;padding:12px;background-color:var(--va-background-element);color:var(--va-primary)">No Config Keys</p>
		</div>
		<!-- proxys -->
		<div
			v-if="get_app_status&&(config_proxy_instance=='proxy'||config_proxy_instance=='')" 
			style="width:100%;display:flex;align-items:center;margin:1px 0;cursor:pointer"
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
			<va-button
				v-if="canwrite()"
				style="height:30px;margin-right:20px"
				size="small"
				gradient
				@mouseover.stop=""
				@mouseout.stop=""
				@click.stop="reset_add_proxy();optype='add_proxy';ing=true"
			>
				ADD
			</va-button>
		</div>
		<!-- paths -->
		<div v-if="get_app_status&&config_proxy_instance=='proxy'&&proxys.size>0" style="width:100%;overflow-y:auto;flex:1;display:flex;flex-direction:column">
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
							update_proxy_permission_read=proxys.get(proxy)!.read
							update_proxy_permission_write=proxys.get(proxy)!.write
							update_proxy_permission_admin=proxys.get(proxy)!.admin
							respdata=''
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
						<va-textarea style="flex:1" v-model.trim="reqdata" :readonly="respdata!=''" autosize :resize="false" />
						<div style="width:100%;display:flex">
							<va-button
								style="width:60px;height:30px;margin:2px 0"
								size="small"
								gradient
								@click="optype='proxy';ing=true"
								:disabled="respdata!=''||!is_json_obj(reqdata)"
							>
								Proxy
							</va-button>
							<div style="flex:1"></div>
							<va-switch
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
							<va-switch
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
							<va-switch
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
							<va-button
								v-if="canwrite()"
								style="width:60px;height:30px;margin:2px 10px"
								size="small"
								gradient
								:disabled="respdata!=''||!update_proxy_able()"
								@click="optype='update_proxy';ing=true"
							>
								Update
							</va-button>
							<va-button
								v-if="canwrite()"
								style="width:60px;height:30px;margin:2px 10px"
								size="small"
								gradient
								:disabled="respdata!=''"
								@click.stop="optype='del_proxy';ing=true"
							>
								DEL
							</va-button>
						</div>
					</div>
					<va-divider v-if="respdata!=''" vertical style="margin:0 4px" />
					<div v-if="respdata!=''" style="flex:1;display:flex;flex-direction:column">
						<va-textarea style="flex:1" v-model="respdata" readonly autosize :resize="false" />
						<va-button style="align-self:center;width:60px;height:30px;margin:2px" size="small" gradient @click="respdata=''">OK</va-button>
					</div>
				</div>
			</template>
		</div>
		<div style="width:100%" v-if="get_app_status&&config_proxy_instance=='proxy'&&proxys.size==0">
			<p style="margin:1px 10px;padding:12px;background-color:var(--va-background-element);color:var(--va-primary)">No Proxy Paths</p>
		</div>
		<!-- instances -->
		<!-- <div -->
		<!-- 	v-if="get_app_status&&(config_proxy_instance=='instance'||config_proxy_instance=='')" -->
		<!-- 	style="display:flex;align-items:center;margin:1px 0;cursor:pointer" -->
		<!-- 	:style="{'background-color':t_instances_hover?'var(--va-shadow)':'var(--va-background-element)'}" -->
		<!-- 	@click="()=>{ -->
		<!-- 		if(config_proxy_instance==''){ -->
		<!-- 			instances=[] -->
		<!-- 			config_proxy_instance='instance' -->
		<!-- 			get_instances() -->
		<!-- 		}else{ -->
		<!-- 			config_proxy_instance='' -->
		<!-- 		} -->
		<!-- 	}" -->
		<!-- 	@mouseover="t_instances_hover=true" -->
		<!-- 	@mouseout="t_instances_hover=false" -->
		<!-- > -->
		<!-- 	<span style="width:40px;padding:12px 20px;color:var(--va-primary)">{{ config_proxy_instance=='instance'?'-':'+' }}</span> -->
		<!-- 	<span style="flex:1;padding:12px;color:var(--va-primary)">Instances</span> -->
		<!-- 	<va-button -->
		<!-- 		v-if="config_proxy_instance=='instance'" -->
		<!-- 		style="width:60px;height:30px" -->
		<!-- 		size="small" -->
		<!-- 		@mouseover.stop="" -->
		<!-- 		@mouseout.stop="" -->
		<!-- 		@click.stop="get_instances" -->
		<!-- 	> -->
		<!-- 		refresh -->
		<!-- 	</va-button> -->
		<!-- 	<span style="width:60px;padding:12px 20px;color:var(--va-primary)">{{ config_proxy_instance?'▲':'▼' }}</span> -->
		<!-- </div> -->
		<!-- <div v-if="config_proxy_instance=='instance'&&instances.length>0" style="display:flex;flex-wrap:wrap"> -->
		<!-- 	<div v-for="instance of instances" style="border:1px solid var(--va-primary);border-radius:5px;margin:5px;display:flex;flex-direction:column;align-items:center"> -->
		<!-- 		<div style="margin:1px;display:flex"> -->
		<!-- 			<span style="width:85px">Host IP</span> -->
		<!-- 			<va-divider vertical /> -->
		<!-- 			<span style="width:200px;word-break:break-all">{{instance.host_ip}}</span> -->
		<!-- 		</div> -->
		<!-- 		<div style="margin:1px;display:flex"> -->
		<!-- 			<span style="width:85px">Host Name</span> -->
		<!-- 			<va-divider vertical /> -->
		<!-- 			<span style="width:200px;word-break:break-all">{{instance.host_name}}</span> -->
		<!-- 		</div> -->
		<!-- 		<div style="margin:1px;display:flex"> -->
		<!-- 			<span style="width:85px">CPU Num</span> -->
		<!-- 			<va-divider vertical /> -->
		<!-- 			<span style="width:200px">{{instance.cpu_num}}</span> -->
		<!-- 		</div> -->
		<!-- 		<div style="margin:1px;display:flex"> -->
		<!-- 			<span style="width:85px">CPU Use</span> -->
		<!-- 			<va-divider vertical /> -->
		<!-- 			<span style="width:200px">{{(instance.cpu_usage*100).toFixed(2)}}%</span> -->
		<!-- 		</div> -->
		<!-- 		<div style="margin:1px;display:flex"> -->
		<!-- 			<span style="width:85px">Mem Total</span> -->
		<!-- 			<va-divider vertical /> -->
		<!-- 			<span style="width:200px">{{instance.mem_total.toFixed(2)}}MB</span> -->
		<!-- 		</div> -->
		<!-- 		<div style="margin:1px;display:flex"> -->
		<!-- 			<span style="width:85px">Mem Use</span> -->
		<!-- 			<va-divider vertical /> -->
		<!-- 			<span style="width:200px">{{(instance.mem_usage*100).toFixed(2)}}%</span> -->
		<!-- 		</div> -->
		<!-- 		<va-divider style="width:100%"/> -->
		<!-- 		<va-button style="margin-bottom:3px" @click="get_pprof(instance.host_ip)">PPROF</va-button> -->
		<!-- 	</div> -->
		<!-- </div> -->
		<!-- <div v-if="config_proxy_instance=='instance'&&instances.length==0"> -->
		<!-- 	<div style="border:1px solid var(--va-primary);border-radius:5px;margin:5px;width:300px;height:150px;display:flex;justify-content:center;align-items:center">No Instances</div> -->
		<!-- </div> -->
	</div>
</template>
