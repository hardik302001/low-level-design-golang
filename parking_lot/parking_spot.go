package main

import (
	"fmt"
	"lld_go_parking_lot/vehicles"
	"sync"
)

type ParkingSpot struct {
	SpotID         int
	VehicleType    vehicles.VehicleType
	CurrentVehicle *vehicles.VehicleInterface
	lock           sync.Mutex
}

func NewParkingSpot(spotID int, vehicleType vehicles.VehicleType) *ParkingSpot {
	return &ParkingSpot{SpotID: spotID, VehicleType: vehicleType}
}

func (ps *ParkingSpot) IsParkingSpotFree() bool {
	return ps.CurrentVehicle == nil
}

func (ps *ParkingSpot) ParkVehicle(vehicle vehicles.VehicleInterface) error {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	if vehicle.GetVehicleType() != ps.VehicleType {
		return fmt.Errorf("vehicle type mismatch: expected %s, got %s", ps.VehicleType, vehicle.GetVehicleType())
	}
	if ps.CurrentVehicle != nil {
		return fmt.Errorf("parking spot already occupied")
	}

	ps.CurrentVehicle = &vehicle
	return nil
}

func (ps *ParkingSpot) RemoveVehicle() {
	ps.CurrentVehicle = nil
}
