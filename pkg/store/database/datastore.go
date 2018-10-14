/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package datastore

import (
	"time"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"

	"github.com/angao/recommender/pkg/store"

	"github.com/golang/glog"
)

// datastore is an implementation of a model.Store built on top
// of the sql/database driver with a relational database backend.
type datastore struct {
	Engine *xorm.Engine
	driver string
}

// New creates a database connection for the given driver and datasource
// and returns a new Store.
func New(driver, config string) store.Store {
	return &datastore{
		Engine: create(driver, config),
		driver: driver,
	}
}

// create opens a new database connection with the specified
// driver and connection string and returns a store.
func create(driver, config string) *xorm.Engine {
	engine, err := xorm.NewEngine(driver, config)
	if err != nil {
		glog.Fatalf("database connection failed: %#v", err)
	}

	if err := engine.Ping(); err != nil {
		glog.Fatalf("database ping attempts failed: %#v", err)
	}
	// engine.ShowSQL(true)

	go pingDatabase(engine)
	engine.SetTableMapper(core.NewPrefixMapper(core.SnakeMapper{}, "t_"))
	if driver == "mysql" {
		// per issue https://github.com/go-sql-driver/mysql/issues/257
		engine.SetMaxIdleConns(0)
	}
	return engine
}

func pingDatabase(engine *xorm.Engine) {
	timer := time.NewTicker(10 * time.Minute)
	for _ = range timer.C {
		err := engine.Ping()
		if err != nil {
			glog.Errorf("database ping failed. retry in 10m. error: %#v", err)
		}
	}
}
