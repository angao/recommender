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

package prometheus

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Helper types used for parsing json returned by Prometheus.
// It would be nice if there was a public go library with these, but currently
// there's none. The server side implementation is at:
// https://github.com/prometheus/prometheus/blob/2d73d2b892853e95dbf157561e9df56ac220875e/web/api/v1/api.go#L92

// This is the top-level structure of the response.

type responseType struct {
	// Should be "success".
	Status      string   `json:status`
	Data        dataType `json:data`
	ErrorType   string   `json:errorType`
	ErrorString string   `json:error`
}

// Holds all the data returned.
type dataType struct {
	// For range vectors, this will be "matrix". Other possibilities are:
	// "vector","scalar","string".
	ResultType string `json:resultType`
	// This has different types depending on ResultType.
	Result json.RawMessage `json:result`
}

type vectorType struct {
	// Labels of the timeseries.
	Metric map[string]string `json:metric`
	// List of samples. Each sample is represented as a two-item list with
	// floating point timestamp in seconds and a string holding the value
	// of the metric.
	Value []interface{} `json:value`
}

func decodeVectorSamples(input []interface{}) (Sample, error) {
	var sample Sample
	if len(input) != 2 {
		return sample, fmt.Errorf("invalid length: %d", len(input))
	}
	ts, ok := input[0].(float64)
	if !ok {
		return sample, fmt.Errorf("invalid time: %v", input[0])
	}
	stringVal, ok := input[1].(string)
	if !ok {
		return sample, fmt.Errorf("invalid value: %v", input[1])
	}
	var val float64
	fmt.Sscan(stringVal, &val)
	sample.Value = val
	sample.Timestamp = time.Unix(int64(ts), 0)
	return sample, nil
}

// Decodes timeseries from a Prometheus response.
func decodeTimeseriesFromResponse(input io.Reader) ([]Timeseries, error) {
	var resp responseType
	err := json.NewDecoder(input).Decode(&resp)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse response: %v", err)
	}
	if resp.Status != "success" || resp.Data.ResultType != "vector" {
		return nil, fmt.Errorf("invalid response status: %s or type: %s", resp.Status, resp.Data.ResultType)
	}
	var vectors []vectorType
	err = json.Unmarshal(resp.Data.Result, &vectors)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse response vector: %v", err)
	}
	res := make([]Timeseries, 0)
	for _, vector := range vectors {
		sample, err := decodeVectorSamples(vector.Value)
		if err != nil {
			return []Timeseries{}, fmt.Errorf("error decoding sample: %v", err)
		}
		res = append(res, Timeseries{Labels: vector.Metric, Sample: sample})
	}
	return res, nil

}
