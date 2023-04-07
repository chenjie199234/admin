import { reactive,ref } from 'vue'
import * as userAPI from '../../api/user_browser_toc'
import * as initializeAPI from '../../api/initialize_browser_toc'
import * as permissionAPI from '../../api/permission_browser_toc'

//-------------------------------------------------------------------------------
export const load = reactive<{
	ing:boolean
}>({
	ing:false,
})
//return true:set success,false:set failed(already setted)
export function set_load(): boolean{
	if(!load.ing){
		load.ing = true
		return true
	}
	return false
}
export function clear_load(){
	load.ing = false
}

//-------------------------------------------------------------------------------
export const alert = reactive<{
	ing:boolean
	title:string
	code:number
	msg:string
}>({
	ing:false,
	title:"",
	code:0,
	msg:"",
})
export function set_error(title: string,code: number,msg: string){
	alert.ing = true
	alert.title = title
	alert.code = code
	alert.msg = msg
}
export function clear_error(){
	alert.ing = false
}

//-------------------------------------------------------------------------------
export const user = reactive<{
	root:boolean
	token:string
	info:userAPI.UserInfo
}>({
	root:false,
	token:"",
	info:null,
})
export function login(token:string,info:userAPI.UserInfo){
	user.token=token
	user.info=info
}
export function avatar():string{
	if(user.root){
		return "R"
	}
	if(user.info){
		return user.info.user_name.substr(0,1)
	}
	return ""
}
export function logout(){
	user.token=""
	user.info=null

	project.all=[]
	project.cur=""
	project.nodes=[]

	page.node=null
}

//-------------------------------------------------------------------------------
export const project = reactive<{
	all:initializeAPI.projectInfo[]
	cur:initializeAPI.ProjectInfo|string
	nodes:permissionAPI.NodeInfo[]
	ing:boolean
	optype:string// 'add' or 'del' or 'update'
	new_project_name:string
}>({
	all:[],
	cur:"",
	nodes:[],
	ing:false,
	optype:"",
	new_project_name:"",
})
export function set_project(optype: string){
	project.ing=true
	project.optype=optype
}
export function clear_project(){
	project.ing=false
	project.optype=""
	project.new_project_name=""
}

//-------------------------------------------------------------------------------
export const node = reactive<{
	ing:boolean
	target:permissionAPI.NodeInfo
	new_node_name:string
	new_node_url:string
	optype:string// 'add' or 'del' or 'update'
}>({
	ing:false,
	target:null,
	new_node_name:"",
	new_node_url:"",
	optype:"",
})
//optype: 'add' or 'del' or 'update'
export function set_node(target:permissionAPI.NodeInfo,optype:string){
	node.ing = true
	node.target = target
	node.optype = optype
}
export function clear_node(){
	node.ing=false
	node.node_id=[]
	node.node_name=""
	node.new_node_name=""
	node.new_node_url=""
	node.optype=""
}

//-------------------------------------------------------------------------------
export const page = reactive<{
	node:permissionAPI.NodeInfo
}>({
	node:null,
})
export function set_page(node:permissionAPI.NodeInfo){
	page.node=node
}
export function clear_page(){
	page.node=null
}
