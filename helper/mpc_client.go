package helper

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/metis-seq/themis/helper/tss"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

var mpcClientConn *grpc.ClientConn

func mpcConnect() error {
	var err error
	mpcClientConn, err = grpc.Dial(
		conf.MpcServerURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	return err
}

func GetMpcConn() *grpc.ClientConn {
	if mpcClientConn == nil && mpcClientConn.GetState() != connectivity.Ready {
		mpcConnect()
	}
	return mpcClientConn
}

func GetMpcKey(mpcId string) ([]byte, uint64, error) {
	mpcClient := tss.NewTssServiceClient(GetMpcConn())
	keyResp, err := mpcClient.GetKey(context.TODO(), &tss.GetKeyRequest{
		KeyId: mpcId,
	})
	if err != nil {
		return nil, 0, err
	}

	if keyResp == nil {
		return nil, 0, errors.New("mpc key not exist")
	}

	return keyResp.PublicKey, uint64(keyResp.Threshold), nil
}

func MpcCreate(reqeust *tss.KeyGenRequest) ([]byte, error) {
	mpcClient := tss.NewTssServiceClient(GetMpcConn())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	keyGenResp, err := mpcClient.KeyGen(ctx, reqeust)
	if err != nil {
		return nil, err
	}

	if keyGenResp == nil {
		return nil, errors.New("mpc key not exist")
	}

	return keyGenResp.PublicKey, nil
}

func MpcSign(signID, mpcID string, signMsg []byte) ([]byte, error) {
	mpcClient := tss.NewTssServiceClient(GetMpcConn())

	var signature []byte
	Retry(5, 1*time.Second, func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		request := &tss.KeySignRequest{
			KeyId:   mpcID,
			SignMsg: signMsg,
			SignId:  signID,
		}
		Logger.Info("MpcSign request", "mpcId", mpcID, "signId", signID)
		Logger.Debug("MpcSign request data", request.String())

		keySignResp, err := mpcClient.KeySign(ctx, request)
		if err != nil {
			Logger.Error("MpcSign failed", "err", err.Error())
			return err
		}

		if keySignResp == nil {
			Logger.Error("MpcSign result is nil")
			return errors.New("MpcSign result is nil")
		}

		signature = ConvertSignature(keySignResp.SignatureR, keySignResp.SignatureS, keySignResp.SignatureV)
		return nil
	})

	Logger.Info("MpcSign success", "signID", signID, "signature", fmt.Sprintf("%x", signature))
	return signature, nil
}

func ConvertSignature(r, s, v []byte) []byte {
	var signature []byte
	signature = append(signature, r...)
	signature = append(signature, s...)
	signature = append(signature, v...)

	if signature[crypto.RecoveryIDOffset] <= 1 {
		signature[crypto.RecoveryIDOffset] += 27 // Transform yellow paper V address 27/28 to 0/1
	}

	return signature
}

func ParseSignature(signature []byte) ([]byte, []byte, []byte) {
	var r, s, v []byte
	r = append(r, signature[:32]...)
	s = append(s, signature[32:64]...)
	v = append(v, signature[64:]...)

	if v[0] > 1 {
		v[0] -= 27
	}
	return r, s, v
}
