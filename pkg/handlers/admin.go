package handlers

import (
	"encoding/json"
	"fmt"
	"task/pkg/admin"
	"task/pkg/errorsForProject"
	"task/pkg/forms"
	"task/pkg/session"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"

	// "io"
	"net/http"
)

var TokenSecret = []byte("my_secret_key")

type AdminHandler struct {
	Logger         *zap.SugaredLogger
	AdminRepo      admin.AdminsRepo
	SessionManager session.SessionRepo
}

func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	fd := &forms.LoginForm{}
	fd.Login= r.URL.Query().Get("login")
	fd.Password= r.URL.Query().Get("password")
	fmt.Println(fd)
	u, err := h.AdminRepo.Authorize(fd.Login, fd.Password)

	if err == admin.ErrNoAdmin {
		resp, errMarshal := json.Marshal(map[string]interface{}{
			"message": "Admin not found",
		})
		if errMarshal != nil {
			JsonError(w, http.StatusBadRequest, "Login: "+ ErrCantMarshal.Error(), h.Logger)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		w.Write(resp)
		return
	}

	if err == admin.ErrBadPass {
		resp, errMarshal := json.Marshal(map[string]interface{}{
			"message": "invalid password",
		})
		if errMarshal != nil {
			JsonError(w, http.StatusBadRequest, "JsonError: "+ErrCantMarshal.Error(), h.Logger)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		w.Write(resp)
		return
	}
	_, err = h.SessionManager.Create(w, u.ID, r.URL.Path)
	if err != nil {
		h.Logger.Infof("can't create session")
		JsonError(w, http.StatusBadRequest, "JsonError: "+"can't create session", h.Logger)
		return
	}

	resp := GetToken(w, *fd, fmt.Sprint(u.ID), h.Logger)

	w.Write(resp)
}


func (h *AdminHandler) Register(w http.ResponseWriter, r *http.Request) {
	fd := &forms.LoginForm{}
	err := json.NewDecoder(r.Body).Decode(&fd)
	if err != nil {
		JsonError(w, http.StatusBadRequest, "bad request", h.Logger)
		return
	}
	_, err = h.AdminRepo.FindAdmin(fd.Login)

	if err == nil {
		resp, errMarshal := json.Marshal(map[string]interface{}{
			"errors": []errorsForProject.RegisterError{{
				Msg:      "already exists",
				Location: "body",
				Value:    fd.Login,
				Param:    "Adminname",
			},
			}})
		if errMarshal != nil {
			JsonError(w, http.StatusBadRequest, "Register: "+errorsForProject.ErrCantMarshal.Error(), h.Logger)
			return
		}

		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(resp)
		return
	}
	id := h.AdminRepo.NewAdminID()
	h.AdminRepo.Add(&admin.Admin{
		ID:       id,
		Login:    fd.Login,
		Password: fd.Password,
	})
	h.AdminRepo.FindAdmin(fd.Login)

	_, err = h.SessionManager.Create(w, id, r.URL.Path)
	if err != nil {
		h.Logger.Infof("can't create session")
		JsonError(w, http.StatusBadRequest, "JsonError: "+"can't create session", h.Logger)
		return
	}

	resp := GetToken(w, *fd, fmt.Sprint(id), h.Logger)

	w.Write(resp)
}


func GetToken(w http.ResponseWriter, fd forms.LoginForm, id string, Logger *zap.SugaredLogger) (resp []byte) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Admin": jwt.MapClaims{
			"Adminname": fd.Login,
			"id":       id,
		},
		"iat": time.Now().Local().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Local().Unix(),
	})

	tokenString, err := token.SignedString(TokenSecret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//fmt.Println("token ", tokenString)
	resp, err = json.Marshal(map[string]interface{}{
		"token": tokenString,
	})
	if err != nil {
		JsonError(w, http.StatusBadRequest, "Get Token: "+errorsForProject.ErrCantMarshal.Error(), Logger)
		return
	}

	return resp
}