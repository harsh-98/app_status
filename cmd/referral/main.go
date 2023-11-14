package main

import (
	"fmt"
	"net/http"

	"github.com/Gearbox-protocol/sdk-go/log"
	"github.com/Gearbox-protocol/sdk-go/utils"
	"github.com/gorilla/mux"
)

var PORT = 8080

func middleware(fn func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS,PUT")
		// w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		fn(w, r)
	}
}
func writeErr(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(utils.ToJsonBytes(map[string]string{"message": err.Error()}))
}
func writeSucess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(utils.ToJsonBytes(map[string]interface{}{"data": data}))
}

func generate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	writeSucess(w)
}

func main() {
	router := mux.NewRouter() //
	mux.HandleFunc("/generate", middleware(generate))

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", PORT),
		Handler: router,
	}
	//
	log.Infof("Starting web server at :%d", PORT)
	srv.ListenAndServe()
}
