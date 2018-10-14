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

package server

import (
	"github.com/angao/recommender/pkg/store"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	CreateApplication(c *gin.Context)
	GetApplication(c *gin.Context)
	ListApplications(c *gin.Context)
	DeleteApplication(c *gin.Context)

	GetResource(c *gin.Context)
	CreateResource(c *gin.Context)
	ListResource(c *gin.Context)
	ListTimeframeResource(c *gin.Context)
	GetTimeframeResource(c *gin.Context)

	CreateTimeframe(c *gin.Context)
	GetTimeframe(c *gin.Context)
	ListTimeframes(c *gin.Context)
	DeleteTimeframe(c *gin.Context)
}

type httpController struct {
	store store.Store
}

func NewController(store store.Store) Controller {
	return &httpController{
		store: store,
	}
}
