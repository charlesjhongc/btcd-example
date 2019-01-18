package main

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

func addOutput(tx *wire.MsgTx, address string, amount int64, net *chaincfg.Params) error {
	destAddr, err := btcutil.DecodeAddress(address, net)
	if err != nil {
		fmt.Printf("Failed to decode dest addr, %s\n", err)
		return err
	}
	// lib will auto determine type of addr and return coresponded locking script which is awesome!!!!
	lockingScript, err := txscript.PayToAddrScript(destAddr)
	if err != nil {
		fmt.Printf("Failed to create locking script, err : %s\n", err)
		return err
	}
	out := wire.NewTxOut(amount, lockingScript)
	tx.AddTxOut(out)
	return nil
}

func addChangeOutput(tx *wire.MsgTx, changeAddr string, change int64, net *chaincfg.Params) error {
	err := addOutput(tx, changeAddr, int64(change), &chaincfg.TestNet3Params)
	if err != nil {
		fmt.Printf("Failed to add output, %s\n", err)
		return err
	}
	return nil
}

func dumpTxOut(tx *wire.MsgTx) {
	for i, txout := range tx.TxOut {
		scriptClass, addrs, sigNeed, err := txscript.ExtractPkScriptAddrs(txout.PkScript, &chaincfg.TestNet3Params)
		if err != nil {
			fmt.Printf("Failed to parse txout, err : %s\n", err)
			return
		}
		fmt.Printf("################################\n")
		fmt.Printf("#%d of input\n", i)
		fmt.Printf("Script Class : %v\n", scriptClass.String())
		for _, addr := range addrs {
			fmt.Printf("Addr : %v\n", addr.EncodeAddress())
			//		scriptStr, err := txscript.DisasmString(addr.ScriptAddress())
			//		if err != nil {
			//			panic(err)
			//		}
			//		fmt.Printf("Script : %s", scriptStr)
		}
		fmt.Printf("Sig Need : %v\n", sigNeed)
		fmt.Printf("################################\n")
	}
	return
}
