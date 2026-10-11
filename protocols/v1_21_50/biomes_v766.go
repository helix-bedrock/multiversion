package v1_21_50

import (
	v766 "github.com/helix-bedrock/multiversion/data/v766"
	"github.com/helix-bedrock/protocol/packet"
)

func marshalBiomeDefinitionList(io *wireIO, _ packet.Packet) {
	definitions := v766.BiomeDefinitions()
	io.Bytes(&definitions)
}
