package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

var jwtKey = []byte("your_secret_key_here") //토큰에 사용할 서명 비밀키 이건 추후 변경예정

func (setup *OrgSetup) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Login request")
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	valid, role := IsValidUser(loginRequest.Username, loginRequest.Password)
	if !valid {
		fmt.Println("ID 및 비밀번호 인증 실패")
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// 역할에 따라 토큰 생성
	fmt.Printf("ID 및 비밀번호 인증 완료. 토큰 생성 (역할: %s)\n", role)
	token, err := GenerateToken(loginRequest.Username, role)
	if err != nil {
		fmt.Println("토큰 생성 오류:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	fmt.Println("토큰:", token)
	response := map[string]string{"token": token}
	JsonResponse(w, http.StatusOK, response)
	fmt.Println("응답 전송 완료")
}

func (setup *OrgSetup) Loadverify(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received verify request")

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		fmt.Println("Only POST requests are allowed")
		return
	}
	type TokenRequest struct {
		Token string `json:"token"`
	}
	type CustomClaims struct {
		jwt.StandardClaims
	}

	var tokenReq TokenRequest

	if err := json.NewDecoder(r.Body).Decode(&tokenReq); err != nil {
		http.Error(w, "Failed to decode token from request body", http.StatusBadRequest)
		fmt.Println("Failed to decode token from request body")
		return
	}
	token, err := jwt.ParseWithClaims(tokenReq.Token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {
		http.Error(w, "Failed to parse token", http.StatusUnauthorized)
		fmt.Println("Failed to parse token")
		return
	}
	if _, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		response := map[string]interface{}{
			"message": "토큰 검증 성공",
		}
		JsonResponse(w, http.StatusOK, response)
		fmt.Println("토큰 검증 성공")
	} else {
		response := map[string]interface{}{
			"message": "토큰이 유효하지 않습니다. 다시 로그인하세요.",
		}
		JsonResponse(w, http.StatusUnauthorized, response)
		fmt.Println("토큰이 유효하지않습니다")
	}
}

func IsValidUser(username, password string) (bool, string) {
	validUsers := map[string]string{
		"master":   "pass",
		"employee": "pass",
		"sensor":   "pass",
		"haccp":    "pass",
	}
	roles := map[string]string{
		"master":   "admin",
		"employee": "employee",
		"sensor":   "sensor",
		"haccp":    "haccp",
	}

	if pwd, exists := validUsers[username]; exists {
		if pwd == password {
			return true, roles[username]
		}
	}
	return false, ""
}

func JsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func GenerateToken(username string, role string) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix()

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

func JWTAndRoleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Redirect(w, r, "/", http.StatusFound) // 로그인 페이지로 리디렉션
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// 토큰 검증 및 역할 확인
		claims, err := VerifyToken(tokenStr)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound) // 로그인 페이지로 리디렉션
			return
		}

		// 필요한 경우 요청 컨텍스트에 클레임 추가
		ctx := context.WithValue(r.Context(), "username", claims.Username)
		ctx = context.WithValue(ctx, "role", claims.Role)

		// 토큰에 있는 역할이 허용된 역할인지 확인
		switch r.URL.Path {
		case "/query":
			if claims.Role == "employee" || claims.Role == "admin" {
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				http.Error(w, "Forbidden", http.StatusForbidden)
			}
		case "/data":
			if claims.Role == "sensor" {
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				http.Error(w, "Forbidden", http.StatusForbidden)
			}
		case "/verify":
			if claims.Role == "haccp" || claims.Role == "admin" {
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				http.Error(w, "Forbidden", http.StatusForbidden)
			}
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	})
}
