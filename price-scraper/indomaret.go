/*
Author : Miswar
*/

package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
	"time"
)

func crawlerIndomaret(storeId string) {

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

	//clone collector
	cSetRegion := c.Clone()
	cGetPrice := c.Clone()
	cGetCategory := c.Clone()

	// inisialisasi variabel
	domain := "http://www.klikindomaret.com"
	kota := "Bandung"
	i := 1

	// akses domain setelah melakukan set lokasi
	cSetRegion.OnResponse(func(r *colly.Response) {
		cGetCategory.Visit(domain)
	})

	// Check lokasi
	cGetCategory.OnHTML("a.setlocation > b.geolocation.block", func(e *colly.HTMLElement) {
		fmt.Printf("Location : %s \n", e.Text)
	})

	cGetCategory.OnHTML("div.megamenu > ul > li", func(e *colly.HTMLElement) {
		// mendapatkan kategori
		catPage := e.ChildAttr("a", "href")
		// akses ke halaman kategori tersebut
		cGetPrice.Visit(domain + catPage)
	})

	// Mendapatkan data product
	cGetPrice.OnHTML("ul.list-unstyled > li", func(e *colly.HTMLElement) {
		name := strings.TrimSpace(e.ChildText("a > h5.producttitle"))
		price := e.ChildText("a > span.price > b")
		urlImage := e.ChildAttr("a > span.thumbnail  > img", "data-original")
		urlProduct := e.ChildAttr("a.comparable", "href")
		typeShipping := strings.TrimSpace(e.ChildText("a > small.shipping"))
		id := e.ChildAttr("input[id=itemId]", "value")

		// add domain in urlProduct
		if urlProduct != "" {
			urlProduct = domain + urlProduct
		}

		// menghilangkan tulisan Rp dan .
		price = strings.Replace(price, "Rp", "", -1)
		price = strings.Replace(price, ".", "", -1)

		//convert string to int
		strconv.Atoi(price)

		if name != "" && strings.ToLower(typeShipping) == "dikirim dari toko" {

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
		}

	})

	cGetPrice.OnHTML(".bm-25 > div.right-floated", func(e *colly.HTMLElement) {
		// mendapatkan url halaman berikutnya
		nextPage := e.ChildAttr("button.nextpage", "data-href")
		// akses ke halaman berikutnya
		cGetPrice.Visit(nextPage)

	})

	// Memulai Prosess Crawler
	fmt.Println("\nProcess Crawler Indomaret Dimulai")

	cSetRegion.Visit(domain + "/home/postregion?regionName=" + kota)

	fmt.Println("\nProcess Crawler Indomaret Selesai")
}
