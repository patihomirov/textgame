package items

type ItemId int

type Item struct {
	id   ItemId
	name string
}

// Список вещей в каком либо месте (в инвентаре, на столе, на полу)
type ItemsCollection struct {
	Items []ItemId
}

type ItemLibrary struct {
	idToData map[ItemId]Item //Карта (словарь) предметов
	ids      []ItemId        //Для некоторых проверок,
}
