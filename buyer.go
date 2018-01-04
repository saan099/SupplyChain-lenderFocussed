package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//creating account for buyer
func (t *Supplychaincode) InitBuyer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) < 2 {
		return nil, errors.New("error:A02 wrong number of arguments")
	}
	var buyerAcc = buyer{}
	buyerAcc.BuyerId = args[0]
	buyerAcc.BuyerName = args[1]
	buyerAcc.EmployeeName = args[2]
	buyerAcc.EmailId = args[3]
	var orders []purchaseOrder
	var invoices []invoice
	var Dinvoices []disbursementInvoice
	buyerAcc.DInvoices = Dinvoices
	buyerAcc.Invoices = invoices
	buyerAcc.PurchaseOrders = orders
	buyerAsbytes, err := json.Marshal(buyerAcc)
	err = stub.PutState(args[0], buyerAsbytes)
	if err != nil {
		return nil, errors.New("error:S02 account initialization failed")
	}

	return nil, nil
}

func (t *Supplychaincode) CreatePO(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) < 0 {
		return nil, errors.New("error:A05 wrong number of arguments")
	}
	order := purchaseOrder{}

	purchaseId := args[0]
	numOfProducts, _ := strconv.Atoi(args[1])
	currentTime := time.Now().Local()
	order.Date = currentTime.Format("02-01-2006")
	order.Time = currentTime.Format("3:04PM")
	order.OrderId = purchaseId
	var i int = 2
	var total float64
	var products []product
	for i = 2; i < numOfProducts+2; i++ {
		pro := product{}
		err := json.Unmarshal([]byte(args[i]), &pro)
		if err != nil {
			return nil, errors.New("error:U05 unamrshalling error")
		}

		total = total + pro.Value
		products = append(products, pro)

	}
	order.Products = products

	order.Supplier = args[i]

	order.TotalValue = total
	i = i + 1
	order.CreditPeriod, _ = strconv.Atoi(args[i])
	i = i + 1
	order.Buyer = args[i]
	order.Status = "processing"

	buyerAsbytes, _ := stub.GetState(order.Buyer)
	supplierAsbytes, _ := stub.GetState(order.Supplier)
	buyerAcc := buyer{}
	supplierAcc := supplier{}
	err := json.Unmarshal(buyerAsbytes, &buyerAcc)
	if err != nil {
		return nil, errors.New("error:U06 unmarshalling error")
	}
	err = json.Unmarshal(supplierAsbytes, &supplierAcc)
	if err != nil {
		return nil, errors.New("error:U07 unmarshalling error")
	}
	buyerAcc.PurchaseOrders = append(buyerAcc.PurchaseOrders, order)
	supplierAcc.PurchaseOrders = append(supplierAcc.PurchaseOrders, order)

	buyerAsNewbytes, _ := json.Marshal(buyerAcc)
	if err != nil {
		return nil, errors.New("error:M04 unmarshalling error")
	}
	supplierAsNewbytes, _ := json.Marshal(supplierAcc)
	if err != nil {
		return nil, errors.New("error:M05 unmarshalling error")
	}
	err = stub.PutState(order.Buyer, buyerAsNewbytes)
	err = stub.PutState(order.Supplier, supplierAsNewbytes)

	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *Supplychaincode) UpdateInvoiceStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("error:A08 wrong number of arguments")
	}
	var supplierId string
	buyerId := args[0]
	invoiceId := args[1]
	status := args[2]
	buyerAcc := buyer{}
	inv := invoice{}
	supplierAcc := supplier{}
	buyerAsbytes, err := stub.GetState(buyerId)
	if err != nil {
		return nil, errors.New("error:S0 failed")
	}
	err = json.Unmarshal(buyerAsbytes, &buyerAcc)

	for i := range buyerAcc.Invoices {
		if invoiceId == buyerAcc.Invoices[i].InvoiceId {
			supplierId = buyerAcc.Invoices[i].Supplier
			buyerAcc.Invoices[i].Status = status
			inv = buyerAcc.Invoices[i]
		}
	}
	newBuyerAsbytes, err := json.Marshal(buyerAcc)
	if err != nil {
		return nil, errors.New("error:M marshalling error")
	}
	err = stub.PutState(buyerId, newBuyerAsbytes)

	supplierAsbytes, err := stub.GetState(supplierId)
	err = json.Unmarshal(supplierAsbytes, &supplierAcc)

	for i := range supplierAcc.Invoices {
		if invoiceId == supplierAcc.Invoices[i].InvoiceId {
			supplierAcc.Invoices[i].Status = status
		}
	}
	newSupplierAsbytes, err := json.Marshal(supplierAcc)
	if err != nil {
		return nil, errors.New("error:M marshalling error")
	}
	err = stub.PutState(supplierId, newSupplierAsbytes)
	if status == "approved" {
		var bankInvoices []invoice
		bankInvoicesAsbytes, _ := stub.GetState(bankInvoicesKey)
		_ = json.Unmarshal(bankInvoicesAsbytes, &bankInvoices)
		bankInvoices = append(bankInvoices, inv)
		newBankInvoicesAsbytes, _ := json.Marshal(bankInvoices)
		_ = stub.PutState(bankInvoicesKey, newBankInvoicesAsbytes)
	}
	/*var banks []bank
	banks = t.readAllBankers(stub, args)
	for i := range bank {
		banks[i].Invoices = append(banks[i].Invoices, inv)
	}*/
	return nil, nil
}
