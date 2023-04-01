package operation

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"restaurant/dbHelper"
	"restaurant/model"
)

func CreateTable(w http.ResponseWriter, r *http.Request) {

	var table model.ResTable
	err := json.NewDecoder(r.Body).Decode(&table)
	if err != nil {
		logrus.Error("error in decode create table", err)
		return
	}
	// validate the table field
	validate := validator.New()
	err = validate.Struct(table)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		w.WriteHeader(http.StatusBadRequest)
		responseBody := map[string]string{"error": validationErrors.Error()}
		logrus.Error(responseBody)
		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		}
		return
	}
	// create the table entry
	_, err = dbHelper.CreateTable(&table)
	if err != nil {
		logrus.Error("error in create table ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func GetAllTable(w http.ResponseWriter, r *http.Request) {
	// get all the table list
	tableList, err := dbHelper.GetAllTable()
	if err != nil {
		logrus.Error("error in getAll table query", err)
		return
	}
	err = json.NewEncoder(w).Encode(&tableList)
	if err != nil {
		logrus.Error("error in encode getAll table  ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetTableById(w http.ResponseWriter, r *http.Request) {
	// read the table id from path param
	tableId := mux.Vars(r)["tableId"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(tableId); uuidErr != nil {
		logrus.Error("Failed to parse the user id to uuid", uuidErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// get the table by id
	resp, err := dbHelper.GetTableById(tableId)
	if err != nil {
		logrus.Error("failed to Get table by table id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		logrus.Error("failed to encode the table ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetTableByBookingId(w http.ResponseWriter, r *http.Request) {
	// read the booking id from path param
	bookingId := mux.Vars(r)["bookingId"]
	if bookingId == "" { // checking bookingId empty or not
		logrus.Error("table id is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// get the table by booking id
	resp, err := dbHelper.GetTableByBookingId(bookingId)
	if err != nil {
		logrus.Error("failed to Get table by booking id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		logrus.Error("failed to encode the table ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateTable(w http.ResponseWriter, r *http.Request) {
	var table model.ResTable

	err := json.NewDecoder(r.Body).Decode(&table)
	if err != nil {
		logrus.Error("error in decode update table", err)
		return
	}
	// read the table id from the path params
	id := mux.Vars(r)["id"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(id); uuidErr != nil {
		logrus.Error("Failed to parse the user id to uuid", uuidErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//update the table entry by id
	err = dbHelper.UpdateTable(&table, &id)
	if err != nil {
		logrus.Error("error in  update table", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteTable(w http.ResponseWriter, r *http.Request) {
	// read the table id from the path params
	id := mux.Vars(r)["id"]

	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(id); uuidErr != nil {
		logrus.Error("Failed to parse the user id to uuid", uuidErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// delete the table entry by id from path params
	err := dbHelper.DeleteTable(&id)
	if err != nil {
		logrus.Error("error in delete table query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
