package tools

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"

	"github.com/zhangshanwen/transport/apps/admin/conf"
)

var (
	defaultExpiresTimes = 12 * time.Hour
	defaultTokenType    = ""
	rasPath             = "rsa"
	privateKey          *rsa.PrivateKey
	publicKey           *rsa.PublicKey
	method              = jwt.SigningMethodRS256 //默认256
)

type Payload struct {
	Uid       int64
	TokenType string
}

type Claims struct {
	*jwt.StandardClaims
	Payload
}

func InitJwt(project string) {
	var err error

	var privateBytes, publicBytes []byte
	rasPath += string(os.PathSeparator)
	privateBytes, err = ioutil.ReadFile(rasPath + fmt.Sprintf("%s.rsa", project))
	if err != nil {
		logrus.Panic(err)
	}
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		logrus.Fatalf("初始化jwt失败 %s\n", err)
	}
	publicBytes, err = ioutil.ReadFile(rasPath + fmt.Sprintf("%s.rsa.pub", project))
	if err != nil {
		logrus.Panic(err)
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		logrus.Fatalf("初始化jwt失败: %s\n", err)
	}
	defaultTokenType = project
}
func CreateToken(uid int64) (token string, err error) {
	var expiresAt int64 //过期时间
	if conf.C.Authorization.ExpireHour == 0 {
		expiresAt = time.Now().Add(defaultExpiresTimes).Unix()
	} else {
		expiresAt = time.Now().Add(time.Duration(conf.C.Authorization.ExpireHour) * time.Hour).Unix()
	}
	t := jwt.NewWithClaims(method, Claims{
		&jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
		Payload{uid, defaultTokenType},
	})
	return t.SignedString(privateKey)
}
func VerifyToken(tokenString string) (claims *Claims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return
	}
	claims = token.Claims.(*Claims)
	return
}
