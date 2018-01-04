package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func (t *Supplychaincode) InitBank(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) < 2 {
		return nil, errors.New("wrong number of arguments")
	}
	bankAcc := bank{}
	bankAcc.BankId = args[0]
	bankAcc.BankName = args[1]
	bankAcc.EmployeeName = args[2]
	bankAcc.EmailId = args[3]
	var Dinvoices []disbursementInvoice
	bankAcc.DInvoices = Dinvoices
	bankAsbytes, err := json.Marshal(bankAcc)
	if err != nil {
		return nil, err
	}
	erro := stub.PutState(args[0], bankAsbytes)
	if erro != nil {
		return nil, err
	}
	var bankIds []string
	bankIndicesAsbytes, _ := stub.GetState(bankIdIndices)
	jsonerr := json.Unmarshal(bankIndicesAsbytes, &bankIds)
	if jsonerr != nil {
		return nil, jsonerr
	}
	bankIds = append(bankIds, args[0])
	newBankIdIndices, _ := json.Marshal(bankIds)
	err = stub.PutState(bankIdIndices, newBankIdIndices)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *Supplychaincode) MakeOffer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) < 2 {
		return nil, errors.New("wrong number of arguments")
	}
	invoiceId := args[0]
	bankId := args[1]
	days := args[2]
	rate := args[3]
	offerId := args[4]
	var supplierId string
	var invoices []invoice
	invoicesAsbytes, _ := stub.GetState(bankInvoicesKey)
	_ = json.Unmarshal(invoicesAsbytes, &invoices)
	bankAcc := bank{}
	bankAsbytes, _ := stub.GetState(bankId)
	_ = json.Unmarshal(bankAsbytes, &bankAcc)
	var disInv disbursementInvoice
	currentTime := time.Now().Local()
	for i := range invoices {
		if invoices[i].InvoiceId == invoiceId {
			disInv.Bank = bankId
			disInv.Date = currentTime.Format("02-01-2006")
			disInv.Time = currentTime.Format("3:04PM")
			disInv.Days, _ = strconv.Atoi(days)
			disInv.Details = invoices[i]
			disInv.DisRate, _ = strconv.ParseFloat(rate, 64)
			disInv.DisAmount = disInv.DisRate * disInv.Details.Total / 100
			disInv.DisbursedAmount = disInv.Details.Total - disInv.DisAmount
			supplierId = invoices[i].Supplier
		}
	}
	off := offer{}
	off.Details = disInv
	off.Status = "pending"
	off.Id = offerId

	bankAcc.Offers = append(bankAcc.Offers, off)
	newBankAsbytes, _ := json.Marshal(bankAcc)
	_ = stub.PutState(bankId, newBankAsbytes)

	supplierAcc := supplier{}
	supplierAsbytes, _ := stub.GetState(supplierId)
	_ = json.Unmarshal(supplierAsbytes, &supplierAcc)
	supplierAcc.Offers = append(supplierAcc.Offers, off)
	newSupplierAsbytes, _ := json.Marshal(supplierAcc)
	_ = stub.PutState(supplierId, newSupplierAsbytes)

	return nil, nil

}

func (t *Supplychaincode) DisburseInvoice(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 6 {
		return nil, errors.New("wrong number of arguments")
	}
	invoiceId := args[0]
	bankId := args[1]
	days := args[2]
	rate := args[3]
	penaltyPerDay, e := strconv.ParseFloat(args[4], 64)
	if e != nil {
		return nil, errors.New("C:expecting float penalty per day")
	}
	recourse := args[5]
	bankAcc := bank{}
	bankAsbytes, _ := stub.GetState(bankId)
	err := json.Unmarshal(bankAsbytes, &bankAcc)
	if err != nil {
		return nil, err
	}
	var Dinv = disbursementInvoice{}
	var invoices []invoice
	var supplierId string
	var buyerId string
	invoicesAsbytes, _ := stub.GetState(bankInvoicesKey)
	_ = json.Unmarshal(invoicesAsbytes, &invoices)
	currentTime := time.Now().Local()
	for i := range invoices {
		if invoices[i].InvoiceId == invoiceId {
			if invoices[i].InvoiceId == `disbursed` {
				return nil, nil
			}
			invoices[i].Status = `disbursed`
			Dinv.Details = invoices[i]
			Dinv.Bank = bankId
			Dinv.Date = currentTime.Format("02-01-2006")
			Dinv.Time = currentTime.Format("3:04PM")
			Dinv.Days, _ = strconv.Atoi(days)
			Dinv.DisRate, _ = strconv.ParseFloat(rate, 64)
			Dinv.DisAmount = Dinv.DisRate * Dinv.Details.Total / 100
			Dinv.DisbursedAmount = Dinv.Details.Total - Dinv.DisAmount
			Dinv.PenaltyPerDay = penaltyPerDay
			Dinv.Recourse = recourse
			bankAcc.DInvoices = append(bankAcc.DInvoices, Dinv)
			supplierId = invoices[i].Supplier
			buyerId = invoices[i].Buyer
		}
	}
	newBankAsbytes, _ := json.Marshal(bankAcc)
	_ = stub.PutState(bankId, newBankAsbytes)

	newBankInvoicesAsbytes, _ := json.Marshal(invoices)
	_ = stub.PutState(bankInvoicesKey, newBankInvoicesAsbytes)

	var supplierAcc = supplier{}
	supplierAsbytes, _ := stub.GetState(supplierId)
	_ = json.Unmarshal(supplierAsbytes, &supplierAcc)
	for i := range supplierAcc.Invoices {
		if supplierAcc.Invoices[i].InvoiceId == invoiceId {
			supplierAcc.Invoices[i].Status = `disbursed`
		}
	}
	supplierAcc.DInvoices = append(supplierAcc.DInvoices, Dinv)
	newSupplierAsbytes, _ := json.Marshal(supplierAcc)
	_ = stub.PutState(supplierId, newSupplierAsbytes)

	var buyerAcc = buyer{}
	buyerAsbytes, _ := stub.GetState(buyerId)
	_ = json.Unmarshal(buyerAsbytes, &buyerAcc)

	for i := range buyerAcc.Invoices {
		if buyerAcc.Invoices[i].InvoiceId == invoiceId {
			buyerAcc.Invoices[i].Status = `disbursed`
		}
	}
	buyerAcc.DInvoices = append(buyerAcc.DInvoices, Dinv)
	newBuyerAcc, _ := json.Marshal(buyerAcc)
	_ = stub.PutState(buyerId, newBuyerAcc)

	return nil, nil
}

func (t *Supplychaincode) MarkRepayment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("wrong number of arguments")
	}
	bankId := args[0]
	invId := args[1]

	bankAcc := bank{}
	bankAsbytes, _ := stub.GetState(bankId)
	err := json.Unmarshal(bankAsbytes, &bankAcc)
	if err != nil {
		return nil, errors.New("C: couldnt unmarshal")
	}
	var supplierId, buyerId string
	for i := range bankAcc.DInvoices {
		if bankAcc.DInvoices[i].Details.InvoiceId == invId {
			bankAcc.DInvoices[i].Details.Status = `repayed`
			buyerId = bankAcc.DInvoices[i].Details.Buyer
			supplierId = bankAcc.DInvoices[i].Details.Supplier
		}
	}
	newbankAsbytes, _ := json.Marshal(bankAcc)
	_ = stub.PutState(bankId, newbankAsbytes)

	buyerAcc := buyer{}
	buyerAsbytes, _ := stub.GetState(buyerId)
	err = json.Unmarshal(buyerAsbytes, &buyerAcc)
	if err != nil {
		return nil, errors.New("C: couldnt unmarshal")
	}
	for j := range buyerAcc.DInvoices {
		if buyerAcc.DInvoices[j].Details.InvoiceId == invId {
			buyerAcc.DInvoices[j].Details.Status = `repayed`
		}
	}
	newBuyerAsbytes, _ := json.Marshal(buyerAcc)
	_ = stub.PutState(buyerId, newBuyerAsbytes)

	supplierAcc := supplier{}
	supplierasbytes, _ := stub.GetState(supplierId)
	err = json.Unmarshal(supplierasbytes, &supplierAcc)
	if err != nil {
		return nil, errors.New("C: couldnt unmarshal")
	}
	for k := range supplierAcc.DInvoices {
		if supplierAcc.DInvoices[k].Details.InvoiceId == invId {
			supplierAcc.DInvoices[k].Details.Status = `repayed`
		}
	}
	newSupplierAsbytes, _ := json.Marshal(supplierAcc)
	_ = stub.PutState(supplierId, newSupplierAsbytes)

	return nil, nil
}

func (t *Supplychaincode) readAllBankers(stub shim.ChaincodeStubInterface, args []string) []bank {
	var banks []bank
	var bankIds []string
	bankIdsAsbytes, _ := stub.GetState(bankIdIndices)
	_ = json.Unmarshal(bankIdsAsbytes, &bankIds)
	for i := 0; i < len(bankIds); i++ {
		var b = bank{}
		accBytes, _ := stub.GetState(bankIds[i])
		_ = json.Unmarshal(accBytes, &b)
		banks = append(banks, b)
	}
	var data, _ = json.Marshal(banks)
	_ = stub.PutState("banks", data)
	return banks

}
