package players

import (
	"../items"
	"../rooms"
	"fmt"
)

// Метод перемещения игрока между комнатами
func (p *PlayerType) MoveToRoom(roomName string) (res MoveToRoomResult) {
	roomData := p.getCurrentRoom()
	if roomData == nil {
		res.Msg = fmt.Sprintf("К сожалению, мы внезапно забыли где вы находитесь, поэтому мы не можем проверить как Вам пройти в/на %v", roomName)
		return
	}
	nextRoom := roomData.GetWayIdByRoomName(roomName)
	if nextRoom == nil {
		res.Msg = fmt.Sprintf("нет пути в %v", roomName)
		return
	}

	p.setCurrentRoom(nextRoom)

	res.Success = true
	res.Msg = p.getCurrentRoom().GetMoveInMsg()
	return
}

// Метод принудительно указания в какой комнате находится игрок
func (p *PlayerType) SetCurrentRoom(roomId int) {
	roomIdTmp := rooms.RoomId(roomId)
	_ = roomIdTmp.GetRoom() //Т.к. метод экспортный, то проверим что комната существует
	p.setCurrentRoom(&roomIdTmp)
}

// Метод принудительно указания в какой комнате находится игрок
func (p *PlayerType) setCurrentRoom(roomId *rooms.RoomId) {
	p.currentRoomId = roomId
}

// Метод получения сведений о комнате в которой сейчас находится игрок
func (p *PlayerType) getCurrentRoom() *rooms.Room {
	return p.currentRoomId.GetRoom()
}

// Получить сведения об окружении игрока
func (p *PlayerType) LookAround() (res LookAroundResult) {
	res.Msg = p.getCurrentRoom().GetLookAroundMsg()
	return
}

// Забрать предмет из комнаты
func (p *PlayerType) GetItem(itemName string) (res GetItemFromRoomToPlayerResult) {
	//Пытаемся получить id предмета в комнате
	itemId := p.getCurrentRoom().GetItemIdByItemNameInRoom(itemName)
	if itemId == nil {
		res.Msg = fmt.Sprintf("нет такого")
		return
	}

	//Пробуем забрать предмет из комнаты
	if !p.getCurrentRoom().GetItemFromRoom(*itemId) {
		res.Msg = "Вы пытались взять предмет, но ничего не вышло, вы можете попробовать потом еще раз"
		return
	}
	//Теперь, если не положим предмет в инвентарь - он бесследно исчезнет!

	//Если получилось забрать предмет из комнаты - кладем его в инвентарь
	p.addItemToInventory(itemId)
	res.MovedToInventory = true
	res.Msg = fmt.Sprintf("предмет добавлен в инвентарь: %v", itemId.GetName())
	return
}

// Метод добавления игроку предметов в инвентарь
func (p *PlayerType) addItemToInventory(itemId *items.ItemId) {
	p.inventory.Items = append(p.inventory.Items, *itemId)
}

// Инициализация настроек игрока по умолчанию
// Игрок у нас один
func InitPlayer() PlayerType {
	return PlayerType{
		currentRoomId: rooms.GetDefaultRoom(),
	}
}
