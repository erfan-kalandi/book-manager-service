package Authenticate

import (
	"crypto/rand"
	"errors"
	"library/db"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string
	Password string
}

type claims struct {
	jwt.MapClaims
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	TokenString string
	Expiration  time.Time
}

type Auth struct {
	db *db.GormDB
	// jwtSecretKey is the JWT secret key. Each time the server starts, new key is generated
	jwtSecretKey          []byte
	jwtExpirationDuration time.Duration
	logger                *logrus.Logger
}

// NewAuth creates new instance of Auth for authenticating user accounts
func NewAuth(authDB *db.GormDB, jwtExpirationInMinutes int64, logger *logrus.Logger) (*Auth, error) {
	// Generate secret key
	secretKey, err := generateRandomKey()
	if err != nil {
		return nil, err
	}

	// Check the authDB
	if authDB == nil {
		return nil, errors.New("the authenticate database is essential")
	}

	return &Auth{
		db:                    authDB,
		jwtSecretKey:          secretKey,
		jwtExpirationDuration: time.Duration(int64(time.Minute) * jwtExpirationInMinutes),
		logger:                logger,
	}, nil
}

// Login check input details with database.
// If everything was ok then it creates JWT token.
func (a *Auth) Login(cred Credentials) (t *Token, e error) {
	// Check if the user exists
	account, err := a.db.GetUserByUsername(cred.Username)
	if err != nil {
		return nil, err
	}

	// Check Password
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(cred.Password))
	if err != nil {
		return nil, errors.New("the password is not correct")
	}

	// Create the JWT token
	expirationTime := time.Now().Add(a.jwtExpirationDuration)
	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		Username: cred.Username,
		MapClaims: jwt.MapClaims{
			"expired_at": expirationTime.Unix(),
		},
	})

	// Calculate the signed account string format of JWT key
	tokenString, err := tokenJWT.SignedString(a.jwtSecretKey)
	if err != nil {
		return nil, err
	}

	return &Token{
		TokenString: tokenString,
		Expiration:  expirationTime,
	}, nil
}

// GetUsernameByToken gets the token and returns the related username
func (a *Auth) GetUsernameByToken(token string) (*string, error) {
	// Handle empty access token
	if token == "" {
		return nil, errors.New("access denied: the access token is empty")
	}

	// Validate jwt token
	claim, err := a.checkToken(token)
	if err != nil {
		return nil, errors.New("access denied: the access token is not valid")
	}

	// Get the user entity
	user, err := a.db.GetUserByUsername(claim.Username)
	if err != nil {
		return nil, errors.New("access denied: can not fetch such a user")
	}

	return &user.Username, nil
}

// GenerateToken returns a JWT token
func (a *Auth) GenerateToken(cred Credentials) (t *Token, e error) {
	// Create the JWT token
	expirationTime := time.Now().Add(a.jwtExpirationDuration)
	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		Username: cred.Username,
		MapClaims: jwt.MapClaims{
			"expired_at": expirationTime.Unix(),
		},
	})

	// Calculate the signed string format of JWT key
	tokenString, err := tokenJWT.SignedString(a.jwtSecretKey)
	if err != nil {
		return nil, err
	}

	return &Token{
		TokenString: tokenString,
		Expiration:  expirationTime,
	}, nil
}

// checkToken checks validation for incoming tokens
func (a *Auth) checkToken(tokenStr string) (*claims, error) {
	c := &claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, c, func(token *jwt.Token) (interface{}, error) {
		return a.jwtSecretKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid token")
		}

		a.logger.WithError(err).Warn("can not validate the token of the user")
		return nil, errors.New("bad error in validating the token")
	}

	if !tkn.Valid {
		return nil, errors.New("unauthorized")
	}

	return c, nil
}

// generateRandomKey is called to generate another key each time Auth is initialized
func generateRandomKey() ([]byte, error) {
	jwtKey := make([]byte, 32)
	if _, err := rand.Read(jwtKey); err != nil {
		return nil, err
	}
	return jwtKey, nil
}
