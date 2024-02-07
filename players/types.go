package players

import (
	"../items"
	"../rooms"
)

type PlayerType struct {
	currentRoomId *rooms.RoomId //В каком помещении сейчас находится
	inventory     items.ItemsCollection
}

type MoveToRoomResult struct {
	Success bool
	Msg     string
}

type LookAroundResult struct {
	Msg string
}

type GetItemFromRoomToPlayerResult struct {
	MovedToInventory bool
	Msg              string
}
