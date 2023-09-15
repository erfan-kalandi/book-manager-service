package APIHandlers

import (
	"encoding/json"
	"fmt"
	"io"
	"library/db"
	"net/http"
)

type SignupRequestBody struct {
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"Email"`
	Gender      string `json:"Gender"`
	Password    string `json:"password"`
}

func (s *Server) HandleSignupAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.Logger.Error("http method is not POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("http method should be POST"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.Logger.WithError(err).Warn("can not read the request data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var RQ SignupRequestBody
	err = json.Unmarshal(body, &RQ)
	if err != nil {
		s.Logger.WithError(err).Warn("can not unmarshal the body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := &db.User{
		Username:    RQ.Username,
		Email:       RQ.Email,
		Password:    RQ.Password,
		Firstname:   RQ.FirstName,
		Lastname:    RQ.LastName,
		PhoneNumber: RQ.PhoneNumber,
		Gender:      RQ.Gender,
	}
	fmt.Printf("user: %v\n", user)
	err = s.DB.CreateNewUser(user)
	if err != nil {
		s.Logger.WithError(err).Warn("can not create new user")
		w.WriteHeader(http.StatusBadRequest)
		body := fmt.Sprintf("can not create new user : %s",err.Error())
		w.Write([]byte(body))
		return
	}

	response := map[string]interface{}{
		"message": "user has been created successfully",
	}
	resBody, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusAccepted)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}
