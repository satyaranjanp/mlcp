package common

// Uid represnets the parking floor (level) number and slot number
// the first byte of the 32 bit integer represents level
// the last 3 bytes represent slot number in that level
type Uid uint32

func GetUID(level, slot uint32) uint32 {
	return (level << 24) + slot
}

func ParseUID(uid uint32) (level uint8, slot uint32) {
	return uint8(uint32(uid) >> 24), uint32(uid) & 0Xffffff
}

func ParseSlotId(slot uint32) (uint16, AllocationType) {
	return uint16(slot & 0xffff), AllocationType((slot >> 16) & 0xff)
}

func GetSlotId(pos uint16, allocationType AllocationType) uint32 {
	return uint32(allocationType << 16) + uint32(pos)
}

func SetPos(pos uint16, slot *Slot) uint32 {
	_, t := ParseSlotId(slot.SlotId)
	return GetSlotId(pos, t)
}

func SetAllocationType(allocType AllocationType, slot *Slot) uint32 {
	p, _ := ParseSlotId(slot.SlotId)
	return GetSlotId(p, allocType)
}