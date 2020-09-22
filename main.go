/*
 * @Author: your name
 * @Date: 2020-09-22 10:52:47
 * @LastEditTime: 2020-09-22 14:24:25
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \test\main.go
 */
package main

import (
	"log"
	"net/http"
	"strconv"
	"test/msg"

	"github.com/gin-gonic/gin"
)

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		m := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if m == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}

func doGetEipMessages(c *gin.Context) {
	_ = c.Query("token")
	eip := &msg.Eip{}
	if msgs, err := eip.GetUnread(); err == nil {
		c.JSON(200, msgs)
	} else {
		log.Print(err)
	}
}

func doGetEipMessage(c *gin.Context) {
	_ = c.Query("token")
	eip := &msg.Eip{}
	if idx, err := strconv.Atoi(c.DefaultQuery("id", "0")); err != nil {
		log.Print(err)
	} else {
		if msg, err := eip.GetIndex(idx); err == nil {
			c.JSON(200, msg)
		} else {
			log.Print(err)
		}
	}
}

func getToken(c *gin.Context) {
	c.String(200, "not impl")
}

func main() {
	r := gin.Default()
	r.Use(cors())
	r.GET("/GetEipMessages", func(c *gin.Context) {
		doGetEipMessages(c)
	})
	r.GET("/GetEipMessage", func(c *gin.Context) {
		doGetEipMessage(c)
	})
	r.GET("/GetToken", func(c *gin.Context) {
		getToken(c)
	})

	//do config
	if err := r.Run(":9090"); err != nil {
		log.Fatal(err)
	}
}
