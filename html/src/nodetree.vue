<script setup lang="ts">
import { ref,watch } from 'vue'
import * as permissionAPI from './api/admin_permission_browser'
const props=defineProps<{
	pnode:permissionAPI.NodeInfo|null
	node:permissionAPI.NodeInfo
	deep:number
	disabled:boolean
}>()
const new_canread=ref<boolean>(false)
const new_canwrite=ref<boolean>(false)
const new_admin=ref<boolean>(false)
watch(()=>props.node.admin,(newval,oldval)=>{
	new_canread.value=props.node.canread
	new_canwrite.value=props.node.canwrite
	new_admin.value=props.node.admin
	if(newval){
		if(!props.node.children){
			return
		}
		for(let child of props.node.children){
			if(!child){
				continue
			}
			child.canread=true
			child.canwrite=true
			child.admin=true
		}
	}else if(oldval){
		if(!props.node.children){
			return
		}
		for(let child of props.node.children){
			if(!child){
				continue
			}
			child.canread=false
			child.canwrite=false
			child.admin=false
		}
	}
},{immediate: true})
const open=ref<boolean>(false)
const hover=ref<boolean>(false)
function permission_update(t :string){
	switch(t){
		case "read":{
			if(!new_canread.value){
				new_canwrite.value=false
				new_admin.value=false
			}
			break
		}
		case "write":{
			if(!new_canwrite.value){
				new_admin.value=false
			}else{
				new_canread.value=true
			}
			break
		}
		case "admin":{
			if(new_admin.value){
				new_canread.value=true
				new_canwrite.value=true
			}
			break
		}
	}
}
function permission_same():boolean{
	return new_canread.value==props.node.canread&&new_canwrite.value==props.node.canwrite&&new_admin.value==props.node.admin
}
</script>
<template>
	<div style="flex:1;display:flex;flex-direction:column;margin:0 15px">
		<VaDivider v-if="deep!=0" vertical style="height:50px;align-self:center;border-right-color:var(--va-primary)"/>
		<div style="display:flex;flex-direction:column;align-items:center;padding:5px 15px 0px 15px;border:1px solid var(--va-primary);border-radius:5px">
			<div style="margin:2px;min-width:100px;border:1px solid var(--va-primary);border-radius:3px;padding:10px;text-align:center">{{node.node_name}}</div>
			<VaSwitch
				:disabled="disabled||(pnode!=null&&pnode.admin)"
				off-color="shadow"
				style="margin:2px"
				v-model="new_canread"
				true-inner-label="Read"
				false-inner-label="Read"
				@update:modelValue="permission_update('read')"
			/>
			<VaSwitch
				:disabled="disabled||(pnode!=null&&pnode.admin)"
				off-color="shadow"
				style="margin:2px"
				v-model="new_canwrite"
				true-inner-label="Write"
				false-inner-label="Write"
				@update:modelValue="permission_update('write')"
			/>
			<VaSwitch
				:disabled="disabled||(pnode!=null&&pnode.admin)"
				off-color="shadow"
				style="margin:2px"
				v-model="new_admin"
				true-inner-label="Admin"
				false-inner-label="Admin"
				@update:modelValue="permission_update('admin')"
			/>
			<VaButton
				:disabled="disabled||(pnode!=null&&pnode.admin)||permission_same()"
				style="margin:2px"
				@click="$emit('permissionevent',node,new_canread,new_canwrite,new_admin)"
			>
				Update
			</VaButton>
			<div
				v-if="Boolean(node.children)&&node.children!.length>0"
				style="padding:5px 10px;border-radius:3px;cursor:pointer"
				:style="{'background-color':hover?'var(--va-shadow)':''}"
				@mouseover="hover=true"
				@mouseout="hover=false"
				@click="open=!open"
			>
				{{open?'▲':'▼'}}
			</div>
		</div>
		<div v-if="open&&node.children&&node.children!.length>0" style="flex:1;display:flex">
			<template v-for="child of node.children">
				<nodetree v-if="child"
					:pnode="node"
					:node="child!"
					:deep="deep+1"
					:disabled="disabled"
					@permissionevent="(updatenode,r,w,a)=>{$emit('permissionevent',updatenode,r,w,a)}"/>
			</template>
		</div>
	</div>
</template>
