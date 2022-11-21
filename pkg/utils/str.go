package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"gopkg.in/yaml.v2"
)

// md5加密
func Md5(password string) string {
	h := md5.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

// 随机字符串
func RandStr(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// 随机数字
func RandNumber(l int) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// 检测字符串是否含有中文
func IsChinese(str string) bool {
	var count int
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			count++
			break
		}
	}
	return count > 0
}

// convert-camel-case-string-to-snake-case
func ToSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// uint slice join to string
func JoinUintSlice(slice []uint) string {
	temp := ""
	for _, v := range slice {
		temp += fmt.Sprintf(",%d", v)
	}
	temp = strings.TrimLeft(temp, ",")
	return temp
}

func ConvStrToIntSlice(str string) []int {
	temp := make([]int, 0)
	slice := strings.Split(str, ",")
	for _, v := range slice {
		d, err := strconv.Atoi(v)
		if err != nil {
			return make([]int, 0)
		}
		temp = append(temp, d)
	}
	return temp
}

func HideStar(str string) (result string) {
	if str == "" {
		return "***"
	}
	if strings.Contains(str, "@") {
		res := strings.Split(str, "@")
		if len(res[0]) < 3 {
			resString := "***"
			result = resString + "@" + res[1]
		} else {
			res2 := Substr2(str, 0, 3)
			resString := res2 + "***"
			result = resString + "@" + res[1]
		}
		return result
	} else {
		reg := `^1[0-9]\d{9}$`
		rgx := regexp.MustCompile(reg)
		mobileMatch := rgx.MatchString(str)
		if mobileMatch {
			result = Substr2(str, 0, 3) + "****" + Substr2(str, 7, 11)
		} else {
			nameRune := []rune(str)
			lens := len(nameRune)
			if lens <= 1 {
				result = "***"
			} else if lens == 2 {
				result = string(nameRune[:1]) + "*"
			} else if lens == 3 {
				result = string(nameRune[:1]) + "*" + string(nameRune[2:3])
			} else if lens == 4 {
				result = string(nameRune[:1]) + "**" + string(nameRune[lens-1:lens])
			} else if lens > 4 {
				result = string(nameRune[:2]) + "***" + string(nameRune[lens-2:lens])
			}
		}
		return
	}
}

func Substr2(str string, start int, end int) string {
	rs := []rune(str)
	return string(rs[start:end])
}

// 计算密码分数
func CalPasswordScore(ps string) int {
	score := 0

	if len(ps) <= 4 {
		score += 5
	} else if len(ps) >= 5 && len(ps) <= 7 {
		score += 10
	} else {
		score += 25
	}

	lowerABCCount := GetStringLowerABCCount(ps)
	upperABCCount := GetStringUpperABCCount(ps)
	if (lowerABCCount > 0 && upperABCCount == 0) || (upperABCCount > 0 && lowerABCCount == 0) {
		score += 10
	} else if lowerABCCount > 0 && upperABCCount > 0 {
		score += 20
	}

	numCount := GetStringNumCount(ps)
	if numCount == 1 {
		score += 10
	} else if numCount > 1 {
		score += 20
	}

	symCount := GetStringSymCount(ps)
	if symCount == 1 {
		score += 10
	} else if symCount > 1 {
		score += 25
	}

	level := CheckPasswordLever(ps)
	if level == 2 {
		score += 7
	} else if level == 3 {
		score += 8
	} else if level == 4 {
		score += 10
	}

	return score
}

// 获取字符串包含的数字数量
func GetStringNumCount(s string) int {
	rege := regexp.MustCompile(`\d`)
	return len(rege.FindAllString(s, -1))
}

// 获取字符串包含的数字数量
func GetStringLowerABCCount(s string) int {
	rege := regexp.MustCompile(`[a-z]`)
	return len(rege.FindAllString(s, -1))
}

// 获取字符串包含的数字数量
func GetStringUpperABCCount(s string) int {
	rege := regexp.MustCompile(`[A-Z]`)
	return len(rege.FindAllString(s, -1))
}

// 获取字符串包含的符号数量
func GetStringSymCount(s string) int {
	rege := regexp.MustCompile(`[!@#~$%^&*()+|_]`)
	return len(rege.FindAllString(s, -1))
}

// 密码强度必须为字⺟⼤⼩写+数字+符号，8位以上
func CheckPasswordLever(ps string) int {
	level := 0
	if len(ps) < 8 {
		return level
	}
	num := `[0-9]{1}`
	a_z := `[a-z]{1}`
	A_Z := `[A-Z]{1}`
	symbol := `[!@#~$%^&*()+|_]{1}`

	if b, err := regexp.MatchString(num, ps); !b || err != nil {
	} else {
		level++
	}
	if b, err := regexp.MatchString(a_z, ps); !b || err != nil {
	} else {
		level++
	}
	if b, err := regexp.MatchString(A_Z, ps); !b || err != nil {
	} else {
		level++
	}
	if b, err := regexp.MatchString(symbol, ps); !b || err != nil {
	} else {
		level++
	}
	return level
}

// 转换请求IP
func ConvIP(header http.Header) string {
	xff, ok := header["X-Forwarded-For"]
	if !ok {
		return "127.0.0.1"
	}
	if len(xff) == 0 {
		return "127.0.0.1"
	}
	ipArray := strings.Split(xff[0], ",")
	ip := ipArray[0]
	return ip
}

// 将IP地址转化为二进制String
func IpToBinary(ip string) string {
	str := strings.Split(ip, ".")
	var ipstr string
	for _, s := range str {
		i, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			fmt.Println(err)
		}
		ipstr = ipstr + fmt.Sprintf("%08b", i)
	}
	return ipstr
}

// 判断IP地址和掩码地址是否匹配（变量ip为字符串，例子192.168.56.4 iprange为地址端192.168.56.64/26）
func MatchIPMask(ip, iprange string) bool {
	ipb := IpToBinary(ip)
	ipr := strings.Split(iprange, "/")
	if len(ipr) == 1 {
		return false
	}
	masklen, err := strconv.ParseUint(ipr[1], 10, 32)
	if err != nil {
		return false
	}
	iprb := IpToBinary(ipr[0])
	return strings.EqualFold(ipb[0:masklen], iprb[0:masklen])
}

// FormatJsonOrYml 格式化json或yml
func FormatJsonOrYml(data string) (string, error) {
	t := make(map[string]any)
	var t2 []byte
	var err error
	if data[0:1] == "{" {
		_ = json.Unmarshal([]byte(data), &t)
		t2, err = json.MarshalIndent(t, "", "\t")
	} else {
		_ = yaml.Unmarshal([]byte(data), &t)
		t2, err = yaml.Marshal(t)
	}
	if err != nil {
		return "", err
	}
	return string(t2), nil
}
