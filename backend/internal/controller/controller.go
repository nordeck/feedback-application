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

package controller

import (
	"encoding/json"
	"errors"
	"feedback/internal"
	"feedback/internal/api"
	"feedback/internal/auth"
	"feedback/internal/logger"
	"feedback/internal/repository"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

const (
	TokenPath    = "/token"
	FeedbackPath = "/feedback"
)

var log = logger.Instance()

type Controller struct {
	repo repository.Interface
	serv *auth.OidcAuthentication
}

func New(repo repository.Interface, serv *auth.OidcAuthentication) *Controller {
	return &Controller{repo, serv}
}

func (c *Controller) GetRouter() http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc(TokenPath, c.createToken).Methods(http.MethodGet)
	router.HandleFunc(TokenPath, c.returnOptions).Methods(http.MethodOptions)
	router.HandleFunc(FeedbackPath, c.createFeedback).Methods(http.MethodPost)
	router.HandleFunc(FeedbackPath, c.returnOptions).Methods(http.MethodOptions)
	return router
}

func (c *Controller) createToken(writer http.ResponseWriter, request *http.Request) {
	addAccessControlHeaders(writer)
	jwt, err := auth.New(internal.ConfigurationFromEnv()).Validate(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "text/plain")
	_, err = writer.Write([]byte(*jwt))

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Controller) createFeedback(writer http.ResponseWriter, request *http.Request) {
	addAccessControlHeaders(writer)
	authentication := auth.New(internal.ConfigurationFromEnv())

	tokenString, err, authorized := c.authenticate(authentication, request)
	if err != nil || !authorized {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		log.Debug(err)
		return
	}

	feedback, err := c.parseFeedback(err, request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		log.Debug(err)
		return
	}

	err = c.createOrUpdate(tokenString, feedback)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Debug(err)
		return
	}
}

func (c *Controller) createOrUpdate(tokenString *string, feedback api.Feedback) error {
	fromDatabase, err := c.repo.FindByToken(*tokenString)
	if err == nil {

		if fromDatabase.Jwt == *tokenString {
			log.Debug("token found in database, updating values")
			feedbackToUpdateModel := *repository.MapToFeedbackModel(feedback, *tokenString)
			_, err := c.repo.Update(feedbackToUpdateModel)
			if err != nil {
				return errors.New("update of values failed")
			} else {
				return nil
			}
		}
	}
	return c.repo.Store(repository.MapToFeedbackModel(feedback, *tokenString))
}

func (c *Controller) authenticate(authentication *auth.OidcAuthentication, request *http.Request) (*string, error, bool) {
	tokenString, err := authentication.ExtractTokenFrom(request)
	if err != nil {
		return nil, err, false
	}
	authorized, err := authentication.IsAuthorized(tokenString)
	return tokenString, err, authorized
}

func (c *Controller) parseFeedback(err error, request *http.Request) (api.Feedback, error) {
	var feedback api.Feedback
	body, err := io.ReadAll(request.Body)
	err = json.Unmarshal(body, &feedback)
	return feedback, err
}

func (c *Controller) returnOptions(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", request.Header.Get("Access-Control-Request-Headers"))
	writer.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,PUT,PATCH,POST,DELETE")
	writer.WriteHeader(http.StatusNoContent)
	return
}

func addAccessControlHeaders(writer http.ResponseWriter) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "*")
}
