package v1_16_100

import (
	"fmt"
	"strings"
	"sync"

	v419data "github.com/helix-bedrock/multiversion/data/v419"
	"github.com/helix-bedrock/multiversion/mapping"
)

type targetBlockItemKey struct {
	name string
	meta int16
}

// targetBlockItems holds the protocol 419 block item tables, loaded on first
// use. If the embedded data cannot be loaded, the tables are empty and block
// items are not mapped; the error is kept in err.
var targetBlockItems = sync.OnceValue(loadTargetBlockItems)

type blockItemTables struct {
	metadata map[string]int16
	states   map[targetBlockItemKey]mapping.BlockState
	err      error
}

func loadTargetBlockItems() blockItemTables {
	states, err := v419data.BlockStates()
	if err != nil {
		return blockItemTables{err: fmt.Errorf("load protocol 419 block item states: %w", err)}
	}
	metadata, err := v419data.BlockItemMetadata()
	if err != nil {
		return blockItemTables{err: fmt.Errorf("load protocol 419 block item metadata: %w", err)}
	}
	if len(metadata) != len(states) {
		return blockItemTables{err: fmt.Errorf("protocol 419 block item metadata count %d differs from state count %d", len(metadata), len(states))}
	}
	t := blockItemTables{
		metadata: make(map[string]int16, len(states)),
		states:   make(map[targetBlockItemKey]mapping.BlockState, len(states)),
	}
	for runtimeID, state := range states {
		stateKey, err := mapping.StateKey(state.Name, state.Properties)
		if err != nil {
			return blockItemTables{err: fmt.Errorf("build protocol 419 block item state key %d: %w", runtimeID, err)}
		}
		t.metadata[stateKey] = metadata[runtimeID]
		key := targetBlockItemKey{name: normaliseTargetBlockItemName(state.Name), meta: metadata[runtimeID]}
		if _, exists := t.states[key]; !exists {
			t.states[key] = mapping.BlockState{Name: state.Name, Properties: state.Properties, Version: state.Version}
		}
	}
	return t
}

func targetBlockItemMeta(state mapping.BlockState) (int16, bool) {
	key, err := mapping.StateKey(state.Name, state.Properties)
	if err != nil {
		return 0, false
	}
	metadata, ok := targetBlockItems().metadata[key]
	return metadata, ok
}

func targetBlockItemState(name string, metadata uint32) (mapping.BlockState, bool) {
	state, ok := targetBlockItems().states[targetBlockItemKey{name: normaliseTargetBlockItemName(name), meta: int16(metadata)}]
	return state, ok
}

func normaliseTargetBlockItemName(name string) string {
	if strings.HasPrefix(name, "minecraft:item.") {
		return "minecraft:" + strings.TrimPrefix(name, "minecraft:item.")
	}
	return name
}
