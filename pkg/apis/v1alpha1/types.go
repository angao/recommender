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

package v1alpha1

import (
	"time"
)

type ApplicationStore interface {
	// Application CRUD
	CreateApplication(application *Application) error

	GetApplication(name string) (*Application, error)

	ListApplication() ([]*Application, error)

	UpdateApplication(application *Application) error

	DeleteApplication(application *Application) error
}

// RecommendResourceStore defines store RecommendResource
type RecommendResourceStore interface {
	// RecommendResource CRUD
	GetApplicationResource(name string) (*ApplicationResource, error)

	ListApplicationResource() ([]*ApplicationResource, error)

	AddOrUpdateContainerResource(resource []*ContainerResource) error

	ListTimeframeApplicationResource(name string) ([]*ApplicationResource, error)

	GetTimeframeApplicationResource(name, appName string) (*ApplicationResource, error)
}

type TimeframeStore interface {
	// Timeframe CRUD
	CreateTimeframe(frame *Timeframe) error

	GetTimeframe(name string) (*Timeframe, error)

	ListTimeframe() ([]*Timeframe, error)

	UpdateTimeframe(frame *Timeframe) error

	UpdateTimeframes(frame []*Timeframe) error

	DeleteTimeframe(frame *Timeframe) error
}

// Application defines application info.
type Application struct {
	ID      int64     `json:"id"      form:"id"         xorm:"pk autoincr 'id'"`
	Name    string    `json:"name"    form:"name"       xorm:"name"`
	Created time.Time `json:"created"                   xorm:"created"`
	Updated time.Time `json:"updated"                   xorm:"updated"`
}

// ContainerResource defines container of application resource
type ContainerResource struct {
	ID                     int64     `json:"id"                             xorm:"pk autoincr 'id'"`
	Name                   string    `json:"name"                           xorm:"name"`
	ApplicationID          int64     `json:"application_id"                 xorm:"application_id"`
	TimeframeID            int64     `json:"timeframe_id"                   xorm:"timeframe_id"`
	CPULimit               int64     `json:"cpu_limit"                      xorm:"cpu_limit"`
	MemoryLimit            int64     `json:"memory_limit"                   xorm:"memory_limit"`
	DiskReadIOLimit        int64     `json:"disk_read_io_limit"             xorm:"disk_read_io_limit"`
	DiskWriteIOLimit       int64     `json:"disk_write_io_limit"            xorm:"disk_write_io_limit"`
	NetworkReceiveIOLimit  int64     `json:"network_receive_io_limit"       xorm:"network_receive_io_limit"`
	NetworkTransmitIOLimit int64     `json:"network_transmit_io_limit"      xorm:"network_transmit_io_limit"`
	Created                time.Time `json:"created"                        xorm:"created"`
	Updated                time.Time `json:"updated"                        xorm:"updated"`
}

type StatusName string

const (
	StatusOn  = "on"
	StatusOff = "off"
)

// Timeframe defines query time frame
type Timeframe struct {
	ID          int64     `json:"id"                  xorm:"pk autoincr 'id'"`
	Name        string    `json:"name"                xorm:"name"`
	Start       time.Time `json:"start"               xorm:"start"`
	End         time.Time `json:"end"                 xorm:"end"`
	Status      string    `json:"status"              xorm:"status"`
	Description string    `json:"description"         xorm:"description"`
	Created     time.Time `json:"created"             xorm:"created"`
	Updated     time.Time `json:"updated"             xorm:"updated"`
}

type ApplicationResource struct {
	ID                int64                `json:"id"`
	Name              string               `json:"name"`
	ContainerResource []*ContainerResource `json:"container_resource"`
}
