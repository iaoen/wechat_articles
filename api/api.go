package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/datetime"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/imroc/req/v3"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

func Search(cookie string, query string) (string, error) {
	token := ""
	tokens := strings.Split(cookie, "token=")
	if len(tokens) == 2 {
		token = tokens[1]
	} else {
		return "", errors.New("cookie缺失token")
	}
	c := proxyreq()
	c.SetCommonHeader("cookie", cookie)
	params := map[string]string{
		"action": "search_biz",
		"begin":  "0",
		"count":  "5",
		"query":  query,
		"token":  token,
		"lang":   "zh_CN",
		"f":      "json",
		"ajax":   "1",
	}
	res, err := c.R().SetQueryParams(params).Get("https://mp.weixin.qq.com/cgi-bin/searchbiz")
	if err != nil {
		fmt.Println("searchbiz", err)
		return "", err
	}
	// fmt.Println("search", res.String())
	list := gjson.Get(res.String(), "list").Value()
	listJson, err := convertor.ToJson(list)
	if err != nil {
		fmt.Println("search ToJson", err)
		return "", err
	}
	// fmt.Println(listJson)
	return listJson, nil
}
func Appmsgpublish(cookie, fakeid, page string) (string, error) {
	token := ""
	tokens := strings.Split(cookie, "token=")
	if len(tokens) == 2 {
		token = tokens[1]
	} else {
		return "", errors.New("cookie缺失token")
	}
	c := proxyreq()
	c.SetCommonHeader("cookie", cookie)
	pageInt, err := convertor.ToInt(page)
	if err != nil {
		fmt.Println("appmsg page ToInt", err)
		return "", err
	}
	params := map[string]string{
		"sub":               "list",
		"search_field":      "null",
		"begin":             convertor.ToString((pageInt - 1) * 5),
		"count":             convertor.ToString(pageInt * 5),
		"query":             "",
		"fakeid":            fakeid,
		"type":              "101_1",
		"free_publish_type": "1",
		"sub_action":        "list_ex",
		"token":             token,
		"lang":              "zh_CN",
		"f":                 "json",
		"ajax":              "1",
	}
	// fmt.Println(params)
	res, err := c.R().SetQueryParams(params).Get("https://mp.weixin.qq.com/cgi-bin/appmsgpublish")
	if err != nil {
		fmt.Println("appmsg", err)
		return "", err
	}
	// fmt.Println("appmsg", res.String())
	result := []map[string]string{}
	publish_page := gjson.Get(res.String(), "publish_page").Str
	publish_infos := gjson.Get(publish_page, "publish_list.#.publish_info")
	publish_infos.ForEach(func(key, value gjson.Result) bool {
		appmsgex := gjson.Get(value.String(), "appmsgex")
		appmsgex.ForEach(func(key, value gjson.Result) bool {
			link := value.Get("link").Str
			title := value.Get("title").Str
			time := convertor.ToString(value.Get("update_time").Int())
			// fmt.Println(link, title, time)
			result = append(result, map[string]string{
				"title": title,
				"link":  link,
				"time":  time,
			})
			return true
		})
		return true
	})
	listJson, err := convertor.ToJson(result)
	if err != nil {
		fmt.Println("appmsg ToJson", err)
		return "", err
	}
	// fmt.Println(listJson)
	return listJson, nil
}

func StartLogin() string {
	recount := 0
	uuid := getUUid()
	getQR(uuid)
	cookie := ""
	for {
		recount += 1
		status := getQRStatus(uuid)
		if status == "已扫码" {
			fmt.Println("已扫码, 等待授权.")
		} else if status == "未扫码" {
			fmt.Println("等待扫码中, 30秒失效.")
		} else if status == "已登录" {
			fmt.Println("已登录.")
			cookie = login(uuid)
			break
		} else {
			fmt.Println("未知状态, 请稍后.")
		}
		if recount > 10 {
			fmt.Println("等待扫码超时.")
			break
		}
		time.Sleep(time.Second * 3)
	}
	return cookie
}

func login(uuid string) string {
	c := proxyreq()
	c.AllowGetMethodPayload = false
	c.SetCommonHeader("Cookie", "uuid="+uuid)
	data := map[string]string{
		"userlang":         "zh_CN",
		"redirect_url":     "",
		"cookie_forbidden": "0",
		"cookie_cleaned":   "0",
		"plugin_used":      "0",
		"login_type":       "3",
		"token":            "",
		"lang":             "zh_CN",
		"f":                "json",
		"ajax":             "1",
	}
	res, err := c.R().SetFormData(data).Post("https://mp.weixin.qq.com/cgi-bin/bizlogin?action=login")
	if err != nil {
		fmt.Println("login", err)
		return err.Error()
	}
	result := ""
	cookies := res.Cookies()
	for _, cookie := range cookies {
		if cookie.Value == "EXPIRED" {
			continue
		}
		result += cookie.Name + "=" + cookie.Value + ";"
	}
	redirect_url := gjson.Get(res.String(), "redirect_url").Str
	tokens := strings.Split(redirect_url, "token=")
	if len(tokens) == 2 {
		token := tokens[1]
		result += "token=" + token
	} else {
		return "获取token失败"
	}
	// fmt.Println("cookies", result)
	return result
}

func getQRStatus(uuid string) string {
	c := proxyreq()
	c.SetCommonHeader("Cookie", "uuid="+uuid)
	res, err := c.R().Get("https://mp.weixin.qq.com/cgi-bin/scanloginqrcode?action=ask&token=&lang=zh_CN&f=json&ajax=1")
	if err != nil {
		fmt.Println("getQRStatus", err)
		return err.Error()
	}
	// fmt.Println("getQRStatus", res.String())
	status := gjson.Get(res.String(), "status").Int()
	switch status {
	case 0:
		return "未扫码"
	case 4:
		return "已扫码"
	case 1:
		return "已登录"
	default:
		return "未知状态"
	}
}

func getQR(uuid string) []byte {
	c := proxyreq()
	c.SetCommonHeader("Cookie", "uuid="+uuid)
	res, err := c.R().Get("https://mp.weixin.qq.com/cgi-bin/scanloginqrcode?action=getqrcode&random=" + convertor.ToString(datetime.TimestampMilli()) + "60")
	if err != nil {
		fmt.Println("getSession", err)
		return []byte{}
	}
	fileutil.WriteBytesToFile("qrcode.jpg", res.Bytes())
	return res.Bytes()
}

func getUUid() string {
	c := proxyreq()
	headers := map[string]string{
		"userlang":     "zh_CN",
		"redirect_url": "",
		"login_type":   "3",
		"sessionid":    convertor.ToString(datetime.TimestampMilli()) + "89",
		"token":        "",
		"lang":         "zh_CN",
		"f":            "json",
		"ajax":         "1",
	}
	res, err := c.R().SetFormData(headers).Post("https://mp.weixin.qq.com/cgi-bin/bizlogin?action=startlogin")
	if err != nil {
		fmt.Println("bizlogin getUUid", err)
		return ""
	}
	uuid := ""
	a := res.Cookies()
	for _, cookie := range a {
		if cookie.Name == "uuid" {
			uuid = cookie.Value
			break
		}
	}
	// fmt.Println("uuid", uuid)
	return uuid
}

func proxyreq() *req.Client {
	c := req.C()
	c.SetCommonHeader("origin", "https://mp.weixin.qq.com")
	c.SetCommonHeader("referer", "https://mp.weixin.qq.com/")
	return c
}
