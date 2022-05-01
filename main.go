package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/hahattan/assignment/controller"
)

func main() {
	r := mux.NewRouter()
	// ping
	r.HandleFunc("/api/ping", controller.Ping).Methods(http.MethodGet)

	// admin
	c := controller.NewController()
	r.HandleFunc("/api/admin/token", c.AdminOnly(c.CreateToken)).Methods(http.MethodPost)
	r.HandleFunc("/api/admin/token", c.AdminOnly(c.GetAllTokens)).Methods(http.MethodGet)
	r.HandleFunc("/api/admin/token/{token}/disable", c.AdminOnly(c.DisableToken)).Methods(http.MethodPost)

	//public
	r.HandleFunc("/api/login", c.Login).Methods(http.MethodPost)

	log.Println("Web server starting...")
	http.ListenAndServe(":8080", r)
}
