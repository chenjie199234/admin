// Code generated by protoc-gen-browser. DO NOT EDIT.
// version:
// 	protoc-gen-browser v0.0.138<br />
// 	protoc             v6.31.0<br />
// source: api/admin_user.proto<br />

// @ts-nocheck
export interface LogicError{
	code: number;
	msg: string;
}

export class AddUserRoleReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	user_id: string = ''
	role_name: string = ''
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.user_id){
			tmp["user_id"]=this.user_id
		}
		if(this.role_name){
			tmp["role_name"]=this.role_name
		}
		return tmp
	}
}
export class AddUserRoleResp{
	fromOBJ(_obj:Object){
	}
}
export class CreateRoleReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	role_name: string = ''
	comment: string = ''
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.role_name){
			tmp["role_name"]=this.role_name
		}
		if(this.comment){
			tmp["comment"]=this.comment
		}
		return tmp
	}
}
export class CreateRoleResp{
	fromOBJ(_obj:Object){
	}
}
export class DelRolesReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	role_names: Array<string>|null = null
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.role_names && this.role_names.length>0){
			tmp["role_names"]=this.role_names
		}
		return tmp
	}
}
export class DelRolesResp{
	fromOBJ(_obj:Object){
	}
}
export class DelUserRoleReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	user_id: string = ''
	role_name: string = ''
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.user_id){
			tmp["user_id"]=this.user_id
		}
		if(this.role_name){
			tmp["role_name"]=this.role_name
		}
		return tmp
	}
}
export class DelUserRoleResp{
	fromOBJ(_obj:Object){
	}
}
export class GetOauth2Req{
	src_type: string = ''
	toJSON(){
		let tmp = {}
		if(this.src_type){
			tmp["src_type"]=this.src_type
		}
		return tmp
	}
}
export class GetOauth2Resp{
	url: string = ''
	fromOBJ(obj:Object){
		if(obj["url"]){
			this.url=obj["url"]
		}
	}
}
export class InviteProjectReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	user_id: string = ''
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.user_id){
			tmp["user_id"]=this.user_id
		}
		return tmp
	}
}
export class InviteProjectResp{
	fromOBJ(_obj:Object){
	}
}
export class KickProjectReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	user_id: string = ''
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.user_id){
			tmp["user_id"]=this.user_id
		}
		return tmp
	}
}
export class KickProjectResp{
	fromOBJ(_obj:Object){
	}
}
export class LoginInfoReq{
	toJSON(){
		let tmp = {}
		return tmp
	}
}
export class LoginInfoResp{
	user: UserInfo|null = null
	fromOBJ(obj:Object){
		if(obj["user"]){
			this.user=new UserInfo()
			this.user.fromOBJ(obj["user"])
		}
	}
}
export class ProjectRoles{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	roles: Array<string>|null = null
	fromOBJ(obj:Object){
		if(obj["project_id"] && obj["project_id"].length>0){
			this.project_id=obj["project_id"]
		}
		if(obj["roles"] && obj["roles"].length>0){
			this.roles=obj["roles"]
		}
	}
}
export class RoleInfo{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	role_name: string = ''
	comment: string = ''
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	ctime: number = 0
	fromOBJ(obj:Object){
		if(obj["project_id"] && obj["project_id"].length>0){
			this.project_id=obj["project_id"]
		}
		if(obj["role_name"]){
			this.role_name=obj["role_name"]
		}
		if(obj["comment"]){
			this.comment=obj["comment"]
		}
		if(obj["ctime"]){
			this.ctime=obj["ctime"]
		}
	}
}
export class SearchRolesReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	role_name: string = ''//fuzzy search
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	page: number = 0//page starts from 1,if page is 0,means return all result
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.role_name){
			tmp["role_name"]=this.role_name
		}
		if(this.page){
			tmp["page"]=this.page
		}
		return tmp
	}
}
export class SearchRolesResp{
	roles: Array<RoleInfo|null>|null = null
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	page: number = 0
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	pagesize: number = 0
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	totalsize: number = 0
	fromOBJ(obj:Object){
		if(obj["roles"] && obj["roles"].length>0){
			this.roles=new Array<RoleInfo|null>()
			for(let value of obj["roles"]){
				if(value){
					let tmp=new RoleInfo()
					tmp.fromOBJ(value)
					this.roles.push(tmp)
				}else{
					this.roles.push(null)
				}
			}
		}
		if(obj["page"]){
			this.page=obj["page"]
		}
		if(obj["pagesize"]){
			this.pagesize=obj["pagesize"]
		}
		if(obj["totalsize"]){
			this.totalsize=obj["totalsize"]
		}
	}
}
export class SearchUsersReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	user_name: string = ''//fuzzy search
	//true - search users in the project,this require operator has read permission on this project's user control node
	//false - search all users(include users not in this project),this require operator has admin permission on this project's user control node
	only_project: boolean = false
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	page: number = 0//if page is 0,means return all result
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.user_name){
			tmp["user_name"]=this.user_name
		}
		if(this.only_project){
			tmp["only_project"]=this.only_project
		}
		if(this.page){
			tmp["page"]=this.page
		}
		return tmp
	}
}
export class SearchUsersResp{
	users: Array<UserInfo|null>|null = null//key userid,value userinfo(only contains the required project's roles)
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	page: number = 0
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	pagesize: number = 0
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	totalsize: number = 0
	fromOBJ(obj:Object){
		if(obj["users"] && obj["users"].length>0){
			this.users=new Array<UserInfo|null>()
			for(let value of obj["users"]){
				if(value){
					let tmp=new UserInfo()
					tmp.fromOBJ(value)
					this.users.push(tmp)
				}else{
					this.users.push(null)
				}
			}
		}
		if(obj["page"]){
			this.page=obj["page"]
		}
		if(obj["pagesize"]){
			this.pagesize=obj["pagesize"]
		}
		if(obj["totalsize"]){
			this.totalsize=obj["totalsize"]
		}
	}
}
export class UpdateRoleReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	role_name: string = ''
	new_comment: string = ''//if didn't change,set this with the old value
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.role_name){
			tmp["role_name"]=this.role_name
		}
		if(this.new_comment){
			tmp["new_comment"]=this.new_comment
		}
		return tmp
	}
}
export class UpdateRoleResp{
	fromOBJ(_obj:Object){
	}
}
export class UserInfo{
	user_id: string = ''
	feishu_user_name: string = ''
	dingding_user_name: string = ''
	wxwork_user_name: string = ''
	//Warning!!!Type is uint32,be careful of sign(+) and overflow
	ctime: number = 0//timestamp,uint seconds
	project_roles: Array<ProjectRoles|null>|null = null
	fromOBJ(obj:Object){
		if(obj["user_id"]){
			this.user_id=obj["user_id"]
		}
		if(obj["feishu_user_name"]){
			this.feishu_user_name=obj["feishu_user_name"]
		}
		if(obj["dingding_user_name"]){
			this.dingding_user_name=obj["dingding_user_name"]
		}
		if(obj["wxwork_user_name"]){
			this.wxwork_user_name=obj["wxwork_user_name"]
		}
		if(obj["ctime"]){
			this.ctime=obj["ctime"]
		}
		if(obj["project_roles"] && obj["project_roles"].length>0){
			this.project_roles=new Array<ProjectRoles|null>()
			for(let value of obj["project_roles"]){
				if(value){
					let tmp=new ProjectRoles()
					tmp.fromOBJ(value)
					this.project_roles.push(tmp)
				}else{
					this.project_roles.push(null)
				}
			}
		}
	}
}
export class UserLoginReq{
	src_type: string = ''
	code: string = ''
	toJSON(){
		let tmp = {}
		if(this.src_type){
			tmp["src_type"]=this.src_type
		}
		if(this.code){
			tmp["code"]=this.code
		}
		return tmp
	}
}
export class UserLoginResp{
	token: string = ''
	fromOBJ(obj:Object){
		if(obj["token"]){
			this.token=obj["token"]
		}
	}
}
//timeout's unit is millisecond,it will be used when > 0
function call(timeout: number,url: string,opts: Object,error: (arg: LogicError)=>void,success: (arg: Object)=>void){
	let tid: number|null = null
	if(timeout>0){
		const c = new AbortController()
		opts["signal"] = c.signal
		tid = setTimeout(()=>{c.abort()},timeout)
	}
	let ok=false
	fetch(url,opts)
	.then(r=>{
		ok=r.ok
		if(r.ok){
			return r.json()
		}
		return r.text()
	})
	.then(d=>{
		if(!ok){
			throw d
		}
		success(d)
	})
	.catch(e=>{
		if(e instanceof Error){
			error({code:-1,msg:e.message})
		}else if(e.length>0 && e[0]=='{' && e[e.length-1]=='}'){
			error(JSON.parse(e))
		}else{
			error({code:-1,msg:e})
		}
	})
	.finally(()=>{
		if(tid){
			clearTimeout(tid)
		}
	})
}
const _WebPathUserGetOauth2: string ="/admin.user/get_oauth2";
const _WebPathUserUserLogin: string ="/admin.user/user_login";
const _WebPathUserLoginInfo: string ="/admin.user/login_info";
const _WebPathUserInviteProject: string ="/admin.user/invite_project";
const _WebPathUserKickProject: string ="/admin.user/kick_project";
const _WebPathUserSearchUsers: string ="/admin.user/search_users";
const _WebPathUserCreateRole: string ="/admin.user/create_role";
const _WebPathUserSearchRoles: string ="/admin.user/search_roles";
const _WebPathUserUpdateRole: string ="/admin.user/update_role";
const _WebPathUserDelRoles: string ="/admin.user/del_roles";
const _WebPathUserAddUserRole: string ="/admin.user/add_user_role";
const _WebPathUserDelUserRole: string ="/admin.user/del_user_role";
export class UserBrowserClient {
	constructor(host: string){
		if(!host || host.length==0){
			throw "UserBrowserClient's host missing"
		}
		this.host=host
	}
	//timeout's unit is millisecond,it will be used when > 0
	get_oauth2(header: Object,req: GetOauth2Req,timeout: number,error: (arg: LogicError)=>void,success: (arg: GetOauth2Resp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathUserGetOauth2,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new GetOauth2Resp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	user_login(header: Object,req: UserLoginReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: UserLoginResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathUserUserLogin,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new UserLoginResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	login_info(header: Object,req: LoginInfoReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: LoginInfoResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathUserLoginInfo,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new LoginInfoResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	invite_project(header: Object,req: InviteProjectReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: InviteProjectResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathUserInviteProject,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new InviteProjectResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	kick_project(header: Object,req: KickProjectReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: KickProjectResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathUserKickProject,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new KickProjectResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	search_users(header: Object,req: SearchUsersReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: SearchUsersResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathUserSearchUsers,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new SearchUsersResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	create_role(header: Object,req: CreateRoleReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: CreateRoleResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathUserCreateRole,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new CreateRoleResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	search_roles(header: Object,req: SearchRolesReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: SearchRolesResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathUserSearchRoles,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new SearchRolesResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	update_role(header: Object,req: UpdateRoleReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: UpdateRoleResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathUserUpdateRole,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new UpdateRoleResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	del_roles(header: Object,req: DelRolesReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: DelRolesResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathUserDelRoles,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new DelRolesResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	add_user_role(header: Object,req: AddUserRoleReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: AddUserRoleResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathUserAddUserRole,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new AddUserRoleResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	del_user_role(header: Object,req: DelUserRoleReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: DelUserRoleResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathUserDelUserRole,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new DelUserRoleResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	private host: string
}
