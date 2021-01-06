//   Copyright 2021 Go-Ceres
//   Author https://github.com/go-ceres/go-ceres
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package CeresRedis

import (
	CeresConfig "github.com/go-ceres/ceres-config"
	CeresLogger "github.com/go-ceres/ceres-logger"
	"github.com/go-redis/redis"
	"time"
)

type Config struct {
	logger CeresLogger.Logger
	*redis.UniversalOptions
	debug bool
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		debug:  false,
		logger: CeresLogger.DefaultLogger.With(CeresLogger.FieldMod("ceres-redis")),
		UniversalOptions: &redis.UniversalOptions{
			DB:           0,
			PoolSize:     10,
			MaxRetries:   3,
			MinIdleConns: 100,
			DialTimeout:  time.Second * 3,
			ReadTimeout:  time.Second * 3,
			WriteTimeout: time.Second * 3,
			IdleTimeout:  time.Second * 60,
			ReadOnly:     false,
		},
	}
}

// RawConfig
func RawConfig(key string) *Config {
	def := DefaultConfig()
	if err := CeresConfig.Get(key).Scan(key); err != nil {
		def.logger.Panicd("init redis error", CeresLogger.FieldErr(err))
	}
	return def
}

// ScanConfig
func ScanConfig(name string) *Config {
	return RawConfig("ceres.redis." + name)
}

// 设置自己的logger
func (c *Config) WithLogger(logger CeresLogger.Logger) *Config {
	c.logger = logger
	return c
}

// 构建redis
func (c *Config) Build() *Redis {
	r := &Redis{
		config:          c,
		UniversalClient: redis.NewUniversalClient(c.UniversalOptions),
	}
	return r
}
