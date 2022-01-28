package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	Id        int
	IncBegin  string
	IncEnd    string
	IncLength int
	Desc      string
	Comments  string
}

type JsonResponse struct {
	Type    string
	Message string
	Data    []DB
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/incidents", getAllIncidents).Methods("GET")
	router.HandleFunc("/incidents", deleteRows).Methods("DELETE")
	router.HandleFunc("/incidents", writeIncident).Methods("POST")

	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"POST", "GET", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*/*"})

	fmt.Println("Server at 8000")
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(credentials, methods, origins)(router)))

}

func setupDB() *sql.DB {
	db, err := sql.Open("mysql", "root:my_name_is_ivan@/incidents")
	checkErr(err)

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}

func writeIncident(w http.ResponseWriter, r *http.Request) {

	db := setupDB()
	var response JsonResponse

	id := r.FormValue("id")
	incBegin := r.FormValue("incBegin")
	incEnd := r.FormValue("incEnd")
	incLength := r.FormValue("incLength")
	desc := r.FormValue("desc")
	comments := r.FormValue("comments")

	if id == "" || incBegin == "" || incEnd == "" || incLength == "" {
		response = JsonResponse{Type: "Error", Message: "You are missing one of the parameters"}
	} else {
		fmt.Println("Adding incident to DB")
		_, err := db.Query("INSERT INTO incidents VALUES (?,?,?,?,?,?)", id, incBegin, incEnd, incLength, desc, comments)
		checkErr(err)
		fmt.Println("Done")
		response = JsonResponse{Type: "Success", Message: "Incident has been stored"}
	}
	json.NewEncoder(w).Encode(response)
}

func deleteRows(w http.ResponseWriter, r *http.Request) {

	var response JsonResponse
	params := mux.Vars(r)
	db := setupDB()
	id := params["id"]

	if id == "" {
		response = JsonResponse{Type: "Error", Message: "No such incident"}
	} else {
		fmt.Println("Deleting incident from DB")
		_, err := db.Query("DELETE FROM incidents WHERE id=?", id)
		checkErr(err)
		fmt.Println("Done")

		response = JsonResponse{Type: "success", Message: "Incident was deleted successfully"}
	}

	json.NewEncoder(w).Encode(response)

}

func getAllIncidents(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	fmt.Println("Getting all the incidents in DB")
	incidents, err := db.Query("SELECT * FROM incidents ")
	checkErr(err)
	fmt.Println("Done")

	var result []DB

	for incidents.Next() {
		var id int
		var incBegin string
		var incEnd string
		var incLength int
		var desc string
		var comments string

		err = incidents.Scan(&id, &incBegin, &incEnd, &incLength, &desc, &comments)
		checkErr(err)
		result = append(result, DB{Id: id, IncBegin: incBegin, IncEnd: incEnd, IncLength: incLength, Desc: desc, Comments: comments})
	}
	var response = JsonResponse{Type: "Success", Data: result}
	json.NewEncoder(w).Encode(response)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
