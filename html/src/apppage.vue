<script setup lang="ts">
import { ref,computed,onMounted } from 'vue'
import * as appAPI from '../../api/app_browser_toc'
import * as state from './state'
import * as client from './client'

const all=ref<{[k:string]:string[]}>({})
onMounted(()=>{
	if(!state.page.node.children){
		return
	}
	for(let n of state.page.node.children){
		let pieces:string[] = n.node_name.split(".")
		if(all.value[pieces[0]]){
			all.value[pieces[0]].push(pieces[1])
		}else{
			all.value[pieces[0]]=[pieces[1]]
		}
	}
})
const curg=ref<string>("")
const cura=ref<string>("")
const secret=ref<string>("")
const t_secret=ref<boolean>(false)

const keys=ref<Map<string,KeyConfigInfo>>(new Map())
const t_keys=ref<boolean>(false)
const t_keys_hover=ref<boolean>(false)

const proxys=ref<Map<string,ProxyPathInfo>>(new Map())
const t_proxys=ref<boolean>(false)
const t_proxys_hover=ref<boolean>(false)

function get_app(){
	if(curg.value==""||cura.value==""){
		keys.value=null
		proxys.value=null
		state.set_error("error",-2,"Group and App must be selected!")
		return
	}
	if(!state.set_load()){
		return
	}
	client.appClient.get_app({"Token":state.user.token},{g_name:curg.value,a_name:cura.value,secret:secret.value},client.timeout,(e: appAPI.Error)=>{
		state.clear_load()
		state.set_error("error",e.code,e.msg)
	},(resp: appAPI.GetAppResp)=>{
		state.clear_load()
		if(resp.keys){
			keys.value = new Map([...resp.keys.entries()].sort())
		}else{
			keys.value = new Map()
		}
		t_keys.value=false
		if(resp.proxys){
			proxys.value = new Map([...resp.proxys.entries()].sort())
		}else{
			proxys.value = new Map()
		}
		t_proxys.value=false
	})
}

const ing=ref<boolean>(false)
const optype=ref<string>("")
const new_g=ref<string>("")
const new_a=ref<string>("")
const new_secret=ref<string>("")
const update_g=ref<string>("")
const update_a=ref<string>("")
const update_old_secret=ref<string>("")
const update_new_secret=ref<string>("")
function app_op(){
	if(!state.set_load()){
		return
	}
	switch(optype.value){
		case 'del_app':{
			if(curg.value==""||cura.value==""){
				state.clear_load()
				state.set_error("error",-2,"Group and App must be selected!")
				return
			}
			client.appClient.del_app({"Token":state.user.token},{g_name:curg.value,a_name:cura.value,secret:secret.value},client.timeout,(e: appAPI.Error)=>{
				state.clear_load()
				state.set_error("error",e.code,e.msg)
			},(resp: appAPI.GetAppResp)=>{
				state.clear_load()
				let index=all.value[curg.value].indexOf(cura.value)
				if(index!=-1){
					all.value[curg.value].splice(index,1)
				}
				if(all.value[curg.value].length==0){
					delete all.value[curg.value]
				}
				curg.value=""
				cura.value=""
				secret.value=""
				keys.value=new Map()
				proxys.value=new Map()
				ing.value=false
			})
			break
		}
		case 'add_app':{
			client.appClient.create_app({"Token":state.user.token},{project_id:state.project.cur.project_id,g_name:new_g.value,a_name:new_a.value,secret:new_secret.value},client.timeout,(e: appAPI.Error)=>{
				state.clear_load()
				state.set_error("error",e.code,e.msg)
			},(resp: appAPI.CreateAppResp)=>{
				state.clear_load()
				if(all.value[new_g.value]){
					all.value[new_g.value].push(new_a.value)
				}else{
					all.value[new_g.value]=[new_a.value]
				}
				new_g.value=""
				new_a.value=""
				new_secret.value=""
				ing.value=false
			})
			break
		}
		case 'update_secret':{
			let req = {
				g_name:update_g.value,
				a_name:update_a.value,
				old_secret:update_old_secret.value,
				new_secret:update_new_secret.value
			}
			client.appClient.update_app_secret({"Token":state.user.token},req,client.timeout,(e: appAPI.Error)=>{
				state.clear_load()
				state.set_error("error",e.code,e.msg)
			},(resp: appAPI.UpdateAppSecretResp)=>{
				state.clear_load()
				update_g.value=""
				update_a.value=""
				update_old_secret.value=""
				update_new_secret.value=""
				ing.value=false
			})
			break
		}
		default:{
			state.clear_load()
			state.set_error("error",-2,"unknown operation")
		}
	}
}

</script>
<template>
	<va-modal v-model="ing" attach-element="#app" max-width="600px" hide-default-actions no-dismiss overlay-opacity="0.2" z-index="999">
		<template #default>
			<div v-if="optype=='del_app'" style="display:flex;flex-direction:column">
				<va-card color="primary" gradient style="margin:0 0 5px 0">
					<va-card-title>Warning</va-card-title>
					<va-card-content>
						<p>You are deleting app: {{ cura }} in group: {{ curg }}.</p>
						<p>All data in this app will be deleted.</p>
						<p>Please confirm!</p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="app_op" gradient>Del</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='add_app'" style="display:flex;flex-direction:column">
				<va-input type="text" label="Group*" style="margin:2px" v-model="new_g" />
				<va-input type="text" label="App*" style="margin:2px" v-model="new_a" />
				<va-input type="text" label="Secret" style="margin:2px" v-model="new_secret" :max-length="31" />
				<div style="display:flex;justify-content:center">
					<va-button @click="app_op" style="margin:5px">Add</va-button>
					<va-button @click="new_g='';new_a='';new_secret='';ing=false" style="margin:5px">Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='update_secret'" style="display:flex;flex-direction:column">
				<va-select 
					trigger="hover"
					dropdown-icon=""
					label="Group*"
					:options="Object.keys(all)"
					style="width:400px;margin:2px"
					v-model="update_g"
					no-options-text="No Groups"
					@update:model-value="update_a=''"
				/>
				<va-select
					trigger="hover"
					dropdown-icon=""
					label="App*"
					:options="all[update_g]"
					style="width:400px;margin:2px"
					v-model="update_a"
					no-options-text="No Apps"
				/>
				<va-input type="text" label="Old Secret" style="width:400px;margin:2px" v-model="update_old_secret" />
				<va-input type="text" label="New Secret" style="width:400px;margin:2px" v-model="update_new_secret" />
				<div style="display:flex;justify-content:center">
					<va-button @click="app_op" style="margin:5px">Update</va-button>
					<va-button @click="update_g='';update_a='';update_old_secret='';update_new_secret='';ing=false" style="margin:5px">Cancel</va-button>
				</div>
			</div>
		</template>
	</va-modal>
	<div style="display:flex;flex:1;flex-direction:column;margin:1px;width:100%;overflow-y:auto">
		<div style="width:100%;display:flex;margin:1px 0">
			<va-select
				dropdown-icon=""
				outline
				trigger="hover"
				label="Group*"
				no-options-text="No Groups"
				:options="Object.keys(all)"
				v-model="curg"
				style="width:250px;margin-right:1px"
				@update:model-value="cura=''"
			>
				<template #appendInner>
				</template>
			</va-select>
			<va-select
				dropdown-icon=""
				outline
				trigger="hover"
				label="App*"
				no-options-text="No Apps"
				:options="all[curg]"
				v-model="cura"
				style="width:250px;margin:0 1px"
			>
				<template #appendInner>
				</template>
			</va-select>
			<va-input :type="t_secret?'text':'password'" v-model="secret" outline label="Secret" :max-length="31" style="min-width:250px;max-width:250px;margin:0 1px">
				<template #appendInner>
					<va-icon :name="t_secret?'◎':'◉'" size="small" color="var(--va-primary)" @click="t_secret=!t_secret" />
				</template>
			</va-input>
			<va-button style="margin:0 2px" :disabled="curg==''||cura==''" @click="get_app">Search</va-button>
			<va-dropdown  v-if="state.page.node.admin" trigger="hover" style="width:36px;margin-right:4px">
				<template #anchor>
					<va-button>•••</va-button>
				</template>
				<va-dropdown-content>
					<va-popover message="Create New App" :hover-out-timeout="0" :hover-over-timeout="0" color="primary">
						<va-button style="width:36px;margin:0 3px" @click="ing=true;optype='add_app'">+</va-button>
					</va-popover>
					<va-popover message="Update Add Secret" :hover-out-timeout="0" :hover-over-timeout="0" color="primary">
						<va-button style="width:36px;margin:0 3px" @click="ing=true;optype='update_secret'">◉</va-button>
					</va-popover>
					<va-popover message="Delete App" :hover-out-timeout="0" :hover-over-timeout="0" color="primary">
						<va-button style="width:36px;margin:0 3px" :disabled="curg==''||cura==''" @click="ing=true;optype='del_app'">x</va-button>
					</va-popover>
				</va-dropdown-content>
			</va-dropdown>
		</div>
		<!-- configs -->
		<div
			style="width:100%;display:flex;align-items:center;margin:1px 0;cursor:pointer"
			:style="{'background-color':t_keys_hover?'var(--va-shadow)':'var(--va-background-element)'}"
			@click="t_keys=!t_keys"
			@mouseover="t_keys_hover=true"
			@mouseout="t_keys_hover=false"
		>
			<span style="flex:1;padding:12px;color:var(--va-primary)">Configs</span>
			<va-button style="height:30px" size="small" @mouseover.stop="" @mouseout.stop="" @click.stop="">ADD</va-button>
			<span style="width:60px;padding:12px 20px;color:var(--va-primary)">{{ t_keys?'▲':'▼' }}</span>
		</div>
		<!-- keys -->
		 <div v-if="t_keys" style="overflow-y:auto;height:auto;max-height:100%">
			<div v-for="key of keys.keys()" style="margin:1px 0 1px 20px;display:flex;flex-direction:column">
				<div
					style="cursor:pointer;display:flex;align-items:center"
					:style="{'background-color':keys.get(key).hover?'var(--va-shadow)':'var(--va-background-element)'}"
					@click="keys.get(key).open=!keys.get(key).open"
					@mouseover="keys.get(key).hover=true"
					@mouseout="keys.get(key).hover=false"
				>
					<span style="width:35px;padding:12px;color:var(--va-primary)"> {{ keys.get(key).open?'▼':'►' }} </span>
					<span style="padding:12px;color:var(--va-primary)">{{key}}</span>
					<span style="flex:1"></span>
					<span style="padding:12px;color:var(--va-primary)">Current Config ID:  {{ keys.get(key).cur_index }}</span>
					<va-dropdown trigger="hover" style="width:60px;margin-right:40px" prevent-overflow placement="bottom-end">
						<template #anchor>
							<va-button style="width:60px;height:30px" size="small" @mouseover.stop="" @mouseout.stop="" @click.stop="">History</va-button>
						</template>
						<va-dropdown-content>
							<div style="max-height:300px;overflow-y:auto;display:flex;flex-direction:column">
								<va-button
									v-for="index of keys.get(key).max_index"
									size="small"
									style="height:24px;width:42px;padding:5px 0;margin:1px;cursor:pointer"
									:disabled="keys.get(key).cur_index==keys.get(key).max_index-index+1"
								>
									{{keys.get(key).max_index-index+1}}
								</va-button>
							</div>
						</va-dropdown-content>
					</va-dropdown>
				</div>
				<div v-if="keys.get(key).open" style="display:flex;margin:2px 0 0 20px">
					<div style="flex:1;display:flex;flex-direction:column;align-items:center">
						<va-input
							:model-value="JSON.stringify(JSON.parse(keys.get(key).cur_value),null,4)"
							style="width:100%"
							type="textarea"
							outline
							readonly
							:min-rows="15"
							:max-rows="15"
						/>
						<va-button
							v-if="keys.get(key).open"
							style="margin:4px"
							:disabled="keys.get(key).new_cur_value"
							@click="()=>{
								if(keys.get(key).cur_value){
									keys.get(key).new_cur_value=JSON.stringify(JSON.parse(keys.get(key).cur_value),null,4)
								}else{
									keys.get(key).new_cur_value='{}'
								}
							}"
						>
							Edit
						</va-button>
					</div>
					<va-divider v-if="keys.get(key).new_cur_value" vertical />
					<div v-if="keys.get(key).new_cur_value" style="flex:1;display:flex;flex-direction:column;align-items:center">
						<va-input
							v-if="keys.get(key).open"
							v-model="keys.get(key).new_cur_value"
							style="width:100%"
							type="textarea"
							outline
							:min-rows="15"
							:max-rows="15"
						/>
						<div>
							<va-button>Update</va-button>
							<va-button style="margin:4px" @click="keys.get(key).new_cur_value=undefined">Cancel</va-button>
						</div>
					</div>
				</div>
			</div>
		</div>
		<!-- proxys -->
		<div
			style="width:100%;display:flex;align-items:center;margin:1px 0;cursor:pointer"
			:style="{'background-color':t_proxys_hover?'var(--va-shadow)':'var(--va-background-element)'}"
			@click="t_proxys=!t_proxys"
			@mouseover="t_proxys_hover=true"
			@mouseout="t_proxys_hover=false"
		>
			<span style="flex:1;padding:12px;color:var(--va-primary)">Proxys</span>
			<va-button style="height:30px" size="small" @mouseover.stop="" @mouseout.stop="" @click.stop="">ADD</va-button>
			<span style="width:60px;padding:12px 20px;color:var(--va-primary)">{{ t_proxys?'▲':'▼' }}</span>
		</div>
		<!-- paths -->
		<div v-if="t_proxys" style="overflow-y:auto;height:auto;max-height:100%">
			<div v-for="proxy of proxys.keys()" style="margin:1px 0 1px 20px;display:flex;flex-direction:column">
				<div
					style="cursor:pointer;padding:12px"
					:style="{'background-color':proxys.get(proxy).hover?'var(--va-shadow)':'var(--va-background-element)'}"
					@click="proxys.get(proxy).open=!proxys.get(proxy).open"
					@mouseover="proxys.get(proxy).hover=true"
					@mouseout="proxys.get(proxy).hover=false"
				>
					<span style="width:35px;padding:12px;color:var(--va-primary)"> {{ proxys.get(proxy).open?'▼':'►' }} </span>
					<span style="padding:12px 0;color:var(--va-primary)">{{proxy}}</span>
				</div>
			</div>
		</div>
	</div>
</template>
