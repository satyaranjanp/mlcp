package cache

import (
	"mlcp/pkg/common"
	"math"
	"mlcp/pkg/config"
	"time"
)

const (
	listLength         = 100
	defaultVehicleType = config.VehicleTypeCar
)

// map["car"][<level>]<cache>
type localCache struct {
	l_cache map[string][]*cache
}

type cache struct {
	c []*List
}

func setupLocalCache() (*localCache, error) {
	c := &localCache{}
	c.l_cache = make(map[string][]*cache)
	for _, v := range config.VechileTypes {
		c.l_cache[v] = make([]*cache, config.MlcpLevelCount)
		for i, _ := range c.l_cache[v] {
			cnt := math.Ceil(float64(config.SlotsPerLevel / 100))
			c.l_cache[v][i] = &cache{c: make([]*List, int(cnt))}
		}
	}
	return c, nil
}

func (lc *localCache) BuildCacheFromDB() {

}

func (lc *localCache) GetNearestSlot() *common.Slot {
	return lc.getDefaultMinSlot()
}

func (lc *localCache) AssignSlot(slot *common.Slot) *common.Slot {
	slot.SlotId = common.SetAllocationType(common.Parked, slot)
	slot.InTime = time.Now()
	if true == lc.assignSlot(slot) {
		return slot
	}
	return nil
}

func (lc *localCache) FreeUpSlot(slot *common.Slot) *common.Slot {
	slot.SlotId = common.SetAllocationType(common.Free, slot)
	slot.OutTime = time.Now()
	if true == lc.freeUpSlot(slot) {
		return slot
	}
	return nil
}

func (lc *localCache) getDefaultMinSlot() *common.Slot {
	var i uint16 = 0
	for i = 0; i < config.MlcpLevelCount; i++ {
		for _, list := range lc.l_cache[defaultVehicleType][i].c {
			if v := list.getMin(); v > 0 {
				//p, t := common.ParseSlotData(v)
				id := common.GetUID(uint32(i), v)
				return &common.Slot{
					Vehicle:        nil,
					SlotId:			id,
				}
			}
		}
	}

	return nil
}

func (lc *localCache) assignSlot(slot *common.Slot) bool {
	level, s := common.ParseUID(slot.SlotId)
	pos := getPos(s)
	p := int(math.Floor(float64(pos/100)))
	return lc.l_cache[slot.Vehicle.GetType()][level].c[p].remove(s)
}

func (lc *localCache) freeUpSlot(slot *common.Slot) bool {
	level, s := common.ParseUID(slot.SlotId)
	pos := getPos(s)
	p := int(math.Floor(float64(pos/100)))
	return lc.l_cache[defaultVehicleType][level].c[p].add(s)
}

type List struct {
	head *node
}

type node struct {
	slot uint32
	next *node
}

func NewList() *List {
	return &List{
		head:   nil,
	}
}

func newNode(slot uint32) *node {
	return &node {
		slot: slot,
		next: nil,
	}
}

// Add nodes to the list with maintaining ascending order of slot values
func (l *List) add(slot uint32) bool {
	n := newNode(slot)
	if l.head == nil {
		l.head = n
		return true
	}

	temp := l.head
	if temp.next == nil {
		n.next = temp
		l.head = n
	}

	for ; temp.next != nil && temp.next.slot < n.slot; {
		temp = temp.next
	}
	n.next = temp.next
	temp.next = n
	return true
}

func getPos(slot uint32) uint16 {
	p, _ := common.ParseSlotId(slot)
	return p
}

func (l *List) remove(slot uint32) bool {
	if l.head == nil {
		return false
	}

	pos := getPos(slot)
	temp := l.head
	if getPos(temp.slot) == pos {
		l.head = temp.next
		return true
	}
	for ; getPos(temp.slot) < pos && temp.next != nil && getPos(temp.next.slot) < pos; {
		temp = temp.next
	}

	if temp.next != nil && getPos(temp.next.slot) == pos {
		temp.next = temp.next.next
	}
	return true
}

func (l *List) getMin() uint32 {
	if l.head != nil {
		return l.head.slot
	}
	return 0
}

func (l *List) setMin() {
	if l.head != nil {
		l.head = l.head.next
	}
}

/*
func (l *List) GetAllAvailableSlots() []uint32 {
	slots := []uint32{}
	var pos uint16 = 1
	var i uint
	temp := l.head
	for i = 0; i < config.SlotsPerLevel; i++ {
		if temp != nil && temp.slot == pos {
			pos++
			temp = temp.next
			continue
		}
		slots = append(slots, Slot{Pos: pos, AllocationType: free})
	}
	return slots
}

func (l *List) ReadList() []Slot {
	allocated := []Slot{}
	if l.head != nil {
		temp := l.head
		for ; temp.next != nil; {
			allocated = append(allocated, Slot{temp.pos, temp.allocationType})
		}
	}
	return allocated
}

func (l *List) GetAllParkedSlots() []Slot {
	slots := []Slot{}
	if l.head != nil {
		temp := l.head
		for ; temp.next != nil; {
			if temp.allocationType == parked {
				slots = append(slots, Slot{Pos: temp.pos, AllocationType:parked})
			}
		}
	}
	return slots
}

func (l *List) GetAllReservedSlots() []Slot {
	slots := []Slot{}
	if l.head != nil {
		temp := l.head
		for ; temp.next != nil; {
			if temp.allocationType == reserved {
				slots = append(slots, Slot{Pos: temp.pos, AllocationType:reserved})
			}
		}
	}
	return slots
}

func (l *List) GetAllUnavailableSlots() []Slot {
	slots := []Slot{}
	if l.head != nil {
		temp := l.head
		for ; temp.next != nil; {
			if temp.allocationType == unavailable {
				slots = append(slots, Slot{Pos: temp.pos, AllocationType:unavailable})
			}
		}
	}
	return slots
}
*/