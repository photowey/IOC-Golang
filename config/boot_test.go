/*
 * Copyright (c) 2022, Alibaba Group;
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type redisConfig struct {
	Address string
	DB      string
}

func TestLoad(t *testing.T) {
	defer clearEnv()
	tests := []struct {
		name                string
		iocGolangConfigPath string
		wantErr             bool
	}{
		{
			"test load from file",
			"./test/ioc_golang.yaml",
			false,
		},
		{
			"test load from file",
			"./test/none-exist.yaml",
			false,
		},
		{
			"test load from invalid file with error",
			"./test/ioc_golang-bad-config.yaml",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Nil(t, os.Setenv("IOC_GOLANG_CONFIG_PATH", tt.iocGolangConfigPath))
			if err := Load(); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadConfigByPrefix(t *testing.T) {
	defer clearEnv()
	assert.Nil(t, os.Setenv("IOC_GOLANG_CONFIG_PATH", "./test/ioc_golang.yaml"))
	assert.Nil(t, Load())

	t.Run("test with multi redis config prefix", func(t *testing.T) {
		redisConfig := &redisConfig{}

		assert.Nil(t, LoadConfigByPrefix("autowire.normal.<github.com/alibaba/ioc-golang/extension/normal/redis.Impl>.param", redisConfig))
		assert.Equal(t, "0", redisConfig.DB)
		assert.Equal(t, "localhost:6379", redisConfig.Address)

		assert.Nil(t, LoadConfigByPrefix("autowire.normal.<github.com/alibaba/ioc-golang/extension/normal/redis.Impl>.db1-redis.param", redisConfig))
		assert.Equal(t, "1", redisConfig.DB)
		assert.Equal(t, "localhost:16379", redisConfig.Address)

		assert.Nil(t, LoadConfigByPrefix("autowire.normal.<github.com/alibaba/ioc-golang/extension/normal/redis.Impl>.db2-redis.param", redisConfig))
		assert.Equal(t, "2", redisConfig.DB)
		assert.Equal(t, "localhost:26379", redisConfig.Address)
	})

	t.Run("test with int value", func(t *testing.T) {
		intValue := 0
		assert.Nil(t, LoadConfigByPrefix("autowire.config.intValue", &intValue))
		assert.Equal(t, 123, intValue)
	})

	t.Run("test with string value", func(t *testing.T) {
		strValue := ""
		assert.Nil(t, LoadConfigByPrefix("autowire.config.strValue", &strValue))
		assert.Equal(t, "strVal", strValue)
	})

	t.Run("test with map value", func(t *testing.T) {
		mapValue := map[string]interface{}{}
		assert.Nil(t, LoadConfigByPrefix("autowire.config.mapValue", &mapValue))
		assert.Equal(t, "mapValue1", mapValue["mapKey1"])
		assert.Equal(t, "mapValue2", mapValue["mapKey2"])
		assert.Equal(t, "mapValue3", mapValue["mapKey3"])
	})

	t.Run("test with slice value", func(t *testing.T) {
		sliceValue := []string{}
		assert.Nil(t, LoadConfigByPrefix("autowire.config.sliceValue", &sliceValue))
		assert.Equal(t, 3, len(sliceValue))
		assert.Equal(t, "sliceStr1", sliceValue[0])
		assert.Equal(t, "sliceStr2", sliceValue[1])
		assert.Equal(t, "sliceStr3", sliceValue[2])
	})
}

func TestSetConfig(t *testing.T) {
	defer clearEnv()
	type args struct {
		yamlBytes []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetConfig(tt.args.yamlBytes); (err != nil) != tt.wantErr {
				t.Errorf("SetConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_determineConfigProperties(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "test determine config properties",
			args: args{
				key: "autowire.normal.<github.com/alibaba/ioc-golang/extension/normal/redis.Impl>.db1-redis.param",
			},
			want: []string{
				"autowire",
				"normal",
				"github.com/alibaba/ioc-golang/extension/normal/redis.Impl",
				"db1-redis",
				"param",
			},
			wantErr: false,
		},
		{
			name: "test determine config properties: unpair-1",
			args: args{
				key: "autowire.normal.github.com/alibaba/ioc-golang/extension/normal/redis.Impl>.db1-redis.param",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test determine config properties: unpair-2",
			args: args{
				key: "autowire.normal.<github.com/alibaba/ioc-golang/extension/normal/redis.Impl.db1-redis.param",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := determineConfigProperties(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("determineConfigProperties() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equalf(t, tt.want, got, "determineConfigProperties(%v)", tt.args.key)
		})
	}
}
