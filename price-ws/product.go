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

//Mendapatkan daftar product
func GetProducts(w http.ResponseWriter, r *http.Request) {

	datas := []*productData{}

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
	start, limit, store, strStore := consStart, consLimit, "", ""
	strLimit := fmt.Sprintf("LIMIT %s", limit)
	strStart := fmt.Sprintf("OFFSET %s", start)
	bolError := false

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
	if val, ok := params["store"]; ok {
		store = val
		strStore = fmt.Sprintf("where store_id ='%s' ", store)
	}

	if store != "" {
		// cek id store apakah ada di list
		cStore := 0
		strQuery := "select count(*) as cStore from store where id = ? "
		err = db.QueryRow(strQuery, store).Scan(&cStore)
		if err != nil {
			log.Panic(err)
			http.Error(w, "Error, Saat mendapatkan data dari database", 500)
			bolError = true
		} else if cStore == 0 {
			http.Error(w, "Error, Store ("+store+") tidak ada", 500)
			bolError = true
		}
	}

	if !bolError {
		// mendapatkan data dari database
		strQuery := fmt.Sprintf("select store_id,id,name,harga,url,Image,update_date from product %s %s %s", strStore, strLimit, strStart)
		rows, err := db.Query(strQuery)
		if err != nil {
			log.Panic(err)
			http.Error(w, "Error, Saat mendapatkan data dari database", 500)
		}
		defer rows.Close()

		// looping untuk mendapatkan semua data
		for rows.Next() {
			data := new(productData)

			//Read the columns in each row into variables
			err := rows.Scan(&data.Store_id, &data.Id, &data.Name, &data.Harga, &data.Url, &data.Image, &data.Update_date)
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
}

//Mendapatkan product By Id dan store Id
func GetProduct(w http.ResponseWriter, r *http.Request) {

	datas := []*productData{}

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
	product_id, store, strStore := "", "", ""
	bolError := false

	// mendapatkan parameter
	params := mux.Vars(r)
	if val, ok := params["product_id"]; ok {
		product_id = val
	}

	if val, ok := params["store"]; ok {
		store = val
		strStore = fmt.Sprintf("and store_id ='%s' ", store)
	}

	if product_id == "" && store == "" {
		// cek apakah keywordnya ada
		http.Error(w, "Error, mohon masukkan kode product dan store", 500)
		bolError = true
	} else if store != "" {
		// cek id store apakah ada di list
		cStore := 0
		strQuery := "select count(*) as cStore from store where id = ? "
		err = db.QueryRow(strQuery, store).Scan(&cStore)
		if err != nil {
			log.Panic(err)
			http.Error(w, "Error, Saat mendapatkan data dari database", 500)
			bolError = true
		} else if cStore == 0 {
			http.Error(w, "Error, Store ("+store+") tidak ada", 500)
			bolError = true
		}
	}

	if !bolError {
		// mendapatkan data dari database

		strQuery := fmt.Sprintf("select store_id,id,name,harga,url,Image,update_date from product "+
			" where id= '%s'  %s", product_id, strStore)
		rows, err := db.Query(strQuery)
		if err != nil {
			log.Panic(err)
			http.Error(w, "Error, Saat mendapatkan data dari database", 500)
		}
		defer rows.Close()

		// looping untuk mendapatkan semua data
		for rows.Next() {
			data := new(productData)

			//Read the columns in each row into variables
			err := rows.Scan(&data.Store_id, &data.Id, &data.Name, &data.Harga, &data.Url, &data.Image, &data.Update_date)
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
}
