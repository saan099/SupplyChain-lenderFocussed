package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type Supplychaincode struct {
}

var supplierIdIndices string = "SupplierIdIndices"
var bankIdIndices string = "BankIdIndices"
var bankInvoicesKey string = "BankInvoices"

func main() {

	err := shim.Start(new(Supplychaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

func (t *Supplychaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, errors.New("error:A01 wrong number of aguments in initialization")
	}
	var suppliers []string
	supplierIdIndicesAsbytes, _ := json.Marshal(suppliers)
	err := stub.PutState(supplierIdIndices, supplierIdIndicesAsbytes)
	if err != nil {
		return nil, err
	}
	var banks []string
	bankIdAsbytes, _ := json.Marshal(banks)
	err = stub.PutState(bankIdIndices, bankIdAsbytes)
	if err != nil {
		return nil, err
	}

	var bankInvoices []invoice
	bankInvoicesAsbytes, _ := json.Marshal(bankInvoices)
	err = stub.PutState(bankInvoicesKey, bankInvoicesAsbytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

//Invoking functionality
func (t *Supplychaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "init" {
		return t.Init(stub, function, args)
	} else if function == "initBuyer" {
		return t.InitBuyer(stub, args)
	} else if function == "initSupplier" {
		return t.InitSupplier(stub, args)
	} else if function == "createPO" {
		return t.CreatePO(stub, args)
	} else if function == "updatePOStatus" {
		return t.UpdatePOStatus(stub, args)
	} else if function == "generateInvoice" {
		return t.GenerateInvoice(stub, args)
	} else if function == "updateInvoiceStatus" {
		return t.UpdateInvoiceStatus(stub, args)
	} else if function == "initBank" {
		return t.InitBank(stub, args)
	} else if function == "disburseInvoice" {
		return t.DisburseInvoice(stub, args)
	} else if function == "makeOffer" {
		return t.MakeOffer(stub, args)
	} else if function == "updateOfferStatus" {
		return t.UpdateOfferStatus(stub, args)
	} else if function == "markRepayment" {
		return t.MarkRepayment(stub, args)
	}

	return nil, errors.New("error:C01 No function called")

}

// Query data
func (t *Supplychaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "read" {
		return t.Read(stub, args)
	} else if function == "readAllSuppliers" {
		return t.ReadAllSuppliers(stub, args)
	}

	return nil, errors.New("error:C02 No function called")
}

func (t *Supplychaincode) Read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("error:A04 Wrong numer of arguments")
	}

	valAsbytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, err
	}
	return valAsbytes, nil

}

func (t *Supplychaincode) ReadAllSuppliers(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 0 {
		return nil, errors.New("error:A05 wrong number of arguments")
	}

	var idList []string
	idListAsbytes, _ := stub.GetState(supplierIdIndices)
	err := json.Unmarshal(idListAsbytes, &idList)
	var suppliers []supplier
	for i := range idList {
		supplierAcc := supplier{}
		supplierAsbytes, _ := stub.GetState(idList[i])
		err := json.Unmarshal(supplierAsbytes, &supplierAcc)
		if err != nil {
			return nil, errors.New("error:U02 unmarshaliing error")
		}
		suppliers = append(suppliers, supplierAcc)
	}
	supplierData, err := json.Marshal(suppliers)
	if err != nil {
		return nil, errors.New("error:M02 marshalling error")
	}

	return supplierData, nil
}
