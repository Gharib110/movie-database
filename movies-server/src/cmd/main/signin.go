package main

import (
	"encoding/json"
	"fmt"
	"github.com/DapperBlondie/movie-server/src/models"
	"github.com/pascaldekloe/jwt"
	zerolog "github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var user *models.User = &models.User{
	ID:       0,
	Email:    "dapper@me.ir",
	Password: tempCryptor(),
}

func createHashPassword(pass string) (string, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(pass), 12)
	if err != nil {
		return "", err
	}

	return string(hashPass), err
}

func compareHashAndPass(hashPass string, pass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(pass))
	if err != nil {
		return bcrypt.ErrMismatchedHashAndPassword
	}

	return nil
}

type ReqUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Application) signIn(w http.ResponseWriter, r *http.Request) {
	var reqUser *ReqUser = &ReqUser{}

	err := json.NewDecoder(r.Body).Decode(reqUser)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return
	}

	err = compareHashAndPass(user.Password, reqUser.Password)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return
	}
	var claims *jwt.Claims = &jwt.Claims{}
	claims.Subject = fmt.Sprint(user.ID)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(time.Hour * 24))
	claims.Issuer = "localhost:8080"
	claims.Audiences = []string{"localhost:8080"}

	jB, err := claims.HMACSign(jwt.HS256, []byte(app.ConfigApp.Jwt.Secret))
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return
	}

	err = app.writeJSON(w, http.StatusOK, jB, "response")
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return
	}

	return
}
