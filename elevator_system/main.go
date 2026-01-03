package main

import "sync"

func main() {
	building := NewBuilding()
	manager := NewElevatorManager(building)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		elevator := building.Floors[1].HallPanels[1].RequestElevator(manager, Up) // request elevator from floor 1 to go Up -> hall call
		elevator.AddDestination(6)                                                // destination floor 6
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		secondElevator := building.Floors[8].HallPanels[2].RequestElevator(manager, Down)
		secondElevator.AddDestination(7)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		thirdElevator := building.Floors[3].HallPanels[0].RequestElevator(manager, Up)
		thirdElevator.AddDestination(12)
	}()

	wg.Wait() // wait unitl all requests are done

	go manager.OperateAllElevators() // start operating/simulating all elevators

	select {} // it is needed for blocking main for exiting as we have spawned goroutines
}
