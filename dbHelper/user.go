package dbHelper

import (
	"database/sql"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"restaurant/common"
	"restaurant/model"
)

func CreateUser(DB *sqlx.Tx, user *model.Users) (*uuid.UUID, error) {

	// language=sql
	sql := `INSERT INTO users(
								  first_name,
								  middle_name,
	     						  last_name, 
								  phone,
                                  password,
								  email_id
                                     )
                     VALUES($1,$2,$3,$4,$5,$6)
                                  returning id `

	var userID uuid.UUID
	err := DB.QueryRowx(sql, user.FirstName, user.MiddleName, user.LastName, user.Phone, user.Password, user.EmailId).Scan(&userID)
	return &userID, err

}

func UpdateUser(DB *sqlx.Tx, users *model.Users, userID string) error {

	// language=SQL
	updateUser := `update users
                set
                    first_name = COALESCE($1,first_name),
                    middle_name = COALESCE($2,middle_name),
                    last_name = COALESCE($3,last_name),  
                    phone=COALESCE($4,phone),
                    email_id=COALESCE($5,email_id),
                    password=COALESCE($6,password),
                    updated_at=now()
            where id=$7 ::uuid`

	arg := []interface{}{users.FirstName, users.MiddleName, users.LastName, users.Phone, users.EmailId, users.Password, userID}

	_, err := DB.Exec(updateUser, arg...)
	return err
}

func GetAllUser() ([]*model.Users, error) {
	users := make([]*model.Users, 0)

	// language=sql
	SQL := `SELECT
                  id,
                  first_name,
                  middle_name,
                  last_name,
                  phone,
                  email_id
              FROM 
                  users`

	err := common.DB.Select(&users, SQL)
	return users, err

}

func GetUserByID(userID *string) (model.Users, error) {
	var user model.Users

	// language=SQL
	SQL := `SELECT
           id,
               first_name,
               middle_name,
               last_name,
               phone,
               email_id,
               password
            FROM 
                users
            WHERE id = $1::uuid`

	err := common.DB.Get(&user, SQL, userID)
	if err != nil && err != sql.ErrNoRows {
		return model.Users{}, err
	}

	if err == sql.ErrNoRows {
		return model.Users{}, nil
	}
	return user, nil

}

func GetUserByEmail(email string) (model.Users, error) {
	var user model.Users

	// language=SQL
	SQL := `SELECT
                id,
               first_name,
               middle_name,
               last_name,
               phone,
               email_id,
               password
            FROM 
                users
            WHERE email_id = $1`

	err := common.DB.Get(&user, SQL, email)
	if err != nil && err != sql.ErrNoRows {
		return model.Users{}, err
	}

	if err == sql.ErrNoRows {
		return model.Users{}, nil
	}
	return user, nil

}

func GetUserByPhone(phone string) (model.Users, error) {
	var user model.Users

	// language=SQL
	SQL := `SELECT
               id,
               first_name,
               middle_name,
               last_name,
               phone,
               email_id,
               password
            FROM 
                users
            WHERE phone = $1`

	err := common.DB.Get(&user, SQL, phone)
	if err != nil && err != sql.ErrNoRows {
		return model.Users{}, err
	}
	if err == sql.ErrNoRows {
		return model.Users{}, nil
	}
	return user, nil
}

func DeleteUser(DB *sqlx.Tx, userID *string) error {

	// language=sql
	sql := `update
	       users
	       set
	       deleted_at=now()
	       where id=$1::uuid`

	_, err := DB.Exec(sql, userID)
	return err

}
