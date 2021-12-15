package model

import "time"
import "go.mongodb.org/mongo-driver/bson/primitive"


// JwtToken struct
type JwtToken struct {
	// Fresh Token code for the customer
	Token     string `json:"token"`
	// Type of customer properties
	Scope     string `json:"scope"`
	// Expiration time in seconds. The default value is one day.
	ExpiresIn int64  `json:"expires_in"`
}

// UserTemplate struct
type UserTemplate struct {
	Email		 		string 		`json:"email" bson:"email"`
	Username 		string 		`json:"username" bson:"username"`
	Password 		string 		`json:"password" bson:"password"`
	CreatedAt		time.Time `json:"createdAt" bson:"createdAt"`
}

// User struct
type User struct {
	ID    				primitive.ObjectID 	`bson:"_id" json:"id,omitempty"`
	Name  				string    					`json:"name" bson:"name"`
	Email		 			string 							`json:"email" bson:"email"`
	Username 			string 							`json:"username" bson:"username"`
	Password		 	string 							`json:"password" bson:"password"`
	PasswordHash 	string 							`json:"passwordHash" bson:"passwordHash"`
	PasswordSalt 	string 							`json:"passwordSalt" bson:"passwordSalt"`
	Token				 	string 							`json:"token" bson:"token"`
	IsDisabled 		bool 								`json:"isDisabled" bson:"isDisabled"`
	CreatedAt			time.Time 					`json:"createdAt" bson:"createdAt"`
}

// https://stackoverflow.com/questions/25218903/how-are-people-managing-authentication-in-go
// create table UserSession (
// 	SessionKey text primary key,
// 	UserID int not null, -- Could have a hard "references User"
// 	LoginTime <time type> not null,
// 	LastSeenTime <time type> not null
//  )

// ResponseResult struct
type ResponseResult struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}

// TopicTemplate struct
type TopicTemplate struct {
	Name  			string    `json:"name" bson:"name"`
	CreatedAt		time.Time `json:"createdAt" bson:"createdAt"`
}

// Topic struct
type Topic struct {
	ID    			primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name  			string    `json:"name" bson:"name"`
	CreatedAt		time.Time `json:"createdAt" bson:"createdAt"`
}

// EssayTemplate struct
type EssayTemplate struct {
	Topic   		string    `json:"topic" bson:"topic"`
	Title   		string    `json:"title" bson:"title"`
	Body    		string    `json:"body" bson:"body"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
}

// Essay struct
type Essay struct {
	ID    			primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Topic   		string    `json:"topic" bson:"topic"`
	Title   		string    `json:"title" bson:"title"`
	Body    		string    `json:"body" bson:"body"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
}