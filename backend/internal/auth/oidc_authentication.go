package auth

import (
	"encoding/json"
	"errors"
	"feedback/internal"
	"feedback/internal/api"
	"feedback/internal/client"
	"feedback/internal/logger"
	"github.com/golang-jwt/jwt"
	"io"
	"net/http"
	"regexp"
	"time"
)

var log = logger.Instance()

type OidcAuthentication struct {
	config *internal.Configuration
}

func New(config *internal.Configuration) *OidcAuthentication {
	return &OidcAuthentication{config}
}

func (auth OidcAuthentication) Validate(request *http.Request) (*string, error) {
	validationResponse, err := auth.validate(request)
	if err != nil {
		return nil, err
	}
	if validationResponse.Results.User == true && len(validationResponse.UserId) > 0 {
		with, err := auth.generateWith(validationResponse.UserId)
		return &with, err
	} else {
		return nil, err
	}
}

func (auth OidcAuthentication) validate(request *http.Request) (*api.ValidationResponse, error) {
	requestBody, err := createValidationRequestFrom(request)
	if err != nil || len(requestBody) == 0 {
		return nil, err
	}

	response, err := client.Post(auth.config.OidcValidationUrl, requestBody)
	if err != nil {
		return nil, err
	}

	validationResponse, err := mapFrom(response)

	return validationResponse, err
}

func createValidationRequestFrom(request *http.Request) ([]byte, error) {
	token := extractTokenFrom(request)
	if len(token) > 0 {
		requestBody, err := json.Marshal(map[string]string{
			"matrix_server_name": "domain.tld",
			"token":              token,
		})
		if err != nil {
			return []byte(""), err
		}
		return requestBody, err
	}
	return []byte(""), errors.New("token value is empty")
}

func extractTokenFrom(request *http.Request) string {
	authHeaderValue := request.Header.Get("authorization")

	return regexp.MustCompile(`\W`).Split(authHeaderValue, -1)[1]
}

func mapFrom(respBody io.ReadCloser) (*api.ValidationResponse, error) {
	body, err := io.ReadAll(respBody)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	var validationResponse *api.ValidationResponse
	if err := json.Unmarshal(body, &validationResponse); err != nil {
		return nil, err
	}

	return validationResponse, err
}

func (auth OidcAuthentication) generateWith(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"nbf":    time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(auth.config.JwtSignature))

	return tokenString, err
}
