package main

import (
	"net"
)

type User struct{
	Name string
	Addr string
	C chan string
	Conn net.Conn
}

//创建一个用户API
func NewUser(conn net.Conn) *User{ 
	// 获取客户端的地址
	useraddr :=conn.RemoteAddr().String()

	user := &User{
		Addr : useraddr,
		Name : useraddr,
        C:     make(chan string),
		Conn : conn,
	}
	//启动监听当前User channel消息的goroutine
	go user.ListenMessage()
	return user
}

//监听当前User channal的方法，有消息 就发送到客户端  
func (this *User) ListenMessage(){
	for{
		msg := <-this.C

		this.Conn.Write([]byte(msg + "\n"))
	}
}
