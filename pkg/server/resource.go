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
	"fmt"
	"net/http"

	"github.com/angao/recommender/pkg/client/apis"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func (h *httpController) GetResource(c *gin.Context) {
	name := c.Param("name")
	glog.V(4).Infof("namespace: %s name: %s", name)
	resource, err := h.store.GetApplicationResource(name)
	if err != nil {
		glog.Errorf("Internal Server Error: %#v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Internal Server Error",
		})
		return
	}
	if resource == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": fmt.Sprintf("%s not found", name),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    resource,
	})
}

func (h *httpController) ListResource(c *gin.Context) {
	resources, err := h.store.ListApplicationResource()
	if err != nil {
		glog.Errorf("Internal Server Error: %#v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Internal Server Error",
		})
		return
	}

	recommendResources := make([]*apis.RecommendResource, 0)
	if len(resources) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
			"data":    recommendResources,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    resources,
	})
}
