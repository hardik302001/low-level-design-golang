package main

import "fmt"

type Direction string

const (
	Up    Direction = "Up"
	Down  Direction = "Down"
	Still Direction = "Still"
)

type HallPanel struct {
	PanelID              int
	DirectionInstruction Direction
	SourceFloor          int
}

func NewHallPanel(panelID int, sourceFloor int) *HallPanel {
	return &HallPanel{PanelID: panelID, SourceFloor: sourceFloor, DirectionInstruction: Still}
}

func (h *HallPanel) SetDirectionInstructions(directionInstruction Direction) {
	h.DirectionInstruction = directionInstruction
}

func (h *HallPanel) RequestElevator(manager *ElevatorManager, direction Direction) (elevator *Elevator) {
	fmt.Printf("Panel %d requesting elevator with direction %s from floor %d\n", h.PanelID, direction, h.SourceFloor)
	return manager.AssignElevator(h.SourceFloor, direction)
}
