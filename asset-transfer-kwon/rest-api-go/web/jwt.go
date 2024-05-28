package web

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key_here") //토큰에 사용할 서명 비밀키 이건 추후 변경예정

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Redirect(w, r, "/", http.StatusFound) // 로그인 페이지로 리디렉션
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// 토큰 검증
		if !VerifyToken(tokenStr) {
			http.Redirect(w, r, "/", http.StatusFound) // 로그인 페이지로 리디렉션
			return
		}

		// 필요한 경우 요청 컨텍스트에 클레임 추가
		ctx := context.WithValue(r.Context(), "token", tokenStr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func VerifyToken(tokenStr string) bool {
	claims := &jwt.StandardClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return false
	}

	return true
}

func GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix()
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
