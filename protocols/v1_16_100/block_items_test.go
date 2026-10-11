package v1_16_100

import "testing"

func TestTargetBlockItemsLoad(t *testing.T) {
	if err := targetBlockItems().err; err != nil {
		t.Fatalf("load block item tables: %v", err)
	}
	if len(targetBlockItems().metadata) == 0 {
		t.Fatal("block item tables are empty")
	}
}
