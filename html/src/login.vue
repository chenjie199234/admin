<script setup lang="ts">
import {ref,onMounted} from 'vue'

import * as initializeAPI from './api/admin_initialize_browser_toc'
import * as userAPI from './api/admin_user_browser_toc'
import * as state from './state'
import * as client from './client'

onMounted(()=>{
	let localtoken=localStorage.getItem("token")
	if(localtoken){
		var obj = JSON.parse(localtoken)
		if(!obj.root){
			if(state.set_load()){
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
			state.user.root=obj.root
			state.user.token=obj.token
		}
	}
})

const password = ref<string>("")
const t_password = ref<boolean>(false)
function login_root_able():boolean{
	return password.value.length>=10&&password.value.length<32
}
function do_login_root(){
	if(!login_root_able()){
		state.set_alert("error",-2,"Root Password length must in [10,32)!")
		return
	}
	if(!state.set_load()){
		return
	}
	client.initializeClient.root_login({},{password:password.value},client.timeout,(e: initializeAPI.Error)=>{
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
	},(resp: initializeAPI.RootLoginResp)=>{
		//clear loading in get_projects function
		password.value=""
		state.login(resp.token)
		state.clear_load()
	})
}

const oauth2 = ref("")
const oauth2s = ref(["Oauth2 Service Name 1","Oauth2 Service Name 2"])
const oauth2img = ref("")
/* TODO
function do_login_user(){

}
*/
</script>
<template>
	<div style="width:100%;height:100%;display:flex;justify-content:center;align-items:center">
		<div v-if="!state.user.root">
			<va-select
				v-model="oauth2"
				:options="oauth2s"
				noOptionsText="NO Oauth2 Login"
				label="Select Oauth2 Login"
				dropdownIcon=""
				style="width:400px"
			>
				<template #option='{option,selectOption}'>
					<va-hover stateful @click="()=>{
						if(oauth2!=option){
							selectOption(option)
						}
					}">
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
			<va-image style="width:400px;height:400px;margin:5px 0" :src="oauth2img" />
			<va-button style="width:400px;margin:0" @click="state.user.root=true">Switch To Root User Login</va-button>
		</div>
		<div v-else>
			<div>
				<va-input :type="t_password?'text':'password'" style="width:300px" label="Root Password*" v-model="password" @keyup.enter="()=>{if(login_root_able()){do_login_root()}}">
					<template #appendInner>
						<va-icon :name="t_password?'◎':'◉'" size="small" color="var(--va-primary)" @click="t_password=!t_password" />
					</template>
				</va-input>
				<va-button style="width:90px;margin:0 0 0 10px" :disabled="!login_root_able()" @click="do_login_root">Login</va-button>
			</div>
			<va-button style="width:400px;margin:10px 0 0 0" @click="state.user.root=false">Switch To Normal User Login</va-button>
		</div>
	</div>
</template>
