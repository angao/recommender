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
	"github.com/angao/recommender/pkg/apis/v1alpha1"
)

func (db *datastore) CreateApplication(application *v1alpha1.Application) error {
	_, err := db.Engine.Insert(application)
	return err
}

func (db *datastore) GetApplication(name string) (*v1alpha1.Application, error) {
	application := new(v1alpha1.Application)
	b, err := db.Engine.Where("name = ?", name).Limit(1).Get(application)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, nil
	}
	return application, err
}

func (db *datastore) ListApplication() ([]*v1alpha1.Application, error) {
	applications := make([]*v1alpha1.Application, 0)
	err := db.Engine.Find(&applications)
	return applications, err
}

func (db *datastore) UpdateApplication(application *v1alpha1.Application) error {
	_, err := db.Engine.ID(application.ID).Update(application)
	return err
}

func (db *datastore) DeleteApplication(application *v1alpha1.Application) error {
	_, err := db.Engine.Id(application.ID).Delete(application)
	return err
}
