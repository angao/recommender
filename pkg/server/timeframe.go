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
	"time"

	"github.com/angao/recommender/pkg/apis/v1alpha1"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

// Layout parse time
const Layout = "2006-01-02 15:04:05"

type TimeframeForm struct {
	Name        string `json:"name"`
	Start       string `form:"start"`
	End         string `form:"end"`
	Status      string `json:"status"`
	Description string `form:"description"`
}

func (h *httpController) GetTimeframe(c *gin.Context) {
	name := c.Param("name")
	timeframe, err := h.store.GetTimeframe(name)
	if err != nil {
		glog.Errorf("GetTimeframe Internal Server Error: %#v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Internal Server Error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    timeframe,
	})
}

func (h *httpController) CreateTimeframe(c *gin.Context) {
	timeframeForm := new(TimeframeForm)
	if err := c.ShouldBindJSON(timeframeForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	timeframe, err := ParseAndValidate(timeframeForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	frame, err := h.store.GetTimeframe(timeframe.Name)
	if err != nil {
		glog.Errorf("CreateApplication Internal Server Error: %#v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Internal Server Error",
		})
		return
	}
	if frame != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "timeframe has already exist",
		})
		return
	}

	err = h.store.CreateTimeframe(timeframe)
	if err != nil {
		glog.Errorf("CreateTimeframe Internal Server Error: %#v", err)
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

func (h *httpController) ListTimeframes(c *gin.Context) {
	timeframes, err := h.store.ListTimeframe()
	if err != nil {
		glog.Errorf("ListTimeframes Internal Server Error: %#v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Internal Server Error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    timeframes,
	})
}

func (h *httpController) DeleteTimeframe(c *gin.Context) {
	name := c.Param("name")
	timeframe, err := h.store.GetTimeframe(name)
	if err != nil {
		glog.Errorf("DeleteTimeframe Internal Server Error: %#v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Internal Server Error",
		})
		return
	}
	if timeframe == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": fmt.Sprintf("%s not found", name),
		})
		return
	}
	err = h.store.DeleteTimeframe(timeframe)
	if err != nil {
		glog.Errorf("DeleteTimeframe Internal Server Error: %#v", err)
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

func ParseAndValidate(form *TimeframeForm) (*v1alpha1.Timeframe, error) {
	if len(form.Name) == 0 || len(form.Start) == 0 || len(form.End) == 0 {
		return nil, errors.New("name, start or end field cannot be empty")
	}
	var status string
	if len(form.Status) == 0 {
		status = "off"
	}
	if !strings.Contains("on,off", status) {
		return nil, errors.New("status must be 'on' or 'off'")
	}
	start, err := time.ParseInLocation(Layout, form.Start, time.Local)
	if err != nil {
		return nil, err
	}
	end, err := time.ParseInLocation(Layout, form.End, time.Local)
	if err != nil {
		return nil, err
	}
	if start.After(end) {
		return nil, errors.New("start cannot be after end")
	}
	return &v1alpha1.Timeframe{
		Name:        form.Name,
		Start:       start,
		End:         end,
		Status:      status,
		Description: form.Description,
	}, nil
}
