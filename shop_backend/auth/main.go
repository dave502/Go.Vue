package auth

import (
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("Server Version :", Version)
	initDatabase()

	initRedis()

	router := mux.NewRouter()

	//POST handler for /login
	router.HandleFunc("/login", loginHandler).Methods("POST")

	//embed handler for /css path
	ccssontentStatic, _ := fs.Sub(cssEmbed, "css")
	css := http.FileServer(http.FS(ccssontentStatic))
	router.PathPrefix("/css").Handler(http.StripPrefix("/css", css))

	//embed handler for /app path
	contentStatic, _ := fs.Sub(staticEmbed, "static")
	static := http.FileServer(http.FS(contentStatic))
	router.PathPrefix("/app").Handler(securityMiddleware(http.StripPrefix("/app", static)))

	//add /login path
	router.PathPrefix("/login").Handler(securityMiddleware(http.StripPrefix("/login", static)))

	//add /logout path
	router.HandleFunc("/logout", logoutHandler).Methods("GET")

	//root will redirect to /apo
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app", http.StatusPermanentRedirect)
	})

	// Use our basicMiddleware
	router.Use(basicMiddleware)

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:3333",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
