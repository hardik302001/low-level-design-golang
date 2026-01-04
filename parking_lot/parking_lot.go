package main

import (
	"fmt"
	"lld_go_parking_lot/vehicles"
	"sync"
	"time"
)

var (
	parkingLotInstance *ParkingLot
	once               sync.Once
)

type ParkingLot struct {
	Name    string
	floors  []*ParkingFloor
	tickets map[string]*ParkingTicket
}

func GetParkingLotInstance() *ParkingLot {

	once.Do(func() {
		parkingLotInstance = &ParkingLot{
			tickets: make(map[string]*ParkingTicket),
		}
	})
	return parkingLotInstance
}

func (pl *ParkingLot) AddFloor(floorID int) {
	pl.floors = append(pl.floors, NewParkingFloor(floorID))
}

func (pl *ParkingLot) DisplayAvailability() {
	fmt.Printf("Parking Lot: %s\n", pl.Name)
	for _, floor := range pl.floors {
		floor.DisplayFloorStatus()
	}
}

func (pl *ParkingLot) findParkingSpot(vehicleType vehicles.VehicleType) (*ParkingSpot, error) {
	for _, floor := range pl.floors {
		if spot := floor.FindParkingSpot(vehicleType); spot != nil {
			return spot, nil
		}
	}

	return nil, fmt.Errorf("no available parking spot found for %s", vehicleType)
}

func (pl *ParkingLot) ParkVehicle(vehicle vehicles.VehicleInterface) (*ParkingTicket, error) {
	parkingSpot, err := pl.findParkingSpot(vehicle.GetVehicleType())
	if err != nil {
		return nil, err
	}

	err = parkingSpot.ParkVehicle(vehicle)
	if err != nil {
		return nil, err
	}

	parkingTicket := NewParkingTicket(vehicle, parkingSpot)
	pl.tickets[parkingTicket.TicketId] = parkingTicket

	fmt.Printf("%s parked successfully. Ticket: %v\n",
		vehicle.GetLicenceNumber(),
		parkingTicket,
	)

	return parkingTicket, nil
}

func (p *ParkingLot) UnparkVehicle(parkingTicket *ParkingTicket, vehicle vehicles.VehicleInterface) error {
	if parkingTicket == nil {
		return fmt.Errorf("invalid parking ticket")
	}

	if vehicle.GetLicenceNumber() != parkingTicket.Vehicle.GetLicenceNumber() {
		return fmt.Errorf("vehicle licence number does not match the ticket")
	}

	if parkingTicket.Status == Closed {
		return fmt.Errorf("ticket is already used for exit")
	}

	parkingTicket.SetExitTime(time.Now())
	charge, err := parkingTicket.CalculateTotalCharge()
	if err != nil {
		return err
	}

	formattedCharge := fmt.Sprintf("%.2f", charge)
	fmt.Printf("bill for %s = %s\n", parkingTicket.Vehicle.GetLicenceNumber(), formattedCharge)

	paymentSystem := NewPaymentSystem(charge, parkingTicket)
	if err := paymentSystem.ProcessPayment(); err != nil {
		return fmt.Errorf("payment failed: %v. Vehicle is still parked", err)
	}

	parkingTicket.Spot.RemoveVehicle()
	parkingTicket.SetTotalCharge(charge)
	parkingTicket.SetTicketStatus(Closed)

	return nil
}
