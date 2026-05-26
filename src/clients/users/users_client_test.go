package users

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/oauth_errors"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/logger"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const testBaseUrl = "http://localhost:8080/users/login"

func setUp() (*usersClient, func()) {
	restyClient := resty.New().SetTimeout(150 * time.Millisecond)
	httpmock.ActivateNonDefault(restyClient.GetClient())

	loggerCfg := logger.Config{
		Level:       "info",
		OutputPaths: []string{"stdout"},
	}
	log, _ := logger.NewLogger(loggerCfg)

	client := &usersClient{
		restClient: restyClient,
		logger:     log,
		apiBaseUrl: testBaseUrl,
	}

	tearDown := func() {
		httpmock.DeactivateAndReset()
	}

	return client, tearDown
}

func TestLoginUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		client, tearDown := setUp()
		defer tearDown()

		httpmock.RegisterResponder(http.MethodPost, testBaseUrl,
			httpmock.NewJsonResponderOrPanic(http.StatusOK, User{
				Id:    1,
				Email: "test@test.com",
			}),
		)

		userId, err := client.LoginUser(context.Background(), "test@test.com", "password")

		assert.NoError(t, err)
		assert.Equal(t, int64(1), userId)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		client, tearDown := setUp()
		defer tearDown()

		httpmock.RegisterResponder(http.MethodPost, testBaseUrl,
			httpmock.NewJsonResponderOrPanic(http.StatusNotFound, map[string]interface{}{
				"message": "user not found",
				"status":  404,
				"error":   "not_found",
			}),
		)

		userId, err := client.LoginUser(context.Background(), "unknown@test.com", "password")

		assert.Equal(t, int64(0), userId)
		assert.Error(t, err)
		assert.ErrorIs(t, err, oauth_errors.NotFoundErr)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		client, tearDown := setUp()
		defer tearDown()

		httpmock.RegisterResponder(http.MethodPost, testBaseUrl,
			httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, map[string]interface{}{
				"message": "internal server error",
				"status":  500,
			}),
		)

		userId, err := client.LoginUser(context.Background(), "test@test.com", "password")

		assert.Equal(t, int64(0), userId)
		assert.Error(t, err)
		assert.ErrorIs(t, err, oauth_errors.InternalServerErr)
	})

	t.Run("ConnectionError", func(t *testing.T) {
		client, tearDown := setUp()
		defer tearDown()

		httpmock.RegisterResponder(http.MethodPost, testBaseUrl,
			httpmock.NewErrorResponder(assert.AnError),
		)

		userId, err := client.LoginUser(context.Background(), "test@test.com", "password")

		assert.Equal(t, int64(0), userId)
		assert.ErrorIs(t, err, oauth_errors.InternalServerErr)
	})
}
