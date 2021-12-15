package writeon



// func middleware(next http.Handler) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Println("Executing Auth Middleware")
// 		_, user, err := strategy.AuthenticateRequest(r)
// 		if err != nil {
// 			fmt.Println(err)
// 			code := http.StatusUnauthorized
// 			http.Error(w, http.StatusText(code), code)
// 			return
// 		}
// 		log.Printf("User %s Authenticated\n", user.GetUserName())
// 		r = auth.RequestWithUser(user, r)
// 		next.ServeHTTP(w, r)
// 	})
// }