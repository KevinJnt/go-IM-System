package main

import (
	"fmt"
	"net"
)
type Server struct{
	Ip string
	Port int
}

//创建一个server接口
func Newserver (ip string, port int) *Server{
	server := &Server{
		Ip : ip,
		Port : port,
	}
	return server
}

func (this *Server) Handler(conn net.Conn){
	fmt.Println("连接建立成功")
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