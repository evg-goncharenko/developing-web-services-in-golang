package main

import (
	"fmt"
)

type Payer interface {
	Pay(int) error // возвращает ошибку (error)
}

type Wallet struct {
	Cash int
}

func (w *Wallet) Pay(amount int) error {
	if w.Cash < amount {
		return fmt.Errorf("Не хватает денег в кошельке")
	}
	w.Cash -= amount
	return nil
}

// функция Pay() удовлетворяет интерфейсу Payer
// неявное удовлетворение интерфейсу, тк нет упоминания о Payer

// при использовании интерфейса указатель не нужен
func Buy(p Payer) {
	err := p.Pay(10)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Спасибо за покупку через %T\n\n", p)
}

func main() {
	myWallet := &Wallet{Cash: 100}
	Buy(myWallet)
}
