package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kasiss-liu/smile"
)

//引用了gorilla的websocket库
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//ws处理程序
func websocketFunc(c *smile.Combination) error {
	conn, err := upgrader.Upgrade(c.ResponseWriter, c.Request, nil)

	if err != nil {
		log.Println(err)
		return nil
	}

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			// fmt.Println(err.Error())
			break
		}
		fmt.Println("Client say :" + string(message))
		err = conn.WriteMessage(messageType, []byte("hello"))
		if err != nil {
			// fmt.Println(err.Error())
			break
		}
	}
	return nil
}

func main() {
	//获取一个服务器引擎
	engine := smile.NewEngine("./")
	//注册路由
	routeGroup := smile.NewRouteGroup()
	routeGroup.SetWS("ws", websocketFunc)
	engine.SetRouteGroup(routeGroup)
	//注册监视器
	smile.MonitorOn()
	engine.GzipOn()
	engine.Run(":8000")
}
