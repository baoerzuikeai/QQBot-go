package rpcclient

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"

	"github.com/baoerzuikeai/QQBot-go/logs"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/event"
)

type PayloadReceiver struct{}

func (p *PayloadReceiver) Ping(lastHeartbeat *time.Time, reply *string) error {
	logs.Info("Connected successfully")
	*lastHeartbeat = time.Now()
	logs.Info(*lastHeartbeat)
	*reply = "success"
	return nil
}

type HearbeatService struct{}

func (h *HearbeatService) Heartbeat(lastHeartbeat *time.Time, reply *string) error {
	logs.Info("Received heartbeat")
	*lastHeartbeat = time.Now()
	*reply = "success"
	return nil
}

func (p *PayloadReceiver) ReceiverPayload(message *[]byte, reply *string) error {
	fmt.Println(string(*message))
	payload := &dto.WSPayload{}
	_ = json.Unmarshal(*message, payload)
	payload.RawMessage = *message
	if payload.OPCode == dto.WSDispatchEvent {
		// 解析具体事件，并投递给业务注册的 handler
		if err := event.ParseAndHandle(payload); err != nil {
			log.Printf("handle event failed: %s", err)
		}
	}
	*reply = "success"
	return nil
}

func Init() {
	conn, err := net.Dial("tcp", "baoer.icu:10529")
	if err != nil {
		log.Panicln("Failed to connect to public server:", err)
	}
	defer conn.Close()
	logs.Info("Connected to public server")
	rpc.Register(&PayloadReceiver{})
	if err != nil {
		log.Panicln("Failed to register receiver:", err)
	}
	rpc.Register(&HearbeatService{})
	if err != nil {
		log.Panicln("Failed to register receiver:", err)
	}
	rpc.ServeConn(conn)
}
