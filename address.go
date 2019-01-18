package main

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/getamis/sirius/log"
)

func createP2PKAddress(net *chaincfg.Params) (string, error) {
	privateKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		log.Error("Failed to create new ECDSA private key", "err", err)
		return "", err
	}
	fmt.Printf("New Private key was created, private : %X\n", privateKey.Serialize())
	pubKeyHash := btcutil.Hash160(privateKey.PubKey().SerializeCompressed())
	addr, err := btcutil.NewAddressPubKeyHash(pubKeyHash, net)
	if err != nil {
		log.Error("Failed to create P2PK address", "err", err)
		return "", err
	}
	fmt.Printf("New P2PK address was created, address : %s\n", addr.EncodeAddress())
	return addr.EncodeAddress(), nil
}

func createP2PKHSegwits(net *chaincfg.Params) (string, error) {
	return "", nil
}
