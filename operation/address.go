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

func CreateAddress(w http.ResponseWriter, r *http.Request) {
	var address model.Address

	// read the userId from the path param
	userId := mux.Vars(r)["userId"]
	if userId == "" {
		logrus.Error("userId is empty") //checking userId empty or not
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		logrus.Error("error in decode create address", err)
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
		}
		return
	}

	txErr := common.Tx(func(tx *sqlx.Tx) error {
		// create the address entry
		addressID, err := dbHelper.CreateAddress(tx, &address)
		if err != nil {
			return err
		}

		// create the user address relation entry
		_, err = dbHelper.CreateUserAddress(tx, userId, *addressID)
		if err != nil {
			return err
		}
		return nil
	})

	if txErr != nil {
		logrus.Error("failed to create address for the user", txErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetAddressByAddressId(w http.ResponseWriter, r *http.Request) {
	// read user id from path param
	userId := mux.Vars(r)["userId"]
	if userId == "" {
		logrus.Error("userId is empty") //checking userId empty or not
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// read address id from path param
	addressId := mux.Vars(r)["addressId"]
	if addressId == "" {
		logrus.Error("addressId is empty") //checking addressId empty or not
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Get address by the userId and addressId
	result, err := dbHelper.GetAddressByAddressId(&addressId)
	if err != nil {
		logrus.Error("failed to get address", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&result)
	if err != nil {
		logrus.Error("failed to encode the address ", err)
		return
	}
}

func GetAddressByUserId(w http.ResponseWriter, r *http.Request) {
	// read user id from path param
	userId := mux.Vars(r)["userId"]
	if userId == "" {
		logrus.Error("userId is empty") //checking userId empty or not
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get address by the userId and addressId
	result, err := dbHelper.GetAddressByUserId(&userId)
	if err != nil {
		logrus.Error("failed to get address", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(&result)
	if err != nil {
		logrus.Error("failed to encode the address ", err)
		return
	}

}

func UpdateAddress(w http.ResponseWriter, r *http.Request) {
	var address model.Address

	err := json.NewDecoder(r.Body).Decode(&address)
	if err != nil {
		logrus.Error("error in decoding billing update", err)
		return
	}
	// read the userId from path param
	userId := mux.Vars(r)["userId"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(userId); uuidErr != nil {
		logrus.Error("Failed to parse the address id to uuid", uuidErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// read the address id from the path params
	addressId := mux.Vars(r)["addressId"]

	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(addressId); uuidErr != nil {
		logrus.Error("Failed to parse the address id to uuid", uuidErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// update the address by use of the userId and addressId
	err = dbHelper.UpdateAddress(&address, addressId, userId)
	if err != nil {
		logrus.Error("error in update address query", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteAddressByAddressId(w http.ResponseWriter, r *http.Request) {

	// read the address id from path param
	addressId := mux.Vars(r)["addressId"]
	if addressId == "" {
		logrus.Error("addressId is empty") //checking addressId empty or not
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// delete the allAddress for particular user
	err := dbHelper.DeleteAddressById(&addressId)
	if err != nil {
		logrus.Error("failed to delete all address for userId", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteAddressByUserId(w http.ResponseWriter, r *http.Request) {
	// read the user id from the path params
	userId := mux.Vars(r)["userId"]

	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(userId); uuidErr != nil {
		logrus.Error("Failed to parse the user id to uuid", uuidErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// delete the address by use of the id
	err := dbHelper.DeleteAddressByUserId(&userId)
	if err != nil {
		logrus.Error("error in delete address query", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func GetAllAddress(w http.ResponseWriter, r *http.Request) {

	// get all the address list
	addressList, err := dbHelper.GetAllAddress()
	if err != nil {
		logrus.Error("error in getAll ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(&addressList)
	if err != nil {
		logrus.Error("error in encode getAll address", err)
		return
	}

}
