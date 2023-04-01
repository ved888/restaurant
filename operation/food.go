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

func CreateFood(w http.ResponseWriter, r *http.Request) {
	var food model.Food

	err := json.NewDecoder(r.Body).Decode(&food)
	if err != nil {
		logrus.Error("error in decode food", err)
		return
	}
	//validate the food field
	validate := validator.New()
	err = validate.Struct(food)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		w.WriteHeader(http.StatusBadRequest)
		responseBody := map[string]string{"error": validationErrors.Error()}
		logrus.Error(responseBody)
		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		}
		return
	}
	// create the food entry
	err = dbHelper.CreateFood(&food)
	if err != nil {
		logrus.Error("error in create food calling function", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(food)
	if err != nil {
		logrus.Error("error in encoding get all food ", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func GetFoodById(w http.ResponseWriter, r *http.Request) {
	// read the food id from path param
	foodId := mux.Vars(r)["foodId"]
	if foodId == "" { // checking foodId empty or not
		logrus.Error("food id is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// get the food by id
	resp, err := dbHelper.GetFoodById(foodId)
	if err != nil {
		logrus.Error("failed to get food by id", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		logrus.Error("failed to encode the address ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetFoodByOrderItemId(w http.ResponseWriter, r *http.Request) {
	// read the orderItem id from path param
	orderItemId := mux.Vars(r)["orderItemId"]
	if orderItemId == "" { // checking foodId empty or not
		logrus.Error("orderItem id is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// get the food by id
	resp, err := dbHelper.GetFoodByOrderItemId(orderItemId)
	if err != nil {
		logrus.Error("failed to get food by id", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		logrus.Error("failed to encode the address ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetAllFood(w http.ResponseWriter, r *http.Request) {

	//get all the food list
	foodList, err := dbHelper.GetAllFood()
	if err != nil {
		logrus.Error("error in getAll food calling function", err)

	}
	err = json.NewEncoder(w).Encode(foodList)
	if err != nil {
		logrus.Error("error in encoding get all food ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func FoodUpdate(w http.ResponseWriter, r *http.Request) {

	var food model.Food

	err := json.NewDecoder(r.Body).Decode(&food)
	if err != nil {
		logrus.Fatal("error in decode food for update", err)
	}
	// read the food id from the path params
	id := mux.Vars(r)["id"]

	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(id); uuidErr != nil {
		logrus.Error("Failed to parse the user id to uuid", uuidErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//update the food entry by id
	err = dbHelper.UpdateFood(&food, id)
	if err != nil {
		logrus.Fatal("error in update food calling function", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

func FoodDelete(w http.ResponseWriter, r *http.Request) {
	//reade the food id from path param
	id := mux.Vars(r)["id"]
	//check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(id); uuidErr != nil {
		logrus.Error("UserDelete: Failed to parse the user id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// delete the food entry by id
	err := dbHelper.DeleteFood(&id)
	if err != nil {
		logrus.Error("error in delete food", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
