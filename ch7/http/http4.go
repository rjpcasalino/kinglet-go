package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/text/currency"
)

func main() {
	db := database{"shoes": 50, "socks": 5}
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/currency", db.usd)
	log.Fatal(http.ListenAndServe("[::1]:8000", nil))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price, ok := db[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	fmt.Fprintf(w, "%s\n", price)
}

func (db database) usd(w http.ResponseWriter, req *http.Request) {
	t1799, _ := time.Parse("2006-01-02", "1799-01-01")
	for it := currency.Query(currency.Date(t1799)); it.Next(); {
		from := ""
		if t, ok := it.From(); ok {
			from = t.Format("2006-01-02")
		}
		fmt.Fprintf(w, "%v is used in %v since: %v\n", it.Unit(), it.Region(), from)
	}
}
