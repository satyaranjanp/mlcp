package common

import "mlcp/pkg/config"

type Vehicle interface {
	GetRegNo() string
	GetType() string
}

type Car struct {
	RegNo string
	vechileType string
}

func NewCar(regNo string) *Car {
	return &Car{RegNo: regNo, vechileType: config.VehicleTypeCar}
}

func (c *Car) GetRegNo() string {
	return c.RegNo
}

func (c *Car) GetType() string {
	return c.vechileType
}