package APIHandlers

import (
	"encoding/json"
	"fmt"
	"io"
	"library/Authenticate"
	"net/http"
)

type loginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
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

	var RQ loginRequestBody
	err = json.Unmarshal(body, &RQ)
	if err != nil {
		s.Logger.WithError(err).Warn("can not unmarshal the body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := s.Authenticate.Login(Authenticate.Credentials{
		Username: RQ.Username,
		Password: RQ.Password,
	})
	if err != nil {
		s.Logger.WithError(err).Warn("password is not correct")
		w.WriteHeader(http.StatusForbidden)
		body := fmt.Sprintf("password is not correct : %s",err.Error())
		w.Write([]byte(body))
		return
	}

	res := map[string]interface{}{
		"access_token": token.TokenString,
	}
	resBody, _ := json.Marshal(res)
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}
