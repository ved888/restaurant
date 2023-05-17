package operation

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"restaurant/common"
	"restaurant/dbHelper"
	"restaurant/model"
)

func CreateTable(w http.ResponseWriter, r *http.Request) {

	var table model.ResTable
	err := json.NewDecoder(r.Body).Decode(&table)
	if err != nil {
		logrus.Error("CreateTable: Failed to decode create table", err)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "CreateTable: Failed to decode create table", nil)
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
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// create the table entry
	_, err = dbHelper.CreateTable(&table)
	if err != nil {
		logrus.Error("CreateTable: Failed to create table entry ", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "CreateTable: Failed to create table entry", nil)
		return
	}
	common.ReturnResponse(w, "success", http.StatusCreated, "", table)
}

func GetAllTable(w http.ResponseWriter, r *http.Request) {
	// get all the table list
	tableList, err := dbHelper.GetAllTable()
	if err != nil {
		logrus.Error("GetAllTable: Failed to getAll table query", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetAllTable: Failed to getAll table query", nil)
		return
	}

	common.ReturnResponse(w, "success", http.StatusOK, "", tableList)
}

func GetTableById(w http.ResponseWriter, r *http.Request) {
	// read the table id from path param
	tableId := mux.Vars(r)["tableId"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(tableId); uuidErr != nil {
		logrus.Error("GetTableById: Failed to parse the user id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "GetTableById: Failed to parse the user id to uuid", nil)
		return
	}

	// get the table by id
	resp, err := dbHelper.GetTableById(tableId)
	if err != nil {
		logrus.Error("GetTableById: failed to Get table by table id")
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetTableById: failed to Get table by table id", nil)
		return
	}

	common.ReturnResponse(w, "success", http.StatusOK, "", resp)

}

func GetTableByBookingId(w http.ResponseWriter, r *http.Request) {
	// read the booking id from path param
	bookingId := mux.Vars(r)["bookingId"]
	if bookingId == "" { // checking bookingId empty or not
		logrus.Error("GetTableByBookingId: table id is empty")
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "GetTableByBookingId: table id is empty", nil)
		return
	}

	// get the table by booking id
	resp, err := dbHelper.GetTableByBookingId(bookingId)
	if err != nil {
		logrus.Error("GetTableByBookingId: failed to Get table by booking id")
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetTableByBookingId: failed to Get table by booking id", nil)
		return
	}

	common.ReturnResponse(w, "success", http.StatusOK, "", resp)

}

func UpdateTable(w http.ResponseWriter, r *http.Request) {
	var table model.ResTable

	err := json.NewDecoder(r.Body).Decode(&table)
	if err != nil {
		logrus.Error("UpdateTable: Failed to decode resTable", err)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "UpdateTable: Failed to decode resTable", nil)
		return
	}

	// read the table id from the path params
	id := mux.Vars(r)["id"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(id); uuidErr != nil {
		logrus.Error("UpdateTable: Failed to parse the user id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "UpdateTable: Failed to parse the user id to uuid", nil)
		return
	}

	// update the table entry by id
	err = dbHelper.UpdateTable(&table, &id)
	if err != nil {
		logrus.Error("UpdateTable: Failed to update table", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "UpdateTable: Failed to update table", nil)
		return
	}
	common.ReturnResponse(w, "success", http.StatusNoContent, "", table)
}

func DeleteTable(w http.ResponseWriter, r *http.Request) {
	// read the table id from the path params
	id := mux.Vars(r)["id"]

	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(id); uuidErr != nil {
		logrus.Error("DeleteTable: failed to parse the user id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "DeleteTable: failed to parse the user id to uuid", nil)
		return
	}

	// delete the table entry by id from path params
	err := dbHelper.DeleteTable(&id)
	if err != nil {
		logrus.Error("DeleteTable: Failed to delete table entry")
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "DeleteTable: Failed to delete table entry", nil)
		return
	}
	common.ReturnResponse(w, "success", http.StatusNoContent, "", nil)
}
