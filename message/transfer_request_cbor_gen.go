// Code generated by github.com/whyrusleeping/cbor-gen. DO NOT EDIT.

package message

import (
	"fmt"
	"io"

	datatransfer "github.com/filecoin-project/go-data-transfer"
	cbg "github.com/whyrusleeping/cbor-gen"
	xerrors "golang.org/x/xerrors"
)

var _ = xerrors.Errorf

func (t *transferRequest) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{137}); err != nil {
		return err
	}

	// t.BCid (cid.Cid) (struct)

	if t.BCid == nil {
		if _, err := w.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteCid(w, *t.BCid); err != nil {
			return xerrors.Errorf("failed to write cid field t.BCid: %w", err)
		}
	}

	// t.Canc (bool) (bool)
	if err := cbg.WriteBool(w, t.Canc); err != nil {
		return err
	}

	// t.Updt (bool) (bool)
	if err := cbg.WriteBool(w, t.Updt); err != nil {
		return err
	}

	// t.Part (bool) (bool)
	if err := cbg.WriteBool(w, t.Part); err != nil {
		return err
	}

	// t.Pull (bool) (bool)
	if err := cbg.WriteBool(w, t.Pull); err != nil {
		return err
	}

	// t.Stor (typegen.Deferred) (struct)
	if err := t.Stor.MarshalCBOR(w); err != nil {
		return err
	}

	// t.Vouch (typegen.Deferred) (struct)
	if err := t.Vouch.MarshalCBOR(w); err != nil {
		return err
	}

	// t.VTyp (datatransfer.TypeIdentifier) (string)
	if len(t.VTyp) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.VTyp was too long")
	}

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajTextString, uint64(len(t.VTyp)))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(t.VTyp)); err != nil {
		return err
	}

	// t.XferID (uint64) (uint64)

	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(t.XferID))); err != nil {
		return err
	}

	return nil
}

func (t *transferRequest) UnmarshalCBOR(r io.Reader) error {
	br := cbg.GetPeeker(r)

	maj, extra, err := cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 9 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.BCid (cid.Cid) (struct)

	{

		pb, err := br.PeekByte()
		if err != nil {
			return err
		}
		if pb == cbg.CborNull[0] {
			var nbuf [1]byte
			if _, err := br.Read(nbuf[:]); err != nil {
				return err
			}
		} else {

			c, err := cbg.ReadCid(br)
			if err != nil {
				return xerrors.Errorf("failed to read cid field t.BCid: %w", err)
			}

			t.BCid = &c
		}

	}
	// t.Canc (bool) (bool)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajOther {
		return fmt.Errorf("booleans must be major type 7")
	}
	switch extra {
	case 20:
		t.Canc = false
	case 21:
		t.Canc = true
	default:
		return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
	}
	// t.Updt (bool) (bool)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajOther {
		return fmt.Errorf("booleans must be major type 7")
	}
	switch extra {
	case 20:
		t.Updt = false
	case 21:
		t.Updt = true
	default:
		return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
	}
	// t.Part (bool) (bool)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajOther {
		return fmt.Errorf("booleans must be major type 7")
	}
	switch extra {
	case 20:
		t.Part = false
	case 21:
		t.Part = true
	default:
		return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
	}
	// t.Pull (bool) (bool)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajOther {
		return fmt.Errorf("booleans must be major type 7")
	}
	switch extra {
	case 20:
		t.Pull = false
	case 21:
		t.Pull = true
	default:
		return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
	}
	// t.Stor (typegen.Deferred) (struct)

	{

		pb, err := br.PeekByte()
		if err != nil {
			return err
		}
		if pb == cbg.CborNull[0] {
			var nbuf [1]byte
			if _, err := br.Read(nbuf[:]); err != nil {
				return err
			}
		} else {
			t.Stor = new(cbg.Deferred)
			if err := t.Stor.UnmarshalCBOR(br); err != nil {
				return xerrors.Errorf("unmarshaling t.Stor pointer: %w", err)
			}
		}

	}
	// t.Vouch (typegen.Deferred) (struct)

	{

		pb, err := br.PeekByte()
		if err != nil {
			return err
		}
		if pb == cbg.CborNull[0] {
			var nbuf [1]byte
			if _, err := br.Read(nbuf[:]); err != nil {
				return err
			}
		} else {
			t.Vouch = new(cbg.Deferred)
			if err := t.Vouch.UnmarshalCBOR(br); err != nil {
				return xerrors.Errorf("unmarshaling t.Vouch pointer: %w", err)
			}
		}

	}
	// t.VTyp (datatransfer.TypeIdentifier) (string)

	{
		sval, err := cbg.ReadString(br)
		if err != nil {
			return err
		}

		t.VTyp = datatransfer.TypeIdentifier(sval)
	}
	// t.XferID (uint64) (uint64)

	{

		maj, extra, err = cbg.CborReadHeader(br)
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.XferID = uint64(extra)

	}
	return nil
}
