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

package store

import (
	"github.com/angao/recommender/pkg/apis/v1alpha1"
)

type Store interface {
	// Application CRUD
	CreateApplication(application *v1alpha1.Application) error

	GetApplication(name string) (*v1alpha1.Application, error)

	ListApplication() ([]*v1alpha1.Application, error)

	UpdateApplication(application *v1alpha1.Application) error

	DeleteApplication(application *v1alpha1.Application) error

	// RecommendResource CRUD
	GetApplicationResource(name string) (*v1alpha1.ApplicationResource, error)

	ListApplicationResource() ([]*v1alpha1.ApplicationResource, error)

	ListTimeframeApplicationResource(name string) ([]*v1alpha1.ApplicationResource, error)

	GetTimeframeApplicationResource(name, appName string) (*v1alpha1.ApplicationResource, error)

	AddOrUpdateContainerResource(resource []*v1alpha1.ContainerResource) error

	// Timeframe CRUD
	CreateTimeframe(frame *v1alpha1.Timeframe) error

	GetTimeframe(name string) (*v1alpha1.Timeframe, error)

	ListTimeframe() ([]*v1alpha1.Timeframe, error)

	UpdateTimeframe(frame *v1alpha1.Timeframe) error

	DeleteTimeframe(frame *v1alpha1.Timeframe) error

	UpdateTimeframes(timeframes []*v1alpha1.Timeframe) error
}
