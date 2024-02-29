// Code generated by protoc-gen-browser. DO NOT EDIT.
// version:
// 	protoc-gen-browser v0.0.103<br />
// 	protoc             v4.25.3<br />
// source: api/admin_status.proto<br />

export interface LogicError{
	code: number;
	msg: string;
}

export class Pingreq{
	//Warning!!!Type is int64,be careful of sign(+,-)
	timestamp: bigint = BigInt(0)
	toFORM(){
		let query=new Array<string>()
		if(this.timestamp){
			let tmp="timestamp="+this.timestamp
			query.push(tmp)
		}
		return encodeURIComponent(query.join('&'))
	}
}
export class Pingresp{
	//Warning!!!Type is int64,be careful of sign(+,-)
	client_timestamp: bigint = BigInt(0)
	//Warning!!!Type is int64,be careful of sign(+,-)
	server_timestamp: bigint = BigInt(0)
	//Warning!!!Type is uint64,be careful of sign(+)
	total_mem: bigint = BigInt(0)
	//Warning!!!Type is uint64,be careful of sign(+)
	cur_mem_usage: bigint = BigInt(0)
	//Warning!!!Type is uint64,be careful of sign(+)
	max_mem_usage: bigint = BigInt(0)
	cpu_num: number = 0
	cur_cpu_usage: number = 0
	avg_cpu_usage: number = 0
	max_cpu_usage: number = 0
	host: string = ''
	ip: string = ''
	fromOBJ(obj:Object){
		if(obj["client_timestamp"]){
			this.client_timestamp=BigInt(obj["client_timestamp"])
		}
		if(obj["server_timestamp"]){
			this.server_timestamp=BigInt(obj["server_timestamp"])
		}
		if(obj["total_mem"]){
			this.total_mem=BigInt(obj["total_mem"])
		}
		if(obj["cur_mem_usage"]){
			this.cur_mem_usage=BigInt(obj["cur_mem_usage"])
		}
		if(obj["max_mem_usage"]){
			this.max_mem_usage=BigInt(obj["max_mem_usage"])
		}
		if(obj["cpu_num"]){
			this.cpu_num=obj["cpu_num"]
		}
		if(obj["cur_cpu_usage"]){
			this.cur_cpu_usage=obj["cur_cpu_usage"]
		}
		if(obj["avg_cpu_usage"]){
			this.avg_cpu_usage=obj["avg_cpu_usage"]
		}
		if(obj["max_cpu_usage"]){
			this.max_cpu_usage=obj["max_cpu_usage"]
		}
		if(obj["host"]){
			this.host=obj["host"]
		}
		if(obj["ip"]){
			this.ip=obj["ip"]
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
const _WebPathStatusPing: string ="/admin.status/ping";
//ToC means this is for users
export class StatusBrowserClientToC {
	constructor(host: string){
		if(!host || host.length==0){
			throw "StatusBrowserClientToC's host missing"
		}
		this.host=host
	}
	//timeout's unit is millisecond,it will be used when > 0
	ping(header: Object,req: Pingreq,timeout: number,error: (arg: LogicError)=>void,success: (arg: Pingresp)=>void){
		if(!header){
			header={}
		}
		header["Content-Type"] = "application/x-www-form-urlencoded"
		call(timeout,this.host+_WebPathStatusPing+"?"+req.toFORM(),{method:"GET",headers:header},error,function(arg: Object){
			let r=new Pingresp()
			r.fromOBJ(arg)
			success(r)
		})
	}
	private host: string
}
