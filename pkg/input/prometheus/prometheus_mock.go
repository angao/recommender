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
	"math/rand"
	"time"

	"github.com/angao/recommender/pkg/model"
)

func mockHistoryMetrics() map[model.AggregateStateKey]*model.AggregateContainerState {
	rand.Seed(time.Now().Unix())
	aggregateContainerStateMap := make(map[model.AggregateStateKey]*model.AggregateContainerState)

	aggregateStateKeys := mockAggregateStateKey()

	for _, aggregateStateKey := range aggregateStateKeys {
		aggregateContainerStateMap[aggregateStateKey] = &model.AggregateContainerState{
			AggregateCPU:               model.ResourceAmount(rand.Intn(1000)),
			AggregateMemory:            model.ResourceAmount(rand.Intn(1000)),
			AggregateDiskReadIO:        model.ResourceAmount(rand.Intn(1000)),
			AggregateDiskWriteIO:       model.ResourceAmount(rand.Intn(1000)),
			AggregateNetworkReceiveIO:  model.ResourceAmount(rand.Intn(1000)),
			AggregateNetworkTransmitIO: model.ResourceAmount(rand.Intn(1000)),
		}
	}

	return aggregateContainerStateMap
}

func mockAggregateStateKey() []model.AggregateStateKey {
	aggregateStateKeys := make([]model.AggregateStateKey, 0)

	aggregateStateKeys = append(aggregateStateKeys, model.NewAggregateStateKey(model.ApplicationContainer{
		ContainerID: model.ContainerID{
			ApplicationID: model.ApplicationID{Name: "app-3103947192"},
			ContainerName: "test",
		},
		Name: "test-1",
	}))

	aggregateStateKeys = append(aggregateStateKeys, model.NewAggregateStateKey(model.ApplicationContainer{
		ContainerID: model.ContainerID{
			ApplicationID: model.ApplicationID{Name: "app-3103947192"},
			ContainerName: "test",
		},
		Name: "test-2",
	}))

	aggregateStateKeys = append(aggregateStateKeys, model.NewAggregateStateKey(model.ApplicationContainer{
		ContainerID: model.ContainerID{
			ApplicationID: model.ApplicationID{Name: "app-3103947192"},
			ContainerName: "test1",
		},
		Name: "test1-1",
	}))
	aggregateStateKeys = append(aggregateStateKeys, model.NewAggregateStateKey(model.ApplicationContainer{
		ContainerID: model.ContainerID{
			ApplicationID: model.ApplicationID{Name: "app-3103947192"},
			ContainerName: "test1",
		},
		Name: "test1-2",
	}))
	aggregateStateKeys = append(aggregateStateKeys, model.NewAggregateStateKey(model.ApplicationContainer{
		ContainerID: model.ContainerID{
			ApplicationID: model.ApplicationID{Name: "app-3103947192"},
			ContainerName: "test1",
		},
		Name: "test1-3",
	}))
	return aggregateStateKeys
}
