package main

import (
	"./items"
	"./players"
	"./rooms"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var player players.PlayerType

// Итак, мы хотим сделать игру с вводом заданий.
//
// Ключевые команды у нас:
//
// - осмотреться - получить сведения о том, где герой, что он видит (название пространства, набор предметов внутри пространства, возможные ходы)
//
// - идти % - попытаться переместить героя в новое пространство. Проходит проверку доступных зон перехода
//
// - взять % - попытаться переместить объект из пространства в котором находится герой в инвентарь. Проходит проверку наличия предмета в пространстве, помещает предмет в инвентарь, исключает предмет из пространства (возможно, следует действовать исходя из того, что не объекты перемещаются, а меняется их свойство, но для большого числа объектов это затруднительно)
//
// Основные объекты у нас:
//
// - персонаж - имеет свойства: расположение, инвентарь (массив разнородных объектов к которым применимы команды взять, получить имя);
//
// - помещение - имеет свойства: название, инвентарь (стол)(возможно, могут быть и другие), допустимые направления перехода (для каждого помещения свои, т.к. обратных переходов может и не быть, например: вышли и захлопнули дверь;
//
// - объекты - интерфейсные типы, к которым удовлетворяющие командам взять, получить имя;
//
// схема хранения:
// - помещения - карта, ключ - код помещения. (словари имя помещения - код помещения храним отдельной картой, чтобы можно было назвать одно помещения несколькими способами).
// -
func main() {
	/*
		в этой функции можно ничего не писать,
		но тогда у вас не будет работать через go run main.go
		очень круто будет сделать построчный ввод команд тут, хотя это и не требуется по заданию
	*/

	defer func() {
		err := recover()
		if err != nil {
			log.Fatalf("сожалеем, во время работы программы произошла ошибка. работа программы завершена")
		}
	}()

	initAll()
	listenAndServe()
}

// Если конфигурация противоречива - выйдем с паникой и пояснением
func initAll() {
	/*
		эта функция инициализирует игровой мир - все локации
		если что-то было - оно корректно перезатирается
	*/

	defer func() {
		err := recover()
		if err != nil {
			log.Fatalf("настройки инициализации содержат некорректные данные. работа программы завершена")
		}
	}()

	//Создаем предметы
	allItems := items.Init() //Затрет все что было
	allItems.AddItemToLibrary(1, "ключи")
	allItems.AddItemToLibrary(2, "конспекты")
	allItems.AddItemToLibrary(3, "рюкзак")

	//log.Printf("all items=%+v", allItems)

	//У нас есть следующие сценарии:
	//1. Пользователь приходит в пространство, moveToRoomResultText - "ты в своей комнате", "ничего интересного", "кухня, ничего интересного"
	//2. Пользозователь осматривается в пустом пространстве - имя пустого пространства emptyLookAroundMsg - "пустая комната, пустой корридор, ты находишься на кухне, надо собрать рюкзак и идти в универ"
	//3. Пользователь осматривается в заполненном пространстве - имя инвентаря помещения invenoryNameIfNotEmpty, например: "на столе: "

	//Создаем комнаты
	allRooms := rooms.Init() //Затрет все что было
	room := allRooms.AddRoom(1, "кухня")
	room.SetAsDefaultForPlayer() //Либо можем вызвать player.SetCurrentRoom(1)
	room.SetMoveToRoomResultText("кухня, ничего интересного")
	room.SetEmptyLookAroundMsg("ты находишься на кухне, надо собрать рюкзак и идти в универ")

	room = allRooms.AddRoom(2, "коридор")
	room.SetMoveToRoomResultText("ничего интересного")
	room.SetEmptyLookAroundMsg("ты в коридоре, тут ничего интересного")

	room = allRooms.AddRoom(3, "комната")
	room.SetMoveToRoomResultText("ты в своей комнате")
	room.SetEmptyLookAroundMsg("пустая комната")
	room.SetInvenoryNameIfNotEmpty("на столе")
	room.AddItems(1, 2, 3)

	room = allRooms.AddRoom(4, "улица")
	room.SetMoveToRoomResultText("на улице весна")
	room.SetEmptyLookAroundMsg("на улице весна")

	room = allRooms.AddRoom(5, "домой")
	room.SetRoomAsJumpTo(2) //Домой приведет нас в коридор

	//Создаем пути
	allRooms.AddWays(1, 2, true)
	allRooms.AddWays(2, 3, true)
	allRooms.AddWays(2, 4, false)
	allRooms.AddWays(4, 5, false) //Домой это коридор

	//Создаем игрока
	player = players.InitPlayer()
	//player.SetCurrentRoom(1)

}

// Оба решения одинаково дрянно обрабатывают редактирование вводимого из консоли текста, через раз выдавая артефакты
func listenAndServe_() {
	for {
		fmt.Print("Введите команду: ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("read line: %s\n", line)
		fmt.Printf("Результат:%s\n", handleCommand(strings.ReplaceAll(line, "\n", "")))
	}
}

func listenAndServe() {
	for {
		fmt.Print("Введите команду: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Результат:%s\n", handleCommand(strings.ReplaceAll(scanner.Text(), "\n", "")))
	}
}

func handleCommand(command string) string {
	/*
		данная функция принимает команду от "пользователя"
		и наверняка вызывает какой-то другой метод или функцию у "мира" - списка комнат
	*/

	commandArgs := strings.Split(command, " ")

	//log.Printf("len(commandArgs) = %v", len(commandArgs))
	//log.Printf("commandArgs = %v", commandArgs)
	//log.Printf("commandArgs[0] = %v", commandArgs[0])
	//log.Printf("rune(restart) = %v", []rune("restart"))
	//log.Printf("rune(commandArgs[0]) = %v", []rune(commandArgs[0]))

	switch len(commandArgs) {
	default:
		return "неизвестная команда"
	case 1: //Пришла команда из одного слова
		switch commandArgs[0] {
		default:
			return "неизвестная команда"
		case "restart":
			initAll()
			return "вернулись к началу"
		case "осмотреться":
			return player.LookAround().Msg
		}
	case 2: //Пришла команда из двух слов
		switch commandArgs[0] {
		default:
			return "неизвестная команда"
		case "взять":
			return player.GetItem(commandArgs[1]).Msg
		case "идти":
			return player.MoveToRoom(commandArgs[1]).Msg
		}
	}
}
