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

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Orders

	// read the userId from path param
	userId := mux.Vars(r)["userId"]
	if userId == "" {
		logrus.Error("CreateOrder: userId is empty") // checking userId empty or not
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "CreateOrder: userId is empty", nil)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		logrus.Error("CreateOrder: Failed to decode create order", err)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "CreateOrder: Failed to decode create order", nil)
		return
	}
	// validate the order field
	validate := validator.New()
	err = validate.Struct(order)
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

	txErr := common.Tx(func(tx *sqlx.Tx) error {
		// create the order entry
		orderId, err := dbHelper.CreateOrder(tx, &order)
		if err != nil {
			return errors.Wrap(err, "CreateOrder: failed to create the order entry")
		}
		// create the userOrder relation entry
		_, err = dbHelper.CreateUserOrder(tx, userId, *orderId)
		if err != nil {
			return errors.Wrap(err, "CreateOrder: failed to create the user order relation entry")
		}
		return err
	})
	if txErr != nil {
		logrus.Error("CreateOrder: failed to create order for the user", txErr)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "CreateOrder: failed to create order for the user", nil)
		return
	}
	common.ReturnResponse(w, "success", http.StatusCreated, "", order)
}

func GetOrderByOrderId(w http.ResponseWriter, r *http.Request) {
	// read the userId from path param
	userId := mux.Vars(r)["userId"]
	if userId == "" { // checking orderId empty or not
		logrus.Error("GetOrderByOrderId: order userId is empty")
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "GetOrderByOrderId: order userId is empty", nil)
		return
	}

	// read the order id from path param
	orderId := mux.Vars(r)["orderId"]
	if orderId == "" { // checking orderId empty or not
		logrus.Error("GetOrderByOrderId: order orderId is empty")
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "GetOrderByOrderId: order orderId is empty", nil)
		return
	}

	// get the order by id
	resp, err := dbHelper.GetOrderById(&orderId)
	if err != nil {
		logrus.Error("GetOrderByOrderId: failed to get order entry by id", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetOrderByOrderId: failed to get order entry by id", nil)
		return
	}

	common.ReturnResponse(w, "success", http.StatusOK, "", resp)
}

func GetOrderByUserId(w http.ResponseWriter, r *http.Request) {
	// read the userId from path param
	userId := mux.Vars(r)["userId"]
	if userId == "" { // checking orderId empty or not
		logrus.Error("GetOrderByUserId: order userId is empty")
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "GetOrderByUserId: order userId is empty", nil)
		return
	}

	// get the order by id
	resp, err := dbHelper.GetOrderByUserId(&userId)
	if err != nil {
		logrus.Error("GetOrderByUserId: failed to get order by id", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetOrderByUserId: failed to get order by id", nil)
		return
	}

	common.ReturnResponse(w, "success", http.StatusOK, "", resp)
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Orders
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		logrus.Error("UpdateOrder: Failed to decode order ", err)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "UpdateOrder: Failed to decode order", nil)
		return
	}

	// read the userId from path param
	userId := mux.Vars(r)["userId"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(userId); uuidErr != nil {
		logrus.Error("UpdateOrder: Failed to parse the user id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "UpdateOrder: Failed to parse the user id to uuid", nil)
		return
	}

	// read the order id from the path params
	orderId := mux.Vars(r)["orderId"]

	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(orderId); uuidErr != nil {
		logrus.Error("UpdateOrder: Failed to parse the order id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "UpdateOrder: Failed to parse the order id to uuid", nil)
		return
	}

	// update the order entry by id from path params
	err = dbHelper.UpdateOrder(&order, orderId, userId)
	if err != nil {
		logrus.Error("UpdateOrder: Failed to update order entry", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "UpdateOrder: Failed to update order entry", nil)
		return
	}
	common.ReturnResponse(w, "success", http.StatusNoContent, "", order)
}

func DeleteOrderByUserId(w http.ResponseWriter, r *http.Request) {
	// read the userId from path param
	userId := mux.Vars(r)["userId"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(userId); uuidErr != nil {
		logrus.Error("DeleteOrderByUserId:failed to parse the user id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "DeleteOrderByUserId:failed to parse the user id to uuid", nil)
		return
	}

	// read the order id from the path params
	orderId := mux.Vars(r)["orderId"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(orderId); uuidErr != nil {
		logrus.Error("DeleteOrderByUserId: Failed to parse the order id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "DeleteOrderByUserId: Failed to parse the order id to uuid", nil)
		return
	}

	// delete the order entry by id from path params
	err := dbHelper.DeleteOrder(&orderId, &userId)
	if err != nil {
		logrus.Error("DeleteOrderByUserId: Failed to delete order entry", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "DeleteOrderByUserId: Failed to delete order entry", nil)
		return
	}
	common.ReturnResponse(w, "success", http.StatusNoContent, "", nil)
}

func GetAllOrder(w http.ResponseWriter, r *http.Request) {

	// get all the order list
	orderList, err := dbHelper.GetAllOrder()
	if err != nil {
		logrus.Error("GetAllOrder: Failed to get data ", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetAllOrder: Failed to get data ", nil)
		return
	}

	common.ReturnResponse(w, "success", http.StatusOK, " ", orderList)

}
