package record

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

var (
	client *http.Client = &http.Client{Timeout: 3 * time.Minute}
)

func init() {
	// load our environmental variables.
	if err := godotenv.Load("../../.env"); err != nil {
		panic(err)
	}
}

// TODO: call main server in test file
// TODO: get token

// TODO: create test migration with data
// TODO: try validator package for response validation

// TestRecordEndpoints is a function that tests the record module endpoints
func TestRecordEndpoints(t *testing.T) {
	baseURL := fmt.Sprintf("%s:%s", os.Getenv("API_URL_REST"), os.Getenv("API_URL_REST_PORT"))
	t.Run("Generate Token", func(t *testing.T) {
		// create request body // TODO: opt to mock data
		body := map[string]interface{}{
			"email": "test@example.com",
		}
		jsonBody, _ := json.Marshal(body)

		// create request
		req, err := http.NewRequest("POST", baseURL+"/v1/record/token/generate", bytes.NewBuffer(jsonBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// send the request
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// parse response
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		// status code should be 200
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// access token should not be empty
		assert.NotEmpty(t, response["data"].(map[string]interface{})["accessToken"])
	})

	t.Run("Create Record", func(t *testing.T) {
		// first generate token
		tokenReq, err := http.NewRequest("POST", baseURL+"/v1/record/token/generate", bytes.NewBuffer([]byte(`{"email":"test@example.com"}`)))
		assert.NoError(t, err)
		tokenReq.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(tokenReq)
		assert.NoError(t, err)
		defer resp.Body.Close()

		var tokenResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&tokenResponse)
		token := tokenResponse["data"].(map[string]interface{})["accessToken"].(string)

		// TODO: should be struct import from dto controllers
		// create record request
		request := map[string]interface{}{
			"id":   generateID(),
			"data": "Test Description",
		}

		buf := new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(request)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", baseURL+"/v1/record", buf)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err = client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		// status code should be 201
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// check response fields
		assert.NotEmpty(t, response["data"].(map[string]interface{})["id"])
		assert.NotEmpty(t, response["data"].(map[string]interface{})["data"])
		assert.NotEmpty(t, response["data"].(map[string]interface{})["createdAt"])
	})

	t.Run("Get Records", func(t *testing.T) {
		// first generate token
		tokenReq, err := http.NewRequest("POST", baseURL+"/v1/record/token/generate", bytes.NewBuffer([]byte(`{"email":"test@example.com"}`)))
		assert.NoError(t, err)
		tokenReq.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(tokenReq)
		assert.NoError(t, err)
		defer resp.Body.Close()

		var tokenResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&tokenResponse)
		token := tokenResponse["data"].(map[string]interface{})["accessToken"].(string)

		// pagination
		page := 1

		// get records request
		req, err := http.NewRequest("GET", baseURL+"/v1/record/list", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		q := req.URL.Query()
		q.Add("page", strconv.Itoa(page))
		req.URL.RawQuery = q.Encode()

		resp, err = client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		// status code should be 200
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// check fields
		assert.NotEmpty(t, response["data"].(map[string]interface{})["records"])
		assert.NotEmpty(t, response["data"].(map[string]interface{})["total"])
	})

	t.Run("Get Record By ID", func(t *testing.T) {
		// first generate token
		tokenReq, err := http.NewRequest("POST", baseURL+"/v1/record/token/generate", bytes.NewBuffer([]byte(`{"email":"test@example.com"}`)))
		assert.NoError(t, err)
		tokenReq.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(tokenReq)
		assert.NoError(t, err)
		defer resp.Body.Close()

		var tokenResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&tokenResponse)
		token := tokenResponse["data"].(map[string]interface{})["accessToken"].(string)

		// id of the test record to be fetched
		ID := "4"

		// get record request
		req, err := http.NewRequest("GET", baseURL+"/v1/record/"+ID, nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err = client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// parse response
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		// status code should be 200
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// check fields
		assert.NotEmpty(t, response["data"].(map[string]interface{})["id"])
		assert.NotEmpty(t, response["data"].(map[string]interface{})["data"])
		assert.NotEmpty(t, response["data"].(map[string]interface{})["createdAt"])
	})

	t.Run("Update Record", func(t *testing.T) {
		// first generate token
		tokenReq, err := http.NewRequest("POST", baseURL+"/v1/record/token/generate", bytes.NewBuffer([]byte(`{"email":"test@example.com"}`)))
		assert.NoError(t, err)
		tokenReq.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(tokenReq)
		assert.NoError(t, err)
		defer resp.Body.Close()

		var tokenResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&tokenResponse)
		token := tokenResponse["data"].(map[string]interface{})["accessToken"].(string)

		// update record request
		request := map[string]interface{}{
			"data": "Updated Description",
		}

		buf := new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(request)
		assert.NoError(t, err)

		// id of the test record to be updated
		ID := "4"

		req, err := http.NewRequest("PUT", baseURL+"/v1/record/"+ID+"/update", buf)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err = client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// status code should be 200
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Delete Record", func(t *testing.T) {
		// first generate token
		tokenReq, err := http.NewRequest("POST", baseURL+"/v1/record/token/generate", bytes.NewBuffer([]byte(`{"email":"test@example.com"}`)))
		assert.NoError(t, err)
		tokenReq.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(tokenReq)
		assert.NoError(t, err)
		defer resp.Body.Close()

		var tokenResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&tokenResponse)
		token := tokenResponse["data"].(map[string]interface{})["accessToken"].(string)

		// id to be deleted // TODO: opt to mock data
		ID := "2tTG0HRIIyKBhCJCZIyYMKDIUwG"

		// delete record request
		req, err := http.NewRequest("DELETE", baseURL+"/v1/record/"+ID+"/delete", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err = client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// status code should be 200
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// generateID generates unique id
func generateID() string {
	return ksuid.New().String()
}
