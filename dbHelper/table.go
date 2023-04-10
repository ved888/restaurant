package dbHelper

import (
	"restaurant/common"
	"restaurant/model"
)

func CreateTable(table *model.ResTable) (*string, error) {
	// language=sql
	sql := `INSERT INTO res_table(
                            code,
                            capacity)
                     VALUES ($1,$2)
                            returning id`

	var tableId string
	err := common.DB.QueryRowx(sql, table.Code, table.Capacity).Scan(&tableId)
	return &tableId, err
}

func GetAllTable() ([]*model.ResTable, error) {
	table := make([]*model.ResTable, 0)

	// language=sql
	SQL := `SELECT
                  id,
                  code,
                  capacity
              FROM 
                  res_table`

	err := common.DB.Select(&table, SQL)
	return table, err

}

func GetTableById(tableId string) (*model.ResTable, error) {
	var table model.ResTable

	// language sql
	sql := `SELECT
                  id,
                  code,
                  capacity
              FROM 
                  res_table
              where id=$1::uuid`

	err := common.DB.Get(&table, sql, tableId)
	if err != nil {
		return nil, err
	}
	return &table, err
}

func GetTableByBookingId(bookingId string) (*model.ResTable, error) {
	var table model.ResTable

	// language sql
	sql := `SELECT
                  t.id,
                  t.code,
                  t.capacity,
                  bt.booking_id
              FROM 
                  res_table t join booking_table bt on  
                      t.id = bt.rest_table_id
              where bt.booking_id=$1::uuid`

	err := common.DB.Get(&table, sql, bookingId)
	if err != nil {
		return nil, err
	}
	return &table, err
}

func UpdateTable(table *model.ResTable, tableId *string) error {

	// language=sql
	sql := `UPDATE res_table
                   set
                       code=COALESCE($1,code),
                       capacity=COALESCE($2,capacity),
                       updated_at=now()
               where id=$3`

	_, err := common.DB.Exec(sql, table.Code, table.Capacity, tableId)
	return err
}

func DeleteTable(tableId *string) error {

	// language=sql

	sql := `update
                 res_table
                 set
                 deleted_at=now()
                 where id=$1`

	_, err := common.DB.Exec(sql, tableId)
	return err

}
