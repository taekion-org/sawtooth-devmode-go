package main

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/hyperledger/sawtooth-sdk-go/consensus"
	// "github.com/hyperledger/sawtooth-sdk-go/messaging"
	// consensus_pb2 "github.com/hyperledger/sawtooth-sdk-go/protobuf/consensus_pb2"
	// zmq "github.com/pebbe/zmq4"
)

const DEFAULT_WAIT_TIME = 0

type LogGuard struct {
	not_ready_to_summarize bool
	not_ready_to_finalize  bool
}

type devmodeService struct {
	service  consensus.ConsensusService
	logGuard LogGuard
}

//> impl devmodeService

	func (self devmodeService) new(service consensus.ConsensusService) devmodeService {
		return devmodeService{
			service: service,
			logGuard: LogGuard{
				not_ready_to_summarize: false,
				not_ready_to_finalize:  false,
			},
		}
	}

	// fn get_chain_head(&mut self) -> Block {
	//     debug!("Getting chain head");
	//     self.service
	//         .get_chain_head()
	//         .expect("Failed to get chain head")
	// }

	func (self devmodeService) get_chain_head() consensus.Block {
		head, err := self.service.GetChainHead()
		if err != nil {
			panic("Failed to get chain head")
		}
		return head
	}

	// fn get_block(&mut self, block_id: &BlockId) -> Block {
	//     debug!("Getting block {}", to_hex(&block_id));
	//     self.service
	//         .get_blocks(vec![block_id.clone()])
	//         .expect("Failed to get block")
	//         .remove(block_id)
	//         .unwrap()
	// }

	func (self devmodeService) get_block(block_id consensus.BlockId) consensus.Block {
		ids := make([]consensus.BlockId, 0)
		ids = append(ids, block_id)
		blockIdMap, err := self.service.GetBlocks(ids)
		if err != nil {
			panic("Failed to get block")
		}
		block := blockIdMap[block_id]
		return block
	}

	// fn initialize_block(&mut self) {
	//     debug!("Initializing block");
	//     self.service
	//         .initialize_block(None)
	//         .expect("Failed to initialize");
	// }

	func (self devmodeService) initialize_block() {
		err := self.service.InitializeBlock(consensus.BlockId{})
		if err != nil {
			panic("Failed to initialize")
		}
	}

	// fn finalize_block(&mut self) -> BlockId {
	//     debug!("Finalizing block");
	//     let mut summary = self.service.summarize_block();
	//     while let Err(Error::BlockNotReady) = summary {
	//         if !self.log_guard.not_ready_to_summarize {
	//             self.log_guard.not_ready_to_summarize = true;
	//             debug!("Block not ready to summarize");
	//         }
	//         sleep(time::Duration::from_secs(1));
	//         summary = self.service.summarize_block();
	//     }
	//     self.log_guard.not_ready_to_summarize = false;
	//     let summary = summary.expect("Failed to summarize block");
	//     debug!("Block has been summarized successfully");

	//     let consensus: Vec<u8> = create_consensus(&summary);
	//     let mut block_id = self.service.finalize_block(consensus.clone());
	//     while let Err(Error::BlockNotReady) = block_id {
	//         if !self.log_guard.not_ready_to_finalize {
	//             self.log_guard.not_ready_to_finalize = true;
	//             debug!("Block not ready to finalize");
	//         }
	//         sleep(time::Duration::from_secs(1));
	//         block_id = self.service.finalize_block(consensus.clone());
	//     }
	//     self.log_guard.not_ready_to_finalize = false;
	//     let block_id = block_id.expect("Failed to finalize block");
	//     debug!(
	//         "Block has been finalized successfully: {}",
	//         to_hex(&block_id)
	//     );

	//     block_id
	// }

	func (self devmodeService) finalize_block() {
		summary, err := self.service.SummarizeBlock()
	}

	// fn check_block(&mut self, block_id: BlockId) {
	//     debug!("Checking block {}", to_hex(&block_id));
	//     self.service
	//         .check_blocks(vec![block_id])
	//         .expect("Failed to check block");
	// }

	func (self devmodeService) check_block(block_id consensus.BlockId) {
		blocks := []consensus.BlockId{block_id}
		err := self.service.CheckBlocks(blocks)
		if err != nil {
			panic("Failed to check block")
		}
	}

	// fn fail_block(&mut self, block_id: BlockId) {
	//     debug!("Failing block {}", to_hex(&block_id));
	//     self.service
	//         .fail_block(block_id)
	//         .expect("Failed to fail block");
	// }

	func (self devmodeService) fail_block(block_id consensus.BlockId) {
		err := self.service.FailBlock(block_id)
		if err != nil {
			panic("Failed to fail block")
		}
	}

	// fn ignore_block(&mut self, block_id: BlockId) {
	//     debug!("Ignoring block {}", to_hex(&block_id));
	//     self.service
	//         .ignore_block(block_id)
	//         .expect("Failed to ignore block")
	// }

	func (self devmodeService) ignore_block(block_id consensus.BlockId) {
		err := self.service.IgnoreBlock(block_id)
		if err != nil {
			panic("Failed to ignore block")
		}
	}

	// fn commit_block(&mut self, block_id: BlockId) {
	//     debug!("Committing block {}", to_hex(&block_id));
	//     self.service
	//         .commit_block(block_id)
	//         .expect("Failed to commit block");
	// }

	func (self devmodeService) commit_block(block_id consensus.BlockId) {
		err := self.service.CommitBlock(block_id)
		if err != nil {
			panic("Failed to commit block")
		}
	}

	// fn cancel_block(&mut self) {
	//     debug!("Canceling block");
	//     match self.service.cancel_block() {
	//         Ok(_) => {}
	//         Err(Error::InvalidState(_)) => {}
	//         Err(err) => {
	//             panic!("Failed to cancel block: {:?}", err);
	//         }
	//     };
	// }

	func (self devmodeService) cancel_block() {
		err := self.service.CancelBlock()
		if err != nil {
			if err.(consensus.Error).ErrorEnum == consensus.InvalidState {

			}
		}
	}

	// fn broadcast_published_block(&mut self, block_id: BlockId) {
	//     debug!("Broadcasting published block: {}", to_hex(&block_id));
	//     self.service
	//         .broadcast("published", block_id)
	//         .expect("Failed to broadcast published block");
	// }

	func (self devmodeService) broadcast_published_block(block_id consensus.BlockId) {
		err := self.service.Broadcast("published", []byte(block_id))
		if err != nil {
			panic("Failed to broadcast published block")
		}
	}

	// fn send_block_received(&mut self, block: &Block) {
	//     let block = block.clone();

	//     self.service
	//         .send_to(&block.signer_id, "received", block.block_id)
	//         .expect("Failed to send block received");
	// }

	func (self devmodeService) send_block_received(block consensus.Block) {
		err := self.service.SendTo(block.Signer_id, "recieved", []uint8(block.BlockId))
		if err != nil {
			panic("Failed to send block received")
		}
	}

	// fn send_block_ack(&mut self, sender_id: &PeerId, block_id: BlockId) {
	//     self.service
	//         .send_to(&sender_id, "ack", block_id)
	//         .expect("Failed to send block ack");
	// }

	func (self devmodeService) send_block_ack(sender_id consensus.PeerId, block_id consensus.BlockId) {
		err := self.service.SendTo(sender_id, "ack", []uint8(block_id))
		if err != nil {
			panic("Failed to send block ack")
		}
	}

	// // Calculate the time to wait between publishing blocks. This will be a
	// // random number between the settings sawtooth.consensus.min_wait_time and
	// // sawtooth.consensus.max_wait_time if max > min, else DEFAULT_WAIT_TIME. If
	// // there is an error parsing those settings, the time will be
	// // DEFAULT_WAIT_TIME.
	// fn calculate_wait_time(&mut self, chain_head_id: BlockId) -> time::Duration {
	//     let settings_result = self.service.get_settings(
	//         chain_head_id,
	//         vec![
	//             String::from("sawtooth.consensus.min_wait_time"),
	//             String::from("sawtooth.consensus.max_wait_time"),
	//         ],
	//     );

	//     let wait_time = if let Ok(settings) = settings_result {
	//         let ints: Vec<u64> = vec![
	//             &settings["sawtooth.consensus.min_wait_time"],
	//             &settings["sawtooth.consensus.max_wait_time"],
	//         ]
	//         .iter()
	//         .map(|string| string.parse::<u64>())
	//         .map(|result| result.unwrap_or(0))
	//         .collect();

	//         let min_wait_time: u64 = ints[0];
	//         let max_wait_time: u64 = ints[1];

	//         debug!("Min: {:?} -- Max: {:?}", min_wait_time, max_wait_time);

	//         if min_wait_time >= max_wait_time {
	//             DEFAULT_WAIT_TIME
	//         } else {
	//             rand::thread_rng().gen_range(min_wait_time, max_wait_time)
	//         }
	//     } else {
	//         DEFAULT_WAIT_TIME
	//     };

	//     info!("Wait time: {:?}", wait_time);

	//     time::Duration::from_secs(wait_time)
	// }

	func (self devmodeService) calculate_wait_time(chain_head_id consensus.BlockId) time.Duration {
		settings, err := self.service.GetSettings(chain_head_id, []string{"sawtooth.consensus.min_wait_time", "sawtooth.consensus.max_wait_time"})

		wait_time := 0

		if err != nil {

			// get min_wait_time
			min_wait_time_string := settings["sawtooth.consensus.min_wait_time"]
			min_wait_time, err := strconv.Atoi(min_wait_time_string)
			if err != nil {
				min_wait_time = 0
			}

			// get max_wait_time
			max_wait_time_string := settings["sawtooth.consensus.max_wait_time"]
			max_wait_time, err := strconv.Atoi(max_wait_time_string)
			if err != nil {
				max_wait_time = 0
			}

			if min_wait_time >= max_wait_time {
				wait_time = DEFAULT_WAIT_TIME
			} else {
				// wait_time = value between min_wait_time (inclusive) and max_wait_time (exclusive)
				rand_range := max_wait_time - min_wait_time
				wait_time = rand.Intn(rand_range) + min_wait_time
			}

		} else {
			wait_time = DEFAULT_WAIT_TIME
		}

		// Convert WAIT_TIME from seconds to nanoseconds so we can store it in a time.Duration
		var duration time.Duration = time.Duration(wait_time) * 1000000000

		return duration
	}

//<

type DevmodeEngine struct{}

//> impl ConsensusEngineImpl for DevmodeEngine
	func (self DevmodeEngine) Version() string {
		return "0.1"
	}

	func (self DevmodeEngine) Name() string {
		return "Devmode"
	}

	func (self DevmodeEngine) Startup(startupState consensus.StartupState, service consensus.ConsensusService) {}
	func (self DevmodeEngine) Shutdown() {}
	func (self DevmodeEngine) HandlePeerConnected(peerInfo consensus.PeerInfo) {}
	func (self DevmodeEngine) HandlePeerDisconnected(peerInfo consensus.PeerInfo) {}
	func (self DevmodeEngine) HandlePeerMessage(peerMessage consensus.PeerMessage) {}
	func (self DevmodeEngine) HandleBlockNew(block consensus.Block) {}
	func (self DevmodeEngine) HandleBlockValid(blockId consensus.BlockId) {}
	func (self DevmodeEngine) HandleBlockInvalid(blockId consensus.BlockId) {}
	func (self DevmodeEngine) HandleBlockCommit(blockId consensus.BlockId) {}
//<

// // The main worker thread finds an appropriate handler and processes the request
// func worker(context *zmq.Context, uri string, queue <-chan *validator_pb2.Message, done chan<- bool, handlers []TransactionHandler) {
// 	// Connect to the main send/receive thread
// 	connection, err := messaging.NewConnection(context, zmq.DEALER, uri, false)
// 	if err != nil {
// 		logger.Errorf("Failed to connect to main thread: %v", err)
// 		done <- false
// 		return
// 	}
// 	defer connection.Close()
// 	id := connection.Identity()

// 	// Receive work off of the queue until the queue is closed
// 	for msg := range queue {
// 		request := &processor_pb2.TpProcessRequest{}
// 		err = proto.Unmarshal(msg.GetContent(), request)
// 		if err != nil {
// 			logger.Errorf(
// 				"(%v) Failed to unmarshal TpProcessRequest: %v", id, err,
// 			)
// 			break
// 		}

// 		header := request.GetHeader()

// 		// Try to find a handler
// 		handler, err := findHandler(handlers, header)
// 		if err != nil {
// 			logger.Errorf("(%v) Failed to find handler: %v", id, err)
// 			break
// 		}

// 		// Construct a new Context instance for the handler
// 		contextId := request.GetContextId()
// 		context := NewContext(connection, contextId)

// 		// Run the handler
// 		err = handler.Apply(request, context)

// 		// Process the handler response
// 		response := &processor_pb2.TpProcessResponse{}
// 		if err != nil {
// 			switch e := err.(type) {
// 			case *InvalidTransactionError:
// 				logger.Warnf("(%v) %v", id, e)
// 				response.Status = processor_pb2.TpProcessResponse_INVALID_TRANSACTION
// 				response.Message = e.Msg
// 				response.ExtendedData = e.ExtendedData
// 			case *InternalError:
// 				logger.Warnf("(%v) %v", id, e)
// 				response.Status = processor_pb2.TpProcessResponse_INTERNAL_ERROR
// 				response.Message = e.Msg
// 				response.ExtendedData = e.ExtendedData
// 			case *AuthorizationException:
// 				logger.Warnf("(%v) %v", id, e)
// 				response.Status = processor_pb2.TpProcessResponse_INVALID_TRANSACTION
// 				response.Message = e.Msg
// 				response.ExtendedData = e.ExtendedData
// 			default:
// 				logger.Errorf("(%v) Unknown error: %v", id, err)
// 				response.Status = processor_pb2.TpProcessResponse_INTERNAL_ERROR
// 				response.Message = e.Error()
// 			}
// 		} else {
// 			response.Status = processor_pb2.TpProcessResponse_OK
// 		}

// 		responseData, err := proto.Marshal(response)
// 		if err != nil {
// 			logger.Errorf("(%v) Failed to marshal TpProcessResponse: %v", id, err)
// 			break
// 		}

// 		// Send back a response to the validator
// 		err = connection.SendMsg(
// 			validator_pb2.Message_TP_PROCESS_RESPONSE,
// 			responseData, msg.GetCorrelationId(),
// 		)
// 		if err != nil {
// 			logger.Errorf("(%v) Error sending TpProcessResponse: %v", id, err)
// 			break
// 		}
// 	}

// 	done <- true
// }