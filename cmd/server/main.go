package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/souravlayek/file-hosting/internal/api/routes"
	"github.com/souravlayek/file-hosting/internal/database"
	"github.com/souravlayek/file-hosting/utils"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 || args[0] != "server" {
		fmt.Println("Loading ENV")
		utils.LoadENV()
	}

	database.ConnectDB()
	router := routes.NewRouter()
	fmt.Println("Server started at port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
