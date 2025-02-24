package init

import (
	"log"
	"os"
	"strconv"
	"testing"

	"gomora/interfaces/http/rest"

	"github.com/joho/godotenv"
)

// FIXME: won't connect to the actual port
// InitTest is a function that initializes the test environment
func InitTest(m *testing.M) {
	// load our environmental variables.
	if err := godotenv.Load("../../.env"); err != nil {
		panic(err)
	}

	// get the rest port
	restPort, err := strconv.Atoi(os.Getenv("API_URL_REST_PORT"))
	if err != nil {
		log.Fatalf("[SERVER] Invalid port")
	}

	// start the server
	rest.ChiRouter().Serve(restPort)

	// run the tests
	code := m.Run()

	log.Printf("[TEST] connection success")

	// closes the testing environment after the tests are finished
	os.Exit(code)
}
