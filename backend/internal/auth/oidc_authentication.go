package auth

import (
	"encoding/json"
	"feedback/internal"
	"feedback/internal/api"
	"feedback/internal/client"
	"feedback/internal/logger"
	"github.com/golang-jwt/jwt"
	"io"
	"strings"
	"time"
)

var log = logger.Instance()

type OidcAuthentication struct {
	config *internal.Configuration
}

func New(config *internal.Configuration) *OidcAuthentication {
	return &OidcAuthentication{config}
}

func (s OidcAuthentication) ExtractTokenFrom(authHeadVal string) string {

	if authHeadVal == "" {
		err := "authentication header is empty!"
		log.Error(err)
		return ""
	}
	trimmedHeaderValue := strings.Trim(authHeadVal, "Bearer")
	return strings.ReplaceAll(trimmedHeaderValue, " ", "")
}

func (s OidcAuthentication) CreateValidationRequestFrom(token string) ([]byte, error) {
	reqBody, err := json.Marshal(map[string]string{
		"matrix_server_name": "domain.tld",
		"token":              token,
	})
	if err != nil {
		log.Error(err)
		return []byte(""), err
	}

	return reqBody, err
}

func (s OidcAuthentication) Validate(requestBody []byte) (io.ReadCloser, error) {
	return client.Post(s.config.OidcValidationUrl, requestBody)
}

func (s OidcAuthentication) Map(respBody io.ReadCloser) (api.ValidationResponse, error) {
	body, err := io.ReadAll(respBody)
	if err != nil {
		log.Error(err)
		return api.ValidationResponse{}, err
	}

	var valRes api.ValidationResponse
	if err := json.Unmarshal(body, &valRes); err != nil {
		log.Error(err)
	}

	return valRes, err
}

func (s OidcAuthentication) Generate(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"nbf":    time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.config.JwtSignature))

	return tokenString, err
}
