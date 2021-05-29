package eventstore

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/client"
	"github.com/EventStore/EventStore-Client-Go/direction"
	"github.com/EventStore/EventStore-Client-Go/messages"
	"github.com/EventStore/EventStore-Client-Go/streamrevision"
	"github.com/gofrs/uuid"
	"log"
	"math"
)

var esc *client.Client

type Event interface {
	Type() string
	Data() []byte
	Version() uint64
}

type Metadata struct {
	Version uint64 `json:"version"`
}

func Init() {
	config, err := client.ParseConnectionString("esdb://127.0.0.1:2113?tls=false")
	if err != nil {
		panic(err)
	}
	esClient, err := client.NewClient(config)
	if err != nil {
		panic(err)
	}

	err = esClient.Connect()
	if err != nil {
		panic(err)
	}

	esc = esClient
}

func Close() {
	if esc != nil {
		err := esc.Close()
		if err != nil {
			panic(err)
		}
	}
}

func NewStream(ctx context.Context, streamName string, events []Event) error {
	// create messages
	var proposedEvents []messages.ProposedEvent
	for _, event := range events {
		proposedEvents = append(proposedEvents, messages.ProposedEvent{
			EventID:     uuid.Must(uuid.NewV4()),
			EventType:   event.Type(),
			ContentType: "application/json",
			Data:        event.Data(),
		})
	}

	wr, err := esc.AppendToStream(ctx, streamName, streamrevision.StreamRevisionNoStream, proposedEvents)
	if err != nil {
		return err
	}

	log.Printf(
		"Appended events. New commit position [%d], next position [%d]",
		wr.CommitPosition,
		wr.NextExpectedVersion,
	)

	return nil
}

func Append(ctx context.Context, streamName string, events []Event, checkpoint uint64) error {
	// create messages
	var proposedEvents []messages.ProposedEvent
	for _, event := range events {
		proposedEvents = append(proposedEvents, messages.ProposedEvent{
			EventID:     uuid.Must(uuid.NewV4()),
			EventType:   event.Type(),
			ContentType: "application/json",
			Data:        event.Data(),
		})
	}

	wr, err := esc.AppendToStream(ctx, streamName, streamrevision.NewStreamRevision(checkpoint), proposedEvents)
	if err != nil {
		return err
	}

	log.Printf(
		"Appended events. New commit position [%d], next position [%d]",
		wr.CommitPosition,
		wr.NextExpectedVersion,
	)

	return nil
}

func ReplayStream(ctx context.Context, streamName string) ([]messages.RecordedEvent, error) {
	recordedEvents, err := esc.ReadStreamEvents(
		ctx,
		direction.Forwards,
		streamName,
		0,
		math.MaxUint64,
		false,
	)

	return recordedEvents, err
}

func AllStreams(ctx context.Context) ([]string, error) {
	recordedEvents, err := esc.ReadStreamEvents(
		ctx,
		direction.Forwards,
		"$et-CreateAccountCommand",
		0,
		math.MaxUint64,
		false,
	)
	if err != nil {
		return nil, err
	}

	var streamNames []string
	for _, recordedEvent := range recordedEvents {
		streamNames = append(streamNames, string(recordedEvent.Data))
	}

	return streamNames, nil
}
