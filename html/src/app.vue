<script setup lang="ts">
import { ref,onMounted,computed } from 'vue'
import * as appAPI from '../../api/app_browser_toc'
import * as state from './state'
import * as client from './client'

const all=computed(()=>{
	let tmp = {}
	if(!state.page.node.children){
		return tmp
	}
	for(let n of state.page.node.children){
		let pieces:string[]=n.node_name.split(".")
		if(tmp[pieces[0]]){
			tmp[pieces[0]][pieces[1]] = n
		}else{
			tmp[pieces[0]]={}
			tmp[pieces[0]][pieces[1]] = n
		}
	}
	return tmp
})
function is_json_obj(str :string):boolean{
	if(str.length<2){
		return false
	}
	if(str[0]!='{'&&str[str.length-1]!='}'){
		return false
	}
	try{
		let tmp=JSON.parse(str)
	}catch(e){
		return false
	}
	return true
}

const curg=ref<string>("")
const cura=ref<string>("")
const secret=ref<string>("")
const t_secret=ref<boolean>(false)

const key_or_proxy=ref<string>("")

const keys=ref<Map<string,KeyConfigInfo>>(new Map())
const t_keys_hover=ref<boolean>(false)

const proxys=ref<Map<string,ProxyPathInfo>>(new Map())
const t_proxys_hover=ref<boolean>(false)

const get_app_status=ref<boolean>(false)

function get_app(){
	if(curg.value==""||cura.value==""){
		keys.value=null
		proxys.value=null
		key_or_proxy.value=""
		state.set_alert("error",-2,"Group and App must be selected!")
		return
	}
	if(!state.set_load()){
		return
	}
	client.appClient.get_app({"Token":state.user.token},{g_name:curg.value,a_name:cura.value,secret:secret.value},client.timeout,(e: appAPI.Error)=>{
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
	},(resp: appAPI.GetAppResp)=>{
		if(resp.keys){
			keys.value = new Map([...resp.keys.entries()].sort())
		}else{
			keys.value = new Map()
		}
		if(resp.paths){
			proxys.value = new Map([...resp.paths.entries()].sort())
		}else{
			proxys.value = new Map()
		}
		get_app_status.value=true
		state.clear_load()
	})
}

const ing=ref<boolean>(false)
const optype=ref<string>("")

//add app
const new_g=ref<string>("")
const new_a=ref<string>("")
const new_secret=ref<string>("")

//update app secret
const update_g=ref<string>("")
const update_a=ref<string>("")
const update_old_secret=ref<string>("")
const update_new_secret=ref<string>("")

//add key config
const config_value_types=ref<string[]>(["json","raw","yaml","toml"])
const config_value_type=ref<string>("json")
const config_value=ref<string>("{\n}")
const config_key=ref<string>("")

const cur_key=ref<string>("")

//get key index config
const cur_key_index=ref<number>(0)
const cur_key_index_value=ref<string>("")
const cur_key_index_value_type=ref<string>("")

//update key config
const new_cur_key_value=ref<string>("")
const new_cur_key_value_type=ref<string>("")

//add proxy
const new_proxy_path=ref<string>("")
const new_proxy_permission_read=ref<boolean>(false)
const new_proxy_permission_write=ref<boolean>(false)
const new_proxy_permission_admin=ref<boolean>(false)
function new_proxy_permission_update(t :string){
	switch(t){
		case "read":{
			if(!new_proxy_permission_read.value){
				new_proxy_permission_write.value=false
				new_proxy_permission_admin.value=false
			}
			break
		}
		case "write":{
			if(!new_proxy_permission_write.value){
				new_proxy_permission_admin.value=false
			}else{
				new_proxy_permission_read.value=true
			}
			break
		}
		case "admin":{
			if(new_proxy_permission_admin.value){
				new_proxy_permission_read.value=true
				new_proxy_permission_write.value=true
			}
			break
		}
	}
}
function new_proxy_permission_same(proxy :string):boolean{
	if(proxys.value.get(proxy).read!=new_proxy_permission_read.value){
		return false 
	}
	if(proxys.value.get(proxy).write!=new_proxy_permission_write.value){
		return false
	}
	if(proxys.value.get(proxy).admin!=new_proxy_permission_admin.value){
		return false
	}
	return true
}

const cur_proxy=ref<string>("")
const req=ref<string>("")
const resp=ref<string>("")
const respstatus=ref<boolean>(false)

function app_op(){
	if(!state.set_load()){
		return
	}
	switch(optype.value){
		case 'del_app':{
			if(curg.value==""||cura.value==""){
				state.clear_load()
				state.set_alert("error",-2,"Group and App must be selected!")
				return
			}
			client.appClient.del_app({"Token":state.user.token},{g_name:curg.value,a_name:cura.value,secret:secret.value},client.timeout,(e: appAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: appAPI.DelAppResp)=>{
				let node = all.value[curg.value][cura.value]
				for(let i=0;i<state.page.node.children.length;i++){
					if(node == state.page.node.children[i]){
						state.page.node.children.splice(i,1)
						break
					}
				}
				curg.value=""
				cura.value=""
				secret.value=""
				keys.value=new Map()
				proxys.value=new Map()
				get_app_status.value=false
				key_or_proxy.value=""
				ing.value=false
				state.clear_load()
			})
			break
		}
		case 'add_app':{
			client.appClient.create_app({"Token":state.user.token},{project_id:state.project.cur_id,g_name:new_g.value,a_name:new_a.value,secret:new_secret.value},client.timeout,(e: appAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: appAPI.CreateAppResp)=>{
				if(all.value[new_g.value] && all.value[new_g.value][new_a.value]){
					return
				}
				if(state.page.node.children){
					state.page.node.children.push({
						node_id:resp.node_id,
						node_name:new_g.value+"."+new_a.value,
						node_data:"",
						canread:true,
						canwrite:true,
						admin:true,
						children:[],
					})
				}else{
					state.page.node.children=[{
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
				ing.value=false
				state.clear_load()
			})
			break
		}
		case 'update_secret':{
			let req = {
				g_name:update_g.value,
				a_name:update_a.value,
				old_secret:update_old_secret.value,
				new_secret:update_new_secret.value
			}
			client.appClient.update_app_secret({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: appAPI.UpdateAppSecretResp)=>{
				update_g.value=""
				update_a.value=""
				update_old_secret.value=""
				update_new_secret.value=""
				ing.value=false
				state.clear_load()
			})
			break
		}
		case 'get_key':{
			let req = {
				g_name:curg.value,
				a_name:cura.value,
				secret:secret.value,
				key:cur_key.value,
				index:cur_key_index.value,
			}
			client.appClient.get_key_config({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: appAPI.GetKeyConfigResp)=>{
				cur_key_index_value.value=resp.value
				cur_key_index_value_type.value=resp.value_type
				state.clear_load()
			})
			break
		}
		case 'add_key':{
			if(keys.value.has(config_key.value)){
				state.clear_load()
				state.set_alert("error",-2,"key already exist")
				break
			}
			let req = {
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
			},(resp: appAPI.SetKeyConfigResp)=>{
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
				g_name:curg.value,
				a_name:cura.value,
				secret:secret.value,
				key:cur_key.value,
				value:new_cur_key_value.value,
				value_type:new_cur_key_value_type.value,
				new_key:false,
			}
			client.appClient.set_key_config({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: appAPI.SetKeyConfigResp)=>{
				new_cur_key_value_type.value=''
				ing.value=false
				state.clear_load()
				get_app()
			})
			break
		}
		case 'rollback_key':{
			let req = {
				g_name:curg.value,
				a_name:cura.value,
				secret:secret.value,
				key:cur_key.value,
				index:cur_key_index.value,
			}
			client.appClient.rollback({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: appAPI.RollbackResp)=>{
				cur_key_index.value=0
				ing.value=false
				state.clear_load()
				get_app()
			})
			break
		}
		case 'del_key':{
			let req = {
				g_name:curg.value,
				a_name:cura.value,
				secret:secret.value,
				key:cur_key.value,
			}
			client.appClient.del_key({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: appAPI.DelKeyResp)=>{
				keys.value.delete(cur_key.value)
				cur_key.value=""
				ing.value=false
				state.clear_load()
			})
			break
		}
		case 'add_proxy':{
			if(proxys.value.has(new_proxy_path.value)){
				state.clear_load()
				state.set_alert("error",-2,"proxy path already exist")
				break
			}
			let req = {
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
				g_name:curg.value,
				a_name:cura.value,
				secret:secret.value,
				path:cur_proxy.value,
				read:new_proxy_permission_read.value,
				write:new_proxy_permission_write.value,
				admin:new_proxy_permission_admin.value,
				new_path:false,
			}
			client.appClient.set_proxy({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: appAPI.SetProxyResp)=>{
				proxys.value.get(cur_proxy.value).read=new_proxy_permission_read.value
				proxys.value.get(cur_proxy.value).write=new_proxy_permission_write.value
				proxys.value.get(cur_proxy.value).admin=new_proxy_permission_admin.value
				ing.value=false
				state.clear_load()
			})
			break
		}
		case 'del_proxy':{
			let req = {
				g_name:curg.value,
				a_name:cura.value,
				secret:secret.value,
				path:cur_proxy.value,
			}
			client.appClient.del_proxy({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: appAPI.DelProxyResp)=>{
				proxys.value.delete(cur_proxy.value)
				cur_proxy.value=''
				ing.value=false
				state.clear_load()
			})
			break
		}
		case 'proxy':{
			let req = {
				g_name:curg.value,
				a_name:cura.value,
				path:cur_proxy.value,
				data:req.value,
			}
			client.appClient.proxy({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: appAPI.ProxyResp)=>{
				resp.value=resp.data
				respstatus.value=true
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
</script>
<template>
	<va-modal v-model="ing" attach-element="#app" max-width="1000px" max-height="600px" hide-default-actions no-dismiss overlay-opacity="0.2" z-index="999">
		<template #default>
			<div v-if="optype=='del_app'" style="display:flex;flex-direction:column">
				<va-card color="primary" gradient style="margin:0 0 5px 0">
					<va-card-title>Warning</va-card-title>
					<va-card-content>
						<p>You are deleting app: {{ cura }} in group: {{ curg }}.</p>
						<p>All data in this app will be deleted.</p>
						<p>Please confirm!</p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="app_op" gradient>Del</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='add_app'" style="display:flex;flex-direction:column">
				<va-input type="text" label="Group*" style="margin:2px" v-model.trim="new_g" @keyup.enter="()=>{if(new_g!=''&&new_a!=''&&(!Boolean(all[new_g])||!Boolean(all[new_g][new_a]))){app_op()}}" />
				<va-input type="text" label="App*" style="margin:2px" v-model.trim="new_a" @keyup.enter="()=>{if(new_g!=''&&new_a!=''&&(!Boolean(all[new_g])||!Boolean(all[new_g][new_a]))){app_op()}}" />
				<va-input type="text" label="Secret" style="margin:2px" v-model.trim="new_secret" :max-length="31" @keyup.enter="()=>{if(new_g!=''&&new_a!=''&&(!Boolean(all[new_g])||!Boolean(all[new_g][new_a]))){app_op()}}" />
				<div style="display:flex;justify-content:center">
					<va-button @click="app_op" style="margin:5px" :disabled="new_g==''||new_a==''||(Boolean(all[new_g])&&Boolean(all[new_g][new_a]))">Add</va-button>
					<va-button @click="new_g='';new_a='';new_secret='';ing=false" style="margin:5px">Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='update_secret'" style="display:flex;flex-direction:column">
				<div>
					<va-select 
					trigger="hover"
					dropdown-icon=""
					label="Group*"
					:options="Object.keys(all)"
					style="width:198px;margin:2px"
					v-model="update_g"
					no-options-text="No Groups"
					>
						<template #option='{option,index,selectOption}'>
							<va-hover
							stateful
							@click="()=>{
								if(option!=update_g){
									selectOption(option)
									update_a=''
									update_old_secret=''
									update_new_secret=''
								}
							}"
							>
								<template #default="{hover}">
									<div
									style="padding:10px;cursor:pointer"
									:style="{'background-color':hover?'var(--va-background-border)':'',color:hover||update_g==option?'var(--va-primary)':'black'}"
									>
										{{option}}
									</div>
								</template>
							</va-hover>
						</template>
					</va-select>
					<va-select
					trigger="hover"
					dropdown-icon=""
					label="App*"
					:options="update_g==''?[]:Object.keys(all[update_g])"
					style="width:198px;margin:2px"
					v-model="update_a"
					no-options-text="No Apps"
					>
						<template #option='{option,index,selectOption}'>
							<va-hover
							stateful
							@click="()=>{
								if(option!=update_a){
									selectOption(option)
									update_old_secret=''
									update_new_secret=''
								}
							}"
							>
								<template #default="{hover}">
									<div
									style="padding:10px;cursor:pointer"
									:style="{'background-color':hover?'var(--va-background-border)':'',color:hover||update_a==option?'var(--va-primary)':'black'}"
									>
										{{option}}
									</div>
								</template>
							</va-hover>
						</template>
					</va-select>
				</div>
				<va-input type="text" label="Old Secret" style="width:400px;margin:2px" v-model.trim="update_old_secret" @keyup.enter="()=>{if(update_g!=''&&update_a!=''&&update_new_secret!=update_old_secret){app_op()}}" />
				<va-input type="text" label="New Secret" style="width:400px;margin:2px" v-model.trim="update_new_secret" @keyup.enter="()=>{if(update_g!=''&&update_a!=''&&update_new_secret!=update_old_secret){app_op()}}" />
				<div style="display:flex;justify-content:center">
					<va-button @click="app_op" style="margin:5px" :disabled="update_g==''||update_a==''||update_old_secret==update_new_secret">Update</va-button>
					<va-button @click="update_g='';update_a='';update_old_secret='';update_new_secret='';ing=false" style="margin:5px">Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='del_key'" style="display:flex;flex-direction:column">
				<va-card color="primary" gradient style="margin:0 0 5px 0">
					<va-card-title>Warning</va-card-title>
					<va-card-content>
						<p>You are deleting config key: {{ cur_key }}.</p>
						<p>All data in this key will be deleted.</p>
						<p>Please confirm!</p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="app_op" gradient>Del</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='add_key'" style="display:flex;flex-direction:column">
				<va-input type="text" label="Key_Name*" style="margin:1px;width:800px" v-model.trim="config_key" />
				<div style="display:flex;justify-content:space-evenly;align-items:center">
					<va-radio v-for="(option,index) in config_value_types" :key="index" :option="option" v-model="config_value_type" style="margin:4px" disabled />
				</div>
				<va-input type="textarea" label="Content" style="margin:1px;width:800px" :min-rows="15" :max-rows="15" v-model.trim="config_value" />
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="app_op" gradient :disabled="config_key==''||keys.has(config_key)">Add</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="config_key='';config_value='{\n}';config_value_type='json';ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='update_key'">
				<va-card color="primary" gradient style="margin:0 0 5px 0">
					<va-card-title>Warning</va-card-title>
					<va-card-content>
						<p>You are updating config key: {{ cur_key }}.</p>
						<p>Data in this key will be updated.</p>
						<p>Please confirm!</p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="app_op" gradient>Update</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='rollback_key'">
				<va-card color="primary" gradient style="margin:0 0 5px 0">
					<va-card-title>Warning</va-card-title>
					<va-card-content>
						<p>You are rollbacking config key: {{ cur_key }} to config id: {{ cur_key_index }}.</p>
						<p>Data in this key will be updated.</p>
						<p>Please confirm!</p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="app_op" gradient>Rollback</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='add_proxy'">
				<va-input label="Path" style="width:500px" v-model.trim="new_proxy_path"/>
				<div style="display:flex;justify-content:space-around;margin:4px">
					<va-switch v-model="new_proxy_permission_read" true-inner-label="Read" false-inner-label="Read" @update:model-value="new_proxy_permission_update('read')" />
					<va-switch v-model="new_proxy_permission_write" true-inner-label="Write" false-inner-label="Write" @update:model-value="new_proxy_permission_update('write')" />
					<va-switch v-model="new_proxy_permission_admin" true-inner-label="Admin" false-inner-label="Admin" @update:model-value="new_proxy_permission_update('admin')" />
				</div>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="app_op" gradient :disabled="new_proxy_path==''||proxys.has(new_proxy_path)">Add</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="new_proxy_path='';new_proxy_permission_read=false;new_proxy_permission_write=false;new_proxy_permission_admin=false;ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='update_proxy'">
				<va-card color="primary" gradient style="margin:0 0 5px 0">
					<va-card-title>Warning</va-card-title>
					<va-card-content>
						<p>You are updating proxy path: {{ cur_proxy }}.</p>
						<p>Permission required on this path will be updated.</p>
						<p>Please confirm!</p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="app_op" gradient>Update</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='del_proxy'">
				<va-card color="primary" gradient style="margin:0 0 5px 0">
					<va-card-title>Warning</va-card-title>
					<va-card-content>
						<p>You are deleting proxy path: {{ cur_proxy }}.</p>
						<p>Please confirm!</p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="app_op" gradient>Del</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='proxy'">
				<va-card color="primary" gradient style="margin:0 0 5px 0">
					<va-card-title>Warning</va-card-title>
					<va-card-content>
						<p>You are requesting path: {{ cur_proxy }}.</p>
						<p>This request may cause changes in server data.</p>
						<p>Please confirm!</p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="app_op" gradient>Proxy</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="cur_proxy='';ing=false" gradient>Cancel</va-button>
				</div>
			</div>
		</template>
	</va-modal>
	<div style="flex:1;display:flex;flex-direction:column;margin:1px;overflow-y:auto">
		<div style="display:flex;margin:1px 0;align-self:center">
			<va-select
				dropdown-icon=""
				outline
				trigger="hover"
				label="Group*"
				no-options-text="No Groups"
				:options="Object.keys(all)"
				v-model="curg"
				style="width:150px;margin-right:1px"
			>
				<template #option='{option,index,selectOption}'>
					<va-hover
						stateful
						@click="()=>{
							if(curg!=option){
								selectOption(option)
								cura=''
								secret=''
								keys=new Map()
								proxys=new Map()
								get_app_status=false
								key_or_proxy=''
							}
						}"
					>
						<template #default='{hover}'>
							<div
								style="padding:10px;cursor:pointer"
								:style="{'background-color':hover?'var(--va-background-border)':'',color:hover||curg==option?'var(--va-primary)':'black'}"
							>
								{{option}}
							</div>
						</template>
					</va-hover>
				</template>
			</va-select>
			<va-select
				dropdown-icon=""
				outline
				trigger="hover"
				label="App*"
				no-options-text="No Apps"
				:options="curg==''?[]:Object.keys(all[curg])"
				v-model="cura"
				style="width:150px;margin:0 1px"
			>
				<template #option='{option,index,selectOption}'>
					<va-hover
						stateful
						@click="()=>{
							if(cura!=option){
								selectOption(option)
								secret=''
								keys=new Map()
								proxys=new Map()
								get_app_status=false
								key_or_proxy=''
							}
						}"
					>
						<template #default='{hover}'>
							<div
								style="padding:10px;cursor:pointer"
								:style="{'background-color':hover?'var(--va-background-border)':'',color:hover||cura==option?'var(--va-primary)':'black'}"
							>
								{{option}}
							</div>
						</template>
					</va-hover>
				</template>
			</va-select>
			<va-input :type="t_secret?'text':'password'" v-model.trim="secret" outline label="Secret" :max-length="31" style="width:250px;margin:0 1px" @keyup.enter="()=>{if(curg!=''&&cura!=''){get_app()}}">
				<template #appendInner>
					<va-icon :name="t_secret?'◎':'◉'" size="small" color="var(--va-primary)" @click="t_secret=!t_secret" />
				</template>
			</va-input>
			<va-button style="margin:0 2px" :disabled="curg==''||cura==''" @click="get_app">Search</va-button>
			<va-dropdown  v-if="state.page.node.admin" trigger="hover" style="width:36px;margin-right:4px">
				<template #anchor>
					<va-button>•••</va-button>
				</template>
				<va-dropdown-content>
					<va-popover message="Create New App" :hover-out-timeout="0" :hover-over-timeout="0" color="primary" prevent-overflow>
						<va-button style="width:36px;margin:0 3px" @click="optype='add_app';ing=true">+</va-button>
					</va-popover>
					<va-popover message="Update Add Secret" :hover-out-timeout="0" :hover-over-timeout="0" color="primary" prevent-overflow>
						<va-button style="width:36px;margin:0 3px" @click="optype='update_secret';ing=true">◉</va-button>
					</va-popover>
					<va-popover message="Delete App" :hover-out-timeout="0" :hover-over-timeout="0" color="primary" prevent-overflow>
						<va-button style="width:36px;margin:0 3px" :disabled="curg==''||cura==''" @click="optype='del_app';ing=true">x</va-button>
					</va-popover>
				</va-dropdown-content>
			</va-dropdown>
		</div>
		<!-- configs -->
		<div
			v-if="get_app_status&&key_or_proxy!='proxy'"
			style="display:flex;align-items:center;margin:1px 0;cursor:pointer"
			:style="{'background-color':t_keys_hover?'var(--va-shadow)':'var(--va-background-element)'}"
			@click="()=>{
				if(key_or_proxy==''){
					key_or_proxy='key'
					cur_key=''
				}else{
					key_or_proxy=''
				}
			}"
			@mouseover="t_keys_hover=true"
			@mouseout="t_keys_hover=false"
		>
			<span style="flex:1;padding:12px;color:var(--va-primary)">Configs</span>
			<va-button
				v-if="all[curg][cura].canwrite||all[curg][cura].admin"
				style="height:30px"
				size="small"
				@mouseover.stop=""
				@mouseout.stop=""
				@click.stop="optype='add_key';ing=true"
			>
				ADD
			</va-button>
			<span style="width:60px;padding:12px 20px;color:var(--va-primary)">{{ key_or_proxy?'▲':'▼' }}</span>
		</div>
		<!-- keys -->
		<div v-if="key_or_proxy=='key'&&keys.size>0" style="overflow-y:auto;flex:1;display:flex;flex-direction:column">
			<template v-for="key of keys.keys()">
				<div
					v-if="cur_key==''||cur_key==key"
					style="cursor:pointer;display:flex;align-items:center;margin:1px 10px"
					:style="{'background-color':keys.get(key).hover?'var(--va-shadow)':'var(--va-background-element)'}"
					@click="()=>{
						if(cur_key==''){
							cur_key=key
							cur_key_index=0
							new_cur_key_value_type=''
						}else{
							cur_key=''
						}
					}"
					@mouseover="keys.get(key).hover=true"
					@mouseout="keys.get(key).hover=false"
				>
					<span style="width:35px;padding:12px;color:var(--va-primary)"> {{ cur_key!=''&&cur_key==key?'▼':'►' }} </span>
					<span style="padding:12px;color:var(--va-primary)">{{key}}</span>
				</div>
				<div v-if="cur_key==key" style="flex:1;display:flex;margin:1px 20px;overflow-y:auto">
					<div style="flex:1;display:flex;flex-direction:column;overflow-y:auto">
						<textarea readonly style="border:0px;flex:1;resize:none;background-color:var(--va-background-element);padding:10px 20px">{{JSON.stringify(JSON.parse(keys.get(key).cur_value),null,4)}}</textarea>
						<div style="align-self:center;display:flex;align-items:center">
							<b style="color:var(--va-primary);margin-right:10px">Current Config ID:  {{ keys.get(key).cur_index }}</b>
							<va-dropdown trigger="hover" prevent-overflow placement="top">
								<template #anchor>
									<va-button style="width:60px;height:30px;margin:2px" size="small">History</va-button>
								</template>
								<va-dropdown-content>
									<div style="max-height:300px;overflow-y:auto;display:flex;flex-direction:column">
										<va-button
											v-for="index of keys.get(key).max_index"
											size="small"
											style="height:24px;width:42px;padding:5px 0;margin:1px;cursor:pointer"
											:disabled="keys.get(key).cur_index==keys.get(key).max_index-index+1"
											@click="cur_key_index=keys.get(key).max_index-index+1;optype='get_key';app_op()"
										>
											{{keys.get(key).max_index-index+1}}
										</va-button>
									</div>
								</va-dropdown-content>
							</va-dropdown>
							<va-button
								v-if="all[curg][cura].canwrite||all[curg][cura].admin"
								size="small"
								style="width:60px;height:30px;margin:2px"
								:disabled="Boolean(keys.get(key).new_cur_value)"
								@click="()=>{
									if(keys.get(key).cur_value){
										new_cur_key_value=JSON.stringify(JSON.parse(keys.get(key).cur_value),null,4)
										new_cur_key_value_type=keys.get(key).cur_value_type
									}else{
										new_cur_key_value='{\n}'
										new_cur_key_value_type='json'
									}
								}"
							>
								Edit
							</va-button>
							<va-button
								v-if="all[curg][cura].canwrite||all[curg][cura].admin"
								size="small"
								style="width:60px;height:30px;margin:2px"
								@click.stop="optype='del_key';ing=true"
							>
								Del
							</va-button>
						</div>
					</div>
					<va-divider v-if="cur_key_index!=0||new_cur_key_value_type!=''" vertical style="margin:0 4px" />
					<div v-if="cur_key_index!=0||new_cur_key_value_type!=''" style="flex:1;display:flex;flex-direction:column">
						<textarea v-if="cur_key_index!=0" style="border:0;flex:1;resize:none;background-color:var(--va-background-element);padding:10px 20px" readonly >{{JSON.stringify(JSON.parse(cur_key_index_value),null,4)}}</textarea>
						<textarea v-if="new_cur_key_value_type!=''" style="border:0;flex:1;resize:none;background-color:var(--va-background-element);padding:10px 20px" v-model.trim="new_cur_key_value"></textarea>
						<div style="display:flex;align-items:center">
							<va-radio
								v-if="cur_key_index!=0"
								v-for="(option,index) in config_value_types"
								:key="index"
								:option="option"
								v-model.trim="cur_key_index_value_type"
								style="margin:4px"
								disabled
							/>
							<va-radio
								v-if="new_cur_key_value_type!=''"
								v-for="(option,index) in config_value_types"
								:key="index"
								:option="option"
								v-model.trim="new_cur_key_value_type"
								style="margin:4px"
								disabled
							/>
							<span style="flex:1"></span>
							<va-button
								v-if="cur_key_index!=0&&all[curg][cura].canwrite&&all[curg][cura].admin"
								size="small"
								style="width:60px;height:30px;margin:2px"
								@click="optype='rollback_key';ing=true"
							>
								Rollback
							</va-button>
							<va-button
								v-if="new_cur_key_value_type!=''"
								:disabled="!is_json_obj(new_cur_key_value)||JSON.stringify(JSON.parse(new_cur_key_value),null,4)==JSON.stringify(JSON.parse(keys.get(key).cur_value),null,4)"
								size="small"
								style="width:60px;height:30px;margin:2px"
								@click="optype='update_key';ing=true"
							>
								Update
							</va-button>
							<va-button size="small" style="width:60px;height:30px;margin:2px" @click="cur_key_index=0;new_cur_key_value_type=''">Cancel</va-button>
						</div>
					</div>
				</div>
			</template>
		</div>
		<div v-if="key_or_proxy=='key'&&keys.size==0">
			<div style="margin:1px 10px;padding:12px;display:flex;flex-direction:column;background-color:var(--va-background-element);color:var(--va-primary)">No Config Keys</div>
		</div>
		<!-- proxys -->
		<div
			v-if="get_app_status&&key_or_proxy!='key'" 
			style="display:flex;align-items:center;margin:1px 0;cursor:pointer"
			:style="{'background-color':t_proxys_hover?'var(--va-shadow)':'var(--va-background-element)'}"
			@click="()=>{
				if(key_or_proxy==''){
					key_or_proxy='proxy'
					cur_proxy=''
				}else{
					key_or_proxy=''
				}
			}"
			@mouseover="t_proxys_hover=true"
			@mouseout="t_proxys_hover=false"
		>
			<span style="flex:1;padding:12px;color:var(--va-primary)">Proxys</span>
			<va-button
				v-if="all[curg][cura].canwrite||all[curg][cura].admin"
				style="height:30px"
				size="small"
				@mouseover.stop=""
				@mouseout.stop=""
				@click.stop="new_proxy_permission_read=false;new_proxy_permission_write=false;new_proxy_permission_admin=false;optype='add_proxy';ing=true"
			>
				ADD
			</va-button>
			<span style="width:60px;padding:12px 20px;color:var(--va-primary)">{{ key_or_proxy?'▲':'▼' }}</span>
		</div>
		<!-- paths -->
		<div v-if="key_or_proxy=='proxy'&&proxys.size>0" style="overflow-y:auto;flex:1;display:flex;flex-direction:column">
			<template v-for="proxy of proxys.keys()">
				<div
					v-if="cur_proxy==''||cur_proxy==proxy"
					style="cursor:pointer;display:flex;align-items:center;margin:1px 10px"
					:style="{'background-color':proxys.get(proxy).hover?'var(--va-shadow)':'var(--va-background-element)'}"
					@mouseover="proxys.get(proxy).hover=true"
					@mouseout="proxys.get(proxy).hover=false"
					@click="()=>{
						if(cur_proxy==''){
							cur_proxy=proxy
							req='{\n}'
							new_proxy_permission_read=proxys.get(proxy).read
							new_proxy_permission_write=proxys.get(proxy).write
							new_proxy_permission_admin=proxys.get(proxy).admin
							respstatus=false
						}else{
							cur_proxy=''
						}
					}"
				>
					<span style="width:35px;padding:12px;color:var(--va-primary)"> {{ cur_proxy!=''&&cur_proxy==proxy?'▼':'►' }} </span>
					<span style="padding:12px;color:var(--va-primary)">{{proxy}}</span>
				</div>
				<div v-if="cur_proxy==proxy" style="flex:1;display:flex;margin:1px 20px;overflow-y:auto">
					<div style="flex:1;display:flex;flex-direction:column">
						<textarea style="border:0px;flex:1;resize:none;background-color:var(--va-background-element);padding:10px 20px" v-model.trim="req" :readonly="respstatus"></textarea>
						<div style="width:100%;display:flex">
							<va-button style="width:60px;height:30px;margin:2px 0" size="small" @click="optype='proxy';ing=true" :disabled="respstatus||!is_json_obj(req)">Proxy</va-button>
							<div style="flex:1"></div>
							<va-switch
								:disabled="respstatus||!all[curg][cura].canwrite||!all[curg][cura].admin"
								style="margin:2px"
								v-model="new_proxy_permission_read"
								size="small" true-inner-label="Read"
								false-inner-label="Read"
								@update:model-value="new_proxy_permission_update('read')"
							/>
							<va-switch
								:disabled="respstatus||!all[curg][cura].canwrite||!all[curg][cura].admin"
								style="margin:2px"
								v-model="new_proxy_permission_write"
								size="small"
								true-inner-label="Write"
								false-inner-label="Write"
								@update:model-value="new_proxy_permission_update('write')"
							/>
							<va-switch
								:disabled="respstatus||!all[curg][cura].canwrite||!all[curg][cura].admin"
								style="margin:2px"
								v-model="new_proxy_permission_admin"
								size="small"
								true-inner-label="Admin"
								false-inner-label="Admin"
								@update:model-value="new_proxy_permission_update('admin')"
							/>
							<va-button
								v-if="all[curg][cura].canwrite||all[curg][cura].admin"
								style="width:60px;height:30px;margin:2px"
								size="small"
								:disabled="respstatus||new_proxy_permission_same(proxy)"
								@click="optype='update_proxy';ing=true"
							>
								Update
							</va-button>
							<va-button
								v-if="all[curg][cura].canwrite||all[curg][cura].admin"
								:disabled="respstatus"
								size="small"
								style="width:60px;height:30px;margin:2px"
								@click.stop="optype='del_proxy';ing=true"
							>
								DEL
							</va-button>
						</div>
					</div>
					<va-divider v-if="respstatus" vertical style="margin:0 4px" />
					<div v-if="respstatus" style="flex:1;display:flex;flex-direction:column">
						<textarea style="border:0px;flex:1;resize:none;background-color:var(--va-background-element);padding:10px 20px" readonly>{{is_json_obj(resp)?JSON.stringify(JSON.parse(resp),null,4):resp}}</textarea>
						<va-button style="align-self:center;width:60px;height:30px;margin:2px" size="small" @click="respstatus=false">OK</va-button>
					</div>
				</div>
			</template>
		</div>
		<div v-if="key_or_proxy=='proxy'&&proxys.size==0">
			<div style="margin:1px 10px;padding:12px;display:flex;flex-direction:column;background-color:var(--va-background-element);color:var(--va-primary)">No Proxy Paths</div>
		</div>
	</div>
</template>
