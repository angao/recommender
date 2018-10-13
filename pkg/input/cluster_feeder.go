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

package input

import (
	"github.com/angao/recommender/pkg/apis/v1alpha1"
	"github.com/angao/recommender/pkg/utils/work"
	"github.com/golang/glog"

	"github.com/angao/recommender/pkg/input/prometheus"
	"github.com/angao/recommender/pkg/model"
	"github.com/angao/recommender/pkg/store"
)

// ClusterStateFeeder can update state of ClusterState object.
type ClusterStateFeeder interface {
	// LoadApplications loads all applications into clusterState
	LoadApplications()

	LoadVPAs()

	// LoadMetrics loads clusterState with current usage metrics of containers.
	LoadMetrics()

	UpdateResources()
}

// NewClusterStateFeeder creates new ClusterStateFeeder with internal data providers, based on kube client config and a historyProvider.
func NewClusterStateFeeder(store store.Store, prometheusAddress string, clusterState *model.ClusterState) ClusterStateFeeder {
	return &clusterStateFeeder{
		store:        store,
		clusterState: clusterState,
		provider:     prometheus.NewPrometheusHistoryProvider(prometheusAddress),
	}
}

type clusterStateFeeder struct {
	store        store.Store
	clusterState *model.ClusterState
	provider     prometheus.Provider
}

func (feeder *clusterStateFeeder) LoadApplications() {
	applications, err := feeder.store.ListApplication()
	if err != nil {
		glog.Errorf("Cannot list applications. Reason: %+v", err)
	} else {
		glog.V(3).Infof("Fetched %d applications.", len(applications))
	}

	for _, application := range applications {
		if _, exist := feeder.clusterState.Applications[application.Name]; !exist {
			feeder.clusterState.AddApplication(application)
		}
	}

	for name := range feeder.clusterState.Applications {
		exist := false
		for _, application := range applications {
			if name == application.Name {
				exist = true
			}
		}
		if !exist {
			feeder.clusterState.DeleteApplication(name)
		}
	}
}

func (feeder *clusterStateFeeder) LoadVPAs() {
	applications := feeder.clusterState.Applications
	applicationKey := make(map[model.ApplicationID]bool)
	for name := range applications {
		applicationID := model.ApplicationID{Name: name}
		applicationKey[applicationID] = true
		feeder.clusterState.AddOrUpdateVPA(applicationID)
	}

	for vpaID := range feeder.clusterState.Vpas {
		if _, exists := applicationKey[vpaID]; !exists {
			feeder.clusterState.DeleteVPA(vpaID)
		}
	}
}

func (feeder *clusterStateFeeder) loadHistoryMetrics(name string) {
	aggregateContainerState, err := feeder.provider.GetHistoryMetrics(name)
	if err != nil {
		glog.Errorf("Cannot get %s history metrics", name)
	}
	applicationID := model.ApplicationID{Name: name}
	for vpaID, vpa := range feeder.clusterState.Vpas {
		if vpaID == applicationID {
			vpa.SetAggregationContainerState(aggregateContainerState)
		}
	}
}

func (feeder *clusterStateFeeder) LoadMetrics() {
	applications := make([]string, 0)
	for name := range feeder.clusterState.Applications {
		applications = append(applications, name)
	}

	load := func(i int) {
		name := applications[i]
		feeder.loadHistoryMetrics(name)
	}

	work.Parallelize(8, len(applications), load)
}

func (feeder *clusterStateFeeder) UpdateResources() {
	applications := feeder.clusterState.Applications
	containerResources := make([]*v1alpha1.ContainerResource, 0)
	for _, application := range applications {
		applicationID := model.ApplicationID{Name: application.Name}
		vpa := feeder.clusterState.Vpas[applicationID]
		for _, recommendResource := range vpa.Recommendation {
			containerResource := Convert(recommendResource)
			containerResource.ApplicationID = application.ID
			containerResources = append(containerResources, containerResource)
		}
	}
	err := feeder.store.AddOrUpdateContainerResource(containerResources)
	if err != nil {
		glog.Errorf("add or update container resource error: %+v", err)
	}
}

func Convert(recommendResource model.RecommendedContainerResources) *v1alpha1.ContainerResource {
	return &v1alpha1.ContainerResource{
		CPULimit:               int64(recommendResource.CPULimit),
		MemoryLimit:            int64(recommendResource.MemoryLimit),
		DiskReadIOLimit:        int64(recommendResource.DiskReadIOLimit),
		DiskWriteIOLimit:       int64(recommendResource.DiskWriteIOLimit),
		NetworkReceiveIOLimit:  int64(recommendResource.NetworkReceiveIOLimit),
		NetworkTransmitIOLimit: int64(recommendResource.NetworkTransmitIOLimit),
	}
}
