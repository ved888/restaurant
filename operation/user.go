package operation

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"restaurant/common"
	"restaurant/dbHelper"
	"restaurant/model"
	"time"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var userRequest model.UserInterestRequest

	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		logrus.Error("CreateUser: error in decode create user request", err)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "CreateUser: error in decode create user request", nil)
		return
	}

	// validate the user fields
	validate := validator.New()
	err = validate.Struct(userRequest)
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
	existPhone, err := dbHelper.GetUserByPhone(userRequest.User.Phone)
	if existPhone.Phone != "" {
		logrus.Error("CreateUser:user already exist for this phoneNumber")
		common.ReturnResponse(w, "Failed", http.StatusUnprocessableEntity, "CreateUser:user already exist for this phoneNumber", nil)
		return
	}

	existMail, err := dbHelper.GetUserByEmail(userRequest.User.EmailId)
	if existMail.EmailId != "" {
		logrus.Error("CreatUser:user already exist for this emailId")
		common.ReturnResponse(w, "Failed", http.StatusUnprocessableEntity, "CreateUser:user already exist for this emailId", nil)
		return
	}

	txErr := common.Tx(func(tx *sqlx.Tx) error {
		// insert password in database in bcrypt form
		password, err := bcrypt.GenerateFromPassword([]byte(userRequest.User.Password), 14)
		if err != nil {
			logrus.Error("CreateUser:failed to insert password in bcrypt from", err)
			common.ReturnResponse(w, "failed", http.StatusUnprocessableEntity, "CreateUser:failed to insert password in bcrypt from", nil)
			return err
		}
		userRequest.User.Password = string(password)

		// create the user entry
		userID, userErr := dbHelper.CreateUser(tx, &userRequest.User)
		if userErr != nil {
			return errors.Wrap(userErr, "CreateUser: failed to create the user entry")
		}
		userRequest.UserInterest.UsersId = *userID

		// create the interest entry
		interestID, err := dbHelper.CreateInterest(tx, &userRequest.Interest)
		if err != nil {
			return errors.Wrap(err, "CreateUser: failed to create interest entry")
		}
		userRequest.UserInterest.InterestId = *interestID

		// create the user interest relation entry
		_, err = dbHelper.CreateUserInterest(tx, &userRequest.UserInterest)
		return errors.Wrap(err, "CreateUser: failed to create user interest relation entry")
	})
	if txErr != nil {
		logrus.Error(" CreateUser: failed to create user and interest", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, " CreateUser: failed to create user and interest", nil)
		return
	}
	common.ReturnResponse(w, "success", http.StatusCreated, "", userRequest)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var data map[string]string

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		logrus.Error("Login:failed to decoding data", err)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "Login:failed to decoding data", nil)
		return
	}

	if data["email"] == "" {
		logrus.Error("Login:please send the valid email")
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "Login:please send valid email", nil)
		return
	}

	if data["password"] == "" {
		logrus.Error("Login:please send the valid password")
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "Login:please send valid password", nil)
		return
	}
	//
	user, err := dbHelper.GetUserByEmail(data["email"])
	fmt.Println("user data ", user)
	if user.EmailId == "" { //checking this email exist in database or not
		logrus.Error("Login:user not found for this email", err)
		common.ReturnResponse(w, "failed", http.StatusNotFound, "Login:user not found for this email", nil)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		logrus.Error("Login:incorrect password", err)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "Login:incorrect password", nil)
		return
	}

	clams := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    "Restaurant App",
		Id:        user.Id.String(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //token valid for 1 day
	})
	token, err := clams.SignedString([]byte(common.SecretKey))
	if err != nil {
		fmt.Println("error at SignedString ", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "code not login", nil)
		return
	}

	common.ReturnResponse(w, "success", http.StatusOK, "", map[string]string{
		"token ": token,
	})

	return
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var usersRequest model.UserInterestRequest

	err := json.NewDecoder(r.Body).Decode(&usersRequest)
	if err != nil {
		logrus.Error("UpdateUser: failed to decode update user", err)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "UpdateUser: failed to decode update user", nil)
		return
	}

	// read the user id from the path params
	userId := mux.Vars(r)["userId"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(userId); uuidErr != nil {
		logrus.Error("UpdateUser: Failed to parse the user id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "UpdateUser: Failed to parse the user id to uuid", nil)
		return
	}

	txErr := common.Tx(func(tx *sqlx.Tx) error {

		// update password save in bcrypt from
		password, err := bcrypt.GenerateFromPassword([]byte(usersRequest.User.Password), 14)
		if err != nil {
			logrus.Error("UpdateUser:failed to save update password in bcrypt from", err)
			common.ReturnResponse(w, "failed", http.StatusUnprocessableEntity, "UpdateUser:failed to save update password in bcrypt from", nil)
			return err
		}
		usersRequest.User.Password = string(password)

		// update the user entry
		userErr := dbHelper.UpdateUser(tx, &usersRequest.User, userId)
		if userErr != nil {
			return errors.Wrap(userErr, "UpdateUser: failed to update the user entry")
		}

		// update the interest entry
		err = dbHelper.UpdateInterest(tx, &usersRequest.Interest, userId)
		if err != nil {
			return errors.Wrap(err, "updateUser: failed to update interest entry")
		}

		return err
	})

	if txErr != nil {
		logrus.Error(" updateUser: failed to update user and interest", err)
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "UpdateUser:failed to update user and interest", nil)
		return
	}
	common.ReturnResponse(w, "success", http.StatusNoContent, "", usersRequest)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	// reading user id from path parameters
	userId := mux.Vars(r)["userId"]
	// check the id is of uuid type or not
	if _, uuidErr := uuid.Parse(userId); uuidErr != nil {
		logrus.Error("GetUserById: Failed to parse the user id to uuid", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "GetUserById: Failed to parse the user id to uuid", nil)
		return
	}

	// get the user by id
	user, err := dbHelper.GetUserByID(&userId)
	if err != nil {
		logrus.Error("GetUserById: failed to getUser")
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetUserById: failed to getUser", nil)
		return
	}

	// get the userInterest by user_id
	interest, err := dbHelper.GetInterestByUserId(&userId)
	if err != nil {
		logrus.Error("GetUserById: failed to getInterest the user")
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetUserById: failed to getInterest the user", nil)
		return
	}

	resp := model.UserInterestRequest{
		User:     user,
		Interest: *interest,
	}

	common.ReturnResponse(w, "success", http.StatusOK, "", resp)
}

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	// get all the user list
	user, err := dbHelper.GetAllUser()
	if err != nil {
		logrus.Error("GetAllUser: failed to getAll the User list ")
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "GetAllUser: failed to getAll the User list", nil)
		return
	}

	// get all the userInterest list
	interest, err := dbHelper.GetAllInterest()
	if err != nil {
		logrus.Error("GetAllUser: failed to getInterest the user")
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "GetAllUser: failed to getInterest the user", nil)
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

	common.ReturnResponse(w, "success", http.StatusOK, "", userInterestList)
}

func UserDelete(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["userId"]
	if _, uuidErr := uuid.Parse(userId); uuidErr != nil {
		logrus.Error("UserDelete: Failed to parse the user id", uuidErr)
		common.ReturnResponse(w, "failed", http.StatusBadRequest, "UserDelete: Failed to parse the user id", nil)
		return
	}

	// delete the user by id
	txErr := common.Tx(func(tx *sqlx.Tx) error {
		// delete the user entry
		userErr := dbHelper.DeleteUser(tx, &userId)
		if userErr != nil {
			return errors.Wrap(userErr, "DeleteUser: failed to delete the user entry")
		}

		// delete the interest entry
		err := dbHelper.DeleteInterest(tx, &userId)
		if err != nil {
			return errors.Wrap(err, "DeleteUser: failed to delete interest entry")
		}

		return err
	})

	if txErr != nil {
		logrus.Error(" DeleteUser: failed to delete user and interest")
		common.ReturnResponse(w, "failed", http.StatusInternalServerError, "DeleteUser: failed to delete user and interest", nil)
		return
	}
	common.ReturnResponse(w, "success", http.StatusNoContent, "", nil)
}
