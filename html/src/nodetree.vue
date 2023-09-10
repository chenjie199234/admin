<script setup lang="ts">
import { ref,onMounted } from 'vue'
import * as permissionAPI from './api/admin_permission_browser_toc'
const props=defineProps<{
	pnode:permissionAPI.NodeInfo
	deep:number
}>()
const new_canread=ref<boolean>(false)
const new_canwrite=ref<boolean>(false)
const new_admin=ref<boolean>(false)
onMounted(()=>{
	new_canread.value=props.pnode.canread
	new_canwrite.value=props.pnode.canwrite
	new_admin.value=props.pnode.admin
})
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
	return new_canread.value==props.pnode.canread&&new_canwrite.value==props.pnode.canwrite&&new_admin.value==props.pnode.admin
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
				v-model="new_canread"
				true-inner-label="Read"
				false-inner-label="Read"
				@update:model-value="permission_update('read')"
			/>
			<va-switch
				off-color="shadow"
				style="margin:2px"
				v-model="new_canwrite"
				true-inner-label="Write"
				false-inner-label="Write"
				@update:model-value="permission_update('write')"
			/>
			<va-switch
				off-color="shadow"
				style="margin:2px"
				v-model="new_admin"
				true-inner-label="Admin"
				false-inner-label="Admin"
				@update:model-value="permission_update('admin')"
			/>
			<va-button
				:disabled="permission_same()"
				style="margin:2px"
				@click="$emit('permissionevent',pnode,new_canread,new_canwrite,new_admin)"
			>
				Update
			</va-button>
			<div
				v-if="Boolean(pnode.children)&&pnode.children!.length>0"
				style="padding:5px 10px;border-radius:3px;cursor:pointer"
				:style="{'background-color':hover?'var(--va-shadow)':''}"
				@mouseover="hover=true"
				@mouseout="hover=false"
				@click="open=!open"
			>
				{{open?'▲':'▼'}}
			</div>
		</div>
		<div v-if="open&&Boolean(pnode.children)&&pnode.children!.length>0" style="flex:1;display:flex">
			<template v-for="child of pnode.children">
				<nodetree v-if="Boolean(child)" :pnode="child!" :deep="deep+1" @permissionevent="(updatenode,r,w,a)=>{$emit('permissionevent',updatenode,r,w,a)}"/>
			</template>
		</div>
	</div>
</template>
