package rooms

import (
	"../items"
)

type Room struct {
	id                     RoomId
	name                   string
	moveToRoomResultText   *string               //Игрок приходит в пространство - "ты в своей комнате", "ничего интересного", "кухня, ничего интересного"
	emptyLookAroundMsg     *string               //Игрок осматривается в пустом пространстве - имя пустого пространства - "пустая комната, пустой корридор, ты находишься на кухне, надо собрать рюкзак и идти в универ"
	invenoryNameIfNotEmpty *string               //Игрок осматривается в заполненном пространстве - имя инвентаря помещения invenoryNameIfNotEmpty, например: "на столе: %v"
	collection             items.ItemsCollection //Чтобы не усложнять - одна комната - одно место хранения
	availableWays          []RoomId              //Куда можно перейти из текущей комнаты
	autoJumpToRoomId       *RoomId               //Если указано, то при попытке перехода в эту комнату происходит автоматический переход в следующую за ней. архитектурное органичение: следующая не должна быть комнатой пробросом, если заметим в конфигурации - вызовем панику
	roomIsJumpTo           bool                  //Если true, значит комната задана хоть один раз в качестве прыжковой для другой комнаты
}

type RoomId int

type Library struct {
	roomDataByRoomId map[RoomId]*Room
	defaultRoom      *RoomId
}
