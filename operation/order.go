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

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Orders

	//read the userId from path param
	userId := mux.Vars(r)["userId"]
	if userId == "" {
		logrus.Error("userId is empty") //checking userId empty or not
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		logrus.Error("error in decode create order", err)
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
		}
		return
	}

	txErr := common.Tx(func(tx *sqlx.Tx) error {
		// create the order entry
		orderId, err := dbHelper.CreateOrder(tx, &order)
		if err != nil {
			return err
		}
		// create the userOrder relation entry
		_, err = dbHelper.CreateUserOrder(tx, userId, *orderId)
		if err != nil {
			return err
		}
		return err
	})
	if txErr != nil {
		logrus.Error("failed to create order for the user", txErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetOrderByOrderId(w http.ResponseWriter, r *http.Request) {
	//read the userId from path param
	userId := mux.Vars(r)["userId"]
	if userId == "" { // checking orderId empty or not
		logrus.Error("order userId is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// read the order id from path param
	orderId := mux.Vars(r)["orderId"]
	if orderId == "" { // checking orderId empty or not
		logrus.Error("order orderId is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get the order by id
	resp, err := dbHelper.GetOrderById(&orderId)
	if err != nil {
		logrus.Error("failed to get order by id", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		logrus.Error("failed to encode the order ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetOrderByUserId(w http.ResponseWriter, r *http.Request) {
	//read the userId from path param
	userId := mux.Vars(r)["userId"]
	if userId == "" { // checking orderId empty or not
		logrus.Error("order userId is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get the order by id
	resp, err := dbHelper.GetOrderByUserId(&userId)
	if err != nil {
		logrus.Error("failed to get order by id", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		logrus.Error("failed to encode the order ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Orders
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		logrus.Error("error in decode order update", err)
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
	// read the order id from the path params
	orderId := mux.Vars(r)["orderId"]

	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(orderId); uuidErr != nil {
		logrus.Error("Failed to parse the order id to uuid", uuidErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// update the order entry by id from path params
	err = dbHelper.UpdateOrder(&order, orderId, userId)
	if err != nil {
		logrus.Error("error in update order query", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteOrderByUserId(w http.ResponseWriter, r *http.Request) {
	// read the userId from path param
	userId := mux.Vars(r)["userId"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(userId); uuidErr != nil {
		logrus.Error("Failed to parse the address id to uuid", uuidErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// read the order id from the path params
	orderId := mux.Vars(r)["orderId"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(orderId); uuidErr != nil {
		logrus.Error("Failed to parse the user id to uuid", uuidErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//delete the order entry by id from path params
	err := dbHelper.DeleteOrder(&orderId, &userId)
	if err != nil {
		logrus.Error("error in delete order query", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func GetAllOrder(w http.ResponseWriter, r *http.Request) {

	//get all the order list
	orderList, err := dbHelper.GetAllOrder()
	if err != nil {
		logrus.Error("error in get data ", err)
		return
	}
	err = json.NewEncoder(w).Encode(orderList)
	if err != nil {
		logrus.Error("error in encode getAll order", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
