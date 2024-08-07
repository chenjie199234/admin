// Code generated by protoc-gen-markdown. DO NOT EDIT.<br />
// version:<br />
// 	protoc-gen-markdown v0.0.117<br />
// 	protoc              v5.27.2<br />
// source: api/admin_initialize.proto<br />

## initialize
### init_status
//初始化状态
#### Req:
```
Path:         /admin.initialize/init_status
Method:       POST
Content-Type: application/json
------------------------------------------------------------------------------------------------------------
{
}
------------------------------------------------------------------------------------------------------------
```
#### Resp:
```
Fail:    httpcode:4xx/5xx
------------------------------------------------------------------------------------------------------------
{"code":123,"msg":"error message"}
------------------------------------------------------------------------------------------------------------
Success: httpcode:200
------------------------------------------------------------------------------------------------------------
{
	//true-already inited,false-not inited
	"status":true
}
------------------------------------------------------------------------------------------------------------
```
### init
//初始化
#### Req:
```
Path:         /admin.initialize/init
Method:       POST
Content-Type: application/json
------------------------------------------------------------------------------------------------------------
{
	//value length must >= 10
	//value length must <= 32
	"password":"str"
}
------------------------------------------------------------------------------------------------------------
```
#### Resp:
```
Fail:    httpcode:4xx/5xx
------------------------------------------------------------------------------------------------------------
{"code":123,"msg":"error message"}
------------------------------------------------------------------------------------------------------------
Success: httpcode:200
------------------------------------------------------------------------------------------------------------
{
}
------------------------------------------------------------------------------------------------------------
```
### root_login
//登录
#### Req:
```
Path:         /admin.initialize/root_login
Method:       POST
Content-Type: application/json
------------------------------------------------------------------------------------------------------------
{
	//value length must >= 10
	//value length must <= 32
	"password":"str"
}
------------------------------------------------------------------------------------------------------------
```
#### Resp:
```
Fail:    httpcode:4xx/5xx
------------------------------------------------------------------------------------------------------------
{"code":123,"msg":"error message"}
------------------------------------------------------------------------------------------------------------
Success: httpcode:200
------------------------------------------------------------------------------------------------------------
{
	"token":"str"
}
------------------------------------------------------------------------------------------------------------
```
### update_root_password
//更新密码
#### Req:
```
Path:         /admin.initialize/update_root_password
Method:       POST
Content-Type: application/json
------------------------------------------------------------------------------------------------------------
{
	//value length must >= 10
	//value length must <= 32
	"old_password":"str",
	//value length must >= 10
	//value length must <= 32
	"new_password":"str"
}
------------------------------------------------------------------------------------------------------------
```
#### Resp:
```
Fail:    httpcode:4xx/5xx
------------------------------------------------------------------------------------------------------------
{"code":123,"msg":"error message"}
------------------------------------------------------------------------------------------------------------
Success: httpcode:200
------------------------------------------------------------------------------------------------------------
{
}
------------------------------------------------------------------------------------------------------------
```
### create_project
//创建项目
#### Req:
```
Path:         /admin.initialize/create_project
Method:       POST
Content-Type: application/json
------------------------------------------------------------------------------------------------------------
{
	//value length must != 0
	"project_name":"str",
	"project_data":"str"
}
------------------------------------------------------------------------------------------------------------
```
#### Resp:
```
Fail:    httpcode:4xx/5xx
------------------------------------------------------------------------------------------------------------
{"code":123,"msg":"error message"}
------------------------------------------------------------------------------------------------------------
Success: httpcode:200
------------------------------------------------------------------------------------------------------------
{
	//uint32
	"project_id":[1,2]
}
------------------------------------------------------------------------------------------------------------
```
### update_project
//更新项目
#### Req:
```
Path:         /admin.initialize/update_project
Method:       POST
Content-Type: application/json
------------------------------------------------------------------------------------------------------------
{
	//uint32
	//element num must == 2
	"project_id":[1,2],
	//if didn't change,set this with the old value
	//value length must != 0
	"new_project_name":"str",
	//if didn't change,set this with the old value
	"new_project_data":"str"
}
------------------------------------------------------------------------------------------------------------
```
#### Resp:
```
Fail:    httpcode:4xx/5xx
------------------------------------------------------------------------------------------------------------
{"code":123,"msg":"error message"}
------------------------------------------------------------------------------------------------------------
Success: httpcode:200
------------------------------------------------------------------------------------------------------------
{
}
------------------------------------------------------------------------------------------------------------
```
### list_project
//获取项目列表
#### Req:
```
Path:         /admin.initialize/list_project
Method:       POST
Content-Type: application/json
------------------------------------------------------------------------------------------------------------
{
}
------------------------------------------------------------------------------------------------------------
```
#### Resp:
```
Fail:    httpcode:4xx/5xx
------------------------------------------------------------------------------------------------------------
{"code":123,"msg":"error message"}
------------------------------------------------------------------------------------------------------------
Success: httpcode:200
------------------------------------------------------------------------------------------------------------
{
	//object project_info
	"projects":[{},{}]
}
------------------------------------------------------------------------------------------------------------
project_info: {
	//uint32
	"project_id":[1,2],
	"project_name":"str",
	"project_data":"str"
}
------------------------------------------------------------------------------------------------------------
```
### delete_project
//删除项目
#### Req:
```
Path:         /admin.initialize/delete_project
Method:       POST
Content-Type: application/json
------------------------------------------------------------------------------------------------------------
{
	//uint32
	//element num must == 2
	"project_id":[1,2]
}
------------------------------------------------------------------------------------------------------------
```
#### Resp:
```
Fail:    httpcode:4xx/5xx
------------------------------------------------------------------------------------------------------------
{"code":123,"msg":"error message"}
------------------------------------------------------------------------------------------------------------
Success: httpcode:200
------------------------------------------------------------------------------------------------------------
{
}
------------------------------------------------------------------------------------------------------------
```
