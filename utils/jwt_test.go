package utils_test

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/dixiedream/authenticator/utils"
)

func TestGenerateAccessToken(t *testing.T) {
	cases := [...]struct {
		payload  *utils.Payload
		expected error
	}{
		{&utils.Payload{}, nil},
		{&utils.Payload{Hostname: "pippo.com"}, nil},
		{&utils.Payload{Hostname: "pippo.com", Role: 1}, nil},
		{&utils.Payload{Role: 1}, nil},
	}

	for _, c := range cases {
		token, err := utils.GenerateAccessToken(c.payload)
		if err != c.expected {
			t.Log(err.Error())
			t.Fail()
		}

		_, err = utils.AccessTokenIsValid(token)
		if err != nil {
			t.Log("Access token needs to be valid")
			t.Fail()
		}
	}
}

func TestAccessTokenIsValid(t *testing.T) {
	hostname := "pippo.com"
	role := 1
	token, _ := utils.GenerateAccessToken(&utils.Payload{Hostname: hostname, Role: role})

	payload, err := utils.AccessTokenIsValid(token)
	if err != nil {
		t.Log("Access token needs to be valid")
		t.Fail()
	}

	if payload.Hostname != hostname {
		t.Logf("%s expected, %s received", hostname, payload.Hostname)
		t.Fail()
	}

	if payload.Role != role {
		t.Logf("%d expected, %d received", role, payload.Role)
		t.Fail()
	}

	aToken := strings.Split(token, ".")
	dToken, _ := base64.RawURLEncoding.DecodeString(aToken[1])
	claims := &utils.Claim{}
	json.Unmarshal(dToken, claims)
	expiration := claims.StandardClaims.ExpiresAt
	if (expiration - time.Now().Unix()) > (60 * 15) {
		t.Log("Wrong Expiration")
		t.Fail()
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	cases := [...]struct {
		payload  *utils.Payload
		expected error
	}{
		{&utils.Payload{}, nil},
		{&utils.Payload{Hostname: "pippo.com"}, nil},
		{&utils.Payload{Hostname: "pippo.com", Role: 1}, nil},
		{&utils.Payload{Role: 1}, nil},
	}

	for _, c := range cases {
		token, err := utils.GenerateRefreshToken(c.payload)
		if err != c.expected {
			t.Log(err.Error())
			t.Fail()
		}

		_, err = utils.RefreshTokenIsValid(token)
		if err != nil {
			t.Log("Access token needs to be valid")
			t.Fail()
		}
	}
}

func TestRefreshTokenIsValid(t *testing.T) {
	hostname := "pippo.com"
	role := 1
	token, _ := utils.GenerateRefreshToken(&utils.Payload{Hostname: hostname, Role: role})

	payload, err := utils.RefreshTokenIsValid(token)
	if err != nil {
		t.Log("Refresh token needs to be valid")
		t.Fail()
	}

	if payload.Hostname != hostname {
		t.Logf("%s expected, %s received", hostname, payload.Hostname)
		t.Fail()
	}

	if payload.Role != role {
		t.Logf("%d expected, %d received", role, payload.Role)
		t.Fail()
	}

	aToken := strings.Split(token, ".")
	dToken, _ := base64.RawURLEncoding.DecodeString(aToken[1])
	claims := &utils.Claim{}
	json.Unmarshal(dToken, claims)
	expiration := claims.StandardClaims.ExpiresAt
	if (expiration - time.Now().Unix()) > (60 * 60 * 24 * 365) {
		t.Log("Wrong Expiration")
		t.Fail()
	}
}
