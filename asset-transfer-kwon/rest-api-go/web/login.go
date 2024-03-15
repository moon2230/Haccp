package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("your_secret_key_here") //토큰에 사용할 서명 비밀키 이건 추후 변경예정

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

	if isValidUser(loginRequest.Username, loginRequest.Password) {
		fmt.Println("ID 및 비밀번호 인증 완료. 토큰 생성")
		token, err := generateToken(loginRequest.Username)
		if err != nil {
			fmt.Println("토큰 생성 오류:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		fmt.Println("토큰:", token)
		response := map[string]string{"token": token}
		jsonResponse(w, http.StatusOK, response)
		fmt.Println("응답 전송 완료")
	} else {
		fmt.Println("ID 및 비밀번호 인증 실패")
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
	}
}

func (setup *OrgSetup) verifyToken(w http.ResponseWriter, r *http.Request) {
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
		return []byte(secretKey), nil
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
		jsonResponse(w, http.StatusOK, response)
		fmt.Println("토큰 검증 성공")
	} else {
		response := map[string]interface{}{
			"message": "토큰이 유효하지 않습니다. 다시 로그인하세요.",
		}
		jsonResponse(w, http.StatusUnauthorized, response)
		fmt.Println("토큰이 유효하지않습니다")
	}
}

func isValidUser(username, password string) bool {
	return username == "user" && password == "pass"
}

func generateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix()
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func jsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
