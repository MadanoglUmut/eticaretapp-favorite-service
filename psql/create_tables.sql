CREATE TABLE FavoriteList(
    id serial PRIMARY KEY NOT NULL,
    listname VARCHAR(100) NOT NULL,
    createddate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    userid INT NOT NULL
);


CREATE TABLE FavoriteItem(
    itemid INT NOT NULL,
    listid INT NOT NULL,
    createddate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (listid, itemid),
    FOREIGN KEY (listid) REFERENCES FavoriteList(id)
);