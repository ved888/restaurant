package operation

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"net/http"
	"restaurant/common"
	"restaurant/dbHelper"
	"restaurant/model"
)

func CreateBooking(w http.ResponseWriter, r *http.Request) {
	var booking model.Booking
	// read userId from path param
	userId := mux.Vars(r)["userId"]
	if userId == "" {
		logrus.Error("userId is empty") //checking userId empty or not
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// read the tableId from query param
	tableId := r.URL.Query().Get("tableId")
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(tableId); uuidErr != nil {
		logrus.Error("Failed to parse the table id to uuid", uuidErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&booking)
	if err != nil {
		logrus.Error("error in decode create booking")
		return
	}

	// validate the booking field
	validate := validator.New()
	err = validate.Struct(booking)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		w.WriteHeader(http.StatusBadRequest)
		responseBody := map[string]string{"error": validationErrors.Error()}
		logrus.Error(responseBody)
		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		}
		return
	}
	txErr := common.Tx(func(tx *sqlx.Tx) error {
		// create the booking entry
		bookingID, err := dbHelper.CreateBooking(tx, &booking)
		if err != nil {
			return err
		}

		// create the user booking relation entry
		_, err = dbHelper.CreateUserBooking(tx, userId, *bookingID)
		if err != nil {
			return err
		}

		// create the table booking relation entry
		_, err = dbHelper.CreateTableBooking(tx, tableId, *bookingID)
		if err != nil {
			return err
		}
		return nil
	})

	if txErr != nil {
		logrus.Error("failed to create booking for the user", txErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func GetBookingByBookingId(w http.ResponseWriter, r *http.Request) {

	// read the booking id from path param
	bookingId := mux.Vars(r)["bookingId"]
	if bookingId == "" { // checking bookingId empty or not
		logrus.Error("booking id is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// get the booking by id
	resp, err := dbHelper.GetBookingByBookingId(&bookingId)
	if err != nil {
		logrus.Error("failed to get booking by id", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		logrus.Error("failed to encode the booking ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetBookingByUserId(w http.ResponseWriter, r *http.Request) {

	// read userId from path param
	userId := mux.Vars(r)["userId"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(userId); uuidErr != nil {
		logrus.Error("Failed to parse the table id to uuid", uuidErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get the booking by id
	resp, err := dbHelper.GetBookingByUserId(&userId)
	if err != nil {
		logrus.Error("failed to get booking by id", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		logrus.Error("failed to encode the booking ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetAllBooking(w http.ResponseWriter, r *http.Request) {

	//get all the booking list
	bookingList, err := dbHelper.GetAllBooking()
	if err != nil {
		logrus.Error("error in getAll booking ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(bookingList)
	if err != nil {
		logrus.Error("error in encoding getAll booking", err)
		return
	}

}

func UpdateBooking(w http.ResponseWriter, r *http.Request) {
	var booking model.Booking
	err := json.NewDecoder(r.Body).Decode(&booking)
	if err != nil {
		logrus.Error("error in decode update booking")
		return
	}
	// read user id from path param
	userId := mux.Vars(r)["userId"]
	if userId == "" {
		logrus.Error("userId is empty") //checking userId empty or not
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// read the booking id from the path params
	bookingId := mux.Vars(r)["bookingId"]
	if userId == "" {
		logrus.Error("userId is empty") //checking bookingId empty or not
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//update the booking entry by user and booking id
	err = dbHelper.UpdateBooking(&booking, bookingId, userId)
	if err != nil {
		logrus.Error("error in booking update query", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteBooking(w http.ResponseWriter, r *http.Request) {

	// read the booking id from path param
	bookingId := mux.Vars(r)["bookingId"]
	if bookingId == "" {
		logrus.Error("bookingId is empty") //checking bookingId empty or not
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// delete the booking entry by id
	err := dbHelper.DeleteBooking(&bookingId)
	if err != nil {
		logrus.Error("error in deleted booking", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
