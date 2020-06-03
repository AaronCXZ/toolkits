package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type Cors struct {
	Origin      string
	Headers     string
	Credentials string
	Methods     string
	Types       string
}

func DefaultCors() *Cors {
	return &Cors{
		Origin:      "*",
		Headers:     "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token",
		Credentials: "true",
		Methods:     "POST, GET, OPTIONS, PUT, DELETE",
		Types:       "application/json;charset=UTF-8",
	}
}

func NewCors(ops ...Option) *Cors {
	cors := DefaultCors()
	for _, opt := range ops {
		opt(cors)
	}
	return cors
}

type Option func(*Cors)

// 允许的域
func WithOrigin(urls ...string) Option {
	url := strings.Join(urls, ",")
	return func(c *Cors) {
		c.Origin = url
	}
}

// header的类型
func WithHeaders(headers ...string) Option {
	header := strings.Join(headers, ",")
	return func(c *Cors) {
		c.Headers = header
	}
}

// 设置为true，允许ajax异步请求带cookie信息
func WithCredentials(credentials string) Option {
	return func(c *Cors) {
		c.Credentials = credentials
	}
}

// 允许请求方法
func WithMethods(methods ...string) Option {
	method := strings.Join(methods, ",")
	return func(c *Cors) {
		c.Methods = method
	}
}

// 返回数据格式
func WithTypes(typ string) Option {
	return func(c *Cors) {
		c.Types = typ
	}
}

func AccessCors(cors *Cors) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", cors.Origin)
		c.Header("Access-Control-Allow-Headers", cors.Headers)
		c.Header("Access-Control-Allow-Credentials", cors.Credentials)
		c.Header("Access-Control-Allow-Methods", cors.Methods)
		c.Header("content-type", cors.Types)
		c.Next()
	}
}
