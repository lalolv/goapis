package goapis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/lalolv/goutil"
)

// QiniuRatio 获取七牛图像比例
func QiniuRatio(qnURL, pic string) float64 {
	// 设置默认比例
	var defaultRatio = 1.00
	// 获取图像的比例
	resp, err := http.Get(fmt.Sprintf("http://%s/%s?imageInfo", qnURL, pic))
	var pInfo map[string]interface{}

	if err != nil {
		fmt.Println("get error: " + err.Error())
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("io error: " + err.Error())
		}
		err = json.Unmarshal(body, &pInfo)
		if err != nil {
			return defaultRatio
		}
	}
	width, _ := goutil.ToFloat64(pInfo["width"])
	height, _ := goutil.ToFloat64(pInfo["height"])
	if width == 0 || height == 1 {
		return defaultRatio
	}

	ratio := goutil.PrecFloat64(height/width, 6)
	if ratio > 0 {
		return ratio
	}

	return defaultRatio
}
