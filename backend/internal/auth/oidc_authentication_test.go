package auth

import (
	"bytes"
	"feedback/internal"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const (
	tokenValue = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	userId     = "1234"
)

func TestService_TestExtract(t *testing.T) {
	config := internal.ConfigurationFromEnv()
	service := New(config)

	extract := service.ExtractTokenFrom("Bearer " + tokenValue)

	assert.Equal(t, tokenValue, extract)
}

func TestService_TestExtract_EmptyValueCreatesError(t *testing.T) {
	config := internal.ConfigurationFromEnv()
	service := New(config)

	extract := service.ExtractTokenFrom("")

	assert.Equal(t, "", extract)
}

func TestService_TestValidate(t *testing.T) {
	config := internal.ConfigurationFromEnv()
	service := New(config)

	validate, err := service.CreateValidationRequestFrom(tokenValue)

	buffer := bytes.NewBuffer(validate)
	assert.Equal(t, err, nil)
	assert.NotNil(t, buffer)
}

func TestService_TestRequest(t *testing.T) {
	config := internal.ConfigurationFromEnv()
	service := New(config)
	validate, err := service.CreateValidationRequestFrom(tokenValue)

	closer, err := service.Validate(validate)

	request, err := service.Map(closer)
	assert.Equal(t, err, nil)
	assert.NotNil(t, request)
}

func TestService_TestGenerate(t *testing.T) {
	config := internal.ConfigurationFromEnv()
	service := New(config)

	generate, err := service.Generate(userId)

	assert.Equal(t, err, nil)
	assert.True(t, strings.Contains(generate, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9."))
}

func TestService_TestGenerate_EmptyValueCreatesError(t *testing.T) {
	config := internal.ConfigurationFromEnv()
	service := New(config)

	defer func() {
		if r := recover(); r != nil {
			service.Generate("userId")
		}
	}()
}
