// Code generated by protoc-gen-browser. DO NOT EDIT.
// version:
// 	protoc-gen-browser v0.0.120<br />
// 	protoc             v5.28.0<br />
// source: api/admin_initialize.proto<br />

export interface LogicError{
	code: number;
	msg: string;
}

export class CreateProjectReq{
	project_name: string = ''
	project_data: string = ''
	toJSON(){
		let tmp = {}
		if(this.project_name){
			tmp["project_name"]=this.project_name
		}
		if(this.project_data){
			tmp["project_data"]=this.project_data
		}
		return tmp
	}
}
export class CreateProjectResp{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	fromOBJ(obj:Object){
		if(obj["project_id"] && obj["project_id"].length>0){
			this.project_id=obj["project_id"]
		}
	}
}
export class DeleteProjectReq{
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
export class DeleteProjectResp{
	fromOBJ(_obj:Object){
	}
}
export class InitReq{
	password: string = ''
	toJSON(){
		let tmp = {}
		if(this.password){
			tmp["password"]=this.password
		}
		return tmp
	}
}
export class InitResp{
	fromOBJ(_obj:Object){
	}
}
export class InitStatusReq{
	toJSON(){
		let tmp = {}
		return tmp
	}
}
export class InitStatusResp{
	status: boolean = false//true-already inited,false-not inited
	fromOBJ(obj:Object){
		if(obj["status"]){
			this.status=obj["status"]
		}
	}
}
export class ListProjectReq{
	toJSON(){
		let tmp = {}
		return tmp
	}
}
export class ListProjectResp{
	projects: Array<ProjectInfo|null>|null = null
	fromOBJ(obj:Object){
		if(obj["projects"] && obj["projects"].length>0){
			this.projects=new Array<ProjectInfo|null>()
			for(let value of obj["projects"]){
				if(value){
					let tmp=new ProjectInfo()
					tmp.fromOBJ(value)
					this.projects.push(tmp)
				}else{
					this.projects.push(null)
				}
			}
		}
	}
}
export class ProjectInfo{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	project_name: string = ''
	project_data: string = ''
	fromOBJ(obj:Object){
		if(obj["project_id"] && obj["project_id"].length>0){
			this.project_id=obj["project_id"]
		}
		if(obj["project_name"]){
			this.project_name=obj["project_name"]
		}
		if(obj["project_data"]){
			this.project_data=obj["project_data"]
		}
	}
}
export class RootLoginReq{
	password: string = ''
	toJSON(){
		let tmp = {}
		if(this.password){
			tmp["password"]=this.password
		}
		return tmp
	}
}
export class RootLoginResp{
	token: string = ''
	fromOBJ(obj:Object){
		if(obj["token"]){
			this.token=obj["token"]
		}
	}
}
export class UpdateProjectReq{
	//Warning!!!Element type is uint32,be careful of sign(+) and overflow
	project_id: Array<number>|null = null
	new_project_name: string = ''//if didn't change,set this with the old value
	new_project_data: string = ''//if didn't change,set this with the old value
	toJSON(){
		let tmp = {}
		if(this.project_id && this.project_id.length>0){
			tmp["project_id"]=this.project_id
		}
		if(this.new_project_name){
			tmp["new_project_name"]=this.new_project_name
		}
		if(this.new_project_data){
			tmp["new_project_data"]=this.new_project_data
		}
		return tmp
	}
}
export class UpdateProjectResp{
	fromOBJ(_obj:Object){
	}
}
export class UpdateRootPasswordReq{
	old_password: string = ''
	new_password: string = ''
	toJSON(){
		let tmp = {}
		if(this.old_password){
			tmp["old_password"]=this.old_password
		}
		if(this.new_password){
			tmp["new_password"]=this.new_password
		}
		return tmp
	}
}
export class UpdateRootPasswordResp{
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
const _WebPathInitializeInitStatus: string ="/admin.initialize/init_status";
const _WebPathInitializeInit: string ="/admin.initialize/init";
const _WebPathInitializeRootLogin: string ="/admin.initialize/root_login";
const _WebPathInitializeUpdateRootPassword: string ="/admin.initialize/update_root_password";
const _WebPathInitializeCreateProject: string ="/admin.initialize/create_project";
const _WebPathInitializeUpdateProject: string ="/admin.initialize/update_project";
const _WebPathInitializeListProject: string ="/admin.initialize/list_project";
const _WebPathInitializeDeleteProject: string ="/admin.initialize/delete_project";
export class InitializeBrowserClient {
	constructor(host: string){
		if(!host || host.length==0){
			throw "InitializeBrowserClient's host missing"
		}
		this.host=host
	}
	//timeout's unit is millisecond,it will be used when > 0
	init_status(header: Object,req: InitStatusReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: InitStatusResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathInitializeInitStatus,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new InitStatusResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	init(header: Object,req: InitReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: InitResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathInitializeInit,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new InitResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	root_login(header: Object,req: RootLoginReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: RootLoginResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathInitializeRootLogin,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new RootLoginResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	update_root_password(header: Object,req: UpdateRootPasswordReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: UpdateRootPasswordResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathInitializeUpdateRootPassword,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new UpdateRootPasswordResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	create_project(header: Object,req: CreateProjectReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: CreateProjectResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathInitializeCreateProject,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new CreateProjectResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	update_project(header: Object,req: UpdateProjectReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: UpdateProjectResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathInitializeUpdateProject,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new UpdateProjectResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	list_project(header: Object,req: ListProjectReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: ListProjectResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathInitializeListProject,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new ListProjectResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	//timeout's unit is millisecond,it will be used when > 0
	delete_project(header: Object,req: DeleteProjectReq,timeout: number,error: (arg: LogicError)=>void,success: (arg: DeleteProjectResp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/json"
		call(timeout,this.host+_WebPathInitializeDeleteProject,{method:"POST",headers:header,body:JSON.stringify(req)},error,function(arg: Object){
			let r=new DeleteProjectResp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	private host: string
}
