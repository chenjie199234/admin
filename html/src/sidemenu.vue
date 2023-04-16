<script setup lang="ts">
import {ref} from 'vue'
import * as permissionAPI from '../../api/permission_browser_toc'

import * as state from './state'
import * as client from './client'

const props=defineProps<{
	nodes:permissionAPI.NodeInfo[]
	deep:number
}>()

const ing=ref<boolean>(false)
const optype=ref<string>("")
const target=ref<permissionAPI.NodeInfo>(null)
const new_node_name=ref<string>("")
const new_node_url=ref<string>("")

function op(){
	if(state.project.cur_id.length==0){
		return
	}
	if(!state.set_load()){
		return
	}
	switch(optype.value){
		case 'add':{
			let req:permissionAPI.AddNodeReq = {
				pnode_id:target.value.node_id,
				node_name:new_node_name.value,
				node_data:new_node_url.value,
			}
			client.permissionClient.add_node({"Token":state.user.token},req,client.timeout,(e :permissionAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp :permissionAPI.AddNodeResp)=>{
				ing.value=false
				if(target.value.children){
					target.value.children.push({
						node_id:resp.node_id,
						node_name:new_node_name.value,
						node_data:new_node_url.value,
						canread:true,
						canwrite:true,
						admin:true,
						children:[],
					})
				}else{
					target.value.children=[{
						node_id:resp.node_id,
						node_name:new_node_name.value,
						node_data:new_node_url.value,
						canread:true,
						canwrite:true,
						admin:true,
						children:[],
					}]
				}
				target.value=null
				new_node_name.value=""
				new_node_url.value=""
				state.clear_load()
			})
			break
		}
		case 'update':{
			let req:permissionAPI.UpdateNodeReq = {
				node_id:target.value.node_id,
				new_node_name:new_node_name.value,
				new_node_data:new_node_url.value,
			}
			client.permissionClient.update_node({"Token":state.user.token},req,client.timeout,(e :permissionAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp :permissionAPI.UpdateNodeResp)=>{
				ing.value=false
				target.value.node_name=new_node_name.value
				target.value.node_data=new_node_url.value
				target.value=null
				new_node_name.value=""
				new_node_url.value=""
				state.clear_load()
			})
			break
		}
		case 'del':{
			let req:permissionAPI.DelNodeReq = {
				node_id:target.value.node_id,
			}
			client.permissionClient.del_node({"Token":state.user.token},req,client.timeout,(e :permissionAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp :permissionAPI.DelNodeResp)=>{
				ing.value=false
				for(let i=0;i<props.nodes.length;i++){
					if(props.nodes[i]==target.value){
						props.nodes.splice(i,1)
						break
					}
				}
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

function need_button(node: permissionAPI.NodeInfo):boolean{
	if(node.node_id.length>=3&&(node.node_id[2]==1||node.node_id[2]==2)){
		//system node don't need button
		return false
	}
	//only admin permission node need button
	return node.admin
}
function has_children(node: permissionAPI.NodeInfo):boolean{
	if(node.node_id.length>=3&&(node.node_id[2]==1||node.node_id[2]==2)){
		//system node hide children
		return false
	}
	return Boolean(node.children)&&node.children.length>0
}
function jumpable(node: permissionAPI.NodeInfo):boolean{
	if(node.node_id.length==3&&(node.node_id[2]==1||node.node_id[2]==2)){
		//system node can jump
		return true
	}
	return node.node_data!=''
}
function showable(node: permissionAPI.NodeInfo):boolean{
	if(node.node_id.length>3&&(node.node_id[2]==1||node.node_id[2]==2)){
		//system node's child need to be hide
		//but system node self need to be show
		return false
	}
	return true
}
</script>
<template>
	<div v-if="deep==0&&nodes[0].node_id.length==state.project.cur_id.length&&nodes[0].admin" style="text-align:center">
		<va-popover message="Create Main Menu" color="primary" :hover-over-timeout="0" :hover-out-timeout="0" placement="right" prevent-overflow>
			<va-hover>
				<template #default="{hover}">
					<div
					style="padding:10px 15px;cursor:pointer"
					:style="{'background-color':hover?'var(--va-shadow)':undefined}"
					@click="optype='add';target=nodes[0];ing=true"
					>
						<b>+</b>
					</div>
				</template>
			</va-hover>
		</va-popover>
	</div>
	<va-modal v-model="ing" attach-element="#app" max-width="600px" hide-default-actions no-dismiss overlay-opacity="0.2" z-index="999">
		<template #default>
			<div v-if="optype=='del'" style="display:flex;flex-direction:column">
				<va-card color="primary" gradient style="margin:0 0 5px 0">
					<va-card-title>Warning</va-card-title>
					<va-card-content>
						<p>You are deleting node: {{ target.node_name }}.</p>
						<p>All data in this node will be deleted.</p>
						<p>Please confirm!</p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="op" gradient>Del</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='add'" style="display:flex;flex-direction:column">
				<va-input style="width:400px;margin:3px 0" label="New Node Name*" v-model="new_node_name" @keyup.enter="()=>{if(new_node_name!=''){op()}}"></va-input>
				<va-input style="width:400px;margin:3px 0" label="New Node Url?" v-model="new_node_url" @keyup.enter="()=>{if(new_node_name!=''){op()}}"></va-input>
				<div>
					<va-button style="width:80px;margin:3px 3px 0 234px" :disabled="new_node_name==''" @click="op">Add</va-button>
					<va-button style="width:80px;margin:3px 0 0 3px" @click="new_node_name='';new_node_url='';ing=false">Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='update'" style="display:flex;flex-direction:column">
				<va-input style="width:400px;margin:3px 0" label="New Node Name*" v-model="new_node_name" @keyup.enter="()=>{if(new_node_name!=''){op()}}"></va-input>
				<va-input style="width:400px;margin:3px 0" label="New Node Url?" v-model="new_node_url" @keyup.enter="()=>{if(new_node_name!=''){op()}}"></va-input>
				<div>
					<va-button style="width:80px;margin:3px 3px 0 234px" :disabled="new_node_name==''" @click="op">Update</va-button>
					<va-button style="width:80px;margin:3px 0 0 3px" @click="new_node_name='';new_node_url='';ing=false">Cancel</va-button>
				</div>
			</div>
		</template>
	</va-modal>
	<template v-for="node of nodes[0].node_id.length==state.project.cur_id.length?nodes[0].children:nodes">
		<div v-if="showable(node)" style="display:flex;align-items:center">
			<div
			style="flex:1;display:flex;align-items:center"
			:style="{'background-color':node.hover>=1?'var(--va-shadow)':undefined,cursor:jumpable(node)?'pointer':'default'}"
			@mouseover="()=>{
				if(node.hover==undefined){
					node.hover=1
				}else{
			   		node.hover+=1
				}
			}"
			@mouseout="node.hover-=1"
			@click="()=>{
				if(jumpable(node)){
					state.set_page(node)
				}
			}"
			>
				<div style="flex:1;padding:10px 2px" :style="{'padding-left':30*deep+5+'px'}">
					{{node.node_name}}
				</div>
				<div
				v-if="has_children(node)"
				style="padding:5px;cursor:pointer;margin-right:4px;border-radius:5px"
				:style="{'background-color':node.hover>=2?'var(--va-background-element)':undefined}"
				@mouseover="node.hover+=2"
				@mouseout="node.hover-=2"
				@click="node.open=!node.open"
				>
					{{ node.open?'▲':'▼' }}
				</div>
			</div>
			<va-dropdown v-if="need_button(node)" trigger="hover" style="width:36px;height:36px;margin:2px" prevent-overflow>
				<template #anchor>
					<va-button>•••</va-button>
				</template>
				<va-dropdown-content>
					<va-popover message="Add Sub Menu" :hover-out-timeout="0" :hover-over-timeout="0" color="primary" prevent-overflow>
						<va-button style="width:36px;margin:0 3px" @click="optype='add';new_node_name='';new_node_url='';target=node;ing=true">+</va-button>
					</va-popover>
					<va-popover message="Update Menu" :hover-out-timeout="0" :hover-over-timeout="0" color="primary" prevent-overflow>
						<va-button style="width:36px;margin:0 3px" @click="optype='update';new_node_name='';new_node_url='';target=node;ing=true">◉</va-button>
					</va-popover>
					<va-popover message="Delete Menu" :hover-out-timeout="0" :hover-over-timeout="0" color="primary" prevent-overflow>
						<va-button style="width:36px;margin:0 3px" @click="optype='del';target=node;ing=true">x</va-button>
					</va-popover>
				</va-dropdown-content>
			</va-dropdown>
		</div>
		<sidemenu v-if="showable(node)&&node.open&&has_children(node)" :nodes="node.children" :deep="deep+1"></sidemenu>
	</template>
</template>
