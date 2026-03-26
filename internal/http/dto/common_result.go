package dto

import "sandbox-game/internal/enum"

// CommonResult 对齐接口文档中的统一响应体格式。
type CommonResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// Success 返回成功响应体。
func Success(data any) CommonResult {
	return CommonResult{
		Code: enum.SuccessCode,
		Msg:  "",
		Data: data,
	}
}

// Failure 返回失败响应体。
func Failure(code int, msg string) CommonResult {
	return CommonResult{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}
