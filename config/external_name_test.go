package config

import (
	"context"
	"testing"
)

// A ClickPipe has no id before create. An empty id makes the upstream Read GET
// the list endpoint and fail to unmarshal an array, so Create never runs.
func TestClickPipeGetIDFnUsesSentinelBeforeCreate(t *testing.T) {
	getID := ExternalNameConfigs["clickhouse_clickpipe"].GetIDFn
	const real = "3f1c2a9e-5b7d-4e8f-9a0b-1c2d3e4f5a6b"
	for in, want := range map[string]string{
		"":        clickPipeSentinelUUID,
		"my-pipe": clickPipeSentinelUUID,
		real:      real,
	} {
		got, err := getID(context.Background(), in, nil, nil)
		if err != nil {
			t.Fatalf("GetIDFn(%q): %v", in, err)
		}
		if got != want {
			t.Errorf("GetIDFn(%q) = %q, want %q", in, got, want)
		}
	}
}
