A test project for msg sample and as a sample better practices are implemented.

All design architectures use golang style as much,may be...

---

測試頁面  
TODO:製作測試容器方便測試.在綫部署暫不提供,如有需要callme,提供私有雲供内部測試.

| 说明                                     | method | url                            |
| ---------------------------------------- | ------ | ------------------------------ |
| 獲取所有條目數量                         | GET    | /eip/v1/msg/count              |
| 根據ID獲取信息                           | GET    | /eip/v1/msg/id/:id             |
| 批量獲取,page和size可選,空缺默認全部獲取 | GET    | /eip/v1/msg/list/[page],[size] |
| 新增                                     | POST   | /eip/v1/msg                    |
| 簡易測試頁面                             | GET    | /eip/v1/msg/test               |



未涉及到部分見(同)需求部分.

---
以下是需求

TABLE
messages

PROPERTIES  
id: int64
title: string
content: string
createAt: time

APIs

// 獲取單筆訊息通知
METHOD
GET
PATH
/v1/msg/{id} 
RETURN
{
	id,
	title,
	content,
	createAt
}

// 獲取訊息通知列表
METHOD
GET
PATH
/v1/msg/{page}/list
RETURN
[
	{
		id,
		title,
		content,
		createAt
	}
]

// 新增訊息通知
METHOD
POST
PATH
/v1/msg
DATA
{
	title,
	content
}
RETURN
{
	status
}
