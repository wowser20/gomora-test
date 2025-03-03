package main

import (
	"gomora/interfaces/http/rest"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	// load our environmental variables.
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

// FIXME: won't connect to the actual port
// InitTest is a function that initializes the test environment
func InitTest(m *testing.M) {
	// get the rest port
	restPort, err := strconv.Atoi(os.Getenv("API_URL_REST_PORT"))
	if err != nil {
		log.Fatalf("[SERVER] Invalid port")
	}

	// start the server
	rest.ChiRouter().Serve(restPort)
}
