package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	BAIDU_URL    = "https://fanyi-api.baidu.com/api/trans/vip/translate"
	BAIDU_APPID  = "20210403000760696"
	BAIDU_APPKEY = "cXX2Q1PSF7HsFeArJLaR"
)

type BaiduTranslateAPI struct {
	BaiduRequestData
	ResponseData BaiduResponseData
}

type BaiduRequestData struct {
	From string
	To string
	Appid string
	Q string
	Sign string
	Salt string
}

/* 百度翻译API返回的Response参数 */
type BaiduResponseData struct {
	From         string        `json:"from"`
	To           string        `json:"to"`
	Trans_result []interface{} `json:"trans_result"`
	Error_code   int32         `json:"error_code,omitempty"`
}

func (t *BaiduTranslateAPI) SignMaker() {
	signStr := []byte(BAIDU_APPID + t.Q + t.Salt + BAIDU_APPKEY)
	t.Sign = fmt.Sprintf("%x", md5.Sum(signStr)) // 对字符串进行MD5加密，并转化为16进制字符串
}

/* 制作访问百度API的参数列表 */
func (t *BaiduTranslateAPI) PostRequestParams(text string, from string, to string) url.Values {
	// golang官方库中没有uuid相关的，需要使用第三方库; NewUUID返回UUID类型，实际上是[16]byte类型的别名
	id, err := uuid.NewUUID()
	if err != nil {
		fmt.Printf("Get UUID failed, %v!\n", err)
	}
	t.Salt = id.String()
	// 其他参数
	t.From = from
	t.To = to
	t.Appid = BAIDU_APPID
	t.Q = text
	// 制作签名
	t.SignMaker()

	// request参数
	params := url.Values{
		"from":  {t.From},
		"to":    {t.To},
		"appid": {t.Appid},
		"q":     {t.Q},
		"sign":  {t.Sign},
		"salt":  {t.Salt},
	}

	return params
}

func (t *BaiduTranslateAPI) PostMethod(q string, from string, to string) string {
	// 获取请求参数
	params := t.PostRequestParams(q, from, to)
	// 发送请求
	resp, err := http.PostForm(BAIDU_URL, params)
	if err != nil {
		fmt.Printf("Http request error, %v!\n", err)
	}
	// 响应Body的反序列化
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Response body read failed, %v!\n", err)
	}
	err = json.Unmarshal(body, &(t.ResponseData))
	if err != nil {
		fmt.Printf("Response body unmarshal failed, %v!\n", err)
	}

	fmt.Printf("Baidu API called successfully through POST method!\n")

	return t.ResponseData.Trans_result[0].(map[string]interface{})["dst"].(string)
}

func (t *BaiduTranslateAPI) GetRequestUrl(text string, from string, to string) string {
	// golang官方库中没有uuid相关的，需要使用第三方库; NewUUID返回UUID类型，实际上是[16]byte类型的别名
	id, err := uuid.NewUUID()
	if err != nil {
		fmt.Printf("Get UUID failed, %v!\n", err)
	}
	t.Salt = id.String()
	// 其他参数
	t.From = from
	t.To = to
	t.Appid = BAIDU_APPID
	t.Q = text
	// 制作签名
	t.SignMaker()

	// URL
	requestUrl := BAIDU_URL
	requestUrl += "?from=" + t.From
	requestUrl += "&to=" + t.To
	requestUrl += "&appid=" + t.Appid
	requestUrl += "&q=" + url.QueryEscape(t.Q) // 拼接url的时候，翻译的文本需要进行urlencode处理
	requestUrl += "&sign=" + t.Sign
	requestUrl += "&salt=" + t.Salt

	return requestUrl
}

func (t *BaiduTranslateAPI) GetMethod(text string, from string, to string) string {
	// 获取Get请求的URL
	requestUrl := t.GetRequestUrl(text, from, to)
	// 发送请求
	resp, err := http.Get(requestUrl)
	if err != nil {
		fmt.Printf("Http request error, %v!\n", err)
	}
	// 响应Body的反序列化
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Response body read failed, %v!\n", err)
	}
	err = json.Unmarshal(body, &(t.ResponseData))
	if err != nil {
		fmt.Printf("Response body unmarshal failed, %v!\n", err)
	}

	fmt.Printf("Baidu API called successfully through GET method!\n")

	return t.ResponseData.Trans_result[0].(map[string]interface{})["dst"].(string)
}
