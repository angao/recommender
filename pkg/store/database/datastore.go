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

	"github.com/angao/recommender/pkg/store/logger"

	"github.com/angao/recommender/pkg/store"
	"github.com/angao/recommender/pkg/utils"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"

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
func New(driver string, config utils.DatabaseConfig) store.Store {
	return &datastore{
		Engine: create(driver, config),
		driver: driver,
	}
}

// create opens a new database connection with the specified
// driver and connection string and returns a store.
func create(driver string, config utils.DatabaseConfig) *xorm.Engine {
	schema := config.Format()
	engine, err := xorm.NewEngine(driver, schema)
	if err != nil {
		glog.Fatalf("database connection failed: %#v", err)
	}

	if err := engine.Ping(); err != nil {
		glog.Errorf("database ping attempts failed: %#v", err)
	}
	// engine.ShowSQL(true)
	engine.SetLogger(&logger.Logger{})

	go pingDatabase(engine)

	engine.SetTableMapper(core.NewPrefixMapper(core.SnakeMapper{}, "t_"))

	engine.SetMaxIdleConns(config.MaxIdleConns)
	engine.SetMaxOpenConns(config.MaxOpenConns)
	return engine
}

func pingDatabase(engine *xorm.Engine) {
	timer := time.NewTicker(10 * time.Minute)
	for range timer.C {
		err := engine.Ping()
		if err != nil {
			glog.Errorf("database ping failed. retry in 10m. error: %#v", err)
		}
	}
}
