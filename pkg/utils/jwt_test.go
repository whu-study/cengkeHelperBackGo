package utils

import (
	"cengkeHelperBackGo/internal/config"
	"testing"
)

func TestGenerateAndParseUserToken(t *testing.T) {
	// ensure config has a key
	config.Conf.JwtSecurityKey = "testkey123456"
	// generate
	token, err := GenerateUserToken("alice", 1, "uid-123")
	if err != nil {
		t.Fatalf("GenerateUserToken failed: %v", err)
	}

	claims, err := ParseUserJwt(token)
	if err != nil {
		t.Fatalf("ParseUserJwt failed: %v", err)
	}

	if claims.Username != "alice" || claims.UserId != "uid-123" || claims.Role != 1 {
		t.Fatalf("claims mismatch: %+v", claims)
	}
}
