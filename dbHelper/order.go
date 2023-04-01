package dbHelper

import (
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"restaurant/common"
	"restaurant/model"
)

func CreateOrder(db *sqlx.Tx, order *model.Orders) (*string, error) {

	//language=SQL
	sql := `INSERT INTO orders(
                      item_discount,
                      tax,
                      shipping,
                      total)
          VALUES ($1,$2,$3,$4)
                 returning id`

	var orderId string

	err := common.DB.QueryRowx(sql, order.ItemDiscount, order.Tax, order.Shipping, order.Total).Scan(&orderId)
	return &orderId, err

}

func CreateUserOrder(db *sqlx.Tx, userId, orderId string) (*uuid.UUID, error) {

	// language sql
	sql := `insert into user_order(
                       users_id,
                       orders_id
                       )
             values ($1,$2)
             returning id`
	var userOrderId uuid.UUID
	err := db.QueryRowx(sql, userId, orderId).Scan(&userOrderId)
	return &userOrderId, err
}

func GetOrderById(orderId *string) (*model.Orders, error) {
	var order model.Orders

	// language sql
	sql := `SELECT 
                  item_discount,
                  tax,
                  shipping,
                  total
           from
               orders  where id=$1::uuid`

	err := common.DB.Get(&order, sql, orderId)
	if err != nil {
		return nil, err
	}
	return &order, err
}

func GetOrderByUserId(userId *string) (*model.Orders, error) {
	var order model.Orders

	// language sql
	sql := `SELECT 
                  o.item_discount,
                  o.tax,
                  o.shipping,
                  o.total,
                  uo.users_id
           from 
               orders o  join user_order uo on  
                   o.id=uo.orders_id
           where users_id=$1::uuid`

	err := common.DB.Get(&order, sql, userId)
	if err != nil {
		return nil, err
	}
	return &order, err
}

func GetAllOrder() ([]*model.Orders, error) {
	orders := make([]*model.Orders, 0)
	// language=SQL
	SQL := `SELECT 
                  users_id,
                  item_discount,
                  tax,
                  shipping,
                  total
           from
               orders`

	err := common.DB.Select(&orders, SQL)
	return orders, err

}

func UpdateOrder(order *model.Orders, orderId, userId string) error {

	// language=SQL
	sql := `update orders 
                 set
                     item_discount=$1,
                     tax=$2,
                     shipping=$3,
                     total=$4,
                     updated_at= now()
                where id=$5 ::uuid`

	_, err := common.DB.Exec(sql, order.ItemDiscount, order.Tax, order.Shipping, order.Total, orderId)
	return err

}

func DeleteOrder(orderId, userId *string) error {

	// language=SQL
	delOrder := `Update 
                  orders 
                  set 
                  deleted_at = now() 
                  WHERE id=$1`
	_, err := common.DB.Exec(delOrder, orderId)
	return err
}

//func CreateOrder(order *model.Orders) error {

//language=kSQL
//	sql := `INSERT INTO orders(
//          item_discount,
//            tax,
//              shipping,
//                total)
//      VALUES ($1,$2,$3,$4,$5)
//               returning id`

//	var orderId string

//	err := common.DB.QueryRowx(sql, order.ItemDiscount, order.Tax, order.Shipping, order.Total).Scan(&orderId)
//	return err
//
//}
