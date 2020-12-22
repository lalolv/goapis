package goapis

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

// SlackUtil Slack 接口
type SlackUtil struct {
	Token string
}

// PostMessage 发送消息给频道
func (slack *SlackUtil) PostMessage(channelID, text string) error {
	url := "https://slack.com/api/chat.postMessage"
	// 消息体
	buf, _ := json.Marshal(bson.M{
		"text":    text,
		"channel": channelID,
	})
	body := bytes.NewBuffer(buf)
	req, err := http.NewRequest("POST", url, body)
	// 头部认证
	auth := fmt.Sprintf("Bearer %s", slack.Token)
	// Header
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Authorization", auth)
	// resp
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	data := bson.M{}
	// get data
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			r, _ := ioutil.ReadAll(resp.Body)
			json.Unmarshal(r, &data)
		} else {
			return errors.New(resp.Status)
		}
	} else {
		return errors.New("Req error")
	}

	return nil
}
