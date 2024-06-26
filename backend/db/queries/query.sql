-- name: ListShelves :many
-- get all shelves orderded by name 
SELECT *
FROM shop.shelve
ORDER BY Shelve_Name;


-- name: ListProducts :many
-- get all products from database 
SELECT  p.Product_ID, -- table's products
        p.Product_Name,
        p.Product_Price,
        l.Shelve_ID, -- product's location
        l.Is_Main
FROM
        shop.product p,
        shop.product_location l
WHERE   p.Product_ID = l.Product_ID
;

-- name: GetAllOpenedOrdersIds :many
-- get all non finished orders Ids
SELECT  Order_ID
FROM    shop.order
WHERE   shop.order.Is_Closed = FALSE
ORDER BY Order_ID;
;

-- name: GetAllOpenedOrders :many
-- get all non finished orders 
SELECT  o.Order_ID,  -- order
        o.User_ID, 
        t.Product_ID,  -- order's table
        SUM(t.Quantity) as Quantity, 
        -- STRING_AGG
        p.Product_Name, -- table's products
        p.Product_Price
FROM
        shop.order o,
        shop.order_table t,
        shop.product p,
        shop.product_location l
WHERE   o.Is_Closed = FALSE
  AND   o.Order_ID = t.Order_ID
  AND   t.Product_ID = p.Product_ID
GROUP BY o.Order_ID, t.Product_ID, p.Product_Name, Product_Price
ORDER BY o.Order_ID;
;


-- name: GetOrdersProducts :many
-- get all products from non finished orders 
WITH products AS (
        SELECT  t.Product_ID,  
                p.Product_Name, 
                max(p.Product_Price) as price,
                o.Order_ID,
                t.Quantity as quantity,
                --json_object_agg(o.Order_ID, t.Quantity) as orders,
                ARRAY_AGG(s.Shelve_Name) as shelves,
                max(s_main.Shelve_Name) as main_shelve
        FROM    shop.order o
                LEFT JOIN shop.order_table t ON o.Order_ID = t.Order_ID
                LEFT JOIN shop.product p ON t.Product_ID = p.Product_ID
                LEFT JOIN shop.product_location l ON p.Product_ID = l.Product_ID
                LEFT JOIN shop.shelve s ON l.Shelve_ID = s.Shelve_ID
                LEFT JOIN shop.shelve s_main ON l.Shelve_ID = s_main.Shelve_ID AND l.Is_Main=TRUE
        -- WHERE   o.Order_ID IN (...)
        WHERE   o.Order_ID = ANY($1::int[])
        GROUP BY t.Product_ID, p.Product_Name, o.Order_ID, t.Quantity
) 
SELECT 
main_shelve,
JSON_AGG(products) as products
FROM products
GROUP BY main_shelve
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

