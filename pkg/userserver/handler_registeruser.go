package userserver

import (
	"context"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/jmoiron/sqlx"
	"github.com/klaital/comics/pkg/config"
	"github.com/klaital/comics/pkg/filters"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

type RegisterUserRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserResponse struct {
	JWT string
}

// IsValid checks that the email given looks valid, and that the password meets the minimum requirements
func (r *RegisterUserRequest) IsValid() error {
	tokens := strings.Split(r.Email, "@")
	if len(tokens) != 2 || strings.Index(tokens[1], ".") < 0 {
		return fmt.Errorf("invalid email format")
	}

	// Success!
	return nil
}



type UserRecord struct {
	Email string `db:"email""`
	PasswordDigest string `db:"password_digest"`
}

func (u *UserRecord) Save(ctx context.Context, db *sqlx.DB) error {
	logger := filters.GetContextLogger(ctx).WithFields(log.Fields{
		"operation": "UserRecord#Save",
	})

	sqlStmt := db.Rebind(`INSERT INTO users (email, password_digest) VALUES (:email, :password_digest)`)
	_, err := db.NamedExec(sqlStmt, u)
	if err != nil {
		logger.WithError(err).Error("Failed to insert new user")
		return err
	}
	// Success!
	return nil
}
func RegisterUserHandler(cfg *config.Config) func(request *restful.Request, response *restful.Response) {
	return func(request *restful.Request, response *restful.Response) {
		logger := filters.GetRequestLogger(request).WithField("operation", "RegisterUserHandler")
		db, err := cfg.ConnectDB()
		if err != nil {
			logger.WithError(err).Error("Failed to connect to database")
			response.WriteHeader(http.StatusInternalServerError)
			return
		}

		var input RegisterUserRequest
		if err := request.ReadEntity(&input); err != nil {
			logger.WithError(err).Debug("Failed to deserialize request body")
			response.WriteHeader(http.StatusBadRequest)
			return
		}

		// Ensure the password meets minimum spec
		if err := input.IsValid(); err != nil {
			logger.WithError(err).Debug("Invalid email or password")
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		// Hash the password
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.WithError(err).Error("Failed to hash password")
			response.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Save the user record
		user := UserRecord{
			Email:          input.Email,
			PasswordDigest: string(hash),
		}

		err = user.Save(filters.GetRequestContext(request), db)
		if err != nil {
			// TODO: discern whether the error is a duplicate email - that's a 400 error rather than a 500
			logger.WithError(err).Error("Failed to create new user")
			response.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Success!

	}
}
