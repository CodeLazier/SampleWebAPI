/*
 * @Author: your name
 * @Date: 2020-09-22 10:52:47
 * @LastEditTime: 2020-09-24 22:09:13
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \test\main.go
 */
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"test/msg"
	"test/msg/orm"
	"time"

	"github.com/gin-gonic/gin"
)

func verifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if parseToken(c.Query("token")) != nil {
			c.String(http.StatusUnauthorized, "Unauthorized call")
		} else {
			//do business
			c.Next()
		}
	}
}

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

func doGetEipMessages() gin.HandlerFunc {
	return func(c *gin.Context) {
		eip := &msg.Eip{}
		if msgs, err := eip.GetUnread(); err == nil {
			c.JSON(http.StatusOK, msgs)
		} else {
			log.Print(err)
		}
	}
}

func doEipMessagesMarkRead() gin.HandlerFunc {
	return func(c *gin.Context) {
		eip := &msg.Eip{
			Orm: &orm.OrmMock{},
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

func doGetEipMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		eip := &msg.Eip{
			Orm: orm.NewOrmMock(),
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

func getToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusNotImplemented, "not impl")
	}
}

func parseToken(token string) error {
	return nil
}

func shutdown(server *http.Server) {
	//It takes time to close, we give him time
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}

func main() {
	g := gin.Default()
	g.Use(cors())
	eip := g.Group("/eip")

	v1 := eip.Group("/v1")
	v1.Use(verifyToken())

	v1.GET("/getMessages", doGetEipMessages())
	v1.GET("/getMessage/:id", doGetEipMessage())
	v1.POST("/setMessageMarkRead/:id", doEipMessagesMarkRead())
	v1.GET("/getToken", getToken())

	//do config
	server := &http.Server{
		Addr:         ":9090",
		Handler:      g,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-func() <-chan os.Signal {
		q := make(chan os.Signal, 1)
		signal.Notify(q, os.Interrupt)
		return q
	}()

	shutdown(server)

}
