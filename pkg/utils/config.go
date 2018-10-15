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

package utils

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	URL      string `yaml:"url"`
	Port     int    `yaml:"port"`
	// Name defines database name
	Name string `yaml:"name"`
	// MaxIdleConns connection pool conns, default is 2
	MaxIdleConns int `yaml:"maxIdleConns"`
	// MaxOpenConns sets the maximum number of open connections to the database. The default is 10.
	MaxOpenConns int `yaml:"maxOpenConns"`
}

type PrometheusConfig struct {
	Address string `yaml:"address"`
}

type ExtraConfig struct {
	APIPort int    `yaml:"apiPort"`
	History string `yaml:"history"`
}

// GlobalConfig defines global config
type GlobalConfig struct {
	DatabaseConfig   DatabaseConfig   `yaml:"databaseConfig"`
	PrometheusConfig PrometheusConfig `yaml:"prometheusConfig"`
	ExtraConfig      ExtraConfig      `yaml:"extraConfig"`
}

func (d *DatabaseConfig) Format() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", d.Username, d.Password, d.URL, d.Port, d.Name)
}

func Unmarshal(filename string) (*GlobalConfig, error) {
	globalConfig := &GlobalConfig{}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, globalConfig)
	if err != nil {
		return nil, err
	}
	if globalConfig.DatabaseConfig.MaxIdleConns == 0 {
		globalConfig.DatabaseConfig.MaxIdleConns = 2
	}
	if globalConfig.DatabaseConfig.MaxOpenConns == 0 {
		globalConfig.DatabaseConfig.MaxOpenConns = 10
	}
	// setting default http api port
	if globalConfig.ExtraConfig.APIPort == 0 {
		globalConfig.ExtraConfig.APIPort = 9098
	}
	// setting default prometheus fetch history length
	if len(globalConfig.ExtraConfig.History) == 0 {
		globalConfig.ExtraConfig.History = "90d"
	}
	return globalConfig, nil
}
