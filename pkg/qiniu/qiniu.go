package qiniu

import (
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type Qiniu struct {
	AK             string
	SK             string
	Bucket         string
	Domain         string
	TokenExpireSec uint64
}

func New(ak, sk, bucket, domain string, tokenExpireSec uint64) *Qiniu {
	return &Qiniu{
		AK:             ak,
		SK:             sk,
		Bucket:         bucket,
		Domain:         domain,
		TokenExpireSec: tokenExpireSec,
	}
}

func (qn *Qiniu) GetUploadToken() string {
	mac := qbox.NewMac(qn.AK, qn.SK)
	putPolicy := storage.PutPolicy{
		Scope:   qn.Bucket,
		Expires: qn.TokenExpireSec,
	}
	upToken := putPolicy.UploadToken(mac)
	return upToken
}

func (qn *Qiniu) GetDomain() string {
	return qn.Domain
}
