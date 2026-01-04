package main

import (
	"fmt"
	"lld_go_parking_lot/vehicles"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

func main() {
	parkingLot := GetParkingLotInstance()

	parkingLot.Name = "Central Parking Lot"

	parkingLot.AddFloor(0)
	parkingLot.AddFloor(1)

	parkingLot.DisplayAvailability()

	for i := 1; i <= 10; i++ {
		wg.Add(1)

		go parkCar(i, parkingLot)
	}

	wg.Wait()

	parkingLot.DisplayAvailability()

	truck := vehicles.NewTruck("truck-1")
	ticket, err := parkingLot.ParkVehicle(truck)
	if err != nil {
		fmt.Printf("Failed to park %s: %v\n", truck.LicenceNumber, err)
		return
	}

	time.Sleep(3 * time.Second)
	err = parkingLot.UnparkVehicle(ticket.TicketId)
	if err != nil {
		fmt.Printf("Failed to unpark for ticketId %s: %v\n", ticket.TicketId, err)
		return
	}

	// duplicate ticket unpark attempt
	err = parkingLot.UnparkVehicle(ticket.TicketId)
	if err != nil {
		fmt.Printf("Failed to unpark for ticketId %s: %v\n", ticket.TicketId, err)
		return
	}
}

func parkCar(ind int, parkingLot *ParkingLot) {
	defer wg.Done()

	car := vehicles.NewCar(fmt.Sprintf("car-%d", ind))

	_, err := parkingLot.ParkVehicle(car)
	if err != nil {
		fmt.Printf("Failed to park %s: %v\n", car.LicenceNumber, err)
		return
	}
}
