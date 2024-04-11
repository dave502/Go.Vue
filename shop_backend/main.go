package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"shop/shop_backend/internal"
	"shop/shop_backend/internal/api"
	logger "shop/shop_backend/logs"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {

	l := flag.Bool("local", false, "true - send to stdout, false - send to logging server")
	flag.Parse()
	logger.SetLoggingOutput(*l)
	logger.Logger.Debugf("Application logging to stdout = %v", *l)
	logger.Logger.Info("Starting the application...")
	//logger.Logger.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)

	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		internal.GetAsString("DB_USERNAME", "postgres"),
		internal.GetAsString("DB_PASSWORD", "postgres"),
		internal.GetAsString("DB_HOST", "localhost"),
		internal.GetAsInt("DB_PORT", 5432),
		internal.GetAsString("DB_DATABASE", "db"),
	)

	conn, err := sql.Open("postgres", dbURI)
	if err != nil {
		logger.Logger.Errorf("Error opening database : %s", err.Error())
	}

	// Проверка подключения
	if err := conn.Ping(); err != nil {
		logger.Logger.Errorf("Error opening database : %s", err.Error())
	}

	logger.Logger.Info("Successful connection to database")

	//
	// db := shop_db.New(conn)

	// если используются аргументы - просто возвращаем результат
	if args := os.Args[1:]; len(args) > 1 {
		for _, arg := range args {
			fmt.Println(arg)
		}
		return
	}

	// _, err = db.ListShelves(context.Background()) //,
	// // shop_db.CreateUsersParams{
	// // 	UserName:     "testuser",
	// // 	PassWordHash: "hash",
	// // 	Name:         "test",
	// // })
	// if err != nil {
	// 	logger.Logger.Errorf("Error ... : %s", err.Error())
	// }
	server := api.NewServer(internal.GetAsInt("BACK_SERVER_PORT", 9010))
	server.MustStart()
	defer server.Stop()

	defaultMiddleware := []mux.MiddlewareFunc{
		api.JSONMiddleware,
		api.CORSMiddleware(
			internal.GetAsSlice("ALLOWED_HEADERS",
				[]string{"X-Requested-With", "Origin", "Content-Type"},
				","),
			internal.GetAsSlice("CORS_WHITELIST",
				[]string{
					"http://localhost:9000",
					"http://0.0.0.0:9000",
				}, ","),
		),
	}

	server.AddRoute("/", appGET(), http.MethodGet, defaultMiddleware...)
	server.AddRoute("/", appPOST(), http.MethodPost, defaultMiddleware...)

	// Handlers
	server.AddRoute("/login", handleLogin(conn), http.MethodPost, defaultMiddleware...)
	server.AddRoute("/logout", handleLogout(), http.MethodGet, defaultMiddleware...)

	// Our session protected middleware
	protectedMiddleware := append(defaultMiddleware, api.ValidCookieMiddleware(conn))
	server.AddRoute("/checkSecret", checkSecret(conn), http.MethodGet, protectedMiddleware...)

	// Workouts
	server.AddRoute("/order", handleCreateOrder(conn), http.MethodPost, protectedMiddleware...)
	server.AddRoute("/order", handleListOrders(conn), http.MethodGet, protectedMiddleware...)
	server.AddRoute("/order/{order_id}", handleDeleteOrder(conn), http.MethodDelete, protectedMiddleware...)
	server.AddRoute("/order/{order_id}", handleAddProduct(conn), http.MethodPost, protectedMiddleware...)
	server.AddRoute("/order/{order_id}/{product_id}", handleUpdateProduct(conn), http.MethodPut, protectedMiddleware...)

	// Wait for CTRL-C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	// We block here until a CTRL-C / SigInt is received
	// Once received, we exit and the server is cleaned up
	<-sigChan
}
