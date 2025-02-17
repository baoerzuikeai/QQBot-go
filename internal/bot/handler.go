package bot

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/baoerzuikeai/QQBot-go/internal/pixiv"
	"github.com/baoerzuikeai/QQBot-go/my_dto"
	"github.com/tencent-connect/botgo/dto"
)

func PostFile(data *dto.WSGroupATMessageData, filedata string) dto.MediaInfo {
	c := pixiv.InitClient()
	defer c.CloseIdleConnections()
	imagedata := my_dto.PostMedia{
		FileType:   1,
		SrvSendMsg: false,
	}
	imagedata.FileData = filedata
	jsonData, err := json.Marshal(imagedata)
	if err != nil {
		log.Println("映射图片json出错:", err)
	}
	url := "https://sandbox.api.sgroup.qq.com/v2/groups/" + data.GroupID + "/files"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("post图片出错", err)
	}
	req.Header.Set("Authorization", "Bot 102340632.ER0AW2JMtAA1G8PweWfXMjGZOoOCpbXB")
	req.Header.Set("Content-Type", "application/json")

	// 创建 HTTP 客户端
	client := &http.Client{}

	// 发起请求
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var media dto.MediaInfo
	json.Unmarshal(body, &media)
	return media
}
