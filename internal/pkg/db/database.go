// Copyright 2021 INTECH Process Automation Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package db

import (
	"fmt"
	"github.com/lf-edge/ekuiper/internal/pkg/db/redis"
	"github.com/lf-edge/ekuiper/internal/pkg/db/sql/sqlite"
)

type Database interface {
	Connect() error
	Disconnect() error
}

func CreateDatabase(conf Config) (error, Database) {
	var db Database
	var err error
	databaseType := conf.Type
	switch databaseType {
	case "redis":
		r := redis.NewRedisFromConf(conf.Redis)
		db = &r
	case "sqlite":
		err, db = sqlite.NewSqliteDatabase(conf.Sqlite)
		if err != nil {
			return err, nil
		}
	default:
		return fmt.Errorf("unrecognized database type - %s", databaseType), nil
	}
	err = db.Connect()
	if err != nil {
		return err, nil
	}
	return nil, db
}