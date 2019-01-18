package main

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

const (
	btc_node_ip = "10.1.181.126:18332"
	fee_rate    = int64(20)
)

type UTXO struct {
	Tx     *btcutil.Tx
	TxHash string
	Idx    uint32
}

var addrPk = map[string]string{
	// P2PK
	"mi7Xxy8KzyiHfbhajZ5v6oUvr8VkKCaxNs": "private_key",
	"myrkeQzeBxBG6rvGBeyxkeEbJLJ8CZ3Fp5": "private_key",

	// P2PKH
	"mugNaeKrJYTrVpvbjeUCXbkZE127gzMEhD": "private_key",
	"mrSN4Peu5VV8ZVMoWojn4UtvhvarwWCohV": "private_key",
}

func main() {
	utxos := []UTXO{
		UTXO{
			TxHash: "34f920e779c3efd17ce0ab4c2cb2b49b83aa4bc693344e98c17bf52e72c742d2",
			Idx:    uint32(0),
		},
		UTXO{
			TxHash: "07548ae81f1ffd97df3d82dc6da8ee22a198201e9b5738c3a0e8ae973ccc6d2d",
			Idx:    uint32(0),
		},
	}

	client, err := rpcclient.New(&rpcclient.ConnConfig{
		HTTPPostMode: true,
		DisableTLS:   true,
		Host:         btc_node_ip,
		User:         "user",
		Pass:         "nono",
	}, nil)
	if err != nil {
		fmt.Printf("Failed to create bitcoin rpc client, err : %s\n", err)
		return
	}
	block_count, err := client.GetBlockCount()
	if err != nil {
		fmt.Printf("Failed to get block count from btcd, err : %s\n", err)
		return
	}
	fmt.Printf("Block count : %v\n", block_count)

	//txHash, err := chainhash.NewHashFromStr("f94bd770f427ffa9234fa229bbbf7314b82f967346928ebe5bd6e838acba482b")
	//if err != nil {
	//	fmt.Printf("Failed to parse hash str, err : %s\n", err)
	//	return
	//}
	//tx, err := client.GetRawTransaction(txHash)
	//if err != nil {
	//	fmt.Printf("Failed to get raw tx from btcd, err : %s\n", err)
	//	return
	//}
	//dumpTxOut(tx.MsgTx())

	// Start constructing a tx
	// Init a tx scruct
	msgtx := wire.NewMsgTx(wire.TxVersion)
	// Add input
	var inputToalValue int64
	var a, b, c int
	for i := 0; i < len(utxos); i++ {
		utxoHashOjb, err := chainhash.NewHashFromStr(utxos[i].TxHash)
		if err != nil {
			fmt.Printf("Failed to parse hash str, err : %s\n", err)
			return
		}

		// get tx using tx hash (in order to use it's output's script)
		utxos[i].Tx, err = client.GetRawTransaction(utxoHashOjb)
		if err != nil {
			fmt.Printf("Failed to get tx by hash, err : %s\n", err)
			return
		}
		inputToalValue += utxos[i].Tx.MsgTx().TxOut[utxos[i].Idx].Value

		// Add input (need tx hash, idx of output at this point)
		in := wire.NewTxIn(wire.NewOutPoint(utxoHashOjb, utxos[i].Idx), nil, nil)
		msgtx.AddTxIn(in)

		class := txscript.GetScriptClass(utxos[i].Tx.MsgTx().TxOut[utxos[i].Idx].PkScript)
		if class == txscript.PubKeyHashTy {
			a++
		} else if class == txscript.WitnessV0PubKeyHashTy {
			b++
		}
		// TODO how to know wether it's nested p2wpkh or not???
	}

	// Add output
	var outputToalValue int64
	err = addOutput(msgtx, "mrSN4Peu5VV8ZVMoWojn4UtvhvarwWCohV", int64(5566), &chaincfg.TestNet3Params)
	if err != nil {
		fmt.Printf("Failed to add output, %s\n", err)
		return
	}
	outputToalValue += 5566

	// Add change
	estTxSize := EstimateVirtualSize(a, b, c, msgtx.TxOut, true)
	fmt.Printf("EST TX size before add change output : %d\n", estTxSize)
	fee := int64(estTxSize) * fee_rate
	change := inputToalValue - outputToalValue - fee
	err = addChangeOutput(msgtx, "mugNaeKrJYTrVpvbjeUCXbkZE127gzMEhD", change, &chaincfg.TestNet3Params)

	// Sign TX
	for k := 0; k < len(msgtx.TxIn); k++ {
		// the third parameter is a little bit confused, my current conclusion is the index of coresponded input
		/*		sigScript, err := txscript.SignTxOutput(&chaincfg.TestNet3Params,
				msgtx, 0, prevTX.MsgTx().TxOut[1].PkScript, txscript.SigHashAll,
				txscript.KeyClosure(lookupKey), nil, nil)*/
		utxo := utxos[k]
		sigScript, err := txscript.SignTxOutput(&chaincfg.TestNet3Params,
			msgtx, k, utxo.Tx.MsgTx().TxOut[utxo.Idx].PkScript, txscript.SigHashAll, txscript.KeyClosure(lookupKey), nil, nil)
		if err != nil {
			fmt.Printf("Failed to create sigScript, err : %s\n", err)
			return
		}
		msgtx.TxIn[k].SignatureScript = sigScript

		// Verify signature is legit
		flags := txscript.ScriptBip16 | txscript.ScriptVerifyDERSignatures |
			txscript.ScriptStrictMultiSig |
			txscript.ScriptDiscourageUpgradableNops
		vm, err := txscript.NewEngine(utxo.Tx.MsgTx().TxOut[utxo.Idx].PkScript, msgtx, k, flags, nil, nil, -1)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := vm.Execute(); err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Printf("Transaction successfully signed, size : %d\n", msgtx.SerializeSizeStripped())

	// Brocast TX
	//hash, err := client.SendRawTransaction(msgtx, true)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("Fuck yeah !!! tx hash : %s\n", hash.String())
	return
}

func lookupKey(a btcutil.Address) (*btcec.PrivateKey, bool, error) {
	addrStr := a.EncodeAddress()
	pkStr, found := addrPk[addrStr]
	if !found {
		fmt.Printf("We don't have pk of addr %s\n", addrStr)
		return nil, false, errors.New("PK not found")
	}
	pkBytes, err := hex.DecodeString(pkStr)
	if err != nil {
		panic(err)
	}

	private, _ := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)
	return private, true, nil
}
