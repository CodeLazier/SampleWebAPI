/*
 * @Author: your name
 * @Date: 2020-09-22 10:52:47
 * @LastEditTime: 2020-09-24 15:00:28
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
	"time"

	"github.com/gin-gonic/gin"
)

func verifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if parseToken(c.Query("token")) != nil {
			c.Status(http.StatusForbidden)
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

func doGetEipMessages(c *gin.Context) {
	eip := &msg.Eip{}
	if msgs, err := eip.GetUnread(); err == nil {
		c.JSON(http.StatusOK, msgs)
	} else {
		log.Print(err)
	}
}

func doEipMessagesMarkRead(c *gin.Context) {
	eip := &msg.Eip{}
	if idx, err := strconv.Atoi(c.Query("id")); err != nil {
		log.Print(err)
	} else {
		if err := eip.MarkRead(idx); err == nil {
			c.JSON(http.StatusOK, gin.H{"error": 0})
		} else {
			log.Print(err)
		}
	}
}

func doGetEipMessage(c *gin.Context) {
	eip := &msg.Eip{}
	if idx, err := strconv.Atoi(c.DefaultQuery("id", "0")); err != nil {
		log.Print(err)
	} else {
		if msg, err := eip.GetIndex(idx); err == nil {
			c.JSON(http.StatusOK, msg)
		} else {
			log.Print(err)
		}
	}
}

func getToken(c *gin.Context) {
	c.String(http.StatusNotImplemented, "not impl")

}

func parseToken(token string) error {
	return nil
}

func main() {
	r := gin.Default()
	r.Use(cors())
	r.Use(verifyToken())

	group := r.Group("/eip")

	group.GET("/getMessages", func(c *gin.Context) {
		doGetEipMessages(c)
	})
	group.GET("/getMessage", func(c *gin.Context) {
		doGetEipMessage(c)
	})
	group.POST("/setMessageMarkRead", func(c *gin.Context) {
		doEipMessagesMarkRead(c)
	})
	r.GET("/getToken", func(c *gin.Context) {
		getToken(c)
	})

	server := &http.Server{
		Addr:         ":9090",
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		//do config
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	<-ctx.Done()

}
