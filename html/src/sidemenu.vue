<script setup lang="ts">
import * as permissionAPI from '../../api/permission_browser_toc'
import * as state from './state'

defineProps<{
	nodes:permissionAPI.NodeInfo[]
	deep:number
}>()

function node_need_button(node: permissionAPI.NodeInfo):boolean{
	return node.admin&&node.node_id.length>=3&&node.node_id[2]!=1&&node.node_id[2]!=2&&node.node_id[2]!=3
}
function node_has_children(node: permissionAPI.NodeInfo):boolean{
	if(node.node_id.length>=3&&(node.node_id[2]==1||node.node_id[2]==2||node.node_id[2]==3)){
		//system node hide children
		return false
	}
	return node.children&&node.children.length>0
}
function system_node_child(node: permissionAPI.NodeInfo):boolean{
	return node.node_id.length>3&&(node.node_id[2]==1||node.node_id[2]==2||node.node_id[2]==3)
}
function jumpable(node: permissionAPI.NodeInfo):boolean{
	return (node.node_id.length==3&&(node.node_id[2]==1||node.node_id[2]==2||node.node_id[2]==3))||node.node_data!=''
}
</script>
<template>
	<div v-for="node of nodes" style="display:flex;flex-direction:column">
		<div v-if="!system_node_child(node)" style="display:flex;align-items:center">
			<div
			style="height:38px;padding-left:10px;display:flex;align-items:center;flex:1;border-radius:5px"
			:style="{'margin-left':30*deep+2+'px',cursor:jumpable(node)?'pointer':'default','background-color':node.labelhover?'var(--va-shadow)':''}"
			@mouseover="node.labelhover=true"
			@mouseout="node.labelhover=false"
			@click="()=>{if(jumpable(node)){state.set_page(node)}}">
				{{ node.node_name }}
			</div>
			<div
			v-if="node_has_children(node)"
			style="width:30px;height:38px;cursor:pointer;display:flex;justify-content:center;align-items:center;border-radius:5px"
			:style="{'background-color':node.iconhover?'var(--va-shadow)':''}"
			@mouseover="node.iconhover=true"
			@mouseout="node.iconhover=false"
			@click="node.open=!node.open">
				{{ node.open?'▲':'▼' }}
			</div>
			<va-dropdown v-if="node_need_button(node)" trigger="hover" style="width:36px;height:36px;margin:2px" prevent-overflow>
				<template #anchor>
					<va-button>•••</va-button>
				</template>
				<va-dropdown-content>
					<va-popover message="Update Menu" :hover-out-timeout="0" :hover-over-timeout="0" color="primary">
						<va-button style="width:36px;margin:0 3px" @click="state.set_node(node,'update')">◉</va-button>
					</va-popover>
					<va-popover message="Delete Menu" :hover-out-timeout="0" :hover-over-timeout="0" color="primary">
						<va-button style="width:36px;margin:0 3px" @click="state.set_node(node,'del')">x</va-button>
					</va-popover>
					<va-popover message="Add Sub Menu" :hover-out-timeout="0" :hover-over-timeout="0" color="primary">
						<va-button style="width:36px;margin:0 3px" @click="state.set_node(node,'add')">✿</va-button>
					</va-popover>
				</va-dropdown-content>
			</va-dropdown>
		</div>
		<sidemenu v-if="!system_node_child(node)&&node.open&&node_has_children(node)" :nodes="node.children" :deep="deep+1"></sidemenu>
	</div>
</template>
