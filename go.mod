module github.com/sloganking/sawtooth-devmode-go

go 1.18

// import branch
// replace github.com/hyperledger/sawtooth-sdk-go => github.com/taekion-org/sawtooth-sdk-go v0.1.2-0.20220629194300-aa402cec237c

replace github.com/hyperledger/sawtooth-sdk-go => ../sawtooth-sdk-go

require (
	github.com/hyperledger/sawtooth-sdk-go v0.0.0-00010101000000-000000000000
	github.com/pebbe/zmq4 v1.2.5
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/golang/protobuf v1.5.0 // indirect
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b // indirect
)
