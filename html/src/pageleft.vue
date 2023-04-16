<script setup lang="ts">
import {ref,reactive,onMounted} from 'vue'

import * as initializeAPI from '../../api/initialize_browser_toc'
import * as permissionAPI from '../../api/permission_browser_toc'
import * as userAPI from '../../api/user_browser_toc'

import * as state from './state'
import * as client from './client'

import sidemenu from './sidemenu.vue'

onMounted(()=>{
	get_projects()
})

const allprojects=ref<initializeAPI.projectInfo[]>([])

function get_projects(){
	if(!state.set_load()){
		return
	}
	client.initializeClient.list_project({"Token":state.user.token},{},client.timeout,(e: initializeAPI.Error)=>{
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
	},(resp: initializeAPI.ListProjectResp)=>{
		allprojects.value=resp.projects
		//if the selected project doesn't exist,sidemenu need to be reset
		if(state.project.cur_id.length!=0){
			let find: boolean=false
			for(let i=0;i<allprojects.value.length;i++){
				if(same_node_id(allprojects.value[i].project_id,state.project.cur_id)){
					find=true
					state.project.cur_name=allprojects.value[i].project_name
					break
				}
			}
			if(!find){
				state.clear_project()
				allnodes=[]
			}
		}
		state.clear_load()
	})
}

const allnodes=ref<permissionAPI.NodeInfo[]>([])

function select_project(need_set_load: boolean){
	if(!state.set_load()){
		return
	}
	client.permissionClient.list_user_node({"Token":state.user.token},{project_id:state.project.cur_id,user_id:"",need_user_role_node:true},client.timeout,(e: permissionAPI.Error)=>{
		allnodes.value=[]
		state.clear_page()
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
	},(resp: permissionAPI.ListUserNodeResp)=>{
		allnodes.value=resp.nodes
		state.clear_page()
		state.clear_load()
	})
}

const ing=ref<boolean>(false)
const optype=ref<string>("")

const project_name=ref<string>("")

function op(){
	if(!state.set_load()){
		return
	}
	switch(optype.value){
		case 'add':{
			let req = {
				project_name: project_name.value,
				project_data: "",
			}
			client.initializeClient.create_project({"Token":state.user.token},req,client.timeout,(e: initializeAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: initializeAPI.CreateProjectResp)=>{
				project_name.value=''
				ing.value=false
				state.clear_load()
				get_projects()
			})
			break
		}
		case 'update':{
			let req = {
				project_id: state.project.cur_id,
				new_project_name: project_name.value,
				new_project_data: "",
			}
			client.initializeClient.update_project({"Token":state.user.token},req,client.timeout,(e: initializeAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: initializeAPI.CreateProjectResp)=>{
				project_name.value=''
				ing.value=false
				state.clear_load()
				get_projects()
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
				ing.value=false
				state.clear_project()
				allnodes.value=[]
				state.clear_load()
				get_projects()
			})
			break
		}
		default:{
			state.clear_load()
			state.set_alert("error",-2,"unknown operation")
		}
	}
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
	<va-modal v-model="ing" attach-element="#app" max-width="600px" hide-default-actions no-dismiss overlay-opacity="0.2" z-index="999">
		<template #default>
			<div v-if="optype=='add'">
				<va-input type="text" style="width:250px" label="New Project Name*" v-model="project_name" @keyup.enter="()=>{if(project_name!=''){op()}}"></va-input>
				<va-button style="width:80px;margin:0 0 0 5px" :disabled="project_name==''" @click="op">Add</va-button>
				<va-button style="width:80px;margin:0 0 0 5px" @click="project_name='';ing=false">Cancel</va-button>
			</div>
			<div v-else-if="optype=='update'">
				<va-input type="text" style="width:250px" label="New Project Name*" v-model="project_name" @keyup.enter="()=>{if(project_name!=''){op()}}"></va-input>
				<va-button style="width:80px;margin:0 0 0 5px" :disabled="project_name==''" @click="op">Update</va-button>
				<va-button style="width:80px;margin:0 0 0 5px" @click="project_name='';ing=false">Cancel</va-button>
			</div>
			<div v-else-if="optype=='del'" style="display:flex;flex-direction:column">
				<va-card color="primary" gradient style="margin:0 0 5px 0">
					<va-card-title>Warning</va-card-title>
					<va-card-content>
						<p>You are deleting project: {{ state.project.cur_name }}.</p>
						<p>All data in this project will be deleted.</p>
						<p>Please confirm!</p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="op" gradient>Del</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
		</template>
	</va-modal>
	<div style="width:200px;display:flex;flex-direction:column">
		<div style="display:flex;padding:5px 0;background-color:var(--va-background-element)">
			<va-select
			dropdown-icon=""
			trigger="hover"
			outline
			style="flex:1;margin:0 2px"
			:model-value="state.project.cur_name"
			:options="allprojects"
			label="Select Project"
			no-options-text="NO Projects"
			prevent-overflow
			>
				<template #option='{option,index,selectOption}'>
					<va-hover
					stateful
					@click="state.project.cur_id=option.project_id;state.project.cur_name=option.project_name;select_project()"
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
			<va-dropdown v-if="state.user.root" prevent-overflow trigger="hover" style="width:36px;margin-right:2px">
				<template #anchor>
					<va-button>•••</va-button>
				</template>
				<va-dropdown-content>
					<va-popover message="Create New Project" :hover-out-timeout="0" :hover-over-timeout="0" color="primary" prevent-overflow>
						<va-button v-if="state.user.root" style="width:36px;margin:0 3px" @click="optype='add';project_name='';ing=true">+</va-button>
					</va-popover>
					<va-popover message="Rename Project" :hover-out-timeout="0" :hover-over-timeout="0" color="primary" prevent-overflow>
						<va-button v-if="state.user.root&&!same_node_id(state.project.cur_id,[0,1])" style="width:36px;margin:0 3px" @click="optype='update';project_name='';ing=true">◉</va-button>
					</va-popover>
					<va-popover message="Delete Project" :hover-out-timeout="0" :hover-over-timeout="0" color="primary" prevent-overflow>
						<va-button v-if="state.user.root&&!same_node_id(state.project.cur_id,[0,1])" style="width:36px;margin:0 3px" @click="optype='del';ing=true">x</va-button>
					</va-popover>
				</va-dropdown-content>
			</va-dropdown>
		</div>
		<div style="flex:1;overflow-x:hidden;overflow-y:auto;background-color:var(--va-background-element)">
			<sidemenu v-if="allnodes.length>0" :nodes="allnodes" :deep="0" />
		</div>
	</div>

</template>
