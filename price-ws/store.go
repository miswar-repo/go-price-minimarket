/*
Author : Miswar
*/

package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

//Mendapatkan daftar store
func GetStores(w http.ResponseWriter, r *http.Request) {

	datas := []*storeData{}

	//konfigurasi header
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// get connection database

	db, err := database()
	if err != nil {
		log.Panic(err)
		http.Error(w, "Error , Koneksi Database", 500)

	}

	// inisialisasi variabel
	start, limit := consStart, consLimit
	strLimit := fmt.Sprintf("LIMIT %s", limit)
	strStart := fmt.Sprintf("OFFSET %s", start)

	// mendapatkan parameter
	params := mux.Vars(r)
	if val, ok := params["limit"]; ok {
		limit = val
		strLimit = fmt.Sprintf("LIMIT %s", limit)
	}
	if val, ok := params["start"]; ok {
		start = val
		strStart = fmt.Sprintf("OFFSET %s", start)
	}

	// mendapatkan data dari database
	strQuery := fmt.Sprintf("select id,name from store %s %s", strLimit, strStart)
	rows, err := db.Query(strQuery)
	if err != nil {
		log.Panic(err)
		http.Error(w, "Error, Saat mendapatkan data dari database", 500)
	}
	defer rows.Close()

	// looping untuk mendapatkan semua data
	for rows.Next() {
		data := new(storeData)

		//Read the columns in each row into variables
		err := rows.Scan(&data.Id, &data.Name)
		if err != nil {
			log.Panic(err)
		}
		datas = append(datas, data)
	}

	//cek apakah ada data yang didapat
	if len(datas) == 0 {
		http.Error(w, "Error, Data tidak ada", 404)
	} else {
		json.NewEncoder(w).Encode(datas)
	}

	err = rows.Err()
	if err != nil {
		log.Panic(err)
		http.Error(w, "Error, Saat mendapatkan data dari database", 500)

	}
}
