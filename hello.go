package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/nervosnetwork/ckb-sdk-go/crypto/secp256k1"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/nervosnetwork/ckb-sdk-go/transaction"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/nervosnetwork/ckb-sdk-go/utils"
)

func main() {
	client, err := rpc.Dial("http://127.0.0.1:8114")
	if err != nil {
		log.Fatalf("create rpc client error: %v", err)
	}
	// Convert the private key of Alice to an ecdsa private key.
	key, err := secp256k1.HexToKey("08730a367dfabcadb805d69e0e613558d5160eb8bab9d6e326980c2c46a05db2")
	if err != nil {
		log.Fatalf("import private key error: %v", err)
	}
    
	if err != nil {
		log.Fatalf("load system script error: %v", err)
	}
	//Decode the lock args of Bob to bytes.
	toAddress, _ := hex.DecodeString("ecbe30bcf5c6b2f2d8ec2dd229a4603a7e206b99")

	tx := &types.Transaction{
		Version:    0,
		HeaderDeps: []types.Hash{},
		CellDeps: []*types.CellDep{
				{
					OutPoint: &types.OutPoint{
						TxHash: types.HexToHash("0x6ddc6718014b7ad50121b95bb25ff61b4445b6c57ade514e7d08447e025f9f30"),
						Index:  0,
					},
				DepType:  "dep_group",
				},
		},
	}
	tx.Outputs = append(tx.Outputs, &types.CellOutput{
		Capacity: 20000000000,
		Lock: &types.Script{
			CodeHash: types.HexToHash("0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8"),
			HashType: types.HashTypeType,
			Args:     toAddress,
		},
	})
	changeAddress, _ := hex.DecodeString("6407c2ef9bd96e8e14ac4cd15d860e9331802172")
	tx.Outputs = append(tx.Outputs, &types.CellOutput{
		Capacity: 2007786346416,
		Lock: &types.Script{
			CodeHash: types.HexToHash("0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8"),
			HashType: types.HashTypeType,
			Args:     changeAddress, 
		},
	})
	tx.OutputsData = [][]byte{{}, {}}

	group, witnessArgs, err := transaction.AddInputsForTransaction(tx, []*types.CellInput{
		{
            Since: 0,
            PreviousOutput: &types.OutPoint{
                TxHash: types.HexToHash("0xace92ad1595ab435a13c095160b23ad8ea0dbb7cf7b8f7b7ef3540ec34372f94"),
                Index: 0,
            },
        },
	})

	if err != nil {
		log.Fatalf("add inputs to transaction error: %v", err)
	}

	err = transaction.SingleSignTransaction(tx, group, witnessArgs, key)
	if err != nil {
		log.Fatalf("sign transaction error: %v", err)
	}

	hash, err := client.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("send transaction error: %v", err)
	}

	fmt.Println(hash.String())
}