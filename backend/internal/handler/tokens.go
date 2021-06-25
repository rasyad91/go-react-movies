package handler

import (
	"backend/internal/model"
	"backend/internal/util"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
)

var validUser = model.User{
	ID:       10,
	Email:    "me@here.com",
	Password: "$2a$12$zQrsu7zNYRJM.KdhHWyWYeypbzhBmT5kgYK9/x.q0IihhjHwZe3YG",
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (m *Repository) Signin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		m.App.Logger.Println(err)
		util.ErrorJSON(w, errors.New("unathorized"))
		return
	}

	hashedPassword := validUser.Password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password)); err != nil {
		m.App.Logger.Println(err)
		util.ErrorJSON(w, errors.New("unathorized"))
		return
	}

	var claims jwt.Claims
	claims.Subject = fmt.Sprint(validUser.ID)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = "mydomain.com"
	claims.Audiences = []string{"mydomain.com"}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(m.App.JWT.Secret))
	if err != nil {
		m.App.Logger.Println(err)
		util.ErrorJSON(w, errors.New("error signin"))
		return
	}

	if err := util.WriteJSON(w, "response", jwtBytes); err != nil {
		m.App.Logger.Println(err)
		util.ErrorJSON(w, err)
		return
	}

}
