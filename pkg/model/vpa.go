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

// Vpa (Vertical Pod Autoscaler) object is responsible for vertical scaling of
// Pods matching a given label selector.
type Vpa struct {
	ID             ApplicationID
	Recommendation []RecommendedContainerResources
	// All container aggregations that contribute to this VPA.
	aggregateContainerStates aggregateContainerStatesMap
}

type RecommendedContainerResources struct {
	// Name of the container.
	ContainerName          string
	CPULimit               ResourceAmount
	MemoryLimit            ResourceAmount
	DiskReadIOLimit        ResourceAmount
	DiskWriteIOLimit       ResourceAmount
	NetworkReceiveIOLimit  ResourceAmount
	NetworkTransmitIOLimit ResourceAmount
}

// NewVpa returns a new Vpa with a given ID and pod selector. Doesn't set the
// links to the matched aggregations.
func NewVpa(id ApplicationID) *Vpa {
	vpa := &Vpa{
		ID:                       id,
		Recommendation:           make([]RecommendedContainerResources, 0),
		aggregateContainerStates: make(aggregateContainerStatesMap),
	}
	return vpa
}

func (vpa *Vpa) SetAggregationContainerState(aggregateContainerStates aggregateContainerStatesMap) {
	vpa.aggregateContainerStates = aggregateContainerStates
}

// DeleteAggregation deletes aggregation used by this container
func (vpa *Vpa) DeleteAggregation(aggregationKey AggregateStateKey) {
	delete(vpa.aggregateContainerStates, aggregationKey)
}

// AggregateStateByContainerName returns a map from container name to the aggregated state
// of all containers with that name, belonging to pods matched by the VPA.
func (vpa *Vpa) AggregateStateByContainerName() ContainerNameToAggregateStateMap {
	return AggregateStateByContainerName(vpa.aggregateContainerStates)
}
