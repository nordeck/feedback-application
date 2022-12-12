package auth

import (
	"encoding/json"
	"errors"
	"feedback/internal"
	"feedback/internal/api"
	"feedback/internal/client"
	"github.com/golang-jwt/jwt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

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
		token, err := auth.generate()
		return &token, err
	} else {
		return nil, errors.New("user is not valid")
	}
}

func (auth OidcAuthentication) IsAuthorized(tokenString *string) (bool, error) {
	sanitized := strings.Replace(*tokenString, "\"", "", -1)
	parsedJwt, err := auth.parseJwt(&sanitized)
	if err != nil {
		return false, err
	}

	return auth.validateJwt(parsedJwt, err)
}

func (auth OidcAuthentication) validate(request *http.Request) (*api.ValidationResponse, error) {
	requestBody, err := auth.createValidationRequestFrom(request)
	if err != nil || len(requestBody) == 0 {
		return nil, err
	}

	response, err := client.Post(auth.config, requestBody)
	if err != nil {
		return nil, err
	}

	return mapFrom(response)
}

func (auth OidcAuthentication) createValidationRequestFrom(request *http.Request) ([]byte, error) {
	token, err := auth.ExtractTokenFrom(request)
	if err != nil || token == nil {
		return nil, err
	}
	requestBody, err := json.Marshal(map[string]string{
		"matrix_server_name": auth.config.MatrixServerName,
		"token":              *token,
	})
	if err != nil {
		return nil, err
	}
	return requestBody, err
}

func (auth OidcAuthentication) ExtractTokenFrom(request *http.Request) (*string, error) {
	authHeaderValue := request.Header.Get("authorization")
	var bearerRegExp = "^Bearer\\s+(.+)$"

	matched, err := regexp.MatchString(bearerRegExp, authHeaderValue)
	if matched {
		return &strings.Fields(authHeaderValue)[1], err
	}
	err = errors.New("authentication header value has not matched / is not a bearer token")
	return nil, err
}

func (auth OidcAuthentication) validateJwt(token *jwt.Token, err error) (bool, error) {
	if token.Valid {
		return true, err
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			err = errors.New("token malformed")
			return false, err
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			err = errors.New("token is either expired or not active yet")
			return false, err
		} else {
			err = errors.New("couldn't handle this token: " + err.Error())
			return false, err
		}
	} else {
		err = errors.New("couldn't handle this token: " + err.Error())
		return false, err
	}
}

func (auth OidcAuthentication) parseJwt(tokenString *string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(
		*tokenString,
		&jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(auth.config.JwtSecret), nil
		},
	)
	return token, err
}

func mapFrom(respBody io.ReadCloser) (*api.ValidationResponse, error) {
	body, err := io.ReadAll(respBody)

	if err != nil {
		return nil, err
	}

	var validationResponse *api.ValidationResponse
	if err := json.Unmarshal(body, &validationResponse); err != nil {
		return nil, err
	}

	return validationResponse, err
}

func (auth OidcAuthentication) generate() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nbf": time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(auth.config.JwtSecret))

	return tokenString, err
}
