<script setup lang="ts">
import { ref } from 'vue'
import * as initializeAPI from '../../api/initialize_browser_toc'
import * as permissionAPI from '../../api/permission_browser_toc'
import * as userAPI from '../../api/user_browser_toc'
import * as appAPI from '../../api/app_browser_toc'

import sidemenu from './sidemenu.vue'
import apppage from './apppage.vue'
import userpage from './userpage.vue'
import rolepage from './rolepage.vue'
import * as state from './state'

const host: string = "http://10.1.134.245:8000"
const initializeClient: initializeAPI.InitializeBrowserClientToC = new initializeAPI.InitializeBrowserClientToC(host)
const permissionClient: permissionAPI.PermissionBrowserClientToC = new permissionAPI.PermissionBrowserClientToC(host)
const userClient: userAPI.UserBrowserClientToC = new userAPI.UserBrowserClientToC(host)
const appClient: appAPI.AppBrowserClientToC = new appAPI.AppBrowserClientToC(host)

const inited = ref(false)
get_init_status()

const init_access_key=ref("")
const t_init_access_key=ref(false)
const init_password=ref("")
const t_init_password=ref(false)
function init_able():boolean{
	return init_access_key.value && init_password.value.length>=10 && init_password.value.length<32
}
function do_init(){
	if(!init_able()){
		if(init_access_key.value){
			state.set_error("error",-2,"Root Password length must in [10,32)!")
		}else{
			state.set_error("error",-2,"Missing Access Key!")
		}
		return
	}
	if(!state.set_load()){
		return
	}
	initializeClient.init({"Access-Key":init_access_key.value},{password:init_password.value},1000,(e: initializeAPI.Error)=>{
		state.clear_load()
		state.set_error("error",e.code,e.msg)
	},(resp: initializeAPI.InitResp)=>{
		state.clear_load()
		init_access_key.value=""
		init_password.value=""
		inited.value=true
	})
}
function get_init_status(){
	if(!state.set_load()){
		return
	}
	initializeClient.init_status({},{},1000,(e: initializeAPI.Error)=>{
		state.clear_load()
		state.set_error("error",e.code,e.msg)
	},(resp: initializeAPI.InitStatusResp)=>{
		state.clear_load()
		inited.value=resp.status
	})
}

const password = ref("")
const t_password = ref(false)
function login_root_able():boolean{
	return password.value.length>=10&&password.value.length<32
}
function do_login_root(){
	if(!login_root_able()){
		state.set_error("error",-2,"Root Password length must in [10,32)!")
		return
	}
	if(!state.set_load()){
		return
	}
	initializeClient.root_login({},{password:password.value},1000,(e: initializeAPI.Error)=>{
		state.clear_load()
		state.set_error("error",e.code,e.msg)
	},(resp: initializeAPI.RootLoginResp)=>{
		//clear loading in get_projects function
		password.value=""
		state.login(resp.token,null)
		get_projects(false)
	})
}

const password_changing = ref(false)
const oldpassword = ref("")
const t_oldpassword = ref(false)
const newpassword = ref("")
const t_newpassword = ref(false)
function change_root_password_able():boolean{
	return oldpassword.value.length>=10 && oldpassword.value.length<32 && newpassword.value.length>=10 && newpassword.value.length<32
}
function do_change_root_password(){
	if(!state.user.root || state.user.token==""){
		return
	}
	if(!change_root_password_able()){
		state.set_error("error",-2,"Root Password length must in [10,32)!")
		return
	}
	if(!state.set_load()){
		return
	}
	initializeClient.root_password({"Token":state.user.token},{old_password:oldpassword.value,new_password:newpassword.value},1000,(e: initializeAPI.Error)=>{
		state.clear_load()
		state.set_error("error",e.code,e.msg)
	},(resp: initializeAPI.RootPasswordResp)=>{
		state.clear_load()
		oldpassword.value=""
		newpassword.value=""
		password_changing.value=false
		state.logout()
	})
}
const oauth2 = ref("")
const oauth2s = ref(["Oauth2 Service Name 1","Oauth2 Service Name 2"])
const oauth2img = ref("")
function update_oauth2(){

}
function do_login_user(){

}

function get_projects(need_set_load: boolean){
	if(need_set_load){
		if(!state.set_load()){
			return
		}
	}
	initializeClient.list_project({"Token":state.user.token},{},1000,(e: initializeAPI.Error)=>{
		state.clear_load()
		state.set_error("error",e.code,e.msg)
	},(resp: initializeAPI.ListProjectResp)=>{
		state.clear_load()
		state.project.all=resp.projects
		//if the project doesn't exist,selector need to be reset
		if(state.project.cur){
			let find: boolean=false
			for(let p of state.project.all){
				if(state.project.cur.project_id[0]==p.project_id[0]&&state.project.cur.project_id[1]==p.project_id[1]){
					find=true
					state.project.cur=p
					break
				}
			}
			if(!find){
				state.project.cur=""
				state.project.nodes=[]
			}
		}
	})
}
function select_project(need_set_load: boolean){
	if(!state.project.cur){
		state.project.nodes=[]
		return
	}
	if(need_set_load){
		if(!state.set_load()){
			return
		}
	}
	permissionClient.list_user_node({"Token":state.user.token},{project_id:state.project.cur.project_id,user_id:"",need_user_role_node:true},1000,(e: permissionAPI.Error)=>{
		state.clear_load()
		state.clear_page()
		state.project.nodes=[]
		state.set_error("error",e.code,e.msg)
	},(resp: permissionAPI.ListUserNodeResp)=>{
		state.clear_load()
		state.clear_page()
		state.project.nodes=resp.nodes
	})
}
function is_root_node_project(): boolean{
	if(state.project.cur){
		if(state.project.nodes.length!=1){
			return false
		}
		if(state.project.nodes[0].node_id.length!=2){
			return false
		}
		if(state.project.nodes[0].node_id[0]!=state.project.cur.project_id[0]){
			return false
		}
		if(state.project.nodes[0].node_id[1]!=state.project.cur.project_id[1]){
			return false
		}
		return true
	}
	return false
}
function is_project_admin(): boolean{
	return state.project.cur&&state.project.cur.project_id[0]==0&&state.project.cur.project_id[1]==1
}

function project_op(){
	if(!state.project.cur){
		return
	}
	if(!state.set_load()){
		return
	}
	switch(state.project.optype){
		case 'add':{
			let req = {
				project_name: state.project.new_project_name,
				project_data: "",
			}
			initializeClient.create_project({"Token":state.user.token},req,1000,(e: initializeAPI.Error)=>{
				state.clear_load()
				state.set_error("error",e.code,e.msg)
			},(resp: initializeAPI.CreateProjectResp)=>{
				//clear loading in get_projects function
				get_projects(false)
				state.clear_project()
			})
			break
		}
		case 'update':{
			let req = {
				project_id: state.project.cur.project_id,
				new_project_name: state.project.new_project_name,
				new_project_data: state.project.cur.project_data,
			}
			initializeClient.update_project({"Token":state.user.token},req,1000,(e: initializeAPI.Error)=>{
				state.clear_load()
				state.set_error("error",e.code,e.msg)
			},(resp: initializeAPI.CreateProjectResp)=>{
				//clear loading in get_projects function
				get_projects(false)
				state.clear_project()
			})
			break
		}
		case 'del':{
			let req = {
				project_id: state.project.cur.project_id,
			}
			initializeClient.delete_project({"Token":state.user.token},req,1000,(e :initializeAPI.Error)=>{
				state.clear_load()
				state.set_error("error",e.code,e.msg)
			},(resp :initializeAPI.DeleteProjectResp)=>{
				//clear loading in get_projects function
				get_projects(false)
				state.clear_project()
			})
			break
		}
		default:{
			state.set_error("error",-2,"unknown operation")
		}
	}
}

function node_op(){
	if(!state.project.cur){
		return
	}
	if(!state.set_load()){
		return
	}
	switch(state.node.optype){
		case 'add':{
			let req = {
				pnode_id:state.node.target.node_id,
				node_name:state.node.new_node_name,
				node_data:state.node.new_node_url,
			}
			permissionClient.add_node({"Token":state.user.token},req,1000,(e :permissionAPI.Error)=>{
				state.clear_load()
				state.set_error("error",e.code,e.msg)
			},(resp :permissionAPI.AddNodeResp)=>{
				//clear loading in select_project
				select_project(false)
				state.clear_node()
			})
			break
		}
		case 'update':{
			let req = {
				node_id:state.node.target.node_id,
				new_node_name:state.node.new_node_name,
				new_node_data:state.node.new_node_url,
			}
			permissionClient.update_node({"Token":state.user.token},req,1000,(e :permissionAPI.Error)=>{
				state.clear_load()
				state.set_error("error",e.code,e.msg)
			},(resp :permissionAPI.AddNodeResp)=>{
				//clear loading in select_project
				select_project(false)
				state.clear_node()
			})
			break
		}
		case 'del':{
			let req = {
				node_id:state.node.target.node_id,
			}
			permissionClient.del_node({"Token":state.user.token},req,1000,(e :permissionAPI.Error)=>{
				state.clear_load()
				state.set_error("error",e.code,e.msg)
			},(resp :permissionAPI.AddNodeResp)=>{
				//clear loading in select_project
				select_project(false)
				state.clear_node()
			})
			break
		}
		default:{
			state.set_error("error",-2,"unknown operation")
		}
	}
}
function iframeload(){
	console.log("iframe")
}
</script>

<template>
	<va-modal v-model="state.load.ing" hide-default-actions no-dismiss :overlay="false" blur background-color="#0000" z-index="1000">
		<template #default>
			<va-inner-loading icon="❃" loading :size="55"></va-inner-loading>
		</template>
	</va-modal>
	<va-modal v-model="state.alert.ing" max-width="600px" max-height="400px" fixed-layout :title="state.alert.title+':'+state.alert.code" :message="state.alert.msg" hide-default-actions :overlay="false" blur z-index="1000" />
	<va-modal v-model="password_changing" attach-element="#app" max-width="600px" hide-default-actions no-dismiss overlay-opacity="0.2" z-index="999">
		<template #default>
			<div style="display:flex;flex-direction:column">
				<va-input :type="t_oldpassword?'text':'password'" label="Old Root Password*" v-model="oldpassword" style="width:400px;margin:5px 0" @keyup.enter="()=>{if(change_root_password_able()){do_change_root_password()}}">
					<template #appendInner>
						<va-icon :name="t_oldpassword?'◎':'◉'" size="small" color="--va-primary" @click="t_oldpassword=!t_oldpassword" />
					</template>
				</va-input>
				<va-input :type="t_newpassword?'text':'password'" label="New Root Password*" v-model="newpassword" style="width:400px;margin:5px 0" @keyup.enter="()=>{if(change_root_password_able()){do_change_root_password()}}">
					<template #appendInner>
						<va-icon :name="t_newpassword?'◎':'◉'" size="small" color="--va-primary" @click="t_newpassword=!t_newpassword" />
					</template>
				</va-input>
				<div>
					<va-button style="width:100px;margin:5px 5px 0 190px" @click="oldpassword='';newpassword='';password_changing=false">Cancel</va-button>
					<va-button style="width:100px;margin:5px 0 0 5px" :disabled="!change_root_password_able()" @click="do_change_root_password">Change</va-button>
				</div>
			</div>
		</template>
	</va-modal>
	<va-modal v-model="state.project.ing" attach-element="#app" max-width="600px" hide-default-actions no-dismiss overlay-opacity="0.2" z-index="999">
		<template #default>
			<div v-if="state.project.optype=='del'" style="display:flex;flex-direction:column">
				<va-card color="primary" gradient style="margin:0 0 5px 0">
					<va-card-title>Warning</va-card-title>
					<va-card-content>
						<p>You are deleting project: {{ state.project.cur.project_name }}.</p>
						<p>All data in this project will be deleted.</p>
						<p>Please confirm!</p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="project_op" gradient>Del</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="state.clear_project" gradient>Cancel</va-button>
				</div>
			</div>

			<div v-if="state.project.optype=='add'">
				<va-input type="text" style="width:250px" label="New Project Name*" v-model="state.project.new_project_name" @keyup.enter="()=>{if(state.project.new_project_name){project_op()}}"></va-input>
				<va-button style="width:80px;margin:0 0 0 5px" :disabled="!state.project.new_project_name" @click="project_op">Add</va-button>
				<va-button style="width:80px;margin:0 0 0 5px" @click="state.clear_project">Cancel</va-button>
			</div>

			<div v-if="state.project.optype=='update'">
				<va-input type="text" style="width:250px" label="New Project Name*" v-model="state.project.new_project_name" @keyup.enter="()=>{if(state.project.new_project_name){project_op()}}"></va-input>
				<va-button style="width:80px;margin:0 0 0 5px" :disabled="!state.project.new_project_name" @click="project_op">Update</va-button>
				<va-button style="width:80px;margin:0 0 0 5px" @click="state.clear_project">Cancel</va-button>
			</div>
		</template>
	</va-modal>
	<va-modal v-model="state.node.ing" attach-element="#app" max-width="600px" hide-default-actions no-dismiss overlay-opacity="0.2" z-index="999">
		<template #default>
			<div v-if="state.node.optype=='del'" style="display:flex;flex-direction:column">
				<va-card color="primary" gradient style="margin:0 0 5px 0">
					<va-card-title>Warning</va-card-title>
					<va-card-content>
						<p>You are deleting node: {{ state.node.target.node_name }}.</p>
						<p>All data in this node will be deleted.</p>
						<p>Please confirm!</p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="node_op" gradient>Del</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="state.clear_node" gradient>Cancel</va-button>
				</div>
			</div>

			<div v-if="state.node.optype=='add'" style="display:flex;flex-direction:column">
				<va-input style="width:400px;margin:3px 0" label="New Node Name*" v-model="state.node.new_node_name" @keyup.enter="()=>{if(state.node.new_node_name){node_op()}}"></va-input>
				<va-input style="width:400px;margin:3px 0" label="New Node Url?" v-model="state.node.new_node_url" @keyup.enter="()=>{if(state.node.new_node_name){node_op()}}"></va-input>
				<div>
					<va-button style="width:80px;margin:3px 3px 0 234px" :disabled="!state.node.new_node_name" @click="node_op">Add</va-button>
					<va-button style="width:80px;margin:3px 0 0 3px" @click="state.clear_node">Cancel</va-button>
				</div>
			</div>

			<div v-if="state.node.optype=='update'" style="display:flex;flex-direction:column">
				<va-input style="width:400px;margin:3px 0" label="New Node Name*" v-model="state.node.new_node_name" @keyup.enter="()=>{if(state.node.new_node_name){node_op()}}"></va-input>
				<va-input style="width:400px;margin:3px 0" label="New Node Url?" v-model="state.node.new_node_url" @keyup.enter="()=>{if(state.node.new_node_name){node_op()}}"></va-input>
				<div>
					<va-button style="width:80px;margin:3px 3px 0 234px" :disabled="!state.node.new_node_name" @click="node_op">Update</va-button>
					<va-button style="width:80px;margin:3px 0 0 3px" @click="state.clear_node">Cancel</va-button>
				</div>
			</div>
		</template>
	</va-modal>

	<div v-if="!inited" style="display:flex;flex-direction:column;width:100%;justify-content:center;align-items:center">
		<va-card style="width:400px;margin:10px" color="primary" gradient>
			<va-card-title>Warning</va-card-title>
			<va-card-content>System not initialized.</va-card-content>
		</va-card>
		<div style="display:flex;flex-direction:column">
			<va-input :type="t_init_access_key?'text':'password'" label="Access Key*" v-model="init_access_key" style="width:400px;margin:5px 0" @keyup.enter="()=>{if(init_able()){do_init()}}">
				<template #appendInner>
					<va-icon :name="t_init_access_key?'◎':'◉'" size="small" color="--va-primary" @click="t_init_access_key=!t_init_access_key" />
				</template>
			</va-input>
			<va-input :type="t_init_password?'text':'password'" label="Root Password*" v-model="init_password" style="width:400px;margin:5px 0" @keyup.enter="()=>{if(init_able()){do_init()}}">
				<template #appendInner>
					<va-icon :name="t_init_password?'◎':'◉'" size="small" color="--va-primary" @click="t_init_password=!t_init_password" />
				</template>
			</va-input>
			<va-button style="width:100px;margin:5px 0 0 300px" :disabled="!init_able()" @click="do_init" gradient>Init</va-button>
		</div>
	</div>
	<div v-else-if="state.user.token.length==0" style="display:flex;width:100%;justify-content:center;align-items:center">
		<div v-if="!state.user.root">
			<va-select v-model="oauth2" :options="oauth2s" no-options-text="NO Oauth2 Login" placeholder="Select Oauth2 Login" dropdown-icon="" style="width:400px" @update:model-value="update_oauth2"></va-select>
			<va-image style="width:400px;height:400px;margin:5px 0" :src="oauth2img" />
			<va-button style="width:400px;margin:0" @click="state.user.root=true">Switch To Root User Login</va-button>
		</div>
		<div v-else>
			<div>
				<va-input :type="t_password?'text':'password'" style="width:300px" label="Root Password*" v-model="password" @keyup.enter="()=>{if(login_root_able()){do_login_root()}}">
					<template #appendInner>
						<va-icon :name="t_password?'◎':'◉'" size="small" color="--va-primary" @click="t_password=!t_password" />
					</template>
				</va-input>
				<va-button style="width:90px;margin:0 0 0 10px" :disabled="!login_root_able()" @click="do_login_root">Login</va-button>
			</div>
			<va-button style="width:400px;margin:10px 0 0 0" @click="state.user.root=false">Switch To Normal User Login</va-button>
		</div>
	</div>
	<div v-else style="display:flex;width:100%">
		<div style="display:flex;flex-direction:column;width:300px">
			<div style="display:flex;margin:4px 0">
				<va-select
					style="margin:0 2px;flex:1"
					v-model="state.project.cur"
					:options="state.project.all"
					text-by="project_name"
					track-by="project_id"
					placeholder="Select Project"
					dropdown-icon=""
					@update:model-value="select_project(true)"
				/>
				<va-dropdown v-if="state.user.root||(is_root_node_project()&&state.project.nodes[0].admin)" trigger="hover" style="width:36px;margin-right:2px">
					<template #anchor>
						<va-button>•••</va-button>
					</template>
					<va-dropdown-content>
						<va-popover message="Create New Project" :hover-out-timeout="0" :hover-over-timeout="0" color="primary">
							<va-button v-if="state.user.root" style="width:36px;margin-right:2px" @click="state.set_project('add')">+</va-button>
						</va-popover>
						<va-popover message="Rename Project" :hover-out-timeout="0" :hover-over-timeout="0" color="primary">
							<va-button v-if="state.user.root&&!is_project_admin()" style="width:36px;margin:0 3px" @click="state.set_project('update')">◉</va-button>
						</va-popover>
						<va-popover message="Delete Project" :hover-out-timeout="0" :hover-over-timeout="0" color="primary">
							<va-button v-if="state.user.root&&!is_project_admin()" style="width:36px;margin:0 3px" @click="state.set_project('del')">x</va-button>
						</va-popover>
						<va-popover message="Add Menu" :hover-out-timeout="0" :hover-over-timeout="0" color="primary">
							<va-button v-if="!is_project_admin()&&is_root_node_project()&&state.project.nodes[0].admin" style="width:36px;margin:0 3px" @click="state.set_node(state.project.nodes[0],'add')">✿</va-button>
						</va-popover>
					</va-dropdown-content>
				</va-dropdown>
			</div>
			<va-divider style="margin:0" />
			<div style="flex:1;overflow-x:hidden;overflow-y:auto">
				<sidemenu :nodes="is_root_node_project()?state.project.nodes[0].children:state.project.nodes" :deep="0" />
			</div>
		</div>
		<va-divider vertical style="margin:0" />
		<div style="display:flex;flex-direction:column;flex:1">
			<div style="display:flex;margin:4px">
				<div style="display:flex;flex:1">
				</div>
				<va-dropdown trigger="hover" style="width:36px" placement="bottom-end">
					<template #anchor>
						<va-button round>{{ state.avatar() }}</va-button>
					</template>
					<va-dropdown-content>
						<div style="display:flex;flex-direction:column">
							<va-button v-if="state.user.root" style="margin:0 0 3px 0" @click="password_changing=true">ChangePassword</va-button>
							<va-button v-if="state.user.root" style="margin:3px 0 0 0" @click="state.logout">Logout</va-button>
							<va-button v-if="!state.user.root" @click="state.logout">Logout</va-button>
						</div>
					</va-dropdown-content>
				</va-dropdown>
			</div>
			<va-divider style="margin:0" />
			<div v-if="state.page.node&&state.page.node.node_id.length==3&&state.page.node.node_id[2]==1">
			<!-- User -->
				<userpage></userpage>
			</div>
			<div v-else-if="state.page.node&&state.page.node.node_id.length==3&&state.page.node.node_id[2]==2">
			<!-- Role -->
				<rolepage></rolepage>
			</div>
			<div v-else-if="state.page.node&&state.page.node.node_id.length==3&&state.page.node.node_id[2]==3">
			<!-- App -->
				<apppage></apppage>
			</div>
			<iframe v-else-if="state.page.node&&state.page.node.node_data!=''" width="100%" height="100%" frameborder="0" :src="state.page.node.node_data" @load="iframeload"></iframe>
		</div>
	</div>
</template>
