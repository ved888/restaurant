package dbHelper

import (
	_ "database/sql"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"restaurant/common"
	"restaurant/model"
)

func CreateAddress(db *sqlx.Tx, address *model.Address) (*string, error) {
	// language=SQL
	SQL := `INSERT INTO address(
                    line1,
                    line2,
                    pin_code,
                    city,
                   state,
                    country)
            VALUES ($1,$2,$3,$4,$5,$6)
                  returning id`

	var addressID string
	args := []interface{}{address.Line1, address.Line2, address.PinCode, address.City, address.State, address.Country}
	err := db.QueryRowx(SQL, args...).Scan(&addressID)
	return &addressID, err

}

func CreateUserAddress(db *sqlx.Tx, userID, AddressID string) (*uuid.UUID, error) {
	// language sql
	sql := `INSERT INTO user_address(
                         users_id,
                         address_id
                         )
                 values($1,$2)  
                 returning id`

	var userAddressId uuid.UUID
	err := db.QueryRowx(sql, userID, AddressID).Scan(&userAddressId)
	return &userAddressId, err
}

func GetAllAddress() ([]*model.Address, error) {

	address := make([]*model.Address, 0)
	// language=SQL

	SQL := `SELECT 
                 id,
                 line1, 
                 line2,
                 pin_code,
                 city,
                 state,
                 country
            from
                 address`

	err := common.DB.Select(&address, SQL)
	return address, err
}

func GetAddressByAddressId(addressId *string) (*model.Address, error) {
	var address model.Address

	// language sql

	sql := `select  
                 id,
                 line1, 
                 line2,
                 pin_code,
                 city,
                 state,
                 country
       from address
       where id=$1::uuid`

	err := common.DB.Get(&address, sql, addressId)
	if err != nil {
		return nil, err
	}
	return &address, nil

}

func GetAddressByUserId(userId *string) (*model.Address, error) {
	var address model.Address

	// language sql

	sql := `select  
                 a.id,
                 a.line1, 
                 a.line2,
                 a.pin_code,
                 a.city,
                 a.state,
                 a.country,
                 ua.users_id
       from address a join user_address ua on 
           a.id=ua.address_id
       where ua.users_id=$1::uuid`

	err := common.DB.Get(&address, sql, userId)
	if err != nil {
		return nil, err
	}
	return &address, nil

}

func UpdateAddress(address *model.Address, addressId, userId string) error {

	// language=SQL
	SQL := `update address
                  set
                      line1=$1,
                      line2=$2,
                      pin_code=$3,
                      city=$4,
                      state=$5,
                      country=$6,
                      updated_at=now()
               where id = $7`

	args := []interface{}{address.Line1, address.Line2, address.PinCode, address.City, address.State, address.Country, addressId}

	_, err := common.DB.Exec(SQL, args...)
	return err
}
func DeleteAddressById(addressId *string) error {

	// language sql
	sql := `update
	              address
	              set
        	      deleted_at =now()
	              where id=$1::uuid`

	_, err := common.DB.Exec(sql, addressId)
	if err != nil {
		return err
	}
	return err
}

func DeleteAddressByUserId(userId *string) error {

	// language sql
	SQL := `update
            address
            set
            deleted_at =now()
            where id in (select users_id from user_address where users_id=$1) `
	_, err := common.DB.Exec(SQL, userId)
	return err

}
