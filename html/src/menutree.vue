<script setup lang="ts">
import {ref} from 'vue'
import * as permissionAPI from './api/admin_permission_browser_toc'

import * as state from './state'

defineProps<{
	pnode:permissionAPI.NodeInfo
	deep:number
}>()

const open=ref<{[k:string]:boolean}>({})
const hovernode=ref<permissionAPI.NodeInfo|null>(null)
function bindstyle(node :permissionAPI.NodeInfo){
	let style={}
	if(node==hovernode.value && jumpable(node)){
		style["background-color"] = "var(--va-shadow)"
	}else if(state.page.node==node){
		style["background-color"] = "#b6d7a8"
	}
	if(jumpable(node)){
		style["cursor"] = "pointer"
	}
	return style
}
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
		<VaDivider vertical color="shadow" style="margin-right:0" :style="{'margin-left':deep==0?'5px':'15px'}" />
		<div style="flex:1;display:flex;flex-direction:column">
			<template v-for="node of pnode.children">
				<div v-if="showable(node)" style="display:flex;align-items:center">
					<div
						style="flex:1;display:flex;align-items:center"
						:style="bindstyle(node!)"
						@mouseenter="hovernode=node!"
						@mouseover="hovernode=node!"
						@mouseout="hovernode=null"
						@mouseleave="hovernode=null"
						@click="()=>{
							open[node!.node_id!.toString()]=!open[node!.node_id!.toString()]
							if(jumpable(node)){
								state.set_page(node!)
							}
						}">
						<VaDivider color="shadow" style="margin:0;width:15px" />
						<div style="flex:1;display:flex;align-items:center;padding:3px 0">
							<span style="padding:7px 0">{{node!.node_name}}</span>
							<span v-if="jumpable(node)" style="padding-left:5px;font-size:30px" :style="{color:state.page.node==node?'green':'black'}">☞</span>
						</div>
						<VaHover v-if="has_children(node)" stateful>
							<template #default="{hover}">
								<div
									style="padding:5px;border-radius:2px"
									:style="{'background-color':hover?'var(--va-shadow)':undefined}"
									@click.stop="open[node!.node_id!.toString()]=!open[node!.node_id!.toString()]"
								>
									{{open[node!.node_id!.toString()]?'▲':'▼'}}
								</div>
							</template>
						</VaHover>
						<VaHover v-if="need_button(node)" stateful>
							<template #default="{hover}">
								<div v-if="!hover" style="padding:5px 7px">
									•••
								</div>
								<VaPopover
									v-if="hover"
									message="Update Menu"
									:hover-out-timeout="0"
									:hover-over-timeout="0"
									color="primary">
									<VaHover stateful>
										<template #default="{hover}">
											<div
												style="padding:5px 7px;border-radius:2px"
												:style="{'background-color':hover?'var(--va-shadow)':undefined}"
												@click.stop="$emit('nodeevent',pnode,node,'update')">
												<b>◉</b>
											</div>
										</template>
									</VaHover>
								</VaPopover>
								<VaPopover
									v-if="hover"
									message="Delete Menu"
									:hover-out-timeout="0"
									:hover-over-timeout="0"
									color="primary">
									<VaHover stateful>
										<template #default="{hover}">
											<div
												style="padding:5px 9px;border-radius:2px"
												:style="{'background-color':hover?'var(--va-shadow)':undefined}"
												@click.stop="$emit('nodeevent',pnode,node,'del')">
												<b>x</b>
											</div>
										</template>
									</VaHover>
								</VaPopover>
								<VaPopover
									v-if="hover"
									message="Add Sub Menu"
									:hover-out-timeout="0"
									:hover-over-timeout="0"
									color="primary">
									<VaHover stateful>
										<template #default="{hover}">
											<div
												style="padding:5px 9px;border-radius:2px"
												:style="{'background-color':hover?'var(--va-shadow)':undefined}"
												@click.stop="$emit('nodeevent',pnode,node,'add')">
												<b>+</b>
											</div>
										</template>
									</VaHover>
								</VaPopover>
							</template>
						</VaHover>
					</div>
				</div>
				<menutree v-if="showable(node)&&open[node!.node_id!.toString()]&&has_children(node)" :pnode="node!" :deep="deep+1" @nodeevent="(pnode,node,type)=>{$emit('nodeevent',pnode,node,type)}"></menutree>
			</template>
		</div>
	</div>
</template>
