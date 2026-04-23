package service

import "errors"

var (
	ErrAccountExists     = errors.New("该账号已存在")
	ErrAccountNotFound   = errors.New("账号不存在")
	ErrNoFieldsToUpdate  = errors.New("没有可更新的字段")
	ErrUnsupportedPlatform = errors.New("不支持的平台类型")
	ErrAlreadyRunning    = errors.New("任务已在运行中")
	ErrNotRunning        = errors.New("任务未在运行")
	ErrLoginFailed       = errors.New("登录失败")
)
