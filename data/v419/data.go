// Package v419 exposes the locked Minecraft 1.16.100/protocol-419 registry snapshots.
package v419

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/helix-bedrock/protocol/nbt"
)

var (
	//go:embed block_states.nbt
	blockStateData []byte
	//go:embed item_runtime_ids.json
	itemRuntimeIDData []byte
	//go:embed block_item_meta.json
	blockItemMetadataData []byte
	//go:embed biome_definitions.nbt
	biomeDefinitionData []byte
)

// BlockState is one block state in the protocol 419 snapshot.
type BlockState struct {
	Name       string         `nbt:"name"`
	Properties map[string]any `nbt:"states"`
	Version    int32          `nbt:"version"`
}

// ItemEntry is one item in the protocol 419 snapshot.
type ItemEntry struct {
	RuntimeID      int16 `json:"runtime_id"`
	ComponentBased bool  `json:"component_based"`
}

// BlockStates returns the block states of this version, in runtime ID order.
func BlockStates() ([]BlockState, error) {
	reader := bytes.NewReader(blockStateData)
	decoder := nbt.NewDecoder(reader)
	states := make([]BlockState, 0, 1<<13)
	for {
		var state BlockState
		if err := decoder.Decode(&state); err != nil {
			if reader.Len() == 0 {
				return states, nil
			}
			return nil, fmt.Errorf("decode block state %d: %w", len(states), err)
		}
		states = append(states, state)
	}
}

// Items returns the item entries of this version.
func Items() (map[string]ItemEntry, error) {
	items := map[string]ItemEntry{}
	if err := json.Unmarshal(itemRuntimeIDData, &items); err != nil {
		return nil, fmt.Errorf("decode item runtime IDs: %w", err)
	}
	return items, nil
}

// BlockItemMetadata returns the item metadata of each block state, indexed by block runtime ID.
func BlockItemMetadata() ([]int16, error) {
	var metadata []int16
	if err := json.Unmarshal(blockItemMetadataData, &metadata); err != nil {
		return nil, fmt.Errorf("decode block item metadata: %w", err)
	}
	return metadata, nil
}

// BiomeDefinitions returns the embedded biome definitions of this version as NBT.
func BiomeDefinitions() []byte { return bytes.Clone(biomeDefinitionData) }

// RawSnapshot returns the embedded snapshot file with the name passed, byte for byte, and false if there is none.
func RawSnapshot(name string) ([]byte, bool) {
	var data []byte
	switch name {
	case "block_states.nbt":
		data = blockStateData
	case "item_runtime_ids.json":
		data = itemRuntimeIDData
	case "block_item_meta.json":
		data = blockItemMetadataData
	case "biome_definitions.nbt":
		data = biomeDefinitionData
	default:
		return nil, false
	}
	return bytes.Clone(data), true
}
