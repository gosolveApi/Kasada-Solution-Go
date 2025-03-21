package demo

import (
	"encoding/base64"
	"encoding/json"
	"github.com/antchfx/htmlquery"
	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
	"log"
	"strconv"
	"strings"
	"time"
)

type footlockerTest struct {
	client  *req.Client
	proxy   string
	pageurl string
}

func Test() {
	f := new(footlockerTest)
	f.proxy = ""
	f.Step1() //init client
}

func (a *footlockerTest) Step1() {
	a.client = req.C()
	a.client.EnableAutoDecompress()
	a.client.ImpersonateChrome()
	a.client.SetProxyURL(a.proxy)
	urlstr := "https://www.footlocker.com/149e9513-01fa-4fb0-aad4-566afd725d1b/2d206a39-8ed7-437e-a3be-862e0f06eea3/fp?x-kpsdk-v=j-1.0.0"
	headers := map[string]string{
		"sec-ch-ua":                 `"Chromium";v="134", "Not:A-Brand";v="24", "Google Chrome";v="134"`,
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        `"Windows"`,
		"upgrade-insecure-requests": "1",
		"user-agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36",
		"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"sec-fetch-site":            "same-origin",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-dest":            "iframe",
		"referer":                   "https://www.footlocker.com/",
		"accept-encoding":           "gzip, deflate, br, zstd",
		"accept-language":           "zh-CN,zh;q=0.9",
		"priority":                  "u=0, i",
	}
	headerorder := []string{
		"sec-ch-ua",
		"sec-ch-ua-mobile",
		"sec-ch-ua-platform",
		"upgrade-insecure-requests",
		"user-agent",
		"accept",
		"sec-fetch-site",
		"sec-fetch-mode",
		"sec-fetch-dest",
		"referer",
		"accept-encoding",
		"accept-language",
		"priority",
	}
	get, err := a.client.R().SetHeaders(headers).SetHeaderOrder(headerorder...).Get(urlstr)

	if err != nil {
		log.Println(err)
	} else {
		if get.StatusCode == 429 {
			parse, err := htmlquery.Parse(get.Body)
			if err != nil {
				log.Println(err)
			} else {
				find := htmlquery.Find(parse, "//script[@src]")
				for _, node := range find {
					for _, attr := range node.Attr {
						if attr.Key == "src" {
							a.pageurl = "https://www.footlocker.com" + attr.Val
							break
						}
					}
				}
				if a.pageurl != "" {

					a.Step2()
				}
			}
		}
	}

}
func (a *footlockerTest) Step2() { //get JavaScriptSource
	urlstr := a.pageurl
	headers := map[string]string{
		"sec-ch-ua-platform": `"Windows"`,
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36",
		"sec-ch-ua":          `"Chromium";v="134", "Not:A-Brand";v="24", "Google Chrome";v="134"`,
		"sec-ch-ua-mobile":   "?0",
		"accept":             "*/*",
		"sec-fetch-site":     "same-origin",
		"sec-fetch-mode":     "no-cors",
		"sec-fetch-dest":     "script",
		"referer":            "https://www.footlocker.com/149e9513-01fa-4fb0-aad4-566afd725d1b/2d206a39-8ed7-437e-a3be-862e0f06eea3/fp?x-kpsdk-v=j-1.0.0",
		"accept-encoding":    "gzip, deflate, br, zstd",
		"accept-language":    "zh-CN,zh;q=0.9",
		"priority":           "u=1",
	}
	headerorder := []string{
		"sec-ch-ua-platform",
		"user-agent",
		"sec-ch-ua",
		"sec-ch-ua-mobile",
		"accept",
		"sec-fetch-site",
		"sec-fetch-mode",
		"sec-fetch-dest",
		"referer",
		"accept-encoding",
		"accept-language",
		"priority",
	}

	get, err := a.client.R().SetHeaders(headers).SetHeaderOrder(headerorder...).Get(urlstr)
	if err != nil {
		log.Println(err)
	} else {
		if strings.Contains(get.String(), "KPSDK.scriptStart") {
			s := genCt(a.pageurl, get.String())
			if s != "" {
				a.Step3(s)
			}

		}
	}
}
func (a *footlockerTest) Step3(str string) { //post to tl
	urlstr := "https://www.footlocker.com/149e9513-01fa-4fb0-aad4-566afd725d1b/2d206a39-8ed7-437e-a3be-862e0f06eea3/tl"
	headers := map[string]string{
		"x-kpsdk-ct":         "0LRzucs7TkyttLWFyL8UII0d9KTjW6yn6nqeq3ciNmyibnUSnlXK2IVRWyzsv9aTLxRKq1zqXS4qzpXMCWbILgkZjonsKv5K2wdJ4DxeKgwFALgotggfP7Eq1wJ3TZECBuerJQgI1DQjf6TEdHiVKaDXcOw0FTJVPxdSomI5",
		"sec-ch-ua-platform": `"Windows"`,
		"x-kpsdk-dt":         "1b2x41mw5bzcjay32w76x27w1t4z6dw82zdc4y01how9hlya0",
		"sec-ch-ua":          `"Chromium";v="134", "Not:A-Brand";v="24", "Google Chrome";v="134"`,
		"x-kpsdk-im":         "CiRlNThhOTYzNS1mMmJmLTQ2NzUtYjU3YS02MjFjNmE0ODQ1NGQ",
		"sec-ch-ua-mobile":   "?0",
		"x-kpsdk-v":          "j-1.0.0",
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36",
		"content-type":       "application/octet-stream",
		"accept":             "*/*",
		"origin":             "https://www.footlocker.com",
		"sec-fetch-site":     "same-origin",
		"sec-fetch-mode":     "cors",
		"sec-fetch-dest":     "empty",
		"referer":            "https://www.footlocker.com/149e9513-01fa-4fb0-aad4-566afd725d1b/2d206a39-8ed7-437e-a3be-862e0f06eea3/fp?x-kpsdk-v=j-1.0.0",
		"accept-encoding":    "gzip, deflate, br, zstd",
		"accept-language":    "zh-CN,zh;q=0.9",
		"priority":           "u=1, i",
	}
	headerorder := []string{
		"x-kpsdk-ct",
		"sec-ch-ua-platform",
		"x-kpsdk-dt",
		"sec-ch-ua",
		"x-kpsdk-im",
		"sec-ch-ua-mobile",
		"x-kpsdk-v",
		"user-agent",
		"content-type",
		"accept",
		"origin",
		"sec-fetch-site",
		"sec-fetch-mode",
		"sec-fetch-dest",
		"referer",
		"accept-encoding",
		"accept-language",
		"priority",
	}
	parse := gjson.Parse(str)
	parse.Get("headers").ForEach(func(key, value gjson.Result) bool {
		value.ForEach(func(k, v gjson.Result) bool {
			headers[k.String()] = v.String()
			return true
		})
		return true
	})
	payload := parse.Get("payload").String()
	body, _ := base64.StdEncoding.DecodeString(payload)
	post, err := a.client.R().SetHeaders(headers).SetHeaderOrder(headerorder...).SetBodyBytes(body).Post(urlstr)
	if err != nil {
		log.Println(err)
	} else {
		if post.String() == `{"reload":true}` {
			log.Println("gen ct success")
			ctstr := post.GetHeader("x-kpsdk-ct")
			a.status(ctstr, post.GetHeader("x-kpsdk-st"))
		}
	}
}

func (a *footlockerTest) status(ct, cd string) { //verify
	s := genCd(cd)

	urlstr := "https://www.footlocker.com/zgw/accounts-experience/v1/accounts/status"
	headers := map[string]string{
		"x-fl-request-id":    "031f6d40-039f-11f0-8958-d13edc1e6920",
		"x-kpsdk-ct":         ct,
		"sec-ch-ua-platform": `"Windows"`,
		"sec-ch-ua":          `"Chromium";v="134", "Not:A-Brand";v="24", "Google Chrome";v="134"`,
		"sec-ch-ua-mobile":   "?0",
		"x-api-lang":         "en-US",
		"x-kpsdk-v":          "j-1.0.0",
		"x-kpsdk-cd":         s,
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36",
		"accept":             "application/json",
		"content-type":       "application/json",
		"origin":             "https://www.footlocker.com",
		"sec-fetch-site":     "same-origin",
		"sec-fetch-mode":     "cors",
		"sec-fetch-dest":     "empty",
		"referer":            "https://www.footlocker.com/",
		"accept-encoding":    "gzip, deflate, br, zstd",
		"accept-language":    "zh-CN,zh;q=0.9",
		"priority":           "u=1, i",
	}
	headerorder := []string{
		"x-fl-request-id",
		"x-kpsdk-ct",
		"sec-ch-ua-platform",
		"sec-ch-ua",
		"sec-ch-ua-mobile",
		"x-api-lang",
		"x-kpsdk-v",
		"x-kpsdk-cd",
		"user-agent",
		"accept",
		"content-type",
		"origin",
		"sec-fetch-site",
		"sec-fetch-mode",
		"sec-fetch-dest",
		"referer",
		"accept-encoding",
		"accept-language",
		"priority",
	}
	body := `{"email":"adfghrsd111fqwr@hotmail.com"}`

	post, _ := a.client.R().SetHeaders(headers).SetHeaderOrder(headerorder...).SetBodyJsonString(body).Post(urlstr)
	log.Println(post.StatusCode)

}

type ct struct {
	Item      string `json:"item"`
	Ver       string `json:"ver"`
	ClientKey string `json:"clientKey"`
	Task      struct {
		ScriptUrl    string `json:"script_url"`
		Ua           string `json:"ua"`
		ScriptSource string `json:"script_source"`
		Lang         string `json:"lang"`
	} `json:"task"`
}
type cd struct {
	Item      string `json:"item"`
	Ver       string `json:"ver"`
	ClientKey string `json:"clientKey"`
	Task      struct {
		St       string `json:"st"`
		WorkTime string `json:"workTime"`
	} `json:"task"`
}

func genCt(ScriptUrl, ScriptSource string) string {
	c := new(ct)
	c.Item = "kasada"
	c.Ver = "ct"
	c.ClientKey = "xxx-xxx-xxx-xxx-xxx"
	c.Task.Ua = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36"
	c.Task.Lang = "zh-CN"
	c.Task.ScriptUrl = ScriptUrl
	c.Task.ScriptSource = ScriptSource
	marshal, err := json.Marshal(c)
	payload := ""
	if err != nil {
		log.Println(err)
	} else {
		client := req.C()
		post, err := client.R().SetBodyJsonString(string(marshal)).Post("https://api.gosolve.cc/resolve")
		if err != nil {
			log.Println(err)
		} else {
			parse := gjson.Parse(post.String())
			payload = parse.Get("payload").String()
		}
	}
	return payload
}
func genCd(st string) string {
	c := new(cd)
	c.Item = "kasada"
	c.Ver = "cd"
	c.ClientKey = "xxx-xxx-xxx-xxx-xxx"
	c.Task.St = st
	c.Task.WorkTime = strconv.FormatInt(time.Now().UnixMilli()+300, 10)
	marshal, err := json.Marshal(c)
	payload := ""
	if err != nil {
		log.Println(err)
	} else {
		client := req.C()
		log.Println(string(marshal))
		post, err := client.R().SetBodyJsonString(string(marshal)).Post("https://api.gosolve.cc/resolve")
		if err != nil {
			log.Println(err)
		} else {
			log.Println(post.String())
			parse := gjson.Parse(post.String())
			payload = parse.Get("payload").String()
		}
	}
	return payload
}
