// Code generated by protoc-gen-browser. DO NOT EDIT.
// version:
// 	protoc-gen-browser v0.0.95<br />
// 	protoc             v4.25.1<br />
// source: api/admin_status.proto<br />

import Axios from "axios";
import Long from "long";

export interface Error{
	code: number;
	msg: string;
}

export interface Pingreq{
	//Warning!!!Type is int64,be careful of sign(+,-)
	timestamp: Long;
}
function PingreqToForm(msg: Pingreq): string{
	let s: string=""
	//timestamp
	if(msg.timestamp==null||msg.timestamp==undefined){
		throw 'Pingreq.timestamp must be integer'
	}else if(msg.timestamp.lessThan(Long.MIN_VALUE)||msg.timestamp.greaterThan(Long.MAX_VALUE)){
		throw 'Pingreq.timestamp overflow'
	}else{
		s+='timestamp='+msg.timestamp.toString()+'&'
	}
	if(s.length!=0){
		s=s.substr(0,s.length-1)
	}
	return s
}
export interface Pingresp{
	//Warning!!!Type is int64,be careful of sign(+,-)
	client_timestamp: Long;
	//Warning!!!Type is int64,be careful of sign(+,-)
	server_timestamp: Long;
	//Warning!!!Type is uint64,be careful of sign(+)
	total_mem: Long;
	//Warning!!!Type is uint64,be careful of sign(+)
	cur_mem_usage: Long;
	//Warning!!!Type is uint64,be careful of sign(+)
	max_mem_usage: Long;
	cpu_num: number;
	cur_cpu_usage: number;
	avg_cpu_usage: number;
	max_cpu_usage: number;
	host: string;
	ip: string;
}
function JsonToPingresp(jsonobj: { [k:string]:any }): Pingresp{
	let obj: Pingresp={
		client_timestamp:Long.ZERO,
		server_timestamp:Long.ZERO,
		total_mem:Long.ZERO,
		cur_mem_usage:Long.ZERO,
		max_mem_usage:Long.ZERO,
		cpu_num:0,
		cur_cpu_usage:0,
		avg_cpu_usage:0,
		max_cpu_usage:0,
		host:'',
		ip:'',
	}
	//client_timestamp
	if(jsonobj['client_timestamp']!=null&&jsonobj['client_timestamp']!=undefined){
		if(typeof jsonobj['client_timestamp']=='number'){
			if(!Number.isInteger(jsonobj['client_timestamp'])){
				throw 'Pingresp.client_timestamp must be integer'
			}
			let tmp: Long=Long.ZERO
			try{
				tmp=Long.fromNumber(jsonobj['client_timestamp'],false)
			}catch(e){
				throw 'Pingresp.client_timestamp must be integer'
			}
			obj['client_timestamp']=tmp
		}else if(typeof jsonobj['client_timestamp']=='string'){
			let tmp:Long=Long.ZERO
			try{
				tmp=Long.fromString(jsonobj['client_timestamp'],false)
			}catch(e){
				throw 'Pingresp.client_timestamp must be integer'
			}
			if(tmp.toString()!=jsonobj['client_timestamp']){
				throw 'Pingresp.client_timestamp overflow'
			}
			obj['client_timestamp']=tmp
		}else{
			throw 'Pingresp.client_timestamp must be integer'
		}
	}
	//server_timestamp
	if(jsonobj['server_timestamp']!=null&&jsonobj['server_timestamp']!=undefined){
		if(typeof jsonobj['server_timestamp']=='number'){
			if(!Number.isInteger(jsonobj['server_timestamp'])){
				throw 'Pingresp.server_timestamp must be integer'
			}
			let tmp: Long=Long.ZERO
			try{
				tmp=Long.fromNumber(jsonobj['server_timestamp'],false)
			}catch(e){
				throw 'Pingresp.server_timestamp must be integer'
			}
			obj['server_timestamp']=tmp
		}else if(typeof jsonobj['server_timestamp']=='string'){
			let tmp:Long=Long.ZERO
			try{
				tmp=Long.fromString(jsonobj['server_timestamp'],false)
			}catch(e){
				throw 'Pingresp.server_timestamp must be integer'
			}
			if(tmp.toString()!=jsonobj['server_timestamp']){
				throw 'Pingresp.server_timestamp overflow'
			}
			obj['server_timestamp']=tmp
		}else{
			throw 'Pingresp.server_timestamp must be integer'
		}
	}
	//total_mem
	if(jsonobj['total_mem']!=null&&jsonobj['total_mem']!=undefined){
		if(typeof jsonobj['total_mem']=='number'){
			if(!Number.isInteger(jsonobj['total_mem'])){
				throw 'Pingresp.total_mem must be integer'
			}
			if(jsonobj['total_mem']<0){
				throw 'Pingresp.total_mem overflow'
			}
			let tmp: Long=Long.ZERO
			try{
				tmp=Long.fromNumber(jsonobj['total_mem'],true)
			}catch(e){
				throw 'Pingresp.total_mem must be integer'
			}
			obj['total_mem']=tmp
		}else if(typeof jsonobj['total_mem']=='string'){
			let tmp:Long=Long.ZERO
			try{
				tmp=Long.fromString(jsonobj['total_mem'],true)
			}catch(e){
				throw 'Pingresp.total_mem must be integer'
			}
			if(tmp.toString()!=jsonobj['total_mem']){
				throw 'Pingresp.total_mem overflow'
			}
			obj['total_mem']=tmp
		}else{
			throw 'format wrong!Pingresp.total_mem must be integer'
		}
	}
	//cur_mem_usage
	if(jsonobj['cur_mem_usage']!=null&&jsonobj['cur_mem_usage']!=undefined){
		if(typeof jsonobj['cur_mem_usage']=='number'){
			if(!Number.isInteger(jsonobj['cur_mem_usage'])){
				throw 'Pingresp.cur_mem_usage must be integer'
			}
			if(jsonobj['cur_mem_usage']<0){
				throw 'Pingresp.cur_mem_usage overflow'
			}
			let tmp: Long=Long.ZERO
			try{
				tmp=Long.fromNumber(jsonobj['cur_mem_usage'],true)
			}catch(e){
				throw 'Pingresp.cur_mem_usage must be integer'
			}
			obj['cur_mem_usage']=tmp
		}else if(typeof jsonobj['cur_mem_usage']=='string'){
			let tmp:Long=Long.ZERO
			try{
				tmp=Long.fromString(jsonobj['cur_mem_usage'],true)
			}catch(e){
				throw 'Pingresp.cur_mem_usage must be integer'
			}
			if(tmp.toString()!=jsonobj['cur_mem_usage']){
				throw 'Pingresp.cur_mem_usage overflow'
			}
			obj['cur_mem_usage']=tmp
		}else{
			throw 'format wrong!Pingresp.cur_mem_usage must be integer'
		}
	}
	//max_mem_usage
	if(jsonobj['max_mem_usage']!=null&&jsonobj['max_mem_usage']!=undefined){
		if(typeof jsonobj['max_mem_usage']=='number'){
			if(!Number.isInteger(jsonobj['max_mem_usage'])){
				throw 'Pingresp.max_mem_usage must be integer'
			}
			if(jsonobj['max_mem_usage']<0){
				throw 'Pingresp.max_mem_usage overflow'
			}
			let tmp: Long=Long.ZERO
			try{
				tmp=Long.fromNumber(jsonobj['max_mem_usage'],true)
			}catch(e){
				throw 'Pingresp.max_mem_usage must be integer'
			}
			obj['max_mem_usage']=tmp
		}else if(typeof jsonobj['max_mem_usage']=='string'){
			let tmp:Long=Long.ZERO
			try{
				tmp=Long.fromString(jsonobj['max_mem_usage'],true)
			}catch(e){
				throw 'Pingresp.max_mem_usage must be integer'
			}
			if(tmp.toString()!=jsonobj['max_mem_usage']){
				throw 'Pingresp.max_mem_usage overflow'
			}
			obj['max_mem_usage']=tmp
		}else{
			throw 'format wrong!Pingresp.max_mem_usage must be integer'
		}
	}
	//cpu_num
	if(jsonobj['cpu_num']!=null&&jsonobj['cpu_num']!=undefined){
		if(typeof jsonobj['cpu_num']!='number'){
			throw 'Pingresp.cpu_num must be number'
		}
		obj['cpu_num']=jsonobj['cpu_num']
	}
	//cur_cpu_usage
	if(jsonobj['cur_cpu_usage']!=null&&jsonobj['cur_cpu_usage']!=undefined){
		if(typeof jsonobj['cur_cpu_usage']!='number'){
			throw 'Pingresp.cur_cpu_usage must be number'
		}
		obj['cur_cpu_usage']=jsonobj['cur_cpu_usage']
	}
	//avg_cpu_usage
	if(jsonobj['avg_cpu_usage']!=null&&jsonobj['avg_cpu_usage']!=undefined){
		if(typeof jsonobj['avg_cpu_usage']!='number'){
			throw 'Pingresp.avg_cpu_usage must be number'
		}
		obj['avg_cpu_usage']=jsonobj['avg_cpu_usage']
	}
	//max_cpu_usage
	if(jsonobj['max_cpu_usage']!=null&&jsonobj['max_cpu_usage']!=undefined){
		if(typeof jsonobj['max_cpu_usage']!='number'){
			throw 'Pingresp.max_cpu_usage must be number'
		}
		obj['max_cpu_usage']=jsonobj['max_cpu_usage']
	}
	//host
	if(jsonobj['host']!=null&&jsonobj['host']!=undefined){
		if(typeof jsonobj['host']!='string'){
			throw 'Pingresp.host must be string'
		}
		obj['host']=jsonobj['host']
	}
	//ip
	if(jsonobj['ip']!=null&&jsonobj['ip']!=undefined){
		if(typeof jsonobj['ip']!='string'){
			throw 'Pingresp.ip must be string'
		}
		obj['ip']=jsonobj['ip']
	}
	return obj
}
const _WebPathStatusPing: string ="/admin.status/ping";
//ToC means this is used for users
export class StatusBrowserClientToC {
	constructor(host: string){
		if(host==null||host==undefined||host.length==0){
			throw "StatusBrowserClientToC's host missing"
		}
		this.host=host
	}
	//timeout must be integer,timeout's unit is millisecond
	//don't set Content-Type in header
	ping(header: { [k: string]: string },req: Pingreq,timeout: number,errorf: (arg: Error)=>void,successf: (arg: Pingresp)=>void){
		if(!Number.isInteger(timeout)){
			errorf({code:-2,msg:'timeout must be integer'})
			return
		}
		if(header==null||header==undefined){
			header={}
		}
		header["Content-Type"] = "application/x-www-form-urlencoded"
		let form: string=''
		try{
			form=PingreqToForm(req)
		}catch(e){
			errorf({code:-2,msg:''+e})
			return
		}
		let config={
			url:_WebPathStatusPing+'?'+form,
			method: "get",
			baseURL: this.host,
			headers: header,
			timeout: timeout,
		}
		Axios.request(config)
		.then(function(response){
			let obj:Pingresp
			try{
				obj=JsonToPingresp(response.data)
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
