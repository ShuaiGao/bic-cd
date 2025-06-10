package ding

import (
	"bic-cd/pkg/limiter"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	dingWebhook string
	dingSecret  string
	limiterMap  sync.Map
)

// StopLimiter 退出所有限流器
// 注意：退出所有限流器，未发送的钉钉消息会被丢弃
func StopLimiter() {
	limiterMap.Range(func(key, value interface{}) bool {
		value.(limiter.Limiter).Quit()
		limiterMap.Delete(key)
		return true
	})
}

// SetWebhook
// 当 secret 为空时时，对应的是钉钉机器人【自定义关键词】校验
// 当 secret 不为空时，对应的是【加签】校验
func SetWebhook(webhook, secret string) {
	dingWebhook = webhook
	dingSecret = secret
}

type At struct {
	AtMobiles []string `json:"atMobiles,omitempty"`
	AtUserIds []string `json:"atUserIds,omitempty"`
	IsAtAll   bool     `json:"isAtAll"`
}

func SendDingTalkMarkdown(title, content string, at At) {
	SendDingTalk(&Message{
		MsgType: "markdown",
		Markdown: Markdown{
			Title: title,
			Text:  content,
		},
		At: at,
	})
}

func SendDingTalkText(content string, at At) error {
	return SendDingTalk(&Message{
		MsgType: "text",
		Text: map[string]interface{}{
			"content": content,
		},
		At: at,
	})
}

type DingResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func getUrl(webhook, secret string) (string, error) {
	if webhook == "" {
		return "", errors.New("no webhook")
	}
	if secret == "" {
		return webhook, nil
	}
	timestamp := time.Now().UnixMilli()
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(stringToSign))
	code := mac.Sum(nil)
	code64 := base64.StdEncoding.EncodeToString(code)
	return webhook + fmt.Sprintf("&timestamp=%d&sign=%s", timestamp, code64), nil
}

// SendDingTalk 封装钉钉机器人发消息接口
// 支持加签方式发送
// api doc old: https://open.dingtalk.com/document/orgapp/custom-robot-access#title-zob-eyu-qse
// api doc new: https://open.dingtalk.com/document/orgapp/custom-robots-send-group-messages#6a8e23113eggw
func SendDingTalk(msg *Message) error {
	if msg == nil {
		return errors.New("ding message is nil")
	}
	if msg.webhook == "" {
		msg.webhook = dingWebhook
		msg.secret = dingSecret
	}
	switch msg.MsgType {
	case "markdown":
		if len(msg.Markdown.Text) > 3500 {
			msg.Markdown.Text = msg.Markdown.Text[0:3500]
		}
	case "text":
		if content, ok := msg.Text["content"]; ok {
			if contentStr, ook := content.(string); ook {
				if len(contentStr) > 3500 {
					msg.Text["content"] = contentStr[0:3500]
				}
			}
		}
	}
	getLimiter(msg.webhook).Run(func() {
		_ = sendDingTalk(msg)
	})
	return nil
}

func sendDingTalk(msg *Message) error {
	url, err := getUrl(msg.webhook, msg.secret)
	text, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(text)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("post webhook failed %s %d", resp.Status, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	result := DingResponse{}
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}
	if result.ErrCode != 0 {
		return fmt.Errorf("ding msg response error, code[%d],msg[%s]", result.ErrCode, result.ErrMsg)
	}
	return nil
}

func getLimiter(webhook string) limiter.Limiter {
	key := webhook
	if v, ok := limiterMap.Load(key); ok {
		return v.(limiter.Limiter)
	}
	l := limiter.NewLocalLimiter(0.33)
	limiterMap.Store(key, l)
	return l
}
