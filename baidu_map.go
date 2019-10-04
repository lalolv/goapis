package goapis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// BaiduCity 根据经纬度和IP地址，获取省份和城市
func BaiduCity(ak, remoteAddr string, lat, lng float64) (string, string, string) {
	// 经纬度类型
	coor := "bd09ll"
	// 返回变量
	prov, city, dist := "", "", ""

	var ret map[string]interface{}
	// 根据经纬度
	if lat > 0 && lng > 0 {
		url := fmt.Sprintf(
			"http://api.map.baidu.com/geocoder/v2/?ak=%s&location=%v,%v&output=json",
			ak, lat, lng)
		resp, _ := http.Get(url)
		if resp != nil {
			body, _ := ioutil.ReadAll(resp.Body)
			err := json.Unmarshal(body, &ret)
			if err == nil && ret != nil && ret["result"] != nil && ret["status"].(float64) == 0 {
				result := ret["result"].(map[string]interface{})
				addr := result["addressComponent"].(map[string]interface{})
				prov = addr["province"].(string)
				city = addr["city"].(string)
				dist = addr["district"].(string)
			}
			resp.Body.Close()
		}
	}

	// 根据IP定位位置
	if remoteAddr != "" {
		end := strings.LastIndex(remoteAddr, ":")
		ipaddr := remoteAddr[:end]
		url := fmt.Sprintf(
			"http://api.map.baidu.com/location/ip?ak=%s&ip=%s&coor=%s",
			ak, ipaddr, coor)
		resp, _ := http.Get(url)
		if resp != nil {
			body, _ := ioutil.ReadAll(resp.Body)
			err := json.Unmarshal(body, &ret)
			if err == nil && ret != nil && ret["status"] != nil && ret["status"].(float64) == 0 {
				content := ret["content"].(map[string]interface{})
				addr := content["address_detail"].(map[string]interface{})
				prov = addr["province"].(string)
				city = addr["city"].(string)
				dist = addr["district"].(string)
			}
			resp.Body.Close()
		}
	}

	return prov, city, dist
}

// BaiduLoc 根据城市名称，解析经纬度
func BaiduLoc(ak, city string) (float64, float64) {
	// 初始化经纬度
	lat, lng := 0.00, 0.00

	var ret map[string]interface{}
	// 通过百度接口，获取结果
	url := fmt.Sprintf(
		"http://api.map.baidu.com/geocoder/v2/?address=%s&output=json&ak=%s",
		city, ak)
	resp, _ := http.Get(url)
	if resp != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal(body, &ret)
		if err == nil && ret != nil && ret["result"] != nil && ret["status"].(float64) == 0 {
			result := ret["result"].(map[string]interface{})
			loc := result["location"].(map[string]interface{})
			lat = loc["lat"].(float64)
			lng = loc["lng"].(float64)
		}
		resp.Body.Close()
	}

	return lat, lng
}
