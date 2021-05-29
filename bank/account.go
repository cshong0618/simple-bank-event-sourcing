package bank

import (
	"errors"
	"es/eventstore"
	"github.com/shopspring/decimal"
)

type AccountAggregate struct {
	ID      string
	Name    string
	Balance decimal.Decimal
	Version uint64

	events []eventstore.Event
}

func (a *AccountAggregate) BuildFromEvents(events []eventstore.Event) error {
	for _, event := range events {
		switch event.Type() {
		case AccountCreatedEventType:
			{
				e := event.(*AccountCreatedEvent)
				a.ID = e.ID
				a.Name = e.Name
			}
		case MoneyDepositedEventType:
			{
				e := event.(*MoneyDepositedEvent)
				a.Balance = a.Balance.Add(e.Amount)
			}
		case MoneyWithdrawnEventType:
			{
				e := event.(*MoneyWithdrawnEvent)
				a.Balance = a.Balance.Sub(e.Amount)
			}
		case NameChangedEventType:
			{
				e := event.(*NameChangedEvent)
				a.Name = e.Name
			}
		default:
			return errors.New("invalid event type")
		}
		a.Version = event.Version()
	}

	return nil
}

func (a *AccountAggregate) Deposit(amount decimal.Decimal) {
	a.events = append(a.events, MoneyDepositedEvent{
		Amount: amount,
	})
	a.Balance = a.Balance.Add(amount)
}

func (a *AccountAggregate) Withdraw(amount decimal.Decimal) error {
	if amount.GreaterThan(a.Balance) {
		return errors.New("not enough balance to withdraw")
	}

	a.events = append(a.events, MoneyWithdrawnEvent{
		Amount: amount,
	})
	a.Balance = a.Balance.Sub(amount)

	return nil
}

func (a *AccountAggregate) ChangeName(newName string) {
	a.events = append(a.events, NameChangedEvent{
		Name: newName,
	})
	a.Name = newName
}
