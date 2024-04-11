package api

import (
	"context"
	"database/sql"
	"mime"
	"net/http"
	dbconn "shop/shop_db/gen"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/sessions"
)


// ***** COOKIE VARS *****
type UserSession struct {
	UserID int64
}

// We define this so it can't clash outside our package with anything else.
type customKey string
const sessionKey customKey = "unique-session-key"
var CookieStore = sessions.NewCookieStore([]byte("DemoShop"))

func init() {
	CookieStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
	}
}


// JSON middleware will help us only handle JSON
// in and out
func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		contentType := req.Header.Get("Content-Type")

		if strings.TrimSpace(contentType) == "" {
			var parseError error
			contentType, _, parseError = mime.ParseMediaType(contentType)
			if parseError != nil {
				JSONError(wr, http.StatusBadRequest, "Bad or no content-type header found")
				return
			}
		}

		if contentType != "application/json" {
			JSONError(wr, http.StatusUnsupportedMediaType, "Content-Type not application/json")
			return
		}
		// Tell the client we're talking JSON as well.
		wr.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(wr, req)
	})
}

func CORSMiddleware(headers []string, origins []string) func(http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{
			"X-Requested-With", "Origin", "Content-Type",
			// "Authorization", "Access-Control-Allow-Methods",
			// "Access-Control-Allow-Origin", , "Accept",
		}),
		handlers.AllowedOrigins(origins),
		handlers.AllowCredentials(),
		handlers.AllowedMethods([]string{
			http.MethodPost,
			//http.MethodPatch,
			//http.MethodPut,
			http.MethodGet,
			http.MethodDelete,
		}),
	)
}

func ValidCookieMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
			session, err := CookieStore.Get(req, "session-name")
			if err != nil {
				JSONError(wr, http.StatusInternalServerError, "Session Error")
				return
			}

			userID, userIDOK := session.Values["userID"].(int64)
			isAuthd, isAuthdOK := session.Values["userAuthenticated"].(bool)

			// We could put with the above but lets keep our logic simple
			if !userIDOK || !isAuthdOK {
				JSONError(wr, http.StatusInternalServerError, "Session Error")
				return
			}

			if !isAuthd || userID < 1 {
				JSONError(wr, http.StatusForbidden, "Bad Credentials")
				return
			}

			query := dbconn.New(db)
			user, err := query.GetUser(req.Context(), int64(userID))
			if err != nil || user.UserID < 1 {
				JSONError(wr, http.StatusForbidden, "Bad Credentials")
				return
			}

			ctx := context.WithValue(req.Context(), sessionKey, UserSession{
				UserID: user.UserID,
			})
			h.ServeHTTP(wr, req.WithContext(ctx))
		})
	}
}


// JWTProtectedMiddleware verifies a valid JWT exists in
// our cookie and if not, encourages the consumer to login
// again.
func JWTProtectedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	    // Grab jwt-token cookie
	    jwtCookie, err := r.Cookie("jwt-token")
	    if err != nil {
			log.Println("Error occurred reading cookie", err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(
				struct { Message string `json:"message,omitempty"` }
				{ Message: "Your session is not valid - please login" }
			)
	      	return
	    }
	    // Decode and validate JWT if there is one
	    userEmail, err := decodeJWTToUser(jwtCookie.Value)
	    if userEmail == "" || err != nil {
	    	log.Println("Error decoding token", err)
	      	w.WriteHeader(http.StatusUnauthorized)
	      	json.NewEncoder(w).Encode(
				struct  { Message string `json:"message,omitempty"` }
						{ Message: "Your session is not valid - please login" }
		  	)
		  	return
		}
		// If it's good, update the expiry time
		freshToken := createJWTTokenForUser(userEmail)
		// Set the new cookie and continue into the handler
		w.Header().Add("Content-Type", "application/json")
		http.SetCookie(w, authCookie(freshToken))
		// continue and indicate the successful handling of a request
		next.ServeHTTP(w, r)
	})
}			

func UserFromSession(req *http.Request) (UserSession, bool) {
	session, ok := req.Context().Value(sessionKey).(UserSession)
	if session.UserID < 1 {
		// Shouldnt happen
		return UserSession{}, false
	}
	return session, ok
}
