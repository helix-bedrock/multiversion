package v1_21_110

import (
	"github.com/helix-bedrock/protocol"
	"github.com/helix-bedrock/protocol/packet"
)

func marshalAddPlayer(io *wireIO, raw packet.Packet) {
	pk := raw.(*packet.AddPlayer)
	io.UUID(&pk.UUID)
	io.String(&pk.Username)
	io.Varuint64(&pk.EntityRuntimeID)
	io.String(&pk.PlatformChatID)
	io.Vec3(&pk.Position)
	io.Vec3(&pk.Velocity)
	io.Float32(&pk.Pitch)
	io.Float32(&pk.Yaw)
	io.Float32(&pk.HeadYaw)
	io.ItemInstance(&pk.HeldItem)
	io.Varint32(&pk.GameType)
	io.EntityMetadata(&pk.EntityMetadata)
	protocol.Single(io.directional(), &pk.EntityProperties)
	marshalAbilityData(io, &pk.AbilityData)
	protocol.Slice(io.directional(), &pk.EntityLinks)
	io.String(&pk.DeviceID)
	io.Int32(&pk.BuildPlatform)
}

func marshalClientCheatAbility(io *wireIO, raw packet.Packet) {
	marshalAbilityData(io, &raw.(*packet.ClientCheatAbility).AbilityData)
}

func marshalUpdateAbilities(io *wireIO, raw packet.Packet) {
	marshalAbilityData(io, &raw.(*packet.UpdateAbilities).AbilityData)
}

func marshalCreativeContent(io *wireIO, raw packet.Packet) {
	pk := raw.(*packet.CreativeContent)
	protocol.FuncIOSlice(io.directional(), &pk.Groups, func(io protocol.IO, group *protocol.CreativeGroup) {
		legacy := asWireIO(io)
		category := int32(group.Category)
		legacy.Int32(&category)
		group.Category = byte(category)
		legacy.String(&group.Name)
		legacy.Item(&group.Icon)
	})
	protocol.Slice(io.directional(), &pk.Items)
}

func marshalAvailableCommands(io *wireIO, raw packet.Packet) {
	marshalAvailableCommands844(io, raw.(*packet.AvailableCommands))
}
