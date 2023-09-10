<script setup lang="ts">
import {ref,onMounted} from  'vue'
import * as initializeAPI from './api/admin_initialize_browser_toc'
import * as state from './state'
import * as client from './client'

onMounted(()=>{
	if(!state.set_load()){
		return
	}
	client.initializeClient.init_status({},{},client.timeout,(e: initializeAPI.Error)=>{
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
	},(resp: initializeAPI.InitStatusResp)=>{
		state.clear_load()
		state.inited.value=resp.status
	})
})

const access_key=ref<string>("")
const t_access_key=ref<boolean>(false)
const password=ref<string>("")
const t_password=ref<boolean>(false)
function init_able():boolean{
	return access_key.value!="" && password.value.length>=10 && password.value.length<32
}
function do_init(){
	if(!init_able()){
		if(access_key.value){
			state.set_alert("error",-2,"Root Password length must in [10,32)!")
		}else{
			state.set_alert("error",-2,"Missing Access Key!")
		}
		return
	}
	if(!state.set_load()){
		return
	}
	client.initializeClient.init({"Access-Key":access_key.value},{password:password.value},client.timeout,(e: initializeAPI.Error)=>{
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
	},(_resp: initializeAPI.InitResp)=>{
		state.clear_load()
		access_key.value=""
		password.value=""
		state.inited.value=true
	})
}
</script>
<template>
	<div style="width:100%;height:100%;display:flex;flex-direction:column;justify-content:center;align-items:center">
		<va-card style="width:400px;margin:10px" color="primary" gradient>
			<va-card-title>Warning</va-card-title>
			<va-card-content>System not initialized.</va-card-content>
		</va-card>
		<div style="display:flex;flex-direction:column">
			<va-input :type="t_access_key?'text':'password'" label="Access Key*" v-model="access_key" style="width:400px;margin:5px 0" @keyup.enter="()=>{if(init_able()){do_init()}}">
				<template #appendInner>
					<va-icon :name="t_access_key?'◎':'◉'" size="small" color="var(--va-primary)" @click="t_access_key=!t_access_key" />
				</template>
			</va-input>
			<va-input :type="t_password?'text':'password'" label="Root Password*" v-model="password" style="width:400px;margin:5px 0" @keyup.enter="()=>{if(init_able()){do_init()}}">
				<template #appendInner>
					<va-icon :name="t_password?'◎':'◉'" size="small" color="var(--va-primary)" @click="t_password=!t_password" />
				</template>
			</va-input>
			<va-button style="width:100px;margin:5px 0 0 300px" :disabled="!init_able()" @click="do_init" gradient>Init</va-button>
		</div>
	</div>
</template>
