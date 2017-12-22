# go-price-minimarket

### Prerequisites

* MySQL, Proses instalasi menggunakan perintah go get -u github.com/go-sql-driver/mysql
*	Mux (untuk Routing), Proses instalasi menggunakan perintah go get -u github.com/gorilla/mux
*	Go-Colly, Proses instalasi menggunakan perintah go get -u github.com/gocolly/colly/

### Web Scrapper 
Folder (price-scraper)

Cara Instalasi aplikasi 
```sh
cd price-scraper
go install price-scraper
```

Cara Menjalankan aplikasi 
```sh
./price-scrapper [idm/alf]
```

idm = Indomaret, alf = alfamart


### Web Service
Folder (price-ws)

Cara Instalasi aplikasi 
```sh
cd price-ws
go install price-ws
```

Cara Menjalankan aplikasi 
* Daftar Toko
```sh
curl http://localhost:8080/store/limit/{limit}/start/{start}
```

* Daftar Produk
```sh
curl http://localhost:8080/product/store/{store}/limit/{limit}/start/{start}
```


* Pencarian Produk
```sh
curl http://localhost:8080/product/key/{key}/store/{store}/limit/{limit}/start/{start}
```


* Harga Produk
```sh
curl http://localhost:8080/price/product_name/{product_name}/store/{store}
```
atau
```sh
curl http://localhost:8080/price/product_id/{product_id}/store/{store}
```

* Histori Harga Produk
```sh
curl http://localhost:8080/price/his/product_name/{product_name}/store/{store}/limit/{limit}
```
atau
```sh
curl http://localhost:8080/price/his/product_id/{product_id}/store/{store}/limit/{limit}
```

* Komparasi Harga Produk
```sh
curl http://localhost:8080/compare/product_id/{product_id_1}/{product_id_2}/store/{store_1}/{store_2}
```
atau
```sh
curl http://localhost:8080/compare/product_name/{product_name}/{product_id_2}/store/{store_1}/{store_2}
```
