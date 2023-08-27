<script setup lang="ts">
import {ref,onMounted} from 'vue'

import * as initializeAPI from './api/initialize_browser_toc'
import * as permissionAPI from './api/permission_browser_toc'

import * as state from './state'
import * as client from './client'

import menutree from './menutree.vue'

onMounted(()=>{
	get_projects()
})

const allprojects=ref<initializeAPI.ProjectInfo[]>([])

function get_projects(){
	if(!state.set_load()){
		return
	}
	let req = {}
	client.initializeClient.list_project({"Token":state.user.token},req,client.timeout,(e: initializeAPI.Error)=>{
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
	},(resp: initializeAPI.ListProjectResp)=>{
		if(resp.projects){
			let tmp: initializeAPI.ProjectInfo[] = []
			for(let i=0;i<resp.projects.length;i++){
				if(resp.projects[i]){
					tmp.push(resp.projects[i]!)
				}
			}
			allprojects.value = tmp
		}else{
			allprojects.value=[]
		}
		//if the selected project doesn't exist,sidemenu need to be reset
		if(state.project.info){
			let find: boolean=false
			for(let i=0;i<allprojects.value.length;i++){
				if(!allprojects.value[i].project_id){
					continue
				}
				if(same_node_id(allprojects.value[i].project_id!,state.project.info!.project_id!)){
					find=true
					state.project.info = allprojects.value[i]
					break
				}
			}
			if(!find){
				state.clear_project()
				projectnodes.value=null
			}
		}
		state.clear_load()
	})
}

const projectnodes=ref<permissionAPI.NodeInfo|null>(null)

function select_project(){
	if(!state.set_load()){
		return
	}
	let req = {
		project_id:state.project.info!.project_id,
		user_id:"",
		need_user_role_node:true,
	}
	client.permissionClient.list_user_node({"Token":state.user.token},req,client.timeout,(e: permissionAPI.Error)=>{
		projectnodes.value=null
		state.clear_project()
		state.clear_page()
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
	},(resp: permissionAPI.ListUserNodeResp)=>{
		if(resp.node){
			projectnodes.value=resp.node!
		}else{
			projectnodes.value=null
		}
		state.clear_page()
		state.clear_load()
	})
}

const node_ing=ref<boolean>(false)
const project_ing=ref<boolean>(false)
const optype=ref<string>("")

const project_name=ref<string>("")

function project_op(){
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
			},(_resp: initializeAPI.CreateProjectResp)=>{
				project_name.value=''
				project_ing.value=false
				state.clear_load()
				get_projects()
			})
			break
		}
		case 'update':{
			let req = {
				project_id: state.project.info!.project_id,
				new_project_name: project_name.value,
				new_project_data: "",
			}
			client.initializeClient.update_project({"Token":state.user.token},req,client.timeout,(e: initializeAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(_resp: initializeAPI.UpdateProjectResp)=>{
				project_name.value=''
				project_ing.value=false
				state.clear_load()
				get_projects()
			})
			break
		}
		case 'del':{
			let req = {
				project_id: state.project.info!.project_id,
			}
			client.initializeClient.delete_project({"Token":state.user.token},req,client.timeout,(e :initializeAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(_resp :initializeAPI.DeleteProjectResp)=>{
				project_ing.value=false
				state.clear_project()
				projectnodes.value=null
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

const ptarget=ref<permissionAPI.NodeInfo|null>(null)
const target=ref<permissionAPI.NodeInfo|null>(null)
const node_name=ref<string>("")
const node_url=ref<string>("")

function node_op(){
	if(!state.project.info){
		return
	}
	if(!state.set_load()){
		return
	}
	switch(optype.value){
		case 'add':{
			let req:permissionAPI.AddNodeReq = {
				pnode_id:target.value!.node_id,
				node_name:node_name.value,
				node_data:node_url.value,
			}
			client.permissionClient.add_node({"Token":state.user.token},req,client.timeout,(e :permissionAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp :permissionAPI.AddNodeResp)=>{
				node_ing.value=false
				if(target.value!.children){
					target.value!.children.push({
						node_id:resp.node_id,
						node_name:node_name.value,
						node_data:node_url.value,
						canread:true,
						canwrite:true,
						admin:true,
						children:[],
					})
				}else{
					target.value!.children=[{
						node_id:resp.node_id,
						node_name:node_name.value,
						node_data:node_url.value,
						canread:true,
						canwrite:true,
						admin:true,
						children:[],
					}]
				}
				target.value=null
				node_name.value=""
				node_url.value=""
				state.clear_load()
			})
			break
		}
		case 'update':{
			let req:permissionAPI.UpdateNodeReq = {
				node_id:target.value!.node_id,
				new_node_name:node_name.value,
				new_node_data:node_url.value,
			}
			client.permissionClient.update_node({"Token":state.user.token},req,client.timeout,(e :permissionAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(_resp :permissionAPI.UpdateNodeResp)=>{
				node_ing.value=false
				target.value!.node_name=node_name.value
				target.value!.node_data=node_url.value
				target.value=null
				node_name.value=""
				node_url.value=""
				state.clear_load()
			})
			break
		}
		case 'del':{
			let req:permissionAPI.DelNodeReq = {
				node_id:target.value!.node_id,
			}
			client.permissionClient.del_node({"Token":state.user.token},req,client.timeout,(e :permissionAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(_resp :permissionAPI.DelNodeResp)=>{
				node_ing.value=false
				for(let i=0;i<ptarget.value!.children!.length;i++){
					if (ptarget.value!.children![i]==target.value!){
						ptarget.value!.children!.splice(i,1)
						break
					}
				}
				ptarget.value=null
				target.value=null
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
function need_create_main_menu_button():boolean{
	if(!state.project.info){
		return false
	}
	if(state.project.info.project_id![1]==1){
		//admin project don't need this
		return false
	}
	//need this button when have admin permission on this project
	return projectnodes.value!=null&&projectnodes.value.admin
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
	<va-modal v-model="project_ing" attach-element="#app" max-width="600px" hide-default-actions no-dismiss overlay-opacity="0.2" z-index="999">
		<template #default>
			<div v-if="optype=='add'">
				<va-input type="text" style="width:250px" label="New Project Name*" v-model="project_name" @keyup.enter="()=>{if(project_name!=''){project_op()}}"></va-input>
				<va-button style="width:80px;margin:0 0 0 5px" :disabled="project_name==''" @click="project_op">Add</va-button>
				<va-button style="width:80px;margin:0 0 0 5px" @click="project_name='';project_ing=false">Cancel</va-button>
			</div>
			<div v-else-if="optype=='update'">
				<va-input type="text" style="width:250px" label="New Project Name*" v-model="project_name" @keyup.enter="()=>{if(project_name!=''){project_op()}}"></va-input>
				<va-button style="width:80px;margin:0 0 0 5px" :disabled="project_name==''" @click="project_op">Update</va-button>
				<va-button style="width:80px;margin:0 0 0 5px" @click="project_name='';project_ing=false">Cancel</va-button>
			</div>
			<div v-else-if="optype=='del'" style="display:flex;flex-direction:column">
				<va-card color="primary" gradient style="margin:0 0 5px 0">
					<va-card-title>Warning</va-card-title>
					<va-card-content>
						<p>You are deleting project: {{ state.project.info!.project_name}}.</p>
						<p>All data in this project will be deleted.</p>
						<p>Please confirm!</p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="project_op" gradient>Del</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="project_ing=false" gradient>Cancel</va-button>
				</div>
			</div>
		</template>
	</va-modal>
	<va-modal v-model="node_ing" attach-element="#app" max-width="600px" hide-default-actions no-dismiss overlay-opacity="0.2" z-index="999">
		<template #default>
			<div v-if="optype=='del'" style="display:flex;flex-direction:column">
				<va-card color="primary" gradient style="margin:0 0 5px 0">
					<va-card-title>Warning</va-card-title>
					<va-card-content>
						<p>You are deleting node: {{ target!.node_name }}.</p>
						<p>All data in this node will be deleted.</p>
						<p>Please confirm!</p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="node_op" gradient>Del</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="node_ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='add'" style="display:flex;flex-direction:column">
				<va-input style="width:400px;margin:3px 0" label="New Node Name*" v-model="node_name" @keyup.enter="()=>{if(node_name!=''){node_op()}}"></va-input>
				<va-input style="width:400px;margin:3px 0" label="New Node Url?" v-model="node_url" @keyup.enter="()=>{if(node_name!=''){node_op()}}"></va-input>
				<div>
					<va-button style="width:80px;margin:3px 3px 0 234px" :disabled="node_name==''" @click="node_op">Add</va-button>
					<va-button style="width:80px;margin:3px 0 0 3px" @click="node_ing=false">Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='update'" style="display:flex;flex-direction:column">
				<va-input style="width:400px;margin:3px 0" label="New Node Name*" v-model="node_name" @keyup.enter="()=>{if(node_name!=''){node_op()}}"></va-input>
				<va-input style="width:400px;margin:3px 0" label="New Node Url?" v-model="node_url" @keyup.enter="()=>{if(node_name!=''){node_op()}}"></va-input>
				<div>
					<va-button style="width:80px;margin:3px 3px 0 234px" :disabled="node_name==''" @click="node_op">Update</va-button>
					<va-button style="width:80px;margin:3px 0 0 3px" @click="node_ing=false">Cancel</va-button>
				</div>
			</div>
		</template>
	</va-modal>
	<div style="height:100%;flex:1;display:flex;flex-direction:column">
		<div style="display:flex;padding:5px 0">
			<va-select
				v-model="state.project.info"
				:options="allprojects"
				noOptionsText="NO Projects"
				label="Select Project"
				dropdownIcon=""
				style="flex:1;margin:0 2px"
				outline
				textBy="project_name"
			>
				<template #option='{option}'>
					<va-hover stateful @click="
						state.project.info=option;
						select_project();
					">
						<template #default="{hover}">
							<div
								style="padding:10px;cursor:pointer"
								:style="{'background-color':hover?'var(--va-background-border)':'',color:state.project.info==option?'green':'black'}"
							>
								{{option.project_name}}
							</div>
						</template>
					</va-hover>
				</template>
			</va-select>
			<va-dropdown v-if="state.user.root" style="width:36px;margin-right:2px">
				<template #anchor>
					<va-button>•••</va-button>
				</template>
				<va-dropdown-content>
					<va-popover message="Create New Project" :hover-out-timeout="0" :hover-over-timeout="0" color="primary" prevent-overflow>
						<va-button v-if="state.user.root" style="width:36px;margin:0 3px" @click="optype='add';project_name='';project_ing=true">+</va-button>
					</va-popover>
					<va-popover message="Rename Project" :hover-out-timeout="0" :hover-over-timeout="0" color="primary" prevent-overflow>
						<va-button v-if="state.user.root&&state.project.info&&!same_node_id(state.project.info.project_id!,[0,1])" style="width:36px;margin:0 3px" @click="optype='update';project_name='';project_ing=true">◉</va-button>
					</va-popover>
					<va-popover message="Delete Project" :hover-out-timeout="0" :hover-over-timeout="0" color="primary" prevent-overflow>
						<va-button v-if="state.user.root&&state.project.info&&!same_node_id(state.project.info.project_id!,[0,1])" style="width:36px;margin:0 3px" @click="optype='del';project_ing=true">x</va-button>
					</va-popover>
				</va-dropdown-content>
			</va-dropdown>
		</div>
		<div v-if="need_create_main_menu_button()" style="text-align:center;background-color:var(--va-background-element)">
			<va-popover message="Create Main Menu" color="primary" :hover-over-timeout="0" :hover-out-timeout="0" placement="right" prevent-overflow>
				<va-hover stateful>
					<template #default="{hover}">
						<div
							style="padding:10px 15px;cursor:pointer"
							:style="{'background-color':hover?'var(--va-shadow)':undefined}"
							@click="node_name='';node_url='';optype='add';ptarget=null;target=projectnodes;node_ing=true"
						>
							<b>+</b>
						</div>
					</template>
				</va-hover>
			</va-popover>
		</div>
		<div style="flex:1;overflow-x:hidden;overflow-y:auto;background-color:var(--va-background-element)">
			<menutree
				v-if="Boolean(projectnodes)&&Boolean(projectnodes!.children)&&projectnodes!.children!.length>0"
				:pnode="projectnodes!"
				:deep="0"
				@nodeevent="(pnode,node,type)=>{
					node_name='';
					node_url='';
					ptarget=pnode;
					target=node;
					optype=type;
					node_ing=true
				}"
			/>
		</div>
	</div>

</template>
