A test project for msg sample and as a sample better practices are implemented.

All design architectures use golang style as much,may be...

This project is written in Go,use VSCode editor or Goland ¯\\_(ツ)_/¯

---

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
