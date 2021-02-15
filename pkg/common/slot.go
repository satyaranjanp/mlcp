package common

import (
	"time"
)

type AllocationType uint8
type SlotType string

const (
	Parked AllocationType = iota+1
	Unavailable
	Reserved
	Free
	DefaultAllocationType = Parked

	SlotTypeCar SlotType = "car slot"
	DefaultSlotType SlotType = SlotTypeCar

)

type Slot struct {
	Vehicle
	SlotId uint32
	Type SlotType
	InTime time.Time
	OutTime time.Time
}

func NewSlot(v Vehicle, id uint32) *Slot {
	return &Slot{
		Vehicle: v,
		SlotId:  id,
		Type: DefaultSlotType,
		InTime:  time.Time{},
		OutTime: time.Time{},
	}
}
