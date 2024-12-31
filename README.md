# Signals

Forked from https://github.com/maniartech/signals.

This implementation adds rollback to the signal. If an error is detected during the handling, provided rollback methods will be called.

## Installation

```bash
go get github.com/imajacket/signals
```

## Usage

```go
package main

import (
  "context"
  "fmt"
  "github.com/imajacket/signals"
)

var RecordCreated = signals.New[Record]()
var RecordUpdated = signals.New[Record]()
var RecordDeleted = signals.New[Record]()

func main() {

  // Add a listener to the RecordCreated signal
  RecordCreated.AddListener(func(ctx context.Context, record Record) error {
    fmt.Println("Record created:", record)
	return nil
  }, nil, "key1") // <- Key is optional useful for removing the listener later

  // Add a listener to the RecordUpdated signal
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

  // Emit the RecordCreated signal
  RecordCreated.Emit(ctx, Record{ID: 1, Name: "John"})

  // Emit the RecordUpdated signal
  RecordUpdated.Emit(ctx, Record{ID: 1, Name: "John Doe"})

  // Emit the RecordDeleted signal
  RecordDeleted.Emit(ctx, Record{ID: 1, Name: "John Doe"})
}
```

## License

![License](https://img.shields.io/badge/license-MIT-blue.svg)
