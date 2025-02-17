package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/baoerzuikeai/QQBot-go/internal/bot"
	deepseek "github.com/baoerzuikeai/QQBot-go/internal/deepseek"
	"github.com/baoerzuikeai/QQBot-go/internal/pixiv"
	rpcclient "github.com/baoerzuikeai/QQBot-go/rpc_client"
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
)

var api openapi.OpenAPI
var pixiv_client = pixiv.InitClient()

func main() {

	credentials := &token.QQBotCredentials{
		AppID:     "102340632",
		AppSecret: "Rblv5FPalw7ITeq2EQco0CPcp2FSft7L",
	}
	tokenSource := token.NewQQBotTokenSource(credentials)
	//启动自动刷新access token协程
	if err := token.StartRefreshAccessToken(context.Background(), tokenSource); err != nil {
		log.Fatalln(err)
	}
	api = botgo.NewSandboxOpenAPI(credentials.AppID, tokenSource).WithTimeout(5 * time.Second).SetDebug(true)
	// 注册事件处理函数
	_ = event.RegisterHandlers(
		// 注册c2c消息处理函数
		C2CMessageEventHandler(),
		GroupATMessageEventHandler(),
	)

	rpcclient.Init()

}

// C2CMessageEventHandler 实现处理 at 消息的回调
func C2CMessageEventHandler() event.C2CMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSC2CMessageData) error {
		api.PostC2CMessage(context.Background(), data.Author.ID, &dto.MessageToCreate{
			MsgID:   data.ID,
			MsgType: 0,
			Content: deepseek.Chat(data.Content),
		})
		return nil
	}
}

func GroupATMessageEventHandler() event.GroupATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSGroupATMessageData) error {
		switch {
		case strings.HasPrefix(data.Content, " /帮助 "):
			api.PostGroupMessage(context.Background(), data.GroupID, &dto.MessageToCreate{
				MsgType: dto.TextMsg,
				Content: ` 
				🌟 功能菜单 🌟

1. /帮助 ❓ - 获取帮助信息，我会告诉你所有的小秘密哦～(｡♥‿♥｡)
2. /随机图片 📸 - 让我为你挑选一张可爱的随机图片，猜猜会是什么呢？(๑•̀ㅂ•́)و✧
3. /搜索图片 keywords 🔍 - 告诉我你想找什么，我会帮你找到最棒的图片！(≧▽≦)
4. 直接@我并发送消息 💬 - 我会立刻回复你，陪你聊天哦～(｡･ω･｡)ﾉ

请输入指令（例如：/帮助）来开始吧！(｡♡‿♡｡)：`,
				MsgID: data.ID,
			})
		case strings.HasPrefix(data.Content, " /随机图片 "):
			filedata := pixiv.Getimage(pixiv_client)
			for len(filedata) < 100 {
				filedata = pixiv.Getimage(pixiv_client)
			}
			fileinfo := bot.PostFile(data, filedata)
			api.PostGroupMessage(context.Background(), data.GroupID, &dto.MessageToCreate{
				MsgID:   data.ID,
				MsgType: dto.RichMediaMsg,
				Content: "随机图片",
				Media:   &fileinfo,
			})
		// case strings.HasPrefix(data.Content, " /搜索图片 "):
		default:
			api.PostGroupMessage(context.Background(), data.GroupID, &dto.MessageToCreate{
				MsgID:   data.ID,
				MsgType: dto.TextMsg,
				Content: deepseek.Chat(data.Content),
			})
		}
		return nil
	}
}
