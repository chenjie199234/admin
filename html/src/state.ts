import { reactive,ref } from 'vue'
import * as userAPI from '../../api/user_browser_toc'
// import * as initializeAPI from '../../api/initialize_browser_toc'
import * as permissionAPI from '../../api/permission_browser_toc'

var obj = JSON.parse(localStorage.getItem(key))
user.root=obj.root
user.token=obj.token

//-------------------------------------------------------------------------------
export const inited = ref<boolean>(false)

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
export function set_alert(title: string,code: number,msg: string){
	alert.ing = true
	alert.title = title
	alert.code = code
	alert.msg = msg
	if(code==10004&&msg=="token wrong"){
		logout()
	}
}
export function clear_alert(){
	alert.ing = false
}
export function get_alert_title():string{
	if(alert.code==0){
		return alert.title
	}
	return alert.title+":"+alert.code
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
export function login(token:string){
	user.token=token
	localStorage.setItem("token",JSON.stringify({root:user.root,token:token}))
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

	localStorage.removeItem("token")

	clear_project()
	clear_page()
}

//-------------------------------------------------------------------------------
export const project = reactive<{
	cur_id:number[]
	cur_name:string
}>({
	cur_id:[],
	cur_name:"",
})

export function set_project(id:number[],name:string){
	project.cur_id=id
	project.cur_name=name
}
export function clear_project(){
	project.cur_id=[]
	project.cur_name=""
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
