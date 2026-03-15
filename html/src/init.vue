<script setup lang="ts">
import {ref,onMounted} from 'vue'
import * as state from './state.ts'

import {InitStatusReq} from '@api/browser_message_admin_InitStatusReq'
import {InitStatusResp} from '@api/browser_message_admin_InitStatusResp'
import {InitStatus} from '@api/browser_method_admin_Initialize_InitStatus'

import {InitReq} from '@api/browser_message_admin_InitReq'
import {InitResp} from '@api/browser_message_admin_InitResp'
import {Init} from '@api/browser_method_admin_Initialize_Init'

onMounted(()=>{
  if(!state.set_load()){
	return
  }
  let req=new InitStatusReq()
  InitStatus(state.baseurl,req,{},state.expire).then(resp => {
	if(resp instanceof InitStatusResp){
	  state.clear_load()
	  state.inited.value=resp.status!
	}else{
	  state.clear_load()
	  state.set_alert("error",resp.code,resp.msg)
	}
  })
})

const access_key=ref<string>("")
const t_access_key=ref<boolean>(false)
const password=ref<string>("")
const t_password=ref<boolean>(false)
function init_able():boolean{
  return access_key.value!="" && password.value.length>=10 && password.value.length<=32
}
function do_init(){
  if(!init_able()){
	if(access_key.value){
	  state.set_alert("error",-2,"Root Password length must in [10,32]!")
	}else{
	  state.set_alert("error",-2,"Missing Access Key!")
	}
	return
  }
  if(!state.set_load()){
	return
  }
  let req = new InitReq()
  req.password=password.value
  Init(state.baseurl,req,{"Access-Key":access_key.value},state.expire).then(resp => {
	if(resp instanceof InitResp){
	  state.clear_load()
	  access_key.value=""
	  password.value=""
	  state.inited.value=true
	}else{
	  state.clear_load()
	  state.set_alert("error",resp.code,resp.msg)
	}
  })
}
</script>
<template>
	<div v-if="state.inited.value==false" style="width:100%;height:100%;display:flex;flex-direction:column;justify-content:center;align-items:center">
		<VaCard style="min-width:350px;width:auto;text-align:center" color="primary" gradient>
			<VaCardContent style="font-size:20px"><b>Initialize Now</b></VaCardContent>
		</VaCard>
		<VaInput :type="t_access_key?'text':'password'" label="Access Key*" v-model.trim="access_key" style="min-width:350px;margin-top:10px">
			<template #appendInner>
				<VaIcon :name="t_access_key?'◎':'◉'" size="small" color="var(--va-primary)" @click="t_access_key=!t_access_key" />
			</template>
		</VaInput>
		<VaInput :type="t_password?'text':'password'" label="Root Password*" v-model="password" style="min-width:350px;margin-top:10px" :max-length="32">
			<template #appendInner>
				<VaIcon :name="t_password?'◎':'◉'" size="small" color="var(--va-primary)" @click="t_password=!t_password" />
			</template>
		</VaInput>
		<VaButton style="width:100px;margin-top:10px" :disabled="!init_able()" @click="do_init" gradient>Init</VaButton>
	</div>
</template>
