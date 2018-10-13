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

package main

import (
	"flag"
	"time"

	"github.com/angao/recommender/pkg/routines"
	"github.com/angao/recommender/pkg/utils"
	util_flag "github.com/angao/recommender/pkg/utils/flag"
	"github.com/angao/recommender/version"

	"github.com/golang/glog"
)

var (
	metricsFetcherInterval = flag.Duration("recommender-interval", 1*time.Minute, `How often metrics should be fetched`)
	prometheusAddress      = flag.String("prometheus-address", "", `Where to reach for Prometheus metrics`)
	apiserverPort          = flag.Int("apiserver-port", 9098, `Specifies the http apiserver port`)
	databaseConfig         = flag.String("db-config-file", "", `Where to reach for MySQL. The db config file type is yaml`)
)

func main() {
	util_flag.InitFlags()
	glog.V(1).Infof("Recommender %s and listen on %d", version.RecommenderVersion, *apiserverPort)

	config, err := utils.Unmarshal(*databaseConfig)
	if err != nil {
		glog.Fatalf("databaseConfig parse failed: %+v", err)
	}

	recommender := routines.NewRecommender(*prometheusAddress, config, *apiserverPort)

	recommender.RunOnce()
	for {
		select {
		case <-time.After(*metricsFetcherInterval):
			{
				recommender.RunOnce()
			}
		}
	}
}
