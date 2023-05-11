package model

import (
	"net"
	"project/group_chat/common/message"
)

//因为在客户端，我们很多地方会使用到curUser,我们将其作为一个全局的变量
type CurUser struct {
	Conn net.Conn
	message.User
} 