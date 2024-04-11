CREATE SCHEMA IF NOT EXISTS shop;

-- ************************************** Клиенты

CREATE TABLE shop.user
(
    User_ID        BIGSERIAL NOT NULL,
    User_Name      text NOT NULL UNIQUE,
    Password_Hash text NOT NULL DEFAULT '$2y$10$khZ.MgCgwhbjLA1MyfuIoOyO0BMeh0CpoqeVu0M5JRE2swwNca4jW',
    Name           text NOT NULL,
    Config         JSONB NOT NULL DEFAULT '{}'::JSONB,
    Created_At     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    Is_Enabled     BOOLEAN NOT NULL DEFAULT TRUE,
    CONSTRAINT PK_user PRIMARY KEY (User_ID),
    CONSTRAINT unique_user_name UNIQUE (User_Name)
);


-- ************************************** Стеллажи

CREATE TABLE shop.shelve
(
    Shelve_ID   bigserial NOT NULL,
    Shelve_Name text NOT NULL,
    CONSTRAINT PK_shelves PRIMARY KEY ( Shelve_ID )
);

-- ************************************** Товар

CREATE TABLE shop.product
(
    Product_ID  bigserial NOT NULL,
    Product_Name text NOT NULL,
    Product_Price numeric(8,2) NOT NULL,
    CONSTRAINT PK_product PRIMARY KEY ( Product_ID ),
    CONSTRAINT check_positive_price CHECK (Product_Price > 0),
    CONSTRAINT unique_product_name UNIQUE (Product_Name)
);


-- ************************************** Заказы

CREATE TABLE shop.order
(
    Order_ID     bigserial NOT NULL,
    User_ID      bigserial NOT NULL,
    Order_Date   timestamp NOT NULL DEFAULT NOW(),
    Is_Closed     BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT PK_order PRIMARY KEY ( Order_ID ),
    CONSTRAINT FK_user FOREIGN KEY ( User_ID ) REFERENCES shop.user ( User_ID )
);

CREATE INDEX index_order_on_user_id ON shop.order (User_ID);


-- ************************************** Табличная часть заказа

CREATE TABLE shop.order_table
(
    Table_ID      bigserial NOT NULL,
    Order_ID 	  bigserial NOT NULL,
    Product_ID 	  bigserial NOT NULL,
    Quantity      int NOT NULL DEFAULT 1,
    CONSTRAINT PK_table PRIMARY KEY ( Table_ID, Order_ID ),
    CONSTRAINT FK_order FOREIGN KEY ( Order_ID ) REFERENCES shop.order ( Order_ID ),
    CONSTRAINT FK_product FOREIGN KEY ( Product_ID ) REFERENCES shop.product ( Product_ID )
);

CREATE INDEX index_order_table_on_order_id ON shop.order_table ( Order_ID );



-- ************************************** Расположение товара

CREATE TABLE shop.product_location
(
    Location_ID  bigserial NOT NULL,
    Product_ID  bigserial NOT NULL,
    Shelve_ID   bigserial NOT NULL,
    Is_Main     boolean NOT NULL DEFAULT TRUE
);
-- https://dba.stackexchange.com/questions/197562/constraint-one-boolean-row-is-true-all-other-rows-false
CREATE UNIQUE INDEX only_one_shelve_is_main 
	ON shop.product_location (Product_ID, Is_Main) WHERE Is_Main;
CREATE INDEX index_product_location_on_product_id ON shop.product_location ( Product_ID );

