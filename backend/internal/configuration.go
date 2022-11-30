/*
 *  Copyright 2022 Nordeck IT + Consulting GmbH
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and  limitations
 *  under the License.
 *
 */

package internal

import (
	"fmt"
	"os"
	"reflect"
)

type Configuration struct {
	DbHost            string // DB_HOST
	DbPort            string // DB_PORT
	DbUser            string // DB_USER
	DbPassword        string // DB_PASSWORD
	DbName            string // DB_NAME
	Sslmode           string "disable" // SSL_MODE
	OidcValidationUrl string // OIDC_VALIDATION_URL
	JwtSignature      string // JWT_SIGNATURE
	MatrixServerName  string // MATRIX_SERVER_NAME
}

func ConfigurationFromEnv() *Configuration {
	config := Configuration{
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("SSL_MODE"),
		os.Getenv("OIDC_VALIDATION_URL"),
		os.Getenv("JWT_SIGNATURE"),
		os.Getenv("MATRIX_SERVER_NAME"),
	}

	elements := reflect.ValueOf(&config).Elem()

	for i := 0; i < elements.NumField(); i++ {
		varValue := elements.Field(i).Interface()
		if varValue == "" {
			panic(fmt.Sprintf("%s not set.", elements.Type().Field(i).Name))
		}
	}

	return &config
}
