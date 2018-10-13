/*
Copyright 2017 The Kubernetes Authors.

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

package model

import (
	"github.com/angao/recommender/pkg/apis/v1alpha1"
)

// ClusterState holds all runtime information about the cluster required for the
// VPA operations, i.e. configuration of resources (pods, containers,
// VPA objects), aggregated utilization of compute resources (CPU, memory) and
// events (container OOMs).
type ClusterState struct {
	// Applications
	Applications map[string]*v1alpha1.Application

	Vpas map[ApplicationID]*Vpa
}

// AggregateStateKey determines the set of containers for which the usage samples
// are kept aggregated in the model.
type AggregateStateKey interface {
	ApplicationName() string
	Name() string
	ContainerName() string
}

// AggregateContainerStatesMap is a map from AggregateStateKey to AggregateContainerState.
type aggregateContainerStatesMap map[AggregateStateKey]*AggregateContainerState

// NewClusterState returns a new ClusterState with no pods.
func NewClusterState() *ClusterState {
	return &ClusterState{
		Applications: make(map[string]*v1alpha1.Application),
		Vpas:         make(map[ApplicationID]*Vpa),
	}
}

func (cluster *ClusterState) AddApplication(application *v1alpha1.Application) {
	name := application.Name
	_, exists := cluster.Applications[name]
	if exists {
		cluster.DeleteApplication(name)
		exists = false
	}
	if !exists {
		cluster.Applications[name] = application
	}
}

func (cluster *ClusterState) DeleteApplication(name string) error {
	if _, exists := cluster.Applications[name]; !exists {
		return NewKeyError(name)
	}
	delete(cluster.Applications, name)
	return nil
}

// MakeAggregateStateKey returns the AggregateStateKey that should be used
// to aggregate usage samples from a container with the given name in a given pod.
func (cluster *ClusterState) MakeAggregateStateKey(app *ApplicationID, containerName string) AggregateStateKey {
	return aggregateStateKey{
		name:          app.Name,
		containerName: containerName,
	}
}

func (cluster *ClusterState) AddOrUpdateVPA(id ApplicationID) {
	_, exist := cluster.Vpas[id]
	if exist {
		cluster.DeleteVPA(id)
		exist = false
	}
	if !exist {
		vpa := NewVpa(id)
		cluster.Vpas[id] = vpa
	}
}

func (cluster *ClusterState) DeleteVPA(id ApplicationID) {
	if _, exist := cluster.Vpas[id]; exist {
		delete(cluster.Vpas, id)
	}
}

// Implementation of the AggregateStateKey interface. It can be used as a map key.
type aggregateStateKey struct {
	applicationName string
	name            string
	containerName   string
}

func NewAggregateStateKey(applicationContainer ApplicationContainer) AggregateStateKey {
	return aggregateStateKey{
		applicationName: applicationContainer.ContainerID.Name,
		containerName:   applicationContainer.ContainerName,
		name:            applicationContainer.Name,
	}
}

func (k aggregateStateKey) ApplicationName() string {
	return k.applicationName
}

// Labels returns the namespace for the aggregateStateKey.
func (k aggregateStateKey) Name() string {
	return k.name
}

// ContainerName returns the name of the container for the aggregateStateKey.
func (k aggregateStateKey) ContainerName() string {
	return k.containerName
}
