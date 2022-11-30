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
	router.HandleFunc(FeedbackPath, c.createFeedback).Methods(http.MethodPost)

	return router
}

func (c *Controller) createToken(writer http.ResponseWriter, request *http.Request) {
	jwt, err := auth.New(internal.ConfigurationFromEnv()).Validate(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(writer).Encode(jwt)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Controller) createFeedback(writer http.ResponseWriter, request *http.Request) {
	var feedback api.Feedback
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &feedback)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		log.Debug(err)
		return
	}

	err = c.repo.Store(repository.MapToFeedbackModel(feedback))

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Debug(err)
		return
	}
}
