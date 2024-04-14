
        SELECT  t.Product_ID,  
                p.Product_Name, 
                max(p.Product_Price) as price,
                o.Order_ID,
                t.Quantity as quantity,
                --jsonb_object_agg(o.Order_ID, t.Quantity) as orders,
                ARRAY_AGG(s.Shelve_Name) as shelves, -- ORDER BY s.Shelve_Name
                max(s_main.Shelve_Name) as main_shelve
        FROM    shop.order o
                LEFT JOIN shop.order_table t ON o.Order_ID = t.Order_ID
                LEFT JOIN shop.product p ON t.Product_ID = p.Product_ID
                LEFT JOIN shop.product_location l ON p.Product_ID = l.Product_ID
                LEFT JOIN shop.shelve s ON l.Shelve_ID = s.Shelve_ID
                LEFT JOIN shop.shelve s_main ON l.Shelve_ID = s_main.Shelve_ID AND l.Is_Main=TRUE
        WHERE   o.Order_ID IN (10, 11, 14, 15)
        GROUP BY t.Product_ID, p.Product_Name, o.Order_ID, t.Quantity

;
