package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct{
	Ip string
	Port int
	
	//在线用户的列表
	OnlineMap map[string]*User
	mapLock sync.RWMutex

	//消息广播的channal
	Message chan string
}

//创建一个server接口
func Newserver (ip string, port int) *Server{
	server := &Server{
		Ip :ip,
		Port : port,
		OnlineMap : make(map[string]*User),
		Message : make(chan string),
	}
	
	return server
}

//监听Messag广播消息，channal的geoutine,一旦有消息就发给全部OnlineUser
func (this *Server) ListenMessage(){
	for{
		msg := <- this.Message

		//将msg发送给全部的OnlineUser
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap{
			cli.C <- msg
		}

		this.mapLock.Unlock()
	}
}

//广播消息方法
func (this *Server) BroadCast(user *User , msg string ){
	sendMsg := "[" +user.Addr + "]" + user.Name + msg

	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn){
	//fmt.Println("连接建立成功")

	user := NewUser(conn)
	
	//用户上线，将用户加入OnlineMap
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	//广播当前用户上线消息
	this.BroadCast(user,"已上线")

	//当前Handler阻塞
	select{}

}

//启动服务器端口
func (this *Server) Start(){
	//listen
	listen, err := net.Listen("tcp",fmt.Sprintf("%s:%d",this.Ip,this.Port))
	if err != nil{
		fmt.Println("Listen error:",err)
		return
	}
	
	//close listen socket
	defer listen.Close()
	
	//启动监听Message的goroutine
	go this.ListenMessage()

	for {
		//accept
		conn,err := listen.Accept()
		if err != nil{
			fmt.Println("accept error:",err)
			continue
		}
		go this.Handler(conn)
	}
}
