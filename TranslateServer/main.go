package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

/* new关键字根据输入的类型分配一块内存空间并返回类型指针 */
var (
	YoudaoAPI = new(YoudaoTranslateAPI)
	BaiduAPI = new(BaiduTranslateAPI)
	TencentAPI = new(TencentTranslateAPI)
)

/* 调用对应的API并返回翻译结果 */
func translate(api TranslateAPI, method string, text string, from string, to string) string {
	if method == "get" {
		return api.GetMethod(text, from, to)
	} else {
		return api.PostMethod(text, from, to)
	}
}

/* GET请求处理函数 */
func GETHandleFUnc(c *gin.Context) {
	// 获取querystring数据
	translateText := c.Query("text") // 需要翻译的文本
	from := c.DefaultQuery("from", "auto") // 源语言
	to := c.DefaultQuery("to", "auto") // 目标语言
	api := c.DefaultQuery("api", "youdao") // 使用哪个翻译API
	method := c.DefaultQuery("method", "post") // 使用哪种方式向API发送请求

	// 错误的目标语言调用
	var err error = nil
	if translateText == "" {
		err = fmt.Errorf("%s", "参数[text]缺失: 翻译文本不能为空")
		c.JSON(http.StatusOK, gin.H{
			"status": "200",
			"error": fmt.Sprintf("%s", err),
		})
		return
	} else if api == "tencent" && to == "auto" {
		err = fmt.Errorf("%s", "参数[to]缺失: 腾讯翻译API对目标语言不支持auto识别, 请明确目标语言")
		c.JSON(http.StatusOK, gin.H{
			"status": "200",
			"error": fmt.Sprintf("%s", err),
		})
		return
	} else if method != "get" && method != "post" {
		err = fmt.Errorf("%s", "参数[method]错误: 有效的method取值为[get, post]")
		c.JSON(http.StatusOK, gin.H{
			"status": "200",
			"error": fmt.Sprintf("%s", err),
		})
		return
	}

	// 调用API
	var result string
	switch api {
	case "youdao":
		result = translate(YoudaoAPI, method, translateText, from, to)
	case "baidu":
		result = translate(BaiduAPI, method, translateText, from, to)
	case "tencent":
		result = translate(TencentAPI, method, translateText, from, to)
	default:
		err = fmt.Errorf("参数[api]错误: 调用了不存在的API:[%s], 有效的API列表:[baidu, youdao, tencent]", api)
	}
	// 处理错误的API名字
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "200",
			"error": fmt.Sprintf("%s", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "200",
		"api": api,
		"called method": method,
		"src": translateText,
		"dst": result,
	})
}

/* POST请求处理函数 */
func POSTHandleFUnc(c *gin.Context) {
	// 获取querystring数据
	translateText := c.PostForm("text") // 需要翻译的文本
	from := c.DefaultPostForm("from", "auto") // 源语言
	to := c.DefaultPostForm("to", "auto") // 目标语言
	api := c.DefaultPostForm("api", "youdao") // 使用哪个翻译API
	method := c.DefaultPostForm("method", "post") // 使用哪种方式向API发送请求

	// 错误处理
	var err error = nil
	if translateText == "" {
		err = fmt.Errorf("%s", "参数[text]缺失: 翻译文本不能为空")
		c.JSON(http.StatusOK, gin.H{
			"status": "200",
			"error": fmt.Sprintf("%s", err),
		})
		return
	} else if api == "tencent" && to == "auto" {
		err = fmt.Errorf("%s", "参数[to]缺失: 腾讯翻译API对目标语言不支持auto识别, 请明确目标语言")
		c.JSON(http.StatusOK, gin.H{
			"status": "200",
			"error": fmt.Sprintf("%s", err),
		})
		return
	} else if method != "get" && method != "post" {
		err = fmt.Errorf("%s", "参数[method]错误: 有效的method取值为[get, post]")
		c.JSON(http.StatusOK, gin.H{
			"status": "200",
			"error": fmt.Sprintf("%s", err),
		})
		return
	}

	// 调用API
	var result string
	switch api {
	case "youdao":
		result = translate(YoudaoAPI, method, translateText, from, to)
	case "baidu":
		result = translate(BaiduAPI, method, translateText, from, to)
	case "tencent":
		result = translate(TencentAPI, method, translateText, from, to)
	default:
		err = fmt.Errorf("参数[api]错误: 调用了不存在的API[%s], 有效的API列表:[baidu, youdao, tencent]", api)
	}
	// 处理错误的API名字
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "200",
			"error": fmt.Sprintf("%s", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "200",
		"api": api,
		"called method": method,
		"src": translateText,
		"dst": result,
	})
}

func main() {
	router := gin.Default()

	// 这个是访问服务器的方式，不是调用API的方式
	router.GET("/translate", GETHandleFUnc)
	router.POST("/translate", POSTHandleFUnc)
	err := router.Run()
	if err != nil {
		fmt.Printf("%v", err)
	}
}
