// Package v2169 exposes the outgoing Minecraft 1.26.45/protocol 2169
// Dragonfly registry snapshots using production-compatible decoders.
package v2169

import (
	"bytes"
	_ "embed"
	"fmt"

	v1001 "github.com/helix-bedrock/multiversion/data/v1001"
	"github.com/helix-bedrock/protocol/nbt"
)

var (
	//go:embed block_states.nbt
	blockStateData []byte
	//go:embed vanilla_items.nbt
	vanillaItemData []byte
	//go:embed creative_items.nbt
	creativeItemData []byte
	//go:embed crafting_data.nbt
	craftingData []byte
	//go:embed potion_data.nbt
	potionData []byte
	//go:embed item_tags.json
	itemTagData []byte
)

type (
	// BlockState is the same as in data/v1001: its layout did not change.
	BlockState = v1001.BlockState
	// ItemEntry is the same as in data/v1001: its layout did not change.
	ItemEntry = v1001.ItemEntry
	// CreativeData is the same as in data/v1001: its layout did not change.
	CreativeData = v1001.CreativeData
	// CraftingData is the same as in data/v1001: its layout did not change.
	CraftingData = v1001.CraftingData
	// PotionData is the same as in data/v1001: its layout did not change.
	PotionData = v1001.PotionData
	// CreativeGroup is the same as in data/v1001: its layout did not change.
	CreativeGroup = v1001.CreativeGroup
	// CreativeItem is the same as in data/v1001: its layout did not change.
	CreativeItem = v1001.CreativeItem
	// RecipeBlockState is the same as in data/v1001: its layout did not change.
	RecipeBlockState = v1001.RecipeBlockState
	// RecipeInput is the same as in data/v1001: its layout did not change.
	RecipeInput = v1001.RecipeInput
	// RecipeOutput is the same as in data/v1001: its layout did not change.
	RecipeOutput = v1001.RecipeOutput
	// ShapedRecipe is the same as in data/v1001: its layout did not change.
	ShapedRecipe = v1001.ShapedRecipe
	// ShapelessRecipe is the same as in data/v1001: its layout did not change.
	ShapelessRecipe = v1001.ShapelessRecipe
	// PotionRecipe is the same as in data/v1001: its layout did not change.
	PotionRecipe = v1001.PotionRecipe
	// PotionContainerChangeRecipe is the same as in data/v1001: its layout did not change.
	PotionContainerChangeRecipe = v1001.PotionContainerChangeRecipe
)

// BlockStates returns the block states of this version, in runtime ID order.
func BlockStates() ([]BlockState, error) {
	reader := bytes.NewReader(blockStateData)
	decoder := nbt.NewDecoder(reader)
	states := make([]BlockState, 0, 1<<16)
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
	if err := nbt.Unmarshal(vanillaItemData, &items); err != nil {
		return nil, fmt.Errorf("decode vanilla items: %w", err)
	}
	return items, nil
}

// Creative returns the creative inventory of this version.
func Creative() (CreativeData, error) {
	var data CreativeData
	if err := nbt.Unmarshal(creativeItemData, &data); err != nil {
		return CreativeData{}, fmt.Errorf("decode creative data: %w", err)
	}
	return data, nil
}

// Crafting returns the crafting recipes of this version.
func Crafting() (CraftingData, error) {
	var data CraftingData
	if err := nbt.Unmarshal(craftingData, &data); err != nil {
		return CraftingData{}, fmt.Errorf("decode crafting data: %w", err)
	}
	return data, nil
}

// Potions returns the potion recipes of this version.
func Potions() (PotionData, error) {
	var data PotionData
	if err := nbt.Unmarshal(potionData, &data); err != nil {
		return PotionData{}, fmt.Errorf("decode potion data: %w", err)
	}
	return data, nil
}

// RawSnapshot returns the embedded snapshot file with the name passed, byte for byte, and false if there is none.
func RawSnapshot(name string) ([]byte, bool) {
	var data []byte
	switch name {
	case "block_states.nbt":
		data = blockStateData
	case "vanilla_items.nbt":
		data = vanillaItemData
	case "creative_items.nbt":
		data = creativeItemData
	case "crafting_data.nbt":
		data = craftingData
	case "potion_data.nbt":
		data = potionData
	case "item_tags.json":
		data = itemTagData
	default:
		return nil, false
	}
	return bytes.Clone(data), true
}
