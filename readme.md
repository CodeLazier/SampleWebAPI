A test project for msg sample and as a sample better practices are implemented.

All design architectures use golang style as much,may be...

---

測試頁面  
TODO:在綫部署暫不提供,如有需要callme,提供私有雲供内部測試.  

# Build & Test:

> Linux:  
> 下載並安裝docker & docker-compose  
> git clone https://gitlab.com/ntsft/tsvc/prework/rain-end.git  
> 或者下載  [離綫包](https://gitlab.com/ntsft/tsvc/prework/rain-end/-/archive/master/rain-end-master.tar.gz)  
> 進入rain-end目錄  
> 執行 docker-compose up --build -d  
> 如果一切順利會啓動二個docker container,並偵聽9090提供webapi服務  
> 瀏覽器開啓 http://宿主IP:9090/eip/v1/msg/test  進入測試頁面  
> *僅提供測試,資料庫未挂在物理盤,所有資料在container stop后消失  
> *~~如果遇到執行權限問題,請給sh script加上.chmod 777 *.sh~~  
> *docker版本默認沒有開啓TLS認證,如有需要可以進入修改配置開啓

> Windows:  
> 下載golang安裝包安裝  
> 下載postgres安裝包安裝  
> 同Linux,clone倉庫進入rain-end目錄  
> 開啓config.toml進行資料庫連綫配置  
> 執行go build編譯,並執行./test.exe  
> 瀏覽器執行同Linux

> MacOS/Other:  
> 未提供



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