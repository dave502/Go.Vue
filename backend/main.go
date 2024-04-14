package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	dbconn "shop/db/gen"
	"shop/internal/api"
	"shop/internal/env"
	logger "shop/logs"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {

	var conn *sql.DB
	var err error

	l := flag.Bool("local", true, "true - send to stdout, false - send to logging server")
	flag.Parse()
	logger.SetLoggingOutput(*l)
	logger.Logger.Debugf("Application logging to stdout = %v", *l)
	logger.Logger.Info("Starting the application...")
	//logger.Logger.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)

	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		env.GetAsString("DB_USERNAME", "root"),
		env.GetAsString("DB_PASSWORD", "root"),
		env.GetAsString("DB_HOST", "localhost"),
		env.GetAsInt("DB_PORT", 5432),
		env.GetAsString("DB_DATABASE", "shop_db"),
	)

	for {
		conn, err = sql.Open("postgres", dbURI)
		if err != nil {
			logger.Logger.Errorf("Error opening database : %s", err.Error())
		}
		// Проверка подключения
		if err = conn.Ping(); err == nil {
			logger.Logger.Info("Successful connection to database")
			break
		} else {
			logger.Logger.Info("Keep trying...")
		}
		time.Sleep(2 * time.Second)
		continue
	}

	// если используются аргументы - просто возвращаем результат
	if str_args := os.Args[1:]; len(str_args) > 1 {

		// структура получаемого товара
		type Product struct {
			ProductId   int      `json:"product_id"`
			ProductName string   `json:"product_name"`
			OrderId     int      `json:"order_id"`
			Quantity    int      `json:"quantity"`
			Price       float32  `json:"price"`
			Shelves     []string `json:"shelves"`
		}

		// конвертация строковых аргументов в числовые
		len_args := len(str_args)
		int_nums := make([]int32, len_args)
		for i := 0; i < len_args; i++ {
			if int_num, err := strconv.Atoi(strings.Trim(str_args[i], ",")); err != nil {
				logger.Logger.Error(err)
			} else {
				int_nums[i] = int32(int_num)
			}
		}

		// выполнение запроса к базе данных
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*100))
		defer cancel()
		query := dbconn.New(conn)
		orders, err := query.GetOrdersProducts(ctx, int_nums)
		if err != nil {
			logger.Logger.Errorf("Ошибка получения данных", orders)
			return
		}

		// вывод результата в терминал
		// заголовок
		if len(orders) > 0 {
			fmt.Printf("=+=+=+=\nСтраница сборки заказов %s \n\n", strings.Join(str_args, ""))
		}

		// для каждого стеллажа
		for _, shelve := range orders {
			fmt.Println("===Стеллаж", shelve.MainShelve)
			// получение списка товаров из json
			var p []Product
			err := json.Unmarshal([]byte(shelve.Products), &p)
			if err != nil {
				logger.Logger.Errorf("Неправильный формат данных %v", err)
				return
			}
			// вывод каждого товара со стеллажа
			for _, product := range p {

				fmt.Printf("%s (id=%d)\n", product.ProductName, product.ProductId)
				fmt.Printf("Заказ %d, %d шт.\n", product.OrderId, product.Quantity)

				// удаление из списка стеллажей главного, чтобы получить список доп. стеллажей
				main_pos := slices.Index(product.Shelves, shelve.MainShelve.(string))
				shelves := slices.Delete(product.Shelves, main_pos, main_pos+1)

				if len(shelves) > 0 {
					fmt.Printf("Доп. стеллаж: %s\n", strings.Join(shelves, ", "))
				}
				fmt.Print("\n")
			}
		}
		return

	}

	server := api.NewServer(env.GetAsInt("BACK_SERVER_PORT", 9010))
	server.MustStart()
	defer server.Stop()

	defaultMiddleware := []mux.MiddlewareFunc{
		//api.JSONMiddleware,
		api.CORSMiddleware(
			env.GetAsSlice("ALLOWED_HEADERS",
				[]string{"X-Requested-With", "Origin", "Content-Type", "Accept"},
				","),
			env.GetAsSlice("CORS_WHITELIST",
				[]string{
					"http://localhost:9000",
					"http://0.0.0.0:9000",
					"http://localhost:5173",
					"http://0.0.0.0:5173",
				}, ","),
		),
	}

	// Handlers
	server.AddRoute("/login", handleLogin(conn), http.MethodPost, defaultMiddleware...)
	server.AddRoute("/logout", handleLogout(), http.MethodGet, defaultMiddleware...)

	// Our session protected middleware
	protectedMiddleware := append(defaultMiddleware, api.ValidCookieMiddleware(conn))
	server.AddRoute("/checkSecret", checkSecret(conn), http.MethodGet, protectedMiddleware...)

	// Products
	server.AddRoute("/product", handleListProducts(conn), http.MethodGet, defaultMiddleware...)
	// Orders
	server.AddRoute("/order", handleCreateOrder(conn), http.MethodPost, protectedMiddleware...)
	server.AddRoute("/order", handleListOrders(conn), http.MethodGet, defaultMiddleware...)        // protectedMiddleware
	server.AddRoute("/order-ids", handleListOrdersIds(conn), http.MethodGet, defaultMiddleware...) // protectedMiddleware
	server.AddRoute("/order/{order_id}", handleDeleteOrder(conn), http.MethodDelete, protectedMiddleware...)
	server.AddRoute("/order/{order_id}/{product_id}", handleUpdateProduct(conn), http.MethodPut, protectedMiddleware...)
	server.AddRoute("/order/{order_id}/{product_id}", handleAddProduct(conn), http.MethodPost, protectedMiddleware...)
	server.AddRoute("/packing-list", handleGetOrdersProducts(conn), http.MethodPost, defaultMiddleware...) // protectedMiddleware

	// Завершение по CTRL-C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}
