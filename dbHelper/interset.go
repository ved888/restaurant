package dbHelper

import (
	"database/sql"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"restaurant/common"
	"restaurant/model"
)

func CreateInterest(db *sqlx.Tx, interest *model.Interest) (*uuid.UUID, error) {

	// language=sql
	sql := `INSERT INTO Interest(
				    name,
                    type )
                    VALUES($1,$2)
                    returning id `
	var interestId uuid.UUID

	err := db.QueryRowx(sql, interest.Name, interest.Type).Scan(&interestId)
	return &interestId, err
}

func UpdateInterest(interest *model.Interest, interestId string) error {

	// language=sql
	SQL := `update interest
                set
                    name = $1,
                    type = $2,
                    updated_at=now()
            where id in (select interest_id from relation_table where users_id = $3)`

	arg := []interface{}{interest.Name, interest.Type, interestId}

	_, err := common.DB.Exec(SQL, arg...)
	return err
}

func GetAllInterest() ([]*model.InterestUser, error) {
	interest := make([]*model.InterestUser, 0)

	// language=sql

	SQL := `SELECT
           i.id,
           i.name,
           i.type,
           ui.users_id as user_id
           FROM
                interest i
           inner join 
                relation_table ui 
           on
                ui.interest_id=i.id`

	err := common.DB.Select(&interest, SQL)
	return interest, err
}

func GetInterestByUserId(userId *string) (*model.Interest, error) {
	var interest model.Interest

	// language=sql

	SQL := `SELECT
           i.id,
           i.name,
           i.type
            FROM
                interest i
           inner join 
              relation_table ui 
            on
               ui.interest_id=i.id
           where ui.users_id =$1`

	err := common.DB.Get(&interest, SQL, userId)
	if err != nil && err == sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &interest, nil
}

func DeleteInterest(interestId *string) error {
	// language=sql
	sql := `update
	        interest
	        set
	        deleted_at=now()
	        where id in (select interest_id from relation_table where users_id=$1) `

	_, err := common.DB.Exec(sql, interestId)
	return err
}
