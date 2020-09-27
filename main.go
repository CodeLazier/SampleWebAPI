/*
 * @Author: your name
 * @Date: 2020-09-22 10:52:47
 * @LastEditTime: 2020-09-27 15:36:28
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
	v "test/v1" //replace vx will if upgrade in the future
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	g.Use(cors.Default()) //Allow *
	eip := g.Group("/eip")

	v1 := eip.Group("/v1")
	v1.Use(v.VerifyToken())

	v1.GET("/getMessages", v.DoGetMessages())
	v1.GET("/getMessage/:id", v.DoGetMessage())
	v1.POST("/setMessageMarkRead/:id", v.DoMessagesMarkRead())
	v1.GET("/getToken", v.GetToken())

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

	func() {
		//It takes time to close, we give him time
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
	}()

}
