package aliyun

import (
	"context"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

// 阿里云短信模板参数
type SmsCaptchaParams struct {
	SignName     string
	TemplateCode string
}

type Aliyun struct {
	SmsClient  *dysmsapi.Client
	SmsCaptcha *SmsCaptchaParams
}

func New(regionId, accessKeyId, accessKeySecret, signName, templateCode string) (*Aliyun, error) {
	// client, err := dysmsapi.NewClientWithAccessKey("cn-qingdao", "<your-access-key-id>", "<your-access-key-secret>")
	// use STS Token
	// client, err := dysmsapi.NewClientWithStsToken("cn-qingdao", "<your-access-key-id>", "<your-access-key-secret>", "<your-sts-token>")

	smsClient, err := dysmsapi.NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret)
	if err != nil {
		return nil, err
	}
	return &Aliyun{
		SmsClient: smsClient,
		SmsCaptcha: &SmsCaptchaParams{
			SignName:     signName,
			TemplateCode: templateCode,
		},
	}, nil
}

// 发送短信
func (aliyun *Aliyun) SendSMS(
	ctx context.Context,
	mobile,
	signName,
	templateCode string,
) (*dysmsapi.SendSmsResponse, error) {
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.PhoneNumbers = mobile                         //手机号码
	request.SignName = aliyun.SmsCaptcha.SignName         //短信签名名称
	request.TemplateCode = aliyun.SmsCaptcha.TemplateCode //短信模板ID

	response, err := aliyun.SmsClient.SendSms(request)
	return response, err
}
