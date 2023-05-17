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

func CreateOrderItem(w http.ResponseWriter, r *http.Request) {
	var orderItem model.OrderItem

	err := json.NewDecoder(r.Body).Decode(&orderItem)
	if err != nil {
		logrus.Error("CreateOrderItem: Failed to decode create orderItem", err)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "CreateOrderItem: Failed to decode create orderItem", nil)
		return
	}

	// read the orderId from query param
	orderId := r.URL.Query().Get("orderId")
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(orderId); uuidErr != nil {
		logrus.Error("CreateOrderItem: Failed to parse the order id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "CreateOrderItem: Failed to parse the order id to uuid", nil)
		return
	}

	// read the foodId from query param
	foodId := r.URL.Query().Get("foodId")
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(foodId); uuidErr != nil {
		logrus.Error("CreateOrderItem: Failed to parse the food id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "CreateOrderItem: Failed to parse the food id to uuid", nil)
		return
	}

	// validate the orderItem field
	validate := validator.New()
	err = validate.Struct(orderItem)
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
		// create the orderItem entry
		orderItemId, orderItemErr := dbHelper.CreateOrderItem(tx, &orderItem)
		if orderItemErr != nil {
			return errors.Wrap(orderItemErr, "CreateUser: failed to create the orderItem entry")
		}

		//create the order orderItem relation entry
		_, err = dbHelper.CreateOrderOrderItem(tx, orderId, *orderItemId)
		if err != nil {
			return errors.Wrap(err, "CreateUser: failed to create the orderItem order relation entry")
		}

		//create the food orderItem relation entry
		_, err = dbHelper.CreateFoodOrderItem(tx, foodId, *orderItemId)
		if err != nil {
			return errors.Wrap(err, "CreateUser: failed to create the orderItem food relation entry")
		}
		return nil
	})
	if txErr != nil {
		logrus.Error("CreateOrderItem: failed to create orderItem entry", txErr)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "CreateOrderItem: failed to create orderItem entry", nil)
		return
	}
	common.ReturnResponse(w, "success", http.StatusCreated, "", orderItem)

}

func GetOrderItemById(w http.ResponseWriter, r *http.Request) {
	// read the orderItem id from path param
	id := mux.Vars(r)["id"]
	if id == "" { // checking orderItem empty or not
		logrus.Error("GetOrderItemById: orderItemId is empty")
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "GetOrderItemById: orderItemId is empty", nil)
		return
	}

	// get the orderItem by id
	resp, err := dbHelper.GetOrderItemById(id)
	if err != nil {
		logrus.Error("GetOrderItemById: failed to get orderItem by id", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetOrderItemById: failed to get orderItem by id", nil)
		return
	}

	common.ReturnResponse(w, "success", http.StatusOK, "", resp)

}

func GetOrderItemByOrderId(w http.ResponseWriter, r *http.Request) {
	// read the orderItem id from path param
	OrderId := mux.Vars(r)["orderId"]
	if OrderId == "" { // checking orderItem empty or not
		logrus.Error("GetOrderItemByOrderId: orderItemId is empty")
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "GetOrderItemByOrderId: orderItemId is empty", nil)
		return
	}

	// get the orderItem by id
	resp, err := dbHelper.GetOrderItemByOrderId(OrderId)
	if err != nil {
		logrus.Error("GetOrderItemByOrderId: failed to get orderItem by id", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetOrderItemByOrderId: failed to get orderItem by id", nil)
		return
	}

	common.ReturnResponse(w, "success", http.StatusOK, "", resp)

}

func GetAllOrderItem(w http.ResponseWriter, r *http.Request) {

	// get all the orderItem list
	orderItemList, err := dbHelper.GetAllOrderItem()
	if err != nil {
		logrus.Error("GetAllOrderItem: failed to getAll orderItem entry", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetAllOrderItem: failed to getAll orderItem entry", nil)
		return
	}

	common.ReturnResponse(w, "success", http.StatusOK, "", orderItemList)

}

func UpdateOrderItem(w http.ResponseWriter, r *http.Request) {
	var orderItem model.OrderItem

	err := json.NewDecoder(r.Body).Decode(&orderItem)
	if err != nil {
		logrus.Error("UpdateOrderItem: failed to decode update orderItem ", err)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "UpdateOrderItem: failed to decode update orderItem", nil)
		return
	}

	// read the orderItem id from the path params
	id := mux.Vars(r)["id"]

	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(id); uuidErr != nil {
		logrus.Error("UpdateOrderItem: Failed to parse the user id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "UpdateOrderItem: Failed to parse the user id to uuid", nil)
		return
	}

	// update the orderItem entry by id
	err = dbHelper.UpdateOrderItem(&orderItem, id)
	if err != nil {
		logrus.Error("UpdateOrderItem: failed to update orderItem ", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "UpdateOrderItem: failed to update orderItem", nil)
		return
	}
	common.ReturnResponse(w, "success", http.StatusNoContent, "", orderItem)
}

func DeleteOrderItem(w http.ResponseWriter, r *http.Request) {
	// read the orderItem id from the path params
	id := mux.Vars(r)["id"]

	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(id); uuidErr != nil {
		logrus.Error("DeleteOrderItem: Failed to parse the user id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "DeleteOrderItem: Failed to parse the user id to uuid", nil)
		return
	}

	// delete the orderItem entry by id from path params
	err := dbHelper.DeleteOrderItem(&id)
	if err != nil {
		logrus.Error("DeleteOrderItem: Failed to delete orderItem", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "DeleteOrderItem: Failed to delete orderItem", nil)
		return
	}
	common.ReturnResponse(w, "success", http.StatusNoContent, "", nil)
}
