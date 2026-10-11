package v1_16_100

import (
	"github.com/helix-bedrock/protocol"
	"github.com/helix-bedrock/protocol/packet"
)

func marshalLevelChunk(io *wireIO, raw packet.Packet) {
	pk := raw.(*packet.LevelChunk)
	io.ChunkPos(&pk.Position)
	count := pk.SubChunkCount
	io.Varuint32(&count)
	pk.SubChunkCount = count
	if io.reading {
		pk.SubChunkLimit = protocol.Optional[int32]{}
	}
	io.Bool(&pk.CacheEnabled)
	if pk.CacheEnabled {
		protocol.FuncSlice(io.directional(), &pk.BlobHashes, io.Uint64)
	} else if io.reading {
		pk.BlobHashes = nil
	}
	io.ByteSlice(&pk.RawPayload)
	if io.reading {
		pk.Dimension = packet.DimensionOverworld
	}
}
