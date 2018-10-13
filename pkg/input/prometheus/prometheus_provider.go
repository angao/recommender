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

package prometheus

import (
	"fmt"
	"net/http"

	"github.com/angao/recommender/pkg/model"
)

const (
	historyLength = "90d"
)

// Provider gives metrics data of all pods in a cluster.
// Consider refactoring to passing ClusterState and create history provider working with checkpoints.
type Provider interface {
	// GetClusterHistory(string) (map[model.PodID]*PodHistory, error)
	GetHistoryMetrics(name string) (map[model.AggregateStateKey]*model.AggregateContainerState, error)
}

type prometheusProvider struct {
	prometheusClient PrometheusClient
}

// NewPrometheusHistoryProvider contructs a history provider that gets data from Prometheus.
func NewPrometheusHistoryProvider(prometheusAddress string) Provider {
	return &prometheusProvider{
		prometheusClient: NewPrometheusClient(&http.Client{}, prometheusAddress),
	}
}

func getApplicationContainerFromLabels(labels map[string]string) (*model.ApplicationContainer, error) {
	applicationName, ok := labels["application_name"]
	if !ok {
		return nil, fmt.Errorf("no pod_name label")
	}
	containerName, ok := labels["container_name"]
	if !ok {
		return nil, fmt.Errorf("no name label on container data")
	}
	name, ok := labels["name"]
	if !ok {
		return nil, fmt.Errorf("no name label on name data")
	}
	return &model.ApplicationContainer{
		ContainerID: model.ContainerID{
			ApplicationID: model.ApplicationID{Name: applicationName},
			ContainerName: containerName,
		},
		Name: name,
	}, nil
}

func (p *prometheusProvider) readResource(res map[model.AggregateStateKey]*model.AggregateContainerState, query string, resource model.ResourceName) error {
	tss, err := p.prometheusClient.GetTimeseries(query)
	if err != nil {
		return fmt.Errorf("cannot get timeseries for %v: %v", resource, err)
	}
	for _, ts := range tss {
		applicationContainer, err := getApplicationContainerFromLabels(ts.Labels)
		if err != nil {
			return fmt.Errorf("cannot get application container from labels: %v", err)
		}
		aggregateContainerKey := model.NewAggregateStateKey(*applicationContainer)
		aggregateContainerState, ok := res[aggregateContainerKey]
		if !ok {
			aggregateContainerState = model.NewAggregateContainerState()
			res[aggregateContainerKey] = aggregateContainerState
		}
		value := ts.Sample.Value
		switch resource {
		case model.ResourceCPU:
			aggregateContainerState.AggregateCPU = model.CPUAmountFromCores(value)
		case model.ResourceMemory:
			aggregateContainerState.AggregateMemory = model.MemoryAmountFromBytes(value)
		case model.ResourceDiskReadIO:
			aggregateContainerState.AggregateDiskReadIO = model.ResourceAmountFromFloat(value)
		case model.ResourceDiskWriteIO:
			aggregateContainerState.AggregateDiskWriteIO = model.ResourceAmountFromFloat(value)
		case model.ResourceNetworkReceiveIO:
			aggregateContainerState.AggregateNetworkReceiveIO = model.ResourceAmountFromFloat(value)
		case model.ResourceNetworkTransmitIO:
			aggregateContainerState.AggregateNetworkTransmitIO = model.ResourceAmountFromFloat(value)
		}
		res[aggregateContainerKey] = aggregateContainerState
	}
	return nil
}

func (p *prometheusProvider) GetHistoryMetrics(name string) (map[model.AggregateStateKey]*model.AggregateContainerState, error) {
	res := make(map[model.AggregateStateKey]*model.AggregateContainerState)
	podSelector := fmt.Sprintf(`job="kubernetes-cadvisor",pod_name=~"^.*$",container_name!="POD",image!="",name=~"^k8s_.*",application_name="%s"`, name)
	err := p.readResource(res, fmt.Sprintf("max_over_time(container_cpu_usage_seconds_total:rate:1m{%s}[%s])", podSelector, historyLength), model.ResourceCPU)
	if err != nil {
		return nil, fmt.Errorf("cannot get cpu usage history: %v", err)
	}
	err = p.readResource(res, fmt.Sprintf("max_over_time(container_memory_usage_bytes{%s}[%s])", podSelector, historyLength), model.ResourceMemory)
	if err != nil {
		return nil, fmt.Errorf("cannot get memory usage history: %v", err)
	}
	err = p.readResource(res, fmt.Sprintf("max_over_time(container_fs_reads_total{%s}[%s])", podSelector, historyLength), model.ResourceDiskReadIO)
	if err != nil {
		return nil, fmt.Errorf("cannot get disk read io history: %v", err)
	}
	err = p.readResource(res, fmt.Sprintf("max_over_time(container_fs_writes_total{%s}[%s])", podSelector, historyLength), model.ResourceDiskWriteIO)
	if err != nil {
		return nil, fmt.Errorf("cannot get disk write io history: %v", err)
	}
	err = p.readResource(res, fmt.Sprintf("max_over_time(container_network_receive_packets_total{%s}[%s])", podSelector, historyLength), model.ResourceNetworkReceiveIO)
	if err != nil {
		return nil, fmt.Errorf("cannot get network receive io history: %v", err)
	}
	err = p.readResource(res, fmt.Sprintf("max_over_time(container_network_transmit_packets_total{%s}[%s])", podSelector, historyLength), model.ResourceNetworkTransmitIO)
	if err != nil {
		return nil, fmt.Errorf("cannot get network transmit io history: %v", err)
	}
	return res, nil
}
