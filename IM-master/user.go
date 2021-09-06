package main

import (
	"net"
	"fmt"
	"strings"
)

type User struct{
	Name string
	Addr string
	C chan string
	conn net.Conn
	server *Server
}

//创建UserAPI
func NewUser(conn net.Conn, server *Server) *User{
	UserAddr := conn.RemoteAddr().String()
	fmt.Println("useraddr:",UserAddr)
	user := &User{
		Name:UserAddr,
		Addr:UserAddr,
		C:make(chan string),
		conn:conn,
		server: server,
	}
	go user.ListenMessage()
	return user
}
//给当前User对应的客户端发送消息
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}
func (this *User) DoMessage(msg string){
	fmt.Println(msg,strings.Replace(msg, "\n", "", -1) == "who")
	if strings.Replace(msg, "\n", "", -1) == "who" {

		for _,user := range this.server.OnlineMap{
			// this.server.mapLock.Lock()
			fmt.Println(user.Name)
			// this.server.BroadCast(this,user.Name+"在线")
			this.SendMsg(user.Name+"在线\n")
			// this.server.mapLock.Unlock()
		}
		
	}
}

//监听当前User channer的方法，有消息发送给客户端
func (this *User) ListenMessage(){
	for{
		msg := <- this.C
		this.conn.Write([]byte(msg+"\n"))
		
	}
}
