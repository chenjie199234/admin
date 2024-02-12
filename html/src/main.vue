<script setup lang="ts">
import * as state from './state'

import init from './init.vue'
import login from './login.vue'
import pageleft from './pageleft.vue'
import pageright from './pageright.vue'
</script>

<template>
	<VaModal v-model="state.load.ing" :mobileFullscreen="false" hideDefaultActions noDismiss blur :overlay="false" noPadding @beforeOpen="(el)=>{el.querySelector('.va-modal__inner').style.minWidth='0px'}">
		<VaInnerLoading icon="❃" loading :size="60" style="width:100px;height:100px"/>
	</VaModal>
	<VaModal v-model="state.alert.ing" :title="state.get_alert_title()" :message="state.alert.msg" :mobileFullscreen="false" hideDefaultActions fixedLayout blur :overlay="false" maxWidth="600px" maxHeight="400px" />
	<init v-if="!state.inited.value" />
	<login v-else-if="state.user.token.length==0" />
	<VaSplit v-else style="width:100%;height:100%;display:flex" stateful :model-value='0' :limits="['300px',50]">
		<template #start>
			<pageleft />
		</template>
		<template #end>
			<pageright />
		</template>
	</VaSplit>
</template>
