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
	"fmt"

	"github.com/angao/recommender/pkg/apis/v1alpha1"
)

func (db *datastore) GetApplicationResource(name string) (*v1alpha1.ApplicationResource, error) {
	application := new(v1alpha1.Application)
	containerResources := make([]*v1alpha1.ContainerResource, 0)
	b, err := db.Engine.Where("name = ?", name).Limit(1).Get(application)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, nil
	}
	err = db.Engine.Where("application_id = ?", application.ID).And("timeframe_id = 0").Find(&containerResources)
	if err != nil {
		return nil, err
	}

	return &v1alpha1.ApplicationResource{
		ID:                application.ID,
		Name:              application.Name,
		ContainerResource: containerResources,
	}, nil
}

func (db *datastore) ListApplicationResource() ([]*v1alpha1.ApplicationResource, error) {
	applications := make([]*v1alpha1.Application, 0)
	containerResources := make([]*v1alpha1.ContainerResource, 0)
	err := db.Engine.Find(&applications)
	if err != nil {
		return nil, err
	}
	err = db.Engine.Where("timeframe_id = 0").Find(&containerResources)
	if err != nil {
		return nil, err
	}
	return Combine(applications, containerResources), nil
}

func (db *datastore) AddOrUpdateContainerResource(resources []*v1alpha1.ContainerResource) error {
	session := db.Engine.NewSession()
	defer session.Close()
	session.Begin()

	for _, resource := range resources {
		resourceCopy := new(v1alpha1.ContainerResource)
		has := false
		var err error
		if resource.TimeframeID == 0 {
			has, err = session.Where("application_id = ?", resource.ApplicationID).
				And("name = ?", resource.Name).And("timeframe_id = 0").Limit(1).Get(resourceCopy)
		} else {
			has, err = session.Where("application_id = ?", resource.ApplicationID).
				And("name = ?", resource.Name).And("timeframe_id = ?", resource.TimeframeID).Limit(1).Get(resourceCopy)
		}

		if err != nil {
			session.Rollback()
			return err
		}
		if has {
			compare(resource, resourceCopy)
			_, err := session.ID(resourceCopy.ID).Update(resource)
			if err != nil {
				session.Rollback()
				return err
			}
		} else {
			_, err = session.Insert(resource)
			if err != nil {
				session.Rollback()
				return err
			}
		}
	}
	return session.Commit()
}

func Combine(applications []*v1alpha1.Application, containerResources []*v1alpha1.ContainerResource) []*v1alpha1.ApplicationResource {
	applicationResources := make([]*v1alpha1.ApplicationResource, 0)

	for _, application := range applications {
		applicationResource := &v1alpha1.ApplicationResource{
			ID:                application.ID,
			Name:              application.Name,
			ContainerResource: make([]*v1alpha1.ContainerResource, 0),
		}
		applicationResources = append(applicationResources, applicationResource)
	}

	for _, appResource := range applicationResources {
		for _, resource := range containerResources {
			if appResource.ID == resource.ApplicationID {
				appResource.ContainerResource = append(appResource.ContainerResource, resource)
			}
		}
	}
	return applicationResources
}

func (db *datastore) ListTimeframeApplicationResource(name string) ([]*v1alpha1.ApplicationResource, error) {
	timeframe := new(v1alpha1.Timeframe)
	b, err := db.Engine.Where("name = ?", name).Limit(1).Get(timeframe)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, nil
	}
	applications := make([]*v1alpha1.Application, 0)
	containerResources := make([]*v1alpha1.ContainerResource, 0)
	err = db.Engine.Find(&applications)
	if err != nil {
		return nil, err
	}
	err = db.Engine.Where("timeframe_id = ?", timeframe.ID).Find(&containerResources)
	if err != nil {
		return nil, err
	}
	return Combine(applications, containerResources), nil
}

func (db *datastore) GetTimeframeApplicationResource(name, appName string) (*v1alpha1.ApplicationResource, error) {
	timeframe := new(v1alpha1.Timeframe)
	b, err := db.Engine.Where("name = ?", name).Limit(1).Get(timeframe)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, nil
	}
	application := new(v1alpha1.Application)
	b, err = db.Engine.Where("name = ?", appName).Limit(1).Get(application)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, nil
	}
	containerResources := make([]*v1alpha1.ContainerResource, 0)
	err = db.Engine.Where("application_id = ?", application.ID).And("timeframe_id = ?", timeframe.ID).Find(&containerResources)
	if err != nil {
		return nil, err
	}

	return &v1alpha1.ApplicationResource{
		ID:                application.ID,
		Name:              application.Name,
		ContainerResource: containerResources,
	}, nil
}

func compare(r1, r2 *v1alpha1.ContainerResource) {
	if r1.CPULimit < r2.CPULimit {
		r1.CPULimit = r2.CPULimit
	}
	if r1.MemoryLimit < r2.MemoryLimit {
		r1.MemoryLimit = r2.MemoryLimit
	}
	if r1.DiskReadIOLimit < r2.DiskReadIOLimit {
		r1.DiskReadIOLimit = r2.DiskReadIOLimit
	}
	if r1.DiskWriteIOLimit < r2.DiskWriteIOLimit {
		r1.DiskWriteIOLimit = r2.DiskWriteIOLimit
	}
	if r1.NetworkReceiveIOLimit < r2.NetworkReceiveIOLimit {
		r1.NetworkReceiveIOLimit = r2.NetworkReceiveIOLimit
	}
	if r1.NetworkTransmitIOLimit < r2.NetworkTransmitIOLimit {
		r1.NetworkTransmitIOLimit = r2.NetworkTransmitIOLimit
	}
}

func (db *datastore) CreateContainerResource(resource *v1alpha1.ContainerResource) error {
	_, err := db.Engine.Insert(resource)
	return err
}

func (db *datastore) DeleteApplicationResource(name string) error {
	session := db.Engine.NewSession()
	defer session.Close()

	session.Begin()

	application := new(v1alpha1.Application)
	has, err := session.Unscoped().Where("name = ?", name).Get(application)
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("%s not found", name)
	}
	containerResources := make([]*v1alpha1.ContainerResource, 0)
	err = session.Where("application_id = ?", application.ID).And("timeframe_id = 0").Find(&containerResources)
	if err != nil {
		return err
	}
	for _, resource := range containerResources {
		_, err := session.ID(resource.ID).Delete(resource)
		if err != nil {
			session.Rollback()
			return err
		}
	}
	return session.Commit()
}

func (db *datastore) DeleteTimeframeResource(name string) error {
	session := db.Engine.NewSession()
	defer session.Close()

	session.Begin()

	timeframe := new(v1alpha1.Timeframe)
	has, err := session.Unscoped().Where("name = ?", name).Get(timeframe)
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("%s not found", name)
	}
	containerResources := make([]*v1alpha1.ContainerResource, 0)
	err = session.Where("timeframe_id = ?", timeframe.ID).Find(&containerResources)
	if err != nil {
		return err
	}
	for _, resource := range containerResources {
		_, err := session.ID(resource.ID).Delete(resource)
		if err != nil {
			session.Rollback()
			return err
		}
	}
	return session.Commit()
}
