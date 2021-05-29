package bank

import (
	"context"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const StreamName = "bank-account-"

func CreateAccount(ctx context.Context, name string) (string, error) {
	accountCreatedEvent := AccountCreatedEvent{
		ID:   uuid.New().String(),
		Name: name,
	}

	err := PublishAccountCreated(ctx, accountCreatedEvent)
	if err != nil {
		return "", err
	}

	return accountCreatedEvent.ID, nil
}

func Deposit(
	ctx context.Context,
	accountID string,
	amount decimal.Decimal,
) error {
	accountAggregate, err := FindAccount(ctx, accountID)
	if err != nil {
		return err
	}

	accountAggregate.Deposit(amount)
	return PublishAccountEvents(ctx, accountAggregate)
}

func Withdraw(
	ctx context.Context,
	accountID string,
	amount decimal.Decimal,
) error {
	accountAggregate, err := FindAccount(ctx, accountID)
	if err != nil {
		return err
	}

	err = accountAggregate.Withdraw(amount)
	if err != nil {
		return err
	}

	return PublishAccountEvents(ctx, accountAggregate)
}

func ProcessEscrow(
	ctx context.Context,
	accountID string,
	amounts []decimal.Decimal,
) error {
	accountAggregate, err := FindAccount(ctx, accountID)
	if err != nil {
		return err
	}

	for _, amount := range amounts {
		accountAggregate.Deposit(amount)
	}

	return PublishAccountEvents(ctx, accountAggregate)
}

func ChangeName(ctx context.Context, accountID string, newName string) error {
	accountAggregate, err := FindAccount(ctx, accountID)
	if err != nil {
		return err
	}

	accountAggregate.ChangeName(newName)
	return PublishAccountEvents(ctx, accountAggregate)
}
