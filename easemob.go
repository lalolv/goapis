package goapis

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/lalolv/goutil"
)

// Easemob 环信
type Easemob struct {
	ID     string
	Secret string
}

// ERRSTR 错误代码字符串
const ERRSTR string = "error"

// 生成环信token
func (p *Easemob) generateToken() bson.M {
	// body params
	params := bson.M{
		"grant_type":    "client_credentials",
		"client_id":     p.ID,
		"client_secret": p.Secret,
	}
	resp := p.doRequest(params, "https://a1.easemob.com/fashionmii/fashionmii/token", "", "POST")
	if resp["errcode"] == 0 {
		return resp["data"].(bson.M)
	}

	return bson.M{"access_token": "error", "expires_in": -1}
}

// 执行post请求，返回json
// 执行请求，返回json数据
func (p *Easemob) doRequest(requestData bson.M, requestURL, aToken, method string) bson.M {
	buf, _ := json.Marshal(requestData)
	body := bytes.NewBuffer(buf)
	req, err := http.NewRequest(method, requestURL, body)
	if err != nil {
		fmt.Println("请求错误:" + err.Error())
	}
	// Header
	req.Header.Add("Content-Type", "application/json")
	if aToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", aToken))
	}

	// resp
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	data := bson.M{}
	errcode := 0
	// get data
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			r, _ := ioutil.ReadAll(resp.Body)
			err = json.Unmarshal(r, &data)
			if err != nil {
				errcode = 2
			}
		} else {
			data = bson.M{"HTTP_CODE": resp.StatusCode}
		}
	} else {
		errcode = 1
	}
	return bson.M{"data": data, "errcode": errcode}
}

// EasemobSignupSingle 环信注册单个用户
// @uid
// @passwd
// @nickName
func (p *Easemob) EasemobSignupSingle(uid int64, passwd, nickName string) bool {
	// body params
	params := bson.M{
		"username": uid,
		"password": passwd,
		"nickname": nickName,
	}
	// 生成token
	aToken := p.generateToken()
	if aToken["access_token"] == ERRSTR {
		fmt.Println("生成token失败")
		return false
	}

	data := p.doRequest(
		params,
		"https://a1.easemob.com/fashionmii/fashionmii/users",
		aToken["access_token"].(string), "POST",
	)

	errcode, _ := goutil.ToInt(data["errcode"])

	return errcode == 0
}

// EasemobNickModify 环信修改用户昵称
// @userName
// @passwd
// @nickName
// @aToken
func (p *Easemob) EasemobNickModify(uid int64, nickName string) bool {
	// body params
	params := bson.M{
		"nickname": nickName,
	}

	aToken := p.generateToken()
	if aToken["access_token"] == ERRSTR {
		fmt.Println("生成token失败")
		return false
	}

	data := p.doRequest(
		params,
		fmt.Sprintf("https://a1.easemob.com/fashionmii/fashionmii/users/%d", uid),
		aToken["access_token"].(string), "PUT")

	errcode, _ := goutil.ToInt(data["errcode"])

	return errcode == 0
}

// PushMessage 发送消息
func (p *Easemob) PushMessage(reqURL, msg string, target map[string]string) (interface{}, error) {
	// body params
	params := bson.M{
		"target_type": "users",
		"target":      target,
		"msg": map[string]string{
			"type": "txt",
			"msg":  msg,
		},
		"from": "systerm",
	}

	aToken := p.generateToken()
	if aToken["access_token"] == ERRSTR {
		fmt.Println("生成token失败")
		return nil, errors.New("生成token失败")
	}

	data := p.doRequest(params, reqURL, aToken["access_token"].(string), "POST")

	return data, nil
}
