package main

import (
    "encoding/json"
    "fmt"
    
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Asset struct {
    DealerID    string `json:"dealerID"`
    MSISDN      string `json:"msisdn"`
    MPIN        string `json:"mpin"`
    Balance     int    `json:"balance"`
    Status      string `json:"status"`
    TransAmount int    `json:"transAmount"`
    TransType   string `json:"transType"`
    Remarks     string `json:"remarks"`
}
type SmartContract struct {
    contractapi.Contract
}

// CreateAsset - Create a new asset
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, dealerID string, msisdn string, mpin string, balance int, status string, transAmount int, transType string, remarks string) error {
    asset := Asset{
        DealerID:    dealerID,
        MSISDN:      msisdn,
        MPIN:        mpin,
        Balance:     balance,
        Status:      status,
        TransAmount: transAmount,
        TransType:   transType,
        Remarks:     remarks,
    }

    assetJSON, err := json.Marshal(asset)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(msisdn, assetJSON)
}

// ReadAsset - Read an asset
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, msisdn string) (*Asset, error) {
    assetJSON, err := ctx.GetStub().GetState(msisdn)
    if err != nil {
        return nil, err
    }
    if assetJSON == nil {
        return nil, fmt.Errorf("Asset %s does not exist", msisdn)
    }

    var asset Asset
    err = json.Unmarshal(assetJSON, &asset)
    if err != nil {
        return nil, err
    }

    return &asset, nil
}

// UpdateAsset - Update an existing asset
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, dealerID string, msisdn string, mpin string, balance int, status string, transAmount int, transType string, remarks string) error {
    asset, err := s.ReadAsset(ctx, msisdn)
    if err != nil {
        return err
    }

    asset.DealerID = dealerID
    asset.MPIN = mpin
    asset.Balance = balance
    asset.Status = status
    asset.TransAmount = transAmount
    asset.TransType = transType
    asset.Remarks = remarks

    assetJSON, err := json.Marshal(asset)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(msisdn, assetJSON)
}

// GetAllAssets - Query all assets
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
    queryString := fmt.Sprintf("{\"selector\":{}}")

    resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
    if err != nil {
        return nil, err
    }
    defer resultsIterator.Close()

    var assets []*Asset
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }

        var asset Asset
        err = json.Unmarshal(queryResponse.Value, &asset)
        if err != nil {
            return nil, err
        }
        assets = append(assets, &asset)
    }

    return assets, nil
}
