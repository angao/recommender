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
	metricsFetcherInterval = flag.Duration("recommender-interval", 2*time.Hour, `How often metrics should be fetched`)
	globalConfig           = flag.String("config-file", "", `Specifies global config file. The config file type is yaml`)
)

func main() {
	util_flag.InitFlags()
	globalConfig, err := utils.Unmarshal(*globalConfig)
	if err != nil {
		glog.Fatalf("global globalConfig parse failed: %+v", err)
	}

	glog.V(1).Infof("Recommender %s", version.RecommenderVersion)
	recommender := routines.NewRecommender(globalConfig)

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
