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

// Compare Harga 2 Store
func GetComparePrice(w http.ResponseWriter, r *http.Request) {

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
	product_name, product_id_1, product_id_2, store_1, store_2 := "", "", "", "", ""
	store := make(map[string]string)

	bolError := false

	// mendapatkan parameter
	params := mux.Vars(r)
	if val, ok := params["product_name"]; ok {
		product_name = strings.Replace(val, "+", " ", -1)
	}

	if val, ok := params["product_id_1"]; ok {
		product_id_1 = val
	}
	if val, ok := params["product_id_2"]; ok {
		product_id_2 = val
	}

	if val, ok := params["store_1"]; ok {
		store_1 = val
	}

	if val, ok := params["store_2"]; ok {
		store_2 = val
	}

	//cek kelengkapan parameter
	if product_name == "" && product_id_1 == "" && product_id_2 == "" {
		// cek apakah product name dan product id ada
		http.Error(w, "Error, mohon masukkan Nama Produk atau ID Produk yang akan di komparasi", 500)
		bolError = true
	} else if product_name == "" && (product_id_1 == "" || product_id_2 == "") {
		// cek apakah ada 2  product id
		http.Error(w, "Error, mohon masukkan 2 ID Produk yang akan di komparasi", 500)
		bolError = true
	} else if store_1 == "" || store_2 == "" {
		// cek apakah ada 2  id store
		http.Error(w, "Error, mohon masukkan 2 kode store yang akan di komparasi", 500)
		bolError = true
	}

	if store_1 != "" && !bolError {
		// cek id store apakah ada di list
		cStore := 0
		strQuery := "select count(*) as cStore from store where id = ? "
		err = db.QueryRow(strQuery, store_1).Scan(&cStore)
		if err != nil {
			log.Panic(err)
			http.Error(w, "Error, Saat mendapatkan data dari database", 500)
			bolError = true
		} else if cStore == 0 {
			http.Error(w, "Error, Store ke-1 ("+store_1+") tidak ada", 500)
			bolError = true
		} else {
			store[store_1] = product_id_1
		}
	}

	if store_2 != "" && !bolError {
		// cek id store apakah ada di list
		cStore := 0
		strQuery := "select count(*) as cStore from store where id = ? "
		err = db.QueryRow(strQuery, store_2).Scan(&cStore)
		if err != nil {
			log.Panic(err)
			http.Error(w, "Error, Saat mendapatkan data dari database", 500)
			bolError = true
		} else if cStore == 0 {
			http.Error(w, "Error, Store ke-2 ("+store_2+") tidak ada", 500)
			bolError = true
		} else {
			store[store_2] = product_id_2
		}
	}

	if product_name != "" && !bolError {
		// mencari id produk
		i := 0
		for key, value := range store {
			strQueryGet := fmt.Sprintf("select id from product "+
				" where MATCH(name) AGAINST('%s') ORDER BY MATCH(name) AGAINST('%s') "+
				" and store_id='%s' DESC  LIMIT 1", product_name, product_name, key)

			err := db.QueryRow(strQueryGet).Scan(&value)

			if err != nil {
				log.Panic(err)
				http.Error(w, "Error, Saat mendapatkan data dari database", 500)
				bolError = true
			} else if value != "" {
				i = i + 1
				store[key] = value
			}
		}

		if i == 0 {
			http.Error(w, "Error, Data tidak ada untuk nama "+product_name, 404)
			bolError = true
		}

	}

	if !bolError {
		// mendapatkan data dari database
		for key, value := range store {

			if value != "" {
				data := new(productData)
				strQuery := fmt.Sprintf("select store_id,id,name,harga,url,Image,update_date from product "+
					" where id ='%s' and store_id ='%s' ", value, key)

				err := db.QueryRow(strQuery).Scan(&data.Store_id, &data.Id, &data.Name, &data.Harga, &data.Url, &data.Image, &data.Update_date)
				if err != nil {
					log.Panic(err)
				} else {
					datas = append(datas, data)
				}
			}
		}

		//cek apakah ada data yang didapat
		if len(datas) == 0 {
			http.Error(w, "Error, Data tidak ada", 404)
		} else {
			json.NewEncoder(w).Encode(datas)
		}
	}
}
