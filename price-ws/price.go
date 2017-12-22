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
	"strings"
)

//Mendapatkan harga
func GetPriceProduct(w http.ResponseWriter, r *http.Request) {

	datas := []*priceData{}

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
	product_name, product_id, store, strStore := "", "", "", ""
	bolError := false

	// mendapatkan parameter
	params := mux.Vars(r)
	if val, ok := params["product_name"]; ok {
		product_name = strings.Replace(val, "+", " ", -1)
	}
	if val, ok := params["product_id"]; ok {
		product_id = val
	}
	if val, ok := params["store"]; ok {
		store = val
		strStore = fmt.Sprintf("and store_id ='%s' ", store)
	}

	//cek kelengkapan parameter
	if product_name == "" && product_id == "" {
		// cek apakah product name dan product id ada
		http.Error(w, "Error, mohon masukkan Nama Produk atau ID Produk", 500)
		bolError = true
	} else if store == "" {
		// cek id store apakah ada di list
		http.Error(w, "Error, mohon masukkan kode store ", 500)
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
		strQuery := ""
		if product_name != "" {
			strQuery = fmt.Sprintf("select name,harga,update_date from product "+
				" where MATCH(name) AGAINST('%s') %s ORDER BY MATCH(name) AGAINST('%s') DESC  LIMIT 1", product_name, strStore, product_name)
		} else {
			strQuery = fmt.Sprintf("select name,harga,update_date from product "+
				" where id ='%s' %s LIMIT 1", product_id, strStore)
		}

		rows, err := db.Query(strQuery)
		if err != nil {
			log.Panic(err)
			http.Error(w, "Error, Saat mendapatkan data dari database", 500)
		}
		defer rows.Close()

		// looping untuk mendapatkan semua data
		for rows.Next() {
			data := new(priceData)

			//Read the columns in each row into variables
			err := rows.Scan(&data.Name, &data.Harga, &data.Update_date)
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

//Mendapatkan harga
func GetPriceProductHis(w http.ResponseWriter, r *http.Request) {

	datas := []*priceData{}

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
	limit, product_name, product_id, store, strStore := consLimit, "", "", "", ""
	strLimit := fmt.Sprintf("LIMIT %s", limit)
	bolError := false

	// mendapatkan parameter
	params := mux.Vars(r)
	if val, ok := params["limit"]; ok {
		limit = val
		strLimit = fmt.Sprintf("LIMIT %s", limit)
	}
	if val, ok := params["product_name"]; ok {
		product_name = strings.Replace(val, "+", " ", -1)
	}
	if val, ok := params["product_id"]; ok {
		product_id = val
	}
	if val, ok := params["store"]; ok {
		store = val
		strStore = fmt.Sprintf("and store_id ='%s' ", store)
	}

	//cek kelengkapan parameter
	if product_name == "" && product_id == "" {
		// cek apakah product name dan product id ada
		http.Error(w, "Error, mohon masukkan Nama Produk atau ID Produk", 500)
		bolError = true
	} else if store == "" {
		// cek id store apakah ada di list
		http.Error(w, "Error, mohon masukkan kode store ", 500)
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
	} else if product_name != "" {

		// mencari id produk
		strQueryGet := fmt.Sprintf("select id from product "+
			" where MATCH(name) AGAINST('%s') %s ORDER BY MATCH(name) AGAINST('%s') DESC  LIMIT 1", product_name, strStore, product_name)

		err := db.QueryRow(strQueryGet).Scan(&product_id)
		if err != nil {
			log.Panic(err)
			http.Error(w, "Error, Saat mendapatkan data dari database", 500)
			bolError = true
		} else if product_id == "" {
			http.Error(w, "Error, Data tidak ada", 404)
			bolError = true
		}
	}

	if !bolError {
		// mendapatkan data dari database

		strQuery := fmt.Sprintf("select name,harga,update_date from product_his "+
			" where id ='%s' %s order by update_date desc %s", product_id, strStore, strLimit)

		rows, err := db.Query(strQuery)
		if err != nil {
			log.Panic(err)
			http.Error(w, "Error, Saat mendapatkan data dari database", 500)
		}
		defer rows.Close()

		// looping untuk mendapatkan semua data
		for rows.Next() {
			data := new(priceData)

			//Read the columns in each row into variables
			err := rows.Scan(&data.Name, &data.Harga, &data.Update_date)
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
