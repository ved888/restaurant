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
								  email_id
                                     )
                     VALUES($1,$2,$3,$4,$5)
                                  returning id `
	var userID uuid.UUID

	err := DB.QueryRowx(sql, user.FirstName, user.MiddleName, user.LastName, user.Phone, user.EmailId).Scan(&userID)
	return &userID, err

}
func UpdateUser(users *model.Users, userId string) error {

	//language=sql
	updateUser := `update users
                set
                    first_name = $1,
                    middle_name = $2,
                    last_name = $3,  
                    phone=$4,
                    email_id=$5 ,
                    updated_at=now()
            where id=$6 ::uuid`
	arg := []interface{}{users.FirstName, users.MiddleName, users.LastName, users.Phone, users.EmailId, &userId}

	_, err := common.DB.Exec(updateUser, arg...)
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

func DeleteUser(userID *string) error {

	// language=sql
	sql := `update
	      users
	      set
	       deleted_at=now()
	       where id=$1::uuid`

	_, err := common.DB.Exec(sql, userID)
	return err

}

func GetUserByID(userID *string) (*model.Users, error) {
	var user model.Users

	// language=SQL
	SQL := `SELECT
           id,
           first_name,
           middle_name,
           last_name,
           phone,
           email_id
            FROM 
                users
            WHERE id = $1::uuid
            `
	err := common.DB.Get(&user, SQL, userID)
	if err != nil && err == sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, nil

}
