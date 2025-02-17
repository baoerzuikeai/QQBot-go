package voice

import (
	"encoding/base64"
	"log"
	"os"
)

func GetVoiceBase64() string {
	voice, err := os.ReadFile("./audio/帅哥要睡觉.ntsilk")
	if err != nil {
		log.Println("读取voice文件出错", err)
		return ""
	}
	voicebase64 := base64.StdEncoding.EncodeToString(voice)
	return voicebase64
}
