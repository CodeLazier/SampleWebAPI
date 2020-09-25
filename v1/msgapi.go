/*
 * @Author: your name
 * @Date: 2020-09-25 09:08:54
 * @LastEditTime: 2020-09-25 09:46:14
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \pre_work\v1\msgapi.go
 */
package v1

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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
		} else {
			//do business
			c.Next()
		}
	}
}

func DoGetMessages() gin.HandlerFunc {
	return func(c *gin.Context) {
		eip := &msg.Eip{
			Control: orm.NewOrmMock(),
		}
		if msgs, err := eip.GetAll(); err == nil {
			c.JSON(http.StatusOK, NewResponseData(msgs, err))
		} else {
			log.Print(err)
			c.JSON(http.StatusOK, NewResponseData(nil, err))
		}
	}
}

func DoMessagesMarkRead() gin.HandlerFunc {
	return func(c *gin.Context) {
		eip := &msg.Eip{
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
		eip := &msg.Eip{
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
