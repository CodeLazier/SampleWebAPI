A test project for msg sample and as a sample better practices are implemented.

All design architectures use golang style as much,may be...

#### Possibility TODO:

  **大部分概念和實驗級证明,并不考慮穩定和效能**

* 構建CD:Register Docker Hub or gitlab docker hub工作流可以自動pull image方便部署和測試
* 對部署進行benchmark測試已檢測DB和API瓶頸,自動采樣pprof,對優化調校提供依據
* 接入api geteway,並對鏈路監控,預警,追溯,指標采樣,日志,融斷,降級,負載平衡,服務發現...等等高可用實現
* 日志落地和分割
* 緩存和分佈式
* DB橫向擴展
* 🔥 Hot reload

# Build & Test:

> Linux:  
>
> > 1. 下載並安裝docker & docker-compose  已有則略過
> > 2. - `git clone https://gitlab.com/ntsft/tsvc/prework/rain-end.git`
> >    - 或  [下載離綫包](https://gitlab.com/ntsft/tsvc/prework/rain-end/-/archive/master/rain-end-master.tar.gz)  
> > 3. 進入repository目錄  
> > 4.    執行 `docker-compose up --build -d`  
> >       如果一切順利會成功啓動docker container,並偵聽9090對外提供webapi服務  
> > 5. 瀏覽器開啓 http://IP:9090/eip/v1/msg/test  進入測試頁面
>
> > **僅提供測試,資料庫落地區未挂在物理盤,所有資料在container stop后消失**  
> > **docker版本默認沒有開啓TLS認證,如有需要可以進入修改配置開啓**



> Windows:(不推薦,可以虛擬機裏安裝)  
>
> > 1. ~~下載golang安裝包安裝~~ [去往Pipeline下載](https://gitlab.com/ntsft/tsvc/prework/rain-end/-/pipelines)
> > 2. 下載postgres安裝包安裝  [Download](https://sbp.enterprisedb.com/getfile.jsp?fileid=12851&_ga=2.269118450.286541361.1602680538-371199612.1601476970)
> > 3. 同Linux,clone倉庫進入repository目錄  
> > 4. 編修config.toml進行資料庫連綫配置  
> > 5. 執行`./test.exe`  
> > 6. 瀏覽器執行同Linux



> MacOS/Other:  
>
> > 同Windows,注意下載對應版本



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

```javascript
id: int64
title: string
content: string
createAt: time
```
APIs  
// 獲取單筆訊息通知  
METHOD  
GET  
PATH  
/v1/msg/{id}   
RETURN  

```javascript
{
  id,
  title,
  content,
  createAt
}
```



// 獲取訊息通知列表  
METHOD  
GET  
PATH  
/v1/msg/{page}/list  
RETURN  

```javascript
[
  {  
    id,
    title,
    content,
    createAt
  }
]
```

// 新增訊息通知  
METHOD  
POST  
PATH  
/v1/msg  
DATA  

```javascript
{
  title,
  content
}  
```

RETURN  

```javascript
{
  status
}
```