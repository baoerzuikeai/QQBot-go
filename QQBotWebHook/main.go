package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/baoer/QQBotWebHook/logs"
	rpcserver "github.com/baoer/QQBotWebHook/rpc_server"
	mywebhook "github.com/baoer/QQBotWebHook/webhook"
	"github.com/tencent-connect/botgo/token"
)

var success = make(chan string, 1)
var newmessage = make(chan []byte, 1)

func main() {
	//创建oauth2标准token source
	credentials := &token.QQBotCredentials{
		AppID:     "102340632",
		AppSecret: "Rblv5FPalw7ITeq2EQco0CPcp2FSft7L",
	}
	tokenSource := token.NewQQBotTokenSource(credentials)
	//启动自动刷新access token协程
	if err := token.StartRefreshAccessToken(context.Background(), tokenSource); err != nil {
		log.Fatalln(err)
	}
	//启动 rpc_server
	go rpcserver.Init(success)
	for range success {
		logs.Info("Local server connected~")
		break
	}
	http.HandleFunc("/qqbot", func(writer http.ResponseWriter, request *http.Request) {
		payload := mywebhook.MyHTTPHandler(writer, request, credentials)
		jsondata, err := json.Marshal(payload)
		if err != nil {
			log.Println("json marshal error:", err)
		}
		newmessage <- jsondata
	})
	go func() {
		for {
			message := <-newmessage
			rpcserver.SendPayloadToLocalServer(&message)
		}
	}()
	// 启动http服务监听端口
	if err := http.ListenAndServeTLS(":443", "baoer.icu_bundle.crt", "baoer.icu.key", nil); err != nil {
		log.Fatal("setup server fatal:", err)
	}
}
