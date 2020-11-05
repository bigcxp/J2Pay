// json web token 令牌
package util

import (
	"github.com/dgrijalva/jwt-go"
	"j2pay-server/model"
	"j2pay-server/model/response"
	"j2pay-server/pkg/logger"
	"j2pay-server/pkg/setting"
	"time"
)

type Claims struct {
	Username string             `json:"username"`
	Secret   string             `json:"secret"`
	Role     []response.CasRole `json:"role"`
	Auth     []model.Auth       `json:"auth"`
	ID       int64              `json:"id"`
	jwt.StandardClaims
}

// 类型转换
var JwtKey = []byte(setting.JwtConf.Key)

// 生成令牌
func MakeToken(account model.Account) (string, error) {
	// 过期时间
	expTime := time.Now().Add(time.Duration(setting.JwtConf.ExpTime) * time.Hour)
	tokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Username: account.UserName,
		ID:       account.ID,
		Secret:   account.Secret,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
			Subject:   "j2pay-server",
		},
	})
	return tokenClaim.SignedString(JwtKey)
}

// 解析令牌
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	logger.Logger.Error("解析jwt出错 : ", err)
	return nil, err
}
