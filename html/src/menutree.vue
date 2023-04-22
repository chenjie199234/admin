<script setup lang="ts">
import * as permissionAPI from '../../api/permission_browser_toc'

import * as state from './state'

const emit = defineEmits(['nodeevent'])

const props=defineProps<{
	pnode:permissionAPI.NodeInfo
	deep:number
}>()

function need_button(node: permissionAPI.NodeInfo):boolean{
	if(node.node_id.length>=3&&(node.node_id[2]==1||node.node_id[2]==2)){
		//system node don't need button
		return false
	}
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
	return node.canread&&node.node_data!=''
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
	<div style="display:flex">
		<va-divider vertical style="margin-right:0" :style="{'margin-left':deep==0?'5px':'15px'}" />
		<div style="flex:1;display:flex;flex-direction:column">
			<template v-for="node of pnode.children">
				<div v-if="showable(node)" style="display:flex;align-items:center">
					<div
						style="flex:1;display:flex;align-items:center"
						:style="{'background-color':node.hover?'var(--va-shadow)':undefined,cursor:jumpable(node)?'pointer':'default'}"
						@mouseover="node.hover=true"
						@mouseout="node.hover=false"
						@click="()=>{
							node.open=!node.open
							if(jumpable(node)){
								state.set_page(node)
							}
						}"
					>
						<va-divider style="margin:0;width:15px" />
						<div style="flex:1;padding:10px 0">{{node.node_name}}  {{state.page.node==node?'☞':''}}</div>
						<div v-if="has_children(node)" style="margin-right:5px;padding:5px;border-radius:2px">{{node.open?'▲':'▼'}}</div>
					</div>
					<va-dropdown v-if="need_button(node)" trigger="hover" style="width:36px;height:36px;margin:2px" prevent-overflow>
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
				<menutree v-if="showable(node)&&node.open&&has_children(node)" :pnode="node" :deep="deep+1" @nodeevent="(pnode,node,type)=>{$emit('nodeevent',pnode,node,type)}"></menutree>
			</template>
		</div>
	</div>
</template>
