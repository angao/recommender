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

package client

import (
	"net/http"

	"github.com/angao/recommender/pkg/client/middleware/logger"

	"github.com/angao/recommender/pkg/server"

	"github.com/angao/recommender/pkg/client/middleware/header"

	"github.com/angao/recommender/version"
	"github.com/gin-gonic/gin"
)

// Load defines HTTP API
func Load(s server.Controller, middleware ...gin.HandlerFunc) http.Handler {
	e := gin.New()
	gin.SetMode(gin.ReleaseMode)

	e.Use(gin.Recovery())

	e.Use(header.NoCache)
	e.Use(header.Options)

	e.Use(middleware...)
	e.Use(logger.Logger())

	e.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "request not found",
		})
	})

	app := e.Group("/api/v1")
	{
		app.GET("/application/:name", s.GetApplication)
		app.GET("/applications", s.ListApplications)
		app.POST("/application", s.CreateApplication)
		app.DELETE("/application/:name", s.DeleteApplication)

		app.GET("/resource/:name", s.GetResource)
		app.DELETE("/resource/:name", s.DeleteResource)
		app.GET("/resources", s.ListResource)
		app.GET("/resources/timeframe/:name", s.ListTimeframeResource)
		app.DELETE("/resources/timeframe/:name", s.DeleteTimeframeResource)
		app.GET("/resources/timeframe/:name/:appName", s.GetTimeframeResource)
		// app.POST("/resource", s.CreateResource)

		app.POST("/timeframe", s.CreateTimeframe)
		app.GET("/timeframes", s.ListTimeframes)
		app.GET("/timeframe/:name", s.GetTimeframe)
		app.PUT("/timeframe", s.UpdateTimeframe)
		app.DELETE("/timeframe/:name", s.DeleteTimeframe)
	}

	e.GET("/version", versionCtrl)
	return e
}

func versionCtrl(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":    "Recommender",
		"version": version.RecommenderVersion,
	})
}
