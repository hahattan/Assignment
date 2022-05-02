package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/hahattan/assignment/application"
	"github.com/hahattan/assignment/interfaces"
)

type Controller struct {
	dbClient interfaces.DBClient
	username string
	password string
}

func NewController(client interfaces.DBClient) *Controller {
	username := os.Getenv("ADMIN_USERNAME")
	if username == "" {
		username = "root"
	}
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		password = "root"
	}

	return &Controller{
		dbClient: client,
		username: username,
		password: password,
	}
}

func (c *Controller) CreateToken(w http.ResponseWriter, _ *http.Request) {
	token, err := application.CreateToken(c.dbClient)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	bytes, err := json.Marshal(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	log.Printf("token %s created", token)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func (c *Controller) DisableToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	exist, err := c.dbClient.TokenExists(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Failed to check token existence: ", err.Error())
		fmt.Fprint(w, err.Error())
		return
	}

	if !exist {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("token %s not existed", token)
		fmt.Fprint(w, "token not existed")
		return
	}

	err = application.DisableToken(token, c.dbClient)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Failed to update token:", err.Error())
		fmt.Fprint(w, err.Error())
		return
	}

	log.Printf("token %s updated", token)
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) GetAllTokens(w http.ResponseWriter, _ *http.Request) {
	res, err := application.GetAllTokens(c.dbClient)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Failed to get tokens:", err.Error())
		fmt.Fprint(w, err.Error())
		return
	}

	bytes, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}

	var token string
	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	exist, err := c.dbClient.TokenExists(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Failed to check token existence: ", err.Error())
		fmt.Fprint(w, err.Error())
		return
	}

	if !exist {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("token %s not existed", token)
		fmt.Fprint(w, "token not existed")
		return
	}

	err = application.ValidateToken(token, c.dbClient)
	if err != nil {
		if err.Error() != "token disabled" && err.Error() != "token expired" {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Failed to get token's expiry:", err.Error())
			fmt.Fprint(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("Failed to parse basic authentication info")
			fmt.Fprint(w, "Error parsing basic auth")
			return
		}

		if username != c.username || password != c.password {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
