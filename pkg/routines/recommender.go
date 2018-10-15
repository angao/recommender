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

package routines

import (
	"fmt"
	"net/http"
	"time"

	"github.com/angao/recommender/pkg/client"
	"github.com/angao/recommender/pkg/input"
	"github.com/angao/recommender/pkg/logic"
	"github.com/angao/recommender/pkg/model"
	"github.com/angao/recommender/pkg/server"
	"github.com/angao/recommender/pkg/store/database"
	"github.com/angao/recommender/pkg/utils"

	"github.com/golang/glog"
)

const (
	// Driver defines which database to use.
	Driver = "mysql"
)

// Recommender recommend resources for certain containers, based on utilization periodically got from metrics api.
type Recommender interface {
	// RunOnce performs one iteration of recommender duties followed by update of recommendations in VPA objects.
	RunOnce()
	// GetClusterState returns ClusterState used by Recommender
	GetClusterState() *model.ClusterState
	// GetClusterStateFeeder returns ClusterStateFeeder used by Recommender
	GetClusterStateFeeder() input.ClusterStateFeeder
}

type recommender struct {
	clusterState        *model.ClusterState
	clusterStateFeeder  input.ClusterStateFeeder
	resourceRecommender logic.ResourceRecommender
}

func (r *recommender) GetClusterState() *model.ClusterState {
	return r.clusterState
}

func (r *recommender) GetClusterStateFeeder() input.ClusterStateFeeder {
	return r.clusterStateFeeder
}

func (r *recommender) RunOnce() {
	start := time.Now()
	glog.V(3).Infof("Recommender Run")
	defer glog.V(3).Infof("Recommender Finished: %v", time.Since(start))
	r.clusterStateFeeder.LoadApplications()
	r.clusterStateFeeder.LoadTimeframes()
	r.clusterStateFeeder.LoadVPAs()
	r.clusterStateFeeder.LoadTimeframeVPAs()
	r.clusterStateFeeder.LoadMetrics()
	r.clusterStateFeeder.LoadTimeframeMetrics()
	r.updateVPAs()
	r.clusterStateFeeder.UpdateResources()
}

func (r *recommender) updateVPAs() {
	for _, vpa := range r.clusterState.Vpas {
		resources := r.resourceRecommender.GetRecommendedResources(vpa)
		vpa.Recommendation = resources
	}
	for _, timeframeVPA := range r.clusterState.TimeframeVpas {
		for _, vpa := range timeframeVPA {
			resources := r.resourceRecommender.GetRecommendedResources(vpa)
			vpa.Recommendation = resources
		}
	}
}

// NewRecommender creates a new recommender instance,
// which can be run in order to provide continuous resource recommendations for containers.
// It requires cluster configuration object and duration between recommender intervals.
func NewRecommender(globalConfig *utils.GlobalConfig) Recommender {
	store := datastore.New(Driver, globalConfig.DatabaseConfig)
	clusterState := model.NewClusterState()
	recommender := &recommender{
		clusterState:        clusterState,
		clusterStateFeeder:  input.NewClusterStateFeeder(store, globalConfig, clusterState),
		resourceRecommender: logic.CreateResourceRecommender(),
	}
	glog.V(3).Infof("New Recommender created %+v", recommender)

	s := server.NewController(store)
	startHTTPServer(s, globalConfig.ExtraConfig.APIPort)

	return recommender
}

func startHTTPServer(ctrl server.Controller, port int) {
	handler := client.Load(ctrl)
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}
	go func() {
		if err := s.ListenAndServe(); err != nil {
			glog.Fatalf("http server start failed: %v", err)
		}
	}()
	glog.V(2).Infof("HTTP Server started and listening on %d.", port)
}
