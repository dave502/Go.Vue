package auth

// securityMiddleware is middleware to make sure all request has a
// valid session and authenticated
func securityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//first of all all request MUST have a valid session
		if sessionValid(w, r) {
			//login path will be let through, otherwise it won't be served
			//to the front end
			if r.URL.Path == "/login" {
				next.ServeHTTP(w, r)
				return
			}
		}

		//if it does have a valid session make sure it has been authenticated
		if hasBeenAuthenticated(w, r) {
			next.ServeHTTP(w, r)
			return
		}

		//otherwise it will need to be redirected to /login
		storeAuthenticated(w, r, false)
		http.Redirect(w, r, "/login", 307)
	})
}

// sessionValid check whether the session is a valid session
func sessionValid(w http.ResponseWriter, r *http.Request) bool {
	session, _ := store.Get(r, "session_token")
	return !session.IsNew
}

// logoutHandler handles logout operation
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if hasBeenAuthenticated(w, r) {
		session, _ := store.Get(r, "session_token")
		session.Options.MaxAge = -1
		err := session.Save(r, w)
		if err != nil {
			log.Println("failed to delete session", err)
		}
	}

	http.Redirect(w, r, "/login", 307)
}

// loginHandler handles authentication
func loginHandler(w http.ResponseWriter, r *http.Request) {
	result := "Login "
	r.ParseForm()

	if validateUser(r.FormValue("username"), r.FormValue("password")) {
		storeAuthenticated(w, r, true)
		result = result + "successfull"
	} else {
		result = result + "unsuccessful"
	}

	renderFiles("msg", w, result)
}

// hasBeenAuthenticated checks whether the session contain the flag to indicate
// that the session has gone through authentication process
func hasBeenAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	session, _ := store.Get(r, "session_token")
	a, _ := session.Values["authenticated"]

	if a == nil {
		return false
	}

	return a.(bool)
}

// storeAuthenticated to store authenticated value
func storeAuthenticated(w http.ResponseWriter, r *http.Request, v bool) {
	session, _ := store.Get(r, "session_token")

	session.Values["authenticated"] = v
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// validateUser check whether username/password exist in database
func validateUser(username string, password string) bool {
	//query the data from database
	ctx := context.Background()
	u, _ := dbQuery.GetUserByName(ctx, username)

	//username does not exist
	if u.UserName != username {
		return false
	}

	return pkg.CheckPasswordHash(password, u.PassWordHash)
}

func basicMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		log.Println("Middleware called on", req.URL.Path)
		// do stuff
		h.ServeHTTP(wr, req)
	})
}
