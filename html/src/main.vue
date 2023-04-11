<script setup lang="ts">
import { ref } from 'vue'
import * as initializeAPI from '../../api/initialize_browser_toc'
import * as permissionAPI from '../../api/permission_browser_toc'
import * as userAPI from '../../api/user_browser_toc'
import * as state from './state'
import * as client from './client'

import sidemenu from './sidemenu.vue'
import app from './app.vue'
import userrole from './userrole.vue'

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
			state.set_alert("error",-2,"Root Password length must in [10,32)!")
		}else{
			state.set_alert("error",-2,"Missing Access Key!")
		}
		return
	}
	if(!state.set_load()){
		return
	}
	client.initializeClient.init({"Access-Key":init_access_key.value},{password:init_password.value},client.timeout,(e: initializeAPI.Error)=>{
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
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
	client.initializeClient.init_status({},{},client.timeout,(e: initializeAPI.Error)=>{
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
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
		state.set_alert("error",-2,"Root Password length must in [10,32)!")
		return
	}
	if(!state.set_load()){
		return
	}
	client.initializeClient.root_login({},{password:password.value},client.timeout,(e: initializeAPI.Error)=>{
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
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
		state.set_alert("error",-2,"Root Password length must in [10,32)!")
		return
	}
	if(!state.set_load()){
		return
	}
	client.initializeClient.root_password({"Token":state.user.token},{old_password:oldpassword.value,new_password:newpassword.value},client.timeout,(e: initializeAPI.Error)=>{
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
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
function do_login_user(){

}

function get_projects(need_set_load: boolean){
	if(need_set_load){
		if(!state.set_load()){
			return
		}
	}
	client.initializeClient.list_project({"Token":state.user.token},{},client.timeout,(e: initializeAPI.Error)=>{
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
	},(resp: initializeAPI.ListProjectResp)=>{
		state.clear_load()
		state.project.all=resp.projects
		//if the project doesn't exist,selector need to be reset
		if(state.project.cur_id.length!=0){
			let find: boolean=false
			for(let p of state.project.all){
				if(same_node_id(p.project_id,state.project.cur_id)){
					find=true
					break
				}
			}
			if(!find){
				state.project.cur_id=[]
				state.project.cur_name=""
				state.project.nodes=[]
			}
		}
	})
}
function select_project(need_set_load: boolean){
	if(state.project.cur_id.length==0){
		//this is impossible
		state.project.nodes=[]
		return
	}
	if(need_set_load){
		if(!state.set_load()){
			return
		}
	}
	client.permissionClient.list_user_node({"Token":state.user.token},{project_id:state.project.cur_id,user_id:"",need_user_role_node:true},client.timeout,(e: permissionAPI.Error)=>{
		state.clear_load()
		state.clear_page()
		state.project.nodes=[]
		state.set_alert("error",e.code,e.msg)
	},(resp: permissionAPI.ListUserNodeResp)=>{
		state.clear_load()
		state.clear_page()
		state.project.nodes=resp.nodes
	})
}

function project_op(){
	if(state.project.cur_id.length==0){
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
			client.initializeClient.create_project({"Token":state.user.token},req,client.timeout,(e: initializeAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: initializeAPI.CreateProjectResp)=>{
				//clear loading in get_projects function
				get_projects(false)
				state.clear_project()
			})
			break
		}
		case 'update':{
			let req = {
				project_id: state.project.cur_id,
				new_project_name: state.project.new_project_name,
				new_project_data: "",
			}
			client.initializeClient.update_project({"Token":state.user.token},req,client.timeout,(e: initializeAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: initializeAPI.CreateProjectResp)=>{
				//clear loading in get_projects function
				get_projects(false)
				state.clear_project()
			})
			break
		}
		case 'del':{
			let req = {
				project_id: state.project.cur_id,
			}
			client.initializeClient.delete_project({"Token":state.user.token},req,client.timeout,(e :initializeAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp :initializeAPI.DeleteProjectResp)=>{
				//clear loading in get_projects function
				get_projects(false)
				state.clear_project()
			})
			break
		}
		default:{
			state.clear_load()
			state.set_alert("error",-2,"unknown operation")
		}
	}
}

function node_op(){
	if(state.project.cur_id.length==0){
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
			client.permissionClient.add_node({"Token":state.user.token},req,client.timeout,(e :permissionAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
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
			client.permissionClient.update_node({"Token":state.user.token},req,client.timeout,(e :permissionAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
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
			client.permissionClient.del_node({"Token":state.user.token},req,client.timeout,(e :permissionAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp :permissionAPI.AddNodeResp)=>{
				//clear loading in select_project
				select_project(false)
				state.clear_node()
			})
			break
		}
		default:{
			state.clear_load()
			state.set_alert("error",-2,"unknown operation")
		}
	}
}
function iframeload(){
	console.log("iframe")
}
function same_node_id(a:number[],b:number[]):boolean{
	if(!Boolean(a)&&!Boolean(b)){
		return true
	}else if(Boolean(a)&&!Boolean(b)&&a.length==0){
		return true
	}else if(!Boolean(a)&&Boolean(b)&&b.length==0){
		return true
	}
	if(a.length!=b.length){
		return false
	}
	for(let i=0;i<a.length;i++){
		if(a[i]!=b[i]){
			return false
		}
	}
	return true
}
</script>

<template>
	<va-modal v-model="state.load.ing" hide-default-actions no-dismiss :overlay="false" blur background-color="#0000" z-index="1000">
		<template #default>
			<va-inner-loading icon="❃" loading :size="55"></va-inner-loading>
		</template>
	</va-modal>
	<va-modal v-model="state.alert.ing" max-width="600px" max-height="400px" fixed-layout :title="state.get_alert_title()" :message="state.alert.msg" hide-default-actions :overlay="false" blur z-index="1000" />
	<va-modal v-model="password_changing" attach-element="#app" max-width="600px" hide-default-actions no-dismiss overlay-opacity="0.2" z-index="999">
		<template #default>
			<div style="display:flex;flex-direction:column">
				<va-input :type="t_oldpassword?'text':'password'" label="Old Root Password*" v-model="oldpassword" style="width:400px;margin:5px 0" @keyup.enter="()=>{if(change_root_password_able()){do_change_root_password()}}">
					<template #appendInner>
						<va-icon :name="t_oldpassword?'◎':'◉'" size="small" color="var(--va-primary)" @click="t_oldpassword=!t_oldpassword" />
					</template>
				</va-input>
				<va-input :type="t_newpassword?'text':'password'" label="New Root Password*" v-model="newpassword" style="width:400px;margin:5px 0" @keyup.enter="()=>{if(change_root_password_able()){do_change_root_password()}}">
					<template #appendInner>
						<va-icon :name="t_newpassword?'◎':'◉'" size="small" color="var(--va-primary)" @click="t_newpassword=!t_newpassword" />
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
						<p>You are deleting project: {{ state.project.cur_name}}.</p>
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

	<div v-if="!inited" style="width:100%;height:100%;display:flex;flex-direction:column;justify-content:center;align-items:center">
		<va-card style="width:400px;margin:10px" color="primary" gradient>
			<va-card-title>Warning</va-card-title>
			<va-card-content>System not initialized.</va-card-content>
		</va-card>
		<div style="display:flex;flex-direction:column">
			<va-input :type="t_init_access_key?'text':'password'" label="Access Key*" v-model="init_access_key" style="width:400px;margin:5px 0" @keyup.enter="()=>{if(init_able()){do_init()}}">
				<template #appendInner>
					<va-icon :name="t_init_access_key?'◎':'◉'" size="small" color="var(--va-primary)" @click="t_init_access_key=!t_init_access_key" />
				</template>
			</va-input>
			<va-input :type="t_init_password?'text':'password'" label="Root Password*" v-model="init_password" style="width:400px;margin:5px 0" @keyup.enter="()=>{if(init_able()){do_init()}}">
				<template #appendInner>
					<va-icon :name="t_init_password?'◎':'◉'" size="small" color="var(--va-primary)" @click="t_init_password=!t_init_password" />
				</template>
			</va-input>
			<va-button style="width:100px;margin:5px 0 0 300px" :disabled="!init_able()" @click="do_init" gradient>Init</va-button>
		</div>
	</div>
	<div v-else-if="state.user.token.length==0" style="width:100%;height:100%;display:flex;justify-content:center;align-items:center">
		<div v-if="!state.user.root">
			<va-select v-model="oauth2" :options="oauth2s" no-options-text="NO Oauth2 Login" label="Select Oauth2 Login" dropdown-icon="" style="width:400px" trigger="hover">
				<template #option='{option,index,selectOption}'>
					<va-hover
					stateful
					@click="()=>{
						if(oauth2!=option){
							selectOption(option)
						}
					}">
						<template #default="{hover}">
							<div
							style="padding:10px;cursor:pointer"
							:style="{'background-color':hover?'var(--va-background-border)':'',color:hover||oauth2==option?'var(--va-primary)':'black'}"
							>
								{{option}}
							</div>
						</template>
					</va-hover>
				</template>
			</va-select>
			<va-image style="width:400px;height:400px;margin:5px 0" :src="oauth2img" />
			<va-button style="width:400px;margin:0" @click="state.user.root=true">Switch To Root User Login</va-button>
		</div>
		<div v-else>
			<div>
				<va-input :type="t_password?'text':'password'" style="width:300px" label="Root Password*" v-model="password" @keyup.enter="()=>{if(login_root_able()){do_login_root()}}">
					<template #appendInner>
						<va-icon :name="t_password?'◎':'◉'" size="small" color="var(--va-primary)" @click="t_password=!t_password" />
					</template>
				</va-input>
				<va-button style="width:90px;margin:0 0 0 10px" :disabled="!login_root_able()" @click="do_login_root">Login</va-button>
			</div>
			<va-button style="width:400px;margin:10px 0 0 0" @click="state.user.root=false">Switch To Normal User Login</va-button>
		</div>
	</div>
	<div v-else style="width:100%;height:100%;display:flex">
		<div style="width:200px;display:flex;flex-direction:column;overflow-x:auto">
			<div style="display:flex;padding:5px 0;background-color:var(--va-background-element)">
				<va-select
				dropdown-icon=""
				trigger="hover"
				outline
				style="flex:1;margin:0 2px"
				:model-value="state.project.cur_name"
				:options="state.project.all"
				label="Select Project"
				no-options-text="NO Projects"
				>
					<template #option='{option,index,selectOption}'>
						<va-hover
						stateful
						@click="state.project.cur_id=option.project_id;state.project.cur_name=option.project_name;select_project(true)"
						>
							<template #default="{hover}">
								<div
								style="padding:10px;cursor:pointer"
								:style="{'background-color':hover?'var(--va-background-border)':'',color:hover||same_node_id(state.project.cur_id,option.project_id)?'var(--va-primary)':'black'}"
								>
									{{option.project_name}}
								</div>
							</template>
						</va-hover>
					</template>
				</va-select>
				<va-dropdown v-if="state.user.root||(same_node_id(state.project.cur_id,state.project.nodes[0].node_id)&&state.project.nodes[0].admin)" trigger="hover" style="width:36px;margin-right:2px">
					<template #anchor>
						<va-button>•••</va-button>
					</template>
					<va-dropdown-content>
						<va-popover message="Create New Project" :hover-out-timeout="0" :hover-over-timeout="0" color="primary">
							<va-button v-if="state.user.root" style="width:36px;margin:0 3px" @click="state.set_project('add')">+</va-button>
						</va-popover>
						<va-popover message="Rename Project" :hover-out-timeout="0" :hover-over-timeout="0" color="primary">
							<va-button v-if="state.user.root&&!same_node_id(state.project.cur_id,[0,1])" style="width:36px;margin:0 3px" @click="state.set_project('update')">◉</va-button>
						</va-popover>
						<va-popover message="Delete Project" :hover-out-timeout="0" :hover-over-timeout="0" color="primary">
							<va-button v-if="state.user.root&&!same_node_id(state.project.cur_id,[0,1])" style="width:36px;margin:0 3px" @click="state.set_project('del')">x</va-button>
						</va-popover>
						<va-popover message="Add Menu" :hover-out-timeout="0" :hover-over-timeout="0" color="primary">
							<va-button
							v-if="!same_node_id(state.project.cur_id,[0,1])&&same_node_id(state.project.cur_id,state.project.nodes[0].node_id)&&state.project.nodes[0].admin"
							style="width:36px;margin:0 3px"
							@click="state.set_node(state.project.nodes[0],'add')"
							>
								✿
							</va-button>
						</va-popover>
					</va-dropdown-content>
				</va-dropdown>
			</div>
			<div style="flex:1;overflow-x:hidden;overflow-y:auto;background-color:var(--va-background-element)">
				<sidemenu v-if="state.project.nodes.length>0" :nodes="same_node_id(state.project.cur_id,state.project.nodes[0].node_id)?state.project.nodes[0].children:state.project.nodes" :deep="0" />
			</div>
		</div>
		<div style="flex:1;display:flex;flex-direction:column;overflow-x:auto">
			<div style="display:flex;padding:5px;background-color:var(--va-background-element)">
				<div style="display:flex;flex:1"></div>
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
			<userrole v-if="state.page.node&&state.page.node.node_id.length==3&&state.page.node.node_id[2]==1"></userrole>
			<app v-else-if="state.page.node&&state.page.node.node_id.length==3&&state.page.node.node_id[2]==2"></app>
			<iframe v-else-if="state.page.node&&state.page.node.node_data!=''" width="100%" height="100%" frameborder="0" :src="state.page.node.node_data" @load="iframeload"></iframe>
		</div>
	</div>
</template>
