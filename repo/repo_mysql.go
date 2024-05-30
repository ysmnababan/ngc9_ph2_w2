package repo

import (
	"database/sql"
	"fmt"
	"log"
	"ngc9/errhandler"
	"ngc9/model"

	"golang.org/x/crypto/bcrypt"
)

type ProductRepo interface {
	GetAllProducts() (interface{}, error)
	GetProductById(id uint) (interface{}, error)
	CreateProduct(p model.ProductDB) (interface{}, error)
	UpdateProduct(id int, p model.ProductDB) error
	DeleteProduct(id int) error
}

type UserRepo interface {
	Login(u model.User) (model.User, error)
	Register(u model.User) (model.User, error)
}

func (r *MysqlRepo) Register(u model.User) (model.User, error) {
	var isExist bool

	err := r.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email = ?)", u.Email).Scan(&isExist)
	if err != nil {
		return model.User{}, errhandler.ErrQuery
	}

	if isExist {
		return model.User{}, errhandler.ErrUserExists
	}

	hashedpwd, _ := bcrypt.GenerateFromPassword([]byte(u.Pwd), bcrypt.DefaultCost)
	result, err := r.DB.Exec("INSERT INTO users (name, email, pwd) VALUES (?,?,?)", u.Name, u.Email, hashedpwd)
	if err != nil {
		return model.User{}, errhandler.ErrQuery
	}

	id, err := result.LastInsertId()
	if err != nil {
		return model.User{}, errhandler.ErrLastInsertId
	}

	u.ID = uint(id)
	return u, nil
}

func (r *MysqlRepo) Login(u model.User) (model.User, error) {
	var isExist bool
	var user model.User
	err := r.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email = ?)", u.Email).Scan(&isExist)
	if err != nil {
		return model.User{}, errhandler.ErrQuery
	}

	if !isExist {
		return model.User{}, errhandler.ErrNoRows
	}

	err = r.DB.QueryRow("SELECT id, name, pwd FROM users WHERE email = ?", u.Email).Scan(&user.ID, &user.Name, &user.Pwd)
	if err != nil {
		return model.User{}, errhandler.ErrQuery
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Pwd), []byte(u.Pwd))
	if err != nil {
		return model.User{}, errhandler.ErrCredential
	}

	return user, nil
}

type MysqlRepo struct {
	DB *sql.DB
}

func (r *MysqlRepo) IsIDExist(id int) (bool, error) {
	var isExist bool
	err := r.DB.QueryRow("SELECT EXISTS (SELECT 1 from products WHERE id = ?)", id).Scan(&isExist)
	if err != nil {
		log.Println("error querying", err)
		return false, err
	}

	return isExist, nil
}

func (r *MysqlRepo) GetAllProducts() ([]model.Product, error) {
	var products []model.Product
	query := "SELECT id, name, description, img, price, store_name FROM products JOIN stores ON products.store_id = stores.store_id"
	rows, err := r.DB.Query(query)
	if err != nil {
		log.Println("error query", err)
		return nil, errhandler.ErrQuery
	}

	defer rows.Close()

	for rows.Next() {
		var p model.Product
		err := rows.Scan(&p.Id, &p.Name, &p.Desc, &p.Img, &p.Price, &p.Store)
		if err != nil {
			log.Println("error scan row")
			return nil, errhandler.ErrScan
		}

		products = append(products, p)
	}

	if len(products) == 0 {
		log.Println("empty table")
		return nil, errhandler.ErrNoRows
	}

	return products, nil
}

func (r *MysqlRepo) GetProductById(id int) (model.Product, error) {
	var p model.Product
	isExist, err := r.IsIDExist(id)
	if err != nil {
		log.Println("error query", err)
		return model.Product{}, errhandler.ErrQuery
	}

	if !isExist {
		log.Println("id not found")
		return model.Product{}, errhandler.ErrNoRows
	}
	fmt.Println("here")
	query := "SELECT name, description, img, price, store_name FROM products JOIN stores ON products.store_id = stores.store_id WHERE id = ?"
	err = r.DB.QueryRow(query, id).Scan(&p.Name, &p.Desc, &p.Img, &p.Price, &p.Store)
	if err != nil {
		log.Println("error query")
		return model.Product{}, errhandler.ErrQuery
	}

	p.Id = int(id)
	return p, nil
}

func (r *MysqlRepo) CreateProduct(p model.Product) (model.Product, error) {

	query := "INSERT INTO products (name, description, img, price, store_id) VALUES (?,?,?,?,(SELECT store_id FROM stores WHERE store_name = ?))"

	fmt.Println(p)
	result, err := r.DB.Exec(query, p.Name, p.Desc, p.Img, p.Price, p.Store)
	if err != nil {
		log.Println("error query", err)
		return model.Product{}, errhandler.ErrQuery
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		log.Println("error getting last inserted id")
		return model.Product{}, errhandler.ErrLastInsertId
	}

	p.Id = int(lastId)
	return p, nil
}

func (r *MysqlRepo) UpdateProduct(id int, p model.Product) error {
	isExist, err := r.IsIDExist(id)
	if err != nil {
		log.Println("error query", err)
		return errhandler.ErrQuery
	}

	if !isExist {
		log.Println("id not found")
		return errhandler.ErrNoRows
	}

	query := `
	UPDATE products
	SET name =?,
	description = ?, 
	img = ?,
	price = ?,
	store_id = (SELECT store_id FROM stores WHERE store_name = ?)
	WHERE id = ?
	`

	result, err := r.DB.Exec(query, p.Name, p.Desc, p.Img, p.Price, p.Store, id)
	if err != nil {
		log.Println("error query", err)
		return errhandler.ErrQuery
	}

	affectedRow, err := result.RowsAffected()
	if err != nil {
		log.Println("error getting num of rows affected")
		return errhandler.ErrRowsAffected
	}

	if affectedRow == 0 {
		log.Println("")
		return errhandler.ErrNoUpdate
	}
	return nil
}

func (r *MysqlRepo) DeleteProduct(id int) error {
	isExist, err := r.IsIDExist(id)
	if err != nil {
		log.Println("error query", err)
		return errhandler.ErrQuery
	}

	if !isExist {
		log.Println("id not found")
		return errhandler.ErrNoRows
	}

	query := "DELETE FROM products WHERE id = ? "
	result, err := r.DB.Exec(query, id)
	if err != nil {
		log.Println("error query", err)
		return errhandler.ErrQuery
	}

	affectedRow, err := result.RowsAffected()
	if err != nil {
		log.Println("error getting num of rows affected")
		return errhandler.ErrRowsAffected
	}

	if affectedRow == 0 {
		log.Println("")
		return errhandler.ErrNoUpdate
	}

	return nil
}
