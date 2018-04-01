package ethereum

import (
    "github.com/ethereum/go-ethereum/common"
    "fmt"
    "github.com/ethereum/go-ethereum/crypto"
)

func EcRecover(hash, sig []byte) (common.Address, error) {
    if len(sig) != 65 {
        return common.Address{}, fmt.Errorf("signature must be 65 bytes long")
    }
    if sig[64] != 27 && sig[64] != 28 {
        return common.Address{}, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
    }
    sig[64] -= 27 // Transform yellow paper V from 27/28 to 0/1

    rpk, err := crypto.Ecrecover(hash, sig)
    if err != nil {
        return common.Address{}, err
    }
    pubKey := crypto.ToECDSAPub(rpk)
    recoveredAddr := crypto.PubkeyToAddress(*pubKey)
    return recoveredAddr, nil
}
