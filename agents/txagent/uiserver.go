package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

func renderIndex(resWri http.ResponseWriter, requ *http.Request) {
	resWri.Write([]byte("Data diode TX Agent"))
}

// func renderBuffer(resWri http.ResponseWriter, requ *http.Request) {
// 	// buffer.Lock()
// 	// defer buffer.Unlock()
// 	// buffBytes, marErr := yaml.Marshal(buffer.content)
// 	if marErr == nil {
// 		resWri.Write(buffBytes)
// 	} else {
// 		resWri.Write([]byte(marErr.Error()))
// 	}
// }

func renderConfig(resWri http.ResponseWriter, requ *http.Request) {
	confBytes, marErr := yaml.Marshal(configuration)
	if marErr == nil {
		resWri.Write(confBytes)
	} else {
		resWri.Write([]byte(marErr.Error()))
	}
}

func serveStatus() {
	var router = mux.NewRouter().StrictSlash(true)
	netInterface := "127.0.0.1:8686"
	router.HandleFunc("/", renderIndex)
	// router.HandleFunc("/buffer", renderBuffer)
	router.HandleFunc("/config", renderConfig)
	fmt.Println("Running status http server at:", netInterface)
	http.ListenAndServe(netInterface, router)
}
