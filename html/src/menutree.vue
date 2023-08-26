<script setup lang="ts">
import {ref} from 'vue'
import * as permissionAPI from './api/permission_browser_toc'

import * as state from './state'

defineProps<{
	pnode:permissionAPI.NodeInfo
	deep:number
}>()

const open=ref<{[k:string]:boolean}>({})
const hover=ref<string>("")

function need_button(node: permissionAPI.NodeInfo|null|undefined):boolean{
	if(!node){
		return false
	}
	if(!node.node_id){
		return false
	}
	if(node.node_id.length>=3&&(node.node_id[2]==1||node.node_id[2]==2)){
		//system node don't need button
		return false
	}
	return node.admin
}
function has_children(node: permissionAPI.NodeInfo|null|undefined):boolean{
	if(!node){
		return false
	}
	if(!node.node_id){
		return false
	}
	if(node.node_id.length>=3&&(node.node_id[2]==1||node.node_id[2]==2)){
		//system node hide children
		return false
	}
	return node.children!=null&&node.children!=undefined&&node.children.length>0
}
function jumpable(node: permissionAPI.NodeInfo|null|undefined):boolean{
	if(!node){
		return false
	}
	if(!node.node_id){
		return false
	}
	if(node.node_id.length==3&&(node.node_id[2]==1||node.node_id[2]==2)){
		//system node can jump
		return true
	}
	return node.canread&&node.node_data!=''
}
function showable(node: permissionAPI.NodeInfo|null|undefined):boolean{
	if(!node){
		return false
	}
	if(!node.node_id){
		return false
	}
	if(node.node_id.length>3&&(node.node_id[2]==1||node.node_id[2]==2)){
		//system node's child need to be hide
		//but system node self need to be show
		return false
	}
	return true
}
</script>
<template>
	<div style="display:flex">
		<va-divider vertical color="shadow" style="margin-right:0" :style="{'margin-left':deep==0?'5px':'15px'}" />
		<div style="flex:1;display:flex;flex-direction:column">
			<template v-for="node of pnode.children">
				<div v-if="showable(node)" style="display:flex;align-items:center">
					<div
						style="flex:1;display:flex;align-items:center"
						:style="{'background-color':hover==node!.node_id!.toString()?'var(--va-shadow)':state.page.node==node?'#b6d7a8':undefined,cursor:jumpable(node)?'pointer':'default'}"
						@mouseover="hover=node!.node_id!.toString()"
						@mouseout="hover=''"
						@click="()=>{
							open[node!.node_id!.toString()]=!open[node!.node_id!.toString()]
							if(jumpable(node)){
								state.set_page(node!)
							}
						}">
						<va-divider color="shadow" style="margin:0;width:15px" />
						<div style="flex:1;display:flex;align-items:center;padding:3px 0">
							<span style="padding:7px 0">{{node!.node_name}}</span>
							<span v-if="jumpable(node)" style="padding-left:5px;font-size:30px" :style="{color:state.page.node==node?'green':'black'}">☞</span>
						</div>
						<div v-if="has_children(node)" style="margin-right:5px;padding:5px;border-radius:2px">{{open[node!.node_id!.toString()]?'▲':'▼'}}</div>
					</div>
					<va-dropdown v-if="need_button(node)" style="width:36px;height:36px;margin:2px">
						<template #anchor>
							<va-button>•••</va-button>
						</template>
						<va-dropdown-content>
							<va-popover message="Add Sub Menu" :hover-out-timeout="0" :hover-over-timeout="0" color="primary" prevent-overflow>
								<va-button style="width:36px;margin:0 3px" @click="$emit('nodeevent',pnode,node,'add')">+</va-button>
							</va-popover>
							<va-popover message="Update Menu" :hover-out-timeout="0" :hover-over-timeout="0" color="primary" prevent-overflow>
								<va-button style="width:36px;margin:0 3px" @click="$emit('nodeevent',pnode,node,'update')">◉</va-button>
							</va-popover>
							<va-popover message="Delete Menu" :hover-out-timeout="0" :hover-over-timeout="0" color="primary" prevent-overflow>
								<va-button style="width:36px;margin:0 3px" @click="$emit('nodeevent',pnode,node,'del')">x</va-button>
							</va-popover>
						</va-dropdown-content>
					</va-dropdown>
				</div>
				<menutree v-if="showable(node)&&open[node!.node_id!.toString()]&&has_children(node)" :pnode="node!" :deep="deep+1" @nodeevent="(pnode,node,type)=>{$emit('nodeevent',pnode,node,type)}"></menutree>
			</template>
		</div>
	</div>
</template>
