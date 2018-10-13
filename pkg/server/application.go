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
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/angao/recommender/pkg/apis/v1alpha1"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func (h *httpController) GetApplication(c *gin.Context) {
	name := c.Param("name")
	application, err := h.store.GetApplication(name)
	if err != nil {
		glog.Errorf("GetApplication Internal Server Error: %#v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Internal Server Error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    application,
	})
}

func (h *httpController) CreateApplication(c *gin.Context) {
	application := new(v1alpha1.Application)
	if err := c.ShouldBindJSON(application); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	if err := ValidateApplication(application); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	app, err := h.store.GetApplication(application.Name)
	if err != nil {
		glog.Errorf("CreateApplication Internal Server Error: %#v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Internal Server Error",
		})
		return
	}
	if app != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "application has already exist",
		})
		return
	}

	err = h.store.CreateApplication(application)
	if err != nil {
		glog.Errorf("Internal Server Error: %#v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Internal Server Error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
	})
}

func (h *httpController) ListApplications(c *gin.Context) {
	applications, err := h.store.ListApplication()
	if err != nil {
		glog.Errorf("ListApplications Internal Server Error: %#v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Internal Server Error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    applications,
	})
}

func (h *httpController) DeleteApplication(c *gin.Context) {
	name := c.Param("name")
	application, err := h.store.GetApplication(name)
	if err != nil {
		glog.Errorf("DeleteApplication Internal Server Error: %#v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Internal Server Error",
		})
		return
	}
	if application == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": fmt.Sprintf("%s not found", name),
		})
		return
	}
	err = h.store.DeleteApplication(application)
	if err != nil {
		glog.Errorf("DeleteApplication Internal Server Error: %#v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Internal Server Error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
	})
}

func ValidateApplication(application *v1alpha1.Application) error {
	if len(strings.TrimSpace(application.Name)) == 0 {
		return errors.New("name cannot be empty")
	}
	return nil
}
