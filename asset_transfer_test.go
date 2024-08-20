package main

import (
    "encoding/json"
    "testing"

    "github.com/hyperledger/fabric-contract-api-go/contractapi"
    "github.com/hyperledger/fabric-chaincode-go/shim"
    "github.com/stretchr/testify/assert"
)

// MockStub is a mock implementation of the ChaincodeStubInterface for testing
type MockStub struct {
    shim.ChaincodeStubInterface
    state map[string][]byte
}

func (m *MockStub) PutState(key string, value []byte) error {
    m.state[key] = value
    return nil
}

func (m *MockStub) GetState(key string) ([]byte, error) {
    return m.state[key], nil
}

func (m *MockStub) DelState(key string) error {
    delete(m.state, key)
    return nil
}

func (m *MockStub) CreateCompositeKey(objectType string, attributes []string) (string, error) {
    return objectType + "_" + attributes[0], nil
}

// MockTransactionContext provides a mock context for testing
type MockTransactionContext struct {
    contractapi.TransactionContextInterface
    stub shim.ChaincodeStubInterface
}

func (m *MockTransactionContext) GetStub() shim.ChaincodeStubInterface {
    return m.stub
}

func TestCreateAsset(t *testing.T) {
    smartContract := new(SmartContract)
    stub := &MockStub{
        state: make(map[string][]byte),
    }
    ctx := &MockTransactionContext{
        stub: stub,
    }

    dealerID := "dealer1"
    msisdn := "1234567890"
    mpin := "0000"
    balance := 100
    status := "active"
    transAmount := 50
    transType := "credit"
    remarks := "Test asset"

    err := smartContract.CreateAsset(ctx, dealerID, msisdn, mpin, balance, status, transAmount, transType, remarks)
    assert.Nil(t, err)

    assetBytes, err := stub.GetState(msisdn)
    assert.Nil(t, err)

    var asset Asset
    err = json.Unmarshal(assetBytes, &asset)
    assert.Nil(t, err)
    assert.Equal(t, dealerID, asset.DealerID)
    assert.Equal(t, balance, asset.Balance)
    assert.Equal(t, status, asset.Status)
    assert.Equal(t, transAmount, asset.TransAmount)
    assert.Equal(t, transType, asset.TransType)
    assert.Equal(t, remarks, asset.Remarks)
}