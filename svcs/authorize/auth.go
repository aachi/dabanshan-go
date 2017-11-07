package authorize

import (
	"crypto/sha1"
	"fmt"
	"io"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
)

var (
	// SecretKey JWT secret key
	SecretKey = "welcome to dabanshan"
	//JwtMiddleware jwt middleware
	JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
)

func init() {}

// CreateJWT generat a jwt token
func CreateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	return token.SignedString([]byte(SecretKey))
}

func CalculatePassHash(pass, salt string) string {
	h := sha1.New()
	io.WriteString(h, salt)
	io.WriteString(h, pass)
	return fmt.Sprintf("%x", h.Sum(nil))
}
