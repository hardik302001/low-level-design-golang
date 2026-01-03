package main

import (
	"fmt"
	"slices"
	"sync"
)

type Elevator struct {
	ID               int
	CurrentFloor     int
	CurrentDirection Direction
	ElevatorPanel    *ElevatorPanel
	Destinations     []int
	sync.Mutex       // helps for locks over go routines
}

func NewElevator(id int) *Elevator {
	return &Elevator{
		ID:               id,
		CurrentFloor:     1,
		CurrentDirection: Still,
		ElevatorPanel:    NewElevatorPanel(id),
	}
}

func (e *Elevator) AddDestination(destinationFloor int) {
	if e.CurrentFloor == destinationFloor {
		fmt.Printf("Elevator %d is already at floor %d\n", e.ID, destinationFloor)
		return
	}

	if slices.Contains(e.Destinations, destinationFloor) {
		fmt.Printf("Elevator %d already has destination floor %d\n", e.ID, destinationFloor)
		return
	}

	e.Lock()
	defer e.Unlock()
	e.ElevatorPanel.AddDestinationFloor(destinationFloor)
	e.Destinations = append(e.Destinations, destinationFloor)
	fmt.Printf("Elevator %d received destination floor %d\n", e.ID, destinationFloor)
}

func (e *Elevator) RemoveDestination(destinationFloor int) {
	e.Lock()
	defer e.Unlock()
	for i, floor := range e.Destinations {
		if floor == destinationFloor {
			e.Destinations = append(e.Destinations[:i], e.Destinations[i+1:]...)
			e.ElevatorPanel.RemoveDestinationFloor(destinationFloor)
			break
		}
	}
	fmt.Printf("Elevator %d removed destination floor %d\n", e.ID, destinationFloor)
}

func (e *Elevator) UpdateCurrentFloor(newFloor int) {
	e.Lock()
	defer e.Unlock()
	fmt.Printf("Elevator %d moving from floor %d to floor %d\n", e.ID, e.CurrentFloor, newFloor)
	e.CurrentFloor = newFloor
}

func (e *Elevator) UpdateCurrentDirection(newDirection Direction) {
	e.Lock()
	defer e.Unlock()
	fmt.Printf("Elevator %d changing direction from %s to %s\n", e.ID, e.CurrentDirection, newDirection)
	e.CurrentDirection = newDirection
}

func (e *Elevator) HighestDestinationFloor() int {
	maxFloor := 0

	for _, floor := range e.Destinations {
		if floor > maxFloor {
			maxFloor = floor
		}
	}

	return maxFloor
}

func (e *Elevator) LowestDestinationFloor() int {
	minFloor := 100

	for _, floor := range e.Destinations {
		if floor < minFloor {
			minFloor = floor
		}
	}

	return minFloor
}

func (e *Elevator) PrintState() {
	e.Lock()
	defer e.Unlock()

	fmt.Printf(
		"Elevator %d | Current Floor: %d | Direction: %s | Destinations: %v\n",
		e.ID,
		e.CurrentFloor,
		e.CurrentDirection,
		e.Destinations,
	)
}
