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

package apis

type RecommendResource struct {
	PodName           string            `json:"pod_name"`
	ContainerName     string            `json:"container_name"`
	Namespace         string            `json:"namespace"`
	CPUResource       CPUResource       `json:"cpu_resource"`
	MemoryResource    MemoryResource    `json:"memory_resource"`
	DiskIOResource    DiskIOResource    `json:"disk_io_resource"`
	NetworkIOResource NetworkIOResource `json:"network_io_resource"`
}

type CPUResource struct {
	CPURequest float64 `json:"cpu_request"`
	CPULimit   float64 `json:"cpu_limit"`
}

type MemoryResource struct {
	MemoryRequest float64 `json:"memory_request"`
	MemoryLimit   float64 `json:"memory_limit"`
}

type DiskIOResource struct {
	DiskReadIOResource
	DiskWriteIOResource
}

type DiskReadIOResource struct {
	DiskReadIORequest float64 `json:"disk_read_io_request"`
	DiskReadIOLimit   float64 `json:"disk_read_io_limit"`
}

type DiskWriteIOResource struct {
	DiskWriteIORequest float64 `json:"disk_write_io_request"`
	DiskWriteIOLimit   float64 `json:"disk_write_io_limit"`
}

type NetworkIOResource struct {
	NetworkIORequest float64 `json:"network_io_request"`
	NetworkIOLimit   float64 `json:"network_io_limit"`
}
