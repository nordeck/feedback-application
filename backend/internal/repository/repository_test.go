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

package repository

import (
	"context"
	"errors"
	"feedback/internal"
	gormjsonb "github.com/dariubs/gorm-jsonb"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image: "postgres:latest",
		Env: map[string]string{
			"POSTGRES_PASSWORD": "postgres",
		},
		AutoRemove:   true,
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		panic(err)
	}

	err = os.Setenv("DB_PORT", port.Port())
	if err != nil {
		panic(err)
	}

	code := m.Run()

	err = postgresContainer.Terminate(ctx)
	if err != nil {
		panic(err)
	}

	os.Exit(code)
}

func TestRepository_CRU_Roundtrip(t *testing.T) {
	conf := internal.ConfigurationFromEnv()
	repo := New(conf)
	repo.Migrate()

	comment := "total doof"
	rating := 3

	tokenValue := "someJwt"
	anotherTokenValue := "anotherJwt"

	feedback := Feedback{
		Rating:        rating,
		RatingComment: comment,
		Metadata:      gormjsonb.JSONB{"first_key": "first_value", "second_key": "second_value"},
		Jwt:           tokenValue,
	}

	feedback2 := Feedback{
		Rating:        rating,
		RatingComment: comment,
		Metadata:      gormjsonb.JSONB{"first_key": "first_value", "second_key": "second_value"},
		Jwt:           anotherTokenValue,
	}
	// CREATE
	err := repo.Store(&feedback)
	err2 := repo.Store(&feedback2)

	if err != nil || err2 != nil {
		panic(err)
	}

	count := repo.Count()
	assert.Equal(t, count, int64(2))

	// READ
	readBeforeUpdate, err := repo.FindByToken(tokenValue)
	assert.Equal(t, readBeforeUpdate.Rating, rating)
	assert.Equal(t, readBeforeUpdate.RatingComment, comment)
	assert.Equal(t, readBeforeUpdate.Metadata, gormjsonb.JSONB{"first_key": "first_value", "second_key": "second_value"})
	assert.Equal(t, readBeforeUpdate.Jwt, tokenValue)

	feedbackToUpdate := Feedback{
		Rating:        -1,
		RatingComment: comment,
		Metadata:      gormjsonb.JSONB{"first_key": "first_value", "second_key": "second_value"},
		Jwt:           tokenValue,
	}

	// UPDATE
	_, err = repo.Update(feedbackToUpdate)
	if err != nil {
		panic(err)
	}

	// READ
	readAfterUpdate, err := repo.FindByToken(tokenValue)
	assert.Equal(t, readAfterUpdate.Rating, -1)
	assert.Equal(t, readAfterUpdate.RatingComment, comment)
	assert.NotNil(t, readAfterUpdate.CreatedAt)
	assert.Equal(t, readAfterUpdate.Metadata, gormjsonb.JSONB{"first_key": "first_value", "second_key": "second_value"})
	assert.Equal(t, readAfterUpdate.Jwt, tokenValue)

	countAfter := repo.Count()
	assert.Equal(t, countAfter, int64(2))

}

func TestRepository_CRU_Roundtrip_Read_TokenValueNotFound(t *testing.T) {
	conf := internal.ConfigurationFromEnv()
	repo := New(conf)
	repo.Migrate()

	// READ
	_, err := repo.FindByToken("tokenValNotAvailable")

	if err == nil {
		// no error occurred, this should not happen.
		panic(err)
	}
	assert.Equal(t, err, errors.New("no record with token value found in database"))
}

func TestRepository_CRU_Roundtrip_Update_TokenValueNotFound(t *testing.T) {
	conf := internal.ConfigurationFromEnv()
	repo := New(conf)
	repo.Migrate()

	comment := "total doof"
	rating := 3
	metadata := map[string]interface{}{
		"first_key":  "first_value",
		"second_key": "second_value",
	}

	feedbackToUpdate := Feedback{
		Rating:        rating,
		RatingComment: comment,
		Metadata:      metadata,
		Jwt:           "tokenValNotAvailable",
	}

	// READ
	_, err := repo.Update(feedbackToUpdate)

	if err == nil {
		// no error occurred, this should not happen.
		panic(err)
	}
	assert.Equal(t, err, errors.New("no record found for update"))
}
