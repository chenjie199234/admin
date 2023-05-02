<script setup lang="ts">
import {ref,onMounted} from 'vue'
import * as userAPI from '../../api/user_browser_toc'
import * as permissionAPI from '../../api/permission_browser_toc'
import * as state from './state'
import * as client from './client'

import nodetree from './nodetree.vue'

onMounted(()=>{
	if(state.page.node.admin){
		ranges.value=["This Project","All Projects"]
	}else{
		ranges.value=["This Project"]
	}
})
const targets=ref<string[]>(["User","Role"])
const target=ref<string>("User")
const ranges=ref<string[]>([])
const range=ref<string>("This Project")
const search=ref<string>("")

const users=ref<userAPI.UserInfo[]>([])
const roles=ref<userAPI.RoleInfo[]>([])
const page=ref<number>(1)//start from 1
const pagesize=ref<number>(0)
const totalsize=ref<number>(0)

function invited(user: userAPI.UserInfo):boolean{
	if(!user.project_roles){
		return false
	}
	for(let i=0;i<user.project_roles.length;i++){
		if(user.project_roles[i].project_id[1]==state.project.cur_id[1]){
			return true
		}
	}
	return false
}

const cur_user=ref<userAPI.UserInfo>(null)
const invite_kick_user=ref<userAPI.UserInfo>(null)

const update_user=ref<userAPI.UserInfo>(null)
const update_user_new_name=ref<string>("")
const update_user_new_department=ref<string[]>([])

const cur_role=ref<userAPI.RoleInfo>(null)

const create_role_name=ref<string>("")
const create_role_comment=ref<string>("")

const update_role_comment=ref<string>("")

const del_role=ref<userAPI.RoleInfo>(null)

const del_user_role_userid=ref<string>("")
const del_user_role_username=ref<string>("")
const del_user_role_rolename=ref<string>("")

const add_user_role_userid=ref<string>("")
const add_user_role_username=ref<string>("")
const add_user_role_rolename=ref<string>("")

const hover_role=ref<string|null>(null)
const node_role=ref<string|null>(null)
const node=ref<permissionAPI.NodeInfo>(null)
const update_node=ref<permissionAPI.NodeInfo>(null)
const canread=ref<boolean>(false)
const canwrite=ref<boolean>(false)
const admin=ref<boolean>(false)

const ing=ref<boolean>(false)
const optype=ref<string>("")
function op(){
	if(!state.set_load()){
		return
	}
	switch(optype.value){
		case "search":{
			if(target.value=="User"){
				let req = {
					project_id:state.project.cur_id,
					user_name:search.value,
					only_project:range.value=="This Project",
					page:page.value,
				}
				client.userClient.search_users({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
					state.clear_load()
					state.set_alert("error",e.code,e.msg)
				},(resp :userAPI.SearchUsersResp)=>{
					roles.value=[]
					users.value=resp.users
					page.value=resp.page
					pagesize.value=resp.pagesize
					totalsize.value=resp.totalsize
					cur_user.value=null
					cur_role.value=null
					state.clear_load()
				})
			}else{
				let req = {
					project_id:state.project.cur_id,
					role_name:search.value,
					page:page.value,
				}
				client.userClient.search_roles({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
					state.clear_load()
					state.set_alert("error",e.code,e.msg)
				},(resp :userAPI.SearchUsersResp)=>{
					users.value=[]
					roles.value=resp.roles
					page.value=resp.page
					pagesize.value=resp.pagesize
					totalsize.value=resp.totalsize
					cur_role.value=null
					cur_user.value=null
					state.clear_load()
				})
			}
			break
		}
		case "invite":{
			let req = {
				project_id:state.project.cur_id,
				user_id:invite_kick_user.value.user_id,
			}
			client.userClient.invite_project({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: userAPI.InviteProjectResp)=>{
				invite_kick_user.value.project_roles=[{project_id:state.project.cur_id,roles:[]}]
				ing.value=false
				state.clear_load()
			})
			break
		}
		case "kick":{
			let req = {
				project_id:state.project.cur_id,
				user_id:invite_kick_user.value.user_id,
			}
			client.userClient.kick_project({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: userAPI.KickProjectResp)=>{
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
					for(let i=0;i<invite_kick_user.value.project_roles.length;i++){
						if(invite_kick_user.value.project_roles[i].project_id[1]==state.project.cur_id[1]){
							invite_kick_user.value.project_roles.splice(i,1)
							break
						}
					}
				}
				ing.value=false
				state.clear_load()
			})
			break
		}
		case "update_user":{
			let req = {
				user_id: update_user.value.user_id,
				new_user_name: update_user_new_name.value,
				new_department: update_user_new_department.value,
			}
			client.userClient.update_user({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: userAPI.UpdateUserResp)=>{
				update_user.value.user_name = update_user_new_name.value
				update_user.value.department = update_user_new_department.value
				ing.value=false
				state.clear_load()
			})
			break
		}
		case "create_role":{
			let req = {
				project_id:state.project.cur_id,
				role_name:create_role_name.value,
				comment:create_role_comment.value,
			}
			client.userClient.create_role({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: userAPI.CreateRoleResp)=>{
				create_role_name.value=""
				create_role_comment.value=""
				ing.value=false
				state.clear_load()
			})
			break
		}
		case "update_role":{
			let req = {
				project_id:state.project.cur_id,
				role_name:cur_role.value.role_name,
				new_comment:update_role_comment.value,
			}
			client.userClient.update_role({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: userAPI.UpdateRoleResp)=>{
				cur_role.value.comment=update_role_comment.value
				state.clear_load()
				state.set_alert("info","0","update role comment success")
			})
			break
		}
		case "del_role":{
			let req = {
				project_id:state.project.cur_id,
				role_names:[del_role.value.role_name],
			}
			client.userClient.del_roles({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: userAPI.DelRolesResp)=>{
				for(let i=0;i<roles.value.length;i++){
					if(roles.value[i]==del_role.value){
						roles.value.splice(i,1)
						break
					}
				}
				if(del_role.value==cur_role.value){
					cur_role.value=null
				}
				del_role.value=null
				ing.value=false
				state.clear_load()
			})
			break
		}
		case "del_user_role":{
			let req = {
				project_id:state.project.cur_id,
				user_id:del_user_role_userid.value,
				role_name:del_user_role_rolename.value,
			}
			client.userClient.del_user_role({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: userAPI.DelUserRoleResp)=>{
				let index:number=cur_user.value.project_roles.findIndex((project_role)=>{
					return project_role.project_id[1]==state.project.cur_id[1]
				})
				if(index!=-1){
					index = cur_user.value.project_roles[index].roles.findIndex((role)=>{
						return role==del_user_role_rolename.value
					})
					if(index!=-1){
						cur_user.value.project_roles[index].roles.splice(index,1)
					}
				}
				ing.value=false
				state.clear_load()
			})
			break
		}
		case "add_user_role_missinguser":
		case "add_user_role_missingrole":{
			let req = {
				project_id:state.project.cur_id,
				user_id:add_user_role_userid.value,
				role_name:add_user_role_rolename.value,
			}
			client.userClient.add_user_role({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: userAPI.DelUserRoleResp)=>{
				if(cur_user.value){
					if(cur_user.value.project_roles){
						let index:number=cur_user.value.project_roles.findIndex((project_role)=>{
							return project_role.project_id[1]==state.project.cur_id[1]
						})
						if(index==-1){
							cur_user.value.project_roles.push({project_id:state.project.cur_id,roles:[add_user_role_rolename.value]})
						}else if(!cur_user.value.project_roles[index].roles){
							cur_user.value.project_roles[index].roles = [add_user_role_rolename.value]
						}else if(!cur_user.value.project_roles[index].roles.includes(add_user_role_rolename.value)){
							cur_user.value.project_roles[index].roles.push(add_user_role_rolename.value)
						}
					}else{
						cur_user.value.project_roles = [{project_id:state.project.cur_id,roles:[add_user_role_rolename.value]}]
					}
				}
				ing.value=false
				state.clear_load()
			})
			break
		}
		case "update_permission":{
			if(node_role.value==''){
				let req = {
					user_id:cur_user.value.user_id,
					node_id:update_node.value.node_id,
					admin:admin.value,
					canread:canread.value,
					canwrite:canwrite.value,
				}
				client.permissionClient.update_user_permission({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
					state.clear_load()
					state.set_alert("error",e.code,e.msg)
				},(resp :permissionAPI.UpdateUserPermissionResp)=>{
					update_node.value.canread=canread.value
					update_node.value.canwrite=canwrite.value
					update_node.value.admin=admin.value
					update_node.value=null
					ing.value=false
					state.clear_load()
				})
			}else if(node_role.value){
				let req = {
					project_id:state.project.cur_id,
					role_name:node_role.value,
					node_id:update_node.value.node_id,
					admin:admin.value,
					canread:canread.value,
					canwrite:canwrite.value,
				}
				client.permissionClient.update_role_permission({"Token":state.user.token},req,client.timeout,(e :permissionAPI.Error)=>{
					state.clear_load()
					state.set_alert("error",e.code,e.msg)
				},(resp :permissionAPI.UpdateRolePermissionResp)=>{
					update_node.value.canread=canread.value
					update_node.value.canwrite=canwrite.value
					update_node.value.admin=admin.value
					update_node.value=null
					ing.value=false
					state.clear_load()
				})
			}
			break
		}
		default:{
			state.clear_load()
			state.set_alert("error",-2,"unknown operation")
		}
	}
}

let timeid :number = 0
function assign_search(part: string){
	clearTimeout(timeid)
	if(!part){
		if(target.value=="User"){
			roles.value=[]
		}else{
			users.value=[]
		}
		return
	}
	if(target.value=="User"){
		timeid = setTimeout(() => {
			if(!state.set_load()){
				return
			}
			let req = {
				project_id:state.project.cur_id,
				role_name:add_user_role_rolename.value,
				page:0,
			}
			client.userClient.search_roles({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp :userAPI.SearchRolesResp)=>{
				if(resp.roles){
					roles.value=resp.roles
				}else{
					roles.value=[]
				}
				state.clear_load()
			})
		}, 500)
	}else{
		timeid = setTimeout(() => {
			if(!state.set_load()){
				return
			}
			let req = {
				project_id:state.project.cur_id,
				user_name:add_user_role_username.value,
				only_project:true,
				page:0,
			}
			client.userClient.search_users({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp :userAPI.SearchUsersResp)=>{
				if(resp.users){
					users.value=resp.users
				}else{
					users.value=[]
				}
				state.clear_load()
			})
		}, 500)
	}
}
function clear_assign_search(){
	clearTimeout(timeid)
}
function get_user_permission(){
	if(!state.set_load()){
		return
	}
	let req = {
		project_id:state.project.cur_id,
		user_id:cur_user.value.user_id,
		need_user_role_node:false,
	}
	client.permissionClient.list_user_node({"Token":state.user.token},req,client.timeout,(e :permissionAPI.Error)=>{
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
	},(resp: permissionAPI.ListUserNodeResp)=>{
		node_role.value=''
		node.value=resp.node
		state.clear_load()
	})
}
function get_role_permission(role:string){
	if(!state.set_load()){
		return
	}
	let req = {
		project_id:state.project.cur_id,
		role_name:role,
	}
	client.permissionClient.list_role_node({"Token":state.user.token},req,client.timeout,(e :permissionAPI.Error)=>{
		state.clear_load()
		state.set_alert("error",e.code,e.msg)
	},(resp: permissionAPI.ListRoleNodeResp)=>{
		node_role.value=role
		node.value=resp.node
		state.clear_load()
	})
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
						<p>You are inviting user: {{ invite_kick_user.user_name }} join project: {{ state.project.cur_name }}.</p>
						<p>Please confirm!</p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="op" gradient>Invite</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='kick'">
				<va-card color="primary" gradient style="margin:0 0 5px 0">
					<va-card-title>Warning</va-card-title>
					<va-card-content>
						<p>You are kicking user: {{ invite_kick_user.user_name }} out of project: {{ state.project.cur_name }}.</p>
						<p>Please confirm!</p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="op" gradient>Kick</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='update_user'">
				<va-input v-model.trim="update_user_new_name" label="New User Name" style="width:300px;margin:1px 0" @keyup.enter="()=>{if(update_user_new_name!=update_user.user_name){op()}}"></va-input>
				<div style="width:300px;display:flex;margin:1px 0">
				<!-- TODO: department -->
				</div>
				<div style="display:flex;justify-content:center;margin:1px 0">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="op" gradient :disabled="update_user_new_name==update_user.user_name">Update</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='create_role'" style="display:flex;flex-direction:column">
				<va-input v-model.trim="create_role_name" label="New Role Name" style="margin:1px 0;width:300px" @keyup.enter="()=>{if(create_role_name!=''){op()}}" />
				<va-input v-model.trim="create_role_comment" label="New Role Comment" style="margin:1px 0;width:300px" @keyup.enter="()=>{if(create_role_name!=''){op()}}" />
				<div style="display:flex;justify-content:center;margin:1px 0">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="op" gradient :disabled="create_role_name==''">Create</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="create_role_name='';create_role_comment='';ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='del_role'">
				<va-card color="primary" gradient style="margin:0 0 5px 0">
					<va-card-title>Warning</va-card-title>
					<va-card-content>
						<p>You are deleting role: {{ del_role.role_name }}.</p>
						<p>Please confirm!</p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center;margin:1px 0">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="op" gradient>Del</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="del_role=null;ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='add_user_role_missingrole'" style="display:flex;flex-direction:column;align-items:center">
				<va-card color="primary" gradient style="margin:0 0 5px 0;min-width:400px">
					<va-card-content>
						<p>Your are assigning user: {{add_user_role_username}} a role.</p>
					</va-card-content>
				</va-card>
				<div style="width:100%;display:flex">
					<div style="margin-right:2px">
						<va-input label="Role Name" style="width:170px" v-model.trim="add_user_role_rolename" @update:model-value="assign_search" />
						<div style="display:flex;justify-content:center;margin:1px 0">
							<va-button style="width:80px;margin:5px 5px 0 0" @click="op" gradient :disabled="roles.findIndex((role)=>{return role.role_name==add_user_role_rolename})==-1">Assign</va-button>
							<va-button style="width:80px;margin:5px 0 0 5px" @click="clear_assign_search();ing=false" gradient>Cancel</va-button>
						</div>
					</div>
					<div style="margin-left:2px;flex:1;background-color:var(--va-background-element);border-radius:3px;height:77px;overflow-y:auto">
						<div
							v-for="role of roles"
							style="padding:5px;cursor:pointer"
							:style="{'background-color':role.hover?'var(--va-shadow)':undefined}"
							@mouseover="role.hover=true"
							@mouseout="role.hover=false"
							@click="add_user_role_rolename=role.role_name"
						>
							{{role.role_name}}
						</div>
					</div>
				</div>
			</div>
			<div v-else-if="optype=='add_user_role_missinguser'" style="display:flex;flex-direction:column;align-items:center">
				<va-card color="primary" gradient style="margin:0 0 5px 0;min-width:400px">
					<va-card-content>
						<p>Your are assigning role: {{add_user_role_rolename}} to a user.</p>
					</va-card-content>
				</va-card>
				<div style="width:100%;display:flex">
					<div style="margin-right:2px">
						<va-input label="User Name" style="width:170px" v-model.trim="add_user_role_username" @update:model-value="add_user_role_userid='';assign_search($event)" />
						<div style="display:flex;justify-content:center;margin:1px 0">
							<va-button style="width:80px;margin:5px 5px 0 0" @click="op" gradient :disabled="add_user_role_userid==''">Assign</va-button>
							<va-button style="width:80px;margin:5px 0 0 5px" @click="clear_assign_search();ing=false" gradient>Cancel</va-button>
						</div>
					</div>
					<div style="margin-left:2px;flex:1;background-color:var(--va-background-element);border-radius:3px;height:77px;overflow-y:auto">
						<div
							v-for="user of users"
							style="padding:5px;cursor:pointer"
							:style="{'background-color':user.hover?'var(--va-shadow)':undefined}"
							@mouseover="user.hover=true"
							@mouseout="user.hover=false"
							@click="add_user_role_userid=user.user_id;add_user_role_username=user.user_name"
						>
							{{user.user_name}}
						</div>
					</div>
				</div>
			</div>
			<div v-else-if="optype=='del_user_role'">
				<va-card color="primary" gradient style="margin:0 0 5px 0">
					<va-card-title>Warning</va-card-title>
					<va-card-content>
						<p>You are discharging user: {{ del_user_role_username }}'s role: {{ del_user_role_rolename }}.</p>
						<p>Please confirm!</p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center;margin:1px 0">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="op" gradient>Del</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='update_permission'" style="display:flex;flex-direction:column">
				<va-card color="primary" gradient style="margin:0 0 5px 0">
					<va-card-title>Warning</va-card-title>
					<va-card-content>
						<p v-if="node_role==''">You are updating user: {{ cur_user.user_name }}'s permission on node: {{ update_node.node_name }}.</p>
						<p v-else>You are updating role: {{ node_role }}'s permission on node: {{ update_node.node_name }}.</p>
						<p>Please confirm!</p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:5px 10px 0 0" @click="op" gradient>Update</va-button>
					<va-button style="width:80px;margin:5px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
		</template>
	</va-modal>
	<div style="flex:1;display:flex;flex-direction:column;margin:1px;overflow-y:auto">
		<div style="display:flex;margin:1px">
			<div style="flex:1"></div>
			<va-select label="Target" dropdown-icon="" outline style="width:100px;margin-right:1px" :options="targets" v-model="target" trigger="hover">
				<template #option='{option,index,selectOption}'>
					<va-hover stateful @click="()=>{
							if(option!=target){
								page=1
								pagesize=0
								totalsize=0
								users=[]
								roles=[]
								cur_user=null
								cur_role=null
								search=''
								selectOption(option)
							}
						}"
					>
						<template #default="{hover}">
							<div
								style="padding:10px;cursor:pointer"
								:style="{'background-color':hover?'var(--va-background-border)':'',color:hover||target==option?'var(--va-primary)':'black'}"
							>
								{{option}}
							</div>
						</template>
					</va-hover>
				</template>
			</va-select>
			<va-select v-if="target=='User'" label="Search Range" dropdown-icon="" outline style="width:130px;margin:0 1px" :options="ranges" v-model="range" trigger="hover">
				<template #option="{ option ,index , selectOption }">
					<va-hover stateful @click="()=>{
							if(option!=range){
								page=1
								pagesize=0
								totalsize=0
								users=[]
								cur_user=null
								selectOption(option)
							}
						}"
					>
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
			<va-input :label="target=='User'?'User Name':'Role Name'" outline style="max-width:250px;margin:0 1px" v-model.trim="search" @keyup.enter="optype='search';user_op()"></va-input>
			<va-button style="margin:0 1px" @click="optype='search';op()">Search</va-button>
			<div style="flex:1"></div>
			<va-button v-if="target=='Role'" style="margin-left:1px" @click="optype='create_role';ing=true">Create</va-button>
		</div>
		<div v-if="target=='User'" style="flex:1;display:flex;flex-direction:column;margin:1px;overflow-y:auto">
			<template v-for="user of users">
				<div
					v-if="!Boolean(cur_user)||cur_user==user"
					style="display:flex;margin:1px 0;align-items:center;cursor:pointer"
					:style="{'background-color':user.hover?'var(--va-shadow)':'var(--va-background-element)'}"
					@click="()=>{
						if(!cur_user){
							cur_user=user
							hover_role=null
							node_role=null
						}else{
							cur_user=null
						}
					}"
					@mouseover="user.hover=true"
					@mouseout="user.hover=false"
				>
					<span style="width:40px;padding:12px 20px;color:var(--va-primary)">{{cur_user==user?'▲':'▼' }}</span>
					<span style="flex:1;padding:12px 20px;color:var(--va-primary)">{{user.user_name}}</span>
					<span style="padding:12px 0;color:var(--va-primary)">Create Time: {{parsetime(user.ctime)}}</span>
					<va-button
						v-if="state.page.node.canwrite||state.page.node.admin"
						size="small"
						style="width:50px;height:30px;margin:0 10px"
						@mouseover.stop=""
						@mouseout.stop=""
						@click.stop="()=>{
							optype='update_user'
							update_user=user
							update_user_new_name=user.user_name
							update_user_new_department=user.department
							ing=true
						}"
					>
						Update
					</va-button>
					<va-button
						v-if="state.page.node.admin"
						size="small"
						style="width:50px;height:30px;margin-right:10px"
						@mouseover.stop=""
						@mouseout.stop=""
						@click.stop="()=>{
							if(invited(user)){
								optype='kick'
							}else{
								optype='invite'
							}
							invite_kick_user=user
							ing=true
						}"
					>
						{{invited(user)?'Kick':'Invite'}}
					</va-button>
				</div>
				<div v-if="cur_user==user" style="margin:1px 10px;display:flex;justify-content:space-around;background-color:var(--va-background-element);color:var(--va-primary)">
				<!-- basic info -->
					<div style="width:400px;margin:1px;padding:12px 10px;display:flex">
						<span><b>UserID:</b></span>
						<span style="flex:1;text-align:center">{{ user.user_id }}</span>
					</div>
					<div style="width:400px;margin:1px;padding:12px 10px;display:flex">
						<span><b>Department:</b></span>
						<span style="flex:1;text-align:center">{{ Boolean(user.department)?user.department.join('/'):'' }}</span>
					</div>
				</div>
				<div v-if="cur_user==user&&invited(user)" style="margin:0 10px;display:flex;flex:1;overflow-y:auto;color:var(--va-primary)">
				<!-- permission info -->
					<div style="width:200px;margin:0 1px 1px 0;display:flex;flex-direction:column;overflow-y:auto">
						<div
							style="padding:12px 10px;margin:1px 0;cursor:pointer"
							:style="{'background-color':hover_role==''?'var(--va-shadow)':node_role==''?'white':'var(--va-background-element)'}"
							@mouseover.stop="hover_role=''"
							@mouseout.stop="hover_role=null"
							@click="node_role=null;get_user_permission()"
						>
							User Permissions
						</div>
						<div style="display:flex;margin:1px 0;align-items:center;background-color:var(--va-background-element)">
							<div style="flex:1;padding:12px 10px;white-space:nowrap"><b>Role Permissions</b></div>
							<va-button
								v-if="state.page.node.admin"
								size="small"
								style="margin:0 2px"
								@mouseover.stop=""
								@mouseout.stop=""
								@click.stop="
									add_user_role_username=user.user_name;
									add_user_role_userid=user.user_id;
									add_user_role_rolename='';
									roles=[];
									optype='add_user_role_missingrole';
									ing=true"
							>
								+
							</va-button>
						</div>
						<div
							v-if="!Boolean(user.project_roles[0].roles)||user.project_roles[0].roles.length==0"
							style="padding:12px 10px;margin:1px 0;color:var(--va-shadow);background-color:var(--va-background-element)"
						>
							No Roles
						</div>
						<div
							v-for="role of user.project_roles[0].roles"
							style="display:flex;margin:1px 0;align-items:center;cursor:pointer"
							:style="{'background-color':hover_role==role?'var(--va-shadow)':node_role==role?'white':'var(--va-background-element)'}"
							@mouseover.stop="hover_role=role"
							@mouseout.stop="hover_role=null"
							@click="node_role=null;get_role_permission(role)"
						>
							<div style="flex:1;padding:12px 10px;white-space:nowrap">{{role}} Permissions</div>
							<va-button
								v-if="state.page.node.admin"
								size="small"
								style="margin:0 2px"
								@mouseover.stop=""
								@mouseout.stop=""
								@click.stop="
									del_user_role_username=user.user_name;
									del_user_role_userid=user.user_id;
									del_user_role_rolename=role;
									optype='del_user_role';
									ing=true"
							>
								X
							</va-button>
						</div>
					</div>
					<div style="flex:1;margin:1px 0 1px 1px;display:flex;background-color:var(--va-background-element);overflow:auto">
						<nodetree v-if="node_role!=null&&node!=null" :pnode="node" :deep="0" @permissionevent="(updatenode,r,w,a)=>{
								update_node=updatenode;
								canread=r;
								canwrite=w;
								admin=a;
								optype='update_permission';
								ing=true
							}"
						/>
						<div v-else-if="node_role!=null" style="flex:1;display:flex;justify-content:center;align-items:center"><b>No Permissions</b></div>
					</div>
				</div>
			</template>
		</div>
		<div v-if="target=='Role'" style="flex:1;display:flex;flex-direction:column;margin:1px;overflow-y:auto">
			<template v-for="role of roles">
				<div
					v-if="!Boolean(cur_role)||cur_role==role"
					style="display:flex;margin:1px 0;align-items:center;cursor:pointer"
					:style="{'background-color':hover_role==role.role_name?'var(--va-shadow)':'var(--va-background-element)'}"
					@click="()=>{
						if(!cur_role){
							cur_role=role
							hover_role=null
							node_role=null
							update_role_comment=role.comment
							get_role_permission(role.role_name)
						}else{
							cur_role=null
						}
					}"
					@mouseover="hover_role=role.role_name"
					@mouseout="hover_role=null"
				>
					<span style="width:40px;padding:12px 20px;color:var(--va-primary)">{{ cur_role==role?'▲':'▼' }}</span>
					<span style="flex:1;padding:12px 20px;color:var(--va-primary)">{{role.role_name}}</span>
					<span style="padding:12px 0;color:var(--va-primary)">Create Time: {{parsetime(role.ctime)}}</span>
					<va-button
						v-if="state.page.node.admin"
						size="small"
						style="width:50px;height:30px;margin:0 10px"
						@mouseover.stop=""
						@mouseout.stop=""
						@click.stop="
							add_user_role_rolename=role.role_name;
							add_user_role_userid='';
							add_user_role_username='';
							users=[];
							optype='add_user_role_missinguser';
							ing=true"
					>
						Assign
					</va-button>
					<va-button
						v-if="state.page.node.admin"
						size="small"
						style="width:50px;height:30px;margin-right:10px"
						@mouseover.stop=""
						@mouseout.stop=""
						@click.stop="del_role=role;optype='del_role';ing=true"
					>
						Del
					</va-button>
				</div>
				<div v-if="cur_role==role" style="display:flex;margin:1px 10px;align-items:center">
					<va-input type="textarea" autosize :min-rows="2" label="Role Comment" outline v-model.trim="update_role_comment" style="flex:1" :readonly="!state.page.node.canwrite&&!state.page.node.admin" />
					<va-button v-if="state.page.node.canwrite||state.page.node.admin" :disabled="update_role_comment==role.comment" style="margin-left:2px" @click="optype='update_role';op()">Update</va-button>
				</div>
				<div v-if="cur_role==role" style="flex:1;margin:1px 10px;display:flex;background-color:var(--va-background-element);color:var(--va-primary);overflow:auto">
					<nodetree v-if="node_role!=null&&node!=null" :pnode="node" :deep="0" @permissionevent="(updatenode,r,w,a)=>{
							update_node=updatenode;
							canread=r;
							canwrite=w;
							admin=a;
							optype='update_permission';
							ing=true
						}"
					/>
					<div v-else-if="node_role!=null" style="flex:1;display:flex;justify-content:center;align-items:center"><b>No Permissions</b></div>
				</div>
			</template>
		</div>
		<va-pagination
			v-if="!Boolean(cur_user)&&!Boolean(cur_role)"
			v-model="page"
			:pages="totalsize==0||pagesize==0?1:Math.ceil(totalsize/pagesize)"
			:visible-pages="5"
			boundary-numbers
			direction-icon-left="〈"
			direction-icon-right="〉"
			gapped
			style="margin:1px;align-self:center"
			@update:model-value="optype='search';op()"
		/>
	</div>
</template>
