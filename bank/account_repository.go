package bank

import (
	"context"
	"encoding/json"
	"errors"
	"es/eventstore"
)

func PublishAccountCreated(
	ctx context.Context,
	accountCreatedEvent AccountCreatedEvent,
) error {
	return eventstore.NewStream(ctx, StreamName+accountCreatedEvent.ID, []eventstore.Event{
		accountCreatedEvent,
	})
}

func PublishAccountEvents(
	ctx context.Context,
	accountAggregate AccountAggregate,
) error {
	return eventstore.Append(ctx, StreamName+accountAggregate.ID, accountAggregate.events, accountAggregate.Version)
}

func FindAccount(
	ctx context.Context,
	accountID string,
) (AccountAggregate, error) {
	streamName := StreamName + accountID
	recordedEvents, err := eventstore.ReplayStream(ctx, streamName)

	var events []eventstore.Event
	for _, recordedEvent := range recordedEvents {
		var event eventstore.Event
		switch recordedEvent.EventType {
		case AccountCreatedEventType:
			{
				event = &AccountCreatedEvent{}
				err = json.Unmarshal(recordedEvent.Data, &event)
				if err != nil {
					return AccountAggregate{}, err
				}
				event.(*AccountCreatedEvent).Metadata.Version = recordedEvent.EventNumber
			}
		case MoneyDepositedEventType:
			{
				event = &MoneyDepositedEvent{}
				err = json.Unmarshal(recordedEvent.Data, &event)
				if err != nil {
					return AccountAggregate{}, err
				}
				event.(*MoneyDepositedEvent).Metadata.Version = recordedEvent.EventNumber
			}
		case MoneyWithdrawnEventType:
			{
				event = &MoneyWithdrawnEvent{}
				err = json.Unmarshal(recordedEvent.Data, &event)
				if err != nil {
					return AccountAggregate{}, err
				}
				event.(*MoneyWithdrawnEvent).Metadata.Version = recordedEvent.EventNumber
			}
		case NameChangedEventType:
			{
				event = &NameChangedEvent{}
				err = json.Unmarshal(recordedEvent.Data, &event)
				if err != nil {
					return AccountAggregate{}, err
				}
				event.(*NameChangedEvent).Metadata.Version = recordedEvent.EventNumber
			}
		default:
			return AccountAggregate{}, errors.New("invalid event type")
		}
		events = append(events, event)
	}

	accountAggregate := AccountAggregate{}
	err = accountAggregate.BuildFromEvents(events)

	return accountAggregate, err
}
