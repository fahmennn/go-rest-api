package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
)

var db *sql.DB

func main() {
	db = DBConnect()
	defer db.Close()

	router := MakeRouter()
	http.ListenAndServe(":8080", router)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  200,
		"message": "Berhasil",
	})
}

func getSiswa(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	result, err := db.Query("SELECT id, name, height, weight FROM player where id = ?", params["id"])

	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var siswa Siswa
	for result.Next() {
		err := result.Scan(&siswa.ID, &siswa.Name, &siswa.Height, &siswa.Weight)
		if err != nil {
			panic(err.Error())
		}
	}
	RenderJson(w, &siswa)
}

func getSiswas(w http.ResponseWriter, r *http.Request) {
	result, err := db.Query("SELECT id, name, height, weight FROM player")

	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	var siswas []*Siswa

	for result.Next() {
		var siswa Siswa
		err := result.Scan(&siswa.ID, &siswa.Name, &siswa.Height, &siswa.Weight)
		if err != nil {
			panic(err.Error())
		}
		siswas = append(siswas, &siswa)
	}
	RenderJson(w, siswas)
}

func createSiswa(w http.ResponseWriter, r *http.Request) {
	var siswa Siswa
	err := jsonapi.UnmarshalPayload(r.Body, &siswa)
	if err != nil {
		panic(err.Error())
	}

	query, err := db.Prepare("INSERT INTO player (name,height,weight) VALUES (?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	result, err := query.Exec(siswa.Name, siswa.Height, siswa.Weight)
	if err != nil {
		panic(err.Error())
	}
	lastId, err := result.LastInsertId()

	if err != nil {
		panic(err.Error())
	}
	siswa.ID = lastId

	RenderJson(w, &siswa)
}

func updateSiswa(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var siswa Siswa
	err := jsonapi.UnmarshalPayload(r.Body, &siswa)
	if err != nil {
		panic(err.Error())
	}
	query, err := db.Prepare("UPDATE player SET name=?, height=?, weight=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}

	result, err := query.Exec(siswa.Name, siswa.Height, siswa.Weight, params["id"])
	if err != nil {
		panic(err.Error())
	}

	siswaId, err := result.LastInsertId()
	if err != nil {
		panic(err.Error())
	}

	siswa.ID = siswaId

	RenderJson(w, &siswa)

}

func deleteSiswa(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	_, err := db.Query("DELETE FROM player WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	RenderJson(w, params["id"]+"was deleted")
}
