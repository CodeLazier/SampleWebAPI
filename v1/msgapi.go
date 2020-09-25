/*
 * @Author: your name
 * @Date: 2020-09-25 09:08:54
 * @LastEditTime: 2020-09-25 09:13:33
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \pre_work\v1\msgapi.go
 */
package v1

import (
	"log"
	"net/http"
	"strconv"
	"test/msg"
	"test/msg/orm"

	"github.com/gin-gonic/gin"
)

func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if parseToken(c.Query("token")) != nil {
			c.String(http.StatusUnauthorized, "Unauthorized call")
		} else {
			//do business
			c.Next()
		}
	}
}

func DoGetMessages() gin.HandlerFunc {
	return func(c *gin.Context) {
		eip := &msg.Eip{}
		if msgs, err := eip.GetUnread(); err == nil {
			c.JSON(http.StatusOK, msgs)
		} else {
			log.Print(err)
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
