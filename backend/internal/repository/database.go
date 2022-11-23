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
	"database/sql"
	"feedback/internal"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	DATABASE_DRIVER = "postgres"
)

func createGormDBConnection(config *internal.Configuration) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.DbHost, config.DbUser, config.DbPassword, config.DbName, config.DbPort, config.Sslmode)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})

}

func createPlainSqlConnection(config *internal.Configuration) (*sql.DB, error) {
	connectString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.DbUser, config.DbPassword, config.DbHost, config.DbPort, config.DbName, config.Sslmode)
	return sql.Open(DATABASE_DRIVER, connectString)
}
