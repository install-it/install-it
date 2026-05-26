// Package storage_test provides external black-box tests for DeleteEventBus.
// The internal tests in store_test.go cover the same struct under "TestDeleteEventBus_*";
// this file provides equivalent coverage from the external test package perspective.
package storage_test

import (
	"errors"
	"testing"

	"install-it/pkg/storage"
)

func TestEventBus_PublishSubscribeBasic(t *testing.T) {
	t.Parallel()

	bus := storage.NewEventBus()

	var received []string
	bus.Subscribe("drivers", func(ids []string) error {
		received = ids
		return nil
	})

	if err := bus.Publish("drivers", []string{"aabbccdd", "11223344"}); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	if len(received) != 2 {
		t.Fatalf("handler received %d IDs, want 2", len(received))
	}
	if received[0] != "aabbccdd" || received[1] != "11223344" {
		t.Errorf("received %v, want [aabbccdd 11223344]", received)
	}
}

func TestEventBus_MultipleSubscribers(t *testing.T) {
	t.Parallel()

	bus := storage.NewEventBus()

	calls := 0
	for i := 0; i < 3; i++ {
		bus.Subscribe("key", func(ids []string) error {
			calls++
			return nil
		})
	}

	if err := bus.Publish("key", []string{"x"}); err != nil {
		t.Fatalf("Publish: %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 subscriber calls, got %d", calls)
	}
}

func TestEventBus_SubscriberError(t *testing.T) {
	t.Parallel()

	bus := storage.NewEventBus()
	sentinel := errors.New("handler error")

	bus.Subscribe("key", func(ids []string) error {
		return sentinel
	})

	err := bus.Publish("key", []string{"a"})
	if err == nil {
		t.Fatal("Publish should propagate subscriber error, got nil")
	}
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestEventBus_NoSubscribersNoError(t *testing.T) {
	t.Parallel()

	bus := storage.NewEventBus()
	if err := bus.Publish("nobody_subscribed", []string{"id1"}); err != nil {
		t.Errorf("Publish with no subscribers: got %v, want nil", err)
	}
}

func TestEventBus_DifferentStorageKeys(t *testing.T) {
	t.Parallel()

	bus := storage.NewEventBus()

	called := false
	bus.Subscribe("storageA", func(ids []string) error {
		called = true
		return nil
	})

	// Publish to a different key — "storageA" handler must NOT be called.
	if err := bus.Publish("storageB", []string{"x"}); err != nil {
		t.Fatalf("Publish: %v", err)
	}
	if called {
		t.Error("handler for 'storageA' was called when publishing to 'storageB'")
	}
}

func TestEventBus_SubscriberReceivesCorrectIds(t *testing.T) {
	t.Parallel()

	bus := storage.NewEventBus()

	want := []string{"id-1", "id-2", "id-3"}
	var got []string
	bus.Subscribe("store", func(ids []string) error {
		got = append(got, ids...)
		return nil
	})

	if err := bus.Publish("store", want); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("received %d IDs, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("IDs[%d]: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestEventBus_PublishEmptyIdSlice(t *testing.T) {
	t.Parallel()

	bus := storage.NewEventBus()
	called := false
	bus.Subscribe("key", func(ids []string) error {
		called = true
		return nil
	})

	if err := bus.Publish("key", []string{}); err != nil {
		t.Fatalf("Publish empty slice: %v", err)
	}
	if !called {
		t.Error("subscriber was not called when publishing empty id slice")
	}
}
