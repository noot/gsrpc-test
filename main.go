package main

import (
	"fmt"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
	//"github.com/centrifuge/go-substrate-rpc-client/config"
	"github.com/centrifuge/go-substrate-rpc-client/signature"
	"github.com/centrifuge/go-substrate-rpc-client/types"
)

func main() {
	// This sample shows how to create a transaction to make a transfer from one an account to another.

	// Instantiate the API
	//api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	api, err := gsrpc.NewSubstrateAPI("http://localhost:8545")
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	// Create a call, transferring 12345 units to Bob
	bob, err := types.NewAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	if err != nil {
		panic(err)
	}

	c, err := types.NewCall(meta, "Balances.transfer", bob, types.NewUCompactFromUInt(12345))
	if err != nil {
		panic(err)
	}

	// Create the extrinsic
	ext := types.NewExtrinsic(c)
	if err != nil {
		panic(err)
	}

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		panic(err)
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		panic(err)
	}

	key, err := types.CreateStorageKey(meta, "System", "AccountNonce", signature.TestKeyringPairAlice.PublicKey, nil)
	if err != nil {
		panic(err)
	}

	var nonce uint32
	_, err = api.RPC.State.GetStorageLatest(key, &nonce)
	if err != nil {
		panic(err)
	}

	o := types.SignatureOptions{
		BlockHash:   genesisHash,
		Era:         types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash: genesisHash,
		Nonce:       types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion: rv.SpecVersion,
		Tip:         types.NewUCompactFromUInt(0),
	}

	// Sign the transaction using Alice's default account
	err = ext.Sign(signature.TestKeyringPairAlice, o)
	if err != nil {
		panic(err)
	}
	fmt.Println(ext)
	// Send the extrinsic
	hash, err := api.RPC.Author.SubmitExtrinsic(ext)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Transfer sent with hash %#x\n", hash)
}
