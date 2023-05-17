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

func CreateBilling(w http.ResponseWriter, r *http.Request) {
	var billing model.Billing

	err := json.NewDecoder(r.Body).Decode(&billing)
	if err != nil {
		logrus.Error("CreateBilling: failed to decode billing", err)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "CreateBilling: failed to decode billing", nil)
		return
	}

	// read the userId from path param
	userId := mux.Vars(r)["userId"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(userId); uuidErr != nil {
		logrus.Error("CreateBilling: Failed to parse the user id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "CreateBilling: Failed to parse the user id to uuid", nil)
		return
	}

	// read the orderId from query param
	orderId := r.URL.Query().Get("orderId")
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(orderId); uuidErr != nil {
		logrus.Error("CreateBilling: Failed to parse the order id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "CreateBilling: Failed to parse the order id to uuid", nil)
		return
	}

	// validate all billing fields
	validate := validator.New()
	err = validate.Struct(billing)
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
		// create the billing entry
		billingID, err := dbHelper.CreateBilling(tx, &billing)
		if err != nil {
			return errors.Wrap(err, "CreateBilling: failed to create the billing entry")
		}

		// create the user billing relation entry
		_, err = dbHelper.CreateUserBilling(tx, userId, *billingID)
		if err != nil {
			return errors.Wrap(err, "CreateBilling: failed to create the user billing relation entry")
		}
		// create the order billing relation entry
		_, err = dbHelper.CreateOrderBilling(tx, orderId, *billingID)
		if err != nil {
			return errors.Wrap(err, "CreateBilling: failed to create the order billing relation entry")
		}
		return nil
	})
	if txErr != nil {
		logrus.Error("CreateBilling: failed to create billing for the user", txErr)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "CreateBilling: failed to create billing for the user", nil)
		return
	}
	common.ReturnResponse(w, "success", http.StatusCreated, "", billing)
}

func GetBillingById(w http.ResponseWriter, r *http.Request) {
	// read user id from path param
	userId := mux.Vars(r)["userId"]
	if userId == "" {
		logrus.Error("GetBillingBYId: userId is empty") // checking userId empty or not
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "GetBillingBYId: userId is empty", nil)
		return
	}

	// read the billing id from path param
	billingId := mux.Vars(r)["billingId"]
	if billingId == "" { // checking billingId empty or not
		logrus.Error("GetBillingBYId: billingId is empty")
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "GetBillingBYId: billingId is empty", nil)
		return
	}

	// get the billing by id
	resp, err := dbHelper.GetBillingById(&billingId)
	if err != nil {
		logrus.Error("GetBillingBYId: failed to get billing by id", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetBillingBYId: failed to get billing by id", nil)
		return
	}

	common.ReturnResponse(w, "success", http.StatusOK, "", resp)

}

func GetBillingByUserId(w http.ResponseWriter, r *http.Request) {
	// read user id from path param
	userId := mux.Vars(r)["userId"]
	if userId == "" {
		logrus.Error("GetBillingBYUserId: userId is empty") // checking userId empty or not
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "GetBillingBYUserId: userId is empty", nil)
		return
	}

	// get the billing by id
	resp, err := dbHelper.GetBillingByUserId(&userId)
	if err != nil {
		logrus.Error("GetBillingBYUserId: failed to get billing by user id", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetBillingBYUserId: failed to get billing by user id", nil)
		return
	}
	common.ReturnResponse(w, "success", http.StatusOK, "", resp)
}

func GetBillingByOrderId(w http.ResponseWriter, r *http.Request) {
	// read user id from path param
	orderId := r.URL.Query().Get("orderId")
	if orderId == "" {
		logrus.Error("GetBillingByOrderId: orderId is empty") // checking userId empty or not
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "GetBillingByOrderId: orderId is empty", nil)
		return
	}

	// read the billing id from path param
	billingId := mux.Vars(r)["billingId"]
	if billingId == "" { // checking billingId empty or not
		logrus.Error("GetBillingByOrderId: billingId is empty")
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "GetBillingByOrderId: billingId is empty", nil)
		return
	}

	// get the billing by id
	resp, err := dbHelper.GetBillingByOrderId(&orderId)
	if err != nil {
		logrus.Error("GetBillingByOrderId: failed to get billing by order id", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetBillingByOrderId: failed to get billing by order id", nil)
		return
	}
	common.ReturnResponse(w, "success", http.StatusOK, "", resp)

}

func GetAllBilling(w http.ResponseWriter, r *http.Request) {

	// get all the billing list
	billingList, err := dbHelper.GetAllBilling()
	if err != nil {
		logrus.Error("GetAllBilling: failed to billing getAll entry", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetAllBilling: failed to billing getAll entry", nil)
		return
	}

	common.ReturnResponse(w, "success", http.StatusOK, "", billingList)

}

func UpdateBilling(w http.ResponseWriter, r *http.Request) {
	var billing model.Billing

	err := json.NewDecoder(r.Body).Decode(&billing)
	if err != nil {
		logrus.Error("UpdateBilling: failed to decoding billing for update", err)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "UpdateBilling: failed to decoding billing for update", nil)
		return
	}

	// read user id from path param
	userId := mux.Vars(r)["userId"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(userId); uuidErr != nil {
		logrus.Error("UpdateBilling: Failed to parse the user id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "UpdateBilling: Failed to parse the user id to uuid", nil)
		return
	}

	// read the billing id from the path params
	billingId := mux.Vars(r)["billingId"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(billingId); uuidErr != nil {
		logrus.Error("UpdateBilling: Failed to parse the billing id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "UpdateBilling: Failed to parse the billing id to uuid", nil)
		return
	}

	// update the billing entry
	err = dbHelper.UpdateBilling(&billing, billingId, userId)
	if err != nil {
		logrus.Error("UpdateBilling: failed to updating billing entry ", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "UpdateBilling: failed to updating billing entry ", nil)
		return
	}
	common.ReturnResponse(w, "success", http.StatusNoContent, "", billing)
}

func DeleteBillingById(w http.ResponseWriter, r *http.Request) {
	// read the billing id from the path params
	billingId := mux.Vars(r)["billingId"]

	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(billingId); uuidErr != nil {
		logrus.Error("DeleteBillingById: Failed to parse the billing id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "DeleteBillingById: Failed to parse the billing id to uuid ", nil)
		return
	}

	// read the user id from the path params
	userId := mux.Vars(r)["userId"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(userId); uuidErr != nil {
		logrus.Error("DeleteBillingById: Failed to parse the user id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "DeleteBillingById: Failed to parse the user id to uuid", nil)
		return
	}

	// delete the billing entry by id
	err := dbHelper.DeleteBillingById(&billingId, &userId)
	if err != nil {
		logrus.Error("DeleteBillingById: failed to delete billing entry", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "DeleteBillingById: failed to delete billing entry", nil)
		return
	}
	common.ReturnResponse(w, "success", http.StatusNoContent, "", nil)
}

func DeleteBilling(w http.ResponseWriter, r *http.Request) {
	// read the billing id from the path params
	userId := mux.Vars(r)["userId"]

	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(userId); uuidErr != nil {
		logrus.Error("DeleteBilling: Failed to parse the user id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "DeleteBilling: Failed to parse the user id to uuid", nil)
		return
	}

	// delete the billing entry by id
	err := dbHelper.DeleteBilling(&userId)
	if err != nil {
		logrus.Error("DeleteBilling: failed to delete billing query", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "DeleteBilling: failed to delete billing query", nil)
		return
	}
	common.ReturnResponse(w, "success", http.StatusNoContent, "", nil)
}
