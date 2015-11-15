package engine

import (
	"encoding/json"
)

type GlobalConfig struct {
	Pools []UserPoolConfig
}

type AmmoProviderConfig struct {
	AmmoType   string
	AmmoSource string
	AmmoLimit  int
	Loop       int
}

type GunConfig struct {
	GunType    string
	Parameters map[string]interface{}
}

type ResultListenerConfig struct {
	ListenerType string
	Destination  string
}

type LimiterConfig struct {
	LimiterType string
	Parameters  map[string]interface{}
}

type CompositeLimiterConfig struct {
	Steps []LimiterConfig
}

type UserConfig struct {
	Name           string
	Gun            *GunConfig
	AmmoProvider   *AmmoProviderConfig
	ResultListener *ResultListenerConfig
	Limiter        *LimiterConfig
}

type UserPoolConfig struct {
	Name           string
	Gun            *GunConfig
	AmmoProvider   *AmmoProviderConfig
	ResultListener *ResultListenerConfig
	UserLimiter    *LimiterConfig
	StartupLimiter *LimiterConfig
}

func NewConfigFromJson(jsonDoc []byte) (gc GlobalConfig, err error) {
	err = json.Unmarshal(jsonDoc, &gc)
	return
}