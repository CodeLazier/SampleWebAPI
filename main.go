package main

import (
	"context"
	"crypto/tls"
	"log"
	"os"
	"os/signal"
	"syscall"
	xserver "test/server"
	"time"

	"test/handler"
	v "test/v1" //replace vx will if upgrade in the future

	"github.com/gin-contrib/cors"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
)

type Config struct {
	DB     DBConfig
	Server ServerConfig
}

type DBConfig struct {
	Conn  string `toml:"conn"`
	Debug bool   `toml:"debug"`
}

type ServerConfig struct {
	Addr     string `toml:"addr"`
	Debug    bool   `toml:"debug"`
	UseTLS   bool   `toml:"useTLS"`
	CertFile string `toml:"cert"`
	KeyFile  string `toml:"key"`
}

func main() {
	cfg := Config{}
	if _, err := toml.DecodeFile("./config.toml", &cfg); err != nil {
		log.Fatalln(err)
	}
	if !cfg.Server.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	g := gin.Default()
	g.Use(cors.Default()) //Allow *
	g.LoadHTMLFiles("v1/static/test.html")
	eip := g.Group("/eip")

	{
		v1 := eip.Group("/v1")
		v1.Use(v.VerifyToken())

		//get msg record info,ex: msg/id/39
		v1.GET("/msg/id/:id", v.DoGetMessage())
		//fetch msgs info,ex:
		//all msgs -> "msg/list/" eq 0,-1
		//fetch msgs -> "msg/list/page,count"
		//fetch msgs -> "msg/list/page" eq page,30
		v1.GET("/msg/list/*page", v.DoGetMessages())
		//create msg record
		v1.POST("/msg", v.DoNewMessage())
		//get all msgs count
		v1.GET("/msg/count", v.DoGetMessagesCount())
		v1.GET("/getToken", v.GetToken())
		//test
		v1.GET("/msg/test", v.GetTextContent())
		//v1.POST("/setMessageMarkRead/:id", v.DoMessagesMarkRead())
	}

	//init db
	handler.InitDB(cfg.DB.Conn, cfg.DB.Debug)
	if _, err := v.NewEipDBHandler(); err != nil {
		log.Fatalln(err)
	}

	//TODO CSR x.509(PEM),may also need other formats(DER,P12?)
	var server *xserver.Server
	if cfg.Server.UseTLS {
		cert, err := tls.LoadX509KeyPair(cfg.Server.CertFile, cfg.Server.KeyFile)
		if err != nil {
			log.Println(err)
		} else {
			server = xserver.NewTLS(cfg.Server.Addr, cert, g)
		}
	}
	if server == nil {
		server = xserver.New(cfg.Server.Addr, g)
	}

	go func() {
		if err := server.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	<-func() <-chan os.Signal {
		q := make(chan os.Signal, 1)
		signal.Notify(q, os.Interrupt, os.Kill, syscall.SIGTERM)
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
