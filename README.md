# pekoIM

## 基于Go+WebSocket的IM群聊消息广播推送工具

***

### 引入包

* Goland终端执行：`go get github.com/dpwgc/pekoIM`

***

### 服务端使用说明

```
//导入github包
import (
    "github.com/dpwgc/pekoIM/broadcast"
)
```

##### 创建一个广播推送服务
* 运行端口为8080
* 主题为"test"
* 消息通道缓冲区大小为10000
* 处理函数为handler

```
broadcast.New("test",10000)     //创建一个主题为test的广播
broadcast.Run("8080",handler)   //在8080端口启动一个http监听服务
```

##### handler处理函数，可在此验证用户身份，作为参数传入Run()函数

```
//http请求处理函数
func handler(w http.ResponseWriter, r *http.Request) {

    //获取这个http请求的uri
    //uri: /{topic}/{groupId}/{userId}/{token}
    uri := r.RequestURI
    
    //切割uri得到topic、groupId、userId、token
    split := strings.Split(uri,"/")

    topic := split[1]
    groupId := split[2]
    
    //获取到的userId用户id和token令牌可用于身份验证
    userId := split[3]
    token := split[4]
    
    /*
        用户身份验证业务，略。。。
    */
	
    //如果通过身份验证
    //Add 将该客户端的http连接升级为websocket连接，并加入广播服务
    broadcast.Add(w,r)
}
```

##### 注，一个项目内可创建多个不同主题的广播服务

* 例：
```
//http请求处理函数
func handler(w http.ResponseWriter, r *http.Request) {
    
    /*
        用户身份验证业务，略。。。
    */
	
    //如果通过身份验证
    //Add 将该客户端的http连接升级为websocket连接，并加入广播服务
    broadcast.Add(w,r)
}

func main() {

    //创建多个广播服务
    broadcast.New("test_1",10000)
    broadcast.New("test_2",10000)
    broadcast.New("test_3",10000)
    broadcast.New("test_4",10000)
    
    broadcast.Run("8080",handler)   //在8080端口启动一个http监听服务
}
```

***

### 客户端使用说明

##### 客户端websocket链接规则
> ws://{addr}:{port}/{topic}/{groupId}/{userId}/{token}
* addr `服务端IP地址`
* port `端口号`
* topic `广播主题`
* groupId `群组id`
* userId `用户id`
* token `用户令牌（用于自定义验证用户身份功能，可有可无）`

##### 客户端websocket链接示范
> ws://127.0.0.1:8080/test/111/222/token123

该客户端所连接的广播主题为"test"，连接到的群组id为"111"，用户id为"222"，用户登录令牌为"token123"
