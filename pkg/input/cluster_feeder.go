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
	"strconv"
	"time"

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

	LoadTimeframes()

	LoadVPAs()

	LoadTimeframeVPAs()

	// LoadMetrics loads clusterState with current usage metrics of containers.
	LoadMetrics()

	LoadTimeframeMetrics()

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
		feeder.clusterState.AddApplication(application)
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

func (feeder *clusterStateFeeder) LoadTimeframes() {
	timeframes, err := feeder.store.ListTimeframe()
	if err != nil {
		glog.Errorf("Cannot list timeframes. Reason: %+v", err)
	} else {
		glog.V(3).Infof("Fetched %d timeframes.", len(timeframes))
	}
	for _, timeframe := range timeframes {
		if timeframe.Status != v1alpha1.StatusOn {
			continue
		}
		feeder.clusterState.AddTimeframe(timeframe)
	}

	for name := range feeder.clusterState.Timeframes {
		exist := false
		for _, timeframe := range timeframes {
			if timeframe.Status == v1alpha1.StatusOn && name == timeframe.Name {
				exist = true
			}
		}
		if !exist {
			feeder.clusterState.DeleteTimeframe(name)
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

func (feeder *clusterStateFeeder) LoadTimeframeVPAs() {
	timeframes := feeder.clusterState.Timeframes
	frameKey := make(map[string]bool)
	for name := range timeframes {
		frameKey[name] = true
		feeder.clusterState.AddOrUpdateTimeframeVPA(name)
	}

	for name := range feeder.clusterState.TimeframeVpas {
		if _, exist := frameKey[name]; !exist {
			feeder.clusterState.DeleteTimeframeVPAs(name)
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

func (feeder *clusterStateFeeder) getTimeframeMetrics(appName, historyLen, offset string) (map[model.AggregateStateKey]*model.AggregateContainerState, error) {
	aggregateContainerState, err := feeder.provider.GetTimeframeMetrics(appName, historyLen, offset)
	if err != nil {
		glog.Errorf("Cannot get %s timeframe history metrics", appName)
		return nil, err
	}
	return aggregateContainerState, nil
}

type QueryParam struct {
	TimeframeName string
	AppName       string
	HistoryLen    string
	Offset        string
}

func (feeder *clusterStateFeeder) LoadTimeframeMetrics() {
	queryParams := make([]QueryParam, 0)
	now := time.Now()
	for timeframeName, timeframe := range feeder.clusterState.Timeframes {
		timeframeVpa := feeder.clusterState.TimeframeVpas[timeframeName]
		for appID := range timeframeVpa {
			duration := timeframe.End.Sub(timeframe.Start).Minutes()
			if duration <= 0 {
				glog.Errorf("timeframe %s start after end", timeframeName)
				continue
			}
			offset := now.Sub(timeframe.End).Hours()
			if offset < 0 {
				glog.Errorf("timeframe %s end time after now", timeframeName)
				continue
			}
			queryParam := QueryParam{
				TimeframeName: timeframeName,
				AppName:       appID.Name,
				HistoryLen:    strconv.Itoa(int(duration)),
				Offset:        strconv.Itoa(int(offset)),
			}
			queryParams = append(queryParams, queryParam)
		}
	}
	load := func(i int) {
		queryParam := queryParams[i]
		aggregateContainerState, err := feeder.getTimeframeMetrics(queryParam.AppName, queryParam.HistoryLen, queryParam.Offset)
		if err != nil {
			return
		}
		timeframeVPA := feeder.clusterState.TimeframeVpas[queryParam.TimeframeName]
		for appID, vpa := range timeframeVPA {
			if appID.Name == queryParam.AppName {
				vpa.SetAggregationContainerState(aggregateContainerState)
			}
		}
	}
	work.Parallelize(8, len(queryParams), load)
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
	timeframes := make([]*v1alpha1.Timeframe, 0)
	for name, timeframe := range feeder.clusterState.Timeframes {
		timeframeVPA := feeder.clusterState.TimeframeVpas[name]
		for appID, vpa := range timeframeVPA {
			application := feeder.clusterState.Applications[appID.Name]
			for _, recommendResource := range vpa.Recommendation {
				containerResource := Convert(recommendResource)
				containerResource.ApplicationID = application.ID
				containerResource.TimeframeID = timeframe.ID
				containerResources = append(containerResources, containerResource)
			}
		}
		timeframe.Status = v1alpha1.StatusOff
		timeframes = append(timeframes, timeframe)
	}

	err := feeder.store.AddOrUpdateContainerResource(containerResources)
	if err != nil {
		glog.Errorf("add or update container resource error: %+v", err)
	} else {
		err := feeder.store.UpdateTimeframes(timeframes)
		if err != nil {
			glog.Errorf("update timeframe error: %+v", err)
		}
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
