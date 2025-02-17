package mywebhook

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/tencent-connect/botgo/constant"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/interaction/signature"
	"github.com/tencent-connect/botgo/interaction/webhook"
	"github.com/tencent-connect/botgo/log"
	"github.com/tencent-connect/botgo/token"
)

func MyHTTPHandler(w http.ResponseWriter, r *http.Request, credentials *token.QQBotCredentials) *dto.WSPayload {
	defer r.Body.Close()
	body := make([]byte, r.ContentLength)
	if _, err := r.Body.Read(body); err != nil && err != io.EOF {
		log.Errorf("read http callback body error: %s", err)
		return nil
	}
	log.Debugf("http callback body: %s,len:%d", string(body), len(body))
	log.Debugf("http callback header: %v", r.Header)
	traceID := r.Header.Get(constant.HeaderTraceID)
	// 签名验证
	if pass, err := signature.Verify(credentials.AppSecret, r.Header, body); err != nil || !pass {
		log.Errorf("signature verify failed, err: %v, traceID: %s", err, traceID)
		return nil
	}
	// 解析 payload
	payload := &dto.WSPayload{}
	if err := json.Unmarshal(body, payload); err != nil {
		log.Errorf("unmarshal http callback body error: %s, traceID: %s", err, traceID)
		return nil
	}
	log.Info("payload:%+v", payload)
	// 原始数据放入，parse 的时候需要从里面提取 d
	payload.RawMessage = body
	payload.Session = &dto.Session{AppID: credentials.AppID}
	var result string
	if payload.OPCode == dto.HTTPCallbackValidation {
		data, ok := payload.Data.(map[string]interface{})
		if !ok {
			log.Errorf("callback data invalid: %+v, traceID: %s", payload.Data, traceID)
			return nil
		}
		plainToken, ptOk := data["plain_token"].(string)
		eventTs, etOk := data["event_ts"].(string)
		if !ptOk || !etOk {
			log.Errorf("callback data invalid: %+v, traceID: %s", payload.Data, traceID)
		}
		req := &dto.WHValidationReq{
			PlainToken: plainToken,
			EventTs:    eventTs,
		}
		validationRsp := webhook.GenValidationACK(req, r.Header, credentials.AppSecret)
		if validationRsp != nil {
			if _, err := w.Write(validationRsp); err != nil {
				log.Errorf("write http callback response error: %s, traceID: %s", err, traceID)
				return nil
			}
		}
		return nil
	}

	result = parsePayload(payload, traceID)
	if result != "" {
		if _, err := w.Write([]byte(result)); err != nil {
			log.Errorf("write http callback response error: %s, traceID: %s", err, traceID)
			return nil
		}
	}
	return payload
}

func parsePayload(payload *dto.WSPayload, traceID string) string {
	// 处理心跳包
	if payload.OPCode == dto.WSHeartbeat {
		return webhook.GenHeartbeatACK(uint32(payload.Data.(float64)))
	}
	// 处理事件
	if payload.OPCode == dto.WSDispatchEvent {
		// 解析具体事件，并投递给业务注册的 handler
		if err := event.ParseAndHandle(payload); err != nil {
			log.Errorf(
				"parseAndHandle failed, %v, traceID:%s, payload: %v", err,
				traceID, payload,
			)
			return webhook.GenDispatchACK(false)
		}
		return webhook.GenDispatchACK(true)
	}

	return ""
}
