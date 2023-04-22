<script setup lang="ts">
import * as state from './state'

import init from './init.vue'
import login from './login.vue'
import pageleft from './pageleft.vue'
import pageright from './pageright.vue'
</script>

<template>
	<va-modal v-model="state.load.ing" hide-default-actions no-dismiss :overlay="false" blur background-color="#0000" z-index="1000">
		<template #default>
			<va-inner-loading icon="❃" loading :size="55"></va-inner-loading>
		</template>
	</va-modal>
	<va-modal v-model="state.alert.ing" max-width="600px" max-height="400px" fixed-layout :title="state.get_alert_title()" :message="state.alert.msg" hide-default-actions :overlay="false" blur z-index="1000" />
	<init v-if="!state.inited.value" />
	<login v-else-if="state.user.token.length==0" />
	<va-split v-else style="width:100%;height:100%;display:flex" stateful :model-value='0' :limits="['250px',50]">
		<template #start="{containerSize}">
			<pageleft />
		</template>
		<template #end="{containerSize}">
			<pageright />
		</template>
	</va-split>
</template>
