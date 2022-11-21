package sms

import (
	"context"
	"errors"

	"github.com/go-resty/resty/v2"
)

var (
	SmsBaoLocalApi      = "http://api.smsbao.com/sms"
	SmsBaoWorldApi      = "http://api.smsbao.com/wsms"
	SmsBaoSuccessResult = "0"

	SmsBaoErr30 = "错误密码"
	SmsBaoErr40 = "账号不存在"
	SmsBaoErr41 = "余额不足"
	SmsBaoErr43 = "IP地址限制"
	SmsBaoErr50 = "内容含有敏感词"
	SmsBaoErr51 = "手机号码不正确"
)

type SendRange string

var (
	SendRangeLocal SendRange = "local"
	SendRangeWorld SendRange = "world"
)

type SmsBao struct {
	username    string
	password    string
	sendRange   SendRange
	restyClient *resty.Client
}

func NewSmsBao(username, password string, sendRange SendRange) *SmsBao {
	return &SmsBao{
		username:    username,
		password:    password,
		sendRange:   sendRange,
		restyClient: resty.New(),
	}
}

func (sb *SmsBao) Send(
	ctx context.Context,
	mobile, content string,
) error {
	url := ""
	switch sb.sendRange {
	case SendRangeLocal:
		url = SmsBaoLocalApi
	case SendRangeWorld:
		url = SmsBaoWorldApi
	}
	res, err := sb.restyClient.R().
		SetQueryParams(map[string]string{
			"u": sb.username,
			"p": sb.password,
			"m": mobile,
			"c": content,
		}).
		Get(url)
	if err != nil {
		return err
	}
	body := string(res.Body())
	if body != SmsBaoSuccessResult {
		errString := ""
		switch body {
		case "30":
			errString = SmsBaoErr30
		case "40":
			errString = SmsBaoErr40
		case "41":
			errString = SmsBaoErr41
		case "43":
			errString = SmsBaoErr41
		case "50":
			errString = SmsBaoErr50
		case "51":
			errString = SmsBaoErr51
		}
		return errors.New(errString)
	}

	return nil
}
