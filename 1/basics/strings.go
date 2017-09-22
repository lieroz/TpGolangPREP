package main

import "fmt"

func main() {
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

	// rune (uint32) для UTF-8 символов
	var someChinese rune = '茶'

	helloWorld := "Привет Мир"
	// конкатенация строк
	andGoodMorning := helloWorld + " и доброе утро!"

	// строки неизменяемы
	// cannot assign to helloWorld[0]
	helloWorld[0] = 72

	// получение длины строки
	byteLen := len(helloWorld)                    // 19 байт
	symbols := utf8.RuneCountInString(helloWorld) // 10 рун

	// получение подстроки, в байтах, не символах!
	hello := helloWorld[:12] // Привет, 0-11 байты
	H := helloWorld[0]       // byte, 72, не "П"

	// конвертация в слайс байт и обратно
	byteString = []byte(helloWorld)
	helloWorld = string(byteString)
}
