/*
 * @Author: your name
 * @Date: 2020-09-25 09:08:54
 * @LastEditTime: 2020-09-28 15:30:43
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
	"strconv"
	"test/cache"
	"test/msg"
	"test/msg/orm"
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

func DoGetMessages() gin.HandlerFunc {
	return func(c *gin.Context) {
		//use db pool,every db connection is a slow operation
		eip := &msg.EipMsg{
			Control: orm.NewOrmMock(),
		}
		//we take it of the cache
		cc := cache.GetInstance()
		if v, err := cc.Get("username"); err != cache.ErrCacheNotFound {
			c.JSON(http.StatusOK, NewResponseData(v, err))
			log.Println("cached...")
		} else {
			//Cache penetration
			if msgs, err := eip.GetAll(msg.CustomWhere{}, 0, -1); err == nil {
				c.JSON(http.StatusOK, NewResponseData(msgs, err))
				_ = cc.Add("username", cache.NewCacheItem(msgs, 3*time.Second))
				log.Println("cache penetration")
			} else {
				log.Print(err)
				c.JSON(http.StatusOK, NewResponseData(nil, err))
			}
		}
	}
}

func DoMessagesMarkRead() gin.HandlerFunc {
	return func(c *gin.Context) {
		eip := &msg.EipMsg{
			Control: &orm.OrmMock{},
		}
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
		eip := &msg.EipMsg{
			Control: orm.NewOrmMock(),
		}
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
