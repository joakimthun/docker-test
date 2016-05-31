package main

import ( 
	"net/http"
	"os"
    "log"
    _ "github.com/lib/pq"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    _ "github.com/joakimthun/docker-test/redis"
    _ "github.com/joakimthun/docker-test/endpoints"
)

func main() {
    port := os.Getenv("PORT")
    
	if port == "" {
		port = "8080"
	}
	
    log.Println("Listening on port: " + port)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}