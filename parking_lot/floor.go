package main

import (
	"fmt"
	"lld_go_parking_lot/vehicles"
)

const (
	CarSpotCount        = 5
	VanSpotCount        = 3
	TruckSpotCount      = 2
	MotorcycleSpotCount = 10
)

type ParkingFloor struct {
	FloorID      int
	ParkingSpots map[vehicles.VehicleType]map[int]*ParkingSpot
}

func NewParkingFloor(floorID int) *ParkingFloor {
	parkingSpots := make(map[vehicles.VehicleType]map[int]*ParkingSpot)

	parkingSpots[vehicles.CarType] = createParkingSpots(CarSpotCount, vehicles.CarType)
	parkingSpots[vehicles.VanType] = createParkingSpots(VanSpotCount, vehicles.VanType)
	parkingSpots[vehicles.TruckType] = createParkingSpots(TruckSpotCount, vehicles.TruckType)
	parkingSpots[vehicles.MotorcycleType] = createParkingSpots(MotorcycleSpotCount, vehicles.MotorcycleType)

	return &ParkingFloor{FloorID: floorID, ParkingSpots: parkingSpots}
}

func createParkingSpots(count int, vehicleType vehicles.VehicleType) map[int]*ParkingSpot {
	spots := make(map[int]*ParkingSpot)
	for i := 1; i <= count; i++ {
		spots[i] = NewParkingSpot(i, vehicleType)
	}
	return spots
}

func (pf *ParkingFloor) FindParkingSpot(vehicleType vehicles.VehicleType) *ParkingSpot {
	for _, spot := range pf.ParkingSpots[vehicleType] {
		if spot.IsParkingSpotFree() {
			return spot
		}
	}
	return nil
}

func (pf *ParkingFloor) DisplayFloorStatus() {
	fmt.Printf("Floor ID: %d\n", pf.FloorID)

	for vehicleType, spotMap := range pf.ParkingSpots {
		fmt.Printf("%s Spots:\n", vehicleType)
		count := 0

		for _, spot := range spotMap {
			if spot.IsParkingSpotFree() {
				count++
			}
		}

		fmt.Printf("%s Spot: %d Free\n", vehicleType, count)
	}
}
