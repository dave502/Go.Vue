package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	dbconn "shop/db/gen"
	logger "shop/logs"

	"shop/internal/api"
	"shop/internal/crypto"
	"strconv"

	"github.com/gorilla/mux"
)

func handleLogin(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {

		// Thanks to our middleware, we know we have JSON
		// we'll decode it into our request type and see if it's valid
		type loginRequest struct {
			Username string `json:"username,omitempty"`
			Password string `json:"password,omitempty"`
		}

		payload := loginRequest{}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			log.Println("Error decoding the body", err)
			api.JSONError(wr, http.StatusBadRequest, "Error decoding JSON")
			return
		}

		querier := dbconn.New(db)
		user, err := querier.GetUserByName(req.Context(), payload.Username)
		if errors.Is(err, sql.ErrNoRows) || !crypto.CheckPasswordHash(payload.Password, user.PasswordHash) {
			api.JSONError(wr, http.StatusForbidden, "Bad Credentials")
			return
		}
		if err != nil {
			log.Println("Received error looking up user", err)
			api.JSONError(wr, http.StatusInternalServerError, "Couldn't log you in due to a server error")
			return
		}

		// We're valid. Let's tell the user and set a cookie
		session, err := api.CookieStore.Get(req, "session-name")
		if err != nil {
			log.Println("Cookie store failed with", err)
			api.JSONError(wr, http.StatusInternalServerError, "Session Error")
		}
		session.Values["userAuthenticated"] = true
		session.Values["userID"] = user.UserID
		session.Save(req, wr)
	})
}

func checkSecret(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		userDetails, _ := api.UserFromSession(req)

		query := dbconn.New(db)
		user, err := query.GetUser(req.Context(), userDetails.UserID)
		if errors.Is(err, sql.ErrNoRows) {
			api.JSONError(wr, http.StatusForbidden, "User not found")
			return
		}

		api.JSONMessage(wr, http.StatusOK, fmt.Sprintf("Hello there %s", user.Name))
	})
}

func handleLogout() http.HandlerFunc {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		session, err := api.CookieStore.Get(req, "session-name")
		if err != nil {
			log.Println("Cookie store failed with", err)
			api.JSONError(wr, http.StatusInternalServerError, "Session Error")
			return
		}

		session.Options.MaxAge = -1 // deletes
		session.Values["userID"] = int64(-1)
		session.Values["userAuthenticated"] = false

		err = session.Save(req, wr)
		if err != nil {
			api.JSONError(wr, http.StatusInternalServerError, "Session Error")
			return
		}

		api.JSONMessage(wr, http.StatusOK, "logout successful")
	})
}

func handleListProducts(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {

		logger.Logger.Debugf("run handleListProducts")

		query := dbconn.New(db)
		products, err := query.ListProducts(req.Context())
		if err != nil {
			logger.Logger.Errorf("Failed to get products ", err.Error())
			api.JSONError(wr, http.StatusInternalServerError, err.Error())
			return
		}
		json.NewEncoder(wr).Encode(&products)
		logger.Logger.Debugf("Got something like products: ", products)
	})
}

func handleCreateOrder(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		userDetails, ok := api.UserFromSession(req)
		if !ok {
			api.JSONError(wr, http.StatusForbidden, "Bad context")
			return
		}
		query := dbconn.New(db)

		res, err := query.CreateOrder(req.Context(), userDetails.UserID)
		if err != nil {
			api.JSONError(wr, http.StatusInternalServerError, err.Error())
			return
		}

		json.NewEncoder(wr).Encode(&res)

	})
}

func handleListOrders(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {

		query := dbconn.New(db)
		orders, err := query.GetAllOpenedOrders(req.Context())
		if err != nil {
			api.JSONError(wr, http.StatusInternalServerError, err.Error())
			return
		}
		json.NewEncoder(wr).Encode(&orders)
	})
}

func handleListOrdersIds(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {

		query := dbconn.New(db)
		orders, err := query.GetAllOpenedOrdersIds(req.Context())
		if err != nil {
			api.JSONError(wr, http.StatusInternalServerError, err.Error())
			return
		}
		json.NewEncoder(wr).Encode(&orders)
	})
}

func handleGetOrdersProducts(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {

		type OrderIds struct {
			Ids []int32 `json:"ids,omitempty"`
		}

		payload := OrderIds{}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			log.Println("Error decoding the body", err)
			api.JSONError(wr, http.StatusBadRequest, "Error decoding JSON")
			return
		}

		query := dbconn.New(db)
		orders, err := query.GetOrdersProducts(req.Context(), payload.Ids)
		if err != nil {
			api.JSONError(wr, http.StatusInternalServerError, err.Error())
			return
		}
		json.NewEncoder(wr).Encode(&orders)
	})
}

func handleAddProduct(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {

		// ProductID, err := strconv.Atoi(mux.Vars(req)["product_id"])
		// if err != nil {
		// 	api.JSONError(wr, http.StatusBadRequest, "Bad product_id")
		// 	return
		// }

		type ProductToOrderRequest struct {
			Order_ID   int64 `json:"order_id"`
			Product_ID int64 `json:"product_id"`
			Quantity   int32 `json:"count,omitempty"`
		}

		payload := ProductToOrderRequest{}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			log.Println("Error decoding the body", err)
			api.JSONError(wr, http.StatusBadRequest, "Error decoding JSON")
			return
		}

		query := dbconn.New(db)

		set, err := query.AddToOrderTable(req.Context(),
			dbconn.AddToOrderTableParams{
				OrderID:   payload.Order_ID,
				ProductID: payload.Product_ID,
				Quantity:  payload.Quantity,
			})
		if err != nil {
			api.JSONError(wr, http.StatusInternalServerError, err.Error())
			return
		}
		json.NewEncoder(wr).Encode(&set)
	})
}

func handleUpdateProduct(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		// TODO
	})
}

func handleDeleteOrder(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		// userDetails, ok := api.UserFromSession(req)
		// if !ok {
		// 	api.JSONError(wr, http.StatusForbidden, "Bad context")
		// 	return
		// }

		orderID, err := strconv.Atoi(mux.Vars(req)["order_id"])
		if err != nil {
			api.JSONError(wr, http.StatusBadRequest, "Bad order_id")
			return
		}

		query := dbconn.New(db)

		err = query.DeleteOrderByOrderID(req.Context(), int64(orderID))
		if err != nil {
			api.JSONError(wr, http.StatusBadRequest, "Bad order_id")
			return
		}

		err = query.DeleteOrderProductsByOrderID(req.Context(), int64(orderID))
		if err != nil {
			api.JSONError(wr, http.StatusBadRequest, "Bad order_id")
			return
		}

		api.JSONMessage(wr, http.StatusOK, fmt.Sprintf("order %d is deleted", orderID))
	})
}

func appGET() http.HandlerFunc {
	type ResponseBody struct {
		Message string
	}
	return func(rw http.ResponseWriter, req *http.Request) {
		log.Println("GET", req)
		json.NewEncoder(rw).Encode(ResponseBody{
			Message: "Hello World",
		})
	}
}

func appPOST() http.HandlerFunc {
	type RequestBody struct {
		Inbound string
	}
	type ResponseBody struct {
		OutBound string
	}
	return func(rw http.ResponseWriter, req *http.Request) {
		log.Println("POST", req)

		var rb RequestBody
		if err := json.NewDecoder(req.Body).Decode(&rb); err != nil {
			log.Println("apiAdminPatchUser: Decode failed:", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Println("We received an inbound value of", rb.Inbound)
		json.NewEncoder(rw).Encode(ResponseBody{
			OutBound: "got" + rb.Inbound,
		})
	}
}
