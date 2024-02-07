package items

import (
	"../tools"
	"log"
	"strings"
)

// Метод получения содержимого набора предметов в виде строки
func (col *ItemsCollection) String() string {
	var strsTmp []string
	for _, itemId := range col.Items {
		strsTmp = append(strsTmp, lib.idToData[itemId].name)
	}
	if len(strsTmp) > 0 {
		return strings.Join(strsTmp, ", ")
	} else {
		return "ничего"
	}
}

// Затрем все словари которые были
func Init() *ItemLibrary {
	lib.idToData = make(map[ItemId]Item)
	lib.ids = []ItemId{}
	return &lib
}

// Добавляем предметы в словарь
func (lib *ItemLibrary) AddItemToLibrary(itemId ItemId, itemName string) {
	if lib.idToData == nil {
		panic("Попытка вызова AddItemToLibrary без вызова InitItems. Карта itemLibrary.idToData не была инициализирована")
	}
	mustNotExits(itemId)

	lib.idToData[itemId] = Item{
		id:   itemId,
		name: itemName,
	}

	lib.ids = append(lib.ids, itemId) //Заполняем для простоты контроля при обращении
}

func CheckAllItemsExists(ids []ItemId) (allExits bool, badItems []ItemId) {
	return tools.CheckObjectsInStice(lib.ids, ids)
}

func (id ItemId) GetName() string {
	if data, ok := lib.idToData[id]; ok {
		return data.name
	}
	return "неопознанный предмет"
}

// Проверяем есть ли в комнате предмет с таким названием и возвращаем его id или nil
func (col *ItemsCollection) GetIdByName(itemName string) *ItemId {
	for _, itemId := range col.Items {
		if itemId.GetName() == itemName {
			return &itemId
		}
	}
	return nil
}

// Проверяем существует ли комната с таким id, если существует - паникуем
func mustNotExits(id ItemId) {
	if _, exists := lib.idToData[id]; exists {
		log.Panicf("Попытка создания предмета с id=%v, а такой уже существует", id)
	}
}
