package main

import (
	"context"
	"es/bank"
	"es/eventstore"
	"github.com/shopspring/decimal"
	"log"
	"sync"
)

func main() {
	once := sync.Once{}
	once.Do(func() {
		eventstore.Init()
	})
	defer eventstore.Close()

	newID, err := bank.CreateAccount(context.Background(), "Fubuki")
	if err != nil {
		panic(err)
	}
	log.Printf("New bank account id: %s", newID)
	err = bank.Deposit(context.Background(), newID, decimal.NewFromInt(100))
	if err != nil {
		panic(err)
	}
	err = bank.Withdraw(context.Background(), newID, decimal.NewFromFloat(50.50))
	if err != nil {
		panic(err)
	}
	err = bank.ProcessEscrow(context.Background(), newID, []decimal.Decimal{
		decimal.NewFromFloat(23.00),
		decimal.NewFromFloat(100.00),
		decimal.NewFromFloat(300.00),
		decimal.NewFromFloat(5.00),
	})
	if err != nil {
		panic(err)
	}

	accountAggregate, err := bank.FindAccount(context.Background(), newID)
	if err != nil {
		panic(err)
	}
	log.Println(accountAggregate)

	err = bank.Withdraw(context.Background(), newID, decimal.NewFromFloat(1000000))
	if err != nil {
		log.Printf("SHOULD FAIL: %v", err)
	}

	err = bank.ChangeName(context.Background(), newID, "Shirakami Fubuki")
	if err != nil {
		panic(err)
	}

	accountAggregate, err = bank.FindAccount(context.Background(), newID)
	if err != nil {
		panic(err)
	}
	log.Println(accountAggregate)
}
