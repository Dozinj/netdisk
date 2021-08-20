package model

import (
	"github.com/dgrijalva/jwt-go"
)

// Custom claims structure
type CustomClaims struct {
	ID          int		//用户主键id
	Username    string
	jwt.StandardClaims
}




