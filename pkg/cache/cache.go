package cache

import (
	"mlcp/pkg/common"
	"mlcp/pkg/config"
)

type Cache interface {
	GetNearestSlot() *common.Slot
	AssignSlot(slot *common.Slot) *common.Slot
	FreeUpSlot(slot *common.Slot) *common.Slot
}

func SetupCache() (Cache, error) {
	switch config.CacheType {
	case "local":
		return setupLocalCache()
	}
	return nil, nil
}