package ethereum

import (
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/accounts/abi"
    "strings"
    "context"
    "ubex-crowdsale-api/common/models"
    "fmt"
    "github.com/ethereum/go-ethereum"
    "math/big"
)

type Contract struct {
    Wallet  *Wallet
    Address common.Address
    Abi     abi.ABI
    Events  map[string]string
}

func NewContract(addr string) *Contract {
    return &Contract{
        Wallet:  GetWallet(),
        Address: common.HexToAddress(addr),
    }
}

func (c *Contract) InitEvents(contractAbi string) (err error) {
    c.Abi, err = abi.JSON(strings.NewReader(contractAbi))
    if err != nil {
        return
    }

    c.Events = make(map[string]string)
    for _, event := range c.Abi.Events {
        c.Events[event.Id().String()] = event.Name
    }

    return
}

func (c *Contract) GetEventsByTopics(topics [][]common.Hash, fromBlock *big.Int) ([]models.ContractEvent, error) {
    query := ethereum.FilterQuery{
        Addresses: []common.Address{
            c.Address,
        },
        Topics: topics,
        FromBlock: fromBlock,
    }

    res, err := c.Wallet.Connection.FilterLogs(context.TODO(), query)
    if err != nil {
        return []models.ContractEvent{}, err
    }

    events := make([]models.ContractEvent, len(res))
    for key, eventLog := range res {
        name, ok := c.Events[eventLog.Topics[0].String()]
        if !ok {
            return events, fmt.Errorf("unknown event parsed: %s", eventLog.Topics[0].String())
        }

        if len(eventLog.Data)%32 != 0 {
            return events, fmt.Errorf("wrong event unindexed data provided: %d", len(eventLog.Data))
        }

        numOfIndexed := len(eventLog.Topics) - 1
        numOfUnindexed := len(eventLog.Data) / 32
        rawArgs := make([][]byte, numOfIndexed+numOfUnindexed)
        if numOfIndexed > 0 {
            for i, indData := range eventLog.Topics[1:] {
                rawArgs[i] = indData.Bytes()
            }
        }
        for i := 0; i < numOfUnindexed; i++ {
            rawArgs[numOfIndexed+i] = eventLog.Data[i*32:(i+1)*32]
        }

        block, err := c.Wallet.Connection.HeaderByHash(context.TODO(), eventLog.BlockHash)
        if err != nil {
            return events, err
        }

        events[key] = models.ContractEvent{
            Name:        name,
            RawArgs:     rawArgs,
            TxHash:      eventLog.TxHash,
            BlockNumber: eventLog.BlockNumber,
            BlockHash:   eventLog.BlockHash,
            BlockTime:   block.Time.String(),
            TxIndex:     eventLog.TxIndex,
            Removed:     eventLog.Removed,
        }
    }

    return events, nil
}
