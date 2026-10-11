package v1_21_50

import (
	"github.com/helix-bedrock/protocol/packet"
	v766 "github.com/helix-bedrock/multiversion/data/v766"
)

func marshalBiomeDefinitionList(io *wireIO, _ packet.Packet) {
	definitions := v766.BiomeDefinitions()
	io.Bytes(&definitions)
}
