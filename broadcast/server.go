package broadcast

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
	"sync"
)

//连接的客户端,把每个客户端都放进来。Key为Client结构体。Value为websocket连接
var clients = make(map[client]*websocket.Conn)

//广播通道列表，用于广播推送群聊用户发送的消息(带缓冲区，提高并发速率)
var broadcast = make(map[string]chan message)

var cLock sync.RWMutex
var bLock sync.RWMutex

// InitBroadcast 初始化广播推送服务
func InitBroadcast(topic string, buffer int) {

	//初始化该主题的广播通道
	bLock.RLock()
	broadcast[topic] = make(chan message, buffer)
	bLock.RUnlock()

	//启动该主题的广播推送协程
	go push(topic)

	fmt.Printf("%s%s\n", "[start broadcast] ", topic)
}

//websocket连接升级与跨域配置
var upGrader = websocket.Upgrader{
	//跨域设置
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Add 添加websocket客户端
func Add(w http.ResponseWriter, r *http.Request) {

	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		_, err = fmt.Fprintln(w, "error")
		if err != nil {
			fmt.Printf("\033[1;31;40m%s\033[0m\n", err)
			return
		}
	}

	//获取uri，切割得到topic、groupId、userId
	uri := r.RequestURI
	split := strings.Split(uri, "/")

	//uri: /{topic}/{groupId}/{userId}/{token}

	c := client{
		topic:   split[1],
		groupId: split[2],
		userId:  split[3],
	}

	//函数报错返回时关闭websocket连接
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			fmt.Printf("\033[1;31;40m%s\033[0m\n", err)
		}
	}(ws)

	//将该客户端加入客户端列表
	cLock.RLock()
	clients[c] = ws
	cLock.RUnlock()

	fmt.Printf("%s%s]%s%s]%s%s]\n", "[topic:", c.topic, " [groupId:", c.groupId, " [userId:", c.userId)

	//循环监听该客户端发送到广播服务的消息
	for {
		//读取websocket发来的数据
		_, data, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			delete(clients, c) //删除map中的客户端
			break
		}

		//封装消息
		m := message{
			groupId: c.groupId,
			userId:  c.userId,
			data:    data,
		}

		//将消息推送给其他相同主题的websocket客户端
		broadcast[c.topic] <- m
	}
}

//广播推送消息
func push(topic string) {
	for {
		//读取通道中的消息
		m := <-broadcast[topic]

		//轮询现有的websocket客户端
		for c, ws := range clients {

			//匹配客户端，判断该客户端的主题、群组id是否与该消息的主题、群组id一致，如果是，则将该消息投递给该客户端
			if m.groupId == c.groupId && topic == c.topic {
				//发送消息到对应客户端
				err := ws.WriteMessage(1, m.data)
				if err != nil {
					//客户端关闭
					err = ws.Close()
					if err != nil {
						fmt.Printf("\033[1;31;40m%s\033[0m\n", err)
					}
					//删除map中的客户端
					delete(clients, c)
				}
			}
		}
	}
}
