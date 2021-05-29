package bank

import (
	"encoding/json"
	"es/eventstore"
	"github.com/shopspring/decimal"
)

const (
	AccountCreatedEventType = "AccountCreatedEvent"
	MoneyDepositedEventType = "MoneyDepositedEvent"
	MoneyWithdrawnEventType = "MoneyWithdrawnEvent"
	NameChangedEventType    = "NameChangedEvent"
)

type AccountCreatedEvent struct {
	eventstore.Metadata `json:"-"`
	ID                  string `json:"id"`
	Name                string `json:"name"`
}

func (c AccountCreatedEvent) Type() string {
	return AccountCreatedEventType
}

func (c AccountCreatedEvent) Data() []byte {
	bs, _ := json.Marshal(&c)
	return bs
}

func (c AccountCreatedEvent) Version() uint64 {
	return c.Metadata.Version
}

type MoneyDepositedEvent struct {
	eventstore.Metadata `json:"-"`
	Amount              decimal.Decimal `json:"amount"`
}

func (c MoneyDepositedEvent) Type() string {
	return MoneyDepositedEventType
}

func (c MoneyDepositedEvent) Data() []byte {
	bs, _ := json.Marshal(&c)
	return bs
}

func (c MoneyDepositedEvent) Version() uint64 {
	return c.Metadata.Version
}

type MoneyWithdrawnEvent struct {
	eventstore.Metadata `json:"-"`
	Amount              decimal.Decimal `json:"amount"`
}

func (c MoneyWithdrawnEvent) Type() string {
	return MoneyWithdrawnEventType
}

func (c MoneyWithdrawnEvent) Data() []byte {
	bs, _ := json.Marshal(&c)
	return bs
}

func (c MoneyWithdrawnEvent) Version() uint64 {
	return c.Metadata.Version
}

type NameChangedEvent struct {
	eventstore.Metadata `json:"-"`
	Name                string `json:"name"`
}

func (c NameChangedEvent) Type() string {
	return NameChangedEventType
}

func (c NameChangedEvent) Data() []byte {
	bs, _ := json.Marshal(&c)
	return bs
}

func (c NameChangedEvent) Version() uint64 {
	return c.Metadata.Version
}
