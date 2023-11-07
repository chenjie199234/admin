// Code generated by protoc-gen-browser. DO NOT EDIT.
// version:
// 	protoc-gen-browser v0.0.92<br />
// 	protoc             v4.24.4<br />
// source: api/admin_initialize.proto<br />

import Axios from "axios";

export interface Error{
	code: number;
	msg: string;
}

export interface CreateProjectReq{
	project_name: string;
	project_data: string;
}
function CreateProjectReqToJson(msg: CreateProjectReq): string{
	let s: string="{"
	//project_name
	if(msg.project_name==null||msg.project_name==undefined){
		throw 'CreateProjectReq.project_name must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.project_name)
		s+='"project_name":'+vv+','
	}
	//project_data
	if(msg.project_data==null||msg.project_data==undefined){
		throw 'CreateProjectReq.project_data must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.project_data)
		s+='"project_data":'+vv+','
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface CreateProjectResp{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null|undefined;
}
function JsonToCreateProjectResp(jsonobj: { [k:string]:any }): CreateProjectResp{
	let obj: CreateProjectResp={
		project_id:null,
	}
	//project_id
	if(jsonobj['project_id']!=null&&jsonobj['project_id']!=undefined){
		if(!(jsonobj['project_id'] instanceof Array)){
			throw 'CreateProjectResp.project_id must be Array<number>|null|undefined'
		}
		for(let element of jsonobj['project_id']){
			if(typeof element!='number'||!Number.isInteger(element)){
				throw 'element in CreateProjectResp.project_id must be integer'
			}else if(element>4294967295||element<0){
				throw 'element in CreateProjectResp.project_id overflow'
			}
			if(obj['project_id']==null){
				obj['project_id']=new Array<number>
			}
			obj['project_id'].push(element)
		}
	}
	return obj
}
export interface DeleteProjectReq{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null|undefined;
}
function DeleteProjectReqToJson(msg: DeleteProjectReq): string{
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
				throw 'element in DeleteProjectReq.project_id must be integer'
			}
			if(element>4294967295||element<0){
				throw 'element in DeleteProjectReq.project_id overflow'
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
export interface DeleteProjectResp{
}
function JsonToDeleteProjectResp(_jsonobj: { [k:string]:any }): DeleteProjectResp{
	let obj: DeleteProjectResp={
	}
	return obj
}
export interface InitReq{
	password: string;
}
function InitReqToJson(msg: InitReq): string{
	let s: string="{"
	//password
	if(msg.password==null||msg.password==undefined){
		throw 'InitReq.password must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.password)
		s+='"password":'+vv+','
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface InitResp{
}
function JsonToInitResp(_jsonobj: { [k:string]:any }): InitResp{
	let obj: InitResp={
	}
	return obj
}
export interface InitStatusReq{
}
function InitStatusReqToJson(_msg: InitStatusReq): string{
	let s: string="{"
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface InitStatusResp{
	status: boolean;//true-already inited,false-not inited
}
function JsonToInitStatusResp(jsonobj: { [k:string]:any }): InitStatusResp{
	let obj: InitStatusResp={
		status:false,
	}
	//status
	if(jsonobj['status']!=null&&jsonobj['status']!=undefined){
		if(typeof jsonobj['status']!='boolean'){
			throw 'InitStatusResp.status must be boolean'
		}
		obj['status']=jsonobj['status']
	}
	return obj
}
export interface ListProjectReq{
}
function ListProjectReqToJson(_msg: ListProjectReq): string{
	let s: string="{"
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface ListProjectResp{
	projects: Array<ProjectInfo|null|undefined>|null|undefined;
}
function JsonToListProjectResp(jsonobj: { [k:string]:any }): ListProjectResp{
	let obj: ListProjectResp={
		projects:null,
	}
	//projects
	if(jsonobj['projects']!=null&&jsonobj['projects']!=undefined){
		if(!(jsonobj['projects'] instanceof Array)){
			throw 'ListProjectResp.projects must be Array<ProjectInfo>|null|undefined'
		}
		for(let element of jsonobj['projects']){
			if(typeof element!='object'){
				throw 'element in ListProjectResp.projects must be ProjectInfo'
			}
			if(obj['projects']==null){
				obj['projects']=new Array<ProjectInfo>
			}
			obj['projects'].push(JsonToProjectInfo(element))
		}
	}
	return obj
}
export interface ProjectInfo{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null|undefined;
	project_name: string;
	project_data: string;
}
function JsonToProjectInfo(jsonobj: { [k:string]:any }): ProjectInfo{
	let obj: ProjectInfo={
		project_id:null,
		project_name:'',
		project_data:'',
	}
	//project_id
	if(jsonobj['project_id']!=null&&jsonobj['project_id']!=undefined){
		if(!(jsonobj['project_id'] instanceof Array)){
			throw 'ProjectInfo.project_id must be Array<number>|null|undefined'
		}
		for(let element of jsonobj['project_id']){
			if(typeof element!='number'||!Number.isInteger(element)){
				throw 'element in ProjectInfo.project_id must be integer'
			}else if(element>4294967295||element<0){
				throw 'element in ProjectInfo.project_id overflow'
			}
			if(obj['project_id']==null){
				obj['project_id']=new Array<number>
			}
			obj['project_id'].push(element)
		}
	}
	//project_name
	if(jsonobj['project_name']!=null&&jsonobj['project_name']!=undefined){
		if(typeof jsonobj['project_name']!='string'){
			throw 'ProjectInfo.project_name must be string'
		}
		obj['project_name']=jsonobj['project_name']
	}
	//project_data
	if(jsonobj['project_data']!=null&&jsonobj['project_data']!=undefined){
		if(typeof jsonobj['project_data']!='string'){
			throw 'ProjectInfo.project_data must be string'
		}
		obj['project_data']=jsonobj['project_data']
	}
	return obj
}
export interface RootLoginReq{
	password: string;
}
function RootLoginReqToJson(msg: RootLoginReq): string{
	let s: string="{"
	//password
	if(msg.password==null||msg.password==undefined){
		throw 'RootLoginReq.password must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.password)
		s+='"password":'+vv+','
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface RootLoginResp{
	token: string;
}
function JsonToRootLoginResp(jsonobj: { [k:string]:any }): RootLoginResp{
	let obj: RootLoginResp={
		token:'',
	}
	//token
	if(jsonobj['token']!=null&&jsonobj['token']!=undefined){
		if(typeof jsonobj['token']!='string'){
			throw 'RootLoginResp.token must be string'
		}
		obj['token']=jsonobj['token']
	}
	return obj
}
export interface UpdateProjectReq{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null|undefined;
	new_project_name: string;//if didn't change,set this with the old value
	new_project_data: string;//if didn't change,set this with the old value
}
function UpdateProjectReqToJson(msg: UpdateProjectReq): string{
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
				throw 'element in UpdateProjectReq.project_id must be integer'
			}
			if(element>4294967295||element<0){
				throw 'element in UpdateProjectReq.project_id overflow'
			}
			s+=element+','
		}
		s=s.substr(0,s.length-1)+'],'
	}
	//new_project_name
	if(msg.new_project_name==null||msg.new_project_name==undefined){
		throw 'UpdateProjectReq.new_project_name must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.new_project_name)
		s+='"new_project_name":'+vv+','
	}
	//new_project_data
	if(msg.new_project_data==null||msg.new_project_data==undefined){
		throw 'UpdateProjectReq.new_project_data must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.new_project_data)
		s+='"new_project_data":'+vv+','
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface UpdateProjectResp{
}
function JsonToUpdateProjectResp(_jsonobj: { [k:string]:any }): UpdateProjectResp{
	let obj: UpdateProjectResp={
	}
	return obj
}
export interface UpdateRootPasswordReq{
	old_password: string;
	new_password: string;
}
function UpdateRootPasswordReqToJson(msg: UpdateRootPasswordReq): string{
	let s: string="{"
	//old_password
	if(msg.old_password==null||msg.old_password==undefined){
		throw 'UpdateRootPasswordReq.old_password must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.old_password)
		s+='"old_password":'+vv+','
	}
	//new_password
	if(msg.new_password==null||msg.new_password==undefined){
		throw 'UpdateRootPasswordReq.new_password must be string'
	}else{
		//transfer the json escape
		let vv=JSON.stringify(msg.new_password)
		s+='"new_password":'+vv+','
	}
	if(s.length==1){
		s+="}"
	}else{
		s=s.substr(0,s.length-1)+'}'
	}
	return s
}
export interface UpdateRootPasswordResp{
}
function JsonToUpdateRootPasswordResp(_jsonobj: { [k:string]:any }): UpdateRootPasswordResp{
	let obj: UpdateRootPasswordResp={
	}
	return obj
}
const _WebPathInitializeInitStatus: string ="/admin.initialize/init_status";
const _WebPathInitializeInit: string ="/admin.initialize/init";
const _WebPathInitializeRootLogin: string ="/admin.initialize/root_login";
const _WebPathInitializeUpdateRootPassword: string ="/admin.initialize/update_root_password";
const _WebPathInitializeCreateProject: string ="/admin.initialize/create_project";
const _WebPathInitializeUpdateProject: string ="/admin.initialize/update_project";
const _WebPathInitializeListProject: string ="/admin.initialize/list_project";
const _WebPathInitializeDeleteProject: string ="/admin.initialize/delete_project";
//ToC means this is used for users
export class InitializeBrowserClientToC {
	constructor(host: string){
		if(host==null||host==undefined||host.length==0){
			throw "InitializeBrowserClientToC's host missing"
		}
		this.host=host
	}
	//timeout must be integer,timeout's unit is millisecond
	//don't set Content-Type in header
	init_status(header: { [k: string]: string },req: InitStatusReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: InitStatusResp)=>void){
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
			body=InitStatusReqToJson(req)
		}catch(e){
			errorf({code:-2,msg:''+e})
			return
		}
		let config={
			url:_WebPathInitializeInitStatus,
			method: "post",
			baseURL: this.host,
			headers: header,
			data: body,
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			let obj:InitStatusResp
			try{
				obj=JsonToInitStatusResp(response.data)
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
	init(header: { [k: string]: string },req: InitReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: InitResp)=>void){
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
			body=InitReqToJson(req)
		}catch(e){
			errorf({code:-2,msg:''+e})
			return
		}
		let config={
			url:_WebPathInitializeInit,
			method: "post",
			baseURL: this.host,
			headers: header,
			data: body,
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			let obj:InitResp
			try{
				obj=JsonToInitResp(response.data)
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
	root_login(header: { [k: string]: string },req: RootLoginReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: RootLoginResp)=>void){
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
			body=RootLoginReqToJson(req)
		}catch(e){
			errorf({code:-2,msg:''+e})
			return
		}
		let config={
			url:_WebPathInitializeRootLogin,
			method: "post",
			baseURL: this.host,
			headers: header,
			data: body,
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			let obj:RootLoginResp
			try{
				obj=JsonToRootLoginResp(response.data)
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
	update_root_password(header: { [k: string]: string },req: UpdateRootPasswordReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: UpdateRootPasswordResp)=>void){
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
			body=UpdateRootPasswordReqToJson(req)
		}catch(e){
			errorf({code:-2,msg:''+e})
			return
		}
		let config={
			url:_WebPathInitializeUpdateRootPassword,
			method: "post",
			baseURL: this.host,
			headers: header,
			data: body,
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			let obj:UpdateRootPasswordResp
			try{
				obj=JsonToUpdateRootPasswordResp(response.data)
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
	create_project(header: { [k: string]: string },req: CreateProjectReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: CreateProjectResp)=>void){
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
			body=CreateProjectReqToJson(req)
		}catch(e){
			errorf({code:-2,msg:''+e})
			return
		}
		let config={
			url:_WebPathInitializeCreateProject,
			method: "post",
			baseURL: this.host,
			headers: header,
			data: body,
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			let obj:CreateProjectResp
			try{
				obj=JsonToCreateProjectResp(response.data)
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
	update_project(header: { [k: string]: string },req: UpdateProjectReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: UpdateProjectResp)=>void){
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
			body=UpdateProjectReqToJson(req)
		}catch(e){
			errorf({code:-2,msg:''+e})
			return
		}
		let config={
			url:_WebPathInitializeUpdateProject,
			method: "post",
			baseURL: this.host,
			headers: header,
			data: body,
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			let obj:UpdateProjectResp
			try{
				obj=JsonToUpdateProjectResp(response.data)
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
	list_project(header: { [k: string]: string },req: ListProjectReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: ListProjectResp)=>void){
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
			body=ListProjectReqToJson(req)
		}catch(e){
			errorf({code:-2,msg:''+e})
			return
		}
		let config={
			url:_WebPathInitializeListProject,
			method: "post",
			baseURL: this.host,
			headers: header,
			data: body,
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			let obj:ListProjectResp
			try{
				obj=JsonToListProjectResp(response.data)
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
	delete_project(header: { [k: string]: string },req: DeleteProjectReq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: DeleteProjectResp)=>void){
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
			body=DeleteProjectReqToJson(req)
		}catch(e){
			errorf({code:-2,msg:''+e})
			return
		}
		let config={
			url:_WebPathInitializeDeleteProject,
			method: "post",
			baseURL: this.host,
			headers: header,
			data: body,
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			let obj:DeleteProjectResp
			try{
				obj=JsonToDeleteProjectResp(response.data)
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
