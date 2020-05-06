package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Object struct { // то на чем могут лежать предмнеьты стол стул
	name    string
	isItems map[string]bool // изъяты ли
	items   []string        // сами предметы
}

type GoStatus struct {
	status bool
	msg    string
}

type Location struct {
	name        string
	destination string              // префикс при переходе
	ways        map[string]GoStatus // можно ли пройти в локацию
	wayNames    []string            // список локаций
	objectNames []string            // имена объектов. Нужен для сохр порядка
	objects     map[string]Object   // сами объекты
	itemsOnObj  map[string]string   // нужно для определения объекта при взять хендлере
	clothes     map[string]bool     // Вещи, которые можно надеть
	prefixWatch string
}

type Status struct {
	currentLocation string          // текущая локация
	clothes         map[string]bool // какие на нас шмотки. онли рюкзак зато без ифов)
	message         string          // показываем только на кухне. такие тесты...
}

var locations map[string]Location              // переход из названия в данные о локации
var state Status                               // текущая инфа об игроке
var backpacker map[string]bool                 // что в рюкзаке
var doorPosition bool                          // состояние двери
var usable map[string]map[string]func() string // отношения используемости объектов друг к другу + вызов обработчика
var emptyRoomMap map[string]string             // как звучит пустая комната))

const (
	Kitchen            = "кухня"
	Home               = "домой"
	Bedroom            = "комната"
	Corridor           = "коридор"
	Street             = "улица"
	NothingInteresting = "ничего интересного"
	KitchenPrefix      = "ты находишься на кухне, "
	BedroomPrefix      = "ты в своей комнате"
	CorridorPrefix     = NothingInteresting
	StreetPrefix       = "на улице весна"

	EmptyRoom    = "пустая комната"
	FirstMessage = "на столе: чай, надо собрать рюкзак и идти в универ"
	LastMessage  = "на столе: чай, надо идти в универ"

	OnTheTable = "на столе: "
	OnTheChair = "на стуле: "

	Open  = true
	Close = false

	UnknownCommand = "неизвестная команда"
	FictitiousItem = "нет такого"

	NothingToApply = "не к чему применить"
	AbsentItem     = "нет предмета в инвентаре - "

	NoBackpack   = "некуда класть"
	AlreadyTaken = " уже надет"
	NotFound     = " не найден"
	NoWayTo      = "нет пути в "
	Taken        = "вы надели: "
	AddItem      = "предмет добавлен в инвентарь: "

	CloseDoor = "дверь закрыта"
	OpenDoor  = "дверь открыта"

	LookAround = "осмотреться"
	Go         = "идти"
	PutOn      = "надеть"
	Take       = "взять"
	Apply      = "применить"

	Door      = "дверь"
	Backpack  = "рюкзак"
	Keys      = "ключи"
	Abstracts = "конспекты"

	CanPass = "можно пройти - "
)

func doorView() string {
	if doorPosition == Open {
		return OpenDoor
	}
	return CloseDoor
}

func doorStatusChange() string { // вызывается при применении ключ в дверь
	if state.currentLocation != Corridor {
		return NothingToApply
	}
	doorPosition = !doorPosition
	if doorPosition == Open { // при открытии двери открывается локация улица из корридора
		locations[Corridor].ways[Street] = GoStatus{true, OpenDoor}
	} else { // при закрытии двери закрывается локация улица из корридора
		locations[Corridor].ways[Street] = GoStatus{false, CloseDoor}
	}
	return doorView()
}

func kudagoFromcur() string {
	return CanPass + strings.Join(locations[state.currentLocation].wayNames, ", ")
}

// *
// *************************
// * WATCH HANDLER METHODS *
// *************************
// *

func watchObjectsHandler() string { // что мы видим в комнате + сообщение для пользователя
	result := make([]string, 0) // объект + то, что на нем мы видим
	for _, objName := range locations[state.currentLocation].objectNames {
		what_i_watch := make([]string, 0) // что мы видим на текущем объекте
		for _, itemName := range locations[state.currentLocation].objects[objName].items {
			if locations[state.currentLocation].objects[objName].isItems[itemName] == true {
				what_i_watch = append(what_i_watch, itemName)
			}
		}
		if len(what_i_watch) != 0 {
			result = append(result, objName+strings.Join(what_i_watch, ", "))
		}
	}
	if len(result) == 0 {
		if state.currentLocation == Kitchen { // такие тесты.
			return state.message
		}
		return emptyRoomMap[state.currentLocation]
	}
	return strings.Join(result, ", ")
}

// *
// *******************
// * HANDLE COMMANDS *
// *******************
// *

func usedHandler(item1 string, item2 string) string { // item1 испльзуется по отношению к item2
	isItem, ok := backpacker[item1]
	if !state.clothes[Backpack] || !ok || !isItem { // используются вещи тлько из инвентаря
		return AbsentItem + item1
	}
	_, ok = usable[item1]
	if !ok {
		return NothingToApply
	}
	usableFunc, ok := usable[item1][item2]
	if !ok {
		return NothingToApply
	}
	return usableFunc()
}

func getHandler(item string) string { // взять в рюкзак
	if state.clothes[Backpack] {
		isItem, ok := backpacker[item]
		if !ok || isItem {
			// если предмет нельзя брать или он уже взят
			return FictitiousItem
		}
		obj, ok := locations[state.currentLocation].itemsOnObj[item]
		if !ok {
			return FictitiousItem
		}
		locations[state.currentLocation].objects[obj].isItems[item] = false
		backpacker[item] = true

		if backpacker[Keys] && backpacker[Abstracts] {
			state.message = LastMessage
		}
		return AddItem + item
	}
	return NoBackpack
}

func takeHandler(item string) string { // надеть (рюкзак)
	if state.clothes[item] { // если мы уже надели вещицу
		return item + AlreadyTaken
	}
	isItem, ok := locations[state.currentLocation].clothes[item]
	if !isItem || !ok { // либо данная вещь отсутствует либо не надеваемая
		return item + NotFound
	}
	state.clothes[item] = true                                          // надели)
	locations[state.currentLocation].clothes[item] = false              // больше ее в локации нет(она ведь на нас)
	obj := locations[state.currentLocation].itemsOnObj[item]            // узнаем, на каком объекте она лежала
	locations[state.currentLocation].objects[obj].isItems[item] = false // убираем ее оттуда
	return Taken + item
}

func goHandler(location string) string { // вызывается при перемещении среди локаций
	_, ok := locations[location]
	if !ok { // данной локации не существует
		return UnknownCommand
	}
	goStatus, ok := locations[state.currentLocation].ways[location] // пути в локацию нет
	if !ok {
		return NoWayTo + location
	}
	if !goStatus.status { // если путь есть, но по каким-то причинам он закрыт возвращаем сообщение об ошибке
		return goStatus.msg
	}

	state.currentLocation = location // меняем локацию и возвращаем ответ
	return strings.Join([]string{locations[state.currentLocation].destination, kudagoFromcur()}, ". ")
}

func watchHandler() string { // при осмотреться что мы видим
	objects := watchObjectsHandler() // какие объекты + сообщения
	kudago := kudagoFromcur()        // куда можно пройти
	return strings.Join([]string{locations[state.currentLocation].prefixWatch + objects, kudago}, ". ")
}

func handleCommand(command string) string { // общий парсер команд
	words := strings.Split(command, " ")
	if words[0] == LookAround {
		if len(words) == 1 {
			return watchHandler()
		}
	} else if words[0] == Go {
		if len(words) == 2 {
			return goHandler(words[1])
		}
	} else if words[0] == PutOn {
		if len(words) == 2 {
			return takeHandler(words[1])
		}
	} else if words[0] == Take {
		if len(words) == 2 {
			return getHandler(words[1])
		}
	} else if words[0] == Apply {
		if len(words) == 3 {
			return usedHandler(words[1], words[2])
		}
	}
	return UnknownCommand
}

// *
// ***************************
// * CREATE LOCATION METHODS *
// ***************************
// *

func getRoomWays() []string { // пути из комнаты
	res := make([]string, 0)
	res = append(res, Corridor)
	return res
}

func getCorridorWays() []string { // пути из корридора
	res := make([]string, 0)
	res = append(res, Kitchen)
	res = append(res, Bedroom)
	res = append(res, Street)
	return res
}

func getKitchenWays() []string { // пути из кухни
	res := make([]string, 0)
	res = append(res, Corridor)
	return res
}

func getStreetWays() []string { // пути из улицы
	res := make([]string, 0)
	res = append(res, Home)
	return res
}

func getStreetObjects() []Object { // объекты улицы
	return make([]Object, 0)
}

func getCorridorObjects() []Object { // объекты корридора
	return make([]Object, 0)
}

func getKitchenObjects() []Object { // объекты кухни
	return make([]Object, 0)
}

func getRoomObjects() []Object { // объекты комнаты
	table := Object{OnTheTable, make(map[string]bool), make([]string, 0)}
	table.isItems[Keys] = true
	table.isItems[Abstracts] = true
	table.items = append(table.items, Keys)
	table.items = append(table.items, Abstracts)
	//создали стол, положили на него ключи и конспекты

	chair := Object{OnTheChair, make(map[string]bool), make([]string, 0)}
	chair.isItems[Backpack] = true
	chair.items = append(chair.items, Backpack)
	// создали кресло, положили на него рюкзак

	roomObjs := make([]Object, 0)
	roomObjs = append(roomObjs, table)
	roomObjs = append(roomObjs, chair)
	return roomObjs // вернули собственно все объекты
}

func createLocation(name string, dest string, objects []Object, prefixWatch string, wayNames []string) {
	objNames := make([]string, 0) // генерируем из массива объектов массив имен объектов
	for _, val := range objects {
		objNames = append(objNames, val.name)
	}
	locations[name] = Location{ // создаем локацию
		name,
		dest,

		make(map[string]GoStatus),
		wayNames,

		objNames,
		make(map[string]Object),

		make(map[string]string),
		make(map[string]bool),

		prefixWatch,
	}

	for _, val := range objects { // налаживаем связь имя объекта --> объект для обращения
		locations[name].objects[val.name] = val
	}

	for _, val := range wayNames { // заполняем пути
		goStatus := GoStatus{true, ""}
		if val == Street {
			goStatus.status = false
			goStatus.msg = CloseDoor
		}
		locations[name].ways[val] = goStatus
	}
}

// *
// *************
// * INIT GAME *
// *************
// *

func createEmptyRoomMap() { // инициализируем пустые локации
	emptyRoomMap = make(map[string]string)
	emptyRoomMap[Bedroom] = EmptyRoom
	emptyRoomMap[Corridor] = NothingInteresting
	emptyRoomMap[Street] = NothingInteresting
	emptyRoomMap[Kitchen] = NothingInteresting
}

func createUsableMap() { // инициадлизируем отношения применимости объектов
	usable = make(map[string]map[string]func() string)
	usable[Keys] = make(map[string]func() string) // ключ
	usable[Keys][Door] = doorStatusChange         // ставим функцию обработчик при примернинии ключа к двери
}

func createLocations() { // создаем локации
	locations = make(map[string]Location)
	createLocation(
		Kitchen,
		Kitchen+", "+NothingInteresting,
		getKitchenObjects(),
		KitchenPrefix,
		getKitchenWays(),
	)

	createLocation(
		Bedroom,
		BedroomPrefix,
		getRoomObjects(),
		"",
		getRoomWays(),
	)
	createLocation(
		Corridor,
		CorridorPrefix,
		getCorridorObjects(),
		"",
		getCorridorWays(),
	)
	createLocation(
		Street,
		StreetPrefix,
		getStreetObjects(),
		"",
		getStreetWays(),
	)

	locations[Bedroom].itemsOnObj[Keys] = OnTheTable
	locations[Bedroom].itemsOnObj[Abstracts] = OnTheTable
	locations[Bedroom].itemsOnObj[Backpack] = OnTheChair
	locations[Bedroom].clothes[Backpack] = true
}

func createBackpack() { // создаем рюкзак
	backpacker = make(map[string]bool)
	backpacker[Keys] = false      // вещь может быть в рюкзаке, но сейчас ее нет
	backpacker[Abstracts] = false // аналогично
}

func createPlayer() { // создаем игрока
	state = Status{
		Kitchen,
		make(map[string]bool),
		FirstMessage,
	}
}

func initGame() {
	doorPosition = Close
	createPlayer()
	createBackpack()
	createLocations()
	createUsableMap()
	createEmptyRoomMap()
}

func main() {
	fmt.Println("")
	fmt.Println("\t*****************")
	fmt.Println("\t* Game started! *")
	fmt.Println("\t*****************")
	fmt.Println("\t|")

	initGame()
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1]
		if text == "exit" {
			break
		}
		fmt.Println("\t|>\t [ in] :: " + text)
		fmt.Println("\t|>\t [out] :: " + handleCommand(text))
	}

	fmt.Println("\t|")
	fmt.Println("\t*****************")
	fmt.Println("\t* Game stopped! *")
	fmt.Println("\t*****************")
	fmt.Println("")

}
