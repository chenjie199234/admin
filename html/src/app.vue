<script setup lang="ts">
import { ref,onMounted } from 'vue'
import * as appAPI from '../../api/app_browser_toc'
import * as state from './state'
import * as client from './client'

const all=ref<{[k:string]:string[]}>({})
onMounted(()=>{
	if(!state.page.node.children){
		return
	}
	for(let n of state.page.node.children){
		let pieces:string[] = n.node_name.split(".")
		if(all.value[pieces[0]]){
			all.value[pieces[0]].push(pieces[1])
		}else{
			all.value[pieces[0]]=[pieces[1]]
		}
	}
})
const curg=ref<string>("")
const cura=ref<string>("")
const secret=ref<string>("")
const t_secret=ref<boolean>(false)

const keys=ref<Map<string,KeyConfigInfo>>(new Map())
const t_keys=ref<boolean>(false)
const t_keys_hover=ref<boolean>(false)

const proxys=ref<Map<string,ProxyPathInfo>>(new Map())
const t_proxys=ref<boolean>(false)
const t_proxys_hover=ref<boolean>(false)

function get_app(need_set_load: boolean){
	if(curg.value==""||cura.value==""){
		keys.value=null
		proxys.value=null
		state.set_alert("error",-2,"Group and App must be selected!")
		return
	}
	if(need_set_load){
		if(!state.set_load()){
			return
		}
	}
	client.appClient.get_app({"Token":state.user.token},{g_name:curg.value,a_name:cura.value,secret:secret.value},client.timeout,(e: appAPI.Error)=>{
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
	},(resp: appAPI.GetAppResp)=>{
		state.clear_load()
		if(resp.keys){
			keys.value = new Map([...resp.keys.entries()].sort())
		}else{
			keys.value = new Map()
		}
		t_keys.value=true
		if(resp.paths){
			proxys.value = new Map([...resp.paths.entries()].sort())
		}else{
			proxys.value = new Map()
		}
		t_proxys.value=true
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

const cur_key=ref<string>("")

//add key config
const config_value_types=ref<string[]>(["json","raw","yaml","toml"])
const config_value_type=ref<string>("json")
const config_value=ref<string>("{\n}")
const config_key=ref<string>("")

//get key index config
const cur_key_index=ref<number>(0)
const cur_key_index_value=ref<string>("")
const cur_key_index_value_type=ref<string>("")

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

//update proxy
function cur_proxy_permission_update(proxy:string,t:string){
	switch(t){
		case "read":{
			if(!proxys.value.get(proxy).new_read){
				proxys.value.get(proxy).new_write=false
				proxys.value.get(proxy).new_admin=false
			}
			break
		}
		case "write":{
			if(!proxys.value.get(proxy).new_write){
				proxys.value.get(proxy).new_admin=false
			}else{
				proxys.value.get(proxy).new_read=true
			}
			break
		}
		case "admin":{
			if(proxys.value.get(proxy).new_admin){
				proxys.value.get(proxy).new_read=true
				proxys.value.get(proxy).new_write=true
			}
			break
		}
	}
}

const cur_proxy=ref<string>("")

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
				state.clear_load()
				let index=all.value[curg.value].indexOf(cura.value)
				if(index!=-1){
					all.value[curg.value].splice(index,1)
				}
				if(all.value[curg.value].length==0){
					delete all.value[curg.value]
				}
				curg.value=""
				cura.value=""
				secret.value=""
				keys.value=new Map()
				proxys.value=new Map()
				ing.value=false
			})
			break
		}
		case 'add_app':{
			client.appClient.create_app({"Token":state.user.token},{project_id:state.project.cur_id,g_name:new_g.value,a_name:new_a.value,secret:new_secret.value},client.timeout,(e: appAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: appAPI.CreateAppResp)=>{
				state.clear_load()
				if(all.value[new_g.value]){
					all.value[new_g.value].push(new_a.value)
				}else{
					all.value[new_g.value]=[new_a.value]
				}
				new_g.value=""
				new_a.value=""
				new_secret.value=""
				ing.value=false
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
				state.clear_load()
				update_g.value=""
				update_a.value=""
				update_old_secret.value=""
				update_new_secret.value=""
				ing.value=false
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
				state.clear_load()
				cur_key_index_value.value=resp.value
				cur_key_index_value_type.value=resp.value_type
				ing.value=true
			})
			break
		}
		case 'add_key':{
			let req = {
				g_name:curg.value,
				a_name:cura.value,
				secret:secret.value,
				key:config_key.value,
				value:config_value.value,
				value_type:config_value_type.value,
			}
			client.appClient.set_key_config({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: appAPI.SetKeyConfigResp)=>{
				config_key.value = ''
				config_value.value = '{}'
				config_value_type.value = 'json'
				ing.value=false
				get_app(false)
			})
			break
		}
		case 'update_key':{
			let req = {
				g_name:curg.value,
				a_name:cura.value,
				secret:secret.value,
				key:cur_key.value,
				value:keys.value.get(cur_key.value).new_cur_value,
				value_type:keys.value.get(cur_key.value).new_cur_value_type,
			}
			client.appClient.set_key_config({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: appAPI.SetKeyConfigResp)=>{
				ing.value=false
				get_app(false)
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
				ing.value=false
				get_app(false)
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
				ing.value=false
				get_app(false)
			})
			break
		}
		case 'add_proxy':{
			let req = {
				g_name:curg.value,
				a_name:cura.value,
				secret:secret.value,
				path:new_proxy_path.value,
				read:new_proxy_permission_read.value,
				write:new_proxy_permission_write.value,
				admin:new_proxy_permission_admin.value,
			}
			client.appClient.set_proxy({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: appAPI.SetProxyResp)=>{
				ing.value=false
				get_app(false)
			})
			break
		}
		case 'update_proxy':{
			let req = {
				g_name:curg.value,
				a_name:cura.value,
				secret:secret.value,
				path:cur_proxy.value,
				read:proxys.value.get(cur_proxy.value).new_read,
				write:proxys.value.get(cur_proxy.value).new_write,
				admin:proxys.value.get(cur_proxy.value).new_admin,
			}
			client.appClient.set_proxy({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: appAPI.SetProxyResp)=>{
				ing.value=false
				get_app(false)
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
				ing.value=false
				get_app(false)
			})
			break
		}
		case 'proxy':{
			let req = {
				g_name:curg.value,
				a_name:cura.value,
				path:cur_proxy.value,
				data:proxys.value.get(cur_proxy.value).req,
			}
			client.appClient.proxy({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: appAPI.ProxyResp)=>{
				state.clear_load()
				proxys.value.get(cur_proxy.value).resp=resp.data
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
					<va-button style="width:80px;margin:5px 10px 0 0" @click="ing=false;app_op()" gradient>Del</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='add_app'" style="display:flex;flex-direction:column">
				<va-input type="text" label="Group*" style="margin:2px" v-model.trim="new_g" @keyup.enter="()=>{if(new_g!=''&&new_a!=''){app_op()}}" />
				<va-input type="text" label="App*" style="margin:2px" v-model.trim="new_a" @keyup.enter="()=>{if(new_g!=''&&new_a!=''){app_op()}}" />
				<va-input type="text" label="Secret" style="margin:2px" v-model.trim="new_secret" :max-length="31" @keyup.enter="()=>{if(new_g!=''&&new_a!=''){app_op()}}" />
				<div style="display:flex;justify-content:center">
					<va-button @click="app_op" style="margin:5px" :disabled="new_g==''||new_a==''">Add</va-button>
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
					:options="all[update_g]"
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
					<va-button style="width:80px;margin:5px 10px 0 0" @click="ing=false;app_op()" gradient>Del</va-button>
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
					<va-button style="width:80px;margin:5px 10px 0 0" @click="app_op" gradient :disabled="config_key==''" >Add</va-button>
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
					<va-button style="width:80px;margin:5px 10px 0 0" @click="ing=false;app_op()" gradient>Update</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='get_key'">
				<va-input type="textarea" :label="'Content Type:'+cur_key_index_value_type" style="margin:1px;width:800px" :model-value="JSON.stringify(JSON.parse(cur_key_index_value),null,4)" readonly :min-rows="15" :max-rows="15" />
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="optype='rollback_key';app_op()" gradient>Rollback</va-button>
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
					<va-button style="width:80px;margin:5px 10px 0 0" gradient :disabled="new_proxy_path==''" @click="app_op">Add</va-button>
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
					<va-button style="width:80px;margin:5px 10px 0 0" @click="ing=false;app_op()" gradient>Update</va-button>
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
					<va-button style="width:80px;margin:5px 10px 0 0" @click="ing=false;app_op()" gradient>Del</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype='proxy'">
				<va-card color="primary" gradient style="margin:0 0 5px 0">
					<va-card-title>Warning</va-card-title>
					<va-card-content>
						<p>You are requesting path: {{ cur_proxy }}.</p>
						<p>This request may cause changes in server data.</p>
						<p>Please confirm!</p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="ing=false;app_op()" gradient>Proxy</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
		</template>
	</va-modal>
	<div style="flex:1;display:flex;flex-direction:column;margin:1px;width:100%;overflow-y:auto">
		<div style="width:100%;display:flex;margin:1px 0">
			<va-select
			dropdown-icon=""
			outline
			trigger="hover"
			label="Group*"
			no-options-text="No Groups"
			:options="Object.keys(all)"
			v-model="curg"
			style="width:250px;margin-right:1px"
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
			:options="all[curg]"
			v-model="cura"
			style="width:250px;margin:0 1px"
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
			<va-input :type="t_secret?'text':'password'" v-model.trim="secret" outline label="Secret" :max-length="31" style="min-width:250px;max-width:250px;margin:0 1px" @keyup.enter="()=>{if(curg!=''&&cura!=''){get_app(true)}}">
				<template #appendInner>
					<va-icon :name="t_secret?'◎':'◉'" size="small" color="var(--va-primary)" @click="t_secret=!t_secret" />
				</template>
			</va-input>
			<va-button style="margin:0 2px" :disabled="curg==''||cura==''" @click="get_app(true)">Search</va-button>
			<va-dropdown  v-if="state.page.node.admin" trigger="hover" style="width:36px;margin-right:4px">
				<template #anchor>
					<va-button>•••</va-button>
				</template>
				<va-dropdown-content>
					<va-popover message="Create New App" :hover-out-timeout="0" :hover-over-timeout="0" color="primary">
						<va-button style="width:36px;margin:0 3px" @click="ing=true;optype='add_app'">+</va-button>
					</va-popover>
					<va-popover message="Update Add Secret" :hover-out-timeout="0" :hover-over-timeout="0" color="primary">
						<va-button style="width:36px;margin:0 3px" @click="ing=true;optype='update_secret'">◉</va-button>
					</va-popover>
					<va-popover message="Delete App" :hover-out-timeout="0" :hover-over-timeout="0" color="primary">
						<va-button style="width:36px;margin:0 3px" :disabled="curg==''||cura==''" @click="optype='del_app';ing=true">x</va-button>
					</va-popover>
				</va-dropdown-content>
			</va-dropdown>
		</div>
		<!-- configs -->
		<div
			style="width:100%;display:flex;align-items:center;margin:1px 0;cursor:pointer"
			:style="{'background-color':t_keys_hover?'var(--va-shadow)':'var(--va-background-element)'}"
			@click="t_keys=!t_keys"
			@mouseover="t_keys_hover=true"
			@mouseout="t_keys_hover=false"
		>
			<span style="flex:1;padding:12px;color:var(--va-primary)">Configs</span>
			<va-button style="height:30px" size="small" :disabled="curg==''||cura==''" @mouseover.stop="" @mouseout.stop="" @click.stop="optype='add_key';ing=true">ADD</va-button>
			<span style="width:60px;padding:12px 20px;color:var(--va-primary)">{{ t_keys?'▲':'▼' }}</span>
		</div>
		<!-- keys -->
		 <div v-if="t_keys&&keys.size" style="overflow-y:auto;height:auto;max-height:100%">
			<div v-for="key of keys.keys()" style="margin:1px 20px;display:flex;flex-direction:column">
				<div
					style="cursor:pointer;display:flex;align-items:center"
					:style="{'background-color':keys.get(key).hover?'var(--va-shadow)':'var(--va-background-element)'}"
					@click="keys.get(key).open=!keys.get(key).open"
					@mouseover="keys.get(key).hover=true"
					@mouseout="keys.get(key).hover=false"
				>
					<span style="width:35px;padding:12px;color:var(--va-primary)"> {{ keys.get(key).open?'▼':'►' }} </span>
					<span style="padding:12px;color:var(--va-primary)">{{key}}</span>
					<span style="flex:1"></span>
					<span style="padding:12px;color:var(--va-primary)">Current Config ID:  {{ keys.get(key).cur_index }}</span>
					<va-dropdown trigger="hover" style="width:60px;margin-right:10px" prevent-overflow placement="bottom-end">
						<template #anchor>
							<va-button style="width:60px;height:30px" size="small" @mouseover.stop="" @mouseout.stop="" @click.stop="">History</va-button>
						</template>
						<va-dropdown-content>
							<div style="max-height:300px;overflow-y:auto;display:flex;flex-direction:column">
								<va-button
									v-for="index of keys.get(key).max_index"
									size="small"
									style="height:24px;width:42px;padding:5px 0;margin:1px;cursor:pointer"
									:disabled="keys.get(key).cur_index==keys.get(key).max_index-index+1"
									@click="cur_key=key;cur_key_index=keys.get(key).max_index-index+1;optype='get_key';app_op()"
								>
									{{keys.get(key).max_index-index+1}}
								</va-button>
							</div>
						</va-dropdown-content>
					</va-dropdown>
					<va-button style="width:50px;height:30px;margin-right:80px" size="small" @mouseover.stop="" @mouseout.stop="" @click.stop="optype='del_key';cur_key=key;ing=true">DEL</va-button>
				</div>
				<div v-if="keys.get(key).open" style="display:flex;margin:2px 20px">
					<div style="flex:1;display:flex;flex-direction:column;align-items:center">
						<va-input
							:label="'Content Type:'+keys.get(key).cur_value_type"
							:model-value="JSON.stringify(JSON.parse(keys.get(key).cur_value),null,4)"
							style="width:100%"
							type="textarea"
							outline
							readonly
							:min-rows="15"
							:max-rows="15"
						/>
						<va-button
							v-if="keys.get(key).open"
							style="margin:2px"
							:disabled="Boolean(keys.get(key).new_cur_value)"
							@click="()=>{
								if(keys.get(key).cur_value){
									keys.get(key).new_cur_value=JSON.stringify(JSON.parse(keys.get(key).cur_value),null,4)
									keys.get(key).new_cur_value_type=keys.get(key).cur_value_type
								}else{
									keys.get(key).new_cur_value='{\n}'
									keys.get(key).new_cur_value_type='json'
								}
							}"
						>
							Edit
						</va-button>
					</div>
					<va-divider v-if="keys.get(key).new_cur_value" vertical />
					<div v-if="keys.get(key).new_cur_value" style="flex:1;display:flex;flex-direction:column;align-items:center">
						<va-input
							v-model.trim="keys.get(key).new_cur_value"
							style="width:100%"
							type="textarea"
							outline
							:min-rows="15"
							:max-rows="15"
						/>
						<div style="width:100%;display:flex;align-items:center">
							<va-radio v-for="(option,index) in config_value_types" :key="index" :option="option" v-model="keys.get(key).new_cur_value_type" style="margin:4px" disabled />
							<span style="flex:1"></span>
							<va-button style="margin-right:2px" @click="cur_key=key;optype='update_key';ing=true">Update</va-button>
							<va-button style="margin-left:2px" @click="keys.get(key).new_cur_value='';keys.get(key).new_cur_value_type=''">Cancel</va-button>
						</div>
					</div>
				</div>
			</div>
		</div>
		<div v-if="t_keys&&!keys.size">
			<div style="margin:1px 20px;padding:12px;display:flex;flex-direction:column;background-color:var(--va-background-element);color:var(--va-primary)">No Config Keys</div>
		</div>
		<!-- proxys -->
		<div
			style="width:100%;display:flex;align-items:center;margin:1px 0;cursor:pointer"
			:style="{'background-color':t_proxys_hover?'var(--va-shadow)':'var(--va-background-element)'}"
			@click="t_proxys=!t_proxys"
			@mouseover="t_proxys_hover=true"
			@mouseout="t_proxys_hover=false"
		>
			<span style="flex:1;padding:12px;color:var(--va-primary)">Proxys</span>
			<va-button style="height:30px" size="small" :disabled="curg==''||cura==''" @mouseover.stop="" @mouseout.stop="" @click.stop="cur_proxy=proxy;optype='add_proxy';ing=true">ADD</va-button>
			<span style="width:60px;padding:12px 20px;color:var(--va-primary)">{{ t_proxys?'▲':'▼' }}</span>
		</div>
		<!-- paths -->
		<div v-if="t_proxys&&proxys.size" style="overflow-y:auto;height:auto;max-height:100%">
			<div v-for="proxy of proxys.keys()" style="margin:1px 20px;display:flex;flex-direction:column">
				<div
					style="cursor:pointer;display:flex;align-items:center"
					:style="{'background-color':proxys.get(proxy).hover?'var(--va-shadow)':'var(--va-background-element)'}"
					@mouseover="proxys.get(proxy).hover=true"
					@mouseout="proxys.get(proxy).hover=false"
					@click="proxys.get(proxy).open=!proxys.get(proxy).open;proxys.get(proxy).req='{\n}'"
				>
					<span style="width:35px;padding:12px;color:var(--va-primary)"> {{ proxys.get(proxy).open?'▼':'►' }} </span>
					<span style="padding:12px;color:var(--va-primary)">{{proxy}}</span>
					<div style="flex:1"></div>
					<va-dropdown
						style="width:80px;margin-right:10px"
						placement="bottom-end"
						prevent-overflow
						:close-on-content-click="false"
						@mouseover.stop=""
						@mouseout.stop=""
						@click.stop="proxys.get(proxy).new_read=proxys.get(proxy).read;proxys.get(proxy).new_write=proxys.get(proxy).write;proxys.get(proxy).new_admin=proxys.get(proxy).admin"
					>
						<template #anchor>
							<va-button style="width:80px;height:30px" size="small">Permission</va-button>
						</template>
						<va-dropdown-content>
							<div style="display:flex;flex-direction:column">
								<va-switch
									style="margin:2px 0"
									v-model="proxys.get(proxy).new_read"
									true-inner-label="Read"
									false-inner-label="Read"
									@update:model-value="cur_proxy_permission_update(proxy,'read')"
								/>
								<va-switch
									style="margin:2px 0"
									v-model="proxys.get(proxy).new_write"
									true-inner-label="Write"
									false-inner-label="Write"
									@update:model-value="cur_proxy_permission_update(proxy,'write')"
								/>
								<va-switch
									style="margin:2px 0"
									v-model="proxys.get(proxy).new_admin"
									true-inner-label="Admin"
									false-inner-label="Admin"
									@update:model-value="cur_proxy_permission_update(proxy,'admin')"
								/>
								<va-button
									style="margin:2px 0"
									:disabled="proxys.get(proxy).new_read==proxys.get(proxy).read&&proxys.get(proxy).new_write==proxys.get(proxy).write&&proxys.get(proxy).new_admin==proxys.get(proxy).admin"
									@click="cur_proxy=proxy;optype='update_proxy';ing=true"
								>
									Update
								</va-button>
							</div>
						</va-dropdown-content>
					</va-dropdown>
					<va-button size="small" style="width:50px;height:30px;margin-right:80px" @mouseover.stop="" @mouseout.stop="" @click.stop="cur_proxy=proxy;optype='del_proxy';ing=true">DEL</va-button>
				</div>
				<div v-if="proxys.get(proxy).open" style="display:flex;margin:2px 20px">
					<div style="flex:1;display:flex;flex-direction:column;align-items:center">
						<va-input type="textarea" outline label="Request" :min-rows="15" :max-rows="15" style="width:100%" v-model.trim="proxys.get(proxy).req" :readonly="Boolean(proxys.get(proxy).resp)" />
						<va-button style="margin:2px" @click="cur_proxy=proxy;optype='proxy';ing=true" :disabled="Boolean(proxys.get(proxy).resp)">Proxy</va-button>
					</div>
					<va-divider v-if="proxys.get(proxy).resp" vertical />
					<div v-if="proxys.get(proxy).resp" style="flex:1;display:flex;flex-direction:column;align-items:center">
						<va-input type="textarea" outline label="Response" :min-rows="15" :max-rows="15" style="width:100%" v-model.trim="proxys.get(proxy).resp" readonly />
						<va-button style="margin:2px" @click="proxys.get(proxy).resp=''">OK</va-button>
					</div>
				</div>
			</div>
		</div>
		<div v-if="t_proxys&&!proxys.size">
			<div style="margin:1px 20px;padding:12px;display:flex;flex-direction:column;background-color:var(--va-background-element);color:var(--va-primary)">No Proxy Paths</div>
		</div>
	</div>
</template>
