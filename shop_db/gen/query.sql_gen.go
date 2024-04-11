// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: query.sql

package shop_db

import (
	"context"
)

const addToOrderTable = `-- name: AddToOrderTable :one
INSERT INTO shop.order_table (
    Order_ID, Product_ID, Quantity
) VALUES (
    $1, $2, $3
)
RETURNING table_id, order_id, product_id, quantity
`

type AddToOrderTableParams struct {
	OrderID   int64 `db:"order_id" json:"orderID"`
	ProductID int64 `db:"product_id" json:"productID"`
	Quantity  int32 `db:"quantity" json:"quantity"`
}

// create the table part of order
func (q *Queries) AddToOrderTable(ctx context.Context, arg AddToOrderTableParams) (ShopOrderTable, error) {
	row := q.db.QueryRowContext(ctx, addToOrderTable, arg.OrderID, arg.ProductID, arg.Quantity)
	var i ShopOrderTable
	err := row.Scan(
		&i.TableID,
		&i.OrderID,
		&i.ProductID,
		&i.Quantity,
	)
	return i, err
}

const createOrder = `-- name: CreateOrder :one


INSERT INTO shop.order (
    User_ID
) VALUES (
    $1
) 
RETURNING Order_ID
`

// create the requisites of order
// create the requisites of order
func (q *Queries) CreateOrder(ctx context.Context, userID int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, createOrder, userID)
	var order_id int64
	err := row.Scan(&order_id)
	return order_id, err
}

const createUsers = `-- name: CreateUsers :one
INSERT INTO shop.user (
    User_Name, Password_Hash, Name
)
VALUES (
    $1, $2, $3
) 
RETURNING user_id, user_name, password_hash, name, config, created_at, is_enabled
`

type CreateUsersParams struct {
	UserName     string `db:"user_name" json:"userName"`
	PasswordHash string `db:"password_hash" json:"passwordHash"`
	Name         string `db:"name" json:"name"`
}

func (q *Queries) CreateUsers(ctx context.Context, arg CreateUsersParams) (ShopUser, error) {
	row := q.db.QueryRowContext(ctx, createUsers, arg.UserName, arg.PasswordHash, arg.Name)
	var i ShopUser
	err := row.Scan(
		&i.UserID,
		&i.UserName,
		&i.PasswordHash,
		&i.Name,
		&i.Config,
		&i.CreatedAt,
		&i.IsEnabled,
	)
	return i, err
}

const deleteOrderByOrderID = `-- name: DeleteOrderByOrderID :exec
DELETE FROM shop.order
WHERE Order_ID = $1
`

// delete order's requisites
func (q *Queries) DeleteOrderByOrderID(ctx context.Context, orderID int64) error {
	_, err := q.db.ExecContext(ctx, deleteOrderByOrderID, orderID)
	return err
}

const deleteOrderProductsByOrderID = `-- name: DeleteOrderProductsByOrderID :exec
DELETE FROM shop.order_table
WHERE Order_ID = $1
`

// clear order's content
func (q *Queries) DeleteOrderProductsByOrderID(ctx context.Context, orderID int64) error {
	_, err := q.db.ExecContext(ctx, deleteOrderProductsByOrderID, orderID)
	return err
}

const deleteUsers = `-- name: DeleteUsers :exec
DELETE
FROM shop.user
WHERE user_id = $1
`

func (q *Queries) DeleteUsers(ctx context.Context, userID int64) error {
	_, err := q.db.ExecContext(ctx, deleteUsers, userID)
	return err
}

const getAllOpenedOrders = `-- name: GetAllOpenedOrders :many
SELECT  o.Order_ID,  -- order
        o.User_ID, 
        t.Product_ID,  -- order's table
        t.Quantity, 
        p.Product_Name, -- table's products
        p.Product_Price,
        l.Shelve_ID, -- product's location
        l.Is_Main
FROM
        shop.order o,
        shop.order_table t,
        shop.product p,
        shop.product_location l
WHERE   o.Order_ID = t.Order_ID
  AND   t.Product_ID = p.Product_ID
  AND   p.Product_ID = l.Product_ID
`

type GetAllOpenedOrdersRow struct {
	OrderID      int64  `db:"order_id" json:"orderID"`
	UserID       int64  `db:"user_id" json:"userID"`
	ProductID    int64  `db:"product_id" json:"productID"`
	Quantity     int32  `db:"quantity" json:"quantity"`
	ProductName  string `db:"product_name" json:"productName"`
	ProductPrice string `db:"product_price" json:"productPrice"`
	ShelveID     int64  `db:"shelve_id" json:"shelveID"`
	IsMain       bool   `db:"is_main" json:"isMain"`
}

// get all non finished orders
func (q *Queries) GetAllOpenedOrders(ctx context.Context) ([]GetAllOpenedOrdersRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllOpenedOrders)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllOpenedOrdersRow
	for rows.Next() {
		var i GetAllOpenedOrdersRow
		if err := rows.Scan(
			&i.OrderID,
			&i.UserID,
			&i.ProductID,
			&i.Quantity,
			&i.ProductName,
			&i.ProductPrice,
			&i.ShelveID,
			&i.IsMain,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUser = `-- name: GetUser :one
SELECT user_id, user_name, password_hash, name, config, created_at, is_enabled
FROM shop.user
WHERE user_id = $1
`

func (q *Queries) GetUser(ctx context.Context, userID int64) (ShopUser, error) {
	row := q.db.QueryRowContext(ctx, getUser, userID)
	var i ShopUser
	err := row.Scan(
		&i.UserID,
		&i.UserName,
		&i.PasswordHash,
		&i.Name,
		&i.Config,
		&i.CreatedAt,
		&i.IsEnabled,
	)
	return i, err
}

const getUserByName = `-- name: GetUserByName :one
SELECT user_id, user_name, password_hash, name, config, created_at, is_enabled
FROM shop.user
WHERE user_name = $1
`

func (q *Queries) GetUserByName(ctx context.Context, userName string) (ShopUser, error) {
	row := q.db.QueryRowContext(ctx, getUserByName, userName)
	var i ShopUser
	err := row.Scan(
		&i.UserID,
		&i.UserName,
		&i.PasswordHash,
		&i.Name,
		&i.Config,
		&i.CreatedAt,
		&i.IsEnabled,
	)
	return i, err
}

const listShelves = `-- name: ListShelves :many
SELECT shelve_id, shelve_name
FROM shop.shelve
ORDER BY Shelve_Name
`

// get all shelves orderded by name
func (q *Queries) ListShelves(ctx context.Context) ([]ShopShelve, error) {
	rows, err := q.db.QueryContext(ctx, listShelves)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ShopShelve
	for rows.Next() {
		var i ShopShelve
		if err := rows.Scan(&i.ShelveID, &i.ShelveName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUsers = `-- name: ListUsers :many
SELECT user_id, user_name, password_hash, name, config, created_at, is_enabled
FROM shop.user
ORDER BY user_name
`

func (q *Queries) ListUsers(ctx context.Context) ([]ShopUser, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ShopUser
	for rows.Next() {
		var i ShopUser
		if err := rows.Scan(
			&i.UserID,
			&i.UserName,
			&i.PasswordHash,
			&i.Name,
			&i.Config,
			&i.CreatedAt,
			&i.IsEnabled,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}