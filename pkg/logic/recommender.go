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

package logic

import (
	"github.com/angao/recommender/pkg/model"
)

// PodResourceRecommender computes resource recommendation for a Vpa object.
type ResourceRecommender interface {
	GetRecommendedResources(vpa *model.Vpa) []model.RecommendedContainerResources
}

type resourceRecommender struct {
}

// Returns recommended resources for a given Vpa object.
func (r *resourceRecommender) GetRecommendedResources(vpa *model.Vpa) []model.RecommendedContainerResources {
	containerNameToAggregateStateMap := vpa.AggregateStateByContainerName()
	recommendedContainerResources := make([]model.RecommendedContainerResources, 0)

	for containerName, aggregatedContainerState := range containerNameToAggregateStateMap {
		containerResource := model.RecommendedContainerResources{
			ContainerName:          containerName,
			CPULimit:               aggregatedContainerState.AggregateCPU,
			MemoryLimit:            aggregatedContainerState.AggregateMemory,
			DiskReadIOLimit:        aggregatedContainerState.AggregateDiskReadIO,
			DiskWriteIOLimit:       aggregatedContainerState.AggregateDiskWriteIO,
			NetworkReceiveIOLimit:  aggregatedContainerState.AggregateNetworkReceiveIO,
			NetworkTransmitIOLimit: aggregatedContainerState.AggregateNetworkTransmitIO,
		}
		recommendedContainerResources = append(recommendedContainerResources, containerResource)
	}
	return recommendedContainerResources
}

// CreatePodResourceRecommender returns the primary recommender.
func CreateResourceRecommender() ResourceRecommender {
	return &resourceRecommender{}
}
