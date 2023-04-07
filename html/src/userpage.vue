<script setup lang="ts">
import {ref,onMounted } from 'vue'
import * as userAPI from '../../api/user_browser_toc'
import * as state from './state'
import * as client from './client'
onMounted(()=>{
	if(state.page.node.admin){
		ranges.value=["This Project","All Projects"]
	}else{
		ranges.value=["This Project"]
	}
})
const ranges=ref<string[]>([])
const range=ref<string>("This Project")
const name=ref<string>("")

const users=ref<userAPI.UserInfo[]>([])
const page=ref<number>(1)//start from 1
const pagesize=ref<number>(0)
const totalsize=ref<number>(0)

const cur_user=ref<userAPI.UserInfo>(null)
const invite_kick_user=ref<userAPI.UserInfo>(null)

const update_user=ref<userAPI.UserInfo>(null)
const update_user_new_name=ref<string>("")
const update_user_new_department=ref<string[]>([])

const ing=ref<boolean>(false)
const optype=ref<string>("")
function user_op(){
	if(!state.set_load()){
		return
	}
	switch(optype.value){
		case "search":{
			if(!name.value){
				state.clear_load()
				state.set_error("error",-2,"missing user name")
				return
			}
			let req = {
				project_id:state.project.cur.project_id,
				user_name:name.value,
				only_project:range.value=="This Project",
				page:page.value,
			}
			client.userClient.search_users({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_error("error",e.code,e.msg)
			},(resp :userAPI.SearchUsersResp)=>{
				state.clear_load()
				users.value=resp.users
				users.value[0].roles=["a","b","c"]
				page.value=resp.page
				pagesize.value=resp.pagesize
				totalsize.value=resp.totalsize
				cur_user.value=null
			})
			break
		}
		case "invite":{
			let req = {
				project_id:state.project.cur.project_id,
				user_id:invite_kick_user.value.user_id,
			}
			client.userClient.invite_project({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_error("error",e.code,e.msg)
			},(resp: userAPI.InviteProjectResp)=>{
				state.clear_load()
				invite_kick_user.value.invited=true
			})
			break
		}
		case "kick":{
			let req = {
				project_id:state.project.cur.project_id,
				user_id:invite_kick_user.value.user_id,
			}
			client.userClient.kick_project({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_error("error",e.code,e.msg)
			},(resp: userAPI.KickProjectResp)=>{
				state.clear_load()
				if(range.value=="This Project"){
					for(let i=0;i<users.value.length;i++){
						if(users.value[i].user_id==invite_kick_user.value.user_id){
							users.value.splice(i,1)
							break
						}
					}
					if(cur_user.value&&cur_user.value.user_id==invite_kick_user.value.user_id){
						cur_user.value=null
					}
				}else{
					invite_kick_user.value.invited=false
				}
			})
			break
		}
		case "update":{
			let req = {
				user_id: update_user.value.user_id,
				new_user_name: update_user_new_name.value,
				new_department: update_user_new_department.value,
			}
			client.userClient.update_user({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_error("error",e.code,e.msg)
			},(resp: userAPI.UpdateUserResp)=>{
				state.clear_load()
				update_user.user_name = new_user_name.value
				update_user.department = update_user_new_department.value
			})
			break
		}
		default:{
			state.clear_load()
			state.set_error("error",-2,"unknown operation")
		}
	}
}

function parsetime(timestamp :number):string{
	let t=new Date(timestamp*1000)
	let offset=Math.abs(t.getTimezoneOffset())
	let hour=Math.floor(offset/60)
	let min=offset%60
	let result = t.toLocaleString()
	if(t.getTimezoneOffset()<0){
		result+=" UTC+"
	}else{
		result+=" UTC-"
	}
	if(hour<10){
		result+="0"+hour
	}else{
		result+=hour
	}
	result+=":"
	if(min<10){
		result+="0"+min
	}else{
		result+=min
	}
	return result
}
</script>
<template>
	<va-modal v-model="ing" attach-element="#app" max-width="1000px" max-height="600px" hide-default-actions no-dismiss overlay-opacity="0.2" z-index="999">
		<template #default>
			<div v-if="optype=='invite'">
				<va-card color="primary" gradient style="margin:0 0 5px 0">
					<va-card-title>Warning</va-card-title>
					<va-card-content>
						<p>You are inviting user: {{ invite_kick_user.user_name }} join project: {{ state.project.cur.project_name }}.</p>
						<p>Please confirm!</p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="ing=false;user_op()" gradient>Invite</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='kick'">
				<va-card color="primary" gradient style="margin:0 0 5px 0">
					<va-card-title>Warning</va-card-title>
					<va-card-content>
						<p>You are kicking user: {{ invite_kick_user.user_name }} out of project: {{ state.project.cur.project_name }}.</p>
						<p>Please confirm!</p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="ing=false;user_op()" gradient>Kick</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='update'">
				<va-input v-model="update_user_new_name" label="New User Name" style="width:100%;margin:1px 0"></va-input>
				<div style="width:300px;display:flex;margin:1px 0">
				</div>
				<div style="display:flex;justify-content:center;flex:margin:1px 0">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="ing=false;user_op()" gradient>Kick</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
		</template>
	</va-modal>
	<div style="display:flex;flex:1;flex-direction:column;margin:1px;width:100%;overflow-y:auto">
		<div style="display:flex;justify-content:center;margin:1px">
			<va-select label="Search Range" dropdown-icon="" outline style="width:130px;margin-right:1px" :options="ranges" v-model="range" trigger="hover">
				<template #option="{ option ,index , selectOption }">
					<va-hover
					stateful
					@click="()=>{
						if(option!=range){
							page=1
							pagesize=0
							totalsize=0
							users=[]
							cur_user=null
							selectOption(option)
						}
					}">
						<template #default="{hover}">
							<div
							style="padding:10px;cursor:pointer"
							:style="{'background-color':hover?'var(--va-background-border)':'',color:hover||range==option?'var(--va-primary)':'black'}"
							>
								{{option}}
							</div>
						</template>
					</va-hover>
				</template>
			</va-select>
			<va-input label="User Name" outline style="max-width:250px;margin:0 1px" v-model="name" @keyup.enter="optype='search';user_op()"></va-input>
			<va-button style="margin-left:1px" @click="optype='search';user_op()" :disabled="!Boolean(name)">Search</va-button>
		</div>
		<div style="flex:1;display:flex;flex-direction:column;margin:1px;overflow-y:auto">
			<div v-for="user of users" style="display:flex;flex-direction:column;overflow-y:auto" :style="{flex:user.open?1:undefined}">
				<div
					v-if="!Boolean(cur_user)||cur_user.user_id==user.user_id"
					style="display:flex;margin:1px 0;cursor:pointer;align-items:center"
					:style="{'background-color':user.hover?'var(--va-shadow)':'var(--va-background-element)'}"
					@click="()=>{
						user.open=!user.open
						if(user.open){
							cur_user=user
						}else{
							cur_user=null
						}
					}"
					@mouseover="user.hover=true"
					@mouseout="user.hover=false"
				>
					<span style="flex:1;padding:12px;color:var(--va-primary)">{{user.user_name}}</span>
					<va-button
					v-if="state.page.node.canwrite||state.page.node.admin"
					size="small"
					style="width:50px;height:30px;margin:0px 4px"
					@mouseover.stop=""
					@mouseout.stop=""
					@click.stop="()=>{
						optype='update'
						update_user=user
						update_user_new_name=user.user_name
						update_user_new_department=user.department
						ing=true
					}">
						Update
					</va-button>
					<va-button
					v-if="state.page.node.admin"
					size="small"
					style="width:50px;height:30px"
					@mouseover.stop=""
					@mouseout.stop=""
					@click.stop="()=>{
						if(user.invited){
							optype='kick'
						}else{
							optype='invite'
						}
						invite_kick_user=user
						ing=true
					}">
						{{user.invited?'Kick':'Invite'}}
					</va-button>
					<span style="width:60px;padding:12px 20px;color:var(--va-primary)">{{ user.open?'▲':'▼' }}</span>
				</div>
				<div v-if="(!Boolean(cur_user)||cur_user.user_id==user.user_id)&&user.open" style="margin:1px 20px;display:flex;justify-content:space-around;background-color:var(--va-background-element);color:var(--va-primary)">
					<div style="width:400px;margin:1px;padding:12px 10px;display:flex">
						<span><b>JoinTime:</b></span>
						<span style="flex:1;text-align:center">{{ parsetime(user.ctime) }}</span>
					</div>
					<div style="width:400px;margin:1px;padding:12px 10px;display:flex">
						<span><b>Department:</b></span>
						<span style="flex:1;text-align:center">{{ Boolean(user.department)?user.department.join('/'):'' }}</span>
					</div>
				</div>
				<div v-if="(!Boolean(cur_user)||cur_user.user_id==user.user_id)&&user.open" style="margin:0 20px;display:flex;flex:1;overflow-y:auto;color:var(--va-primary)">
					<div style="width:150px;margin:0 1px 1px 0;display:flex;flex-direction:column;overflow-y:auto">
						<div style="padding:12px 10px;margin:1px 0;background-color:var(--va-background-element)"><b>Role Name</b></div>
						<div v-if="Boolean(user.roles)&&user.roles.length>0" style="flex:1;display:flex;flex-direction:column;overflow-y:auto">
							<va-hover stateful v-for="role of user.roles" @click="">
								<template #default="{hover}">
									<div style="padding:12px 10px;margin:1px 0;color:var(--va-primary);cursor:pointer" :style="{'background-color':hover?'var(--va-shadow)':'var(--va-background-element)'}">{{role}}</div>
								</template>
							</va-hover>
						</div>
						<div v-else style="flex:1;display:flex;flex-direction:column;overflow-y:auto">
							<div style="padding:12px 10px;margin:1px 0;background-color:var(--va-background-element);color:var(--va-shadow)">No Roles</div>
						</div>
					</div>
					<div style="flex:1;margin:0 0 1px 1px;display:flex;flex-direction:column;overflow-y:auto">
						<div style="padding:12px 10px;margin:1px 0;background-color:var(--va-background-element)"><b>Role Permission Nodes</b></div>
						<div style="flex:1;margin:1px 0;background-color:var(--va-background-element);overflow-y:auto">
						</div>
					</div>
				</div>
			</div>
		</div>
		<va-pagination v-if="!Boolean(cur_user)" v-model="page" :total="totalsize" :page-size="pagesize" :visible-pages="7" gapped boundary-numbers :direction-links="false" style="margin:1px;align-self:center">
		</va-pagination>
	</div>
</template>
