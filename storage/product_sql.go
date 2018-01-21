package storage

const createTableProductStmt = `
CREATE TABLE IF NOT EXISTS product (
  id     INTEGER PRIMARY KEY AUTOINCREMENT
, name   varchar(50)
, desc   TEXT
, attributes TEXT
, images  TEXT
, created timestamp
);
`

const insertProductStmt = `
INSERT INTO product (
  id
, name
, desc
, attributes
, images
, created
) VALUES (?,?,?,?,?,?)
`

const selectProductsStmt = `
SELECT
  id
, name
, desc
, attributes
, images
, created
FROM users
`

const selectProductByIDStmt = `
SELECT
  id
, name
, desc
, attributes
, images
, created
FROM product
WHERE id=?
`
