<script setup lang="ts">
import { ref,onMounted } from 'vue'
import * as permissionAPI from '../../api/permission_browser_toc'
const props=defineProps<{
	pnode:permissionAPI.NodeInfo
	deep:number
}>()
onMounted(()=>{
	props.pnode.new_canread=props.pnode.canread
	props.pnode.new_canwrite=props.pnode.canwrite
	props.pnode.new_admin=props.pnode.admin
})
const open=ref<boolean>(false)
const hover=ref<boolean>(false)
function permission_update(t :string){
	switch(t){
		case "read":{
			if(!props.pnode.new_canread){
				props.pnode.new_canwrite=false
				props.pnode.new_admin=false
			}
			break
		}
		case "write":{
			if(!props.pnode.new_canwrite){
				props.pnode.new_admin=false
			}else{
				props.pnode.new_canread=true
			}
			break
		}
		case "admin":{
			if(props.pnode.new_admin){
				props.pnode.new_canread=true
				props.pnode.new_canwrite=true
			}
			break
		}
	}
}
function permission_same():boolean{
	return props.pnode.new_canread==props.pnode.canread&&props.pnode.new_canwrite==props.pnode.canwrite&&props.pnode.new_admin==props.pnode.admin
}
</script>
<template>
	<div style="flex:1;display:flex;flex-direction:column;margin:0 15px">
		<va-divider v-if="deep!=0" vertical style="height:50px;align-self:center;border-right-color:var(--va-primary)"/>
		<div style="display:flex;flex-direction:column;align-items:center;padding:5px 15px 0px 15px;border:1px solid var(--va-primary);border-radius:5px">
			<div style="margin:2px;min-width:100px;border:1px solid var(--va-primary);border-radius:3px;padding:10px;text-align:center">{{pnode.node_name}}</div>
			<va-switch
				off-color="shadow"
				style="margin:2px"
				v-model="pnode.new_canread"
				true-inner-label="Read"
				false-inner-label="Read"
				@update:model-value="permission_update('read')"
			/>
			<va-switch
				off-color="shadow"
				style="margin:2px"
				v-model="pnode.new_canwrite"
				true-inner-label="Write"
				false-inner-label="Write"
				@update:model-value="permission_update('write')"
			/>
			<va-switch
				off-color="shadow"
				style="margin:2px"
				v-model="pnode.new_admin"
				true-inner-label="Admin"
				false-inner-label="Admin"
				@update:model-value="permission_update('admin')"
			/>
			<va-button
				:disabled="permission_same()"
				style="margin:2px"
				@click="$emit('permissionevent',pnode,pnode.new_canread,pnode.new_canwrite,pnode.new_admin)"
			>
				Update
			</va-button>
			<div
				v-if="Boolean(pnode.children)&&pnode.children.length>0"
				style="padding:5px 10px;border-radius:3px;cursor:pointer"
				:style="{'background-color':hover?'var(--va-shadow)':''}"
				@mouseover="hover=true"
				@mouseout="hover=false"
				@click="open=!open"
			>
				{{open?'▲':'▼'}}
			</div>
		</div>
		<div v-if="open&&Boolean(pnode.children)&&pnode.children.length>0" style="flex:1;display:flex">
			<nodetree v-for="child of pnode.children" :pnode="child" :deep="deep+1" @permissionevent="(updatenode,r,w,a)=>{$emit('permissionevent',updatenode,r,w,a)}"/>
		</div>
	</div>
</template>
