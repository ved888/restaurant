package dbHelper

import (
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"restaurant/common"
	"restaurant/model"
)

func CreateBilling(db *sqlx.Tx, billing *model.Billing) (*string, error) {

	// language=sql

	SQL := `INSERT INTO billing(
                      type,
                      mode)
            VALUES ($1,$2)
                      returning id`

	var billingId string
	args := []interface{}{billing.Type, billing.Mode}

	err := common.DB.QueryRowx(SQL, args...).Scan(&billingId)
	return &billingId, err

}

func CreateUserBilling(db *sqlx.Tx, userID, BillingId string) (*uuid.UUID, error) {
	// language sql
	sql := `INSERT INTO user_billing(
                         users_id,
                         billing_id
                         )
                 values($1,$2)  
                 returning id`

	var userBillingId uuid.UUID
	err := db.QueryRowx(sql, userID, BillingId).Scan(&userBillingId)
	return &userBillingId, err
}

func CreateOrderBilling(db *sqlx.Tx, orderID, billingId string) (*uuid.UUID, error) {
	// language sql
	sql := `INSERT INTO Order_billing(
                         order_id,
                         billing_id
                         )
                 values($1,$2)  
                 returning id`

	var userBillingId uuid.UUID
	err := db.QueryRowx(sql, orderID, billingId).Scan(&userBillingId)
	return &userBillingId, err
}

// start Get by id query

func GetBillingById(billingId *string) (*model.Billing, error) {
	var billing model.Billing

	// language sql
	sql := `SELECT 
                  id,
                  type,
                  mode       
              from
                  billing where id=$1::uuid`

	err := common.DB.Get(&billing, sql, billingId)
	if err != nil {
		return nil, err
	}
	return &billing, err
}

func GetBillingByUserId(userId *string) (*model.Billing, error) {
	var billing model.Billing

	// language sql
	sql := `SELECT 
                  b.id,
                  b.type,
                  b.mode,
                  ub.users_id
                  from
                      billing b join user_billing ub on 
                          b.id = ub.billing_id
                  where ub.users_id=$1::uuid`

	err := common.DB.Get(&billing, sql, userId)
	if err != nil {
		return nil, err
	}
	return &billing, err
}

func GetBillingByOrderId(orderId *string) (*model.Billing, error) {
	var billing model.Billing

	// language sql
	sql := `SELECT 
                  b.id,
                  b.type,
                  b.mode,
                  ob.order_id
              from billing b join order_billing ob on 
                  b.id=ob.billing_id
                   where order_id=$1::uuid`

	err := common.DB.Get(&billing, sql, orderId)
	if err != nil {
		return nil, err
	}
	return &billing, err
}

func GetAllBilling() ([]*model.Billing, error) {
	billing := make([]*model.Billing, 0)
	// language=SQL

	SQL := `SELECT 
                  id,
                  type,
                  mode    
              from
                  billing`

	err := common.DB.Select(&billing, SQL)
	return billing, err

}

func UpdateBilling(billing *model.Billing, billingId, userId string) error {

	// language=SQL

	sql := `update billing 
                 set
                     type=$1,
                     mode=$2, 
                     updated_at=now()
                 where id=$3::uuid`

	_, err := common.DB.Exec(sql, billing.Type, billing.Mode, &billingId)

	return err

}
func DeleteBillingById(billingId, userId *string) error {

	// language SQL
	SQL := `Update
              billing
              set 
              deleted_at = now() 
             where id=$1::uuid`

	_, err := common.DB.Exec(SQL, billingId)
	return err

}

func DeleteBilling(userId *string) error {

	// language SQL
	SQL := `Update
              billing
              set 
              deleted_at = now() 
             where id in (select users_id from user_billing where users_id=$1)`

	_, err := common.DB.Exec(SQL, userId)
	return err

}

//func GetBillingById(billingId, userId *string) (*model.Billing, error) {
//	var billing model.Billing

// language sql
//	sql1 := `SELECT
//                  id,
//                  type,
//                  mode
//              from
//                  billing where id=$1::uuid`

//	err := common.DB.Get(&billing, sql, billingId)
//	if err != nil {
//		return nil, err
//	}
//	return &billing, err
//}
