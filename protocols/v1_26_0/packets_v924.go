package v1_26_0

import (
	"github.com/helix-bedrock/protocol"
	"github.com/helix-bedrock/protocol/packet"
)

func marshalBossEvent(io *wireIO, raw packet.Packet) {
	pk := raw.(*packet.BossEvent)
	io.Varint64(&pk.BossEntityUniqueID)
	eventType := uint32(pk.EventType)
	io.Varuint32(&eventType)
	pk.EventType = uint8(eventType)
	var screenDarkening uint16
	switch eventType {
	case packet.BossEventShow:
		io.String(&pk.BossBarTitle)
		io.String(&pk.FilteredBossBarTitle)
		io.Float32(&pk.HealthPercentage)
		io.Uint16(&screenDarkening)
		marshalLegacyBossAppearance(io, pk)
	case packet.BossEventRegisterPlayer, packet.BossEventUnregisterPlayer, packet.BossEventRequest:
		var discardedPlayerUniqueID int64
		io.Varint64(&discardedPlayerUniqueID)
	case packet.BossEventHide:
	case packet.BossEventHealthPercentage:
		io.Float32(&pk.HealthPercentage)
	case packet.BossEventTitle:
		io.String(&pk.BossBarTitle)
		io.String(&pk.FilteredBossBarTitle)
	case packet.BossEventAppearanceProperties:
		io.Uint16(&screenDarkening)
		marshalLegacyBossAppearance(io, pk)
	case packet.BossEventTexture:
		marshalLegacyBossAppearance(io, pk)
	default:
		io.UnknownEnumOption(eventType, "boss event type")
	}
}

func marshalLegacyBossAppearance(io *wireIO, pk *packet.BossEvent) {
	colour, overlay := uint32(pk.Colour), uint32(pk.Overlay)
	if !io.reading {
		switch colour {
		case packet.BossEventColourWhite:
			colour = 6
		case packet.BossEventColourRebeccaPurple:
			colour = packet.BossEventColourPurple
		}
	}
	io.Varuint32(&colour)
	io.Varuint32(&overlay)
	if io.reading {
		if colour == 6 {
			colour = packet.BossEventColourWhite
		}
		pk.Colour, pk.Overlay = uint8(colour), uint8(overlay)
	}
}

func marshalClientCacheBlobStatus(io *wireIO, raw packet.Packet) {
	pk := raw.(*packet.ClientCacheBlobStatus)
	missCount, hitCount := uint32(len(pk.MissHashes)), uint32(len(pk.HitHashes))
	io.Varuint32(&missCount)
	io.Varuint32(&hitCount)
	protocol.FuncSliceOfLen(io.directional(), missCount, &pk.MissHashes, io.Uint64)
	protocol.FuncSliceOfLen(io.directional(), hitCount, &pk.HitHashes, io.Uint64)
}

func marshalGraphicsOverrideParameter(io *wireIO, raw packet.Packet) {
	pk := raw.(*packet.GraphicsOverrideParameter)
	protocol.Slice(io.directional(), &pk.Values)
	floatValue, _ := pk.FloatValue.Value()
	io.Float32(&floatValue)
	vec3Value, _ := pk.Vec3Value.Value()
	io.Vec3(&vec3Value)
	io.String(&pk.BiomeIdentifier)
	io.Uint8(&pk.ParameterType)
	io.Bool(&pk.Reset)
	if io.reading {
		pk.FloatValue = protocol.Option(floatValue)
		pk.Vec3Value = protocol.Option(vec3Value)
		pk.PlayerIdentifier = protocol.Optional[string]{}
	}
}
