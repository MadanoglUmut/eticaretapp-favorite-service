INSERT INTO FavoriteList (listname, userid) VALUES
('Alışveriş Listesi', 1),
('Film Listesi', 2),
('Okunacak Kitaplar', 3),
('Seyahat Planları', 4),
('Fitness Hedefleri', 5);

INSERT INTO FavoriteItem (itemid,listid) VALUES
(1, 1), 
(2, 1), 
(3, 2), 
(4, 3), 
(1, 5),
(5, 4);