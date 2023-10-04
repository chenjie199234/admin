<script setup lang="ts">
import {ref,onMounted} from 'vue'

import * as initializeAPI from './api/admin_initialize_browser_toc'
import * as userAPI from './api/admin_user_browser_toc'
import * as state from './state'
import * as client from './client'

onMounted(()=>{
	let localtoken=localStorage.getItem("token")
	if(localtoken){
		let obj = JSON.parse(localtoken)
		if(!obj.root){
			state.user.root=false
			if(!state.set_load()){
				return
			}
			client.userClient.login_info({"Token":obj.token},{},client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp :userAPI.LoginInfoResp)=>{
				state.user.token=obj.token
				if(resp.user){
					state.user.info=resp.user
				}else{
					state.user.info=null
				}
				state.clear_load()
			})
		}else{
			state.user.root=true
			state.user.token=obj.token
		}
	}else if(window.location.search){
		let querys = new URLSearchParams(window.location.search)
		switch(querys.get("state")){
		case "DingTalk":
			if(querys.get("authCode")){
				oauth2.value=querys.get("state")!
				oauth2code.value=querys.get("authCode")!
				do_login_user()
			}else{
				state.set_alert("error",-2,"missng authCode in redirect url")
			}
			break

		case "FeiShu":
			if(querys.get("code")){
				oauth2.value=querys.get("state")!
				oauth2code.value=querys.get("code")!
				do_login_user()
			}else{
				state.set_alert("error",-2,"missing code in redirect url")
			}
			break
		default:
			state.set_alert("error",-2,"unknown oauth2 state in redirect url,must be DingTalk or FeiShu")
		}
	}
})

const password = ref<string>("")
const t_password = ref<boolean>(false)
function login_root_able():boolean{
	return password.value.length>=10&&password.value.length<32
}
function do_login_root(){
	if(!state.set_load()){
		return
	}
	client.initializeClient.root_login({},{password:password.value},client.timeout,(e: initializeAPI.Error)=>{
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
	},(resp: initializeAPI.RootLoginResp)=>{
		password.value=""
		state.login(resp.token)
		state.clear_load()
	})
}

const oauth2 = ref<string>("")
const oauth2code = ref<string>("")
const oauth2s = ref<string[]>(["DingTalk","FeiShu"])
function doauth(){
	if(!state.set_load()){
		return
	}
	let req={
		src_type:oauth2.value,
	}
	client.userClient.get_oauth2({},req,client.timeout,(e :userAPI.Error)=>{
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
	},(resp :userAPI.GetOauth2Resp)=>{
		state.clear_load()
		window.location.href = resp.url
	})
}
function do_login_user(){
	if(!state.set_load()){
		return
	}
	let req = {
		src_type:oauth2.value,
		code:oauth2code.value,
	}
	client.userClient.user_login({},req,client.timeout,(e :userAPI.Error)=>{
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
	},(resp :userAPI.UserLoginResp)=>{
		state.login(resp.token)
		state.clear_load()
		window.location.href = window.location.href.slice(0,window.location.href.indexOf("?"))
	})
}
</script>

<template>
	<div style="width:100%;height:100%;display:flex;justify-content:center;align-items:center">
		<div v-if="!state.user.root">
			<va-card style="text-align:center" color="primary" gradient>
				<va-card-content style="font-size:20px"><b>Normal User Login</b></va-card-content>
			</va-card>
			<div style="display:flex;align-items:end;margin-top:20px">
				<va-select
					v-model="oauth2"
					:options="oauth2s"
					noOptionsText="NO Oauth2 Login"
					label="Select Oauth2 Login*"
					dropdownIcon=""
					trigger="hover"
					:hoverOverTimeout="0"
					:hoverOutTimeout="100"
				>
					<template #option='{option,selectOption}'>
						<va-hover stateful @click="selectOption(option)">
							<template #default="{hover}">
								<div
									style="padding:10px;cursor:pointer"
									:style="{'background-color':hover?'var(--va-background-border)':'',color:oauth2==option?'green':'black'}"
								>
									{{option}}
								</div>
							</template>
						</va-hover>
					</template>
				</va-select>
				<va-button style="width:90px;margin-left:10px" :disabled="oauth2==''" @click="doauth" gradient>Login</va-button>
			</div>
			<va-button style="width:400px;margin:10px 0 0 0" @click="state.user.root=true;oauth2=''" gradient>Switch To Root User Login</va-button>
		</div>
		<div v-if="state.user.root">
			<va-card style="text-align:center" color="primary" gradient>
				<va-card-content style="font-size:20px"><b>Root User Login</b></va-card-content>
			</va-card>
			<div style="display:flex;align-items:end;margin-top:20px">
				<va-input :type="t_password?'text':'password'" style="width:300px" label="Root Password*" v-model="password">
					<template #appendInner>
						<va-icon :name="t_password?'◎':'◉'" size="small" color="var(--va-primary)" @click="t_password=!t_password" />
					</template>
				</va-input>
				<va-button style="width:90px;margin-left:10px" :disabled="!login_root_able()" @click="do_login_root" gradient>Login</va-button>
			</div>
			<va-button style="width:400px;margin:10px 0 0 0" @click="state.user.root=false;oauth2=''" gradient>Switch To Normal User Login</va-button>
		</div>
	</div>
</template>
