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

func CreateAddress(w http.ResponseWriter, r *http.Request) {
	var address model.Address

	// read the userId from the path param
	userId := mux.Vars(r)["userId"]
	if userId == "" {
		logrus.Error("CreateAddress: userId is empty") // checking userId empty or not
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "CreateAddress: userId is empty", nil)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		logrus.Error("CreateAddress: failed to decode address for create", err)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "CreateAddress: failed to decode address for create", nil)
		return
	}

	// validate the address field
	validate := validator.New()

	err = validate.Struct(address)
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
		// create the address entry
		addressID, err := dbHelper.CreateAddress(tx, &address)
		if err != nil {
			return errors.Wrap(err, "CreateAddress: failed to create address entry")
		}

		// create the user address relation entry
		_, err = dbHelper.CreateUserAddress(tx, userId, *addressID)
		if err != nil {
			return errors.Wrap(err, "CreateAddress: failed to create user address relation entry")
		}
		return nil
	})

	if txErr != nil {
		logrus.Error("CreateAddress: failed to create address for the user", txErr)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "CreateAddress: failed to create address for the user", nil)
		return
	}
	common.ReturnResponse(w, "success", http.StatusCreated, "", address)
}

func GetAddressByAddressId(w http.ResponseWriter, r *http.Request) {
	// read user id from path param
	userId := mux.Vars(r)["userId"]
	if userId == "" {
		logrus.Error("GetAddressByAddressId: userId is empty") // checking userId empty or not
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "GetAddressByAddressId: userId is empty", nil)
		return
	}

	// read address id from path param
	addressId := mux.Vars(r)["addressId"]
	if addressId == "" {
		logrus.Error("GetAddressByAddressId: addressId is empty") // checking addressId empty or not
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "GetAddressByAddressId: addressId is empty", nil)
		return
	}

	// Get address by the userId and addressId
	result, err := dbHelper.GetAddressByAddressId(&addressId)
	if err != nil {
		logrus.Error("GetAddressByAddressId: failed to get address", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetAddressByAddressId: failed to get address", nil)
		return
	}

	common.ReturnResponse(w, "success", http.StatusOK, "", result)

}

func GetAddressByUserId(w http.ResponseWriter, r *http.Request) {
	// read user id from path param
	userId := mux.Vars(r)["userId"]
	if userId == "" {
		logrus.Error("GetAddressByUserId: userId is empty") // checking userId empty or not
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "GetAddressByUserId: userId is empty", nil)
		return
	}

	// Get address by the userId and addressId
	result, err := dbHelper.GetAddressByUserId(&userId)
	if err != nil {
		logrus.Error("GetAddressByUserId: failed to get address", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetAddressByUserId: failed to get address", nil)
		return
	}

	common.ReturnResponse(w, "success", http.StatusOK, "", result)

}

func UpdateAddress(w http.ResponseWriter, r *http.Request) {
	var address model.Address

	err := json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		logrus.Error("UpdateAddress: failed to decoding billing update", err)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "UpdateAddress: failed to decoding billing update", nil)
		return
	}

	// read the userId from path param
	userId := mux.Vars(r)["userId"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(userId); uuidErr != nil {
		logrus.Error("UpdateAddress: Failed to parse the user id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "UpdateAddress: Failed to parse the user id to uuid", nil)
		return
	}
	// read the address id from the path params
	addressId := mux.Vars(r)["addressId"]

	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(addressId); uuidErr != nil {
		logrus.Error("UpdateAddress: Failed to parse the address id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "UpdateAddress: Failed to parse the address id to uuid", nil)
		return
	}

	// update the address by use of the userId and addressId
	err = dbHelper.UpdateAddress(&address, addressId, userId)
	if err != nil {
		logrus.Error("UpdateAddress: failed to update address query", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "UpdateAddress: failed to update address query", nil)
		return
	}
	common.ReturnResponse(w, "success", http.StatusNoContent, "", address)
	//w.WriteHeader(http.StatusNoContent)
}

func DeleteAddressByAddressId(w http.ResponseWriter, r *http.Request) {

	// read the address id from path param
	addressId := mux.Vars(r)["addressId"]
	if addressId == "" {
		logrus.Error("DeleteAddressByAddressId: addressId is empty") // checking addressId empty or not
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "DeleteAddressByAddressId: addressId is empty", nil)
		return
	}
	// delete the allAddress for particular user
	err := dbHelper.DeleteAddressById(&addressId)
	if err != nil {
		logrus.Error("DeleteAddressByAddressId: failed to delete all address for userId", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "DeleteAddressByAddressId: failed to delete all address for userId", nil)
		return
	}
	common.ReturnResponse(w, "success", http.StatusNoContent, "", nil)
	//w.WriteHeader(http.StatusNoContent)
}

func DeleteAddressByUserId(w http.ResponseWriter, r *http.Request) {
	// read the user id from the path params
	userId := mux.Vars(r)["userId"]

	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(userId); uuidErr != nil {
		logrus.Error("DeleteAddressByUserId: Failed to parse the user id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "DeleteAddressByUserId: Failed to parse the user id to uuid", nil)
		return
	}
	// delete the address by use of the id
	err := dbHelper.DeleteAddressByUserId(&userId)
	if err != nil {
		logrus.Error("DeleteAddressByUserId: failed to delete address query", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "DeleteAddressByUserId: failed to delete address query", nil)
		return
	}
	common.ReturnResponse(w, "success", http.StatusNoContent, "", nil)
	//w.WriteHeader(http.StatusNoContent)
}

func GetAllAddress(w http.ResponseWriter, r *http.Request) {

	// get all the address list
	addressList, err := dbHelper.GetAllAddress()
	if err != nil {
		logrus.Error("GetAllAddress: failed to getAll ", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetAllAddress: failed to getAll", nil)
		return
	}

	common.ReturnResponse(w, "success", http.StatusOK, "", addressList)
}
