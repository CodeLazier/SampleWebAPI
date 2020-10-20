package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/arl/statsviz"
	"test/config"
	"test/msg"
	"test/xserver"

	"test/handler"
	v "test/v1" //replace vx will if upgrade in the future

	"github.com/gin-contrib/cors"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
)

var CUR_PATH = getCurPath()

func main() {
	cfg := loadConfig()
	initDB(cfg)
	sawServer(cfg, initGin(cfg))
}

func initDB(cfg config.Config) {
	//init and test db
	handler.InitDB(cfg.DB.Conn, cfg.DB.Debug)
	msg.NewEipDBHandler(func(eip *msg.EipMsgHandler) {
		if eip == nil {
			fmt.Println("DB connection is failed")
		} else {
			installDB(cfg)
		}
	})
}

func initGin(cfg config.Config) *gin.Engine {
	if !cfg.Server.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	g := gin.Default()
	g.Use(cors.Default()) //Allow *
	g.LoadHTMLFiles("v1/static/test.html")

	stepEipRouter(g.Group("/eip"))
	return g
}

func loadConfig() config.Config {
	//TODO 增加配置文件的hot reload实现,虽然没有什么好实时改变的,数据库的Debug模式?
	cfgFile := filepath.Join(CUR_PATH, "config.toml")
	if !fileExist(cfgFile) {
		log.Fatalln("cfg file is not found", cfgFile)
	} else {
		fmt.Println("config is reading...", cfgFile)
	}
	cfg := config.Config{}
	if _, err := toml.DecodeFile(cfgFile, &cfg); err != nil {
		log.Fatalln(err)
	}
	return cfg
}

func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func getCurPath() string {
	if path, err := os.Executable(); err != nil {
		return "."
	} else {
		ext := strings.ToLower(filepath.Ext(path))
		if !(ext == "" || ext == ".exe") {
			if r, err := filepath.EvalSymlinks(path); err != nil {
				return "."
			} else {
				return r
			}
		}
		return filepath.Dir(path)
	}
}

func stepEipRouter(eip *gin.RouterGroup) {
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

func sawServer(cfg config.Config, g *gin.Engine) {
	//TODO CSR x.509(PEM),may also need other formats(DER,P12?)
	var server *xserver.Server
	if cfg.Server.UseTLS {
		//test file is exist
		check := func(f string) string {
			if !fileExist(f) {
				if fileExist(filepath.Join(CUR_PATH, f)) {
					return filepath.Join(CUR_PATH, f)
				}
				return ""
			}
			return f
		}
		cfg.Server.CertFile = check(cfg.Server.CertFile)
		cfg.Server.KeyFile = check(cfg.Server.KeyFile)
		if cfg.Server.CertFile != "" && cfg.Server.KeyFile != "" {
			log.Println("use", cfg.Server.CertFile, cfg.Server.KeyFile)
			cert, err := tls.LoadX509KeyPair(cfg.Server.CertFile, cfg.Server.KeyFile)
			if err != nil {
				log.Println(err)
			} else {
				server = xserver.NewTLS(cfg.Server.Addr, cert, g)
			}
		} else {
			log.Println("crt or key file is not exist")
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

	if cfg.Server.Debug {
		registerPprof()
	}

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

func installDB(cfg config.Config) error {
	if dbctrl, err := handler.NewMsgDB(
		handler.MsgDBConfig{
			DBConn: cfg.DB.Conn,
		}); err != nil {
		return err
	} else {
		if !dbctrl.RawDB.Migrator().HasTable(&handler.EipMsg{}) {
			if err = dbctrl.RawDB.AutoMigrate(&handler.EipMsg{}); err != nil {
				log.Fatalln(err)
			}
		}
	}
	return nil
}

func registerPprof() {
	mux := http.NewServeMux()
	statsviz.Register(mux)
	statsviz.RegisterDefault()
	go func() {
		http.ListenAndServe(":9091", mux)
	}()
}
