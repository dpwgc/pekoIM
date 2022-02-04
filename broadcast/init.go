package broadcast

import (
	"fmt"
	"log"
	"net/http"
)

// Run 启动http服务
func Run(port string, handler func(w http.ResponseWriter, r *http.Request)) {

	//监听http请求
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(fmt.Sprintf("%s%s", "127.0.0.1:", port), nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}

// New 新建一个广播
func New(topic string, buffer int) {

	//初始化广播推送服务
	InitBroadcast(topic, buffer)
}
