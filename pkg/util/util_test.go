package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRandomString(t *testing.T) {
	randomString := RandomString(32)
	require.Len(t, randomString, 32)
}

func TestRandomMobile(t *testing.T) {
	mobile := RandomMobile()
	require.Len(t, mobile, 11)
}
