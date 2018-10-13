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

package model

// ContainerNameToAggregateStateMap maps a container name to AggregateContainerState
// that aggregates state of containers with that name.
type ContainerNameToAggregateStateMap map[string]*AggregateContainerState

// AggregateContainerState holds input signals aggregated from a set of containers.
// It can be used as an input to compute the recommendation.
// The CPU and memory distributions use decaying histograms by default
// (see NewAggregateContainerState()).
// Implements ContainerStateAggregator interface.
type AggregateContainerState struct {
	AggregateCPU               ResourceAmount
	AggregateMemory            ResourceAmount
	AggregateDiskReadIO        ResourceAmount
	AggregateDiskWriteIO       ResourceAmount
	AggregateNetworkReceiveIO  ResourceAmount
	AggregateNetworkTransmitIO ResourceAmount
}

// MergeContainerState merges two AggregateContainerStates.
func (a *AggregateContainerState) MergeContainerState(other *AggregateContainerState) {
	if a.AggregateCPU < other.AggregateCPU {
		a.AggregateCPU = other.AggregateCPU
	}
	if a.AggregateMemory < other.AggregateMemory {
		a.AggregateMemory = other.AggregateMemory
	}
	if a.AggregateDiskReadIO < other.AggregateDiskReadIO {
		a.AggregateDiskReadIO = other.AggregateDiskReadIO
	}
	if a.AggregateDiskWriteIO < other.AggregateDiskWriteIO {
		a.AggregateDiskWriteIO = other.AggregateDiskWriteIO
	}
	if a.AggregateNetworkReceiveIO < other.AggregateNetworkReceiveIO {
		a.AggregateNetworkReceiveIO = other.AggregateNetworkReceiveIO
	}
	if a.AggregateNetworkTransmitIO < other.AggregateNetworkTransmitIO {
		a.AggregateNetworkTransmitIO = other.AggregateNetworkTransmitIO
	}
}

// NewAggregateContainerState returns a new, empty AggregateContainerState.
func NewAggregateContainerState() *AggregateContainerState {
	return &AggregateContainerState{}
}

// AggregateStateByContainerName takes a set of AggregateContainerStates and merge them
// grouping by the container name. The result is a map from the container name to the aggregation
// from all input containers with the given name.
func AggregateStateByContainerName(aggregateContainerStateMap aggregateContainerStatesMap) ContainerNameToAggregateStateMap {
	containerNameToAggregateStateMap := make(ContainerNameToAggregateStateMap)
	for aggregationKey, aggregation := range aggregateContainerStateMap {
		containerName := aggregationKey.ContainerName()
		aggregateContainerState, isInitialized := containerNameToAggregateStateMap[containerName]
		if !isInitialized {
			aggregateContainerState = NewAggregateContainerState()
			containerNameToAggregateStateMap[containerName] = aggregateContainerState
		}
		aggregateContainerState.MergeContainerState(aggregation)
	}
	return containerNameToAggregateStateMap
}
