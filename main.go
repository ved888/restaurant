package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"restaurant/common"
	"restaurant/operation"
)

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.ErrorLevel)

	common.DbConnection()

	r := mux.NewRouter()

	operation.Handler(r)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("error in portNumber line", err)
		return
	}
	fmt.Println("start")
}
