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

func CreateBilling(w http.ResponseWriter, r *http.Request) {
	var billing model.Billing

	err := json.NewDecoder(r.Body).Decode(&billing)
	if err != nil {
		logrus.Error("error in decode billing create", err)
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

	// read the orderId from query param
	orderId := r.URL.Query().Get("orderId")
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(orderId); uuidErr != nil {
		logrus.Error("Failed to parse the order id to uuid", uuidErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//validate all billing fields
	validate := validator.New()
	err = validate.Struct(billing)
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
		// create the billing entry
		billingID, err := dbHelper.CreateBilling(tx, &billing)
		if err != nil {
			return err
		}

		// create the user billing relation entry
		_, err = dbHelper.CreateUserBilling(tx, userId, *billingID)
		if err != nil {
			return err
		}
		// create the order billing relation entry
		_, err = dbHelper.CreateOrderBilling(tx, orderId, *billingID)
		if err != nil {
			return err
		}
		return nil
	})

	if txErr != nil {
		logrus.Error("failed to create billing for the user", txErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetBillingById(w http.ResponseWriter, r *http.Request) {
	// read user id from path param
	userId := mux.Vars(r)["userId"]
	if userId == "" {
		logrus.Error("userId is empty") //checking userId empty or not
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// read the billing id from path param
	billingId := mux.Vars(r)["billingId"]
	if billingId == "" { // checking billingId empty or not
		logrus.Error("billingId is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// get the billing by id
	resp, err := dbHelper.GetBillingById(&billingId)
	if err != nil {
		logrus.Error("failed to get billing by id", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		logrus.Error("failed to encode the billing ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetBillingByUserId(w http.ResponseWriter, r *http.Request) {
	// read user id from path param
	userId := mux.Vars(r)["userId"]
	if userId == "" {
		logrus.Error("userId is empty") //checking userId empty or not
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get the billing by id
	resp, err := dbHelper.GetBillingByUserId(&userId)
	if err != nil {
		logrus.Error("failed to get billing by id", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		logrus.Error("failed to encode the billing ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetBillingByOrderId(w http.ResponseWriter, r *http.Request) {
	// read user id from path param
	orderId := r.URL.Query().Get("orderId")
	if orderId == "" {
		logrus.Error("orderId is empty") //checking userId empty or not
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// read the billing id from path param
	billingId := mux.Vars(r)["billingId"]
	if billingId == "" { // checking billingId empty or not
		logrus.Error("billingId is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// get the billing by id
	resp, err := dbHelper.GetBillingByOrderId(&orderId)
	if err != nil {
		logrus.Error("failed to get billing by id", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		logrus.Error("failed to encode the billing ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetAllBilling(w http.ResponseWriter, r *http.Request) {

	// get all the billing list
	billingList, err := dbHelper.GetAllBilling()
	if err != nil {
		logrus.Error("error in billing getAll query", err)
		return
	}
	err = json.NewEncoder(w).Encode(billingList)
	if err != nil {
		logrus.Error("error in billing getAll encoding ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func UpdateBilling(w http.ResponseWriter, r *http.Request) {
	var billing model.Billing

	err := json.NewDecoder(r.Body).Decode(&billing)
	if err != nil {
		logrus.Error("error in decoding billing update", err)
		return
	}
	// read user id from path param
	userId := mux.Vars(r)["userId"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(userId); uuidErr != nil {
		logrus.Error("Failed to parse the user id to uuid", uuidErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// read the billing id from the path params
	billingId := mux.Vars(r)["billingId"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(billingId); uuidErr != nil {
		logrus.Error("Failed to parse the billing id to uuid", uuidErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// update the billing entry
	err = dbHelper.UpdateBilling(&billing, billingId, userId)
	if err != nil {
		logrus.Error("error in updating billing ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteBillingById(w http.ResponseWriter, r *http.Request) {
	// read the billing id from the path params
	billingId := mux.Vars(r)["billingId"]

	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(billingId); uuidErr != nil {
		logrus.Error("Failed to parse the billing id to uuid", uuidErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// read the user id from the path params
	userId := mux.Vars(r)["userId"]

	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(userId); uuidErr != nil {
		logrus.Error("Failed to parse the user id to uuid", uuidErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// delete the billing entry by id
	err := dbHelper.DeleteBillingById(&billingId, &userId)
	if err != nil {
		logrus.Error("error in delete billing query", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteBilling(w http.ResponseWriter, r *http.Request) {
	// read the billing id from the path params
	userId := mux.Vars(r)["userId"]

	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(userId); uuidErr != nil {
		logrus.Error("Failed to parse the user id to uuid", uuidErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// delete the billing entry by id
	err := dbHelper.DeleteBilling(&userId)
	if err != nil {
		logrus.Error("error in delete billing query", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
