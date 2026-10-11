package v1_21_40

import (
	v748 "github.com/helix-bedrock/multiversion/data/v748"
	"github.com/helix-bedrock/protocol/packet"
)

func marshalBiomeDefinitionList(io *wireIO, _ packet.Packet) {
	definitions := v748.BiomeDefinitions()
	io.Bytes(&definitions)
}
