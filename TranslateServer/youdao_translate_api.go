/*
@Title youdaotranslate.go
@Description 有道翻译API调用，封装了GET和POST方法
*/

package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	YOUDAO_URL    = "https://openapi.youdao.com/api"
	YOUDAO_APPID  = "0c48af489fab0a9b"
	YOUDAO_APPKEY = "GDStDZ7SOxMAEANgVOlDeN5vWgKYtP0x"
)

/* 有道翻译API */
type YoudaoTranslateAPI struct {
	YoudaoRequestData
	ResponseData YoudaoResponseData
}

type YoudaoRequestData struct {
	From string
	To string
	SignType string
	Curtime string
	Appkey string
	Q string
	Sign string
	Salt string
}

/* 有道翻译API返回的Response参数, 注释掉的字段在返回值中不一定存在 */
type YoudaoResponseData struct {
	ErrorCode    string      `json:"errorCode"`
	Query        string      `json:"query"`
	Translation  []string    `json:"translation"`
	Basic        interface{} `json:"basic,omitempty"`
	Web          interface{} `json:"web,omitempty"`
	L            string      `json:"l"`
	Dict         interface{} `json:"dict,omitempty"`
	Webdict      interface{} `json:"webdict,omitempty"`
	TSpeakUrl    interface{} `json:"tSpeakUrl,omitempty"`
	SpeakUrl     interface{} `json:"speakUrl,omitempty"`
	ReturnPhrase []string    `json:"returnPhrase"`
}

/* 将翻译文本按照文档要求进行切割 */
func truncate(q string) string {
	text := []rune(q)
	size := len(text)
	if size <= 20 {
		return q
	} else {
		return string(text[0:10]) + strconv.Itoa(size) + string(text[size-10:size])
	}
}

/* 对签名进行SHA256加密并转换为16进制字符串 */
func encrypt(sign string) string {
	// 返回16进制字符串数据
	return fmt.Sprintf("%x", sha256.Sum256([]byte(sign)))
}

func (t *YoudaoTranslateAPI) SignMaker() {
	signStr := YOUDAO_APPID + truncate(t.Q) + t.Salt + t.Curtime + YOUDAO_APPKEY
	t.Sign = encrypt(signStr) // encrypt函数对字符串进行SHA256加密
}

/* 制作访问有道API的参数列表 */
func (t *YoudaoTranslateAPI) PostRequestParams(text string, from string, to string) url.Values {
	// Unix方法返回时间戳，int64类型，不能直接int()
	t.Curtime = strconv.FormatInt(time.Now().Unix(), 10)
	// golang官方库中没有uuid相关的，需要使用第三方库; NewUUID返回UUID类型，实际上是[16]byte类型的别名
	id, err := uuid.NewUUID()
	if err != nil {
		fmt.Printf("Get UUID failed, %v!\n", err)
	}
	t.Salt = id.String()
	// 其他参数
	t.From = from
	t.To = to
	t.SignType = "v3"
	t.Appkey = YOUDAO_APPID
	t.Q = text
	// 最后制作签名
	t.SignMaker()

	// request参数
	params := url.Values{
		"from":     {t.From},
		"to":       {t.To},
		"signType": {t.SignType},
		"curtime":  {t.Curtime},
		"appKey":   {t.Appkey},
		"q":        {t.Q},
		"sign":     {t.Sign},
		"salt":     {t.Salt},
	}

	return params
}

func (t *YoudaoTranslateAPI) PostMethod(q string, from string, to string) string {
	// 获取请求参数
	params := t.PostRequestParams(q, from, to)
	// 发送请求
	resp, err := http.PostForm(YOUDAO_URL, params)
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

	fmt.Printf("Youdao API called successfully through POST method!\n")

	return t.ResponseData.Translation[0]
}

func (t *YoudaoTranslateAPI) GetRequestUrl(text string, from string, to string) string {
	// Unix方法返回时间戳，int64类型，不能直接int()
	t.Curtime = strconv.FormatInt(time.Now().Unix(), 10)
	// golang官方库中没有uuid相关的，需要使用第三方库; NewUUID返回UUID类型，实际上是[16]byte类型的别名
	id, err := uuid.NewUUID()
	if err != nil {
		fmt.Printf("Get UUID failed, %v!\n", err)
	}
	t.Salt = id.String()
	// 其他参数
	t.From = from
	t.To = to
	t.SignType = "v3"
	t.Appkey = YOUDAO_APPID
	t.Q = text
	// 最后制作签名
	t.SignMaker()

	// URL
	requestUrl := YOUDAO_URL
	requestUrl += "?from=" + t.From
	requestUrl += "&to=" + t.To
	requestUrl += "&signType=" + t.SignType
	requestUrl += "&curtime=" + t.Curtime
	requestUrl += "&appKey=" + t.Appkey
	requestUrl += "&q=" + url.QueryEscape(t.Q) // 拼接url的时候，翻译的文本需要进行urlencode处理
	requestUrl += "&sign=" + t.Sign
	requestUrl += "&salt=" + t.Salt

	return requestUrl
}

func (t *YoudaoTranslateAPI) GetMethod(text string, from string, to string) string {
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

	fmt.Printf("Youdao API called successfully through GET method!\n")

	return t.ResponseData.Translation[0]
}
