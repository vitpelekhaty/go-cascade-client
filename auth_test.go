package cascade

import (
	"testing"
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
