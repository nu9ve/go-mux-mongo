package writeon

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/kipulab/arrivedo-api-articles/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log"
	"net/http"
	"strings"
)

// TokenValidationError and other token counters
var (
	TokenValidationError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "invalid_token_errors",
		Help: "The total token validation errors",
	})
	TokenUsage = promauto.NewCounter(prometheus.CounterOpts{
		Name: "token_usage",
		Help: "The total number of token usages",
	})
	TokenParsingError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "token_parsing_errors",
		Help: "The total of Token parsing issues",
	})
	WrongTokenError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "wrong_token_requests",
		Help: "The total of requests with wrong authorization token structure",
	})
	MissingTokenError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "missing_token_requests",
		Help: "The total of requests with no authorization token",
	})
)

// ProtectMiddleware to do the auth process
func ProtectMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				TokenUsage.Inc()
				permissionGranted := true
				var blackList []string
				token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("there was an error")
					}
					return []byte(config.Config.Secret), nil
				})
				if err != nil {
					TokenParsingError.Inc()
					log.Print(err.Error())
					ErrorResponse(w, http.StatusUnauthorized, "Invalid authorization token", "")
					return
				}
				if token.Valid {
					if config.Config.BlackList != nil{
						blackList = config.Config.BlackList
					}
					for _, ip := range blackList {
						// TODO improve this validating more data in the body
						if strings.Split(r.RemoteAddr, ":")[0] == ip {
							permissionGranted = false
						}
					}
					if !permissionGranted {
						ErrorResponse(w, http.StatusForbidden, "Invalid request", "")
						return
					}

					apiRate := token.Claims.(jwt.MapClaims)["rate_limit"]
					if apiRate == nil {
						apiRate = 0
					}
					orgSlugList := token.Claims.(jwt.MapClaims)["organizations_slug"]
					if orgSlugList == nil {
						orgSlugList = []interface{}{}
					}
					context.Set(r, "customer", token.Claims.(jwt.MapClaims)["customer"])
					context.Set(r, "orgList", token.Claims.(jwt.MapClaims)["organizations"])
					context.Set(r, "orgSlugList", orgSlugList)
					context.Set(r, "rateLimit", apiRate)
					vals := r.URL.Query()
					language := config.Config.DefaultLanguage
					if val, ok := vals["language"]; ok {
						language = val[0]
					}
					context.Set(r, "language", language)
					log.Printf("%v %v %v %v", context.Get(r,"customer"), r.URL, r.RemoteAddr, r.Referer())
					next(w, r)
				} else {
					TokenValidationError.Inc()
					ErrorResponse(w, http.StatusUnauthorized, "Invalid authorization token", "")
					return
				}
			} else {
				WrongTokenError.Inc()
				ErrorResponse(w, http.StatusUnauthorized, "Wrong authorization header", "")
				return
			}
		} else {
			MissingTokenError.Inc()
			ErrorResponse(w, http.StatusUnauthorized, "An authorization header is required", "")
			return
		}
	})
}

// AddUserMiddleware doc to be defined
func AddUserMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("there was an error")
					}
					return []byte(config.Config.Secret), nil
				})
				if err != nil {
					http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
					return
				}
				if token.Valid {
					context.Set(r, "customer", token.Claims.(jwt.MapClaims)["customer"])
					next(w, r)
				} else {
					http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
					return
				}
			}
		} else {
			next(w, r)
		}
	})
}
