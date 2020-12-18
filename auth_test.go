package cascade

import (
	"testing"
	"time"
)

var secretTestCases = []struct {
	username string
	password string
	secret   string
}{
	{
		username: "username",
		password: "password",
		secret:   "dXNlcm5hbWU6cGFzc3dvcmQ=",
	},
}

func TestAuth_Secret(t *testing.T) {
	for _, test := range secretTestCases {
		auth := Auth{
			Username: test.username,
			Password: test.password,
		}

		secret := auth.Secret()

		if secret != test.secret {
			t.Errorf(`Auth.Secret(username: %s, password: %s) failed: have secret - %s, want - %s`,
				test.username, test.password, secret, test.secret)
		}
	}
}

var expiredTestCases = []struct {
	timestamp int64
	want      time.Time
}{
	{
		timestamp: 598492,
		want:      time.Date(1970, 1, 7, 22, 14, 52, 0, time.UTC),
	},
	{
		timestamp: 1586335909,
		want:      time.Date(2020, 4, 8, 8, 51, 49, 0, time.UTC),
	},
}

func TestLoginResponse_Expired(t *testing.T) {
	for _, test := range expiredTestCases {
		login := &token{expiresIn: test.timestamp}

		have := login.expired(time.UTC)

		if have != test.want {
			t.Errorf("token(timestamp: %d) failed: have %v, want %v", test.timestamp, have, test.want)
		}
	}
}
