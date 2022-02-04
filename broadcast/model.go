package broadcast

type client struct {
	topic   string //用户所属主题
	groupId string //用户所属群组id
	userId  string //用户id
}

type message struct {
	groupId string //消息所属群组id
	userId  string //消息所属用户id
	data    []byte //要推送的消息数据
}
