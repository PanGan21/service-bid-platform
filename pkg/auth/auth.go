package auth

import (
	"errors"
	"fmt"

	"github.com/PanGan21/pkg/entity"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mitchellh/mapstructure"
)

type AuthService interface {
	SignJWT(sessionId string, user entity.PublicUser, path string, roles ...string) (string, error)
	VerifyJWT(encoded string, route string) (AuthTokenData, error)
}

type authService struct {
	secret []byte
}

func NewAuthService(secret []byte) AuthService {
	return &authService{secret: secret}
}

func (s *authService) SignJWT(sessionId string, user entity.PublicUser, path string, roles ...string) (string, error) {
	internalPath, _ := SplitPath(path)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["service"] = internalPath.Service
	claims["route"] = internalPath.Route
	claims["sessionId"] = sessionId
	claims["user"] = user
	claims["roles"] = roles

	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) VerifyJWT(encoded string, route string) (AuthTokenData, error) {
	var user entity.PublicUser

	authTokenData := AuthTokenData{
		Service:   "",
		Route:     "",
		SessionId: "",
		User:      user,
		Roles:     make([]string, 0),
	}

	token, err := jwt.Parse(encoded, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return authTokenData, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		authTokenData.Service = claims["service"].(string)
		authTokenData.Route = claims["route"].(string)
		authTokenData.SessionId = claims["sessionId"].(string)
		mapstructure.Decode(claims["user"], &authTokenData.User)
		parsedRoles := claims["roles"].([]interface{})

		for _, role := range parsedRoles {
			authTokenData.Roles = append(authTokenData.Roles, fmt.Sprintf("%s", role))
		}
	}

	if authTokenData.Route != route {
		return authTokenData, errors.New("incorrect route")
	}

	return authTokenData, nil
}
