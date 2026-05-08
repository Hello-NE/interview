package token

import (
	"interview/db/util"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	if err != nil {
		t.Fatal(err)
	}
	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	if err != nil {
		t.Fatal(err)
	}
	payload, err := maker.VerifyToken(token)
	if err != nil {
		t.Fatal(err)
	}
	if payload.Username != username {
		t.Errorf("expected username %s, got %s", username, payload.Username)
	}

	require.NotZero(t, payload.ID)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

// TestJWTMakerExpiredToken tests the expired token case.
func TestJWTMakerExpiredToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	if err != nil {
		t.Fatal(err)
	}
	username := util.RandomOwner()
	duration := -time.Minute

	token, err := maker.CreateToken(username, duration)
	if err != nil {
		t.Fatal(err)
	}
	payload, err := maker.VerifyToken(token)
	if err == nil {
		t.Error("expected error, got nil")
	}
	require.Nil(t, payload)
}

func TestJWTMakerInvalidToken(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
}
