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
	"time"

	"github.com/angao/recommender/pkg/apis/v1alpha1"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

// Layout parse time
const Layout = "2006-01-02 15:04:05"

type TimeframeForm struct {
	ID          int64  `json:"id"`
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
			"message": err.Error(),
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

	timeframe, err := ParseAndValidate(timeframeForm, "add")
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
			"message": err.Error(),
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
			"message": err.Error(),
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
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    timeframes,
	})
}

func (h *httpController) UpdateTimeframe(c *gin.Context) {
	timeframeForm := new(TimeframeForm)
	if err := c.ShouldBindJSON(timeframeForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	timeframe, err := ParseAndValidate(timeframeForm, "update")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	timeframeCopy, err := h.store.GetTimeframe(timeframe.Name)
	if err != nil {
		glog.Errorf("UpdateTimeframe Internal Server Error: %#v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	if timeframeCopy != nil && timeframeCopy.ID != timeframe.ID {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "name is already exist",
		})
		return
	}
	err = h.store.UpdateTimeframe(timeframe)
	if err != nil {
		glog.Errorf("UpdateTimeframe Internal Server Error: %#v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
	})
}

func (h *httpController) DeleteTimeframe(c *gin.Context) {
	name := c.Param("name")
	timeframe, err := h.store.GetTimeframe(name)
	if err != nil {
		glog.Errorf("DeleteTimeframe Internal Server Error: %#v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	if timeframe == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    404,
			"message": fmt.Sprintf("%s not found", name),
		})
		return
	}
	err = h.store.DeleteTimeframe(timeframe)
	if err != nil {
		glog.Errorf("DeleteTimeframe Internal Server Error: %#v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
	})
}

func ParseAndValidate(form *TimeframeForm, flag string) (*v1alpha1.Timeframe, error) {
	if flag == "add" {
		if len(form.Name) == 0 || len(form.Start) == 0 || len(form.End) == 0 {
			return nil, errors.New("name, start or end field cannot be empty")
		}
		status := form.Status
		if len(form.Status) == 0 {
			status = "off"
		}
		if status != "on" && status != "off" {
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
	if flag == "update" {
		timeframe := &v1alpha1.Timeframe{}
		if form.ID == 0 {
			return nil, errors.New("update: id cannot be empty")
		}
		timeframe.ID = form.ID
		if len(form.Name) != 0 {
			timeframe.Name = form.Name
		}
		if len(form.Start) != 0 {
			start, err := time.ParseInLocation(Layout, form.Start, time.Local)
			if err != nil {
				return nil, err
			}
			timeframe.Start = start
		}
		if len(form.End) != 0 {
			end, err := time.ParseInLocation(Layout, form.End, time.Local)
			if err != nil {
				return nil, err
			}
			timeframe.End = end
		}
		if len(form.Status) != 0 {
			if form.Status != "on" && form.Status != "off" {
				return nil, errors.New("status must be 'on' or 'off'")
			}
			timeframe.Status = form.Status
		}
		return timeframe, nil
	}
	return nil, nil
}
