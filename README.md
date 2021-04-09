# TranslateAPI
## 概览
这是一个基于Golang及Gin框架，实现请求调用翻译API的项目，目前包含: Youdao, Baidu, Tencent.
## 用法
### 主机
默认本地: 127.0.0.1:8080
### 请求参数

参数名 | 默认值 | 可选值 | 含义 | 注意事项
---- | ----- | ------ | ------- | --------
text | 无 | 无 | 需翻译的文本 | 请求必须包含的参数
from | "auto" | 请参照对应API文档 | 源语言 | API可自动识别
to | "auto" | 请参照对应API文档 |目标语言 | API可自动识别，但调用腾讯翻译必须明确语种，例如英文"en"
method | "post" | "post"/"get" | 请求方式 | 实现了get与post两种方式，默认为post
api | "youdao" | "youdao"/"baidu"/"tencent" | 调用的API名称 | 默认调用有道翻译
