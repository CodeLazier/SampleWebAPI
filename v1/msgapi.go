/*
 * @Author: your name
 * @Date: 2020-09-25 09:08:54
 * @LastEditTime: 2020-10-02 17:24:36
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \pre_work\v1\msgapi.go
 * 概念性的代码,缺少严谨和必要的重构,由于还没有实现鉴权所以缺少必要的如user等字段来获取资料
 */
package v1

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"test/handler"
	"test/msg"

	"time"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Id     int         `json:"id"` //stub
	ErrNo  int         `json:"errno"`
	ErrMsg string      `json:"err"`
	ByTime int64       `json:"bytime"`
	Result interface{} `json:"result"`
}

type NewEipMsg struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	result  chan interface{}
}

type PostMsgHandler func(msg NewEipMsg) interface{}

var post_eipmsg PostMsgHandler = func() PostMsgHandler {
	msgChan := make(chan NewEipMsg)
	one := sync.Once{}
	func() {
		one.Do(func() {
			/*
				*饥饿竞态的简单实现,如果饱和会阻塞提交直到db被缓解,如果单台db撑不了需要做分布式和负债平衡
				*设置每CPU开启8 grouting,根据CPU频率可以调高.但上限是DB的并发处理能力,设为1则为队列形式,但吞吐量下降
				*导致未充分压榨db并发和多核潜力且造成主观上的post性能效率下降
				*本机实测10000请求,1000并发量 平均每个请求响应在0.3s内,99%的请求量均控制在0.5s以下
				*具体见压测
				!此写瓶颈在于DB,提高DB能有效提高写入负载.超大并发年和流量请求因写入缓存系统并进入Job队列...
			*/
			for i := 0; i < runtime.NumCPU()*8; i++ {
				go func() {
					for {
						msg := <-msgChan
						if eip, err := NewEipDBHandler(true); err == nil {
							if err := eip.New(handler.EipMsg{
								Title:   msg.Title,
								Content: msg.Content,
							}); err != nil {
								msg.result <- err
							} else {
								msg.result <- 0
							}
						} else {
							log.Fatalln(err)
						}
					}
				}()
			}
		})
	}()
	return func(msg NewEipMsg) interface{} {
		msgChan <- msg
		return <-msg.result
	}
}()

func InitEipDBHandler() {
	handler.InitDB("user=postgres password=sasa dbname=postgres port=5432", false)

}

//call init first
func NewEipDBHandler(useCache bool) (*msg.EipMsgHandler, error) {
	if dbctl, err := handler.GetInstance(); err != nil {
		return nil, err
	} else {
		return &msg.EipMsgHandler{
			Control:  dbctl,
			UseCache: useCache,
		}, nil
	}
}

func NewResponseData(r interface{}, err error) ResponseData {
	result := ResponseData{
		ByTime: time.Now().Unix(),
		Id:     0, //stub
	}
	if err != nil {
		result.ErrNo = -1
		result.ErrMsg = fmt.Sprint(err)
	} else if r != nil {
		result.Result = r
	}
	return result
}

// func wrapResponseData(res ResponseData) (string, error) {
// 	if b, err := json.Marshal(&res); err != nil {
// 		return "", err
// 	} else {
// 		return string(b), nil
// 	}
// }

func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if parseToken(c.Query("token")) != nil {
			c.JSON(http.StatusUnauthorized, NewResponseData(nil, fmt.Errorf("Unauthorized call")))
			c.Abort()
		} else {
			//do business
			c.Next()
		}
	}
}

//TODO Pagination
func DoGetMessages() gin.HandlerFunc {
	return func(c *gin.Context) {
		//use db pool,every db connection is a slow operation
		eip, _ := NewEipDBHandler(true)
		if msgs, err := eip.GetAll(0, -1); err == nil {
			c.JSON(http.StatusOK, NewResponseData(msgs, err))
		} else {
			log.Println(err)
			c.JSON(http.StatusOK, NewResponseData(nil, err))
		}
	}
}

func DoNewMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		msg := NewEipMsg{}
		if err := c.BindJSON(&msg); err != nil {
			log.Println(err)
			c.JSON(http.StatusOK, gin.H{"status": -1})
		} else {
			msg.result = make(chan interface{})
			switch v := post_eipmsg(msg).(type) {
			case error:
				log.Println(v)
				c.JSON(http.StatusOK, gin.H{"status": v.Error()})
			case int:
				c.JSON(http.StatusOK, v)
			}
		}
	}
}

//TODO Pagination
func DoMessagesMarkRead() gin.HandlerFunc {
	return func(c *gin.Context) {
		eip, _ := NewEipDBHandler(true)
		if idx, err := strconv.Atoi(c.Param("id")); err != nil {
			log.Print(err)
		} else {
			if err := eip.MarkRead(idx); err == nil {
				c.JSON(http.StatusOK, gin.H{"error": 0})
			} else {
				log.Print(err)
			}
		}
	}
}

func DoGetMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		eip, _ := NewEipDBHandler(true)
		if idx, err := strconv.Atoi(c.Param("id")); err != nil {
			log.Print(err)
		} else {
			if msg, err := eip.GetIndex(idx); err == nil {
				c.JSON(http.StatusOK, msg)
			} else {
				log.Print(err)
			}
		}
	}
}

func GetToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusNotImplemented, "not impl")
	}
}

func parseToken(token string) error {
	return nil
}
