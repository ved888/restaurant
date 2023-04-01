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

func CreateOrderItem(w http.ResponseWriter, r *http.Request) {
	var orderItem model.OrderItem

	err := json.NewDecoder(r.Body).Decode(&orderItem)
	if err != nil {
		logrus.Error("error in decode create orderItem", err)
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

	// read the foodId from query param
	foodId := r.URL.Query().Get("foodId")
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(foodId); uuidErr != nil {
		logrus.Error("Failed to parse the food id to uuid", uuidErr)
		w.WriteHeader(http.StatusBadRequest)
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
		}
		return
	}
	txErr := common.Tx(func(tx *sqlx.Tx) error {
		// create the orderItem entry
		orderItemId, err := dbHelper.CreateOrderItem(tx, &orderItem)
		if err != nil {
			return err
		}

		//create the order orderItem relation entry
		_, err = dbHelper.CreateOrderOrderItem(tx, orderId, *orderItemId)
		if err != nil {
			return err
		}

		//create the food orderItem relation entry
		_, err = dbHelper.CreateFoodOrderItem(tx, foodId, *orderItemId)
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

func GetOrderItemById(w http.ResponseWriter, r *http.Request) {
	// read the orderItem id from path param
	id := mux.Vars(r)["id"]
	if id == "" { // checking orderItem empty or not
		logrus.Error("orderItemId is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// get the orderItem by id
	resp, err := dbHelper.GetOrderItemById(id)
	if err != nil {
		logrus.Error("failed to get orderItem by id", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		logrus.Error("failed to encode the orderItem ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetOrderItemByOrderId(w http.ResponseWriter, r *http.Request) {
	// read the orderItem id from path param
	OrderId := mux.Vars(r)["orderId"]
	if OrderId == "" { // checking orderItem empty or not
		logrus.Error("orderItemId is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// get the orderItem by id
	resp, err := dbHelper.GetOrderItemByOrderId(OrderId)
	if err != nil {
		logrus.Error("failed to get orderItem by id", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		logrus.Error("failed to encode the orderItem ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetAllOrderItem(w http.ResponseWriter, r *http.Request) {

	// get all the orderItem list
	orderItemList, err := dbHelper.GetAllOrderItem()
	if err != nil {
		logrus.Error("error in getAll orderItem query", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(orderItemList)
	if err != nil {
		logrus.Error("error in encode getAll orderItem", err)
		return
	}

}
func UpdateOrderItem(w http.ResponseWriter, r *http.Request) {
	var orderItem model.OrderItem

	err := json.NewDecoder(r.Body).Decode(&orderItem)
	if err != nil {
		logrus.Error("error in decode update orderItem ", err)
		return
	}
	// read the orderItem id from the path params
	id := mux.Vars(r)["id"]

	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(id); uuidErr != nil {
		logrus.Error("Failed to parse the user id to uuid", uuidErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// update the orderItem entry by id
	err = dbHelper.UpdateOrderItem(&orderItem, id)
	if err != nil {
		logrus.Error("error in update orderItem ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteOrderItem(w http.ResponseWriter, r *http.Request) {
	// read the orderItem id from the path params
	id := mux.Vars(r)["id"]

	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(id); uuidErr != nil {
		logrus.Error("Failed to parse the user id to uuid", uuidErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// delete the orderItem entry by id from path params
	err := dbHelper.DeleteOrderItem(&id)
	if err != nil {
		logrus.Error("error in delete orderItem ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
