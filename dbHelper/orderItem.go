package dbHelper

import (
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"restaurant/common"
	"restaurant/model"
)

func CreateOrderItem(db *sqlx.Tx, orderItem *model.OrderItem) (*string, error) {
	// language=SQL
	sql := `INSERT INTO order_item(
                       price,
                       quantity )
            VALUES ($1,$2)
                       returning id`

	var orderItemId string

	err := common.DB.QueryRowx(sql, orderItem.Price, orderItem.Quantity).Scan(&orderItemId)
	return &orderItemId, err
}

func CreateOrderOrderItem(db *sqlx.Tx, orderItemId, orderId string) (*uuid.UUID, error) {
	// language=SQL
	sql := `INSERT INTO order_orderItem(
                       order_id,
                       orderItem_id 
                       )
            VALUES ($1,$2)
                       returning id`

	var orderItemOrderId uuid.UUID

	err := db.QueryRowx(sql, orderItemId, orderId).Scan(&orderItemOrderId)
	return &orderItemOrderId, err
}

func CreateFoodOrderItem(db *sqlx.Tx, foodId, orderItemId string) (*uuid.UUID, error) {

	//language sql

	sql := `insert into order_item_food(
                           food_id,
                           order_item_Id
                           )
                values ($1,$2)
                           returning id`

	var foodOrderItemId uuid.UUID
	err := db.QueryRowx(sql, foodId, orderItemId).Scan(&foodOrderItemId)
	return &foodOrderItemId, err

}

func GetOrderItemById(OrderItemId string) (*model.OrderItem, error) {
	var orderItem model.OrderItem

	// language sql
	sql := `SELECT 
                  price,
                  quantity
          from
                  order_item 
          where id=$1::uuid`

	err := common.DB.Get(&orderItem, sql, OrderItemId)
	if err != nil {
		return nil, err
	}
	return &orderItem, err
}

func GetOrderItemByOrderId(OrderItemId string) (*model.OrderItem, error) {
	var orderItem model.OrderItem

	// language sql
	sql := `SELECT 
                  oi.price,
                  oi.quantity,
                  oio.order_id
                  from
                  order_item oi join order_orderItem oio
                  on oi.id=oio.orderItem_id 
          where order_id=$1::uuid`

	err := common.DB.Get(&orderItem, sql, OrderItemId)
	if err != nil {
		return nil, err
	}
	return &orderItem, err
}

func GetAllOrderItem() ([]*model.OrderItem, error) {
	orderItem := make([]*model.OrderItem, 0)
	// language=sql
	SQL := `SELECT 
                  price,
                  quantity
          from
               order_item`

	err := common.DB.Select(&orderItem, SQL)
	return orderItem, err

}

func UpdateOrderItem(orderItem *model.OrderItem, Id string) error {

	// language=sql
	sql := `update order_item 
                 set
                     price=$1,
                     quantity=$2,
                     updated_at=now()
               where id=$3`

	_, err := common.DB.Exec(sql, orderItem.Price, orderItem.Quantity, Id)
	return err
}

func DeleteOrderItem(id *string) error {

	// language=SQL
	sql := `Update
                    order_item 
                    set
                    deleted_at = now()
                    WHERE id=$1`

	_, err := common.DB.Exec(sql, id)
	return err

}
