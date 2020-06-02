package datatransfer

import (
	"context"
	"time"

	"github.com/ipfs/go-cid"
	ipldformat "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/ipld/go-ipld-prime"
	"github.com/libp2p/go-libp2p-core/peer"
	"golang.org/x/xerrors"

	datatransfer "github.com/filecoin-project/go-data-transfer"
	"github.com/filecoin-project/go-data-transfer/channels"
)

// This file implements a VERY simple, incomplete version of the data transfer
// module that allows us to make the necessary insertions of data transfer
// functionality into the storage market
// It does not:
// -- actually validate requests
// -- support Push requests
// -- support multiple subscribers
// -- do any actual network coordination or use Graphsync

type dagserviceImpl struct {
	dag        ipldformat.DAGService
	subscriber datatransfer.Subscriber
}

// NewDAGServiceDataTransfer returns a data transfer manager based on
// an IPLD DAGService
func NewDAGServiceDataTransfer(dag ipldformat.DAGService) datatransfer.Manager {
	return &dagserviceImpl{dag, nil}
}

// open a data transfer that will send data to the recipient peer and
// transfer parts of the piece that match the selector
func (impl *dagserviceImpl) OpenPushDataChannel(ctx context.Context, to peer.ID, voucher datatransfer.Voucher, baseCid cid.Cid, Selector ipld.Node) (datatransfer.ChannelID, error) {
	return datatransfer.ChannelID{}, xerrors.Errorf("not implemented")
}

// open a data transfer that will request data from the sending peer and
// transfer parts of the piece that match the selector
func (impl *dagserviceImpl) OpenPullDataChannel(ctx context.Context, to peer.ID, voucher datatransfer.Voucher, baseCid cid.Cid, Selector ipld.Node) (datatransfer.ChannelID, error) {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()
		err := merkledag.FetchGraph(ctx, baseCid, impl.dag)
		event := datatransfer.Event{Timestamp: time.Now()}
		if err != nil {
			event.Code = datatransfer.Error
			event.Message = err.Error()
		} else {
			event.Code = datatransfer.Complete
		}
		impl.subscriber(event, channels.ChannelState{Channel: channels.NewChannel(0, baseCid, Selector, voucher, to, "", 0)})
	}()
	return datatransfer.ChannelID{}, nil
}

// close an open channel (effectively a cancel)
func (impl *dagserviceImpl) CloseDataTransferChannel(x datatransfer.ChannelID) {}

// get status of a transfer
func (impl *dagserviceImpl) TransferChannelStatus(x datatransfer.ChannelID) datatransfer.Status {
	return datatransfer.ChannelNotFoundError
}

// get notified when certain types of events happen
func (impl *dagserviceImpl) SubscribeToEvents(subscriber datatransfer.Subscriber) datatransfer.Unsubscribe {
	impl.subscriber = subscriber
	return func() {}
}

// get all in progress transfers
func (impl *dagserviceImpl) InProgressChannels() map[datatransfer.ChannelID]datatransfer.ChannelState {
	return nil
}

// RegisterVoucherType registers a validator for the given voucher type
// will error if voucher type does not implement voucher
// or if there is a voucher type registered with an identical identifier
func (impl *dagserviceImpl) RegisterVoucherType(voucherType datatransfer.Voucher, validator datatransfer.RequestValidator) error {
	return nil
}

// RegisterRevalidator registers a revalidator for the given voucher type
// Note: this is the voucher type used to revalidate. It can share a name
// with the initial validator type and CAN be the same type, or a different type.
// The revalidator can simply be the sampe as the original request validator,
// or a different validator that satisfies the revalidator interface.
func (impl *dagserviceImpl) RegisterRevalidator(voucherType datatransfer.Voucher, revalidator datatransfer.Revalidator) error {
	return nil
}

// SendVoucher sends an intermediate voucher as needed when the receiver sends a request for revalidation
func (impl *dagserviceImpl) SendVoucher(ctx context.Context, channelID datatransfer.ChannelID, voucher datatransfer.Voucher) error {
	return nil
}

// RegisterVoucherResultType allows deserialization of a voucher result,
// so that a listener can read the metadata
func (impl *dagserviceImpl) RegisterVoucherResultType(resultType datatransfer.VoucherResult) error {
	return nil
}
