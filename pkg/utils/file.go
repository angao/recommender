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

// DatabaseConfig defines database connection config file
type DatabaseConfig struct {
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	URL      string `yaml:"URL"`
	Port     int    `yaml:"Port"`
	DBName   string `yaml:"DBName"`
}

func (d *DatabaseConfig) Format() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", d.Username, d.Password, d.URL, d.Port, d.DBName)
}

func Unmarshal(filename string) (*DatabaseConfig, error) {
	databaseConfig := &DatabaseConfig{}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, databaseConfig)
	if err != nil {
		return nil, err
	}
	return databaseConfig, nil
}
