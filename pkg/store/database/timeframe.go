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

func (db *datastore) CreateTimeframe(frame *v1alpha1.Timeframe) error {
	_, err := db.Engine.Insert(frame)
	return err
}

func (db *datastore) GetTimeframe(name string) (*v1alpha1.Timeframe, error) {
	timeframe := new(v1alpha1.Timeframe)
	b, err := db.Engine.Where("name = ?", name).Limit(1).Get(timeframe)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, nil
	}
	return timeframe, nil
}

func (db *datastore) ListTimeframe() ([]*v1alpha1.Timeframe, error) {
	frames := make([]*v1alpha1.Timeframe, 0)
	err := db.Engine.Find(&frames)
	return frames, err
}

func (db *datastore) UpdateTimeframe(frame *v1alpha1.Timeframe) error {
	_, err := db.Engine.Id(frame.ID).Update(frame)
	return err
}

func (db *datastore) DeleteTimeframe(frame *v1alpha1.Timeframe) error {
	_, err := db.Engine.Id(frame.ID).Delete(frame)
	return err
}
