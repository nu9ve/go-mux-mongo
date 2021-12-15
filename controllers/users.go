package writeon

import (
  "encoding/json"
	"fmt"
	"net/http"
	// "log"
	"time"
	// "io/ioutil"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	jwt "github.com/dgrijalva/jwt-go"

	"writeon/model"
)

const SECRET_KEY = "jfdsfjsdk"
const TokenExpiration = 300


func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var res model.ResponseResult
	existingUser := &model.User{}
	var userTemplate model.UserTemplate
	json.NewDecoder(r.Body).Decode(&userTemplate)
	// check error 
	if userTemplate.Username == "" {
		res.Error = "User needs username"
		json.NewEncoder(w).Encode(res)
		return
	}
	if userTemplate.Email == "" {
		res.Error = "User needs email"
		json.NewEncoder(w).Encode(res)
		return
	}
	if userTemplate.Password == "" {
		res.Error = "User needs password"
		json.NewEncoder(w).Encode(res)
		return
	}

	err := mh.GetUser(existingUser, bson.D{{"username", userTemplate.Username}})
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			hash, err := bcrypt.GenerateFromPassword([]byte(userTemplate.Password), 5)

			if err != nil {
				res.Error = "Error While Hashing Password, Try Again"
				json.NewEncoder(w).Encode(res)
				return
			}
			userTemplate.Password = string(hash)
			userTemplate.CreatedAt = time.Now()

			inserted, insertErr := mh.CreateUser(&userTemplate)
			if insertErr != nil {
				res.Error = "Error While Creating User, Try Again"
				json.NewEncoder(w).Encode(res)
				return
			}

			finalUser := &model.User{}
			finalUser.ID = inserted.InsertedID.(primitive.ObjectID)
			finalUser.CreatedAt = userTemplate.CreatedAt

			jwtToken, err := GenerateJWT(finalUser)
			if err != nil{
				res.Error = err.Error()
				json.NewEncoder(w).Encode(res)
				return
			}

			finalUser.Token = jwtToken.Token
			json.NewEncoder(w).Encode(finalUser)
			return
		}

		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	fmt.Println(err)

	res.Result = "Username already Exists!!"
	json.NewEncoder(w).Encode(res)
	return
}


func GenerateJWT(user *model.User) (model.JwtToken, error) {
	currentDate := time.Now()
	expDate := currentDate.Add(time.Second * time.Duration(TokenExpiration)).Unix()
	claims := jwt.MapClaims{
		"user":      user.ID,
		// "role": user.UserType,
		// "rate_limit":    user.RateLimit,
		"exp_date":      expDate,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := SECRET_KEY
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return model.JwtToken{}, err
	}

	customerToken := model.JwtToken{
		Token:     "bearer " + tokenString,
		// Scope:     customer.UserType,
		ExpiresIn: expDate - time.Now().Unix(),
	}
	return customerToken, nil
}



func LoginUser(w http.ResponseWriter, r *http.Request) {

	var res model.ResponseResult
	existingUser := &model.User{}
	var user model.UserTemplate
	json.NewDecoder(r.Body).Decode(&user)
	err := mh.GetUser(existingUser, bson.D{{"username", user.Username}})

	if err != nil {
		res.Error = "Invalid username"
		json.NewEncoder(w).Encode(res)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))

	if err != nil {
		res.Error = "Invalid credentials"
		json.NewEncoder(w).Encode(res)
		return
	}

	jwtToken, err := GenerateJWT(existingUser)
	if err != nil{
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}


  // var resp = map[string]interface{}{"status": false, "message": "logged in"}
  // resp["token"] = tokenString
  // resp["tk"] = tk


	existingUser.Token = jwtToken.Token
	existingUser.Password = ""

	json.NewEncoder(w).Encode(existingUser)

}


// LogoutUser function
func LogoutUser(w http.ResponseWriter, r *http.Request) {

	var res model.ResponseResult
	existingUser := &model.User{}
	var user model.UserTemplate
	json.NewDecoder(r.Body).Decode(&user)
	err := mh.GetUser(existingUser, bson.D{{"username", user.Username}})

	if err != nil {
		res.Error = "Invalid username"
		json.NewEncoder(w).Encode(res)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))

	if err != nil {
		res.Error = "Invalid password"
		json.NewEncoder(w).Encode(res)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  existingUser.ID,
		"username":  existingUser.Username,
		"name":  existingUser.Name,
	})

	tokenString, err := token.SignedString([]byte("secret"))
	
	if err != nil {
		res.Error = "Error while generating token,Try again"
		json.NewEncoder(w).Encode(res)
		return
	}

	existingUser.Token = tokenString
	existingUser.Password = ""

	json.NewEncoder(w).Encode(existingUser)

}

// UserProfile function
func UserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	var result model.User
	var res model.ResponseResult
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result.Username = claims["username"].(string)
		result.Name = claims["name"].(string)
		// result.ID = claims["id"].(string)

		json.NewEncoder(w).Encode(result)
		return
	}
	res.Error = err.Error()
	json.NewEncoder(w).Encode(res)
	return

}


