<script setup lang="ts">
import {ref,computed} from 'vue'
import * as userAPI from './api/admin_user_browser_toc'
import * as permissionAPI from './api/admin_permission_browser_toc'
import * as state from './state'
import * as client from './client'

import nodetree from './nodetree.vue'

const targets=ref<string[]>(["User","Role"])
const ranges=computed(()=>{
	if(state.page.node!.admin){
		return ["This Project","All Projects"]
	}else{
		return ["This Project"]
	}
})
const target=ref<string>("User")
const range=ref<string>("This Project")
const search=ref<string>("")

const users=ref<userAPI.UserInfo[]>([])
const userhover=ref<userAPI.UserInfo|null>(null)
function user_bindstyle(user: userAPI.UserInfo){
	let style={}
	if(user==userhover.value&&invited(user)){
		style["background-color"]="var(--va-shadow)"
	}else{
		style["background-color"]="var(--va-background-element)"
	}
	if(invited(user)){
		style["cursor"]="pointer"
	}
	return style
}
const roles=ref<userAPI.RoleInfo[]>([])
const rolehover=ref<userAPI.RoleInfo|null>(null)
const page=ref<number>(1)//start from 1
const pagesize=ref<number>(0)
const totalsize=ref<number>(0)
function user_has_role(user: userAPI.UserInfo,role: userAPI.RoleInfo):boolean{
	if(!user.project_roles![0]!.roles){
		return false
	}
	return user.project_roles![0]!.roles!.includes(role.role_name)
}

function invited(user: userAPI.UserInfo):boolean{
	if(!user){
		return false
	}
	if(!user!.project_roles){
		return false
	}
	for(let i=0;i<user!.project_roles!.length;i++){
		if(!user!.project_roles![i]){
			return false
		}
		if(!user!.project_roles[i]!.project_id){
			return false
		}
		if(user!.project_roles[i]!.project_id![1]==state.project.info!.project_id![1]){
			return true
		}
	}
	return false
}

//user
const cur_user=ref<userAPI.UserInfo|null>(null)
const invite_kick_user=ref<userAPI.UserInfo|null>(null)
const update_user_delete_role_rolename=ref<string>("")

//role
const create_role_name=ref<string>("")
const create_role_comment=ref<string>("")

const cur_role=ref<userAPI.RoleInfo|null>(null)
const del_role=ref<userAPI.RoleInfo|null>(null)
const update_role=ref<userAPI.RoleInfo|null>(null)
const update_role_comment=ref<string>("")

//add user role
const add_user_role_search=ref<string>("")
const add_user_role_user=ref<userAPI.UserInfo|null>(null)
const add_user_role_role=ref<userAPI.RoleInfo|null>(null)

//permission node
const node_from=ref<string|null>(null)//empty means from current user,not empty means from specific role
const user_node=ref<permissionAPI.NodeInfo|null>(null)
const role_node=ref<permissionAPI.NodeInfo|null>(null)
const update_node=ref<permissionAPI.NodeInfo|null>(null)
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
		case "search_user":{
			let req = {
				project_id:state.project.info!.project_id,
				user_name:search.value,
				only_project:range.value=="This Project",
				page:page.value,
			}
			client.userClient.search_users({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp :userAPI.SearchUsersResp)=>{
				roles.value=[]
				if(resp.users){
					let tmp:userAPI.UserInfo[]=[]
					for(let i=0;i<resp.users.length;i++){
						if(resp.users[i]){
							tmp.push(resp.users[i]!)
						}
					}
					users.value=tmp
				}else{
					users.value=[]
				}
				page.value=resp.page
				pagesize.value=resp.pagesize
				totalsize.value=resp.totalsize
				cur_user.value=null
				cur_role.value=null
				state.clear_load()
			})
			break
		}
		case "search_role":{
			let req = {
				project_id:state.project.info!.project_id,
				role_name:search.value,
				page:page.value,
			}
			client.userClient.search_roles({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp :userAPI.SearchRolesResp)=>{
				users.value=[]
				if(resp.roles){
					let tmp:userAPI.RoleInfo[]=[]
					for(let i=0;i<resp.roles.length;i++){
						if(resp.roles[i]){
							tmp.push(resp.roles[i]!)
						}
					}
					roles.value=tmp
				}else{
					roles.value=[]
				}
				page.value=resp.page
				pagesize.value=resp.pagesize
				totalsize.value=resp.totalsize
				cur_role.value=null
				cur_user.value=null
				state.clear_load()
			})
			break
		}
		case "invite":{
			let req = {
				project_id:state.project.info!.project_id,
				user_id:invite_kick_user.value!.user_id,
			}
			client.userClient.invite_project({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(_resp: userAPI.InviteProjectResp)=>{
				invite_kick_user.value!.project_roles=[{project_id:state.project.info!.project_id,roles:[]}]
				ing.value=false
				state.clear_load()
			})
			break
		}
		case "kick":{
			let req = {
				project_id:state.project.info!.project_id,
				user_id:invite_kick_user.value!.user_id,
			}
			client.userClient.kick_project({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(_resp: userAPI.KickProjectResp)=>{
				if(range.value=="This Project"){
					for(let i=0;i<users.value.length;i++){
						if(users.value[i].user_id==invite_kick_user.value!.user_id){
							users.value.splice(i,1)
							break
						}
					}
					if(cur_user.value&&cur_user.value.user_id==invite_kick_user.value!.user_id){
						cur_user.value=null
					}
				}else{
					for(let i=0;i<invite_kick_user.value!.project_roles!.length;i++){
						if(!invite_kick_user.value!.project_roles![i]){
							continue
						}
						if(!invite_kick_user.value!.project_roles![i]!.project_id){
							continue
						}
						if(invite_kick_user.value!.project_roles![i]!.project_id![1]==state.project.info!.project_id![1]){
							invite_kick_user.value!.project_roles!.splice(i,1)
							break
						}
					}
				}
				ing.value=false
				state.clear_load()
			})
			break
		}
		case "create_role":{
			let req = {
				project_id:state.project.info!.project_id,
				role_name:create_role_name.value,
				comment:create_role_comment.value,
			}
			client.userClient.create_role({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(_resp: userAPI.CreateRoleResp)=>{
				create_role_name.value=""
				create_role_comment.value=""
				ing.value=false
				state.clear_load()
			})
			break
		}
		case "update_role":{
			let req = {
				project_id:state.project.info!.project_id,
				role_name:update_role.value!.role_name,
				new_comment:update_role_comment.value,
			}
			client.userClient.update_role({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(_resp: userAPI.UpdateRoleResp)=>{
				update_role.value!.comment=update_role_comment.value
				update_role_comment.value=""
				ing.value=false
				state.clear_load()
			})
			break
		}
		case "del_role":{
			let req = {
				project_id:state.project.info!.project_id,
				role_names:[del_role.value!.role_name],
			}
			client.userClient.del_roles({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(_resp: userAPI.DelRolesResp)=>{
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
				project_id:state.project.info!.project_id,
				user_id:cur_user.value!.user_id,
				role_name:update_user_delete_role_rolename.value,
			}
			client.userClient.del_user_role({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(_resp: userAPI.DelUserRoleResp)=>{
				let index=cur_user.value!.project_roles![0]!.roles!.indexOf(update_user_delete_role_rolename.value)
				if(index!=-1){
					cur_user.value!.project_roles![0]!.roles!.splice(index,1)
				}
				if(node_from.value==update_user_delete_role_rolename.value){
					node_from.value=null
				}
				ing.value=false
				state.clear_load()
			})
			break
		}
		case "add_user_role_missinguser":
		case "add_user_role_missingrole":{
			let req = {
				project_id:state.project.info!.project_id,
				user_id:add_user_role_user.value!.user_id,
				role_name:add_user_role_role.value!.role_name,
			}
			client.userClient.add_user_role({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(_resp: userAPI.DelUserRoleResp)=>{
				if(!add_user_role_user.value!.project_roles![0]!.roles){
					add_user_role_user.value!.project_roles![0]!.roles=[add_user_role_role.value!.role_name]
				}else if(!add_user_role_user.value!.project_roles![0]!.roles!.includes(add_user_role_role.value!.role_name)){
					add_user_role_user.value!.project_roles![0]!.roles!.push(add_user_role_role.value!.role_name)
				}
				add_user_role_user.value=null
				add_user_role_role.value=null
				ing.value=false
				state.clear_load()
			})
			break
		}
		case "get_user_permission":{
			let req = {
				project_id:state.project.info!.project_id,
				user_id:cur_user.value!.user_id,
				need_user_role_node:false,
			}
			client.permissionClient.list_user_node({"Token":state.user.token},req,client.timeout,(e :permissionAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: permissionAPI.ListUserNodeResp)=>{
				if(resp.node){
					user_node.value=resp.node
				}else{
					user_node.value=null
				}
				node_from.value=""
				state.clear_load()
			})
			break
		}
		case "update_user_permission":{
			let req = {
				user_id:cur_user.value!.user_id,
				node_id:update_node.value!.node_id,
				admin:admin.value,
				canread:canread.value,
				canwrite:canwrite.value,
			}
			client.permissionClient.update_user_permission({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(_resp :permissionAPI.UpdateUserPermissionResp)=>{
				update_node.value!.canread=canread.value
				update_node.value!.canwrite=canwrite.value
				update_node.value!.admin=admin.value
				update_node.value=null
				ing.value=false
				state.clear_load()
			})
			break
		}
		case "get_role_permission":{
			let req = {
				project_id:state.project.info!.project_id,
				role_name:cur_role.value!.role_name,
			}
			client.permissionClient.list_role_node({"Token":state.user.token},req,client.timeout,(e :permissionAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(resp: permissionAPI.ListRoleNodeResp)=>{
				if(resp.node){
					role_node.value=resp.node
				}else{
					role_node.value=null
				}
				node_from.value=cur_role.value!.role_name
				state.clear_load()
			})
			break
		}
		case "update_role_permission":{
			let req = {
				project_id:state.project.info!.project_id,
				role_name:cur_role.value!.role_name,
				node_id:update_node.value!.node_id,
				admin:admin.value,
				canread:canread.value,
				canwrite:canwrite.value,
			}
			client.permissionClient.update_role_permission({"Token":state.user.token},req,client.timeout,(e :permissionAPI.Error)=>{
				state.clear_load()
				state.set_alert("error",e.code,e.msg)
			},(_resp :permissionAPI.UpdateRolePermissionResp)=>{
				update_node.value!.canread=canread.value
				update_node.value!.canwrite=canwrite.value
				update_node.value!.admin=admin.value
				update_node.value=null
				ing.value=false
				state.clear_load()
			})
			break
		}
		default:{
			state.clear_load()
			state.set_alert("error",-2,"unknown operation")
		}
	}
}

let timeid :number = 0
function assign_search_users(part: string){
	clearTimeout(timeid)
	add_user_role_user.value=null
	users.value=[]
	if(!part){
		return
	}
	timeid = setTimeout(() => {
		if(!state.set_load()){
			return
		}
		let req = {
			project_id:state.project.info!.project_id,
			user_name:add_user_role_search.value,
			only_project:true,
			page:0,
		}
		client.userClient.search_users({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
			state.clear_load()
			state.set_alert("error",e.code,e.msg)
		},(resp :userAPI.SearchUsersResp)=>{
			if(resp.users){
				let tmp:userAPI.UserInfo[]=[]
				for(let i=0;i<resp.users.length;i++){
					if(resp.users[i]&&!user_has_role(resp.users[i]!,add_user_role_role.value!)){
						tmp.push(resp.users[i]!)
					}
				}
				users.value=tmp
			}else{
				users.value=[]
			}
			state.clear_load()
		})
	}, 500)
}
function assign_search_roles(part: string){
	clearTimeout(timeid)
	add_user_role_role.value=null
	roles.value=[]
	if(!part){
		return
	}
	timeid = setTimeout(() => {
		if(!state.set_load()){
			return
		}
		let req = {
			project_id:state.project.info!.project_id,
			role_name:add_user_role_search.value,
			page:0,
		}
		client.userClient.search_roles({"Token":state.user.token},req,client.timeout,(e :userAPI.Error)=>{
			state.clear_load()
			state.set_alert("error",e.code,e.msg)
		},(resp :userAPI.SearchRolesResp)=>{
			if(resp.roles){
				let tmp: userAPI.RoleInfo[]=[]
				for(let i=0;i<resp.roles.length;i++){
					if(resp.roles[i]&&!user_has_role(add_user_role_user.value!,resp.roles[i]!)){
						tmp.push(resp.roles[i]!)
					}
				}
				roles.value=tmp
			}else{
				roles.value=[]
			}
			state.clear_load()
		})
	}, 500)
}
function clear_assign_search(){
	clearTimeout(timeid)
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
	<va-modal v-model="ing" attach-element="#app" max-width="800px" max-height="600px" hide-default-actions no-dismiss overlay-opacity="0.2" z-index="999">
		<template #default>
			<div v-if="optype=='invite'" style="display:flex;flex-direction:column">
				<va-card  style="min-width:350px;width:auto;text-align:center" color="primary" gradient>
					<va-card-content style="font-size:20px">
						<p v-if="invite_kick_user!.feishu_user_name">
							<b>Invite user: {{ invite_kick_user!.feishu_user_name }}(feishu) join project: {{ state.project.info!.project_name}}</b>
						</p>
						<p v-if="invite_kick_user!.dingtalk_user_name">
							<b>Invite user: {{ invite_kick_user!.dingtalk_user_name }}(dingtalk) join project: {{ state.project.info!.project_name}}</b>
						</p>
						<p><b>Please confirm</b></p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:10px 10px 0 0" @click="op" gradient>Invite</va-button>
					<va-button style="width:80px;margin:10px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='kick'" style="display:flex;flex-direction:column">
				<va-card style="min-width:350px;width:auto;text-align:center" color="primary" gradient>
					<va-card-content style="font-size:20px">
						<p v-if="invite_kick_user!.feishu_user_name">
							<b>Kick user: {{ invite_kick_user!.feishu_user_name }}(feishu) out of project: {{ state.project.info!.project_name}}</b>
						</p>
						<p v-if="invite_kick_user!.dingtalk_user_name">
							<b>Kick user: {{ invite_kick_user!.dingtalk_user_name }}(feishu) out of project: {{ state.project.info!.project_name}}</b>
						</p>
						<p><b>Please confirm</b></p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:10px 10px 0 0" @click="op" gradient>Kick</va-button>
					<va-button style="width:80px;margin:10px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='create_role'" style="display:flex;flex-direction:column">
				<va-card style="min-width:350px;width:auto;text-align:center" color="primary" gradient>
					<va-card-content style="font-size:20px"><b>Create Role</b></va-card-content>
				</va-card>
				<va-input v-model.trim="create_role_name" label="New Role Name*" style="margin-top:10px" />
				<va-input v-model.trim="create_role_comment" label="New Role Comment" style="margin-top:10px" />
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:10px 10px 0 0" @click="op" gradient :disabled="create_role_name==''">Create</va-button>
					<va-button style="width:80px;margin:10px 0 0 10px" @click="create_role_name='';create_role_comment='';ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='update_role'" style="display:flex;flex-direction:column">
				<va-card style="min-width:350px;width:auto;text-align:center" color="primary" gradient>
					<va-card-content style="font-size:20px"><b>Update Role: {{ update_role!.role_name }}</b></va-card-content>
				</va-card>
				<va-input v-model.trim="update_role_comment" label="New Role Comment" style="margin-top:10px" />
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:10px 10px 0 0" :disabled="update_role_comment==update_role!.comment" @click="op" gradient>Update</va-button>
					<va-button style="width:80px;margin:10px 0 0 10px" @click="update_role=null;update_role_comment='';ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='del_role'" style="display:flex;flex-direction:column">
				<va-card style="min-width:350px;width:auto;text-align:center" color="primary" gradient>
					<va-card-content style="font-size:20px">
						<p><b>Delete role: {{ del_role!.role_name }}</b></p>
						<p><b>Please confirm</b></p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:10px 10px 0 0" @click="op" gradient>Del</va-button>
					<va-button style="width:80px;margin:10px 0 0 10px" @click="del_role=null;ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='add_user_role_missingrole'" style="display:flex;flex-direction:column">
				<va-card style="min-width:350px;width:auto;text-align:center" color="primary" gradient>
					<va-card-content style="font-size:20px">
						<p v-if="add_user_role_user!.feishu_user_name">
							<b>Assign user: {{add_user_role_user!.feishu_user_name}}(feishu) a role</b>
						</p>
						<p v-if="add_user_role_user!.dingtalk_user_name">
							<b>Assign user: {{add_user_role_user!.dingtalk_user_name}}(dingtalk) a role</b>
						</p>
					</va-card-content>
				</va-card>
				<va-input
					placeholder="Role Name*"
					style="margin-top:10px"
					v-model.trim="add_user_role_search"
					@update:modelValue="assign_search_roles($event)" />
				<div style="margin-top:10px;background-color:var(--va-background-element);height:200px;overflow-y:auto">
					<div v-if="!roles||roles.length==0" style="padding:5px">No More Roles</div>
					<va-hover stateful v-for="role of roles">
						<template #default="{hover}">
							<div
								style="padding:5px;cursor:pointer"
								:style="{'background-color':hover?'var(--va-shadow)':add_user_role_role==role?'#b6d7a8':undefined}"
								@click="add_user_role_role=role"
							>
								{{role.role_name}}
							</div>
						</template>
					</va-hover>
				</div>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:10px 10px 0 0" @click="op" :disabled="!add_user_role_role" gradient>Assign</va-button>
					<va-button style="width:80px;margin:10px 0 0 10px" @click="clear_assign_search();ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='add_user_role_missinguser'" style="display:flex;flex-direction:column">
				<va-card style="min-width:350px;width:auto;text-align:center" color="primary" gradient>
					<va-card-content style="font-size:20px"><b>Assign role: {{add_user_role_role!.role_name}} to a user</b></va-card-content>
				</va-card>
				<va-input
					placeholder="User Name*"
					style="margin-top:10px"
					v-model.trim="add_user_role_search"
					@update:modelValue="assign_search_users($event)" />
				<div style="margin-top:10px;background-color:var(--va-background-element);height:200px;overflow-y:auto">
					<div v-if="!users||users.length==0">No More Users</div>
					<va-hover stateful v-for="user of users">
						<template #default="{hover}">
							<div
								style="padding:5px;cursor:pointer"
								:style="{'background-color':hover?'var(--va-shadow)':add_user_role_user==user?'#b6d7a8':undefined}"
								@click="add_user_role_user=user"
							>
								<span v-if="user.feishu_user_name">{{user.feishu_user_name}}(feishu)</span>
								<span v-if="user.dingtalk_user_name" style="margin-left:5px">{{user.dingtalk_user_name}}(dingtalk)</span>
								<span style="color:green;margin-left:10px">{{user.user_id}}</span>
							</div>
						</template>
					</va-hover>
				</div>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:10px 10px 0 0" @click="op" :disabled="!add_user_role_user" gradient>Assign</va-button>
					<va-button style="width:80px;margin:10px 0 0 10px" @click="clear_assign_search();ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='del_user_role'" style="display:flex;flex-direction:column">
				<va-card style="min-width:350px;width:auto;text-align:center" color="primary" gradient>
					<va-card-content style="font-size:20px">
						<p v-if="cur_user!.feishu_user_name">
							<b>Remove user: {{ cur_user!.feishu_user_name }}(feishu)'s role: {{ update_user_delete_role_rolename }}</b>
						</p>
						<p v-if="cur_user!.dingtalk_user_name">
							<b>Remove user: {{ cur_user!.dingtalk_user_name }}(dingtalk)'s role: {{ update_user_delete_role_rolename }}</b>
						</p>
						<p><b>Please confirm</b></p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:10px 10px 0 0" @click="op" gradient>Del</va-button>
					<va-button style="width:80px;margin:10px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='update_role_permission'" style="display:flex;flex-direction:column">
				<va-card style="min-width:350px;width:auto;text-align:center" color="primary" gradient>
					<va-card-content style="font-size:20px">
						<p><b>Update role: {{ cur_role!.role_name }}'s permission on node: {{ update_node!.node_name }}</b></p>
						<p><b>Please confirm</b></p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:10px 10px 0 0" @click="op" gradient>Update</va-button>
					<va-button style="width:80px;margin:10px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
			<div v-else-if="optype=='update_user_permission'" style="display:flex;flex-direction:column">
				<va-card style="min-width:350px;width:auto;text-align:center" color="primary" gradient>
					<va-card-content style="font-size:20px">
						<p v-if="cur_user!.feishu_user_name">
							<b>Update user: {{ cur_user!.feishu_user_name}}(feishu)'s permission on node: {{ update_node!.node_name }}</b>
						</p>
						<p v-if="cur_user!.dingtalk_user_name">
							<b>Update user: {{ cur_user!.dingtalk_user_name}}(dingtalk)'s permission on node: {{ update_node!.node_name }}</b>
						</p>
						<p><b>Please confirm</b></p>
					</va-card-content>
				</va-card>
				<div style="display:flex;justify-content:center">
					<va-button style="width:80px;margin:10px 10px 0 0" @click="op" gradient>Update</va-button>
					<va-button style="width:80px;margin:10px 0 0 10px" @click="ing=false" gradient>Cancel</va-button>
				</div>
			</div>
		</template>
	</va-modal>
	<div style="flex:1;display:flex;flex-direction:column;margin:1px;overflow-y:auto">
		<div style="display:flex;margin:1px;align-self:center">
			<va-select
				v-model="target"
				:options="targets"
				dropdownIcon=""
				style="width:100px;margin-right:1px"
				outline
				trigger="hover"
				:hoverOverTimeout="0"
				:hoverOutTimeout="100"
			>
				<template #option='{option,selectOption}'>
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
							}
							selectOption(option)
						}">
						<template #default="{hover}">
							<div
								style="padding:10px;cursor:pointer"
								:style="{'background-color':hover?'var(--va-background-border)':'',color:target==option?'green':'black'}"
							>
								{{option}}
							</div>
						</template>
					</va-hover>
				</template>
			</va-select>
			<va-select
				v-if="target=='User'"
				v-model="range"
				:options="ranges"
				dropdownIcon=""
				style="width:150px;margin:0 1px"
				outline
				trigger="hover"
				:hoverOverTimeout="0"
				:hoverOutTimeout="100"
			>
				<template #option="{ option,selectOption }">
					<va-hover stateful @click="()=>{
							if(option!=range){
								page=1
								pagesize=0
								totalsize=0
								users=[]
								cur_user=null
							}
							selectOption(option)
						}">
						<template #default="{hover}">
							<div
								style="padding:10px;cursor:pointer"
								:style="{'background-color':hover?'var(--va-background-border)':'',color:range==option?'green':'black'}"
							>
								{{option}}
							</div>
						</template>
					</va-hover>
				</template>
			</va-select>
			<va-input :placeholder="target=='User'?'User Name':'Role Name'" outline style="max-width:250px;margin:0 1px" v-model.trim="search" />
			<va-button v-if="target=='User'" style="margin:0 1px" @click="optype='search_user';op()" gradient>Search</va-button>
			<va-button v-if="target=='Role'" style="margin:0 1px" @click="optype='search_role';op()" gradient>Search</va-button>
			<va-button v-if="target=='Role'" style="margin-left:1px" @click="optype='create_role';ing=true" gradient>Create</va-button>
		</div>
		<div v-if="target=='User'" style="flex:1;display:flex;flex-direction:column;margin:1px;overflow-y:auto">
			<template v-for="user of users">
				<div
					v-if="!cur_user||cur_user==user"
					style="display:flex;align-items:center;margin:1px 0"
					:style="user_bindstyle(user)"
					@click="()=>{
						if(!cur_user&&invited(user)){
							cur_user=user
							node_from=null
						}else{
							cur_user=null
							user_node=null
						}
					}"
					@mouseover="userhover=user"
					@mouseout="userhover=null"
				>
					<span style="width:40px;padding:12px 20px;color:var(--va-primary)">{{cur_user==user?'-':invited(user)?'+':' ' }}</span>
					<span v-if="user.feishu_user_name" style="padding:12px 0px 12px 20px;color:var(--va-primary)">{{user.feishu_user_name}}(feishu)</span>
					<span v-if="user.dingtalk_user_name" style="padding:12px 0px 12px 20px;color:var(--va-primary)">{{user.dingtalk_user_name}}(dingtalk)</span>
					<span style="padding:12px 20px;color:green">{{user.user_id}}</span>
					<span style="flex:1"></span>
					<span style="padding:12px;color:green">Create Time: {{parsetime(user.ctime)}}</span>
					<va-button
						v-if="state.page.node!.admin&&invited(user)"
						style="min-width:60px;height:30px;margin:0 10px"
						size="small"
						gradient
						@mouseover.stop=""
						@mouseout.stop=""
						@click.stop="add_user_role_search='';add_user_role_user=user;roles=[];optype='add_user_role_missingrole';ing=true">
						AddRole
					</va-button>
					<va-button
						v-if="state.page.node!.admin"
						style="min-width:60px;height:30px;margin:0 10px"
						size="small"
						gradient
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
				<va-split
					v-if="cur_user==user"
					style="margin:2px 10px;display:flex;flex:1;overflow-y:auto;color:var(--va-primary)"
					stateful
					:model-value='0'
					:limits="['250px',50]">
					<template #start>
						<div style="height:99%;display:flex;flex-direction:column;overflow-y:auto">
							<va-hover stateful style="flex:1;margin:2px 0">
								<template #default="{hover}">
									<div
										style="padding:12px;cursor:pointer"
										:style="{'background-color':hover?'var(--va-shadow)':node_from==''?'#b6d7a8':'var(--va-background-element)'}"
										@click="optype='get_user_permission';op()">
										User Self Permissions
									</div>
								</template>
							</va-hover>
							<va-hover v-for="rolename of user.project_roles![0]!.roles!" stateful style="flex:1;margin:2px 0">
								<template #default="{hover}">
									<div
										style="display:flex;align-items:center;cursor:pointer"
										:style="{'background-color':hover?'var(--va-shadow)':node_from==rolename?'#b6d7a8':'var(--va-background-element)'}"
										@click="cur_role={project_id:state.project.info!.project_id,role_name:rolename,comment:'',ctime:0};
											optype='get_role_permission';
											op()">
										<div style="flex:1;padding:12px;white-space:nowrap">Role Permissions: {{rolename}}</div>
										<va-button
											v-if="state.page.node!.admin"
											size="small"
											style="margin:0 2px"
											gradient
											@mouseenter.stop=""
											@mouseover.stop=""
											@mouseout.stop=""
											@mouseleave.stop=""
											@click.stop="update_user_delete_role_rolename=rolename;optype='del_user_role';ing=true">
											X
										</va-button>
									</div>
								</template>
							</va-hover>
						</div>
					</template>
					<template #end>
						<div v-if="node_from==''&&user_node"
							style="height:99%;margin:2px 0;display:flex;background-color:var(--va-background-element);color:var(--va-primary);overflow:auto">
							<nodetree
								:pnode="null"
								:node="user_node"
								:deep="0"
								:disabled="false"
								@permissionevent="(updatenode,r,w,a)=>{
									update_node=updatenode;
									canread=r;
									canwrite=w;
									admin=a;
									optype='update_user_permission';
									ing=true
								}"/>
						</div>
						<div v-else-if="node_from==''"
							style="height:99%;margin:2px 0;display:flex;justify-content:center;align-items:center;background-color:var(--va-background-element)">
							<b>No permission</b>
						</div>
						<div v-if="node_from!=null&&node_from!=''&&role_node&&role_node.children&&role_node.children.length>0"
							style="height:99%;margin:2px 0;display:flex;background-color:var(--va-background-element);color:var(--va-primary);overflow:auto">
							<template v-for="child of role_node.children">
								<nodetree v-if="child" :pnode="role_node" :node="child" :deep="0" disabled/>
							</template>
						</div>
						<div v-else-if="node_from!=null&&node_from!=''"
							style="height:99%;margin:2px 0;display:flex;justify-content:center;align-items:center;background-color:var(--va-background-element)">
							<b>No permission</b>
						</div>
					</template>
				</va-split>
			</template>
		</div>
		<div v-if="target=='Role'" style="flex:1;display:flex;flex-direction:column;margin:1px;overflow-y:auto">
			<template v-for="role of roles">
				<div
					v-if="!Boolean(cur_role)||cur_role==role"
					style="display:flex;align-items:center;margin:1px 0;cursor:pointer"
					:style="{'background-color':rolehover==role?'var(--va-shadow)':'var(--va-background-element)'}"
					@click="()=>{
						if(!cur_role){
							cur_role=role
							optype='get_role_permission'
							op()
						}else{
							cur_role=null
							role_node=null
						}
					}"
					@mouseover="rolehover=role"
					@mouseout="rolehover=null"
				>
					<span style="width:40px;padding:12px 20px;color:var(--va-primary)">{{ cur_role==role?'-':'+' }}</span>
					<span style="padding:12px 20px;color:var(--va-primary)">{{role.role_name}}</span>
					<span style="flex:1"></span>
					<span style="padding:12px;color:green">Create Time: {{parsetime(role.ctime)}}</span>
					<va-button
						v-if="state.page.node!.admin"
						style="width:60px;height:30px;margin:0 10px"
						size="small"
						gradient
						@mouseover.stop=""
						@mouseout.stop=""
						@click.stop="add_user_role_search='';add_user_role_role=role;users=[];optype='add_user_role_missinguser';ing=true">
						Assign
					</va-button>
					<va-button
						v-if="state.page.node!.canwrite||state.page.node!.admin"
						style="width:60px;height:30px;margin:0 10px"
						size="small"
						gradient
						@mouseover.stop=""
						@mouseout.stop=""
						@click.stop="update_role=role;update_role_comment=role.comment;optype='update_role';ing=true"
					>
						Update
					</va-button>
					<va-button
						v-if="state.page.node!.admin"
						style="width:60px;height:30px;margin:0 10px"
						size="small"
						gradient
						@mouseover.stop=""
						@mouseout.stop=""
						@click.stop="del_role=role;optype='del_role';ing=true"
					>
						Del
					</va-button>
				</div>
				<textarea
					v-if="cur_role==role"
					style="border:1px solid var(--va-background-element);border-radius:5px;margin:1px 10px;resize:none"
					readonly
					v-model="role.comment" />
				<div v-if="cur_role==role&&role_node&&role_node.children&&role_node.children.length>0"
					style="flex:1;margin:2px 10px;display:flex;background-color:var(--va-background-element);color:var(--va-primary);overflow:auto">
					<template v-for="child of role_node.children">
						<nodetree
							v-if="child"
							:pnode="role_node"
							:node="child"
							:deep="0"
							:disabled="false"
							@permissionevent="(updatenode,r,w,a)=>{
								update_node=updatenode;
								canread=r;
								canwrite=w;
								admin=a;
								optype='update_role_permission';
								ing=true
							}"/>
					</template>
				</div>
				<div v-else-if="cur_role==role"
					style="flex:1;margin:1px 10px;display:flex;justify-content:center;align-items:center;background-color:var(--va-background-element);color:var(--va-primary)">
					<b>No permission</b>
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
