package token

import (
	"interview/db/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
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
