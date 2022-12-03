package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"okc/config"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jordan-wright/email"
	uuid2 "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

type UserClaims struct {
	Id       int    `json:"id"`       // 用户ID
	Account  string `json:"account"`  // 账号
	Password string `json:"password"` // 密码
	jwt.RegisteredClaims
}

// token加解密keys
var tokenKeys = []byte("b2tjCg==")

// GenerateMd5
// 生成Md5
func GenerateMd5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return fmt.Sprintf("%x", m.Sum(nil))
}

// GenerateHashPassword
// 生成哈希密码
func GenerateHashPassword(password string, mode bool) string {

	if mode {
		salt := "ABCDEFG"
		for _, s := range password {
			salt += GenerateMd5(string(s))
		}
		return GenerateMd5(salt)
	}

	return "TPSHOP" + GenerateMd5(password)
}

// GenerateSha256
// 生成Sha256
func GenerateSha256(str string) string {
	return fmt.Sprintf("%X", sha256.New().Sum([]byte(str)))
}

// GenerateRandExtensionCode
// 生成随机推广码
func GenerateRandExtensionCode(lens int) string {
	rand.Seed(time.Now().UnixNano())

	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXZY0123456789"
	result := ""
	for i := 0; i < lens; i++ {
		num := rand.Intn(len(str))
		result += string(str[num])
	}

	return result
}

// GenerateUUID
// 生成uuid
func GenerateUUID() string {
	uuid := uuid2.NewV4().String()
	return uuid
}

// GenerateToken
// 生成token
func GenerateToken(Id int, account string, password string) (string, error) {
	UserClaims := &UserClaims{
		Id:               Id,
		Account:          account,
		Password:         password,
		RegisteredClaims: jwt.RegisteredClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims)
	tokenString, err := token.SignedString(tokenKeys)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// AnalyseToken
// 解析token
func AnalyseToken(token string) (*UserClaims, error) {
	UserClaims := new(UserClaims)
	claims, err := jwt.ParseWithClaims(token, UserClaims,
		func(token *jwt.Token) (interface{}, error) {
			return tokenKeys, nil
		})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("analyse Token Error:%v", err)
	}
	return UserClaims, nil
}

// AnalyseUUID
// 解析uuid
func AnalyseUUID(uuid string) (uuid2.UUID, error) {
	uid, err := uuid2.FromString(uuid)
	if err != nil {
		return uid, err
	}

	return uid, err
}

// GenerateCode
// 生成六位数验证码
func GenerateCode() string {
	rand.Seed(time.Now().UnixNano())
	code := ""
	for i := 0; i < 6; i++ {
		code += strconv.Itoa(rand.Intn(10))
	}
	return code
}

// SendEmailCode
// 发送邮件验证码
func SendEmailCode(toUserEmail string, code string) error {
	conf := config.Config().EMAIL.(map[interface{}]interface{})

	e := email.NewEmail()
	e.From = fmt.Sprintf("Admin <%s>", conf["USERNAME"])
	e.To = []string{toUserEmail}
	e.Subject = "验证码已发送，请查收"
	e.HTML = []byte("您的验证码：<b>" + code + "</b>")

	//return e.SendWithTLS("smtp.gmail.com:465",
	//	smtp.PlainAuth("", "watterwtli@gmail.com", "Lf0720@.", "smtp.gmail.com"),
	//	&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.gmail.com"})

	return e.SendWithTLS(conf["SERVER_HOST_PORT"].(string),
		smtp.PlainAuth("", conf["USERNAME"].(string), conf["PASSWORD"].(string), conf["SERVER_HOST"].(string)),
		&tls.Config{InsecureSkipVerify: true, ServerName: conf["SERVER_HOST"].(string)})
}

// RefStructGetFieldPtr
// 反射结构体获取字段指针
func RefStructGetFieldPtr(stc interface{}, field ...interface{}) ([]interface{}, error) {
	stcTypeOf := reflect.TypeOf(stc)
	stcValueOf := reflect.ValueOf(stc)
	fieldPtrSlice := make([]interface{}, 0)
	if stcTypeOf.Elem().Kind() != reflect.Struct {
		log.Println("stcTypeOf.Kind() :", stcTypeOf.Kind())
		return nil, errors.New("parameter type is not Struct")
	}

	if stcTypeOf.Kind() != reflect.Ptr {
		log.Println("stcTypeOf.Kind() :", stcTypeOf.Kind())
		return nil, errors.New("parameter type is not Ptr")
	}

	stcFieldNum := stcValueOf.Elem().NumField()

	for i := 0; i < stcFieldNum; i++ {
		fieldName := stcTypeOf.Elem().Field(i).Name

		for _, v := range field {
			if v == "*" || fieldName == v {
				fieldPointer := reflect.NewAt(stcValueOf.Elem().Field(i).Type(), unsafe.Pointer(stcValueOf.Elem().Field(i).UnsafeAddr())).Pointer()
				fieldKind := stcValueOf.Elem().Field(i).Type().Kind()
				switch fieldKind {
				case reflect.Int:
					fieldContext := (*int)(unsafe.Pointer(fieldPointer))
					fieldPtrSlice = append(fieldPtrSlice, fieldContext)
					break
				case reflect.String:
					fieldContext := (*string)(unsafe.Pointer(fieldPointer))
					fieldPtrSlice = append(fieldPtrSlice, fieldContext)
					break
				case reflect.Float64:
					fieldContext := (*float64)(unsafe.Pointer(fieldPointer))
					fieldPtrSlice = append(fieldPtrSlice, fieldContext)
					break
				}
			}
		}

	}
	return fieldPtrSlice, nil
}

// RefStructChangeMap
// 反射结构体转换成map
func RefStructChangeMap(structPtr interface{}) (map[string]interface{}, error) {

	typeOf := reflect.TypeOf(structPtr)
	valueOf := reflect.ValueOf(structPtr)
	kind := typeOf.Kind()

	if kind != reflect.Ptr {
		return nil, errors.New(" Parameter &Type Not *Ptr ")
	}
	if typeOf.Elem().Kind() != reflect.Struct {
		return nil, errors.New(" Parameter *Type Not Struct ")
	}

	container := make(map[string]interface{})
	fieldNum := typeOf.Elem().NumField()
	for i := 0; i < fieldNum; i++ {
		fieldName := typeOf.Elem().Field(i).Name
		fieldTag := typeOf.Elem().Field(i).Tag.Get("json")

		if strings.Index(fieldTag, ",") == -1 {
			container[fieldTag] = valueOf.Elem().FieldByName(fieldName).Interface()
		} else {
			keys := strings.Split(fieldTag, ",")
			container[keys[0]] = valueOf.Elem().FieldByName(fieldName).Interface()
		}
	}
	log.Println("kind : ", kind)
	log.Println("*kind : ", typeOf.Elem().Kind())
	log.Println("data : ", container)
	return container, nil
}

// TrimRightZeroByFloatStr
// 根据浮点型字符串去除零
func TrimRightZeroByFloatStr(floatStr string) string {

	if strings.Index(floatStr, ".") != -1 {
		floatStr = strings.TrimRight(floatStr, "0")
		floatStr = strings.TrimRight(floatStr, ".")
	}

	return floatStr
}

// InCheckIntOrIntSlice
// 检查int类型是否存在int类型切片中
func InCheckIntOrIntSlice(index int, checkValue []int) bool {
	for _, v := range checkValue {
		if index == v {
			return true
		}
	}

	return false
}

// MatchRandBoolSlicePop
// 随机弹出Bool类型切片中的值
func MatchRandBoolSlicePop(boolSlice []bool) bool {

	var res bool
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < len(boolSlice); i++ {
		index := rand.Intn(len(boolSlice))
		res = boolSlice[index]
	}

	return res
}

// MatchRandFloatBySectionAndPec
// 根据定制区间随机生成浮点字符串
func MatchRandFloatBySectionAndPec(min, max string, pec int) (string, error) {

	minParseFloat, err := strconv.ParseFloat(min, 64)
	if err != nil {
		return "", err
	}
	maxParseFloat, err := strconv.ParseFloat(max, 64)
	if err != nil {
		return "", err
	}

	rand.Seed(time.Now().UnixNano())
	num := rand.Float64()*(maxParseFloat-minParseFloat) + minParseFloat

	return strconv.FormatFloat(num, 'f', pec, 64), nil
}

// CheckDecimalLenByString
// 根据字符串检查是否有小数点位数(必须是浮点字符串)
func CheckDecimalLenByString(decimalStr string) int {

	decimals := 0
	indexOf := strings.Index(decimalStr, ".")

	if indexOf != -1 {
		decimalStr = strings.TrimRight(decimalStr, "0")
		decimalStr = strings.TrimRight(decimalStr, ".")
		decimalOf := strings.Index(decimalStr, ".")

		if decimalOf != -1 {
			decimals = len(decimalStr) - decimalOf - 1
		}
	}

	return decimals
}

// BcSub
// 两个任意精度数字的减法
func BcSub(num1, num2 string) (float64, error) {

	num1ParseFloat, err := strconv.ParseFloat(num1, 64)
	if err != nil {
		return 0, err
	}
	num2ParseFloat, err := strconv.ParseFloat(num2, 64)
	if err != nil {
		return 0, err
	}

	num1Decimal := decimal.NewFromFloat(num1ParseFloat)
	num2Decimal := decimal.NewFromFloat(num2ParseFloat)

	res, _ := num1Decimal.Sub(num2Decimal).Float64()

	return res, nil
}

// BcAdd
// 两个任意精度数字的加法计算
func BcAdd(num1, num2 string) (float64, error) {

	num1ParseFloat, err := strconv.ParseFloat(num1, 64)
	if err != nil {
		return 0, err
	}
	num2ParseFloat, err := strconv.ParseFloat(num2, 64)
	if err != nil {
		return 0, err
	}

	num1Decimal := decimal.NewFromFloat(num1ParseFloat)
	num2Decimal := decimal.NewFromFloat(num2ParseFloat)

	res, _ := num1Decimal.Add(num2Decimal).Float64()

	return res, nil
}

// BcMul
// 两个任意精度数字乘法计算
func BcMul(num1, num2 string) (float64, error) {
	num1ParseFloat, err := strconv.ParseFloat(num1, 64)
	if err != nil {
		return 0, err
	}
	num2ParseFloat, err := strconv.ParseFloat(num2, 64)
	if err != nil {
		return 0, err
	}

	num1Decimal := decimal.NewFromFloat(num1ParseFloat)
	num2Decimal := decimal.NewFromFloat(num2ParseFloat)

	res, _ := num1Decimal.Mul(num2Decimal).Float64()

	return res, nil
}

// BcDiv
// 两个任意精度的数字除法计算
func BcDiv(num1, num2 string) (float64, error) {
	num1ParseFloat, err := strconv.ParseFloat(num1, 64)
	if err != nil {
		return 0, err
	}
	num2ParseFloat, err := strconv.ParseFloat(num2, 64)
	if err != nil {
		return 0, err
	}

	num1Decimal := decimal.NewFromFloat(num1ParseFloat)
	num2Decimal := decimal.NewFromFloat(num2ParseFloat)

	res, _ := num1Decimal.Div(num2Decimal).Float64()

	return res, nil
}

// TimeStrChangeSecond
// 时间转换秒
func TimeStrChangeSecond(times string) (second uint64) {

	times = strings.Trim(times, " ")

	if strings.Index(times, "m") != -1 { // 分

		timeSplit := strings.Split(times, "m")
		times := timeSplit[0]
		times = strings.Trim(times, " ")
		timesParseInt, _ := strconv.ParseInt(times, 10, 64)
		timeNumber := uint64(timesParseInt) * 60

		second = timeNumber

	} else if strings.Index(times, "h") != -1 { // 时

		timeSplit := strings.Split(times, "h")
		times := timeSplit[0]
		times = strings.Trim(times, " ")
		timesParseInt, _ := strconv.ParseInt(times, 10, 64)
		timeNumber := uint64(timesParseInt) * 60 * 60

		second = timeNumber

	} else if strings.Index(times, "d") != -1 { // 天

		timeSplit := strings.Split(times, "d")
		times := timeSplit[0]
		times = strings.Trim(times, " ")
		timesParseInt, _ := strconv.ParseInt(times, 10, 64)
		timeNumber := uint64(timesParseInt) * 24 * 60 * 60

		second = timeNumber
	}

	return
}
