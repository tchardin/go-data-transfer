package datatransfer

import (
	"context"
	"time"

	"github.com/filecoin-project/go-data-transfer/encoding"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	"github.com/libp2p/go-libp2p-core/peer"
)

// TypeIdentifier is a unique string identifier for a type of encodable object in a
// registry
type TypeIdentifier string

// Registerable is a type of object in a registry. It must be encodable and must
// have a single method that uniquely identifies its type
type Registerable interface {
	encoding.Encodable
	// Type is a unique string identifier for this voucher type
	Type() TypeIdentifier
}

// Voucher is used to validate
// a data transfer request against the underlying storage or retrieval deal
// that precipitated it. The only requirement is a voucher can read and write
// from bytes, and has a string identifier type
type Voucher Registerable

// VoucherResult is used to provide option additional information about a
// voucher being rejected or accepted
type VoucherResult Registerable

// RevalidationRequest is used to request revalidation of a data transfer
type RevalidationRequest Registerable

// Status is the status of transfer for a given channel
type Status int

const (
	// Ongoing means the data transfer is in progress
	Ongoing Status = iota

	// Completed means the data transfer is completed successfully
	Completed

	// Failed means the data transfer failed
	Failed

	// ChannelNotFoundError means the searched for data transfer does not exist
	ChannelNotFoundError
)

// TransferID is an identifier for a data transfer, shared between
// request/responder and unique to the requester
type TransferID uint64

// ChannelID is a unique identifier for a channel, distinct by both the other
// party's peer ID + the transfer ID
type ChannelID struct {
	Initiator peer.ID
	ID        TransferID
}

// Channel represents all the parameters for a single data transfer
type Channel interface {
	// TransferID returns the transfer id for this channel
	TransferID() TransferID

	// BaseCID returns the CID that is at the root of this data transfer
	BaseCID() cid.Cid

	// Selector returns the IPLD selector for this data transfer (represented as
	// an IPLD node)
	Selector() ipld.Node

	// Voucher returns the voucher for this data transfer
	Voucher() Voucher

	// Sender returns the peer id for the node that is sending data
	Sender() peer.ID

	// Recipient returns the peer id for the node that is receiving data
	Recipient() peer.ID

	// TotalSize returns the total size for the data being transferred
	TotalSize() uint64
}

// ChannelState is channel parameters plus it's current state
type ChannelState interface {
	Channel

	// Sent returns the number of bytes sent
	Sent() uint64

	// Received returns the number of bytes received
	Received() uint64
}

// EventCode is a name for an event that occurs on a data transfer channel
type EventCode int

const (
	// Open is an event occurs when a channel is first opened
	Open EventCode = iota

	// Accepted is an event that occurs when the data transfer is first accepted
	Accepted

	// Progress is an event that gets emitted every time more data is transferred
	Progress

	// Error is an event that emits when an error occurs in a data transfer
	Error

	// Complete is emitted when a data transfer is complete
	Complete
)

// Event is a struct containing information about a data transfer event
type Event struct {
	Code      EventCode // What type of event it is
	Message   string    // Any clarifying information about the event
	Timestamp time.Time // when the event happened
}

// Subscriber is a callback that is called when events are emitted
type Subscriber func(event Event, channelState ChannelState)

// Unsubscribe is a function that gets called to unsubscribe from data transfer events
type Unsubscribe func()

// RequestValidator is an interface implemented by the client of the
// data transfer module to validate requests
type RequestValidator interface {
	// ValidatePush validates a push request received from the peer that will send data
	ValidatePush(
		sender peer.ID,
		voucher Voucher,
		baseCid cid.Cid,
		selector ipld.Node) (VoucherResult, error)
	// ValidatePull validates a pull request received from the peer that will receive data
	ValidatePull(
		receiver peer.ID,
		voucher Voucher,
		baseCid cid.Cid,
		selector ipld.Node) (VoucherResult, error)
}

// Revalidator is a request validator revalidates in progress requests
// by requesting request additional vouchers, and resuming when it receives them
type Revalidator interface {
	RequestValidator
	// OnPullDataSent is called on the responder side when more bytes are sent
	// for a given pull request. It should return a RevalidationRequest or nil
	// to continue uninterrupted, and err if the request should be terminated
	OnPullDataSent(ChannelID, ChannelState) (RevalidationRequest, error)
	// OnPushDataReceived is called on the responder side when more bytes are received
	// for a given push request. It should return a RevalidationRequest or nil
	// to continue uninterrupted, and err if the request should be terminated
	OnPushDataReceived(ChannelID, ChannelState) (RevalidationRequest, error)
}

// RevalidationFulfiller produces new vouchers in response to revalidation requests
type RevalidationFulfiller interface {
	// FulfillPushVoucherRequest is called on the initiator side when a request
	// for revalidation is received for a push request
	FulfillPushRevalidationRequest(
		sender peer.ID,
		voucher RevalidationRequest,
		baseCid cid.Cid,
		selector ipld.Node) (Voucher, error)
	// FulfillPushVoucherRequest is called on the initiator side when a request
	// for revalidation is received for a push request
	FulfillPullRevalidationRequest(
		sender peer.ID,
		voucher RevalidationRequest,
		baseCid cid.Cid,
		selector ipld.Node) (Voucher, error)
}

// Manager is the core interface presented by all implementations of
// of the data transfer sub system
type Manager interface {
	// RegisterVoucherType registers a validator for the given voucher type
	// will error if voucher type does not implement voucher
	// or if there is a voucher type registered with an identical identifier
	RegisterVoucherType(voucherType Voucher, validator RequestValidator) error

	// RegisterRevalidator registers a revalidator for the given voucher type
	// Note: this is the voucher type used to revalidate. It can share a name
	// with the initial validator type and CAN be the same type, or a different type.
	// The revalidator can simply be the sampe as the original request validator,
	// or a different validator that satisfies the revalidator interface.
	RegisterRevalidator(voucherType Voucher, revalidator Revalidator) error

	// RegisterRevalidationFulfiller registers a fulfiller for revalidation requests
	RegisterRevalidationFulfiller(revalidationRequestType RevalidationRequest, fulfiller RevalidationFulfiller) error

	// RegisterVoucherResultType allows deserialization of a voucher result,
	// so that a listener can read the metadata
	RegisterVoucherResultType(resultType VoucherResult) error

	// open a data transfer that will send data to the recipient peer and
	// transfer parts of the piece that match the selector
	OpenPushDataChannel(ctx context.Context, to peer.ID, voucher Voucher, baseCid cid.Cid, selector ipld.Node) (ChannelID, error)

	// open a data transfer that will request data from the sending peer and
	// transfer parts of the piece that match the selector
	OpenPullDataChannel(ctx context.Context, to peer.ID, voucher Voucher, baseCid cid.Cid, selector ipld.Node) (ChannelID, error)

	// close an open channel (effectively a cancel)
	CloseDataTransferChannel(x ChannelID)

	// get status of a transfer
	TransferChannelStatus(x ChannelID) Status

	// get notified when certain types of events happen
	SubscribeToEvents(subscriber Subscriber) Unsubscribe

	// get all in progress transfers
	InProgressChannels() map[ChannelID]ChannelState
}
