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

func dumpTxOut(tx *wire.MsgTx) {
	fmt.Printf("Start dumping output of tx %s\n", tx.TxHash().String())
	for i, txout := range tx.TxOut {
		scriptClass, addrs, sigNeed, err := txscript.ExtractPkScriptAddrs(txout.PkScript, &chaincfg.TestNet3Params)
		if err != nil {
			fmt.Printf("Failed to parse txout, err : %s\n", err)
			return
		}
		fmt.Printf("#%d of output :\n", i)
		fmt.Printf(">> Script Class : %s\n", scriptClass.String())
		fmt.Printf(">> Sig Need : %d\n", sigNeed)
		for _, addr := range addrs {
			fmt.Printf(">> To Addr : %s\n", addr.EncodeAddress())
		}
	}
	return
}

func getWalletUTXOs(wallet_id int, amount int64) []UTXO {
	return []UTXO{
		UTXO{
			TxHash: "34f920e779c3efd17ce0ab4c2cb2b49b83aa4bc693344e98c17bf52e72c742d2",
			Idx:    uint32(0),
		},
		UTXO{
			TxHash: "07548ae81f1ffd97df3d82dc6da8ee22a198201e9b5738c3a0e8ae973ccc6d2d",
			Idx:    uint32(0),
		},
	}
}

// scriptStr, err := txscript.DisasmString(addr.ScriptAddress())
// 			if err != nil {
// 				panic(err)
// 			}
// 			fmt.Printf(">> Script : %s", scriptStr)
