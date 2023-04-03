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

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var userRequest model.UserInterestRequest

	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		logrus.Error("CreateUser: error in decode create user request", err)
		return
	}
	//validate the user fields
	validate := validator.New()
	err = validate.Struct(userRequest)
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
		// create the user entry
		userID, userErr := dbHelper.CreateUser(tx, &userRequest.User)
		if userErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return userErr
		}
		userRequest.UserInterest.UsersId = *userID

		// create the interest entry
		interestID, err := dbHelper.CreateInterest(tx, &userRequest.Interest)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
		userRequest.UserInterest.InterestId = *interestID

		// create the user interest relation entry
		_, err = dbHelper.CreateUserInterest(tx, &userRequest.UserInterest)
		return err
	})

	if txErr != nil {
		logrus.Error(" CreateUser: failed to create user and interest", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var usersRequest model.UserInterestRequest

	err := json.NewDecoder(r.Body).Decode(&usersRequest)
	if err != nil {
		logrus.Error("error in decode update user", err)
		return
	}
	// read the user id from the path params
	id := mux.Vars(r)["id"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(id); uuidErr != nil {
		logrus.Error("Failed to parse the user id to uuid", uuidErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// update the user entry by id
	err = dbHelper.UpdateUser(&usersRequest.User, id)
	if err != nil {
		logrus.Error("failed to update the user entry", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// update the interest entry by userId
	err = dbHelper.UpdateInterest(&usersRequest.Interest, id)
	if err != nil {
		logrus.Error("failed to update the interest entry", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	// reading user id from path parameters
	userId := mux.Vars(r)["userId"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(userId); uuidErr != nil {
		logrus.Error("Failed to parse the user id to uuid", uuidErr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// get the user by id
	user, err := dbHelper.GetUserByID(&userId)
	if err != nil {
		logrus.Error("failed to getUser")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// get the userInterest by user_id
	interest, err := dbHelper.GetInterestByUserId(&userId)
	if err != nil {
		logrus.Error("failed to getInterest the user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := model.UserInterestRequest{
		User:     *user,
		Interest: *interest,
	}
	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		logrus.Error("failed to encode the  getUserInterest ")

		return
	}
}

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	// get all the user list
	user, err := dbHelper.GetAllUser()
	if err != nil {
		logrus.Error("failed to getAll the User list ")
		return
	}
	// get all the userInterest list
	interest, err := dbHelper.GetAllInterest()
	if err != nil {
		logrus.Error("failed to getInterest the user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userMap := make(map[string]model.Users)
	for i := range user {
		userMap[user[i].Id.String()] = *user[i]
	}

	userInterestList := make([]model.UserInterestRequest, 0)
	for i := range interest {
		if val, ok := userMap[interest[i].UserID.String()]; ok {
			userInterestList = append(userInterestList, model.UserInterestRequest{
				User: val,
				Interest: model.Interest{
					ID:   interest[i].ID,
					Name: interest[i].Name,
					Type: interest[i].Type,
				},
			})
		}
	}

	err = json.NewEncoder(w).Encode(&userInterestList)
	if err != nil {
		logrus.Error("failed to encode the  getAllUserInterest ")

		return
	}

}

func UserDelete(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["userId"]
	if _, uuidErr := uuid.Parse(userId); uuidErr != nil {
		logrus.Error("UserDelete: Failed to parse the user id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// delete the user by id
	err := dbHelper.DeleteUser(&userId)
	if err != nil {
		logrus.Error("Failed to  delete user", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// delete the userInterest by userId
	err = dbHelper.DeleteInterest(&userId)
	if err != nil {
		logrus.Error("Failed to  delete interest", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
