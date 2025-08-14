package rest

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

type slowTransport struct {
    original http.RoundTripper
    delay    time.Duration
}

func (s *slowTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    time.Sleep(s.delay) 
    return s.original.RoundTrip(req)
}

func TestMain(m *testing.M) {
	fmt.Println("==> Running TestMain")
    gock.InterceptClient(usersRestClient.GetClient()) 
    code := m.Run()
    gock.Off()
    os.Exit(code)
}

func TestLoginUserTimeoutFromApi(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bookstore.com").
		Post("/users/login")

	clientTransport := &slowTransport{
        original: http.DefaultTransport,
        delay:    200 * time.Millisecond,
    }
	originalTransport := usersRestClient.GetClient().Transport
	defer usersRestClient.SetTransport(originalTransport)
	usersRestClient.SetTransport(clientTransport)

	repository := usersRepository{}
	response, err := repository.LoginUsesr("test@gmail.com", "test")

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Message, "Client.Timeout"))
}

func TestLoginUserNetworkError(t *testing.T) {
	defer gock.Off()
	gock.DisableNetworking()

	repository := usersRepository{}
	response, err := repository.LoginUsesr("test@gmail.com", "test")

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Message, "context deadline exceeded"))
}

func TestLoginUserInvalidLoginCredentials(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bookstore.com").
		Post("/users/login").
		Reply(http.StatusNotFound).
		JSON(map[string]interface{}{
    		"message": "invalid login credentials",
    		"status":  404,
    		"error":   "not_found",
		})

	repository := usersRepository{}
	response, err := repository.LoginUsesr("test@gmail.com", "test")


	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "invalid login credentials", err.Message)
	assert.EqualValues(t, http.StatusNotFound, err.Status)
}

func TestLoginUserNoError(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bookstore.com").
		Post("/users/login").
		Reply(http.StatusOK).
		JSON(map[string]interface{}{
    		"id": 1,
    		"first_name":  "test",
    		"last_name":   "Test",
			"email": "test@gmail.com",
		})

	repository := usersRepository{}
	user, err := repository.LoginUsesr("test@gmail.com", "test")


	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, 1, user.Id)
	assert.EqualValues(t, "test", user.FirstName)
	assert.EqualValues(t, "Test", user.LastName)
	assert.EqualValues(t, "test@gmail.com", user.Email)
}
