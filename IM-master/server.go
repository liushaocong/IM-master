package main

import (
	"net"
	"log"
	"bufio"
	"fmt"
	"sync"
)

type Server struct{
	Ip string
	Port int

	OnlineMap map[string] *User
	mapLock sync.RWMutex
	Message chan string
}

// 创建一个server
func NewServer(ip string,port int) *Server{
	server := &Server{
		Ip:ip,
		Port:port,
		OnlineMap:make(map[string]*User),
		Message:make(chan string),
	}
	return server
} 

//启动server
func (this *Server) Start(){
	// socket listen
	l, err := net.Listen("tcp", ":8888")
	fmt.Println("开始监听端口：","8888")
    if err != nil {
        log.Fatal(err)
    }
	
	defer l.Close()
	//启动监听消息的goroutine
	go this.ListenMessager()
	//socket accept
    for {
        conn, err := l.Accept()
		
        if err != nil {
            log.Fatal(err)
        }
        go this.HandleConnection(conn)
    }
	
}
//监听消息 发送给User
func (this *Server) ListenMessager(){
	for{
		msg := <-this.Message
		//将message发送给全部User
		this.mapLock.Lock()
		for _,cli := range this.OnlineMap{
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

//广播用户上线
func (this *Server) BroadCast(user *User,msg string){
	sendmsg := "["+user.Addr+"]"+user.Name + ":"+msg

	this.Message <- sendmsg
}



func (this *Server) HandleConnection(conn net.Conn){
	// defer conn.Close()
	fmt.Println("连接建立成功")

	// 将用户加入到User
	user := NewUser(conn,this)
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	fmt.Println("OnlineMap长度：",len(this.OnlineMap))

	this.BroadCast(user,"已上线")
	data,err := bufio.NewReader(conn).ReadString('\n')
	if err != nil{
		log.Fatal(err)
	}

	fmt.Println("from:",user.Name,"发送的消息：",string(data))
	// this.BroadCast(user,"from:"+user.Name+"发送的消息："+string(data))
	user.DoMessage(string(data))
	select {}
}