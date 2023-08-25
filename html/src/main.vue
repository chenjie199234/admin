<script setup lang="ts">
import * as state from './state'

import init from './init.vue'
import login from './login.vue'
import pageleft from './pageleft.vue'
import pageright from './pageright.vue'
</script>

<template>
	<va-modal v-model="state.load.ing" hideDefaultActions noDismiss :overlay="false" blur z-index="1000" :mobileFullscreen="false" backgroundColor="#0000">
		<template #default>
			<va-inner-loading icon="❃" loading :size="55" />
		</template>
	</va-modal>
	<va-modal v-model="state.alert.ing" maxWidth="600px" maxHeight="400px" fixedLayout :title="state.get_alert_title()" :message="state.alert.msg" hideDefaultActions :overlay="false" blur z-index="1000"  :mobileFullscreen="false"/>
	<init v-if="!state.inited.value" />
	<login v-else-if="state.user.token.length==0" />
	<va-split v-else style="width:100%;height:100%;display:flex" stateful :model-value='0' :limits="['250px',50]">
		<template #start>
			<pageleft />
		</template>
		<template #end>
			<pageright />
		</template>
	</va-split>
</template>
