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

	"github.com/brianvoe/gofakeit/v6"
	"github.com/joho/godotenv"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"

	"gomora/interfaces/http/rest/viewmodels"
	types "gomora/module/record/interfaces/http"
)

var (
	client *http.Client = &http.Client{Timeout: 3 * time.Minute}
)

var (
	token  string
	testID string
)

// TODO: create test migration with data
// TODO: try validator package for response validation

func init() {
	// load our environmental variables.
	if err := godotenv.Load("../../.env"); err != nil {
		panic(err)
	}
}

// TestRecordEndpoints is a function that tests the record module endpoints
func TestRecordEndpoints(t *testing.T) {
	baseURL := fmt.Sprintf("%s:%s", os.Getenv("API_URL_REST"), os.Getenv("API_URL_REST_PORT"))
	t.Run("Generate Token", func(t *testing.T) {
		body := map[string]interface{}{
			"email": gofakeit.Email(),
		}

		req := createHTTPRequest(t, "POST", baseURL+"/v1/record/token/generate", nil, body)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		var response viewmodels.HTTPResponseVM
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		/// success should be true
		assert.Equal(t, response.Success, true)

		expectedData := types.GenerateTokenResponse{}

		jsonData, err := json.Marshal(response.Data)
		assert.NoError(t, err)

		actualData := types.GenerateTokenResponse{}
		err = json.Unmarshal(jsonData, &actualData)
		assert.NoError(t, err)

		/// fetch access token
		token = actualData.AccessToken

		/// expected response struct type should match the actual response struct type
		assert.IsType(t, expectedData, actualData)

		/// actual response should not be empty
		assert.NotEmpty(t, actualData.AccessToken)
	})

	t.Run("Create Record", func(t *testing.T) {
		request := types.CreateRecordRequest{
			ID:   generateID(),
			Data: gofakeit.Cat(),
		}

		req := createHTTPRequest(t, "POST", baseURL+"/v1/record", &token, request)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		var response viewmodels.HTTPResponseVM
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		/// success should be true
		assert.Equal(t, response.Success, true)

		expectedData := types.CreateRecordResponse{}

		jsonData, err := json.Marshal(response.Data)
		assert.NoError(t, err)

		actualData := types.CreateRecordResponse{}
		err = json.Unmarshal(jsonData, &actualData)
		assert.NoError(t, err)

		/// fetch id
		testID = actualData.ID

		/// expected response struct type should match the actual response struct type
		assert.IsType(t, expectedData, actualData)

		/// actual response should not be empty
		assert.NotEmpty(t, actualData.ID)
		assert.NotEmpty(t, actualData.Data)
		assert.NotEmpty(t, actualData.CreatedAt)
	})

	t.Run("Get Records", func(t *testing.T) {
		page := 1

		req := createHTTPRequest(t, "GET", baseURL+"/v1/record/list", &token, nil)

		q := req.URL.Query()
		q.Add("page", strconv.Itoa(page))
		req.URL.RawQuery = q.Encode()

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		var response viewmodels.HTTPResponseVM
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		/// success should be true
		assert.Equal(t, response.Success, true)

		expectedData := types.GetRecordResponse{}

		jsonData, err := json.Marshal(response.Data)
		assert.NoError(t, err)

		actualData := types.GetRecordResponse{}
		err = json.Unmarshal(jsonData, &actualData)
		assert.NoError(t, err)

		/// expected response struct type should match the actual response struct type
		assert.IsType(t, expectedData, actualData)

		/// actual response should not be empty
		assert.NotEmpty(t, actualData)
	})

	t.Run("Get Record By ID", func(t *testing.T) {
		req := createHTTPRequest(t, "GET", baseURL+"/v1/record/"+testID, &token, nil)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		var response viewmodels.HTTPResponseVM
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		/// success should be true
		assert.Equal(t, response.Success, true)

		expectedData := types.GetRecordResponse{}

		jsonData, err := json.Marshal(response)
		assert.NoError(t, err)

		actualData := types.GetRecordResponse{}
		err = json.Unmarshal(jsonData, &actualData)
		assert.NoError(t, err)

		/// expected response struct type should match the actual response struct type
		assert.IsType(t, expectedData, actualData)

		/// actual response should not be empty
		assert.NotEmpty(t, actualData.ID)
		assert.NotEmpty(t, actualData.Data)
		assert.NotEmpty(t, actualData.CreatedAt)
	})

	t.Run("Update Record", func(t *testing.T) {
		request := types.UpdateRecordRequest{
			ID:   testID,
			Data: gofakeit.BuzzWord(),
		}

		buf := new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(request)
		assert.NoError(t, err)

		req := createHTTPRequest(t, "PUT", baseURL+"/v1/record/"+request.ID+"/update", &token, request)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		var response viewmodels.HTTPResponseVM
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		/// success should be true
		assert.Equal(t, response.Success, true)
	})

	t.Run("Delete Record", func(t *testing.T) {
		req := createHTTPRequest(t, "DELETE", baseURL+"/v1/record/"+testID+"/delete", &token, nil)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		var response viewmodels.HTTPResponseVM
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		/// success should be true
		assert.Equal(t, response.Success, true)
	})
}

func generateID() string {
	return ksuid.New().String()
}

func createHTTPRequest(t *testing.T, method string, url string, token *string, body interface{}) *http.Request {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(body)
	assert.NoError(t, err)

	req, err := http.NewRequest(method, url, buf)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if token != nil {
		req.Header.Set("Authorization", "Bearer "+*token)
	}
	return req
}
