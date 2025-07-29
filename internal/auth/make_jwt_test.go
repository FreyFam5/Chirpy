package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)


var testSecretToken = "thisIsASecret"

func TestMakeJWT(t *testing.T) {
	_, err := MakeJWT(uuid.New(), testSecretToken, time.Duration(5 * time.Second))
	if err != nil {
		t.Errorf("Couldn't make JWT: %s\n", err)
	}
}

func TestValidateJMT(t *testing.T) {
	id := uuid.New()

	tokenString, err := MakeJWT(id, testSecretToken, time.Duration(5 * time.Second))
	if err != nil {
		t.Errorf("Couldn't make JWT: %s\n", err)
	}

	valID, err := ValidateJWT(tokenString, testSecretToken)
	if err != nil {
		t.Errorf("Couldn't validate JMT: %s\n", err)
	}

	if valID != id {
		t.Error("Id's do not match\n")
	}
}	

func TestExpiredTokens(t *testing.T) {
	id := uuid.New()

	tokenString, err := MakeJWT(id, testSecretToken, time.Duration(time.Millisecond))
	if err != nil {
		t.Errorf("Couldn't make JMT: %s", err)
	}

	<- time.Tick(time.Duration(2 * time.Millisecond))

	_, err = ValidateJWT(tokenString, testSecretToken)
	if err != nil {
		if err.Error() != "token has invalid claims: token is expired" {
			t.Errorf("Couldn't validate JMT: %s", err)
		}
	}
}

func TestDifferentSecretTokens(t *testing.T) {
	id := uuid.New()

	tokenString, err := MakeJWT(id, testSecretToken, time.Duration(time.Nanosecond))
	if err != nil {
		t.Errorf("Couldn't make JMT: %s", err)
	}

	_, err = ValidateJWT(tokenString, "This isn't the secret token!")
	if err != nil {
		if err.Error() != "token signature is invalid: signature is invalid" {
			t.Errorf("Couldn't validate JMT: %s", err)
		}
	}
}