package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// 加密
func RsaEncrypt(origData, publicKey []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密
func RsaDecrypt(ciphertext, privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	parseResult, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("parse error!")
	}
	priv := parseResult.(*rsa.PrivateKey)
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

var DefaultRSAPrivateKey = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDaCudsvnJsf82n
vNdt9Jq+j4EXJohOJ5FIy01kSNeBEO99AIdN4KIqEVe0LeSu+5IPcPAZVbLtzvTu
f3URzL3lI/pWbU5p4x3jardJQmbI/kluYUNAi++SSy+LAYB7l1tItNNNBxHvKmwR
aBJyeX2ZJyF31LKqVpwghVJJ5HdRkZbFbfSL8W19M6yiL156bKg6Kt+BS+nWQXJB
bj2H80fTSAkykWrpy4FKUi4DhJTd2jtkd9F9L2kUxB8gGhLjOgcZWdXx+K2LhErz
Nganuy4NoMElFhIVfKmOSbE9HHhu569oGqf0GGHp0pDdD6yeZ0wlrE6tC7ZlxyPD
W97E2qr3AgMBAAECggEAVYLNO+J8sXx1xQDUY/T38yAIenHMJwpxeeDxWxEOIznn
Eljwt2QPdPwULz+IXe+fWprqLqVjZMwzVo159h8bog+4D8kSZYiCojup4fs+oDjJ
x5Yxb9Dxhagi3xsZrl6vmBYCgETzjZ5Pf4wXH/nxfezQvVQXNaJ93Gss2HnXKY97
yUCNsFnb3v++LlgSDsp4GaK7OSe8tJf0fo12YvgM/m6Vk2LpIF9fVNthUFNTCleu
DRuLY7IX7G5PkezxixhF8JDmEtw+B6/YWLNZrgufxfcpuk99D9wLZOlB56aG4YZT
5vl/q2F44uxsfbw+TE+6ywbvbfpQMBECbfTvSgPsgQKBgQD4qoXFCXbdylPwREEv
+MNlTdtQ9TWXO54tLqp+N9HUTw6adDUSfRIu/EVF+HN/qusPIWhKFQjwMgAQjynX
wIXDtkRB0fgrt1TF33LfQboRXYjynM8XZASF7J5t/TcrU9q476p4eN3vJ5ZrXGf2
xSSBALDgy2jmlIKSGcYbi33lFwKBgQDgeSsCijFKgX6cadIMu4FkMIK+iYID1afw
uWvJ06mHiCw+vDRW9SzdNYp7/yXYrvPzLnd9X9k5Sz9JMLJc4XdGoDFDDOo1A/8S
kWedZPSajOx7giZGngQHfnuQiFe9auuZ90x7NfkL0yyiV/MwdGHwOqiTJM+Xembo
P28vEJzVIQKBgQCAdo/PABmxcPI2QPywTMKdFkDELTm3XGxWCTK5LBsxpHn13yz/
1S9MqLUc9cKtZN40ndyj0QQiGqKf62YBeQth1Uqj+lZMN1ULOGm+3tTCXeD+/XWb
LueLTHd4eQVEU/i968rUnBSDlZ8G7eEjwisenf3C1DLoVDa0Ra5r0n+ClwKBgFmF
2XJc3MWjGXSV+3Cag0MK2cnVm2WeGyk1Odi3MoBb/ZFTi+g2RZs/VCiZnGVreN0+
Zec5h6+C5A1zf17tiJ1BHARqrSlRm7OzC8jIz4intVSYll1Jfb/jYLJGvf9MGgRA
jV8CKn3dzYo9Wz6y27BsJHjykFwQM+RiEByMGpAhAoGBANP3HiI/oci85u49nnut
NclHf3vGiO09AyQFeHXIU2fJ98b3B14K/bRvVHe3W+FFYXnPBlwmq/jBAdSlj4NC
uGWlZCpMmmYMOX+00amh0jHcnDVIhsyoIMexU9B5qbl+Tr0AqvwGYu4BFlxE5fqt
m0DEMGHDQ//loJQb+onypSPA
-----END PRIVATE KEY-----`

var DefaultRSAPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2grnbL5ybH/Np7zXbfSa
vo+BFyaITieRSMtNZEjXgRDvfQCHTeCiKhFXtC3krvuSD3DwGVWy7c707n91Ecy9
5SP6Vm1OaeMd42q3SUJmyP5JbmFDQIvvkksviwGAe5dbSLTTTQcR7ypsEWgScnl9
mSchd9SyqlacIIVSSeR3UZGWxW30i/FtfTOsoi9eemyoOirfgUvp1kFyQW49h/NH
00gJMpFq6cuBSlIuA4SU3do7ZHfRfS9pFMQfIBoS4zoHGVnV8fiti4RK8zYGp7su
DaDBJRYSFXypjkmxPRx4buevaBqn9Bhh6dKQ3Q+snmdMJaxOrQu2Zccjw1vexNqq
9wIDAQAB
-----END PUBLIC KEY-----`
