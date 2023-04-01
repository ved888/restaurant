package dbHelper

import (
	"restaurant/common"
	"restaurant/model"
)

func CreateFood(food *model.Food) error {
	// language=SQL
	sql := `INSERT INTO food(
                
                 name,
                 price,
                 type
                 )
    VALUES ($1,$2,$3)
    returning id`
	var foodId string

	err := common.DB.QueryRowx(sql, food.Name, food.Price, food.Type).Scan(&foodId)
	return err
}

func GetAllFood() ([]*model.Food, error) {
	food := make([]*model.Food, 0)
	// language=sql

	SQL := `SELECT
           id,
           name,
           price,
           type
    FROM 
        food`

	err := common.DB.Select(&food, SQL)
	return food, err
}

func GetFoodById(foodId string) (*model.Food, error) {
	var food model.Food

	// language sql
	sql := `SELECT
           id,
           name,
           price,
           type
    FROM 
        food where id=$1`

	err := common.DB.Get(&food, sql, foodId)
	if err != nil {
		return nil, err
	}
	return &food, err
}

func GetFoodByOrderItemId(orderItemId string) (*model.Food, error) {
	var food model.Food

	// language sql
	sql := `SELECT
           f.id,
           f.name,
           f.price,
           f.type,
           fo.order_item_id
    FROM 
        food f join order_item_food fo on 
        f.id=fo.food_id
    where fo.order_item_id=$1`

	err := common.DB.Get(&food, sql, orderItemId)
	if err != nil {
		return nil, err
	}
	return &food, err
}

func UpdateFood(food *model.Food, id string) error {

	// language=sql

	sql := `update food
               set
                   name=$1,
                   price=$2,
                   type=$3,
                   updated_at=now()
               where id=$4 ::uuid `

	arg := []interface{}{food.Name, food.Price, food.Type, id}

	_, err := common.DB.Exec(sql, arg...)
	return err

}

func DeleteFood(foodId *string) error {

	// language=SQL
	sql := `Update
                  food 
                  set 
                 deleted_at = now() 
            WHERE id=$1`

	_, err := common.DB.Exec(sql, foodId)
	return err

}
