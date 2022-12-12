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
	"embed"
	"errors"
	"feedback/internal"
	"feedback/internal/logger"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrations embed.FS
var log = logger.Instance()

type Interface interface {
	Store(value interface{}) error
	Read(tokenValue string) (Feedback, error)
	Update(feedbackToUpdate Feedback, tokenValue string) (Feedback, error)
}

type Repository struct {
	config *internal.Configuration
	db     *gorm.DB
}

func New(config *internal.Configuration) *Repository {
	db, err := createGormDBConnection(config)
	if err != nil {
		panic(err)
	}

	return &Repository{config, db}
}

func (repo *Repository) Migrate() {
	db, err := createPlainSqlConnection(repo.config)
	if err != nil {
		log.Panic(err)
	}

	goose.SetLogger(logger.GooseLoggerWrapper(logger.Instance()))
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}
}

func (repo *Repository) Store(value interface{}) error {
	tx := repo.db.Create(value)
	tx.Commit()

	return repo.db.Error
}

func (repo *Repository) Read(tokenValue string) (Feedback, error) {
	if repo.checkIfFeedbackExists(tokenValue) == false {
		return Feedback{}, errors.New("token not found")
	}
	var feedback = Feedback{}
	repo.db.Find(&feedback, "Jwt = ?", tokenValue)

	return feedback, nil
}

func (repo *Repository) Update(feedbackToUpdate Feedback, tokenValue string) (Feedback, error) {
	if repo.checkIfFeedbackExists(tokenValue) == false {
		return Feedback{}, errors.New("token not found")
	}
	read, _ := repo.Read(tokenValue)
	read.Rating = feedbackToUpdate.Rating
	read.RatingComment = feedbackToUpdate.RatingComment
	read.Metadata = feedbackToUpdate.Metadata
	repo.db.Save(&read)
	return feedbackToUpdate, nil
}

func (repo *Repository) checkIfFeedbackExists(tokenValue string) bool {
	var feedback = &Feedback{}
	repo.db.Find(&feedback, "Jwt = ?", tokenValue)
	if feedback.ID == 0 {
		return false
	}
	return true
}

func (repo *Repository) Count() int64 {
	var count int64
	repo.db.Table("feedbacks").Count(&count)
	return count
}
