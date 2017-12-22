/*
Author : Miswar
*/

package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func crawlerAlfamart(storeId string) {

	// mendapatkan koneksi database
	db, err := database()
	if err != nil {
		log.Panic(err)
	}

	// inisialisasi crawler
	c := colly.NewCollector()
	// set timeout 60 detik
	c.SetRequestTimeout(60 * time.Second)
	// set delay 5 detik
	c.Limit(&colly.LimitRule{
		Delay: 5 * time.Second,
	})

	// inisialisasi variabel
	i := 1
	domain := "https://www.alfacart.com"
	link := domain + "/seller/alfamart-5659"
	bodyFormTemplate := "page=[page]&tampilkan=20&tab=tab1a&categoryValue=&brandValue=&minPriceHide=800&minPrice=&maxPriceHide=283500&maxPrice=&ratingValue=&score=5&score=4&score=3&score=2&score=1&totalOthers=1&codeOthers0=aol_size&othersValue0=&sorting=desc%2Crelevancy&gridList=&pagereview=1&tampilkanreview=4&score=4"

	// inisialisasi header agar menerima form post
	hdr := http.Header{}
	hdr.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	// inisialisasi untuk halaman ke 1
	page := 1
	bodyForm := strings.Replace(bodyFormTemplate, "[page]", strconv.Itoa(page), -1)

	// get url product
	c.OnHTML("div[id=products] > div.item > div.thumbnail ", func(e *colly.HTMLElement) {

		// mendapatkan data nama, id, price, url image dan url product
		name := e.ChildText(".caption > h4 > a[href]")
		id := e.ChildAttr(".caption > input", "value")
		price := e.ChildText(".vdivide > .col-xs-8 > .lead")
		urlImage := e.ChildAttr("a > img", "data-src")
		urlProduct := e.ChildAttr(".caption > h4 > a[href]", "href")

		// menghilangkan tulisan Rp dan .
		price = strings.Replace(price, "Rp", "", -1)
		price = strings.Replace(price, ".", "", -1)

		//convert string to int
		strconv.Atoi(price)

		// jika nama mengandung ..., maka sesuaikan nama
		if strings.Contains(name, "...") == true {
			name = strings.Replace(urlProduct, "/product/", "", -1)
			name = strings.Replace(name, "-", " ", -1)
			name = strings.Replace(name, id, " ", -1)
		}

		// add domain in urlProduct
		if urlProduct != "" {
			urlProduct = domain + urlProduct
		}
		fmt.Printf("\n %d  , %s , %s ", i, name, price)

		// memulai menyimpan di database
		tx, err := db.Begin()

		// mempersiapkan query
		stmt, err := tx.Prepare("REPLACE INTO product (store_id,id,name,harga,url,image,update_date)VALUES(?,?,?,?,?,?,now());")
		stmt_his, err := tx.Prepare("REPLACE INTO product_his (store_id,id,name,harga,url,image,update_date)VALUES(?,?,?,?,?,?,now());")

		if err != nil {
			log.Fatal(err)
		}

		//execute dengan parameter data
		_, err = stmt.Exec(storeId, id, name, price, urlProduct, urlImage)
		_, err = stmt_his.Exec(storeId, id, name, price, urlProduct, urlImage)
		if err != nil {
			log.Panic(err)
		}

		//commit eksekusi
		tx.Commit()

		i += 1

	})

	c.OnHTML("ul[id=paging] > li ", func(e *colly.HTMLElement) {
		//mencari halaman berikutnya
		if e.ChildText("a[href]") == "Berikutnya" {
			page += 1
			bodyForm := strings.Replace(bodyFormTemplate, "[page]", strconv.Itoa(page), -1)

			//memulai proses halaman berikutnya
			c.Request("POST", link, strings.NewReader(bodyForm), nil, nil)
		}
	})

	//memulai proses halaman pertama
	fmt.Println("\nProses Crawler Alfamart Dimulai")
	c.Request("POST", link, strings.NewReader(bodyForm), nil, nil)

	// proses crawler telah selesai
	fmt.Println("\nProses Crawler Alfamart Selesai")
}
