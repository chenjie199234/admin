<script setup lang="ts">
import {ref,onMounted} from 'vue'
import * as state from './state'

import {NodeInfo} from '@api/browser_message_admin_NodeInfo'
import {ProjectInfo} from '@api/browser_message_admin_ProjectInfo'

import {ListProjectReq} from '@api/browser_message_admin_ListProjectReq'
import {ListProjectResp} from '@api/browser_message_admin_ListProjectResp'
import {ListProject} from '@api/browser_method_admin_Initialize_ListProject'

import {ListUserNodeReq} from '@api/browser_message_admin_ListUserNodeReq'
import {ListUserNodeResp} from '@api/browser_message_admin_ListUserNodeResp'
import {ListUserNode} from '@api/browser_method_admin_Permission_ListUserNode'

import {CreateProjectReq} from '@api/browser_message_admin_CreateProjectReq'
import {CreateProjectResp} from '@api/browser_message_admin_CreateProjectResp'
import {CreateProject} from '@api/browser_method_admin_Initialize_CreateProject'

import {UpdateProjectReq} from '@api/browser_message_admin_UpdateProjectReq'
import {UpdateProjectResp} from '@api/browser_message_admin_UpdateProjectResp'
import {UpdateProject} from '@api/browser_method_admin_Initialize_UpdateProject'

import {DeleteProjectReq} from '@api/browser_message_admin_DeleteProjectReq'
import {DeleteProjectResp} from '@api/browser_message_admin_DeleteProjectResp'
import {DeleteProject} from '@api/browser_method_admin_Initialize_DeleteProject'

import {AddNodeReq} from '@api/browser_message_admin_AddNodeReq'
import {AddNodeResp} from '@api/browser_message_admin_AddNodeResp'
import {AddNode} from '@api/browser_method_admin_Permission_AddNode'

import {UpdateNodeReq} from '@api/browser_message_admin_UpdateNodeReq'
import {UpdateNodeResp} from '@api/browser_message_admin_UpdateNodeResp'
import {UpdateNode} from '@api/browser_method_admin_Permission_UpdateNode'

import {DelNodeReq} from '@api/browser_message_admin_DelNodeReq'
import {DelNodeResp} from '@api/browser_message_admin_DelNodeResp'
import {DelNode} from '@api/browser_method_admin_Permission_DelNode'

import menutree from './menutree.vue'

onMounted(()=>{
  get_projects()
})

const allprojects=ref<ProjectInfo[]>([])

function get_projects(){
  if(!state.set_load()){
	return
  }
  let req = new ListProjectReq()
  ListProject(state.baseurl,req,{"Token":state.user.token},state.expire).then(resp => {
	if(resp instanceof ListProjectResp){
	  if(resp.projects){
		let tmp: ProjectInfo[] = []
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
	}else{
	  state.clear_load()
	  state.set_alert("error",resp.code,resp.msg)
	}
  })
}

function selfproject():boolean{
  if(!state.project.info){
	return false
  }
  if(!state.project.info.project_id){
	return false
  }
  return state.project.info!.project_id![0]==0&&state.project.info!.project_id![1]==1
}

const projectnodes=ref<NodeInfo|null>(null)

function select_project(){
  if(!state.project.info){
	return
  }
  if(!state.set_load()){
	return
  }
  let req = new ListUserNodeReq()
  req.project_id=state.project.info!.project_id
  req.user_id=""
  req.need_user_role_node=true
  ListUserNode(state.baseurl,req,{"Token":state.user.token},state.expire).then(resp => {
	if(resp instanceof ListUserNodeResp){
	  if(resp.node){
		projectnodes.value=resp.node!
	  }else{
		projectnodes.value=null
	  }
	  state.clear_page()
	  state.clear_load()
	}else{
	  projectnodes.value=null
	  state.clear_project()
	  state.clear_page()
	  state.clear_load()
	  state.set_alert("error",resp.code,resp.msg)
	}
  })
}

const node_ing=ref<boolean>(false)
const project_ing=ref<boolean>(false)
const optype=ref<string>("")

const project_name=ref<string>("")

function project_update_able():boolean{
  return project_name.value!=state.project.info!.project_name
}

function project_op(){
  if(!state.set_load()){
	return
  }
  switch(optype.value){
	case 'add':{
	  let req = new CreateProjectReq()
	  req.project_name=project_name.value
	  req.project_data=""
	  CreateProject(state.baseurl,req,{"Token":state.user.token},state.expire).then(resp => {
		if(resp instanceof CreateProjectResp){
		  project_name.value=''
		  project_ing.value=false
		  state.clear_load()
		  get_projects()
		}else{
		  state.clear_load()
		  state.set_alert("error",resp.code,resp.msg)
		}
	  })
	  break
	}
	case 'update':{
	  let req = new UpdateProjectReq()
	  req.project_id=state.project.info!.project_id
	  req.new_project_name=project_name.value
	  req.new_project_data=""
	  UpdateProject(state.baseurl,req,{"Token":state.user.token},state.expire).then(resp => {
		if(resp instanceof UpdateProjectResp){
		  project_name.value=''
		  project_ing.value=false
		  state.clear_load()
		  get_projects()
		}else{
		  state.clear_load()
		  state.set_alert("error",resp.code,resp.msg)
		}
	  })
	  break
	}
	case 'del':{
	  let req = new DeleteProjectReq()
	  req.project_id=state.project.info!.project_id
	  DeleteProject(state.baseurl,req,{"Token":state.user.token},state.expire).then(resp => {
		if(resp instanceof DeleteProjectResp){
		  project_ing.value=false
		  state.clear_project()
		  projectnodes.value=null
		  state.clear_load()
		  get_projects()
		}else{
		  state.clear_load()
		  state.set_alert("error",resp.code,resp.msg)
		}
	  })
	  break
	}
	default:{
	  state.clear_load()
	  state.set_alert("error",-2,"unknown operation")
	}
  }
}

const ptarget=ref<NodeInfo|null>(null)
const target=ref<NodeInfo|null>(null)
const node_name=ref<string>("")
const node_url=ref<string>("")

function node_update_able():boolean{
  if(!target.value){
	return false
  }
  return node_name.value!=target.value!.node_name || node_url.value!=target.value!.node_data
}

function node_op(){
  if(!state.project.info){
	return
  }
  if(!state.set_load()){
	return
  }
  switch(optype.value){
	case 'add':{
	  let req=new AddNodeReq() 
	  req.pnode_id=target.value!.node_id
	  req.node_name=node_name.value
	  req.node_data=node_url.value
	  AddNode(state.baseurl,req,{"Token":state.user.token},state.expire).then(resp => {
		if(resp instanceof AddNodeResp){
		  node_ing.value=false
		  let tmp = new NodeInfo()
		  tmp.node_id=resp.node_id
		  tmp.node_name=node_name.value
		  tmp.node_data=node_url.value
		  tmp.canread=true
		  tmp.canwrite=true
		  tmp.admin=true
		  tmp.children=[]
		  if(!target.value!.children){
			target.value!.children=[]
		  }
		  target.value!.children!.push(tmp)
		  target.value=null
		  node_name.value=""
		  node_url.value=""
		  state.clear_load()
		}else{
		  state.clear_load()
		  state.set_alert("error",resp.code,resp.msg)
		}
	  })
	  break
	}
	case 'update':{
	  let req=new UpdateNodeReq()
	  req.node_id=target.value!.node_id
	  req.new_node_name=node_name.value
	  req.new_node_data=node_url.value
	  UpdateNode(state.baseurl,req,{"Token":state.user.token},state.expire).then(resp => {
		if(resp instanceof UpdateNodeResp){
		  node_ing.value=false
		  target.value!.node_name=node_name.value
		  target.value!.node_data=node_url.value
		  target.value=null
		  node_name.value=""
		  node_url.value=""
		  state.clear_load()
		}else{
		  state.clear_load()
		  state.set_alert("error",resp.code,resp.msg)
		}
	  })
	  break
	}
	case 'del':{
	  let req=new DelNodeReq()
	  req.node_id=target.value!.node_id
	  DelNode(state.baseurl,req,{"Token":state.user.token},state.expire).then(resp => {
		if(resp instanceof DelNodeResp){
		  node_ing.value=false
		  for(let i=0;i<ptarget.value!.children!.length;i++){
			if (ptarget.value!.children![i]==target.value!){
			  ptarget.value!.children!.splice(i,1)
			  break
			}
		  }
		  if(state.page.node){
			let needcleanpage=true
			for(let i=0;i<target.value!.node_id!.length;i++){
			  if(target.value!.node_id![i]!=state.page.node!.node_id![i]){
				needcleanpage=true
				break
			  }
			}
			if(needcleanpage){
			  state.clear_page()
			}
		  }
		  ptarget.value=null
		  target.value=null
		  state.clear_load()
		}else{
		  state.clear_load()
		  state.set_alert("error",resp.code,resp.msg)
		}
	  })
	  break
	}
	default:{
	  state.clear_load()
	  state.set_alert("error",-2,"unknown operation")
	}
  }
}
function create_main_menu_able():boolean{
  if(!state.project.info){
	return false
  }
  if(state.project.info.project_id![1]==1){
	//admin project don't need this
	return false
  }
  //need this button when have admin permission on this project
  return projectnodes.value!=null&&projectnodes.value.admin!
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
	<!-- <VaModal v-model="project_ing" :mobileFullscreen="false" hideDefaultActions noDismiss blur :overlay="false" maxWidth="800px" @beforeOpen="(el:HTMLElement)=>{el.querySelector('.va-modal__dialog').style.width='auto'}"> -->
	<VaModal v-model="project_ing" :mobileFullscreen="false" hideDefaultActions noDismiss blur :overlay="false" maxWidth="800px">
		<template #default>
			<div v-if="optype=='add'" style="display:flex;flex-direction:column">
				<VaCard style="min-width:350px;width:auto;text-align:center" color="primary" gradient>
					<VaCardContent style="font-size:20px"><b>Create Project</b></VaCardContent>
				</VaCard>
				<VaInput type="text" style="margin-top:10px" label="New Project Name*" v-model.trim="project_name" />
				<div style="display:flex;justify-content:center">
					<VaButton style="width:80px;margin:10px 10px 0 0" :disabled="project_name==''" @click="project_op" gradient>Add</VaButton>
					<VaButton style="width:80px;margin:10px 0 0 10px" @click="project_name='';project_ing=false" gradient>Cancel</VaButton>
				</div>
			</div>
			<div v-else-if="optype=='update'" style="display:flex;flex-direction:column">
				<VaCard style="min-width:350px;width:auto;text-align:center" color="primary" gradient>
					<VaCardContent style="font-size:20px"><b>Update Project</b></VaCardContent>
				</VaCard>
				<VaInput type="text" style="margin-top:10px" label="New Project Name*" v-model.trim="project_name" />
				<div style="display:flex;justify-content:center">
					<VaButton style="width:80px;margin:10px 10px 0 0" :disabled="!project_update_able()" @click="project_op" gradient>Update</VaButton>
					<VaButton style="width:80px;margin:10px 0 0 10px" @click="project_name='';project_ing=false" gradient>Cancel</VaButton>
				</div>
			</div>
			<div v-else-if="optype=='del'" style="display:flex;flex-direction:column">
				<VaCard style="min-width:350px;width:auto;text-align:center" color="primary" gradient>
					<VaCardContent style="font-size:20px">
						<p><b>Delete project: {{ state.project.info!.project_name }}</b></p>
						<p><b>Please confirm</b></p>
					</VaCardContent>
				</VaCard>
				<div style="display:flex;justify-content:center">
					<VaButton style="width:80px;margin:10px 10px 0 0" @click="project_op" gradient>Del</VaButton>
					<VaButton style="width:80px;margin:10px 0 0 10px" @click="project_ing=false" gradient>Cancel</VaButton>
				</div>
			</div>
		</template>
	</VaModal>
	<!-- <VaModal v-model="node_ing" :mobileFullscreen="false" hideDefaultActions noDismiss blur :overlay="false" maxWidth="800px" @beforeOpen="(el:HTMLElement)=>{el.querySelector('.va-modal__dialog').style.width='auto'}"> -->
	<VaModal v-model="node_ing" :mobileFullscreen="false" hideDefaultActions noDismiss blur :overlay="false" maxWidth="800px">
		<template #default>
			<div v-if="optype=='del'" style="display:flex;flex-direction:column">
				<VaCard style="min-width:350px;width:auto;text-align:center" color="primary" gradient>
					<VaCardContent style="font-size:20px">
						<p><b>Delete node: {{ target!.node_name }}</b></p>
						<p><b>Please confirm</b></p>
					</VaCardContent>
				</VaCard>
				<div style="display:flex;justify-content:center">
					<VaButton style="width:80px;margin:10px 10px 0 0" @click="node_op" gradient>Del</VaButton>
					<VaButton style="width:80px;margin:10px 0 0 10px" @click="node_ing=false" gradient>Cancel</VaButton>
				</div>
			</div>
			<div v-else-if="optype=='add'" style="display:flex;flex-direction:column">
				<VaCard style="min-width:350px;width:auto;text-align:center" color="primary" gradient>
					<VaCardContent style="font-size:20px"><b>Create Node</b></VaCardContent>
				</VaCard>
				<VaInput style="margin-top:10px" label="New Node Name*" v-model.trim="node_name" />
				<VaInput style="margin-top:10px" label="New Node Url?" v-model.trim="node_url" />
				<div style="display:flex;justify-content:center">
					<VaButton style="width:80px;margin:10px 10px 0 0" :disabled="node_name==''" @click="node_op" gradient>Add</VaButton>
					<VaButton style="width:80px;margin:10px 0 0 10px" @click="node_ing=false" gradient>Cancel</VaButton>
				</div>
			</div>
			<div v-else-if="optype=='update'" style="display:flex;flex-direction:column">
				<VaCard style="min-width:350px;width:auto;text-align:center" color="primary" gradient>
					<VaCardContent style="font-size:20px"><b>Update Node</b></VaCardContent>
				</VaCard>
				<VaInput style="margin-top:10px" label="New Node Name*" v-model.trim="node_name" />
				<VaInput style="margin-top:10px" label="New Node Url?" v-model.trim="node_url" />
				<div style="display:flex;justify-content:center">
					<VaButton style="width:80px;margin:10px 10px 0 0" :disabled="!node_update_able()" @click="node_op" gradient>Update</VaButton>
					<VaButton style="width:80px;margin:10px 0 0 10px" @click="node_ing=false" gradient>Cancel</VaButton>
				</div>
			</div>
		</template>
	</VaModal>
	<div style="height:100%;flex:1;display:flex;flex-direction:column">
		<div style="display:flex;padding:5px 0">
			<VaSelect
				v-model="state.project.info"
				textBy="project_name"
				:options="allprojects"
				noOptionsText="NO Projects"
				placeholder="Select Project*"
				dropdownIcon=""
				style="flex:1;padding-left:5px"
				outline
				trigger="hover"
				:hoverOverTimeout="0"
				:hoverOutTimeout="100"
				@update:modelValue="select_project()"
			>
				<template #option='{option,selectOption}'>
					<VaHover stateful @click="selectOption(option)">
						<template #default="{hover}">
							<div
								style="padding:10px;cursor:pointer"
								:style="{'background-color':hover?'var(--va-background-border)':'',color:state.project.info&&option&&state.project.info.project_name==(option as ProjectInfo)['project_name']?'green':'black'}"
							>
								{{option?(option as ProjectInfo)['project_name']:''}}
							</div>
						</template>
					</VaHover>
				</template>
			</VaSelect>
			<VaPopover v-if="state.user.root" message="Create New Project" :hover-out-timeout="0" :hover-over-timeout="0" color="primary">
				<VaButton style="margin:0 3px" @click="optype='add';project_name='';project_ing=true">+</VaButton>
			</VaPopover>
		</div>
		<div v-if="state.project.info" style="text-align:center;background-color:var(--va-background-element)">
			<VaPopover v-if="state.user.root&&!selfproject()" message="Update Project" color="primary" :hover-over-timeout="0" :hover-out-timeout="0">
				<VaHover stateful>
					<template #default="{hover}">
						<div
							style="padding:10px 15px;cursor:pointer"
							:style="{'background-color':hover?'var(--va-shadow)':undefined}"
							@click="optype='update';project_name=state.project.info.project_name!;project_ing=true"
						>
							<b>◉</b>
						</div>
					</template>
				</VaHover>
			</VaPopover>
			<VaPopover v-if="state.user.root&&!selfproject()" message="Delete Project" color="primary" :hover-over-timeout="0" :hover-out-timeout="0">
				<VaHover stateful>
					<template #default="{hover}">
						<div
							style="padding:10px 15px;cursor:pointer"
							:style="{'background-color':hover?'var(--va-shadow)':undefined}"
							@click="optype='del';project_ing=true"
						>
							<b>x</b>
						</div>
					</template>
				</VaHover>
			</VaPopover>
			<VaPopover v-if="create_main_menu_able()" message="Create Main Menu" color="primary" :hover-over-timeout="0" :hover-out-timeout="0">
				<VaHover stateful>
					<template #default="{hover}">
						<div
							style="padding:10px 15px;cursor:pointer"
							:style="{'background-color':hover?'var(--va-shadow)':undefined}"
							@click="node_name='';node_url='';optype='add';ptarget=null;target=projectnodes;node_ing=true"
						>
							<b>+</b>
						</div>
					</template>
				</VaHover>
			</VaPopover>
		</div>
		<div style="flex:1;overflow-x:hidden;overflow-y:auto;background-color:var(--va-background-element)">
			<menutree
				v-if="Boolean(projectnodes)&&Boolean(projectnodes!.children)&&projectnodes!.children!.length>0"
				:pnode="projectnodes!"
				:deep="0"
				@nodeevent="(pnode:NodeInfo,node:NodeInfo,type:string)=>{
					if(type=='update'){
						node_name=node.node_name!;
						node_url=node.node_data!;
					}else{
						node_name='';
						node_url='';
					}
					ptarget=pnode;
					target=node;
					optype=type;
					node_ing=true;
				}"
			/>
		</div>
	</div>
</template>
