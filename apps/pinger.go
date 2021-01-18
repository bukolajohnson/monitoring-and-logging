package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	hello = map[string]string{
		"en": "Hello",
		"es": "Hola",
		"de": "Hallo",
		"ch": "你好",
		"ru": "Привет",
	}

	port = os.Getenv("PORT")
)

func handler(w http.ResponseWriter, r *http.Request) {
	lang := strings.TrimPrefix(r.URL.RequestURI(), "/")
	greeting, ok := hello[lang]
	if !ok {
		fmt.Printf("unknown language: %s\n", lang)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, greeting)
}

func selfPing() {
	start := time.Now()
	oscillationFactor := func() float64 {
		return 20 + math.Sin(math.Sin(2*math.Pi*float64(time.Since(start))/float64(10*time.Minute)))
	}

	langs := []string{"en", "es", "de", "ch", "ru"}

	for {
		time.Sleep(time.Duration(100*oscillationFactor()) * time.Millisecond)

		index := rand.Intn(len(langs))
		lang := langs[index]

		resp, err := http.Get(fmt.Sprintf("http://0.0.0.0:%s/%s", port, lang))
		if err != nil {
			log.Println("got error", err)
			continue
		}

		io.Copy(os.Stdout, resp.Body)
		fmt.Println()
		resp.Body.Close()
	}
}

func main() {
	rand.Seed(time.Now().Unix())

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Println("starting server on", addr)

	go selfPing()

	log.Fatal(http.ListenAndServe(addr, mux))
}
