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

// @title restaurant API
// @version 1.0
// @description This is a demo http server to serve rest apis.
// @termsOfService http://swagger.io/terms/

// @contact.name Ved Prakash Verma
// @contact.url http://www.swagger.io/support
// @contact.email vedprakashv888@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes https

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

