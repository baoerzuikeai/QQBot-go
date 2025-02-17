package rpcserver

import (
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"

	"github.com/baoer/QQBotWebHook/logs"
)

var clientConn *rpc.Client
var mutex sync.Mutex
var disconneted = make(chan struct{})

func SendHeartbeat() {
	for {
		if clientConn == nil {
			log.Println("Local server not connected")
			return
		}
		var reply string
		var heartbeattime time.Time
		err := clientConn.Call("HearbeatService.Heartbeat", &heartbeattime, &reply)
		if err != nil {
			log.Println("Failed to send heartbeat to local server:", err)
			closeClient()
			disconneted <- struct{}{}
			return
		} else {
			logs.Info("Server Reply:", reply)
		}
		time.Sleep(30 * time.Second)
	}
}

func handleClient(conn net.Conn, success chan string) {
	defer conn.Close()
	mutex.Lock()
	if clientConn != nil {
		clientConn.Close()
		clientConn = nil
	}

	clientConn = rpc.NewClient(conn)
	mutex.Unlock()
	//确保此链接断开时候，释放资源
	defer func() {
		if clientConn != nil {
			closeClient()
		}
	}()
	//发送检测链接消息
	var reply string
	var heartbeattime time.Time
	err := clientConn.Call("PayloadReceiver.Ping", &heartbeattime, &reply)
	if err != nil {
		log.Println("Failed to send payload to local server:", err)
		closeClient()
	}
	success <- reply
	//启动心跳检测
	go SendHeartbeat()
	<-disconneted
}

func Init(success chan string) {
	listener, err := net.Listen("tcp", ":10529")
	if err != nil {
		log.Println("Failed to start local server:", err)
	}
	logs.Info("local server started")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		go handleClient(conn, success)
	}
}

func SendPayloadToLocalServer(message *[]byte) {
	if clientConn == nil {
		log.Println("Local server not connected")
		return
	}
	var reply string
	err := clientConn.Call("PayloadReceiver.ReceiverPayload", message, &reply)
	if err != nil {
		log.Println("Failed to send payload to local server:", err)
		closeClient()
	} else {
		log.Println("Server Reply:", reply)
	}

}

func closeClient() {
	mutex.Lock()
	clientConn.Close()
	clientConn = nil
	mutex.Unlock()
}
