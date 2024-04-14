INSERT INTO shop.user (User_ID, User_Name, Name) VALUES 
    (1, 'undefined', 'Незарегистрированный пользователь'),
    (2, 'client', 'Иванов Олег'),
    (3, 'serj123', 'Жилин Сергей Владимирович'),
    (4, 'nata', 'Симова Наталья');
    
INSERT INTO shop.shelve (Shelve_Name) VALUES 
    ('А'),('А'),('А'),('Б'),('Б'),('В'),('Г'),('Д'),('Е'),('Ж'),('З'),('И');
    
INSERT INTO shop.product VALUES 
    (1, 'Ноутбук', 60000.00),
    (2, 'Телевизор', 123000.00),
    (3, 'Телефон', 85000.00),
    (4, 'Системный блок', 333000.00),
    (5, 'Часы', 67000.00),
    (6, 'Микрофон', 32500.00);


INSERT INTO shop.product_location (Product_ID, Shelve_ID, Is_Main) VALUES 
-- основные стеллажи
    (1, 1, TRUE),   -- Ноутбук - A(1)
    (2, 1, TRUE),   -- Телевизор - A(1)
    (3, 4, TRUE),   -- Телефон - Б(1)
    (4, 10, TRUE),  -- Системный блок - Ж
    (5, 10, TRUE),  -- Часы - Ж
    (6, 10, TRUE),  -- Микрофон - Ж
-- дополнительные стеллажи
    (3, 6, FALSE),  -- Телефон - В
    (3, 11, FALSE), -- Телефон - З
    (5, 2, FALSE);  -- Часы - A(2)

-- Заказ № 10
DO $$
DECLARE new_order_id integer;
BEGIN
  INSERT INTO shop.order (Order_ID, User_ID) VALUES (10, 2) 
    RETURNING Order_ID INTO new_order_id;
  INSERT INTO shop.order_table (Order_ID, Product_ID, Quantity) VALUES 
    (new_order_id, 1, 2), -- Ноутбук - 2 шт.
    (new_order_id, 3, 1), -- Телефон - 1 шт.
    (new_order_id, 6, 1)  -- Микрофон - 1 шт.
    ;
END $$;

-- Заказ № 11
DO $$
DECLARE new_order_id integer;
BEGIN
  INSERT INTO shop.order (Order_ID, User_ID) VALUES (11, 1) 
    RETURNING Order_ID INTO new_order_id;
  INSERT INTO shop.order_table (Order_ID, Product_ID, Quantity) VALUES 
    (new_order_id, 2, 3) -- Телевизор - 3 шт.
    ;
END $$;

-- Заказ № 14
DO $$
DECLARE new_order_id integer;
BEGIN
  INSERT INTO shop.order (Order_ID, User_ID) VALUES (14, 4) 
    RETURNING Order_ID INTO new_order_id;
  INSERT INTO shop.order_table (Order_ID, Product_ID, Quantity) VALUES 
    (new_order_id, 1, 3), -- Ноутбук - 3 шт.
    (new_order_id, 4, 4)  -- Системный блок - 4 шт.
    ;
END $$;

-- Заказ № 15
DO $$
DECLARE new_order_id integer;
BEGIN
  INSERT INTO shop.order (Order_ID, User_ID) VALUES (15, 3) 
    RETURNING Order_ID INTO new_order_id;
  INSERT INTO shop.order_table (Order_ID, Product_ID, Quantity) VALUES 
    (new_order_id, 5, 1) -- Часы - 1 шт.
    ;
END $$;
