package rooms

import (
	"../items"
	"../tools"
	"fmt"
	"log"
	"strings"
)

// Возвращает указатель на комнату по ее id
func (r RoomId) GetRoom() *Room {
	room := lib.roomDataByRoomId[r]
	if room == nil {
		log.Panicf("Попытка обращения к комнате id=%v, но она не была создана, сначала создайте ее", r)
	}
	return room
}

func (r *Room) GetName() string {
	if r != nil {
		return r.name
	} else {
		return "Неизвестно где"
	}
}

func (r RoomId) GetName() string {
	if roomData := r.GetRoom(); roomData != nil {
		return roomData.name
	} else {
		return "Неизвестно где"
	}
}

// Удаление предмета из комнаты
// Возвращает булевый признак был ли предмет найден в комнате ("взяли в руки")
// Если мы не сможем распорядится предметом - он бесследно исчезнет
func (r *Room) GetItemFromRoom(id items.ItemId) (getInArms bool) {
	r.collection.Items, getInArms = tools.RemoveObjectFromStice(r.collection.Items, id)
	return
}

// Сбросим хранилище комнат, оно одно у нас в проекте
// Возвращаем указатель на хранилище
func Init() *Library {
	lib.roomDataByRoomId = make(map[RoomId]*Room)
	lib.defaultRoom = nil
	return &lib
}

func (r *Room) SetAsDefaultForPlayer() {
	lib.defaultRoom = &r.id
}

// Добавление комнат
func (lib *Library) AddRoom(id RoomId, name string) *Room {
	if lib.roomDataByRoomId == nil {
		log.Panicf("Попытка вызова AddItemToLibrary без вызова InitItems. Карта itemLibrary.idToData не была инициализирована")
	}
	mustNotExits(id)

	room := Room{
		id:            id,
		name:          name,
		collection:    items.ItemsCollection{},
		availableWays: nil,
	}
	lib.roomDataByRoomId[id] = &room

	return &room //Вернем указатель на комнату чтобы с комнатой можно было работать
}

// Проверяем существует ли комната с таким id, если существует - паникуем
func mustNotExits(id RoomId) {
	if _, exists := lib.roomDataByRoomId[id]; exists {
		log.Panicf("Попытка создания комнаты с id=%v, а такая уже существует", id)
	}
}

// Добавление предметов в комнату
func (r *Room) AddItems(ids ...items.ItemId) {
	//Проверка что предметы существуют
	ok, lostItems := items.CheckAllItemsExists(ids)
	if !ok {
		log.Panicf("items %v not found in library", lostItems)
	}

	//Непосредственно добавляем
	r.collection.Items = append(r.collection.Items, ids...)
}

// Результат "Идти" если переход успешен
func (r *Room) SetMoveToRoomResultText(msg string) {
	r.moveToRoomResultText = &msg
}

// Результат "Осмотреться" если пространстов пустое
func (r *Room) SetEmptyLookAroundMsg(msg string) {
	r.emptyLookAroundMsg = &msg
}

// Название инвентаря в котором в комнате предметы если пространстов не пустое
func (r *Room) SetInvenoryNameIfNotEmpty(msg string) {
	r.invenoryNameIfNotEmpty = &msg
}

// Добавление маршрута между комнатами. Если bothWays==false, то машрут допустим только id1 -> id2
func (lib *Library) AddWays(id1 RoomId, id2 RoomId, bothWays bool) {
	if lib.roomDataByRoomId == nil {
		log.Panicf("Попытка вызова AddWays без вызова InitItems. Карта itemLibrary.idToData не была инициализирована")
	}

	room1 := id1.GetRoom()
	room2 := id2.GetRoom()

	room1.availableWays = append(room1.availableWays, id2)
	if bothWays {
		room2.availableWays = append(room2.availableWays, id1)
	}
}

func GetDefaultRoom() *RoomId {
	if lib.defaultRoom == nil {
		log.Panicf("Попытка вызова GetDefaultRoom без вызова room.SetAsDefaultForPlayer()")
	}
	return lib.defaultRoom
}

func (r *Room) GetAvaibleWaysString() string {
	var strsTmp []string
	for _, way := range r.availableWays {
		strsTmp = append(strsTmp, way.GetName())
	}
	if len(strsTmp) == 0 {
		return "никуда"
	}
	return strings.Join(strsTmp, ", ")
}

// Получим id предмета в текущей комнате, если есть
// Предметы могут быть одни и те же в разных комнатах, а может с одинаковым названием. Вернем id первого попавшегося в комнате
func (r *Room) GetItemIdByItemNameInRoom(itemName string) *items.ItemId {
	return r.collection.GetIdByName(itemName)
}

// Метод возвращает id комнаты в которую можно пройти по ее имени (если такая комната доступна), иначе nil
func (r *Room) GetWayIdByRoomName(roomName string) *RoomId {
	for _, wayId := range r.availableWays {
		if wayId.GetName() == roomName {
			jumpToRoomId := wayId.GetRoom().autoJumpToRoomId
			if jumpToRoomId != nil { //Проверяем, если для комнаты предусмотрена проброс, то пробрасываем в целевую комнату
				return jumpToRoomId
			} else {
				return &wayId
			}
		}
	}
	return nil
}

// Метод возвращает специфичные сообщения о результате перемещения в комнату
func (r *Room) GetMoveInMsg() string {
	if r.moveToRoomResultText != nil {
		return fmt.Sprintf("%v. можно пройти - %v", *r.moveToRoomResultText, r.GetAvaibleWaysString())
	} else {
		return fmt.Sprintf("Вы в/на %v. можно пройти - %v", r.name, r.GetAvaibleWaysString())
	}
}

// Метод возвращает специфичные сообщения о том что игрок видит в комнате
func (r *Room) GetLookAroundMsg() string {
	//Тут у нас несколько случаев: если комната не пустая, если комната пустая, если дополнительные сообщения не заданы

	if !r.CheckNotEmpty() {
		if r.emptyLookAroundMsg != nil {
			return fmt.Sprintf("%v. можно пройти - %v", *r.emptyLookAroundMsg, r.GetAvaibleWaysString())
		} else {
			return fmt.Sprintf("%v. можно пройти - %v", "Ничего", r.GetAvaibleWaysString())
		}
	} else {
		if r.invenoryNameIfNotEmpty != nil {
			return fmt.Sprintf("%v: %v. можно пройти - %v", *r.invenoryNameIfNotEmpty, r.collection.String(), r.GetAvaibleWaysString())
		} else {
			return fmt.Sprintf("%v: %v. можно пройти - %v", "Вы видите", r.collection.String(), r.GetAvaibleWaysString())
		}
	}
}

// Проверка есть ли в комнате предметы
func (r *Room) CheckNotEmpty() bool {
	return len(r.collection.Items) > 0
}

// Превращает комнату в комнату проброс в следующую комнату
func (r *Room) SetRoomAsJumpTo(nextRoomId RoomId) {
	//Небольшие проверки
	_ = nextRoomId.GetRoom()                 //Проверим что комната существует. Если не существует, то вызовем панику
	r.checkRoomCanSetAsJumpRoom()            //Проверим что комната может быть назначена прыжковой
	nextRoomId.checkRoomCanSetAsTargetRoom() //Проверим что целевая комната может быть назначена таковой

	r.autoJumpToRoomId = &nextRoomId
	nextRoomId.setRoomsIsRoomJumpTo() //добавим комнате
}

func (r *Room) checkRoomCanSetAsJumpRoom() {
	if r.roomIsJumpTo {
		log.Panicf("Попытка сделать комнату id=%v прыжковой, но она уже целевая", r.id)
	}
}

func (r RoomId) checkRoomCanSetAsTargetRoom() {
	if r.GetRoom().autoJumpToRoomId != nil {
		log.Panicf("Попытка сделать комнату id=%v целевой, но она уже прыжковая", r)
	}
}

// Зададим целевой комнате проброса признак, что она целевая
func (r RoomId) setRoomsIsRoomJumpTo() {
	r.GetRoom().roomIsJumpTo = true
}
