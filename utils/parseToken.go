package utils

import (
"errors"
"github.com/dgrijalva/jwt-go"
"github.com/wonderivan/logger"
)

var JWTToken jwtToken

type jwtToken struct {}

//token解析后对应的结构体，包含自定义信息以及jwt签名信息
type CustomClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

//加解密因子,跟前端对应，前端生成时也要用这个因子
const SECRET = "adoodevops"

//解析token
func (*jwtToken) ParseToken(tokenStr string) (claims *CustomClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token)(interface{}, error){
		return []byte(SECRET), nil
	})
	if err != nil {
		logger.Error("parse token failed", err)
		//处理token解析后的各种报错情况
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("TokenMalformed")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("TokenExpired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("TokenNotValidYet")
			} else {
				return nil, errors.New("TokenInvalid")
			}
		}
	}
	//转换成*CustomClaims类型并返回
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("解析Token失败")
}