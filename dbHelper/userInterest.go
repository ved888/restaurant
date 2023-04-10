package dbHelper

import (
	"database/sql"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"restaurant/common"
	"restaurant/model"
)

func CreateUserInterest(db *sqlx.Tx, userInterest *model.UserInterest) (*uuid.UUID, error) {
	// language sql

	sql := `INSERT INTO relation_table(
								 users_id,
                                interest_id
                                     )
                     VALUES($1,$2)
                    returning id `

	var userInterestId uuid.UUID
	err := db.QueryRowx(sql, userInterest.UsersId, userInterest.InterestId).Scan(&userInterestId)
	return &userInterestId, err

}

func UpdateUserInterest(db *sqlx.Tx, userInterest model.UserInterest, usersId string) error {
	// language sql

	sql := `update  relation_table
    set
								 users_id=$1,
                                interest_id=$2,
                               updated_at=now()      
            where id in (select id from relation_table where users_id = $3)`

	_, err := db.Exec(sql, userInterest.UsersId, userInterest.InterestId, &usersId)
	return err

}
func GetAllUserInterest() ([]*model.Interest, error) {
	var interest []*model.Interest
	// language=sql

	SQL := `SELECT
           id,
           users_id,
           interest_id
            FROM 
                interest`

	err := common.DB.Select(interest, SQL)
	return interest, err

}
func GetUserInterestById(userInterestId *string) (*model.UserInterest, error) {
	var userInterest model.UserInterest
	// language=sql
	SQL := `SELECT
           id,
           interest_id,
           users_id
            FROM 
                relation_table`

	err := common.DB.Get(&userInterest, SQL, userInterestId)
	if err != nil && err == sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &userInterest, nil
}
