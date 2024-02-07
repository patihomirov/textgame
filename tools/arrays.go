package tools

// Метод проверки наличия элемента в слайсе. Применим только к типам comparable
func CheckObjectInStice[T comparable](slice []T, object T) bool {
	for _, val := range slice {
		if val == object {
			return true
		}
	}
	return false
}

// Метод проверки наличия всех элементов в слайсе. Применим только к типам comparable. Если не что-то не совпало - вернем ошибку
func CheckObjectsInStice[T comparable](slice []T, objects []T) (allExits bool, badItems []T) {
	for _, object := range objects {
		if !CheckObjectInStice(slice, object) {
			badItems = append(badItems, object)
		}
	}
	allExits = len(badItems) == 0 //false если есть не найденные
	return
}

// Метод проверки наличия элемента в слайсе. Применим только к типам comparable
// Возвращаем новый слайс и булевый признак был ли элемент найден в слайсе перед удалением
func RemoveObjectFromStice[T comparable](sliceIn []T, object T) (sliceOut []T, founded bool) {
	index := -1

	// Находим индекс числа, которое нужно удалить
	for i, objIter := range sliceIn {
		if objIter == object {
			index = i
			break
		}
	}

	// Если число найдено, удаляем его из среза
	if index != -1 {
		founded = true
		sliceOut = append(sliceIn[:index], sliceIn[index+1:]...)
	} else {
		sliceOut = sliceIn //Иначе вовзращаем что было
	}

	return
}
