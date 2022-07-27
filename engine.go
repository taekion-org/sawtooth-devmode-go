package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"time"

	"github.com/hyperledger/sawtooth-sdk-go/consensus"
	"github.com/hyperledger/sawtooth-sdk-go/logging"
	// "github.com/hyperledger/sawtooth-sdk-go/messaging"
	// consensus_pb2 "github.com/hyperledger/sawtooth-sdk-go/protobuf/consensus_pb2"
	// zmq "github.com/pebbe/zmq4"
)

var logger *logging.Logger = logging.Get()

const DEFAULT_WAIT_TIME = 0

type LogGuard struct {
	not_ready_to_summarize bool
	not_ready_to_finalize  bool
}

type DevmodeService struct {
	service  consensus.ConsensusService
	logGuard LogGuard
}

// Creates a new DevmodeService
func NewDevmodeService(service consensus.ConsensusService) *DevmodeService {
	return &DevmodeService{
		service: service,
		logGuard: LogGuard{
			not_ready_to_summarize: false,
			not_ready_to_finalize:  false,
		},
	}
}

//> impl DevmodeService

	func (self *DevmodeService) getChainHead() consensus.Block {
		logger.Debug("Getting chain head")
		head, err := self.service.GetChainHead()
		if err != nil {
			panic("Failed to get chain head")
		}
		return head
	}

	func (self *DevmodeService) getBlock(block_id consensus.BlockId) consensus.Block {
		logger.Debugf("Getting block %v", block_id)
		ids := make([]consensus.BlockId, 0)
		ids = append(ids, block_id)
		blockIdMap, err := self.service.GetBlocks(ids)
		if err != nil {
			panic("Failed to get block")
		}
		block := blockIdMap[block_id]
		return block
	}

	func (self *DevmodeService) initializeBlock() {
		logger.Debug("Initializing block")
		err := self.service.InitializeBlock(consensus.BLOCK_ID_NULL)
		if err != nil {
			panic("Failed to initialize")
		}
	}

	// fn finalize_block(&mut self) -> BlockId {
	//     debug!("Finalizing block");
	//     let mut summary = self.service.summarize_block();
	//     while let Err(Error::BlockNotReadyError) = summary {
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
	//     while let Err(Error::BlockNotReadyError) = block_id {
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

	func (self *DevmodeService) finalizeBlock() consensus.BlockId {
		logger.Debug("Finalizing block")
		var err error

		summary, err := self.service.SummarizeBlock()
		for err != nil && consensus.IsBlockNotReadyError(err) {
			if !self.logGuard.not_ready_to_summarize {
				self.logGuard.not_ready_to_summarize = true
				logger.Debug("Block not ready to summarize")
			}
			time.Sleep(time.Second * 1)
			summary, err = self.service.SummarizeBlock()
		}
		self.logGuard.not_ready_to_summarize = false
		logger.Debug("Block has been summarized successfully")

		newConsensus := createConsensus(summary)
		blockId, err := self.service.FinalizeBlock(newConsensus)
		for err != nil && consensus.IsBlockNotReadyError(err) {
			if !self.logGuard.not_ready_to_finalize {
				self.logGuard.not_ready_to_finalize = true
				logger.Debug("Block not ready to finalize")
			}
			time.Sleep(time.Second * 1)
			blockId, err = self.service.FinalizeBlock(newConsensus)
		}
		self.logGuard.not_ready_to_finalize = false

		logger.Debugf("Block has been finalized successfully: %v", blockId)

		return blockId
	}

	func (self *DevmodeService) checkBlock(block_id consensus.BlockId) {
		logger.Debugf("Checking block %v", block_id)
		blocks := []consensus.BlockId{block_id}

		err := self.service.CheckBlocks(blocks)
		if err != nil {
			panic("Failed to check block")
		}
	}

	func (self *DevmodeService) failBlock(block_id consensus.BlockId) {
		logger.Debugf("Failing block %v", block_id)
		err := self.service.FailBlock(block_id)
		if err != nil {
			panic("Failed to fail block")
		}
	}

	func (self *DevmodeService) ignoreBlock(block_id consensus.BlockId) {
		logger.Debugf("Ignoring block %v", block_id)
		err := self.service.IgnoreBlock(block_id)
		if err != nil {
			panic("Failed to ignore block")
		}
	}

	func (self *DevmodeService) commitBlock(block_id consensus.BlockId) {
		logger.Debugf("Committing block %v", block_id)
		err := self.service.CommitBlock(block_id)
		if err != nil {
			panic("Failed to commit block")
		}
	}

	// fn cancel_block(&mut self) {
	//     debug!("Canceling block");
	//     match self.service.cancel_block() {
	//         Ok(_) => {}
	//         Err(Error::InvalidStateError(_)) => {}
	//         Err(err) => {
	//             panic!("Failed to cancel block: {:?}", err);
	//         }
	//     };
	// }

	func (self *DevmodeService) cancelBlock() {
		logger.Debug("Canceling block")
		err := self.service.CancelBlock()
		if err != nil && !consensus.IsInvalidStateError(err) {
			panic(fmt.Sprintf("Failed to cancel block: %v", err))
		}
	}

	func (self *DevmodeService) broadcast_published_block(block_id consensus.BlockId) {
		logger.Debugf("Broadcasting published block: %v", block_id)
		err := self.service.Broadcast("published", block_id.AsBytes())
		if err != nil {
			panic("Failed to broadcast published block")
		}
	}

	func (self DevmodeService) sendBlockReceived(block consensus.Block) {
		blockId := block.BlockId()
		err := self.service.SendTo(block.SignerId(), "received", blockId.AsBytes())
		if err != nil {
			panic("Failed to send block received")
		}
	}

	func (self DevmodeService) sendBlockAck(sender_id consensus.PeerId, block_id consensus.BlockId) {
		err := self.service.SendTo(sender_id, "ack", block_id.AsBytes())
		if err != nil {
			panic("Failed to send block ack")
		}
	}

	// Calculate the time to wait between publishing blocks. This will be a
	// random number between the settings sawtooth.consensus.min_wait_time and
	// sawtooth.consensus.max_wait_time if max > min, else DEFAULT_WAIT_TIME. If
	// there is an error parsing those settings, the time will be
	// DEFAULT_WAIT_TIME.
	func (self *DevmodeService) calculateWaitTime(chain_head_id consensus.BlockId) time.Duration {
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

			logger.Debugf("Min: ", min_wait_time, " -- Max: {:?}", max_wait_time)

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

type DevmodeEngineImpl struct {
	startupState      consensus.StartupState
	service           *DevmodeService
	chainHead         consensus.Block
	waitTime          time.Duration
	publishedAtHeight bool
	start             time.Time
}

//> impl ConsensusEngineImpl for DevmodeEngineImpl
	func (self *DevmodeEngineImpl) Version() string {
		return "0.1"
	}

	func (self *DevmodeEngineImpl) Name() string {
		return "Devmode"
	}

	func (self *DevmodeEngineImpl) Start(startupState consensus.StartupState, service consensus.ConsensusService, notifyChan chan consensus.Notification) error {
		self.service = NewDevmodeService(service)
		self.chainHead = startupState.ChainHead()
		self.waitTime = self.service.calculateWaitTime(self.chainHead.BlockId())
		self.publishedAtHeight = false
		self.start = time.Now()

		self.service.initializeBlock()

		for {
			select {
			case n := <-notifyChan:
				switch notification := n.(type) {
				case consensus.NotificationShutdown:
					break
				case consensus.NotificationBlockNew:
					self.HandleBlockNew(notification.Block)
				case consensus.NotificationBlockValid:
					self.HandleBlockValid(notification.BlockId)
				case consensus.NotificationBlockCommit:
					self.HandleBlockCommit(notification.BlockId)
				case consensus.NotificationPeerMessage:
					self.HandlePeerMessage(notification.PeerMessage, notification.SenderId)
				}
			case <-time.After(time.Millisecond * 10):
				if !self.publishedAtHeight && (time.Now().Sub(self.start) > self.waitTime) {
					logger.Info("Timer expired -- publishing block")
					newBlockId := self.service.finalizeBlock()
					self.publishedAtHeight = true
					self.service.broadcast_published_block(newBlockId)
				}
			}
		}

		return nil
	}

	func (self *DevmodeEngineImpl) HandlePeerMessage(peerMessage consensus.PeerMessage, senderId consensus.PeerId) {
		messageType := peerMessage.Header().MessageType()

		switch messageType {
		case "published":
			logger.Infof("Received block published message from ", senderId, ": ", hex.EncodeToString(peerMessage.Content()))
		case "received":
			logger.Infof("Received block received message from ", senderId, ": ", hex.EncodeToString(peerMessage.Content()))
			self.service.sendBlockAck(senderId, consensus.NewBlockIdFromBytes(peerMessage.Content()))
		case "ack":
			logger.Infof("Received ack message from ", senderId, ": ", hex.EncodeToString(peerMessage.Content()))
		default:
			panic("HandlePeerMessage() recieved an invalid message type")
		}
	}

	func (self *DevmodeEngineImpl) HandleBlockNew(block consensus.Block) {
		logger.Infof("Checking consensus data: %v", block)

		if block.PreviousId() == consensus.BLOCK_ID_NULL {
			logger.Warn("Received genesis block; ignoring")
			return
		}

		if checkConsensus(block) {
			logger.Infof("Passed consensus check: %v", block)
			self.service.checkBlock(block.BlockId())
		} else {
			logger.Infof("Failed consensus check: %v", block)
			self.service.failBlock(block.BlockId())
		}
	}
	func (self *DevmodeEngineImpl) HandleBlockValid(blockId consensus.BlockId) {
		block := self.service.getBlock(blockId)

		self.service.sendBlockReceived(block)

		chainHead := self.service.getChainHead()

		logger.Infof("Choosing between chain heads -- current: %v -- new: %v", chainHead, block)

		// blockBlockIdGreater = if block.BlockId() > chainHead.BlockId()
		blockBlockId := block.BlockId()
		chainHeadId := chainHead.BlockId()
		blockBlockIdGreater := bytes.Compare(blockBlockId.AsBytes(), chainHeadId.AsBytes()) > 0

		// advance the chain if possible
		if block.BlockNum() > chainHead.BlockNum() || (block.BlockNum() == chainHead.BlockNum() && blockBlockIdGreater) {
			logger.Infof("Committing %v", block)
			self.service.commitBlock(blockId)
		} else if block.BlockNum() < chainHead.BlockNum() {
			chainBlock := chainHead
			for {
				chainBlock = self.service.getBlock(chainBlock.PreviousId())
				if chainBlock.BlockNum() == block.BlockNum() {
					break
				}
			}

			// if block.BlockId() > chainBlock.BlockId()
			blockBlockId := block.BlockId()
			chainBlockId := chainBlock.BlockId()
			if bytes.Compare(blockBlockId.AsBytes(), chainBlockId.AsBytes()) > 0 {
				logger.Infof("Switching to new fork %v", block)
				self.service.commitBlock(blockId)
			} else {
				logger.Infof("Ignoring fork %v", block)
				self.service.ignoreBlock(blockId)
			}
		} else {
			logger.Infof("Ignoring %v", block)
			self.service.ignoreBlock(blockId)
		}
	}

	// The chain head was updated, so abandon the
	// block in progress and start a new one.
	func (self *DevmodeEngineImpl) HandleBlockCommit(newChainHead consensus.BlockId) {
		logger.Infof("Chain head updated to %v, abandoning block in progress", newChainHead)

		self.service.cancelBlock()

		self.waitTime = self.service.calculateWaitTime(newChainHead)
		self.publishedAtHeight = false
		self.start = time.Now()
		self.service.initializeBlock()
	}

//<

func checkConsensus(block consensus.Block) bool {
	return reflect.DeepEqual(block.Payload(), createConsensus(block.Summary()))
}

func createConsensus(summary []byte) []byte {
	// create a byte slice from the ascii values of a string
	consensusSlice := []byte("Devmode")

	// concatinate the two byte slices
	consensusSlice = append(consensusSlice, summary...)

	return consensusSlice
}
