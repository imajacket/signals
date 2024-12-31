package example

import (
	"context"
	"fmt"
	"log"
)

func RunAsync() {

	// Add a listener to the RecordCreatedAsync signal
	RecordCreated.AddListener(func(ctx context.Context, record Record) error {
		fmt.Println("Record created:", record)
		return nil
	}, nil, "key1") // <- Key is optional useful for removing the listener later

	// Add a listener to the RecordUpdatedAsync signal
	RecordUpdated.AddListener(func(ctx context.Context, record Record) error {
		fmt.Println("Record updated:", record)
		return nil
	}, nil)

	// Add a listener to the RecordDeleted signal
	RecordDeleted.AddListener(func(ctx context.Context, record Record) error {
		fmt.Println("Record deleted:", record)
		return nil
	}, nil)

	ctx := context.Background()

	// Emit the RecordCreatedAsync signal
	err := RecordCreated.Emit(ctx, Record{ID: 3, Name: "Record C"})
	if err != nil {
		log.Fatal(err)
	}

	// Emit the RecordUpdated signal
	err = RecordUpdated.Emit(ctx, Record{ID: 2, Name: "Record B"})
	if err != nil {
		log.Fatal(err)
	}

	// Emit the RecordDeleted signal
	err = RecordDeleted.Emit(ctx, Record{ID: 1, Name: "Record A"})
	if err != nil {
		log.Fatal(err)
	}
}
