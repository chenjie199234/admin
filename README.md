# admin
```
admin是一个微服务.
运行cmd脚本可查看使用方法.windows下将./cmd.sh换为cmd.bat
./cmd.sh help 输出帮助信息
./cmd.sh pb 解析proto文件,生成桩代码
./cmd.sh sub 在该项目中创建一个新的子服务
./cmd.sh kube 新建kubernetes的配置
./cmd.sh html 新建前端html代码模版
```

### 服务端口
```
6060                                    MONITOR AND PPROF
8000                                    WEB
9000                                    CRPC
10000                                   GRPC
```

## 环境变量
```
LOG_TARGET                              日志等级,debug,info(default),warning,error
LOG_TRACE                               是否开启链路追踪,1-开启,0-关闭(default)
LOG_TARGET                              日志输出目标,std-输出到标准输出,file-输出到文件(可执行文件相同目录),both-两者都输出
GROUP                                   该项目所属的group(k8s的namespace),如果不使用k8s需要手动指定,如果使用k8s,需修改项目根目录的deployment.yaml中的<GROUP>
RUN_ENV                                 当前运行环境,如:test,pre,prod
DEPLOY_ENV                              部署环境,如:ali-kube-shanghai-1,ali-host-hangzhou-1
MONITOR                                 是否开启系统监控采集,0关闭,1开启
CONFIG_TYPE                             配置类型
                                        0-使用本地配置
                                        1-使用远程配置中心config服务
REMOTE_CONFIG_MONGO_URL                 当CONFIG_TYPE为1时,配置中心mongodb的url,[mongodb/mongodb+srv]://[username:password@]host1,...,hostN[/dbname][?param1=value1&...&paramN=valueN]
REMOTE_CONFIG_SECRET                    当CONFIG_TYPE为1时,配置中心配置的密钥
```

## 配置文件
```
AppConfig.json该文件配置了该服务需要使用的业务配置,可热更新
SourceConfig.json该文件配置了该服务需要使用的资源配置,不热更新
```

## DB
### Mongo(ReplicaSet mode)(Version >= 4.4)
#### config
```
database: app

collection: config
{
	"_id":ObjectId("xxxx"),
	"group":"",
	"app":"",
	"key":"",//always empty
	"index":0,//always be 0
	"paths":{
		"base64(/path1)":{
			"permission_node_id":"",
			"permission_read":true,
			"permission_write":true
		},
		"base64(/path2)":{
			"permission_node_id":"",
			"permission_read":true,
			"permission_write":true
		}
	},
	"keys":{
		"config_key1":{
			"cur_index":0,
			"max_index":0,//auto increment(every time insert a new config log)
			"cur_version":0,//auto increment(every time insert or rollback)
			"cur_value":"xxx"
		},
		"config_key2":{
			"cur_index":0,
			"max_index":0,//auto increment(every time insert a new config log)
			"cur_version":0,//auto increment(every time insert or rollback)
			"cur_value":"xxx"
		}
	},
	"value":"",//this is a random str + it's sha512 sign,this is used to check the secret
	"permission_node_id":"",
}//summary
{
	"group":"",
	"app":"",
	"_id":ObjectId("xxx"),
	"key":"config_key1",//always not empty
	"index":1,//always > 0
	"value":""
}//log
//手动创建数据库
use app;
db.createCollection("config");
db.config.createIndex({group:1,app:1,key:1,index:1},{unique:true});
db.config.createIndex({permission_node_id:1},{sparse:true,unique:true});
```
#### user
```
database: user

collection: user
{
	"_id":ObjectId("xxx"),//userid,if this is empty,means this is the super admin user
	"user_name":"",
	"password":"",
	"department":["",""],
	"ctime":123,//unixtimestamp,unit second
	"projects":["project1","project2"],
	"roles":["project1:role_name1","project2:role_name2"]
}
//手动创建数据库
use user;
db.createCollection("user");
db.user.createIndex({user_name:1});
db.user.createIndex({projects:1});
db.user.createIndex({roles:1});

collection: role
{
	"project":"",
	"role_name":"",
	"comment":"",
	"ctime":123,//unixtimestamp,unit second
}
//手动创建数据库
use user;
db.createCollection("role");
db.role.createIndex({project:1,role_name:1},{unique:true});
```
#### permission
```
database: permission

collection: node
{
	"node_id":"",
	"node_name":"",
	"node_data":"",
	"cur_node_index":0,
}
//手动创建数据库
use permission;
db.createCollection("node");
db.node.createIndex({node_id:1},{unique:true});

collection: usernode
{
	"user_id":ObjectId("xxx"),
	"node_id":"",
	"r":true,//can read
	"w":true,//can write
	"x":true,//admin
}
//手动创建数据库
use permission;
db.createCollection("usernode");
db.usernode.createIndex({user_id:1,node_id:1},{unique:true});
db.usernode.createIndex({node_id:1});

collection: rolenode
{
	"project":"",
	"role_name":"",
	"node_id":"",
	"r":true,//can read
	"w":true,//can write
	"x":true,//admin
}
//手动mongo创建数据库
use permission;
db.createCollection("rolenode");
db.rolenode.createIndex({project:1,role_name:1,node_id:1},{unique:true});
db.rolenode.createIndex({node_id:1});
```
