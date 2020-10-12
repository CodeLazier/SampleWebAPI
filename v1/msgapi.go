package v1

import (
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"test/handler"
	"test/msg"

	"github.com/gin-gonic/gin"
)

//type ResponseData struct {
//	Id     int         `json:"id"` //stub
//	ErrNo  int         `json:"errno"`
//	ErrMsg string      `json:"err"`
//	ByTime int64       `json:"bytime"`
//	Result interface{} `json:"result"`
//}

type NewEipMsg struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	result  chan interface{}
}

//test
func Test_PostNewEipMsg() error {
	e := NewEipMsg{Title: "test", Content: "test content"}
	e.result = make(chan interface{})
	switch v := postMsg(e).(type) {
	case error:
		return v
	case int:
		return nil
	}
	return nil
}

var postMsg = func() func(msg NewEipMsg) interface{} {
	msgChan := make(chan NewEipMsg)
	one := sync.Once{}
	func() {
		one.Do(func() {
			for i := 0; i < runtime.NumCPU()*8; i++ {
				go func() {
					for {
						xmsg := <-msgChan
						msg.NewEipDBHandler(func(eip *msg.EipMsgHandler) {
							if eip != nil {
								if err := eip.New(handler.EipMsg{
									Title:   xmsg.Title,
									Content: xmsg.Content,
								}); err != nil {
									xmsg.result <- err
								} else {
									xmsg.result <- 0 // success tag
								}
							}
						})
						//if eip, err := NewEipDBHandler(); err == nil {
						//	if err := eip.New(handler.EipMsg{
						//		Title:   msg.Title,
						//		Content: msg.Content,
						//	}); err != nil {
						//		msg.result <- err
						//	} else {
						//		msg.result <- 0 // success tag
						//	}
						//} else {
						//	log.Fatalln(err)
						//}
					}
				}()
			}
		})
	}()
	return func(msg NewEipMsg) interface{} {
		msgChan <- msg
		return <-msg.result
	}
}()

//func NewEipDBHandler(f func(*msg.EipMsgHandler)) {
//	dbctl := handler.GetMsgDB()
//	if dbctl == nil {
//		log.Println(fmt.Errorf("db conn is error"))
//	}
//	if f != nil {
//		func() {
//			defer func() {
//				handler.PutMsgDB(dbctl)
//			}()
//			f(
//				&msg.EipMsgHandler{
//					Control: dbctl,
//				})
//		}()
//	}
//}

//func NewResponseData(r interface{}, err error) ResponseData {
//	result := ResponseData{
//		ByTime: time.Now().Unix(),
//		Id:     0, //stub
//	}
//	if err != nil {
//		result.ErrNo = -1
//		result.ErrMsg = fmt.Sprint(err)
//	} else if r != nil {
//		result.Result = r
//	}
//	return result
//}

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
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized call"}) //  NewResponseData(nil, fmt.Errorf("unauthorized call")))
			c.Abort()
		} else {
			//do business
			c.Next()
		}
	}
}

func processCacheReq(c *gin.Context, f func()) (useCache bool, cacheTime time.Duration, fun func()) {
	useCache = false
	cacheTime = 0 * time.Second
	fun = f
	cc := strings.ToLower(c.GetHeader("x-cache"))
	if strings.Contains(cc, "x-expire") {
		sp := strings.Split(cc, "=")
		if len(sp) > 0 {
			useCache = true
			cacheTime = 5 * time.Second
			if v, err := strconv.Atoi(strings.TrimSpace(sp[1])); err != nil {
				log.Print(err)
			} else {
				cacheTime = time.Duration(v) * time.Second
			}
		}
	} else {
		useCache = !(cc == "" || cc == "no-cache" || cc == "no-store")
	}
	return
}

func DoGetMessagesCount() gin.HandlerFunc {
	return func(c *gin.Context) {
		msg.NewEipDBHandler(func(eip *msg.EipMsgHandler) {
			eip.UseCache(processCacheReq(c, func() {
				if r, err := eip.GetCount(); err != nil {
					log.Fatalln(err)
					//always return error code information instead of http status code
					c.Status(http.StatusInternalServerError)
				} else {
					c.JSON(http.StatusOK, gin.H{"count": r})
				}
			}))
		})
	}
}

func DoGetMessages() gin.HandlerFunc {
	return func(c *gin.Context) {
		var page, size int
		var err error
		size = -1
		p := c.Param("page")
		if p != "" && p != "/" {
			s := strings.Split(p[1:], ",")
			if len(s) > 0 && s[0] != "" {
				size = 30
				if page, err = strconv.Atoi(s[0]); err != nil {
					//
				} else if len(s) > 1 && s[1] != "" {
					if size, err = strconv.Atoi(s[1]); err != nil {
						//
					}
				}
			}
		}
		msg.NewEipDBHandler(func(eip *msg.EipMsgHandler) {
			eip.UseCache(processCacheReq(c, func() {
				if msgs, err := eip.GetAll(page*size, size); err != nil {
					log.Println(err)
					c.JSON(http.StatusOK, gin.H{"status": err.Error()})
				} else {
					c.JSON(http.StatusOK, msgs)
				}
			}))
		})
	}
}

func DoNewMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		emsg := NewEipMsg{}
		if err := c.BindJSON(&emsg); err != nil {
			log.Println(err)
			c.JSON(http.StatusOK, gin.H{"status": -1})
		} else {
			if emsg.Title == "" {
				c.JSON(http.StatusOK, gin.H{"status": -1})
			} else {
				emsg.result = make(chan interface{})
				switch v := postMsg(emsg).(type) {
				case error:
					log.Println(v)
					c.JSON(http.StatusOK, gin.H{"status": v.Error()})
				case int:
					c.JSON(http.StatusOK, gin.H{"status": v})
				}
			}
		}
	}
}

func DoMessagesMarkRead() gin.HandlerFunc {
	return func(c *gin.Context) {
		msg.NewEipDBHandler(func(eip *msg.EipMsgHandler) {
			if idx, err := strconv.Atoi(c.Param("id")); err != nil {
				log.Print(err)
			} else {
				if err := eip.MarkRead(idx); err != nil {
					log.Fatalln(err)
				} else {
					c.JSON(http.StatusOK, gin.H{"error": 0})
				}
			}
		})
	}
}

func DoGetMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		msg.NewEipDBHandler(func(eip *msg.EipMsgHandler) {
			if idx, err := strconv.Atoi(c.Param("id")); err != nil {
				log.Fatalln(err)
			} else {
				eip.UseCache(processCacheReq(c, func() {
					if emsg, err := eip.GetIndex(idx); err != nil {
						log.Print(err)
						//always return error code information instead of http status code
						c.Status(http.StatusNotFound)
					} else {
						c.JSON(http.StatusOK, emsg)
					}
				}))
			}
		})
	}
}

func GetToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusNotImplemented, "not impl")
	}
}

func GetTextContent() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(200, "test.html", nil)
	}
}

func parseToken(token string) error {
	_ = token
	return nil
}
