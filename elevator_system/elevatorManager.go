package main

import (
	"fmt"
	"sort"
)

type ElevatorManager struct {
	Building *Building
}

func NewElevatorManager(building *Building) *ElevatorManager {
	return &ElevatorManager{Building: building}
}

func (em *ElevatorManager) OperateAllElevators() {
	for _, elevator := range em.Building.Elevators {
		elevator.PrintState()
		go em.OperateElevator(elevator) // go thread for each elevator
	}
}

func (em *ElevatorManager) OperateElevator(elevator *Elevator) {
	for { // infinte loop to keep elevator operating(real world scenario -> event can come any time)
	
		if len(elevator.Destinations) == 0 {
			elevator.CurrentDirection = Still
			continue
		}

		sort.Ints(elevator.Destinations)

		if elevator.CurrentDirection == Up {
			em.MoveElevatorUp(elevator)
		} else if elevator.CurrentDirection == Down {
			em.MoveElevatorDown(elevator)
		} else {
			em.DecideDirection(elevator)
		}

	}
}

// move to nearest request
func (em *ElevatorManager) DecideDirection(elevator *Elevator) {
	if len(elevator.Destinations) == 0 {
		return
	}

	currentFloor := elevator.CurrentFloor
	nearestDestination := elevator.Destinations[0]
	if nearestDestination > currentFloor {
		elevator.UpdateCurrentDirection(Up)
		em.MoveElevatorUp(elevator)
	} else {
		elevator.UpdateCurrentDirection(Down)
		em.MoveElevatorDown(elevator)
	}
}

func (em *ElevatorManager) MoveElevatorUp(elevator *Elevator) {
	for i := 0; i < len(elevator.Destinations); i++ {
		destination := elevator.Destinations[i]
		if destination >= elevator.CurrentFloor {
			elevator.UpdateCurrentFloor(destination)
			elevator.RemoveDestination(destination)
			i-- // because we have removed an element from the slice
		} else {
			// skip destinations below current floor
		}
	}

	if len(elevator.Destinations) == 0 {
		elevator.UpdateCurrentDirection(Still)
	} else {
		elevator.UpdateCurrentDirection(Down)
	}
}

func (em *ElevatorManager) MoveElevatorDown(elevator *Elevator) {
	for i := len(elevator.Destinations) - 1; i >= 0; i-- {
		destination := elevator.Destinations[i]
		if destination <= elevator.CurrentFloor {
			elevator.UpdateCurrentFloor(destination)
			elevator.RemoveDestination(destination)
		} else {
			// skip destinations above current floor
		}
	}

	if len(elevator.Destinations) == 0 {
		elevator.UpdateCurrentDirection(Still)
	} else {
		elevator.UpdateCurrentDirection(Up)
	}
}

// manager will assign the best elevator for the hall call request
func (em *ElevatorManager) AssignElevator(floor int, direction Direction) (bestElevator *Elevator) {
	bestElevator = em.FindClosestElevator(floor, direction)
	if bestElevator != nil {
		bestElevator.AddDestination(floor)
		fmt.Printf("Elevator %d assigned to floor %d with direction %s\n", bestElevator.ID, floor, direction)
	}
	return bestElevator
}

// mamnager will find best elevator for the hall call request
func (em *ElevatorManager) FindClosestElevator(floor int, direction Direction) *Elevator {
	var closestElevator *Elevator
	minDistance := int(1e9)

	for _, elevator := range em.Building.Elevators {
		elevator.Lock()
		distance := em.calculateDistance(elevator, floor, direction)
		elevator.Unlock()

		if distance < minDistance {
			minDistance = distance
			closestElevator = elevator
		}
	}

	return closestElevator
}

func (em *ElevatorManager) calculateDistance(elevator *Elevator, floor int, direction Direction) int {
	currentFloor := elevator.CurrentFloor
	currentDirection := elevator.CurrentDirection

	fmt.Println("Calculating distance for Elevator", elevator.ID, "at floor", currentFloor, "going", currentDirection, "to floor", floor, "going", direction)

	// Case 1: Elevator is idle
	if currentDirection == Still || len(elevator.Destinations) == 0 {
		return abs(floor - currentFloor)
	}

	// Case 2: Elevator moving in same direction
	if currentDirection == direction {
		if (direction == Up && floor >= currentFloor) || (direction == Down && floor <= currentFloor) {
			// Request is ahead in same direction → pick up immediately
			return abs(floor - currentFloor)
		} else {
			// Request is behind → calculate distance after finishing current sweep
			if direction == Up {
				farthest := elevator.FarthestDestination()
				return abs(farthest - currentFloor) + abs(farthest - floor)
			} else { // Down
				nearest := elevator.NearestDestination()
				return abs(currentFloor - nearest) + abs(floor - nearest)
			}
		}
	}

	// Case 3: Elevator moving in opposite direction
	if (currentDirection == Up && direction == Down) || (currentDirection == Down && direction == Up) {
		if currentDirection == Up {
			farthest := elevator.FarthestDestination()
			return abs(farthest - currentFloor) + abs(farthest - floor)
		} else { // Down
			nearest := elevator.NearestDestination()
			return abs(currentFloor - nearest) + abs(floor - nearest)
		}
	}

	// Fallback: very far / unlikely elevator
	return 1000
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
