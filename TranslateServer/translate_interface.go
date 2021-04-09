package main

import "net/url"

/* GET方法 */
type GetAPI interface {
	GetRequestUrl(text string, from string, to string) string
	GetMethod(text string, from string, to string) string
}

/* POST方法 */
type PostAPI interface {
	PostRequestParams(text string, from string, to string) url.Values
	PostMethod(text string, from string, to string) string
}

/* 同时实现GET和POST方法 */
type TranslateAPI interface {
	GetAPI
	PostAPI
}
