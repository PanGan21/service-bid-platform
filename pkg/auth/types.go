package auth

import (
	"errors"
	"regexp"
	"strings"

	"github.com/PanGan21/pkg/entity"
)

type AuthTokenData struct {
	Service   string
	Route     string
	SessionId string
	User      entity.PublicUser
	Roles     []string
}

type InternalPath struct {
	Service string
	Route   string
}

func SplitPath(path string) (InternalPath, error) {
	var internalPath = InternalPath{Service: "", Route: "/"}

	path = strings.Split(path, "?")[0]
	regex, err := regexp.Compile("/([^/]*)(.*)")
	if err != nil {
		return internalPath, err
	}

	match := regex.FindStringSubmatch(path)

	if len(match) > 0 {
		internalPath.Service = match[1]
	}

	if len(match) > 1 {
		internalPath.Route = match[2]
	}

	if internalPath.Service == "" {
		return internalPath, errors.New("missing service")
	}

	return internalPath, nil
}
