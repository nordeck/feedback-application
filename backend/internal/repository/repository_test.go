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

	os.Setenv("DB_PORT", port.Port())

	code := m.Run()

	err = postgresContainer.Terminate(ctx)
	if err != nil {
		panic(err)
	}

	os.Exit(code)
}

func TestRepository_Store(t *testing.T) {
	conf := internal.ConfigurationFromEnv()
	repo := New(conf)
	repo.Migrate()

	comment := "total doof"
	rating := 3
	metadata := map[string]interface{}{
		"first_key":  "first_value",
		"second_key": "second_value",
	}

	feedback := Feedback{
		Rating:        rating,
		RatingComment: comment,
		Metadata:      metadata,
	}
	err := repo.Store(&feedback)
	if err != nil {
		panic(err)
	}

	dbFeedback := Feedback{}
	repo.db.Find(&dbFeedback)

	assert.Equal(t, dbFeedback.ID, uint(1))
	assert.Equal(t, dbFeedback.Rating, rating)
	assert.Equal(t, dbFeedback.RatingComment, comment)
	assert.NotNil(t, dbFeedback.CreatedAt)
	assert.Equal(t, dbFeedback.Metadata, gormjsonb.JSONB{"first_key": "first_value", "second_key": "second_value"})
}
