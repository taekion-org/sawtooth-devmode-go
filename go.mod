module github.com/sloganking/sawtooth-devmode-go

go 1.18

// import branch
// replace github.com/hyperledger/sawtooth-sdk-go => github.com/taekion-org/sawtooth-sdk-go v0.1.2-0.20220629194300-aa402cec237c

replace github.com/hyperledger/sawtooth-sdk-go => ../sawtooth-sdk-go

require (
	github.com/hyperledger/sawtooth-sdk-go v0.0.0-00010101000000-000000000000
	google.golang.org/protobuf v1.28.0
)
