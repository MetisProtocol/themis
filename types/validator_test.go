package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// valInput struct is used to seed data for testing
// if the need arises it can be ported to the main build
type valInput struct {
	id         ValidatorID
	startBatch uint64
	endBatch   uint64
	power      int64
	pubKey     PubKey
	signer     ThemisAddress
}

func TestNewValidator(t *testing.T) {
	t.Parallel()

	// valCase created so as to pass it to assertPanics func,
	// ideally would like to get rid of this and pass the function directly
	tc := []struct {
		in  valInput
		out *Validator
		msg string
	}{
		{
			in: valInput{
				id:     ValidatorID(uint64(0)),
				signer: BytesToThemisAddress([]byte("12345678909876543210")),
			},
			out: &Validator{Signer: BytesToThemisAddress([]byte("12345678909876543210"))},
			msg: "testing for exact ThemisAddress",
		},
		{
			in: valInput{
				id:     ValidatorID(uint64(0)),
				signer: BytesToThemisAddress([]byte("1")),
			},
			out: &Validator{Signer: BytesToThemisAddress([]byte("1"))},
			msg: "testing for small ThemisAddress",
		},
		{
			in: valInput{
				id:     ValidatorID(uint64(0)),
				signer: BytesToThemisAddress([]byte("123456789098765432101")),
			},
			out: &Validator{Signer: BytesToThemisAddress([]byte("123456789098765432101"))},
			msg: "testing for excessively long ThemisAddress, max length is supposed to be 20",
		},
	}
	for _, c := range tc {
		out := NewValidator(c.in.id, c.in.startBatch, c.in.endBatch, 1, c.in.power, c.in.pubKey, c.in.signer)
		assert.Equal(t, c.out, out)
	}
}

// TestSortValidatorByAddress am populating only the signer as that is the only value used in sorting
func TestSortValidatorByAddress(t *testing.T) {
	t.Parallel()

	tc := []struct {
		in  []Validator
		out []Validator
		msg string
	}{
		{
			in: []Validator{
				{Signer: BytesToThemisAddress([]byte("3"))},
				{Signer: BytesToThemisAddress([]byte("2"))},
				{Signer: BytesToThemisAddress([]byte("1"))},
			},
			out: []Validator{
				{Signer: BytesToThemisAddress([]byte("1"))},
				{Signer: BytesToThemisAddress([]byte("2"))},
				{Signer: BytesToThemisAddress([]byte("3"))},
			},
			msg: "reverse sorting of validator objects",
		},
	}
	for i, c := range tc {
		out := SortValidatorByAddress(c.in)
		assert.Equal(t, c.out, out, fmt.Sprintf("i: %v, case: %v", i, c.msg))
	}
}

func TestValidateBasic(t *testing.T) {
	t.Parallel()

	neg1, uNeg1 := uint64(1), uint64(0)
	uNeg1 = uNeg1 - neg1
	tc := []struct {
		in  Validator
		out bool
		msg string
	}{
		{
			in:  Validator{StartBatch: 1, EndBatch: 5, Nonce: 0, PubKey: NewPubKey([]byte("nonZeroTestPubKey")), Signer: BytesToThemisAddress([]byte("3"))},
			out: true,
			msg: "Valid basic validator test",
		},
		{
			in:  Validator{StartBatch: 1, EndBatch: 5, Nonce: 0, PubKey: NewPubKey([]byte("")), Signer: BytesToThemisAddress([]byte("3"))},
			out: false,
			msg: "Invalid PubKey \"\"",
		},
		{
			in:  Validator{StartBatch: 1, EndBatch: 5, Nonce: 0, PubKey: ZeroPubKey, Signer: BytesToThemisAddress([]byte("3"))},
			out: false,
			msg: "Invalid PubKey",
		},

		//		{
		//			in:  Validator{StartBatch: uNeg1, EndBatch: 5, PubKey: NewPubKey([]byte("nonZeroTestPubKey")), Signer: BytesToThemisAddress([]byte("3"))},
		//			out: false,
		//			msg: "Invalid StartBatch",
		//		},
		{
			// do we allow for endBatch to be smaller than startBatch ??
			in:  Validator{StartBatch: 1, EndBatch: uNeg1, Nonce: 0, PubKey: NewPubKey([]byte("nonZeroTestPubKey")), Signer: BytesToThemisAddress([]byte("3"))},
			out: false,
			msg: "Invalid endBatch",
		},
		{
			// in:  Validator{StartBatch: 1, EndBatch: 1, PubKey: NewPubKey([]byte("nonZeroTestPubKey")), Signer: ThemisAddress(BytesToThemisAddress([]byte(string(""))))},
			in:  Validator{StartBatch: 1, EndBatch: 1, Nonce: 0, PubKey: NewPubKey([]byte("nonZeroTestPubKey")), Signer: BytesToThemisAddress([]byte(""))},
			out: false,
			msg: "Invalid Signer",
		},
		{
			in:  Validator{},
			out: false,
			msg: "Valid basic validator test",
		},
	}

	for _, c := range tc {
		out := c.in.ValidateBasic()
		assert.Equal(t, c.out, out, c.msg)
	}
}
