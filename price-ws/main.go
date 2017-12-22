/*
Author : Miswar
*/

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	// silahkan sesuaikan configurasi koneksi ke database dengan konfigurasi environment anda
	host     = "192.168.73.100"
	port     = "3306"
	user     = "root"
	password = "password"
	dbname   = "pricedb"

	consLimit = "25"
	consStart = "0"
)

//Struct Type
//Store
type storeData struct {
	Id   string
	Name string
}

//Product
type productData struct {
	Store_id    string
	Id          string
	Name        string
	Harga       int
	Url         string
	Image       string
	Update_date string
}

//Product
type priceData struct {
	Name        string
	Harga       int
	Update_date string
}

//main function
func main() {
	//inisialisasi port
	port := 8080

	// inisialisasi router
	router := mux.NewRouter()

	// configurasi routing

	// store
	router.HandleFunc("/store", GetStores).Methods("GET")
	router.HandleFunc("/store/", GetStores).Methods("GET")
	router.HandleFunc("/store/limit/{limit}", GetStores).Methods("GET")
	router.HandleFunc("/store/limit/{limit}/start/{start}", GetStores).Methods("GET")

	// daftar produk
	router.HandleFunc("/product", GetProducts).Methods("GET")
	router.HandleFunc("/product/", GetProducts).Methods("GET")
	router.HandleFunc("/product/limit/{limit}", GetProducts).Methods("GET")
	router.HandleFunc("/product/limit/{limit}/start/{start}", GetProducts).Methods("GET")

	router.HandleFunc("/product/store/{store}", GetProducts).Methods("GET")
	router.HandleFunc("/product/store/{store}/limit/{limit}", GetProducts).Methods("GET")
	router.HandleFunc("/product/store/{store}/limit/{limit}/start/{start}", GetProducts).Methods("GET")

	// produk berdasarkan id
	router.HandleFunc("/product/product_id/{product_id}/store/{store}", GetProduct).Methods("GET")

	// produk berdasarakan Keyword
	router.HandleFunc("/product/key/{key}", GetProductsByKey).Methods("GET")
	router.HandleFunc("/product/key/{key}/limit/{limit}", GetProductsByKey).Methods("GET")
	router.HandleFunc("/product/key/{key}/limit/{limit}/start/{start}", GetProductsByKey).Methods("GET")

	router.HandleFunc("/product/key/{key}/store/{store}", GetProductsByKey).Methods("GET")
	router.HandleFunc("/product/key/{key}/store/{store}/limit/{limit}", GetProductsByKey).Methods("GET")
	router.HandleFunc("/product/key/{key}/store/{store}/limit/{limit}/start/{start}", GetProductsByKey).Methods("GET")

	// harga/price
	router.HandleFunc("/price/product_id/{product_id}/store/{store}", GetPriceProduct).Methods("GET")
	router.HandleFunc("/price/product_name/{product_name}/store/{store}", GetPriceProduct).Methods("GET")

	// history harga
	router.HandleFunc("/price/his/product_id/{product_id}/store/{store}", GetPriceProductHis).Methods("GET")
	router.HandleFunc("/price/his/product_name/{product_name}/store/{store}", GetPriceProductHis).Methods("GET")
	router.HandleFunc("/price/his/product_id/{product_id}/store/{store}/limit/{limit}", GetPriceProductHis).Methods("GET")
	router.HandleFunc("/price/his/product_name/{product_name}/store/{store}/limit/{limit}", GetPriceProductHis).Methods("GET")

	// komparasi harga
	router.HandleFunc("/compare/product_id/{product_id_1}/{product_id_2}/store/{store_1}/{store_2}", GetComparePrice).Methods("GET")
	router.HandleFunc("/compare/product_name/{product_name}/store/{store_1}/{store_2}", GetComparePrice).Methods("GET")

	// Start Server
	log.Printf("Server Webservice Price starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), router))
}

//Koneksi Database
func database() (*sql.DB, error) {

	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)
	db, err := sql.Open("mysql", connString)

	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
