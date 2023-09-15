package APIHandlers

import (
	"encoding/json"
	"net/http"
)

type userInfo struct {
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
}

func (s *Server) HandleProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.Logger.Error("http method is not GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("http method should be GET"))
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		s.Logger.Error("can not find Authorization")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("there is no token for Authorization"))
		return
	}

	username, err := s.Authenticate.GetUsernameByToken(token)
	if err != nil {
		s.Logger.Error("can not find user by this Authorization")
		w.WriteHeader(http.StatusUnauthorized)
		s.Logger.Error("can not find user by this token")
		return
	}

	user, err := s.DB.GetUserByUsername(*username)
	if err != nil {
		s.Logger.WithError(err).Error("some internal error in find user")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Create the response body
	res, _ := json.Marshal(&userInfo{
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		FirstName:   user.Firstname,
		LastName:    user.Lastname,
	})
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
