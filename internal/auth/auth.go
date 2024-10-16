package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"strings"
	"time"
)

// Claims — структура утверждений, которая включает стандартные утверждения
// и одно пользовательское — UserID
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// Утверждения для текущего пользователя
var UserClaims Claims

const TOKEN_EXP = time.Hour * 3
const SECRET_KEY = "Ec#9<8gc,/7zu*vTeX)=q96JQw+I8|]6/+*'8YWqx\"G06Yy\"H;)wwn`K+*Z;C(i"

// BuildJWTString создаёт токен и возвращает его в виде строки.
func BuildJWTString() (string, error) {
	//Создаём ИД пользователя
	UserID, err := CreateUserID(32)
	if err != nil {
		return ``, err
	}
	UserClaims = Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		// собственное утверждение
		UserID: UserID,
	}

	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims)

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return ``, err
	}

	// возвращаем строку токена
	return tokenString, nil
}

// Генерация случайной строки
func CreateUserID(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return ``, err
	}
	return hex.EncodeToString(b), nil
}

// Получаем текущий ИД пользователя (созданный или полученный из куков)
func GetCurrentUserID() string {
	return UserClaims.UserID
}

// Получаем ИД пользователя из токена, который получаем из куков
func GetUserID(r *http.Request) string {
	//Получаем куки
	cookie, err := r.Cookie("Authorization")
	if err != nil || cookie.Value == "" {
		return ""
	}
	//Получаем токен из куков
	splitToken := strings.Split(cookie.Value, "Bearer ")
	if len(splitToken) < 2 {
		return ""
	}
	tokenString := splitToken[1]
	//Парсим токен
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return "", fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SECRET_KEY), nil
		})
	if err != nil {
		log.Print(err)
		return ""
	}

	if !token.Valid {
		log.Print("Token invalid")
		return ""
	}
	return UserClaims.UserID
}

// Авторизация для роутера Chi
func Authenticator() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			//Если не нашли UserID - создаём новую куку
			if GetUserID(r) == "" {
				token, err := BuildJWTString()
				if err != nil {
					panic(err)
				}
				cookie := &http.Cookie{
					Name:   "Authorization",
					Value:  "Bearer " + token,
					MaxAge: int(TOKEN_EXP),
				}
				http.SetCookie(w, cookie)
			}
			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}
