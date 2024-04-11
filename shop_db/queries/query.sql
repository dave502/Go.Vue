-- name: ListShelves :many
-- get all shelves orderded by name 
SELECT *
FROM shop.shelve
ORDER BY Shelve_Name;


-- name: GetAllOpenedOrders :many
-- get all non finished orders 
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
;

-- name: CreateOrder :one
-- create the requisites of order 


-- name: CreateOrder :one
-- create the requisites of order 
INSERT INTO shop.order (
    User_ID
) VALUES (
    $1
) 
RETURNING Order_ID;

-- name: AddToOrderTable :one
-- create the table part of order 
INSERT INTO shop.order_table (
    Order_ID, Product_ID, Quantity
) VALUES (
    $1, $2, $3
)
RETURNING *
;

-- name: DeleteOrderByOrderID :exec
-- delete order's requisites
DELETE FROM shop.order
WHERE Order_ID = $1;


-- name: DeleteOrderProductsByOrderID :exec
-- clear order's content
DELETE FROM shop.order_table
WHERE Order_ID = $1;


-- name: CreateUsers :one
INSERT INTO shop.user (
    User_Name, Password_Hash, Name
)
VALUES (
    $1, $2, $3
) 
RETURNING *;

-- name: ListUsers :many
SELECT *
FROM shop.user
ORDER BY user_name;

-- name: GetUser :one
SELECT *
FROM shop.user
WHERE user_id = $1;

-- name: GetUserByName :one
SELECT *
FROM shop.user
WHERE user_name = $1;

-- name: DeleteUsers :exec
DELETE
FROM shop.user
WHERE user_id = $1;

