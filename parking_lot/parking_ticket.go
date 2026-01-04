package main

import (
	"fmt"
	"lld_go_parking_lot/vehicles"
	"time"

	"github.com/google/uuid"
)

const baseCharge = 100.00

type TicketStatus string

const (
	Active TicketStatus = "ACTIVE"
	Closed TicketStatus = "CLOSED"
)

type ParkingTicket struct {
	TicketId    string
	EntryTime   time.Time
	ExitTime    *time.Time
	Vehicle     vehicles.VehicleInterface
	Spot        *ParkingSpot
	TotalCharge float64
	Status      TicketStatus
}

func NewParkingTicket(vehicle vehicles.VehicleInterface, spot *ParkingSpot) *ParkingTicket {
	return &ParkingTicket{
		TicketId:  generateTicketID(),
		EntryTime: time.Now(),
		Vehicle:   vehicle,
		Spot:      spot,
		Status:    Active,
	}
}

func generateTicketID() string {
	return "TICKET-" + uuid.New().String()
}

func (p *ParkingTicket) SetExitTime(exitTime time.Time) {
	p.ExitTime = &exitTime
}

func (p *ParkingTicket) SetTotalCharge(totalCharge float64) {
	p.TotalCharge = totalCharge
}

func (p *ParkingTicket) SetTicketStatus(status TicketStatus) {
	p.Status = status
}

func (p *ParkingTicket) CalculateTotalCharge() (float64, error) {
	if p.ExitTime == nil {
		return 0, fmt.Errorf("exit time not set")
	}

	duration := p.ExitTime.Sub(p.EntryTime)
	hours := duration.Hours()
	additionalCharge := hours * p.Vehicle.GetVehicleCost()
	p.TotalCharge = baseCharge + additionalCharge

	return p.TotalCharge, nil
}
