package auth

import (
	"errors"
	"testing"
)

var secret = "mockSecret"
var service = "mockService"
var route = "/route/web"
var path = "/" + service + route + "?query=something"
var sessionId = "mockSessionId"
var userId = "mockUserId"
var roles = []string{"mockRole1", "mockRole2"}
var mockJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlcyI6WyJtb2NrUm9sZTEiLCJtb2NrUm9sZTIiXSwicm91dGUiOiIvcm91dGUvd2ViIiwic2VydmljZSI6Im1vY2tTZXJ2aWNlIiwic2Vzc2lvbklkIjoibW9ja1Nlc3Npb25JZCIsInVzZXJJZCI6Im1vY2tVc2VySWQifQ.KE0GOFD4fk4b3C0CSl7zxoAMAxl7v2FvkYS-LgUQmRs"

func TestSingJWT(t *testing.T) {
	authService := NewAuthService([]byte(secret))

	token, err := authService.SignJWT(sessionId, userId, path, roles...)
	if err != nil {
		t.Fatal(err)
	}

	if token != mockJWT {
		t.Fatal("incorrect jwt")
	}
}

func TestVerifyJWT(t *testing.T) {
	authService := NewAuthService([]byte(secret))
	authData, err := authService.VerifyJWT(mockJWT, route)
	if err != nil {
		t.Fatal(err)
	}

	if authData.Service != service || authData.UserId != userId || authData.Route != route || authData.SessionId != sessionId {
		t.Fatal(errors.New("incorrect verification"))
	}
}

func TestMatchRoute(t *testing.T) {
	internalPath, err := SplitPath(path)
	if err != nil {
		t.Fatal(err)
	}

	if internalPath.Service != service {
		t.Fatal("service does not match")
	}

	if internalPath.Route != route {
		t.Fatal("route does not match")
	}
}
