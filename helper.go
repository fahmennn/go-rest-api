package main

import (
	"database/sql"
	"net/http"

	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
)

type Siswa struct {
	ID     int64  `jsonapi:"primary,siswas"`
	Name   string `jsonapi:"attr,name"`
	Height int    `jsonapi:"attr,height"`
	Weight int    `jsonapi:"attr,weight"`
}

func MakeRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", handleHome).Methods("GET")
	router.HandleFunc("/api/siswa", getSiswas).Methods("GET")
	router.HandleFunc("/api/siswa/{id}", getSiswa).Methods("GET")
	router.HandleFunc("/api/siswa", createSiswa).Methods("POST")
	router.HandleFunc("/api/siswa{id}", updateSiswa).Methods("PUT")
	router.HandleFunc("/api/siswa{id}", deleteSiswa).Methods("DELETE")
	return router
}

func DBConnect() *sql.DB {
	username := "root"
	password := ""
	server := "localhost"
	name := "fpro"
	db, err := sql.Open("mysql", username+":"+password+"@tcp("+server+")/"+name)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func RenderJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("content-type", "application/json")

	jsonapi.MarshalPayload(w, data)
}
