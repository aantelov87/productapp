package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	app "github.com/productapp"
)

var _ app.ProductStore = &mysqlDB{}

type mysqlDB struct {
	*sql.DB
	get    *sql.Stmt
	list   *sql.Stmt
	insert *sql.Stmt
}

func NewMysqlDB(db *sql.DB) (app.ProductStore, error) {
	m := mysqlDB{DB: db}
	var err error
	if m.list, err = db.Prepare(selectProductsStmt); err != nil {
		return nil, fmt.Errorf("mysql: prepare list: %v", err)
	}
	if m.get, err = db.Prepare(selectProductByIDStmt); err != nil {
		return nil, fmt.Errorf("mysql: prepare get: %v", err)
	}
	if m.insert, err = db.Prepare(insertProductStmt); err != nil {
		return nil, fmt.Errorf("mysql: prepare insert: %v", err)
	}
	return &m, nil
}

func (db *mysqlDB) List() ([]*app.Product, error) {
	rows, err := db.list.Query()
	if err != nil {
		return nil, err
	}
	var pl []*app.Product
	for rows.Next() {
		p, err := ScanProduct(rows)
		if err != nil {
			return nil, err
		}
		pl = append(pl, p)
	}
	return pl, nil
}

func (db *mysqlDB) Lookup(id string) (*app.Product, error) {
	row := db.get.QueryRow(id)
	return ScanProduct(row)
}

func (db *mysqlDB) Create(p *app.Product) error {
	_, err := db.insert.Exec(&p.ID,
		&p.Name,
		&p.Desc,
		&p.Attributes,
		&p.Images,
		&p.Created,
	)
	return err
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func ScanProduct(row rowScanner) (*app.Product, error) {
	var id string
	var name string
	var desc string
	var attributes []byte
	var images []byte
	var created time.Time

	err := row.Scan(
		&id,
		&name,
		&desc,
		&attributes,
		&images,
		&created,
	)
	if err != nil {
		return nil, err
	}
	p := &app.Product{}
	p.ID = id
	p.Name = name
	p.Desc = desc
	if err := json.Unmarshal(attributes, p.Attributes); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(images, p.Images); err != nil {
		return nil, err
	}
	p.Created = created
	return p, nil
}
