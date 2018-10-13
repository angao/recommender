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

package model

// ResourceName represents the name of the resource monitored by recommender.
type ResourceName string

// ResourceAmount represents quantity of a certain resource within a container.
// Note this keeps CPU in millicores (which is not a standard unit in APIs)
// and memory in bytes.
// Allowed values are in the range from 0 to MaxResourceAmount.
type ResourceAmount int64

// Resources is a map from resource name to the corresponding ResourceAmount.
type Resources map[ResourceName]ResourceAmount

const (
	// ResourceCPU represents CPU in millicores (1core = 1000millicores).
	ResourceCPU ResourceName = "cpu"
	// ResourceMemory represents memory, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024).
	ResourceMemory ResourceName = "memory"
	// ResourceReadDiskIO represents disk read iops
	ResourceDiskReadIO ResourceName = "disk-read-io"
	// ResourceWriteDiskIO represents disk write iops
	ResourceDiskWriteIO ResourceName = "disk-write-io"
	// ResourceNetworkIO represents network receive iops
	ResourceNetworkReceiveIO ResourceName = "network-receive-io"
	// ResourceNetworkTransmitIO represents network transmit iops
	ResourceNetworkTransmitIO ResourceName = "network-transmit-io"
	// MaxResourceAmount is the maximum allowed value of resource amount.
	MaxResourceAmount = ResourceAmount(1e14)
)

// CPUAmountFromCores converts CPU cores to a ResourceAmount.
func CPUAmountFromCores(cores float64) ResourceAmount {
	return ResourceAmountFromFloat(cores * 1000.0)
}

// CoresFromCPUAmount converts ResourceAmount to number of cores expressed as float64.
func CoresFromCPUAmount(cpuAmount ResourceAmount) float64 {
	return float64(cpuAmount) / 1000.0
}

// MemoryAmountFromBytes converts memory bytes to a ResourceAmount.
func MemoryAmountFromBytes(bytes float64) ResourceAmount {
	return ResourceAmountFromFloat(bytes)
}

// BytesFromMemoryAmount converts ResourceAmount to number of bytes expressed as float64.
func BytesFromMemoryAmount(memoryAmount ResourceAmount) float64 {
	return float64(memoryAmount)
}

// ResourceAmountMax returns the larger of two resource amounts.
func ResourceAmountMax(amount1, amount2 ResourceAmount) ResourceAmount {
	if amount1 > amount2 {
		return amount1
	}
	return amount2
}

func ResourceAmountFromFloat(amount float64) ResourceAmount {
	if amount < 0 {
		return ResourceAmount(0)
	} else if amount > float64(MaxResourceAmount) {
		return MaxResourceAmount
	} else {
		return ResourceAmount(amount)
	}
}

type ApplicationID struct {
	Name string
}

// ContainerID contains information needed to identify a Container within a cluster.
type ContainerID struct {
	ApplicationID
	// ContainerName is the name of the container, unique within a pod.
	ContainerName string
}

type ApplicationContainer struct {
	ContainerID
	Name string
}
