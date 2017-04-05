通用权限系统

**interface
[get] /user/list?idx=&size=
output:
status: 200
	{
		"users":[{"id":xx,"name":"xxx"},...]
		"next_idx":xxx,
		"total_size":xxx
	}
ERROR:
status: 500 
{
"err_msg":"xxx"
}

[get] /user/get
input: 
{
"id":xxx,	       //2选一，优先使用id查询
"code":"xxx"  //2选一
}
output:
status: 200 
{
"id":xxx,
"code":"xxx",
"name":"xxx",
"authoritys": [
{"id":xxx, "code":"xxx", "name":"xxx", "has_sub":"xxx", "sub_group":{...}},
...
]
}
has_sub字段解析：表示该项权限有没有子组。取值如下
none	没有子组，为叶节点
partial	有子组，但只拥有部分子组权限
all		有子组，且拥有全部子组权限
ERROR:
status: 400   ID do not esxit
status: 500   server error
{
"err_msg":"xxx"
}

[get] /user/auth
input: 
{
"id":xxx,	       //2选一，优先使用id查询
"code":"xxx",  //2选一
"auth":"xxx"   //可选
}
output: 
{
"has_auth":0 or 1  //参数中有auth字段时返回
"authoritys": [
{"id":xxx, "code":"xxx", "name":"xxx", "has_sub":"xxx", "sub_group":{...}},
...
]  //参数中无auth字段时返回
}


[post] /user/register
input:
{
"code":"xxx",
	"name":"xxx"
}
output:
status: 200

[put] /user/authority/grant
input: 
[
{"id":xx, "code":"xxx"}  //id和code二选一
...
]
output:
status: 200

[delete] /user/delete
intput:
{
"id":xxx,
"code":"xxx"  //id和code二选一
}
output:
status:200

[delete] /user/delete/id/{id}
[delete] /user/delete/code/{code}

[post] /authority/register
intput:
[{
"code": "xxx",
"name": "xxx",
"group_id": xxx   //可选, 没有该参数code必须是绝对路径
},...
]
output:
status:200
ERROR
status:400  code为绝对路径，且有group_id参数，但code路径和group_id不匹配
status:500
{
"err_msg":"xxx"
}

[delete] /authority/delete
intput:
[{
"id":xxx.
"code":"xxx"  //二选一
}，...
]

[get] /authority/get[/{group_id}]
返回某一层的权限
input:
{
"group_id":xxx    //0或者没有设置该参数则取顶层权限
}
output:
[
{
"id":xxx,
"code": "xxx",
"name": "xxx",
"group_id": xxx   //可选, 没有该参数code必须是绝对路径
},...
]


db_authority
tb_user
id	bigint 	primary
code varchar(256)   uniquie
name varchar(64) 

tb_authority
id	bigint	primary
code		varchar(256)	uniquie
name	varchar(64)
group_id	bigint

tb_user_authority
id	bigint	primary
user_id	bigint
authority_id	bigint


