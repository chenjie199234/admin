// Code generated by protoc-gen-browser. DO NOT EDIT.
// version:
// 	protoc-gen-browser v0.0.81<br />
// 	protoc             v4.24.1<br />
// source: api/admin_permission.proto<br />

import Axios from "axios";

export interface Error{
	code: number;
	msg: string;
}

export interface AddNodeReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	pnode_id: Array<number>|null|undefined;
	node_name: string;
	node_data: string;
}
function AddNodeReqToJson(msg: AddNodeReq): string{
	let s: string="{"
	//pnode_id
	if(msg.pnode_id==null||msg.pnode_id==undefined){
		s+='"pnode_id":null,'
	}else if(msg.pnode_id.length==0){
		s+='"pnode_id":[],'
	}else{
		s+='"pnode_id":['
		for(let element of msg.pnode_id){
			if(element==null||element==undefined||!Number.isInteger(element)){
				throw 'element in AddNodeReq.pnode_id must be integer'
			}
			if(element>4294967295||element<0){
				throw 'element in AddNodeReq.pnode_id overflow'
			}
			s+=element+','
		}
		s=s.substr(0,s.length-1)+'],'
	}
	//node_name
	if(msg.node_name==null||msg.node_name==undefined){
		throw 'AddNodeReq.node_name must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.node_name)
		s+='"node_name":'+vv+','
	}
	//node_data
	if(msg.node_data==null||msg.node_data==undefined){
		throw 'AddNodeReq.node_data must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.node_data)
		s+='"node_data":'+vv+','
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface AddNodeResp{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	node_id: Array<number>|null|undefined;
}
function JsonToAddNodeResp(jsonobj: { [k:string]:any }): AddNodeResp{
	let obj: AddNodeResp={
		node_id:null,
	}
	//node_id
	if(jsonobj['node_id']!=null&&jsonobj['node_id']!=undefined){
		if(!(jsonobj['node_id'] instanceof Array)){
			throw 'AddNodeResp.node_id must be Array<number>|null|undefined'
		}
		for(let element of jsonobj['node_id']){
			if(typeof element!='number'||!Number.isInteger(element)){
				throw 'element in AddNodeResp.node_id must be integer'
			}else if(element>4294967295||element<0){
				throw 'element in AddNodeResp.node_id overflow'
			}
			if(obj['node_id']==null){
				obj['node_id']=new Array<number>
			}
			obj['node_id'].push(element)
		}
	}
	return obj
}
export interface DelNodeReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	node_id: Array<number>|null|undefined;
}
function DelNodeReqToJson(msg: DelNodeReq): string{
	let s: string="{"
	//node_id
	if(msg.node_id==null||msg.node_id==undefined){
		s+='"node_id":null,'
	}else if(msg.node_id.length==0){
		s+='"node_id":[],'
	}else{
		s+='"node_id":['
		for(let element of msg.node_id){
			if(element==null||element==undefined||!Number.isInteger(element)){
				throw 'element in DelNodeReq.node_id must be integer'
			}
			if(element>4294967295||element<0){
				throw 'element in DelNodeReq.node_id overflow'
			}
			s+=element+','
		}
		s=s.substr(0,s.length-1)+'],'
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface DelNodeResp{
}
function JsonToDelNodeResp(_jsonobj: { [k:string]:any }): DelNodeResp{
	let obj: DelNodeResp={
	}
	return obj
}
export interface GetUserPermissionReq{
	user_id: string;
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	node_id: Array<number>|null|undefined;
}
function GetUserPermissionReqToJson(msg: GetUserPermissionReq): string{
	let s: string="{"
	//user_id
	if(msg.user_id==null||msg.user_id==undefined){
		throw 'GetUserPermissionReq.user_id must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.user_id)
		s+='"user_id":'+vv+','
	}
	//node_id
	if(msg.node_id==null||msg.node_id==undefined){
		s+='"node_id":null,'
	}else if(msg.node_id.length==0){
		s+='"node_id":[],'
	}else{
		s+='"node_id":['
		for(let element of msg.node_id){
			if(element==null||element==undefined||!Number.isInteger(element)){
				throw 'element in GetUserPermissionReq.node_id must be integer'
			}
			if(element>4294967295||element<0){
				throw 'element in GetUserPermissionReq.node_id overflow'
			}
			s+=element+','
		}
		s=s.substr(0,s.length-1)+'],'
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface GetUserPermissionResp{
	canread: boolean;
	canwrite: boolean;
	admin: boolean;
}
function JsonToGetUserPermissionResp(jsonobj: { [k:string]:any }): GetUserPermissionResp{
	let obj: GetUserPermissionResp={
		canread:false,
		canwrite:false,
		admin:false,
	}
	//canread
	if(jsonobj['canread']!=null&&jsonobj['canread']!=undefined){
		if(typeof jsonobj['canread']!='boolean'){
			throw 'GetUserPermissionResp.canread must be boolean'
		}
		obj['canread']=jsonobj['canread']
	}
	//canwrite
	if(jsonobj['canwrite']!=null&&jsonobj['canwrite']!=undefined){
		if(typeof jsonobj['canwrite']!='boolean'){
			throw 'GetUserPermissionResp.canwrite must be boolean'
		}
		obj['canwrite']=jsonobj['canwrite']
	}
	//admin
	if(jsonobj['admin']!=null&&jsonobj['admin']!=undefined){
		if(typeof jsonobj['admin']!='boolean'){
			throw 'GetUserPermissionResp.admin must be boolean'
		}
		obj['admin']=jsonobj['admin']
	}
	return obj
}
export interface ListProjectNodeReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null|undefined;
}
function ListProjectNodeReqToJson(msg: ListProjectNodeReq): string{
	let s: string="{"
	//project_id
	if(msg.project_id==null||msg.project_id==undefined){
		s+='"project_id":null,'
	}else if(msg.project_id.length==0){
		s+='"project_id":[],'
	}else{
		s+='"project_id":['
		for(let element of msg.project_id){
			if(element==null||element==undefined||!Number.isInteger(element)){
				throw 'element in ListProjectNodeReq.project_id must be integer'
			}
			if(element>4294967295||element<0){
				throw 'element in ListProjectNodeReq.project_id overflow'
			}
			s+=element+','
		}
		s=s.substr(0,s.length-1)+'],'
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface ListProjectNodeResp{
	//this will only return the node name,node data and children
	//other node's info will not return
	node: NodeInfo|null|undefined;
}
function JsonToListProjectNodeResp(jsonobj: { [k:string]:any }): ListProjectNodeResp{
	let obj: ListProjectNodeResp={
		node:null,
	}
	//node
	if(jsonobj['node']!=null&&jsonobj['node']!=undefined){
		if(typeof jsonobj['node']!='object'){
			throw 'ListProjectNodeResp.node must be NodeInfo'
		}
		obj['node']=JsonToNodeInfo(jsonobj['node'])
	}
	return obj
}
export interface ListRoleNodeReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null|undefined;
	role_name: string;
}
function ListRoleNodeReqToJson(msg: ListRoleNodeReq): string{
	let s: string="{"
	//project_id
	if(msg.project_id==null||msg.project_id==undefined){
		s+='"project_id":null,'
	}else if(msg.project_id.length==0){
		s+='"project_id":[],'
	}else{
		s+='"project_id":['
		for(let element of msg.project_id){
			if(element==null||element==undefined||!Number.isInteger(element)){
				throw 'element in ListRoleNodeReq.project_id must be integer'
			}
			if(element>4294967295||element<0){
				throw 'element in ListRoleNodeReq.project_id overflow'
			}
			s+=element+','
		}
		s=s.substr(0,s.length-1)+'],'
	}
	//role_name
	if(msg.role_name==null||msg.role_name==undefined){
		throw 'ListRoleNodeReq.role_name must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.role_name)
		s+='"role_name":'+vv+','
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface ListRoleNodeResp{
	node: NodeInfo|null|undefined;
}
function JsonToListRoleNodeResp(jsonobj: { [k:string]:any }): ListRoleNodeResp{
	let obj: ListRoleNodeResp={
		node:null,
	}
	//node
	if(jsonobj['node']!=null&&jsonobj['node']!=undefined){
		if(typeof jsonobj['node']!='object'){
			throw 'ListRoleNodeResp.node must be NodeInfo'
		}
		obj['node']=JsonToNodeInfo(jsonobj['node'])
	}
	return obj
}
export interface ListUserNodeReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null|undefined;
	user_id: string;//if this is empty means return self's
	need_user_role_node: boolean;//false - only return user's base node,true - return user's base node and user's roles' node
}
function ListUserNodeReqToJson(msg: ListUserNodeReq): string{
	let s: string="{"
	//project_id
	if(msg.project_id==null||msg.project_id==undefined){
		s+='"project_id":null,'
	}else if(msg.project_id.length==0){
		s+='"project_id":[],'
	}else{
		s+='"project_id":['
		for(let element of msg.project_id){
			if(element==null||element==undefined||!Number.isInteger(element)){
				throw 'element in ListUserNodeReq.project_id must be integer'
			}
			if(element>4294967295||element<0){
				throw 'element in ListUserNodeReq.project_id overflow'
			}
			s+=element+','
		}
		s=s.substr(0,s.length-1)+'],'
	}
	//user_id
	if(msg.user_id==null||msg.user_id==undefined){
		throw 'ListUserNodeReq.user_id must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.user_id)
		s+='"user_id":'+vv+','
	}
	//need_user_role_node
	if(msg.need_user_role_node==null||msg.need_user_role_node==undefined){
		throw 'ListUserNodeReq.need_user_role_node must be boolean'
	}else{
		s+='"need_user_role_node":'+msg.need_user_role_node+','
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface ListUserNodeResp{
	node: NodeInfo|null|undefined;
}
function JsonToListUserNodeResp(jsonobj: { [k:string]:any }): ListUserNodeResp{
	let obj: ListUserNodeResp={
		node:null,
	}
	//node
	if(jsonobj['node']!=null&&jsonobj['node']!=undefined){
		if(typeof jsonobj['node']!='object'){
			throw 'ListUserNodeResp.node must be NodeInfo'
		}
		obj['node']=JsonToNodeInfo(jsonobj['node'])
	}
	return obj
}
export interface MoveNodeReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	node_id: Array<number>|null|undefined;
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	pnode_id: Array<number>|null|undefined;
}
function MoveNodeReqToJson(msg: MoveNodeReq): string{
	let s: string="{"
	//node_id
	if(msg.node_id==null||msg.node_id==undefined){
		s+='"node_id":null,'
	}else if(msg.node_id.length==0){
		s+='"node_id":[],'
	}else{
		s+='"node_id":['
		for(let element of msg.node_id){
			if(element==null||element==undefined||!Number.isInteger(element)){
				throw 'element in MoveNodeReq.node_id must be integer'
			}
			if(element>4294967295||element<0){
				throw 'element in MoveNodeReq.node_id overflow'
			}
			s+=element+','
		}
		s=s.substr(0,s.length-1)+'],'
	}
	//pnode_id
	if(msg.pnode_id==null||msg.pnode_id==undefined){
		s+='"pnode_id":null,'
	}else if(msg.pnode_id.length==0){
		s+='"pnode_id":[],'
	}else{
		s+='"pnode_id":['
		for(let element of msg.pnode_id){
			if(element==null||element==undefined||!Number.isInteger(element)){
				throw 'element in MoveNodeReq.pnode_id must be integer'
			}
			if(element>4294967295||element<0){
				throw 'element in MoveNodeReq.pnode_id overflow'
			}
			s+=element+','
		}
		s=s.substr(0,s.length-1)+'],'
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface MoveNodeResp{
}
function JsonToMoveNodeResp(_jsonobj: { [k:string]:any }): MoveNodeResp{
	let obj: MoveNodeResp={
	}
	return obj
}
export interface NodeInfo{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	node_id: Array<number>|null|undefined;
	node_name: string;
	node_data: string;
	canread: boolean;
	canwrite: boolean;
	admin: boolean;
	children: Array<NodeInfo|null|undefined>|null|undefined;
}
function JsonToNodeInfo(jsonobj: { [k:string]:any }): NodeInfo{
	let obj: NodeInfo={
		node_id:null,
		node_name:'',
		node_data:'',
		canread:false,
		canwrite:false,
		admin:false,
		children:null,
	}
	//node_id
	if(jsonobj['node_id']!=null&&jsonobj['node_id']!=undefined){
		if(!(jsonobj['node_id'] instanceof Array)){
			throw 'NodeInfo.node_id must be Array<number>|null|undefined'
		}
		for(let element of jsonobj['node_id']){
			if(typeof element!='number'||!Number.isInteger(element)){
				throw 'element in NodeInfo.node_id must be integer'
			}else if(element>4294967295||element<0){
				throw 'element in NodeInfo.node_id overflow'
			}
			if(obj['node_id']==null){
				obj['node_id']=new Array<number>
			}
			obj['node_id'].push(element)
		}
	}
	//node_name
	if(jsonobj['node_name']!=null&&jsonobj['node_name']!=undefined){
		if(typeof jsonobj['node_name']!='string'){
			throw 'NodeInfo.node_name must be string'
		}
		obj['node_name']=jsonobj['node_name']
	}
	//node_data
	if(jsonobj['node_data']!=null&&jsonobj['node_data']!=undefined){
		if(typeof jsonobj['node_data']!='string'){
			throw 'NodeInfo.node_data must be string'
		}
		obj['node_data']=jsonobj['node_data']
	}
	//canread
	if(jsonobj['canread']!=null&&jsonobj['canread']!=undefined){
		if(typeof jsonobj['canread']!='boolean'){
			throw 'NodeInfo.canread must be boolean'
		}
		obj['canread']=jsonobj['canread']
	}
	//canwrite
	if(jsonobj['canwrite']!=null&&jsonobj['canwrite']!=undefined){
		if(typeof jsonobj['canwrite']!='boolean'){
			throw 'NodeInfo.canwrite must be boolean'
		}
		obj['canwrite']=jsonobj['canwrite']
	}
	//admin
	if(jsonobj['admin']!=null&&jsonobj['admin']!=undefined){
		if(typeof jsonobj['admin']!='boolean'){
			throw 'NodeInfo.admin must be boolean'
		}
		obj['admin']=jsonobj['admin']
	}
	//children
	if(jsonobj['children']!=null&&jsonobj['children']!=undefined){
		if(!(jsonobj['children'] instanceof Array)){
			throw 'NodeInfo.children must be Array<NodeInfo>|null|undefined'
		}
		for(let element of jsonobj['children']){
			if(typeof element!='object'){
				throw 'element in NodeInfo.children must be NodeInfo'
			}
			if(obj['children']==null){
				obj['children']=new Array<NodeInfo>
			}
			obj['children'].push(JsonToNodeInfo(element))
		}
	}
	return obj
}
export interface UpdateNodeReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	node_id: Array<number>|null|undefined;
	new_node_name: string;//if didn't change,set this with the old value
	new_node_data: string;//if didn't change,set this with the old value
}
function UpdateNodeReqToJson(msg: UpdateNodeReq): string{
	let s: string="{"
	//node_id
	if(msg.node_id==null||msg.node_id==undefined){
		s+='"node_id":null,'
	}else if(msg.node_id.length==0){
		s+='"node_id":[],'
	}else{
		s+='"node_id":['
		for(let element of msg.node_id){
			if(element==null||element==undefined||!Number.isInteger(element)){
				throw 'element in UpdateNodeReq.node_id must be integer'
			}
			if(element>4294967295||element<0){
				throw 'element in UpdateNodeReq.node_id overflow'
			}
			s+=element+','
		}
		s=s.substr(0,s.length-1)+'],'
	}
	//new_node_name
	if(msg.new_node_name==null||msg.new_node_name==undefined){
		throw 'UpdateNodeReq.new_node_name must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.new_node_name)
		s+='"new_node_name":'+vv+','
	}
	//new_node_data
	if(msg.new_node_data==null||msg.new_node_data==undefined){
		throw 'UpdateNodeReq.new_node_data must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.new_node_data)
		s+='"new_node_data":'+vv+','
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface UpdateNodeResp{
}
function JsonToUpdateNodeResp(_jsonobj: { [k:string]:any }): UpdateNodeResp{
	let obj: UpdateNodeResp={
	}
	return obj
}
export interface UpdateRolePermissionReq{
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null|undefined;
	role_name: string;
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	node_id: Array<number>|null|undefined;
	//if admin is true,canread and canwrite will be ignore
	admin: boolean;
	//if admin is false,and canread is false too,means delete this user from this node
	//if admin is false,and canread is false and node_id's length is 1,means delete this user completely
	canread: boolean;
	//if canwrite is true,canread must be true too
	canwrite: boolean;
}
function UpdateRolePermissionReqToJson(msg: UpdateRolePermissionReq): string{
	let s: string="{"
	//project_id
	if(msg.project_id==null||msg.project_id==undefined){
		s+='"project_id":null,'
	}else if(msg.project_id.length==0){
		s+='"project_id":[],'
	}else{
		s+='"project_id":['
		for(let element of msg.project_id){
			if(element==null||element==undefined||!Number.isInteger(element)){
				throw 'element in UpdateRolePermissionReq.project_id must be integer'
			}
			if(element>4294967295||element<0){
				throw 'element in UpdateRolePermissionReq.project_id overflow'
			}
			s+=element+','
		}
		s=s.substr(0,s.length-1)+'],'
	}
	//role_name
	if(msg.role_name==null||msg.role_name==undefined){
		throw 'UpdateRolePermissionReq.role_name must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.role_name)
		s+='"role_name":'+vv+','
	}
	//node_id
	if(msg.node_id==null||msg.node_id==undefined){
		s+='"node_id":null,'
	}else if(msg.node_id.length==0){
		s+='"node_id":[],'
	}else{
		s+='"node_id":['
		for(let element of msg.node_id){
			if(element==null||element==undefined||!Number.isInteger(element)){
				throw 'element in UpdateRolePermissionReq.node_id must be integer'
			}
			if(element>4294967295||element<0){
				throw 'element in UpdateRolePermissionReq.node_id overflow'
			}
			s+=element+','
		}
		s=s.substr(0,s.length-1)+'],'
	}
	//admin
	if(msg.admin==null||msg.admin==undefined){
		throw 'UpdateRolePermissionReq.admin must be boolean'
	}else{
		s+='"admin":'+msg.admin+','
	}
	//canread
	if(msg.canread==null||msg.canread==undefined){
		throw 'UpdateRolePermissionReq.canread must be boolean'
	}else{
		s+='"canread":'+msg.canread+','
	}
	//canwrite
	if(msg.canwrite==null||msg.canwrite==undefined){
		throw 'UpdateRolePermissionReq.canwrite must be boolean'
	}else{
		s+='"canwrite":'+msg.canwrite+','
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface UpdateRolePermissionResp{
}
function JsonToUpdateRolePermissionResp(_jsonobj: { [k:string]:any }): UpdateRolePermissionResp{
	let obj: UpdateRolePermissionResp={
	}
	return obj
}
export interface UpdateUserPermissionReq{
	user_id: string;
	//first element must be 0
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	node_id: Array<number>|null|undefined;
	//if admin is true,canread and canwrite will be ignore
	admin: boolean;
	//if admin is false,and canread is false too,means delete this user from this node
	//if admin is false,and canread is false and node_id's length is 1,means delete this user completely
	canread: boolean;
	//if canwrite is true,canread must be true too
	canwrite: boolean;
}
function UpdateUserPermissionReqToJson(msg: UpdateUserPermissionReq): string{
	let s: string="{"
	//user_id
	if(msg.user_id==null||msg.user_id==undefined){
		throw 'UpdateUserPermissionReq.user_id must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.user_id)
		s+='"user_id":'+vv+','
	}
	//node_id
	if(msg.node_id==null||msg.node_id==undefined){
		s+='"node_id":null,'
	}else if(msg.node_id.length==0){
		s+='"node_id":[],'
	}else{
		s+='"node_id":['
		for(let element of msg.node_id){
			if(element==null||element==undefined||!Number.isInteger(element)){
				throw 'element in UpdateUserPermissionReq.node_id must be integer'
			}
			if(element>4294967295||element<0){
				throw 'element in UpdateUserPermissionReq.node_id overflow'
			}
			s+=element+','
		}
		s=s.substr(0,s.length-1)+'],'
	}
	//admin
	if(msg.admin==null||msg.admin==undefined){
		throw 'UpdateUserPermissionReq.admin must be boolean'
	}else{
		s+='"admin":'+msg.admin+','
	}
	//canread
	if(msg.canread==null||msg.canread==undefined){
		throw 'UpdateUserPermissionReq.canread must be boolean'
	}else{
		s+='"canread":'+msg.canread+','
	}
	//canwrite
	if(msg.canwrite==null||msg.canwrite==undefined){
		throw 'UpdateUserPermissionReq.canwrite must be boolean'
	}else{
		s+='"canwrite":'+msg.canwrite+','
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface UpdateUserPermissionResp{
}
function JsonToUpdateUserPermissionResp(_jsonobj: { [k:string]:any }): UpdateUserPermissionResp{
	let obj: UpdateUserPermissionResp={
	}
	return obj
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
//ToC means this is used for users
export class PermissionBrowserClientToC {
	constructor(host: string){
		if(host==null||host==undefined||host.length==0){
			throw "PermissionBrowserClientToC's host missing"
		}
		this.host=host
	}
	//timeout must be integer,timeout's unit is millisecond
	//don't set Content-Type in header
	get_user_permission(header: { [k: string]: string },req: GetUserPermissionReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: GetUserPermissionResp)=>void){
		if(!Number.isInteger(timeout)){
			errorf({code:-2,msg:'timeout must be integer'})
			return
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let body: string=''
		try{
			body=GetUserPermissionReqToJson(req)
		}catch(e){
			errorf({code:-2,msg:''+e})
			return
		}
		let config={
			url:_WebPathPermissionGetUserPermission,
			method: "post",
			baseURL: this.host,
			headers: header,
			data: body,
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			let obj:GetUserPermissionResp
			try{
				obj=JsonToGetUserPermissionResp(response.data)
			}catch(e){
				let err:Error={code:-1,msg:'response body decode failed'}
				errorf(err)
				return
			}
			try{
				successf(obj)
			}catch(e){
				let err:Error={code:-1,msg:'success callback run failed'}
				errorf(err)
			}
		})
		.catch(function(error){
			if(error.response==undefined){
				errorf({code:-2,msg:error.message})
				return
			}
			let respdata=error.response.data
			let err:Error={code:-1,msg:''}
			if(respdata.code==undefined||typeof respdata.code!='number'||!Number.isInteger(respdata.code)||respdata.msg==undefined||typeof respdata.msg!='string'){
				err.msg=respdata
			}else{
				err.code=respdata.code
				err.msg=respdata.msg
			}
			errorf(err)
		})
	}
	//timeout must be integer,timeout's unit is millisecond
	//don't set Content-Type in header
	update_user_permission(header: { [k: string]: string },req: UpdateUserPermissionReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: UpdateUserPermissionResp)=>void){
		if(!Number.isInteger(timeout)){
			errorf({code:-2,msg:'timeout must be integer'})
			return
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let body: string=''
		try{
			body=UpdateUserPermissionReqToJson(req)
		}catch(e){
			errorf({code:-2,msg:''+e})
			return
		}
		let config={
			url:_WebPathPermissionUpdateUserPermission,
			method: "post",
			baseURL: this.host,
			headers: header,
			data: body,
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			let obj:UpdateUserPermissionResp
			try{
				obj=JsonToUpdateUserPermissionResp(response.data)
			}catch(e){
				let err:Error={code:-1,msg:'response body decode failed'}
				errorf(err)
				return
			}
			try{
				successf(obj)
			}catch(e){
				let err:Error={code:-1,msg:'success callback run failed'}
				errorf(err)
			}
		})
		.catch(function(error){
			if(error.response==undefined){
				errorf({code:-2,msg:error.message})
				return
			}
			let respdata=error.response.data
			let err:Error={code:-1,msg:''}
			if(respdata.code==undefined||typeof respdata.code!='number'||!Number.isInteger(respdata.code)||respdata.msg==undefined||typeof respdata.msg!='string'){
				err.msg=respdata
			}else{
				err.code=respdata.code
				err.msg=respdata.msg
			}
			errorf(err)
		})
	}
	//timeout must be integer,timeout's unit is millisecond
	//don't set Content-Type in header
	update_role_permission(header: { [k: string]: string },req: UpdateRolePermissionReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: UpdateRolePermissionResp)=>void){
		if(!Number.isInteger(timeout)){
			errorf({code:-2,msg:'timeout must be integer'})
			return
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let body: string=''
		try{
			body=UpdateRolePermissionReqToJson(req)
		}catch(e){
			errorf({code:-2,msg:''+e})
			return
		}
		let config={
			url:_WebPathPermissionUpdateRolePermission,
			method: "post",
			baseURL: this.host,
			headers: header,
			data: body,
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			let obj:UpdateRolePermissionResp
			try{
				obj=JsonToUpdateRolePermissionResp(response.data)
			}catch(e){
				let err:Error={code:-1,msg:'response body decode failed'}
				errorf(err)
				return
			}
			try{
				successf(obj)
			}catch(e){
				let err:Error={code:-1,msg:'success callback run failed'}
				errorf(err)
			}
		})
		.catch(function(error){
			if(error.response==undefined){
				errorf({code:-2,msg:error.message})
				return
			}
			let respdata=error.response.data
			let err:Error={code:-1,msg:''}
			if(respdata.code==undefined||typeof respdata.code!='number'||!Number.isInteger(respdata.code)||respdata.msg==undefined||typeof respdata.msg!='string'){
				err.msg=respdata
			}else{
				err.code=respdata.code
				err.msg=respdata.msg
			}
			errorf(err)
		})
	}
	//timeout must be integer,timeout's unit is millisecond
	//don't set Content-Type in header
	add_node(header: { [k: string]: string },req: AddNodeReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: AddNodeResp)=>void){
		if(!Number.isInteger(timeout)){
			errorf({code:-2,msg:'timeout must be integer'})
			return
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let body: string=''
		try{
			body=AddNodeReqToJson(req)
		}catch(e){
			errorf({code:-2,msg:''+e})
			return
		}
		let config={
			url:_WebPathPermissionAddNode,
			method: "post",
			baseURL: this.host,
			headers: header,
			data: body,
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			let obj:AddNodeResp
			try{
				obj=JsonToAddNodeResp(response.data)
			}catch(e){
				let err:Error={code:-1,msg:'response body decode failed'}
				errorf(err)
				return
			}
			try{
				successf(obj)
			}catch(e){
				let err:Error={code:-1,msg:'success callback run failed'}
				errorf(err)
			}
		})
		.catch(function(error){
			if(error.response==undefined){
				errorf({code:-2,msg:error.message})
				return
			}
			let respdata=error.response.data
			let err:Error={code:-1,msg:''}
			if(respdata.code==undefined||typeof respdata.code!='number'||!Number.isInteger(respdata.code)||respdata.msg==undefined||typeof respdata.msg!='string'){
				err.msg=respdata
			}else{
				err.code=respdata.code
				err.msg=respdata.msg
			}
			errorf(err)
		})
	}
	//timeout must be integer,timeout's unit is millisecond
	//don't set Content-Type in header
	update_node(header: { [k: string]: string },req: UpdateNodeReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: UpdateNodeResp)=>void){
		if(!Number.isInteger(timeout)){
			errorf({code:-2,msg:'timeout must be integer'})
			return
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let body: string=''
		try{
			body=UpdateNodeReqToJson(req)
		}catch(e){
			errorf({code:-2,msg:''+e})
			return
		}
		let config={
			url:_WebPathPermissionUpdateNode,
			method: "post",
			baseURL: this.host,
			headers: header,
			data: body,
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			let obj:UpdateNodeResp
			try{
				obj=JsonToUpdateNodeResp(response.data)
			}catch(e){
				let err:Error={code:-1,msg:'response body decode failed'}
				errorf(err)
				return
			}
			try{
				successf(obj)
			}catch(e){
				let err:Error={code:-1,msg:'success callback run failed'}
				errorf(err)
			}
		})
		.catch(function(error){
			if(error.response==undefined){
				errorf({code:-2,msg:error.message})
				return
			}
			let respdata=error.response.data
			let err:Error={code:-1,msg:''}
			if(respdata.code==undefined||typeof respdata.code!='number'||!Number.isInteger(respdata.code)||respdata.msg==undefined||typeof respdata.msg!='string'){
				err.msg=respdata
			}else{
				err.code=respdata.code
				err.msg=respdata.msg
			}
			errorf(err)
		})
	}
	//timeout must be integer,timeout's unit is millisecond
	//don't set Content-Type in header
	move_node(header: { [k: string]: string },req: MoveNodeReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: MoveNodeResp)=>void){
		if(!Number.isInteger(timeout)){
			errorf({code:-2,msg:'timeout must be integer'})
			return
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let body: string=''
		try{
			body=MoveNodeReqToJson(req)
		}catch(e){
			errorf({code:-2,msg:''+e})
			return
		}
		let config={
			url:_WebPathPermissionMoveNode,
			method: "post",
			baseURL: this.host,
			headers: header,
			data: body,
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			let obj:MoveNodeResp
			try{
				obj=JsonToMoveNodeResp(response.data)
			}catch(e){
				let err:Error={code:-1,msg:'response body decode failed'}
				errorf(err)
				return
			}
			try{
				successf(obj)
			}catch(e){
				let err:Error={code:-1,msg:'success callback run failed'}
				errorf(err)
			}
		})
		.catch(function(error){
			if(error.response==undefined){
				errorf({code:-2,msg:error.message})
				return
			}
			let respdata=error.response.data
			let err:Error={code:-1,msg:''}
			if(respdata.code==undefined||typeof respdata.code!='number'||!Number.isInteger(respdata.code)||respdata.msg==undefined||typeof respdata.msg!='string'){
				err.msg=respdata
			}else{
				err.code=respdata.code
				err.msg=respdata.msg
			}
			errorf(err)
		})
	}
	//timeout must be integer,timeout's unit is millisecond
	//don't set Content-Type in header
	del_node(header: { [k: string]: string },req: DelNodeReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: DelNodeResp)=>void){
		if(!Number.isInteger(timeout)){
			errorf({code:-2,msg:'timeout must be integer'})
			return
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let body: string=''
		try{
			body=DelNodeReqToJson(req)
		}catch(e){
			errorf({code:-2,msg:''+e})
			return
		}
		let config={
			url:_WebPathPermissionDelNode,
			method: "post",
			baseURL: this.host,
			headers: header,
			data: body,
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			let obj:DelNodeResp
			try{
				obj=JsonToDelNodeResp(response.data)
			}catch(e){
				let err:Error={code:-1,msg:'response body decode failed'}
				errorf(err)
				return
			}
			try{
				successf(obj)
			}catch(e){
				let err:Error={code:-1,msg:'success callback run failed'}
				errorf(err)
			}
		})
		.catch(function(error){
			if(error.response==undefined){
				errorf({code:-2,msg:error.message})
				return
			}
			let respdata=error.response.data
			let err:Error={code:-1,msg:''}
			if(respdata.code==undefined||typeof respdata.code!='number'||!Number.isInteger(respdata.code)||respdata.msg==undefined||typeof respdata.msg!='string'){
				err.msg=respdata
			}else{
				err.code=respdata.code
				err.msg=respdata.msg
			}
			errorf(err)
		})
	}
	//timeout must be integer,timeout's unit is millisecond
	//don't set Content-Type in header
	list_user_node(header: { [k: string]: string },req: ListUserNodeReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: ListUserNodeResp)=>void){
		if(!Number.isInteger(timeout)){
			errorf({code:-2,msg:'timeout must be integer'})
			return
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let body: string=''
		try{
			body=ListUserNodeReqToJson(req)
		}catch(e){
			errorf({code:-2,msg:''+e})
			return
		}
		let config={
			url:_WebPathPermissionListUserNode,
			method: "post",
			baseURL: this.host,
			headers: header,
			data: body,
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			let obj:ListUserNodeResp
			try{
				obj=JsonToListUserNodeResp(response.data)
			}catch(e){
				let err:Error={code:-1,msg:'response body decode failed'}
				errorf(err)
				return
			}
			try{
				successf(obj)
			}catch(e){
				let err:Error={code:-1,msg:'success callback run failed'}
				errorf(err)
			}
		})
		.catch(function(error){
			if(error.response==undefined){
				errorf({code:-2,msg:error.message})
				return
			}
			let respdata=error.response.data
			let err:Error={code:-1,msg:''}
			if(respdata.code==undefined||typeof respdata.code!='number'||!Number.isInteger(respdata.code)||respdata.msg==undefined||typeof respdata.msg!='string'){
				err.msg=respdata
			}else{
				err.code=respdata.code
				err.msg=respdata.msg
			}
			errorf(err)
		})
	}
	//timeout must be integer,timeout's unit is millisecond
	//don't set Content-Type in header
	list_role_node(header: { [k: string]: string },req: ListRoleNodeReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: ListRoleNodeResp)=>void){
		if(!Number.isInteger(timeout)){
			errorf({code:-2,msg:'timeout must be integer'})
			return
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let body: string=''
		try{
			body=ListRoleNodeReqToJson(req)
		}catch(e){
			errorf({code:-2,msg:''+e})
			return
		}
		let config={
			url:_WebPathPermissionListRoleNode,
			method: "post",
			baseURL: this.host,
			headers: header,
			data: body,
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			let obj:ListRoleNodeResp
			try{
				obj=JsonToListRoleNodeResp(response.data)
			}catch(e){
				let err:Error={code:-1,msg:'response body decode failed'}
				errorf(err)
				return
			}
			try{
				successf(obj)
			}catch(e){
				let err:Error={code:-1,msg:'success callback run failed'}
				errorf(err)
			}
		})
		.catch(function(error){
			if(error.response==undefined){
				errorf({code:-2,msg:error.message})
				return
			}
			let respdata=error.response.data
			let err:Error={code:-1,msg:''}
			if(respdata.code==undefined||typeof respdata.code!='number'||!Number.isInteger(respdata.code)||respdata.msg==undefined||typeof respdata.msg!='string'){
				err.msg=respdata
			}else{
				err.code=respdata.code
				err.msg=respdata.msg
			}
			errorf(err)
		})
	}
	//timeout must be integer,timeout's unit is millisecond
	//don't set Content-Type in header
	list_project_node(header: { [k: string]: string },req: ListProjectNodeReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: ListProjectNodeResp)=>void){
		if(!Number.isInteger(timeout)){
			errorf({code:-2,msg:'timeout must be integer'})
			return
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/json"
		let body: string=''
		try{
			body=ListProjectNodeReqToJson(req)
		}catch(e){
			errorf({code:-2,msg:''+e})
			return
		}
		let config={
			url:_WebPathPermissionListProjectNode,
			method: "post",
			baseURL: this.host,
			headers: header,
			data: body,
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			let obj:ListProjectNodeResp
			try{
				obj=JsonToListProjectNodeResp(response.data)
			}catch(e){
				let err:Error={code:-1,msg:'response body decode failed'}
				errorf(err)
				return
			}
			try{
				successf(obj)
			}catch(e){
				let err:Error={code:-1,msg:'success callback run failed'}
				errorf(err)
			}
		})
		.catch(function(error){
			if(error.response==undefined){
				errorf({code:-2,msg:error.message})
				return
			}
			let respdata=error.response.data
			let err:Error={code:-1,msg:''}
			if(respdata.code==undefined||typeof respdata.code!='number'||!Number.isInteger(respdata.code)||respdata.msg==undefined||typeof respdata.msg!='string'){
				err.msg=respdata
			}else{
				err.code=respdata.code
				err.msg=respdata.msg
			}
			errorf(err)
		})
	}
	private host: string
}
