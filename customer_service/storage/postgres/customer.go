package postgres

import (
	"database/sql"
	pb "exam/customer_service/genproto/customer"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type customerRepo struct {
	db *sqlx.DB
}

func NewCustomerRepo(db *sqlx.DB) *customerRepo {
	return &customerRepo{db: db}
}

func (r *customerRepo) Create(user *pb.CustomerRequest) (*pb.CustomerResponse, error) {
	tr, _ := r.db.Begin()
	defer tr.Rollback()
	respCustom := pb.CustomerResponse{}
	err := tr.QueryRow(`INSERT INTO customers(
		id,
		full_name,  
		bio, email,
		password, 
		refresh_token) 
		values($1, $2, $3, $4, $5, $6) 
		RETURNING 
		id, 
		full_name, 
		bio, email,
		password, 
		refresh_token,
		created_at, updated_at`,
		user.Id,
		user.FullName,
		user.Bio, user.Email,
		user.Password,
		user.RefreshToken,
	).Scan(
		&respCustom.Id,
		&respCustom.FullName,
		&respCustom.Bio,
		&respCustom.Email,
		&respCustom.Password,
		&respCustom.RefreshToken,
		&respCustom.CreatedAt,
		&respCustom.UpdatedAt,
	)
	if err != nil {
		tr.Rollback()
		fmt.Println("error while inserting customers")
		return &pb.CustomerResponse{}, err
	}
	var addresses []*pb.Address
	for _, address := range user.Addresses {
		addresResp := pb.Address{}
		err := tr.QueryRow(`INSERT INTO addresses(
			id,
			owner_id, 
			country, 
			street) 
			values($1, $2, $3, $4) 
			RETURNING 
			id, owner_id, 
			country, street`,
			address.Id,
			respCustom.Id,
			address.Country,
			address.Street).Scan(
			&addresResp.Id,
			&addresResp.OwnerId,
			&addresResp.Country,
			&addresResp.Street)
		if err != nil {
			tr.Rollback()
			fmt.Println("error while inserting addresses")
			return &pb.CustomerResponse{}, err
		}
		addresses = append(addresses, &addresResp)
	}
	respCustom.Addresses = addresses
	if err = tr.Commit(); err != nil {
		fmt.Println("error tr.commit", err)
	}
	return &respCustom, nil
}

func (r *customerRepo) GetCustomerInfo(req *pb.CustomerID) (*pb.CustomerInfo, error) {
	customer := pb.CustomerInfo{}
	err := r.db.QueryRow(`SELECT 
		id, 
		full_name, 
		bio, 
		email,
		password,
		created_at, 
		updated_at FROM customers WHERE id=$1 AND deleted_at IS NULL`, req.Id).Scan(
		&customer.Id,
		&customer.FullName,
		&customer.Bio,
		&customer.Email,
		&customer.Password,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return &pb.CustomerInfo{}, err
	}
	if err != nil {
		fmt.Println("error while selecting customers")
		return &pb.CustomerInfo{}, err
	}
	rows, err := r.db.Query(`SELECT id, owner_id, country, street FROM addresses WHERE owner_id=$1`, req.Id)
	if err == sql.ErrNoRows {
		return &pb.CustomerInfo{}, err
	}
	if err != nil {
		fmt.Println("error while selecting address")
		return &pb.CustomerInfo{}, err
	}
	defer rows.Close()

	for rows.Next() {
		address := pb.Address{}
		err = rows.Scan(
			&address.Id,
			&address.OwnerId,
			&address.Country,
			&address.Street,
		)
		if err != nil {
			fmt.Println("error while scanning address ")
			return &pb.CustomerInfo{}, err
		}
		customer.Addresses = append(customer.Addresses, &address)
	}
	return &customer, nil
}

func (r *customerRepo) UpdateCustomer(req *pb.CustomerUpdate) (*pb.CustomerResponse, error) {
	user := pb.CustomerResponse{}
	_, err := r.db.Exec(`UPDATE customers SET updated_at=NOW(),
	full_name=$1,
	bio=$2, 
	email=$3,
	password=$4 WHERE id=$5 AND deleted_at IS NULL`, req.FullName, req.Bio, req.Email, req.Password, req.Id)
	if err != nil {
		fmt.Println("error while updating customers")
		return &pb.CustomerResponse{}, err
	}

	err = r.db.QueryRow(`SELECT 
	id, 
	full_name, 
	bio,
	email, 
	password,
	created_at, 
	updated_at 
	FROM customers WHERE id=$1`, req.Id).Scan(
		&user.Id,
		&user.FullName,
		&user.Bio,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt)

	if err != nil {
		fmt.Println("error while getting customers update")
		return &pb.CustomerResponse{}, err
	}
	rowsAddress, err := r.db.Query(`SELECT id, owner_id, country, street FROM addresses WHERE owner_id=$1`, req.Id)
	if err != nil {
		fmt.Println("error while getting customers addresses update")
		return &pb.CustomerResponse{}, err
	}
	defer rowsAddress.Close()

	for rowsAddress.Next() {
		addressRes := pb.Address{}
		err = rowsAddress.Scan(&addressRes.Id, &addressRes.OwnerId, &addressRes.Country, &addressRes.Street)
		if err != nil {
			fmt.Println("error while scanning customer addresses update")
			return &pb.CustomerResponse{}, err
		}
		user.Addresses = append(user.Addresses, &addressRes)
	}

	return &user, nil
}

func (r *customerRepo) DeleteCustomer(req *pb.CustomerID) (*pb.Empty, error) {
	usersRepo := pb.Empty{}
	err := r.db.QueryRow(`update customers set deleted_at=NOW() where id=$1 and deleted_at is null`, req.Id).Err()
	if err != nil {
		fmt.Println("error while deleting customers")
		return &pb.Empty{}, err
	}
	return &usersRepo, nil
}

func (r *customerRepo) CheckFiedld(req *pb.CheckFieldReq) (*pb.CheckFieldResp, error) {
	query := fmt.Sprintf("SELECT 1 FROM customers WHERE %s=$1", req.Field)
	res := &pb.CheckFieldResp{}
	temp := -1
	err := r.db.QueryRow(query, req.Value).Scan(&temp)
	if err != nil {
		res.Exists = false
		return res, nil
	}
	if temp == 0 {
		res.Exists = true
	} else {
		res.Exists = false
	}
	return res, nil
}

func (r *customerRepo) GetByEmail(req *pb.EmailReq) (*pb.LoginResponse, error) {
	res := pb.LoginResponse{}
	err := r.db.QueryRow(`SELECT 
		id, 
		full_name, 
		bio, email, 
		password,
		refresh_token,
		created_at, 
		updated_at 
		FROM customers 
		WHERE email=$1 AND deleted_at IS NULL`, req.Email).Scan(
		&res.Id, &res.FullName, &res.Bio, &res.Email, &res.Password,
		&res.RefreshToken, &res.CreatedAt, &res.UpdatedAt,
	)
	if err != nil {
		fmt.Println("error while getting user login")
		return &pb.LoginResponse{}, err
	}
	rows, err := r.db.Query(`SELECT id, owner_id, country, street FROM addresses WHERE owner_id=$1`, res.Id)
	if err != nil {
		fmt.Println("error while getting addresses login")
	}
	for rows.Next() {
		address := pb.Address{}
		err = rows.Scan(&address.Id, &address.OwnerId, &address.Country, &address.Street)
		if err != nil {
			fmt.Println("error while scanning address")
			return &pb.LoginResponse{}, err
		}
		res.Addresses = append(res.Addresses, &address)
	}
	return &res, nil
}

func (r *customerRepo) GetCustomerBySearchOrder(user *pb.GetListUserRequest) (*pb.CustomerAll, error) {
	offset := (user.Page - 1) * user.Limit
	search := fmt.Sprintf("%s ILIKE $1", user.Search.Field)
	order := fmt.Sprintf("ORDER BY %s %s", user.Orders.Field, user.Orders.Value)

	query := `SELECT 
	full_name, 
	bio, 
	email
	FROM customers WHERE deleted_at IS NULL AND ` + search + " " + order + " "
	rowCustomer, err := r.db.Query(query+"LIMIT $2 OFFSET $3", "%"+user.Search.Value+"%", user.Limit, offset)
	if err == sql.ErrNoRows {
		return &pb.CustomerAll{}, err
	}
	if err != nil {
		fmt.Println("error while getting customers in search")
		return &pb.CustomerAll{}, err
	}
	customers := pb.CustomerAll{}
	for rowCustomer.Next() {
		customer := pb.CustomerListRes{}
		err = rowCustomer.Scan(
			&customer.FullName,
			&customer.Bio,
			&customer.Email)
		if err != nil {
			fmt.Println("error while scanning customer for search")
			return &pb.CustomerAll{}, err
		}
		customers.Customers = append(customers.Customers, &customer)
	}
	return &customers, nil
}
