package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
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
	Addr  string `toml:"addr"`
	Debug bool   `toml:"debug"`
}

//func Listen(addr string) (net.Listener, error) {
//	var ln net.Listener
//	if strings.HasPrefix(addr, "systemd:") {
//		name := addr[8:]
//		listeners, _ := activation.ListenersWithNames()
//		listener, ok := listeners[name]
//		if !ok {
//			return nil, fmt.Errorf("listen systemd %s: socket not found", name)
//		}
//		ln = listener[0]
//	} else {
//		var err error
//		ln, err = net.Listen("tcp", addr)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	return ln, nil
//}
func NewServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		//std sec for internet
		TLSConfig: &tls.Config{
			NextProtos:       []string{"h2", "http/1.1"},
			MinVersion:       tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{tls.CurveP256, tls.X25519},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
			PreferServerCipherSuites: true,
		},
	}
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
	if _, err := v.NewEipDBHandler(true); err != nil {
		log.Fatalln(err)
	}

	server := NewServer(cfg.Server.Addr, g)

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
