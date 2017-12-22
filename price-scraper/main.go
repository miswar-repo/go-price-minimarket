/*
Author : Miswar
*/

package main

import (
	"fmt"
	//"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
)

// silahkan sesuaikan configurasi koneksi ke database dengan konfigurasi environment anda
const (
	host     = "192.168.73.100"
	port     = "3306"
	user     = "root"
	password = "password"
	dbname   = "pricedb"
)

func main() {

	//inisialisasi list store
	storeList := []string{"Indomaret", "Alfamart"}
	//inisialisasi variabel untuk menyimpan store yang dipilih
	storeSelected := ""

	// Parameter hanya terdiri dari satu atau tidak sama sekali (menjalankan semua store)
	if len(os.Args) > 2 && len(os.Args) < 1 {
		fmt.Printf("Sorry, Parameternya hanya satu ( %s ). price-scrapper [nama store] \n", strings.Join(storeList, ", "))
		return
	} else if len(os.Args) == 2 {
		// Menyimpan argument ke variabel storeSelected
		storeSelected = strings.ToLower(os.Args[1])

		if storeSelected == "indomaret" {
			//run crawler Indomaret
			storeIdIndo := "idm"
			crawlerIndomaret(storeIdIndo)
		} else if storeSelected == "alfamart" {
			//run crawler Alfamart
			storeIdAlfa := "alf"
			crawlerAlfamart(storeIdAlfa)
		} else {
			// Parameter Store parameter yang didukung
			fmt.Printf("Sorry, Parameter ( %s ) tidak tersedia. Parameter yang tersedia ( %s ) \n", storeSelected, strings.Join(storeList, ", "))
			return
		}
	} else {
		//run crawler alfamart & Indomaret
		storeIdAlfa := "alf"
		crawlerAlfamart(storeIdAlfa)

		storeIdIndo := "idm"
		crawlerIndomaret(storeIdIndo)
	}

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
