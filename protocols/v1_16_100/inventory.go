package v1_16_100

import (
	"fmt"

	"github.com/helix-bedrock/protocol"
	"github.com/helix-bedrock/protocol/packet"
)

func marshalTransactionDataType(io *wireIO, value *protocol.InventoryTransactionData) {
	var typeID uint32
	if !io.reading {
		switch (*value).(type) {
		case nil, *protocol.NormalTransactionData:
			typeID = protocol.InventoryTransactionTypeNormal
		case *protocol.MismatchTransactionData:
			typeID = protocol.InventoryTransactionTypeMismatch
		case *protocol.UseItemTransactionData:
			typeID = protocol.InventoryTransactionTypeUseItem
		case *protocol.UseItemOnEntityTransactionData:
			typeID = protocol.InventoryTransactionTypeUseItemOnEntity
		case *protocol.ReleaseItemTransactionData:
			typeID = protocol.InventoryTransactionTypeReleaseItem
		default:
			io.UnknownEnumOption(fmt.Sprintf("%T", *value), "inventory transaction data type")
			return
		}
	}
	io.Varuint32(&typeID)
	if io.reading {
		switch typeID {
		case protocol.InventoryTransactionTypeNormal:
			*value = &protocol.NormalTransactionData{}
		case protocol.InventoryTransactionTypeMismatch:
			*value = &protocol.MismatchTransactionData{}
		case protocol.InventoryTransactionTypeUseItem:
			*value = &protocol.UseItemTransactionData{}
		case protocol.InventoryTransactionTypeUseItemOnEntity:
			*value = &protocol.UseItemOnEntityTransactionData{}
		case protocol.InventoryTransactionTypeReleaseItem:
			*value = &protocol.ReleaseItemTransactionData{}
		default:
			io.UnknownEnumOption(typeID, "inventory transaction data type")
			return
		}
	}
}

func marshalInventoryTransactionPacket(io *wireIO, raw packet.Packet) {
	pk := raw.(*packet.InventoryTransaction)
	io.Varint32(&pk.LegacyRequestID)
	if pk.LegacyRequestID != 0 {
		protocol.Slice(io.directional(), &pk.LegacySetItemSlots)
	}
	marshalTransactionDataType(io, &pk.TransactionData)
	var hasNetworkIDs bool
	io.Bool(&hasNetworkIDs)
	protocol.FuncIOSlice(io.directional(), &pk.Actions, func(io protocol.IO, action *protocol.InventoryAction) {
		marshalProtocol419InventoryAction(asWireIO(io), action, hasNetworkIDs)
	})
	marshalLegacyInventoryTransactionData(io, pk.TransactionData)
}

func marshalProtocol419InventoryAction(io *wireIO, action *protocol.InventoryAction, hasNetworkIDs bool) {
	io.Varuint32(&action.SourceType)
	switch action.SourceType {
	case protocol.InventoryActionSourceContainer, protocol.InventoryActionSourceTODO:
		windowID, _ := action.WindowID.Value()
		legacyWindowID := int32(windowID)
		io.Varint32(&legacyWindowID)
		if io.reading {
			action.WindowID = protocol.Option(int8(legacyWindowID))
		}
	case protocol.InventoryActionSourceWorld:
		flags, _ := action.SourceFlags.Value()
		io.Varuint32(&flags)
		if io.reading {
			action.SourceFlags = protocol.Option(flags)
		}
	}
	io.Varuint32(&action.InventorySlot)
	marshalItem(io, &action.OldItem.Stack)
	marshalItem(io, &action.NewItem.Stack)
	if hasNetworkIDs {
		io.Varint32(&action.NewItem.StackNetworkID)
	} else if io.reading {
		action.OldItem.StackNetworkID = 0
		action.NewItem.StackNetworkID = 0
	}
}

func marshalLegacyInventoryTransactionData(io *wireIO, value protocol.InventoryTransactionData) {
	switch data := value.(type) {
	case nil, *protocol.NormalTransactionData, *protocol.MismatchTransactionData:
	case *protocol.UseItemTransactionData:
		io.Varuint32(&data.ActionType)
		marshalUnsignedBlockPos(io, &data.BlockPosition)
		io.Varint32(&data.BlockFace)
		io.Varint32(&data.HotBarSlot)
		marshalItem(io, &data.HeldItem.Stack)
		io.Vec3(&data.Position)
		io.Vec3(&data.ClickedPosition)
		io.Varuint32(&data.BlockRuntimeID)
		if io.reading {
			data.TriggerType = 0
			data.ClientPrediction = 0
			data.ClientCooldownState = 0
		}
	case *protocol.UseItemOnEntityTransactionData:
		io.Varuint64(&data.TargetEntityRuntimeID)
		actionType := uint32(data.ActionType)
		io.Varuint32(&actionType)
		data.ActionType = int32(actionType)
		io.Varint32(&data.HotBarSlot)
		marshalItem(io, &data.HeldItem.Stack)
		io.Vec3(&data.Position)
		io.Vec3(&data.ClickedPosition)
	case *protocol.ReleaseItemTransactionData:
		actionType := uint32(data.ActionType)
		io.Varuint32(&actionType)
		data.ActionType = int32(actionType)
		io.Varint32(&data.HotBarSlot)
		marshalItem(io, &data.HeldItem.Stack)
		io.Vec3(&data.HeadPosition)
	default:
		io.UnknownEnumOption(fmt.Sprintf("%T", value), "inventory transaction data type")
	}
}

func marshalPlayerInventoryAction(io *wireIO, data *protocol.UseItemTransactionData) {
	io.Varint32(&data.LegacyRequestID)
	if data.LegacyRequestID < -1 && data.LegacyRequestID&1 == 0 {
		slots, _ := data.LegacySetItemSlots.Value()
		protocol.FuncIOSlice(io.directional(), &slots, func(io protocol.IO, slot *protocol.LegacySetItemSlot) {
			io.Uint8(&slot.ContainerID)
			io.ByteSlice(&slot.Slots)
		})
		data.LegacySetItemSlots = protocol.Option(slots)
	} else if io.reading {
		data.LegacySetItemSlots = protocol.Optional[[]protocol.LegacySetItemSlot]{}
	}
	actions := data.Actions
	protocol.FuncIOSlice(io.directional(), &actions, marshalLegacyPlayerInventoryAction)
	data.Actions = actions
	io.Varuint32(&data.ActionType)
	io.BlockPos(&data.BlockPosition)
	io.Varint32(&data.BlockFace)
	io.Varint32(&data.HotBarSlot)
	marshalItemInstance(io, &data.HeldItem)
	io.Vec3(&data.Position)
	io.Vec3(&data.ClickedPosition)
	io.Varuint32(&data.BlockRuntimeID)
	if io.reading {
		data.TriggerType = 0
		data.ClientPrediction = 0
		data.ClientCooldownState = 0
	}
}

func marshalLegacyPlayerInventoryAction(io protocol.IO, action *protocol.InventoryAction) {
	legacy := asWireIO(io)
	legacy.Varuint32(&action.SourceType)
	switch action.SourceType {
	case protocol.InventoryActionSourceContainer, protocol.InventoryActionSourceTODO:
		windowID, _ := action.WindowID.Value()
		protocol.IntegerFunc(&windowID, legacy.Varint32)
		action.WindowID = protocol.Option(windowID)
	case protocol.InventoryActionSourceWorld:
		flags, _ := action.SourceFlags.Value()
		legacy.Varuint32(&flags)
		action.SourceFlags = protocol.Option(flags)
	}
	legacy.Varuint32(&action.InventorySlot)
	marshalItemInstance(legacy, &action.OldItem)
	marshalItemInstance(legacy, &action.NewItem)
}
