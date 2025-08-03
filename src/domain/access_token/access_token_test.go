package accesstoken

import (
	"testing"
	"time"
)

func TestGetNewAccessToken(t *testing.T) {
	at := GetAccessToken()
	if at.IsExpired() {
		t.Error("brand new access token should not be expired")
	}

	if at.AccessToken != "" {
		t.Error("new access token should not have defined access token id")
	}

	if at.UserId != 0 {
		t.Error("new access token should not have an associated user id ")
	}
}

func TestAccessTokenIsExpired(t *testing.T) {
	at := AccessToken{}
	if !at.IsExpired() {
		t.Error("empty access token should be empty by default")

		at.Expiers = time.Now().UTC().Add(3 * time.Hour).Unix()
		if at.IsExpired() {
			t.Error("access token expiring three hours from now should NOT be expired")
		}
	}
}