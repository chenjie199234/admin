<script setup lang="ts">
import {ref} from 'vue'

import * as initializeAPI from './api/admin_initialize_browser_toc'
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
	let req=new initializeAPI.UpdateRootPasswordReq()
	req.old_password=oldpassword.value
	req.new_password=newpassword.value
	client.initializeClient.update_root_password({"Token":state.user.token},req,client.timeout,(e: initializeAPI.LogicError)=>{
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
	},(_resp: initializeAPI.UpdateRootPasswordResp)=>{
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
	<VaModal v-model="password_changing" :mobileFullscreen="false" hideDefaultActions noDismiss blur :overlay="false" maxWidth="800px" @beforeOpen="(el)=>{el.querySelector('.va-modal__inner').style.minWidth='0px'}">
		<template #default>
			<div style="display:flex;flex-direction:column">
				<VaCard style="min-width:350px;width:auto;text-align:center" color="primary" gradient>
					<VaCardContent style="font-size:20px"><b>Change Root Password</b></VaCardContent>
				</VaCard>
				<VaInput :type="t_oldpassword?'text':'password'" label="Old Root Password*" v-model="oldpassword" style="margin-top:10px">
					<template #appendInner>
						<VaIcon :name="t_oldpassword?'◎':'◉'" size="small" color="var(--va-primary)" @click="t_oldpassword=!t_oldpassword" />
					</template>
				</VaInput>
				<VaInput :type="t_newpassword?'text':'password'" label="New Root Password*" v-model="newpassword" style="margin-top:10px">
					<template #appendInner>
						<VaIcon :name="t_newpassword?'◎':'◉'" size="small" color="var(--va-primary)" @click="t_newpassword=!t_newpassword" />
					</template>
				</VaInput>
				<div style="display:flex;justify-content:center">
					<VaButton style="width:100px;margin:10px 10px 0 0" :disabled="!change_root_password_able()" @click="do_change_root_password">Change</VaButton>
					<VaButton style="width:100px;margin:10px 0 0 10px" @click="oldpassword='';newpassword='';t_oldpassword=false;t_newpassword=false;password_changing=false">Cancel</VaButton>
				</div>
			</div>
		</template>
	</VaModal>
	<div style="height:100%;flex:1;display:flex;flex-direction:column;overflow:auto">
		<div style="display:flex;align-items:center;padding:5px;background-color:var(--va-background-element)">
			<div style="display:flex;flex:1"></div>
			<p v-if="!state.user.root&&state.user.info" style="color:green;margin-right:10px">{{state.user.info.user_id}}</p>
			<VaDropdown style="width:36px" trigger="hover" :hoverOverTimeout="0" :hoverOutTimeout="100" placement="bottom-end">
				<template #anchor>
					<VaButton round>{{ state.avatar() }}</VaButton>
				</template>
				<VaDropdownContent>
					<div style="display:flex;flex-direction:column">
						<VaButton v-if="state.user.root" style="margin:0 0 3px 0" @click="password_changing=true">ChangePassword</VaButton>
						<VaButton style="margin:3px 0 0 0" @click="state.logout">Logout</VaButton>
					</div>
				</VaDropdownContent>
			</VaDropdown>
		</div>
		<userrole v-if="Boolean(state.page.node)&&Boolean(state.page.node!.node_id)&&state.page.node!.node_id!.length==3&&state.page.node!.node_id![2]==1"></userrole>
		<app v-else-if="Boolean(state.page.node)&&Boolean(state.page.node!.node_id)&&state.page.node!.node_id!.length==3&&state.page.node!.node_id![2]==2"></app>
		<iframe v-else-if="Boolean(state.page.node)&&state.page.node!.node_data!=''" width="100%" height="100%" frameborder="0" :src="state.page.node!.node_data" @load="iframeload"></iframe>
	</div>
</template>
