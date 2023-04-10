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

func CreateFood(w http.ResponseWriter, r *http.Request) {
	var food model.Food

	err := json.NewDecoder(r.Body).Decode(&food)
	if err != nil {
		logrus.Error("CreateFood: Failed to decode food", err)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "CreateFood: Failed to decode food", nil)
		return
	}

	// validate the food field
	validate := validator.New()
	err = validate.Struct(food)
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

	// create the food entry
	err = dbHelper.CreateFood(&food)
	if err != nil {
		logrus.Error("CreateFood: Failed to create food entry", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "CreateFood: Failed to create food entry", nil)
		return
	}

	common.ReturnResponse(w, "success", http.StatusCreated, "", food)
}

func GetFoodById(w http.ResponseWriter, r *http.Request) {
	// read the food id from path param
	foodId := mux.Vars(r)["foodId"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(foodId); uuidErr != nil {
		logrus.Error("GetFoodById: Failed to parse the food id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "GetFoodById: Failed to parse the food id to uuid", nil)
		return
	}

	// get the food by id
	resp, err := dbHelper.GetFoodById(foodId)
	if err != nil {
		logrus.Error("GetFoodById: failed to get food by food id", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetFoodById: failed to get food by food id", nil)
		return
	}

	common.ReturnResponse(w, "success", http.StatusOK, "", resp)

}

func GetFoodByOrderItemId(w http.ResponseWriter, r *http.Request) {
	// read the orderItem id from path param
	orderItemId := mux.Vars(r)["orderItemId"]
	if orderItemId == "" { // checking foodId empty or not
		logrus.Error("GetFoodByOrderItemId: orderItem id is empty")
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "GetFoodByOrderItemId: orderItem id is empty", nil)
		return
	}

	// get the food by id
	resp, err := dbHelper.GetFoodByOrderItemId(orderItemId)
	if err != nil {
		logrus.Error("GetFoodByOrderItemId: failed to get food by OrderItem id", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetFoodByOrderItemId: failed to get food by OrderItem id", nil)
		return
	}

	common.ReturnResponse(w, "success", http.StatusOK, "", resp)

}

func GetAllFood(w http.ResponseWriter, r *http.Request) {

	// get all the food list
	foodList, err := dbHelper.GetAllFood()
	if err != nil {
		logrus.Error("GetAllFood: Failed to getAll food ", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetAllFood: Failed to getAll food ", nil)
		return
	}

	common.ReturnResponse(w, "success", http.StatusOK, "", foodList)

}

func UpdateFood(w http.ResponseWriter, r *http.Request) {

	var food model.Food

	err := json.NewDecoder(r.Body).Decode(&food)
	if err != nil {
		logrus.Fatal("UpdateFood: Failed to decode food for update", err)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "UpdateFood: Failed to decode food for update", nil)
		return
	}

	// read the food id from the path params
	id := mux.Vars(r)["id"]

	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(id); uuidErr != nil {
		logrus.Error("UpdateFood: Failed to parse the food id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "UpdateFood: Failed to parse the food id to uuid", nil)
		return
	}

	// update the food entry by id
	err = dbHelper.UpdateFood(&food, id)
	if err != nil {
		logrus.Fatal("UpdateFood: Failed to update food entry", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "UpdateFood: Failed to update food entry", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteFood(w http.ResponseWriter, r *http.Request) {
	// reade the food id from path param
	id := mux.Vars(r)["id"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(id); uuidErr != nil {
		logrus.Error("DeleteFood: Failed to parse the  id")
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "DeleteFood: Failed to parse the  id", nil)
		return
	}

	// delete the food entry by id
	err := dbHelper.DeleteFood(&id)
	if err != nil {
		logrus.Error("DeleteFood: Failed to delete food entry", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "DeleteFood: Failed to delete food entry", nil)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
