package operation

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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
		logrus.Error("CreateBooking: userId is empty") // checking userId empty or not
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "CreateBooking: userId is empty", nil)
		return
	}

	// read the tableId from query param
	tableId := r.URL.Query().Get("tableId")
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(tableId); uuidErr != nil {
		logrus.Error("CreateBooking: Failed to parse the table id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "CreateBooking: Failed to parse the table id to uuid", nil)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&booking)
	if err != nil {
		logrus.Error("CreateBooking: Failed to decode for create booking")
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "CreateBooking: Failed to decode for create booking", nil)
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
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}

	txErr := common.Tx(func(tx *sqlx.Tx) error {
		// create the booking entry
		bookingID, err := dbHelper.CreateBooking(tx, &booking)
		if err != nil {
			return errors.Wrap(err, "CreateBooking: failed to create the booking entry")
		}

		// create the user booking relation entry
		_, err = dbHelper.CreateUserBooking(tx, userId, *bookingID)
		if err != nil {
			return errors.Wrap(err, "CreateBooking: failed to create the user booking relation entry")
		}

		// create the table booking relation entry
		_, err = dbHelper.CreateTableBooking(tx, tableId, *bookingID)
		if err != nil {
			return errors.Wrap(err, "CreateBooking: failed to create the table booking relation entry")
		}
		return nil
	})

	if txErr != nil {
		logrus.Error("CreateBooking: failed to create booking for the user", txErr)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "CreateBooking: failed to create booking for the user", nil)
		return
	}
	common.ReturnResponse(w, "success", http.StatusCreated, "", booking)

}

func GetBookingByBookingId(w http.ResponseWriter, r *http.Request) {

	// read the booking id from path param
	bookingId := mux.Vars(r)["bookingId"]
	if bookingId == "" { // checking bookingId empty or not
		logrus.Error("GetBookingByBookingId: booking id is empty")
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "GetBookingByBookingId: booking id is empty", nil)
		return
	}

	// get the booking by id
	resp, err := dbHelper.GetBookingByBookingId(&bookingId)
	if err != nil {
		logrus.Error("GetBookingByBookingId: failed to get booking by id", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetBookingByBookingId: failed to get booking by id", nil)
		return
	}

	//err = json.NewEncoder(w).Encode(&resp)
	//if err != nil {
	//	logrus.Error("GetBookingByBookingId: failed to encode the booking ", err)
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	common.ReturnResponse(w, "success", http.StatusOK, "", resp)

}

func GetBookingByUserId(w http.ResponseWriter, r *http.Request) {

	// read userId from path param
	userId := mux.Vars(r)["userId"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(userId); uuidErr != nil {
		logrus.Error("GetBookingByUserId: Failed to parse the table id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "GetBookingByUserId: Failed to parse the table id to uuid", nil)
		return
	}

	// get the booking by id
	resp, err := dbHelper.GetBookingByUserId(&userId)
	if err != nil {
		logrus.Error("GetBookingByUserId: failed to get booking by id", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetBookingByUserId: failed to get booking by id", nil)
		return
	}

	//err = json.NewEncoder(w).Encode(&resp)
	//if err != nil {
	//	logrus.Error("GetBookingByUserId: failed to encode the booking ", err)
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	common.ReturnResponse(w, "success", http.StatusOK, "", resp)

}

func GetAllBooking(w http.ResponseWriter, r *http.Request) {

	// get all the booking list
	bookingList, err := dbHelper.GetAllBooking()
	if err != nil {
		logrus.Error("GetAllBooking: failed to getAll booking ", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetAllBooking: failed to getAll booking", nil)
		return
	}

	//err = json.NewEncoder(w).Encode(bookingList)
	//if err != nil {
	//	logrus.Error("GetAllBooking: failed to encoding getAll booking", err)
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	common.ReturnResponse(w, "success", http.StatusOK, "", bookingList)

}

func UpdateBooking(w http.ResponseWriter, r *http.Request) {
	var booking model.Booking
	err := json.NewDecoder(r.Body).Decode(&booking)
	if err != nil {
		logrus.Error("UpdatingBooking: failed to decode update booking")
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "UpdatingBooking: failed to decode update booking", nil)
		return
	}

	// read user id from path param
	userId := mux.Vars(r)["userId"]
	if userId == "" {
		logrus.Error("UpdatingBooking: userId is empty") // checking userId empty or not
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "UpdatingBooking: userId is empty", nil)
		return
	}

	// read the booking id from the path params
	bookingId := mux.Vars(r)["bookingId"]
	if userId == "" {
		logrus.Error("UpdatingBooking: booking is empty") // checking bookingId empty or not
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "UpdatingBooking: booking is empty", nil)
		return
	}

	// update the booking entry by user and booking id
	err = dbHelper.UpdateBooking(&booking, bookingId, userId)
	if err != nil {
		logrus.Error("UpdatingBooking: failed to booking update entry", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "UpdatingBooking: failed to booking update entry", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteBooking(w http.ResponseWriter, r *http.Request) {

	// read the booking id from path param
	bookingId := mux.Vars(r)["bookingId"]
	if bookingId == "" {
		logrus.Error("DeleteBooking: bookingId is empty") // checking bookingId empty or not
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "DeleteBooking: bookingId is empty", nil)
		return
	}

	// delete the booking entry by id
	err := dbHelper.DeleteBooking(&bookingId)
	if err != nil {
		logrus.Error("DeleteBooking: failed to deleted booking", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "DeleteBooking: failed to deleted booking", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
