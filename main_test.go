package main

import (
	"testing"
)

type gameCase struct {
	step    int
	command string
	answer  string
}

var game0cases = [][]gameCase{
	[]gameCase{
		{1, "осмотреться", "ты находишься на кухне, надо собрать рюкзак и идти в универ. можно пройти - коридор"}, // действие осмотреться
		{2, "идти коридор", "ничего интересного. можно пройти - кухня, комната, улица"},                           // действие идти
		{3, "идти комната", "ты в своей комнате. можно пройти - коридор"},
		{4, "осмотреться", "на столе: ключи, конспекты, рюкзак. можно пройти - коридор"},
		{5, "взять рюкзак", "предмет добавлен в инвентарь: рюкзак"}, // действие взять
		{6, "взять ключи", "предмет добавлен в инвентарь: ключи"},   // действие взять
		{7, "взять конспекты", "предмет добавлен в инвентарь: конспекты"},
		{8, "идти коридор", "ничего интересного. можно пройти - кухня, комната, улица"},
		{9, "идти улица", "на улице весна. можно пройти - домой"},
	},

	[]gameCase{
		{1, "осмотреться", "ты находишься на кухне, надо собрать рюкзак и идти в универ. можно пройти - коридор"},
		{2, "завтракать", "неизвестная команда"},  // придёт топать в универ голодным :(
		{3, "идти комната", "нет пути в комната"}, // через стены ходить нельзя
		{4, "идти коридор", "ничего интересного. можно пройти - кухня, комната, улица"},
		{5, "идти комната", "ты в своей комнате. можно пройти - коридор"},
		{6, "осмотреться", "на столе: ключи, конспекты, рюкзак. можно пройти - коридор"},
		{7, "взять рюкзак", "предмет добавлен в инвентарь: рюкзак"},              // действие взять
		{8, "осмотреться", "на столе: ключи, конспекты. можно пройти - коридор"}, // состояние изменилось
		{9, "взять ключи", "предмет добавлен в инвентарь: ключи"},
		{10, "взять телефон", "нет такого"},                                // неизвестный предмет
		{11, "взять ключи", "нет такого"},                                  // предмента уже нет в комнатеы - мы его взяли
		{12, "осмотреться", "на столе: конспекты. можно пройти - коридор"}, // состояние изменилось
		{13, "взять конспекты", "предмет добавлен в инвентарь: конспекты"},
		{14, "осмотреться", "пустая комната. можно пройти - коридор"}, // состояние изменилось
		{15, "идти коридор", "ничего интересного. можно пройти - кухня, комната, улица"},
		{16, "идти кухня", "кухня, ничего интересного. можно пройти - коридор"},
		{17, "осмотреться", "ты находишься на кухне, надо собрать рюкзак и идти в универ. можно пройти - коридор"}, // состояние изменилось
		{18, "идти коридор", "ничего интересного. можно пройти - кухня, комната, улица"},
		{19, "идти улица", "на улице весна. можно пройти - домой"},
	},
}

func TestGame0(t *testing.T) {
	for caseNum, commands := range game0cases {
		//log.Printf("caseNum=%v, commands=%v", caseNum, commands)
		initAll()
		//log.Printf("init done")
		for _, item := range commands {
			//log.Printf("item=%v", item)
			answer := handleCommand(item.command)
			if answer != item.answer {
				//log.Printf("answer != item.answer")
				t.Error("case:", caseNum, item.step,
					"\n\tcmd:", item.command,
					"\n\tresult:  ", answer,
					"\n\texpected:", item.answer)
			}
		}
	}
}
