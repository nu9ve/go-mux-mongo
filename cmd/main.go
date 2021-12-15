
package main

import (
	"os"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"writeon/controllers"	
)


//ConfigureRouter setup the router
func ConfigureRouter() *mux.Router {
	router := mux.NewRouter()

	router.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// router.HandleFunc("/", homeHandler)
	router.HandleFunc("/auth/register", writeon.RegisterUser).Methods("POST")
	router.HandleFunc("/auth/login", writeon.LoginUser).Methods("POST")
	router.HandleFunc("/auth/logout", writeon.LogoutUser).Methods("POST")

	router.HandleFunc("/user/profile", writeon.UserProfile).Methods("GET")

	router.HandleFunc("/topic", writeon.GetAllTopics).Methods("GET")
	router.HandleFunc("/topic", writeon.AddTopic).Methods("POST")
	
	router.HandleFunc("/essay", writeon.GetAllEssays).Methods("GET")
	router.HandleFunc("/essay/topic/{topicID}", writeon.GetAllTopicEssays).Methods("GET")
	router.HandleFunc("/essay/{essayID}", writeon.GetEssay).Methods("GET")
	router.HandleFunc("/essay", writeon.AddEssay).Methods("POST")
	// router.HandleFunc("/essay", writeon.AddEssay).Methods("PUT")
	// router.HandleFunc("/essay", writeon.AddEssay).Methods("DELETE")

	// secure := router.PathPrefix("/auth").Subrouter()
	// secure.Use(auth.JwtVerify)
	// secure.HandleFunc("/api", middleware.ApiHandler).Methods("GET")

	return router
}


func main() {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // All origins
		AllowedMethods: []string{"GET","POST","PUT"}, // Allowing only get, just an example
	})
	//Create an instance of MongoHander with the connection string provided
	mongoDbConnection := "mongodb://localhost:27017"
	writeon.NewMongoHandler(mongoDbConnection) 
	router := ConfigureRouter()
	port := "8088"
	log.Println("Listening on port ", port)
	// handler := cors.Default().Handler(r)
	log.Fatal(http.ListenAndServe(":" + port, handlers.LoggingHandler(os.Stdout, c.Handler(router))))
}
