package main

import (
	"github.com/hyperledger/sawtooth-sdk-go/consensus"
	"math/rand"
	"time"
)

func init() {
	// initialize RNG
	rand.Seed(time.Now().UnixNano())
}

func main() {
	engine := consensus.NewConsensusEngine("tcp://tesla:5050", &DevmodeEngineImpl{})
	engine.Start()
}
