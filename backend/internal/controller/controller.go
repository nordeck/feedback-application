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
	router.HandleFunc(TokenPath, c.returnOptions).Methods(http.MethodOptions)
	router.HandleFunc(FeedbackPath, c.createFeedback).Methods(http.MethodPost)
	router.HandleFunc(FeedbackPath, c.returnOptions).Methods(http.MethodOptions)
	return router
}

func (c *Controller) createToken(writer http.ResponseWriter, request *http.Request) {
	addAccessControlHeaders(writer, request)
	jwt, err := auth.New(internal.ConfigurationFromEnv()).Validate(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(writer).Encode(jwt)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Controller) createFeedback(writer http.ResponseWriter, request *http.Request) {
	addAccessControlHeaders(writer, request)
	authentication := auth.New(internal.ConfigurationFromEnv())
	tokenString, err := authentication.ExtractTokenFrom(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		return
	}
	authorized, err := authentication.IsAuthorized(tokenString)
	if err != nil || !authorized {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		return
	}
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

	fromDatabase, err := c.repo.Read(*tokenString)
	if fromDatabase.Jwt == *tokenString {
		_, err := c.repo.Update(fromDatabase, *tokenString)
		if err == nil {
			return
		} else {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			log.Debug(err)
			return
		}
	}
	err = c.repo.Store(repository.MapToFeedbackModel(feedback, *tokenString))

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Debug(err)
		return
	}
}

func (c *Controller) returnOptions(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", request.Header.Get("Access-Control-Request-Headers"))
	writer.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,PUT,PATCH,POST,DELETE")
	writer.WriteHeader(http.StatusNoContent)
	return
}

func addAccessControlHeaders(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "*")
}
