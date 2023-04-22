<script setup lang="ts">
const props=defineProps<{
	pnode:permissionAPI.NodeInfo
}>()
</script>
<template>
	<div style="display:flex">
		<div style="witdh:100px">
		</div>
		<div v-if="Boolean(nodes)&&nodes.length>0" style="flex:1;display:flex;flex-direction:column">
			<template v-for="node of nodes">
				<nodetree :pnode="node" :nodes="node.children" />
			</template>
		</div>
	</div>
</template>
