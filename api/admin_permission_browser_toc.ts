// Code generated by protoc-gen-browser. DO NOT EDIT.
// version:
// 	protoc-gen-browser v0.0.113<br />
// 	protoc             v5.26.1<br />
// source: api/admin_permission.proto<br />

export interface LogicError{
	code: number;
	msg: string;
}

export class AddNodeReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	pnode_id: Array<number>|null = null
	node_name: string = ''
	node_data: string = ''
	toJSON(){
		let tmp = {}
		if(this.pnode_id && this.pnode_id.length>0){
			tmp["pnode_id"]=this.pnode_id
		}
		if(this.node_name){
			tmp["node_name"]=this.node_name
		}
		if(this.node_data){
			tmp["node_data"]=this.node_data
		}
		return tmp
	}
}
export class AddNodeResp{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	node_id: Array<number>|null = null
	fromOBJ(obj:Object){
		if(obj["node_id"] && obj["node_id"].length>0){
			this.node_id=obj["node_id"]
		}
	}
}
export class DelNodeReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	node_id: Array<number>|null = null
	toJSON(){
		let tmp = {}
		if(this.node_id && this.node_id.length>0){
			tmp["node_id"]=this.node_id
		}
		return tmp
	}
}
export class DelNodeResp{
	fromOBJ(_obj:Object){
	}
}
export class GetUserPermissionReq{
	user_id: string = ''
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	node_id: Array<number>|null = null
	toJSON(){
		let tmp = {}
		if(this.user_id){
			tmp["user_id"]=this.user_id
		}
		if(this.node_id && this.node_id.length>0){
			tmp["node_id"]=this.node_id
		}
		return tmp
	}
}
export class GetUserPermissionResp{
	canread: boolean = false
	canwrite: boolean = false
	admin: boolean = false
	fromOBJ(obj:Object){
		if(obj["canread"]){
			this.canread=obj["canread"]
		}
		if(obj["canwrite"]){
			this.canwrite=obj["canwrite"]
		}
		if(obj["admin"]){
			this.admin=obj["admin"]
		}
	}
}
export class ListProjectNodeReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		return tmp
	}
}
export class ListProjectNodeResp{
	//this will only return the node name,node data and children
	//other node's info will not return
	node: NodeInfo|null = null
	fromOBJ(obj:Object){
		if(obj["node"]){
			this.node=new NodeInfo()
			this.node.fromOBJ(obj["node"])
		}
	}
}
export class ListRoleNodeReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	role_name: string = ''
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.role_name){
			tmp["role_name"]=this.role_name
		}
		return tmp
	}
}
export class ListRoleNodeResp{
	node: NodeInfo|null = null
	fromOBJ(obj:Object){
		if(obj["node"]){
			this.node=new NodeInfo()
			this.node.fromOBJ(obj["node"])
		}
	}
}
export class ListUserNodeReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	user_id: string = ''//if this is empty means return self's
	need_user_role_node: boolean = false//false - only return user's base node,true - return user's base node and user's roles' node
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.user_id){
			tmp["user_id"]=this.user_id
		}
		if(this.need_user_role_node){
			tmp["need_user_role_node"]=this.need_user_role_node
		}
		return tmp
	}
}
export class ListUserNodeResp{
	node: NodeInfo|null = null
	fromOBJ(obj:Object){
		if(obj["node"]){
			this.node=new NodeInfo()
			this.node.fromOBJ(obj["node"])
		}
	}
}
export class MoveNodeReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	node_id: Array<number>|null = null
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	pnode_id: Array<number>|null = null
	toJSON(){
		let tmp = {}
		if(this.node_id && this.node_id.length>0){
			tmp["node_id"]=this.node_id
		}
		if(this.pnode_id && this.pnode_id.length>0){
			tmp["pnode_id"]=this.pnode_id
		}
		return tmp
	}
}
export class MoveNodeResp{
	fromOBJ(_obj:Object){
	}
}
export class NodeInfo{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	node_id: Array<number>|null = null
	node_name: string = ''
	node_data: string = ''
	canread: boolean = false
	canwrite: boolean = false
	admin: boolean = false
	children: Array<NodeInfo|null>|null = null
	fromOBJ(obj:Object){
		if(obj["node_id"] && obj["node_id"].length>0){
			this.node_id=obj["node_id"]
		}
		if(obj["node_name"]){
			this.node_name=obj["node_name"]
		}
		if(obj["node_data"]){
			this.node_data=obj["node_data"]
		}
		if(obj["canread"]){
			this.canread=obj["canread"]
		}
		if(obj["canwrite"]){
			this.canwrite=obj["canwrite"]
		}
		if(obj["admin"]){
			this.admin=obj["admin"]
		}
		if(obj["children"] && obj["children"].length>0){
			this.children=new Array<NodeInfo|null>()
			for(let value of obj["children"]){
				if(value){
					let tmp=new NodeInfo()
					tmp.fromOBJ(value)
					this.children.push(tmp)
				}else{
					this.children.push(null)
				}
			}
		}
	}
}
export class UpdateNodeReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	node_id: Array<number>|null = null
	new_node_name: string = ''//if didn't change,set this with the old value
	new_node_data: string = ''//if didn't change,set this with the old value
	toJSON(){
		let tmp = {}
		if(this.node_id && this.node_id.length>0){
			tmp["node_id"]=this.node_id
		}
		if(this.new_node_name){
			tmp["new_node_name"]=this.new_node_name
		}
		if(this.new_node_data){
			tmp["new_node_data"]=this.new_node_data
		}
		return tmp
	}
}
export class UpdateNodeResp{
	fromOBJ(_obj:Object){
	}
}
export class UpdateRolePermissionReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	role_name: string = ''
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	node_id: Array<number>|null = null
	//if admin is true,canread and canwrite will be ignore
	admin: boolean = false
	//if admin is false,and canread is false too,means delete this user from this node
	//if admin is false,and canread is false and node_id's length is 1,means delete this user completely
	canread: boolean = false
	//if canwrite is true,canread must be true too
	canwrite: boolean = false
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.role_name){
			tmp["role_name"]=this.role_name
		}
		if(this.node_id && this.node_id.length>0){
			tmp["node_id"]=this.node_id
		}
		if(this.admin){
			tmp["admin"]=this.admin
		}
		if(this.canread){
			tmp["canread"]=this.canread
		}
		if(this.canwrite){
			tmp["canwrite"]=this.canwrite
		}
		return tmp
	}
}
export class UpdateRolePermissionResp{
	fromOBJ(_obj:Object){
	}
}
export class UpdateUserPermissionReq{
	user_id: string = ''
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	node_id: Array<number>|null = null
	//if admin is true,canread and canwrite will be ignore
	admin: boolean = false
	//if admin is false,and canread is false too,means delete this user from this node
	//if admin is false,and canread is false and node_id's length is 1,means delete this user completely
	canread: boolean = false
	//if canwrite is true,canread must be true too
	canwrite: boolean = false
	toJSON(){
		let tmp = {}
		if(this.user_id){
			tmp["user_id"]=this.user_id
		}
		if(this.node_id && this.node_id.length>0){
			tmp["node_id"]=this.node_id
		}
		if(this.admin){
			tmp["admin"]=this.admin
		}
		if(this.canread){
			tmp["canread"]=this.canread
		}
		if(this.canwrite){
			tmp["canwrite"]=this.canwrite
		}
		return tmp
	}
}
export class UpdateUserPermissionResp{
	fromOBJ(_obj:Object){
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
const _WebPathPermissionGetUserPermission: string ="/admin.permission/get_user_permission";
const _WebPathPermissionUpdateUserPermission: string ="/admin.permission/update_user_permission";
const _WebPathPermissionUpdateRolePermission: string ="/admin.permission/update_role_permission";
const _WebPathPermissionAddNode: string ="/admin.permission/add_node";
const _WebPathPermissionUpdateNode: string ="/admin.permission/update_node";
const _WebPathPermissionMoveNode: string ="/admin.permission/move_node";
const _WebPathPermissionDelNode: string ="/admin.permission/del_node";
const _WebPathPermissionListUserNode: string ="/admin.permission/list_user_node";
const _WebPathPermissionListRoleNode: string ="/admin.permission/list_role_node";
const _WebPathPermissionListProjectNode: string ="/admin.permission/list_project_node";
//ToC means this is for users
export class PermissionBrowserClientToC {
	constructor(host: string){
		if(!host || host.length==0){
			throw "PermissionBrowserClientToC's host missing"
		}
		this.host=host
	}
	//timeout's unit is millisecond,it will be used when > 0
	get_user_permission(header: Object,req: GetUserPermissionReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: GetUserPermissionResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathPermissionGetUserPermission,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new GetUserPermissionResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	update_user_permission(header: Object,req: UpdateUserPermissionReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: UpdateUserPermissionResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathPermissionUpdateUserPermission,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new UpdateUserPermissionResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	update_role_permission(header: Object,req: UpdateRolePermissionReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: UpdateRolePermissionResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathPermissionUpdateRolePermission,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new UpdateRolePermissionResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	add_node(header: Object,req: AddNodeReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: AddNodeResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathPermissionAddNode,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new AddNodeResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	update_node(header: Object,req: UpdateNodeReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: UpdateNodeResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathPermissionUpdateNode,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new UpdateNodeResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	move_node(header: Object,req: MoveNodeReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: MoveNodeResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathPermissionMoveNode,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new MoveNodeResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	del_node(header: Object,req: DelNodeReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: DelNodeResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathPermissionDelNode,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new DelNodeResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	list_user_node(header: Object,req: ListUserNodeReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: ListUserNodeResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathPermissionListUserNode,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new ListUserNodeResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	list_role_node(header: Object,req: ListRoleNodeReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: ListRoleNodeResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathPermissionListRoleNode,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new ListRoleNodeResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	list_project_node(header: Object,req: ListProjectNodeReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: ListProjectNodeResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathPermissionListProjectNode,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new ListProjectNodeResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	private host: string
}
