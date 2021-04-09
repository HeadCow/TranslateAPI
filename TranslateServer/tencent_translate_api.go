package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"time"
)

var (
	TENCENT_URL    = "https://tmt.tencentcloudapi.com/"
	TENCENT_APPID  = ""
	TENCENT_APPKEY = ""
)

type TencentTranslateAPI struct {
	// 请求方法
	RequestMethod string
	// 请求参数，因为要排序所以不能匿名
	RequestData TencentRequestData
	// 响应参数
	ResponseData TencentResponseData
}

type TencentRequestData struct {
	Action     string
	Region     string
	Timestamp  string
	Nonce      string
	SecretId   string
	Signature  string
	Version    string
	SourceText string
	Source     string
	Target     string
	ProjectId  string
}

/* 腾讯翻译API返回的Response参数，Response嵌套其他数据 */
type TencentResponseData struct {
	Response map[string]interface{} `json:"Response"`
}

func (t *TencentTranslateAPI) SignatureMaker() {
	params := make(map[string]interface{})
	var keyStr []string
	// 通过反射将struct转换为map
	keys := reflect.TypeOf(t.RequestData)
	values := reflect.ValueOf(t.RequestData)
	for i := 0; i < keys.NumField(); i++ {
		if keys.Field(i).Name == "Signature" {
			continue
		}
		params[keys.Field(i).Name] = values.Field(i).Interface()
		keyStr = append(keyStr, keys.Field(i).Name)
	}
	// 按key排序
	sort.Strings(keyStr)
	// 拼接参数
	var paramList string
	for _, v := range keyStr {
		paramList += v + "=" + params[v].(string) + "&"
	}
	l := len(paramList)
	paramList = paramList[:l-1]
	// 拼接请求字符串，规则为：请求方法 + 请求主机 +请求路径 + ? + 请求字符串
	signURL := t.RequestMethod + "tmt.tencentcloudapi.com/?" + paramList
	// 对上述签名原文进行HMAC-SHA1加密，并使用Base64编码，得到最终的签名
	mac := hmac.New(sha1.New, []byte(TENCENT_APPKEY))
	mac.Write([]byte(signURL))
	// 获得最终的签名
	t.RequestData.Signature = base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

/* 制作访问腾讯云API的参数列表 */
func (t *TencentTranslateAPI) PostRequestParams(text string, from string, to string) url.Values {
	t.RequestMethod = "POST"
	// Unix时间戳，int64类型，需要转换为int类型
	t.RequestData.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	// 随机正整数
	t.RequestData.Nonce = strconv.Itoa(rand.Intn(10000))
	// 其他参数
	t.RequestData.Action = "TextTranslate"
	t.RequestData.Region = "ap-guangzhou"
	t.RequestData.SecretId = TENCENT_APPID
	t.RequestData.Version = "2018-03-21"
	t.RequestData.SourceText = text
	t.RequestData.Source = from
	t.RequestData.Target = to
	t.RequestData.ProjectId = "0"

	// 制作签名
	t.SignatureMaker()

	// request参数
	params := url.Values{
		"Action":     {t.RequestData.Action},
		"Region":     {t.RequestData.Region},
		"Timestamp":  {t.RequestData.Timestamp},
		"Nonce":      {t.RequestData.Nonce},
		"SecretId":   {t.RequestData.SecretId},
		"Signature":  {t.RequestData.Signature},
		"Version":    {t.RequestData.Version},
		"SourceText": {t.RequestData.SourceText},
		"Source":     {t.RequestData.Source},
		"Target":     {t.RequestData.Target},
		"ProjectId":  {t.RequestData.ProjectId},
	}

	return params
}

/* POST方法 */
func (t *TencentTranslateAPI) PostMethod(q string, from string, to string) string {
	// 获取请求参数
	params := t.PostRequestParams(q, from, to)
	// 发送请求
	// 特别注意：由于PostForm方法内部已经实现参数的URLEncode编码，不要自己进行URLEncode，否则会失效！
	resp, err := http.PostForm(TENCENT_URL, params)
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

	fmt.Printf("Tencent API called successfully through POST method!\n")

	return t.ResponseData.Response["TargetText"].(string)
}

func (t *TencentTranslateAPI) GetRequestUrl(text string, from string, to string) string {
	t.RequestMethod = "GET"
	// Unix时间戳，int64类型，需要转换为int类型
	t.RequestData.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	// 随机正整数
	t.RequestData.Nonce = strconv.Itoa(rand.Intn(10000))
	// 其他参数
	t.RequestData.Action = "TextTranslate"
	t.RequestData.Region = "ap-guangzhou"
	t.RequestData.SecretId = TENCENT_APPID
	t.RequestData.Version = "2018-03-21"
	t.RequestData.SourceText = text
	t.RequestData.Source = from
	t.RequestData.Target = to
	t.RequestData.ProjectId = "0"

	// 制作签名
	t.SignatureMaker()

	requestUrl := TENCENT_URL
	requestUrl += "?Action=" + t.RequestData.Action
	requestUrl += "&Region=" + t.RequestData.Region
	requestUrl += "&Timestamp=" + t.RequestData.Timestamp
	requestUrl += "&Nonce=" + t.RequestData.Nonce
	requestUrl += "&SecretId=" + t.RequestData.SecretId
	requestUrl += "&Signature=" + url.QueryEscape(t.RequestData.Signature) // 签名需要urlencode编码
	requestUrl += "&Version=" + t.RequestData.Version
	requestUrl += "&SourceText=" + url.QueryEscape(t.RequestData.SourceText) // 翻译文本需要urlencode编码
	requestUrl += "&Source=" + t.RequestData.Source
	requestUrl += "&Target=" + t.RequestData.Target
	requestUrl += "&ProjectId=" + t.RequestData.ProjectId

	return requestUrl
}

/* GET方法，特别注意签名和翻译文本需要进行URLEncode编码 */
func (t *TencentTranslateAPI) GetMethod(text string, from string, to string) string {
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

	fmt.Printf("Tentent API called successfully through GET method!\n")

	return t.ResponseData.Response["TargetText"].(string)
}
