import { reactive,ref } from 'vue'
import * as userAPI from './api/admin_user_browser'
import * as permissionAPI from './api/admin_permission_browser'
import * as initializeAPI from './api/admin_initialize_browser'

//-------------------------------------------------------------------------------
export const inited = ref<boolean|null>(null)

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
	if(code==20104&&msg=="user not exist"){
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
	oauth2:string
	token:string
	info:userAPI.UserInfo|null
}>({
	root:false,
	oauth2:"",
	token:"",
	info:null,
})
export function login(oauth2:string,token:string){
	user.oauth2=oauth2
	user.token=token
	localStorage.setItem("token",JSON.stringify({root:user.root,oauth2:oauth2,token:token}))
}
export function avatar():string{
	if(user.root){
		return "R"
	}
	if(user.info){
		if(user.oauth2 == "FeiShu"){
			return user.info.feishu_user_name.substr(0,1)
		}else if(user.oauth2 == "DingDing"){
			return user.info.dingding_user_name.substr(0,1)
		}else if(user.oauth2 == "WXWork"){
			return user.info.wxwork_user_name.substr(0,1)
		}
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
	info:initializeAPI.ProjectInfo|null
}>({
	info:null,
})
export function clear_project(){
	project.info=null
}

//-------------------------------------------------------------------------------
export const page = reactive<{
	node:permissionAPI.NodeInfo|null
}>({
	node:null,
})
export function set_page(node:permissionAPI.NodeInfo){
	page.node=node
}
export function clear_page(){
	page.node=null
}
