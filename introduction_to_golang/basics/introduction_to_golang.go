package main

import (
	"fmt"
	"unicode/utf8"
)

// UserID - новый тип
type UserID int

func main() {

	/*---------------------------------------------------*/
	/*                  ПЕРЕМЕННЫЕ                       */
	/*---------------------------------------------------*/

	// значение по умолчанию
	var num0 int // num0 = 0

	// значение при инициализации
	var num1 int = 1

	// пропуск типа
	var num2 = 20

	// короткое объявление переменной, но только для новых переменных!
	num := 30

	num++ // num += 1 тоже самое
	// ++num нет

	// userIndex - принятый стиль
	userIndex := 10
	// user_index - не принято

	// объявление нескольких переменных
	var weight, height int = 10, 20

	// присваивание в существующие переменные
	weight, height = 11, 21

	// при коротком присваивании хотя-бы одна переменная должна быть новой!
	weight, age := 12, 22

	fmt.Println("\nVars:")
	fmt.Println(num0, num1, num2, num, userIndex, weight, height, age)

	// int - платформозависимый тип, 32/64 бита
	var i int = 10

	// автоматически выбранный int
	var autoInt = -10

	// int8, int16, int32, int64
	var bigInt int64 = 1<<32 - 1

	// платформозависимый тип, 32/64
	var unsignedInt uint = 100500

	// uint8, unit16, uint32, unit64
	var unsignedBigInt uint64 = 1<<64 - 1

	// float32, float64
	var pi float32 = 3.141
	var e = 2.718
	goldenRatio := 1.618

	// bool
	var b bool // false по умолчанию
	var isOk bool = true
	var success = true
	cond := true

	// complex64, complex128
	var c complex128 = -1.1 + 7.12i
	c2 := -1.1 + 7.12i

	fmt.Println(i, autoInt, bigInt, unsignedInt, unsignedBigInt, pi, e, goldenRatio, b, isOk, success, cond, c, c2)

	/*---------------------------------------------------*/
	/*                    СТРОКИ                         */
	/*---------------------------------------------------*/

	// пустая строка по-умолчанию
	var str string

	// со спец символами
	var hello string = "Привет\n\t"

	// без спец символов
	var world string = `Мир\n\t`

	// UTF-8 из коробки
	var helloWorld = "Привет, Мир!"
	hi := "你好，世界"

	// одинарные кавычки для байт (uint8)
	var rawBinary byte = '\x27'
	// rawBinary = 39

	// rune (uint32) для UTF-8 символов
	var someChinese rune = '茶'
	// someChinese = 33590

	// конкатенация строк
	helloWorld = "Привет Мир"
	andGoodMorning := helloWorld + " и доброе утро!"
	// andGoodMorning = "Привет Мир и доброе утро!""

	// строки неизменяемы!

	// получение длины строки
	byteLen := len(helloWorld)                    // 19 байт
	symbols := utf8.RuneCountInString(helloWorld) // 10 рун

	// получение подстроки, в байтах, не символах! (оператор slice)
	hello = helloWorld[:12] // Привет, 0-11 байты
	H := helloWorld[0]      // byte, 72, не "П"

	// конвертация в слайс байт и обратно
	byteString := []byte(helloWorld)
	// byteString = [208 159 209 128 208 184 208 178 208 181 209 130 32 208 156 208 184 209 128]
	helloWorld = string(byteString)
	// helloWorld = "Привет Мир"

	fmt.Println("\nStrings:")
	fmt.Println(str, hello, ";", H, ";", world, ";", helloWorld, ";", hi, ";", rawBinary, someChinese)
	fmt.Println(andGoodMorning, ";", byteLen, ";", symbols, ";", byteString, ";", helloWorld)

	/*---------------------------------------------------*/
	/*                    МАССИВЫ                        */
	/*---------------------------------------------------*/

	// размер массива является частью его типа
	// инициализация значениями по-умолчанию
	var a1 [3]int // a1 = [0,0,0]
	fmt.Println("\nArray:")
	fmt.Println("short", a1)     // [0 0 0]
	fmt.Printf("full %#v\n", a1) // [3]int{0, 0, 0}

	const size = 2
	var a2 [2 * size]bool // a2 = [false,false,false,false]

	// определение размера при объявлении
	a3 := [...]int{1, 2, 3} // a3 = [1 2 3]

	fmt.Println(a2, a3)

	/*---------------------------------------------------*/
	/*                     СЛАЙСЫ                        */
	/*---------------------------------------------------*/

	// создание
	var buf0 []int             // len=0, cap=0
	buf1 := []int{}            // len=0, cap=0
	buf2 := []int{42}          // len=1, cap=1
	buf3 := make([]int, 0)     // len=0, cap=0
	buf4 := make([]int, 5)     // len=5, cap=5
	buf5 := make([]int, 5, 10) // len=5, cap=10

	// обращение к элементам
	someInt := buf2[0] // someInt = 42

	// добавление элементов
	var buf []int            // len=0, cap=0
	buf = append(buf, 9, 10) // len=2, cap=2
	buf = append(buf, 12)    // len=3, cap=4
	// на i-м append-e при len == cap, cap = 2^i

	// добавление другого слайса
	otherBuf := make([]int, 3)     // otherBuf = [0,0,0]
	buf = append(buf, otherBuf...) // len=6, cap=8

	// просмотр информации о слайсе
	var bufLen, bufCap int = len(buf), cap(buf)

	bufMain := []int{1, 2, 3, 4, 5}

	// получение среза, указывающего на ту же память
	sl1 := bufMain[1:4] // [2, 3, 4]
	sl2 := bufMain[:2]  // [1, 2]
	sl3 := bufMain[2:]  // [3, 4, 5]

	newBuf := bufMain[:] // [1, 2, 3, 4, 5]
	newBuf[0] = 9
	// bufMain = [9, 2, 3, 4, 5], т.к. та же память

	// newBuf теперь указывает на другие данные
	newBuf = append(newBuf, 6)

	// bufMain    = [9, 2, 3, 4, 5], не изменился
	// newBuf = [1, 2, 3, 4, 5, 6], изменился
	newBuf[0] = 1

	// копирование одного слайса в другой
	var emptyBuf []int // len=0, cap=0
	// неправильно!
	copied := copy(emptyBuf, bufMain) // copied = 0, скопирует меньшее (по len) из 2-х слайсов
	// правильно!
	newBuf = make([]int, len(bufMain), len(bufMain))
	copy(newBuf, bufMain)

	// можно копировать в часть существующего слайса
	ints := []int{1, 2, 3, 4}
	copy(ints[1:3], []int{5, 6}) // ints = [1, 5, 6, 4]

	fmt.Println("\nSlice:")
	println(buf0, buf1, buf2, buf3, buf4, buf5)
	fmt.Println(someInt, buf, otherBuf, bufLen, bufCap, bufMain, sl1, sl2, sl3)
	fmt.Println(bufMain, newBuf, copied, emptyBuf, ints)

	/*---------------------------------------------------*/
	/*                       MAP                         */
	/*---------------------------------------------------*/

	// инициализация при создании
	var user map[string]string = map[string]string{
		"name":     "Eugene",
		"lastName": "Goncharenko",
	}

	// сразу с нужной ёмкостью
	profile := make(map[string]string, 10) // profile = map[]

	// количество элементов
	mapLength := len(user) // mapLength = 2

	// если ключа нет - вернёт значение по умолчанию для типа
	mName := user["middleName"] // mName = " "

	// проверка на существование ключа
	mName, mNameExist := user["middleName"] // mName = " ", mNameExist = false

	// пустая переменная - только проверяем что ключ есть
	_, mNameExist2 := user["middleName"] // mNameExist2 = false

	// удаление ключа
	delete(user, "lastName")

	fmt.Println("\nMap:")
	fmt.Printf("%d %+v\n", mapLength, profile)
	fmt.Println(mName, mNameExist, mNameExist2)
	fmt.Printf("%#v\n", user)

	/*---------------------------------------------------*/
	/*                     УКАЗАТЕЛИ                     */
	/*---------------------------------------------------*/

	objA := 2
	objB := &objA // тоже самое, что: var objB *int = &objA
	*objB = 3     // objA = 3
	objC := &objA // новый указатель на переменную a

	// получение указателя на переменную типа int
	// инициализировано значением по-умолчанию
	objD := new(int)
	*objD = 12
	*objC = *objD // переприсваивание значений, те objC = 12 -> objA = 12
	*objD = 13    // objC и objA не изменились

	objC = objD // теперь objC указывает туда же, куда objD
	*objC = 14  // objC = 14 -> d = 14, objA = 12

	/*---------------------------------------------------*/
	/*                     ТИПЫ                          */
	/*---------------------------------------------------*/

	idx := 1
	var uid UserID = 42

	// даже если базовый тип одинаковый, разные типы несовместимы
	// cannot use uid (type UserID) as type int64 in assignment
	// myID := idx

	myID := UserID(idx)

	fmt.Println("\nTypes:")
	println(uid, myID)
}
