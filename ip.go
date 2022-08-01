package goapis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// IPAddressInfo IP 地址信息
type IPAddressInfo struct {
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	Region      string `json:"region"`
	RegionName  string `json:"regionName"`
	City        string `json:"city"`
	ISP         string `json:"isp"`
}

// IP2Adress IP地址转换为物理地址
// ip-api.com
// http://ip-api.com/json/113.116.28.98?lang=zh-CN
func IP2Adress(ip string) (IPAddressInfo, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s?lang=zh-CN", ip)
	// fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("ip error: " + err.Error())
		return IPAddressInfo{}, err
	}

	// 读取数据资料
	var uInfo IPAddressInfo
	if resp != nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("ip api io error: " + err.Error())
			return IPAddressInfo{}, err
		}

		err = json.Unmarshal(body, &uInfo)
		if err != nil {
			fmt.Println("解析失败: " + err.Error())
			return IPAddressInfo{}, err
		}
	}

	return uInfo, nil
}
