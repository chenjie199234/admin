<script setup lang="ts">
import {ref} from 'vue'

import * as initializeAPI from 'admin/api/initialize_browser_toc'
import * as state from './state'
import * as client from './client'

import app from './app.vue'
import userrole from './userrole.vue'

const password_changing = ref<boolean>(false)
const oldpassword = ref<string>("")
const t_oldpassword = ref<boolean>(false)
const newpassword = ref<string>("")
const t_newpassword = ref<boolean>(false)
function change_root_password_able():boolean{
	return oldpassword.value.length>=10 && oldpassword.value.length<32 && newpassword.value.length>=10 && newpassword.value.length<32
}
function do_change_root_password(){
	if(!state.user.root || state.user.token==""){
		return
	}
	if(!change_root_password_able()){
		state.set_alert("error",-2,"Root Password length must in [10,32)!")
		return
	}
	if(!state.set_load()){
		return
	}
	client.initializeClient.root_password({"Token":state.user.token},{old_password:oldpassword.value,new_password:newpassword.value},client.timeout,(e: initializeAPI.Error)=>{
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
	},(_resp: initializeAPI.RootPasswordResp)=>{
		oldpassword.value=""
		newpassword.value=""
		t_oldpassword.value=false
		t_newpassword.value=false
		password_changing.value=false
		state.logout()
		state.clear_load()
	})
}
function iframeload(){
	console.log("iframe")
}
</script>
<template>
	<va-modal v-model="password_changing" attach-element="#app" max-width="600px" hide-default-actions no-dismiss overlay-opacity="0.2" z-index="999">
		<template #default>
			<div style="display:flex;flex-direction:column">
				<va-input :type="t_oldpassword?'text':'password'" label="Old Root Password*" v-model="oldpassword" style="width:400px;margin:5px 0" @keyup.enter="()=>{if(change_root_password_able()){do_change_root_password()}}">
					<template #appendInner>
						<va-icon :name="t_oldpassword?'◎':'◉'" size="small" color="var(--va-primary)" @click="t_oldpassword=!t_oldpassword" />
					</template>
				</va-input>
				<va-input :type="t_newpassword?'text':'password'" label="New Root Password*" v-model="newpassword" style="width:400px;margin:5px 0" @keyup.enter="()=>{if(change_root_password_able()){do_change_root_password()}}">
					<template #appendInner>
						<va-icon :name="t_newpassword?'◎':'◉'" size="small" color="var(--va-primary)" @click="t_newpassword=!t_newpassword" />
					</template>
				</va-input>
				<div>
					<va-button style="width:100px;margin:5px 5px 0 190px" :disabled="!change_root_password_able()" @click="do_change_root_password">Change</va-button>
					<va-button style="width:100px;margin:5px 0 0 5px" @click="oldpassword='';newpassword='';t_oldpassword=false;t_newpassword=false;password_changing=false">Cancel</va-button>
				</div>
			</div>
		</template>
	</va-modal>
	<div style="height:100%;flex:1;display:flex;flex-direction:column;overflow:auto">
		<div style="display:flex;padding:5px;background-color:var(--va-background-element)">
			<div style="display:flex;flex:1"></div>
			<va-dropdown trigger="hover" :hover-out-timeout="60000" style="width:36px" placement="bottom-end" prevent-overflow>
				<template #anchor>
					<va-button round>{{ state.avatar() }}</va-button>
				</template>
				<va-dropdown-content>
					<div style="display:flex;flex-direction:column">
						<va-button v-if="state.user.root" style="margin:0 0 3px 0" @click="password_changing=true">ChangePassword</va-button>
						<va-button v-if="state.user.root" style="margin:3px 0 0 0" @click="state.logout">Logout</va-button>
						<va-button v-if="!state.user.root" @click="state.logout">Logout</va-button>
					</div>
				</va-dropdown-content>
			</va-dropdown>
		</div>
		<userrole v-if="Boolean(state.page.node)&&Boolean(state.page.node!.node_id)&&state.page.node!.node_id!.length==3&&state.page.node!.node_id![2]==1"></userrole>
		<app v-else-if="Boolean(state.page.node)&&Boolean(state.page.node!.node_id)&&state.page.node!.node_id!.length==3&&state.page.node!.node_id![2]==2"></app>
		<iframe v-else-if="Boolean(state.page.node)&&state.page.node!.node_data!=''" width="100%" height="100%" frameborder="0" :src="state.page.node!.node_data" @load="iframeload"></iframe>
	</div>
</template>
