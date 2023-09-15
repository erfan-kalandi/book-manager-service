package APIHandlers

import (
	"library/Authenticate"
	"library/db"

	"github.com/sirupsen/logrus"
)

type Server struct {
	Authenticate *Authenticate.Auth
	DB           *db.GormDB
	Logger       *logrus.Logger
}
