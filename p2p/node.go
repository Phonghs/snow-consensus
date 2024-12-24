package p2p

import (
	"SnowConsensus/util"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Key struct {
	PublicKey  string
	PrivateKey string // = ""
	Algorithm  string
	Hash       string
}

type Profile struct {
	ID   string
	IP   string
	Port string
	Key  Key
}

type Transaction struct {
	ID                   string // consider the ID as the transaction data.
	Status               *string
	PreferenceStatus     *string
	ConsecutiveSuccesses int
	Mutex                sync.Mutex
}

type Env struct {
	TimeOut           int
	SampleSize        int
	QuorumSize        int
	DecisionThreshold int
}

type SnowNode struct {
	NodeProfile     *Profile
	PeerNodeInfoMap map[string]*Profile     // NodeId -> Node
	Transaction     map[string]*Transaction // TransactionID -> Transaction
	Env             Env
	Mutex           sync.Mutex
}

func (n *SnowNode) ReceiveQuery(m Message) (Message, error) {
	verified, err := n.VerifyMessage(m)
	if err != nil {
		return Message{}, err
	}
	if !verified {
		return Message{}, errors.New("signature not verified")
	}

	if !n.ValidTransaction(m) {
		return Message{}, errors.New("invalid transaction")
	}

	// TODO: Shouldn't save on local memory
	n.Mutex.Lock()
	_, ok := n.Transaction[m.TransactionData]
	if !ok {
		n.Transaction[m.TransactionData] = &Transaction{
			ID:                   m.TransactionData,
			Status:               nil,
			PreferenceStatus:     nil,
			ConsecutiveSuccesses: 0,
		}
	}
	n.Mutex.Unlock()

	if n.Transaction[m.TransactionData].Status == nil {
		n.ProcessSnowConsensus(m)
	}

	mess := Message{
		Id:              n.NodeProfile.ID,
		StatusCode:      200,
		TransactionData: m.TransactionData,
		Status:          n.Transaction[m.TransactionData].Status,
	}

	if n.SignMessage(&mess) != nil {
		log.Println(fmt.Println("Can't sign message", err))
		return Message{}, errors.New("failed to sign message")
	}

	return mess, nil
}

func (n *SnowNode) SendQuery(m Message, destinationID string) (Message, error) {
	peerNode, ok := n.PeerNodeInfoMap[destinationID]
	if !ok {
		return Message{}, errors.New("peer node not found")
	}
	url := fmt.Sprintf("http://%s:%s/transaction/query", peerNode.IP, peerNode.Port)
	method := "POST"
	byteData, err := json.Marshal(&m)
	if err != nil {
		log.Println(fmt.Println("Can't marshal message", err))
		return Message{}, err
	}
	res := Message{}

	code, err := CommunicateHTTP[Message](url, method, nil, byteData, &res, time.Duration(n.Env.TimeOut)*time.Millisecond)
	if err != nil {
		log.Println(fmt.Println(fmt.Sprintf("Can't communicate with peer node %s with error %s", peerNode.ID), err.Error()))
		return Message{}, err
	}

	if code != 200 {
		return Message{}, errors.New("failed to send message")
	}
	return res, nil
}

func (n *SnowNode) ValidTransaction(m Message) bool {
	return true
}

func (n *SnowNode) VerifyMessage(m Message) (bool, error) {
	peerNode, ok := n.PeerNodeInfoMap[m.Id]
	if !ok {
		return false, errors.New("peer node not found")
	}
	return util.VerifySignature(m.TransactionData, m.Signature, peerNode.Key.PublicKey, peerNode.Key.Algorithm, peerNode.Key.Hash)
}

func (n *SnowNode) SignMessage(m *Message) error {
	signature, err := util.SignMessage(m.TransactionData, n.NodeProfile.Key.PrivateKey, n.NodeProfile.Key.Algorithm, n.NodeProfile.Key.Hash)
	if err != nil {
		return err
	}
	m.Signature = signature
	return nil
}

func (n *SnowNode) SelectRandomValidator() []string {
	if n.PeerNodeInfoMap == nil || len(n.PeerNodeInfoMap) == 0 {
		return []string{}
	}
	if n.Env.SampleSize >= len(n.PeerNodeInfoMap) {
		return util.GetKeys(n.PeerNodeInfoMap)
	}
	k := n.Env.SampleSize
	result := make([]string, k)
	i := 0
	for id, _ := range n.PeerNodeInfoMap {
		if i < k {
			result[i] = id
		} else {
			j := rand.Intn(i + 1)
			if j < k {
				result[j] = id
			}
		}
		i++
	}

	return result
}

func (n *SnowNode) ProcessSnowConsensus(m Message) {
	// begin repeated random subsampling
	if n.Transaction[m.TransactionData].ConsecutiveSuccesses < n.Env.DecisionThreshold {
		j := 1
		for {
			log.Println(fmt.Printf("Node %s: Repeated random subsampling for transaction %s in loop %v\n", n.NodeProfile.ID, m.TransactionData, j))
			j++
			if n.Transaction[m.TransactionData].ConsecutiveSuccesses >= n.Env.DecisionThreshold {
				break
			}

			validators := n.SelectRandomValidatorV2(m.Id)
			// query K validators
			if len(validators) == 0 {
				log.Println("no validators found")
				break
			}
			var wg sync.WaitGroup
			var result = make([]Message, len(validators))
			signature, err := util.SignMessage(m.TransactionData, n.NodeProfile.Key.PrivateKey, n.NodeProfile.Key.Algorithm, n.NodeProfile.Key.Hash)
			if err != nil {
				log.Println(err)
				break
			}
			for i, validatorID := range validators {
				wg.Add(1)
				go func(i int, validatorId string) {
					defer wg.Done()
					message, err := n.SendQuery(Message{
						Id:              n.NodeProfile.ID,
						TransactionData: m.TransactionData,
						Status:          nil,
						Signature:       signature,
					}, validatorId)
					if err != nil {
						result[i] = Message{}
					}
					result[i] = message
				}(i, validatorID)
			}
			wg.Wait()

			ans, countAns := n.countAnswer(result)
			if n.Transaction[m.TransactionData].Status != nil {
				break
			}
			n.Transaction[m.TransactionData].Mutex.Lock()
			if countAns >= n.Env.QuorumSize {
				n.Transaction[m.TransactionData].PreferenceStatus = &ans
				if *n.Transaction[m.TransactionData].PreferenceStatus == ans {
					n.Transaction[m.TransactionData].ConsecutiveSuccesses++
				} else {
					n.Transaction[m.TransactionData].ConsecutiveSuccesses = 1
				}
			} else if countAns < n.Env.QuorumSize {
				n.Transaction[m.TransactionData].PreferenceStatus = nil
				n.Transaction[m.TransactionData].ConsecutiveSuccesses = 0
			}
			if n.Transaction[m.TransactionData].ConsecutiveSuccesses > n.Env.DecisionThreshold {
				log.Println(fmt.Printf("Node %s: Decision reached for transaction %s\n", n.NodeProfile.ID, m.TransactionData))
				n.Transaction[m.TransactionData].Status = n.Transaction[m.TransactionData].PreferenceStatus
				n.Transaction[m.TransactionData].Mutex.Unlock()
				break
			}
			n.Transaction[m.TransactionData].Mutex.Unlock()
		}
	}
}

func (n *SnowNode) countAnswer(input []Message) (string, int) {
	maxAns := 0
	ans := ""
	countAnsMap := make(map[string]int)
	for _, message := range input {
		if message.StatusCode == 200 {
			if ans, ok := countAnsMap[message.TransactionData]; ok {
				countAnsMap[message.TransactionData] = ans + 1
			} else {
				countAnsMap[message.TransactionData] = 1
			}
			if countAnsMap[message.TransactionData] > maxAns {
				maxAns = countAnsMap[message.TransactionData]
				ans = message.TransactionData
			}
		}
	}

	return ans, maxAns
}

func (n *SnowNode) CreateTransaction(tranData string) (string, error) {
	// TODO: Shouldn't save on local memory
	n.Mutex.Lock()
	_, ok := n.Transaction[tranData]
	if !ok {
		status := "yes"
		n.Transaction[tranData] = &Transaction{
			ID:                   tranData,
			Status:               &status,
			PreferenceStatus:     nil,
			ConsecutiveSuccesses: 0,
		}
	}
	n.Mutex.Unlock()

	n.ProcessSnowConsensus(Message{
		Id:              n.NodeProfile.ID,
		TransactionData: tranData,
	})

	if n.Transaction[tranData].Status != nil {
		return *n.Transaction[tranData].Status, nil
	} else {
		return "", errors.New("failed to create transaction")
	}

}

func (n *SnowNode) GetID() string {
	return n.NodeProfile.ID
}

func (n *SnowNode) SelectRandomValidatorV2(excludeId string) []string {

	if n.PeerNodeInfoMap == nil || len(n.PeerNodeInfoMap) == 0 {
		return []string{}
	}
	if n.Env.SampleSize >= len(n.PeerNodeInfoMap) {
		return util.GetKeys(n.PeerNodeInfoMap)
	}
	k := n.Env.SampleSize
	result := make([]string, k)

	i := 0

	for id, _ := range n.PeerNodeInfoMap {
		if id == excludeId {
			continue
		}
		if i < k {
			result[i] = id
		} else {
			j := rand.Intn(i + 1)
			if j < k {
				result[j] = id
			}
		}
		i++
	}
	return result
}
