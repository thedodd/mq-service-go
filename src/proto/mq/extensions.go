package mq

// SystemEventMessage is the interface definition use to mark the specific message
// types which can be emitted to the broker's `events` exchnage.
//
// `Reset, String & ProtoMessage` are automatically generated from protoc code generation.
//
// WARNING!!! NOTE: do not arbitrarily extend this interface or add random type conformance.
// This interface is designed to work exactly with the protobuf message types which are actually
// valid `SystemEvent` types.
type SystemEventMessage interface {
	isSystemEvent_Event

	// RoutingKey is this message's routing key.
	RoutingKey() string
}

// WARNING!!! NOTE: do not randomly add type conformance to the `SystemEventMessage` interface.
// The interface is designed to work exactly with the protobuf message types which are actually
// valid `SystemEvent` types.

// RoutingKey is this message's routing key.
func (msg *SystemEvent_PhotoScanUploaded) RoutingKey() string {
	return "events.photoscan.uploaded"
}

// RoutingKey is this message's routing key.
func (msg *SystemEvent_PhotoScanSampled) RoutingKey() string {
	return "events.photoscan.sampled"
}
