package dbHelper

import (
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"restaurant/common"
	"restaurant/model"
)

func CreateBooking(db *sqlx.Tx, booking *model.Booking) (*string, error) {
	// language=SQL

	sql := `INSERT INTO booking(
                               booking_date,
                               pre_advance_booking
                               )
                        VALUES ($1,$2)
                               returning id`

	var bookingId string

	err := common.DB.QueryRowx(sql, booking.BookingDate, booking.PreAdvanceBooking).Scan(&bookingId)
	return &bookingId, err
}

func CreateUserBooking(db *sqlx.Tx, userID, bookingID string) (*uuid.UUID, error) {
	// language sql
	sql := `INSERT INTO user_booking(
                         users_id,
                         booking_id
                         )
                 values($1,$2)  
                 returning id`

	var userBookingId uuid.UUID
	err := db.QueryRowx(sql, userID, bookingID).Scan(&userBookingId)
	return &userBookingId, err
}

func CreateTableBooking(db *sqlx.Tx, userID, tableID string) (*uuid.UUID, error) {
	// language sql
	sql := `INSERT INTO booking_table(
                         rest_table_id,
                         booking_id
                         )
                 values($1,$2)  
                 returning id`

	var userBookingId uuid.UUID
	err := db.QueryRowx(sql, userID, tableID).Scan(&userBookingId)
	return &userBookingId, err
}

func GetBookingByBookingId(bookingId *string) (*model.Booking, error) {
	var booking model.Booking

	// language sql
	sql := `SELECT 
                id,
                booking_date,
                pre_advance_booking
           FROM
                booking where id=$1::uuid`

	err := common.DB.Get(&booking, sql, bookingId)
	if err != nil {
		return nil, err
	}
	return &booking, err
}

func GetBookingByUserId(userId *string) (*model.Booking, error) {
	var booking model.Booking

	// language sql
	sql := `SELECT 
                b.id,
                b.booking_date,
                b.pre_advance_booking,
                ub.users_id
           FROM
                booking b join user_booking ub on  
                    b.id=ub.booking_id
           where ub.users_id=$1::uuid`

	err := common.DB.Get(&booking, sql, userId)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func GetAllBooking() ([]*model.Booking, error) {
	booking := make([]*model.Booking, 0)
	// language=SQL

	SQL := `SELECT 
                id,
                booking_date,
                pre_advance_booking
           FROM
                booking`

	err := common.DB.Select(&booking, SQL)
	return booking, err
}

func UpdateBooking(booking *model.Booking, bookingId, userId string) error {

	// language=SQL
	SQL := `UPDATE booking
             SET 
                 booking_date=$1,
                 pre_advance_booking=$2,
                 updated_at=now()
          where id=$3 ::uuid`

	_, err := common.DB.Exec(SQL, booking.BookingDate, booking.PreAdvanceBooking, bookingId)
	return err
}

func DeleteBooking(bookingId *string) error {

	// language sql
	SQL := `update 
            booking 
            set
            deleted_at =now()
             where id=$1::uuid`
	_, err := common.DB.Exec(SQL, bookingId)
	return err
}
