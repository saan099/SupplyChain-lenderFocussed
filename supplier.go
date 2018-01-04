package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//initialization supplier
func (t *Supplychaincode) InitSupplier(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) < 2 {
		return nil, errors.New("error:A03 wrong number of arguments")
	}
	var supplierAcc = supplier{}
	supplierAcc.SupplierId = args[0]
	supplierAcc.SupplierName = args[1]
	supplierAcc.EmployeeName = args[2]
	supplierAcc.EmailId = args[3]
	var orders []purchaseOrder
	var invoices []invoice
	supplierAcc.Invoices = invoices
	supplierAcc.PurchaseOrders = orders
	supplierAsbytes, err := json.Marshal(supplierAcc)
	err = stub.PutState(args[0], supplierAsbytes)
	if err != nil {
		return nil, errors.New("error:S03 account initialization failed")
	}
	var supplierList []string
	listAsbytes, _ := stub.GetState(supplierIdIndices)
	err = json.Unmarshal(listAsbytes, &supplierList)
	if err != nil {
		return nil, errors.New("error:U01 Unmarshaling error")
	}
	supplierList = append(supplierList, args[0])
	newListAsbytes, err := json.Marshal(supplierList)
	err = stub.PutState(supplierIdIndices, newListAsbytes)
	if err != nil {
		return nil, errors.New("error:S06 failed")
	}
	return supplierAsbytes, nil
}

func (t *Supplychaincode) UpdatePOStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 3 {
		return nil, errors.New("error:A06 wrong number of arguments")
	}
	POid := args[0]
	supplierId := args[1]
	status := args[2]
	var buyerId string
	supplierAcc := supplier{}
	supplierAsbytes, err := stub.GetState(supplierId)
	if err != nil {
		return nil, errors.New("error:S07 failed")
	}
	err = json.Unmarshal(supplierAsbytes, &supplierAcc)
	if err != nil {
		return nil, errors.New("error:U03 unmashalling error")
	}
	for i := range supplierAcc.PurchaseOrders {
		if supplierAcc.PurchaseOrders[i].OrderId == POid {
			supplierAcc.PurchaseOrders[i].Status = status
			buyerId = supplierAcc.PurchaseOrders[i].Buyer

		}

	}
	supplierAsNewbytes, err := json.Marshal(supplierAcc)
	if err != nil {
		return nil, errors.New("error:M03 marshalling error")
	}
	err = stub.PutState(supplierId, supplierAsNewbytes)

	buyerAcc := buyer{}
	buyerAsBytes, err := stub.GetState(buyerId)
	if err != nil {
		return nil, errors.New("error:S08 wrong buyer ID")
	}
	err = json.Unmarshal(buyerAsBytes, &buyerAcc)
	if err != nil {
		return nil, errors.New("error:U04 unmshalling error")
	}
	for i := range buyerAcc.PurchaseOrders {
		if buyerAcc.PurchaseOrders[i].OrderId == POid {
			buyerAcc.PurchaseOrders[i].Status = status
		}
	}
	buyerAsNewbytes, err := json.Marshal(buyerAcc)
	if err != nil {
		return nil, errors.New("error:M03 marshalling error")
	}
	err = stub.PutState(buyerId, buyerAsNewbytes)

	return nil, nil
}

func (t *Supplychaincode) GenerateInvoice(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) < 4 {
		return nil, errors.New("error:A07 wrong number of arguments")
	}

	var err error
	invoiceId := args[0]

	inv := invoice{}

	order := purchaseOrder{}

	purchaseId := args[1]
	numOfProducts, _ := strconv.Atoi(args[2])
	currentTime := time.Now().Local()
	order.Date = currentTime.Format("02-01-2006")
	order.OrderId = purchaseId
	var i int = 3
	var total float64
	var products []product
	for i = 3; i < numOfProducts+3; i++ {
		pro := product{}
		err := json.Unmarshal([]byte(args[i]), &pro)
		if err != nil {
			return nil, errors.New("error:U04 error unmarshalling")
		}

		total = total + pro.Value
		products = append(products, pro)

	}
	order.Products = products

	order.Supplier = args[i]
	inv.Supplier = args[i]
	order.TotalValue = total
	i = i + 1
	order.CreditPeriod, _ = strconv.Atoi(args[i])
	i = i + 1
	order.Buyer = args[i]
	inv.Buyer = args[i]
	order.Status = "processing"

	i = i + 1
	inv.Taxes, err = strconv.ParseFloat(args[i], 64)
	inv.Subtotal = total
	inv.Total = total + inv.Subtotal*inv.Taxes/100
	inv.PurchaseOrders = append(inv.PurchaseOrders, order)
	inv.InvoiceId = invoiceId
	inv.Status = "in process"
	currentT := time.Now().Local()
	inv.Date = currentT.Format("02-01-2006")
	inv.Time = currentTime.Format("3:04PM")

	buyerAcc := buyer{}
	buyerAsbytes, err := stub.GetState(inv.Buyer)
	err = json.Unmarshal(buyerAsbytes, &buyerAcc)
	buyerAcc.Invoices = append(buyerAcc.Invoices, inv)
	newBuyerAsbytes, err := json.Marshal(buyerAcc)
	err = stub.PutState(inv.Buyer, newBuyerAsbytes)

	supplierAcc := supplier{}
	supplierAsbytes, err := stub.GetState(inv.Supplier)
	err = json.Unmarshal(supplierAsbytes, &supplierAcc)
	supplierAcc.Invoices = append(supplierAcc.Invoices, inv)
	newSupplierAsbytes, err := json.Marshal(supplierAcc)
	err = stub.PutState(inv.Supplier, newSupplierAsbytes)

	if err != nil {
		return nil, errors.New("error")
	}
	return nil, nil
}

func (t *Supplychaincode) UpdateOfferStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 3 {
		return nil, nil
	}
	supplierId := args[0]
	offerId := args[1]
	status := args[2]
	var bankId string
	var supplierAcc supplier
	supplierAsbytes, _ := stub.GetState(supplierId)
	_ = json.Unmarshal(supplierAsbytes, &supplierAcc)
	var invoiceId string
	for i := range supplierAcc.Offers {
		if supplierAcc.Offers[i].Id == offerId {
			supplierAcc.Offers[i].Status = status
			bankId = supplierAcc.Offers[i].Details.Bank
			invoiceId = supplierAcc.Offers[i].Details.Details.InvoiceId
		}
	}
	if status == `approved` {
		var bankIds []string
		for z := range supplierAcc.Offers {
			if supplierAcc.Offers[z].Details.Details.InvoiceId == invoiceId && supplierAcc.Offers[z].Id != offerId {
				supplierAcc.Offers[z].Status = `rejected`
				bankIds = append(bankIds, supplierAcc.Offers[z].Details.Bank)
			}
		}

		for x := range bankIds {
			localBankAcc := bank{}
			localBankAsbytes, _ := stub.GetState(bankIds[x])
			e := json.Unmarshal(localBankAsbytes, &localBankAcc)
			if e != nil {
				return nil, errors.New("C: couldnt unmarshal banks")
			}
			for c := range localBankAcc.Offers {
				if localBankAcc.Offers[c].Details.Details.InvoiceId == invoiceId {
					localBankAcc.Offers[c].Status = `rejected`
				}
			}
			newLocalBankAsbytes, _ := json.Marshal(localBankAcc)
			_ = stub.PutState(bankIds[x], newLocalBankAsbytes)
		}
	}
	newSupplierAsbytes, _ := json.Marshal(supplierAcc)
	_ = stub.PutState(supplierId, newSupplierAsbytes)

	var bankAcc bank
	bankAsbytes, _ := stub.GetState(bankId)
	_ = json.Unmarshal(bankAsbytes, &bankAcc)

	for i := range bankAcc.Offers {
		if bankAcc.Offers[i].Id == offerId {
			bankAcc.Offers[i].Status = status
		}
	}

	newBankAsbytes, _ := json.Marshal(bankAcc)
	_ = stub.PutState(bankId, newBankAsbytes)
	if status == `approved` {
		var invoiceStack []invoice
		invoiceStackAsbytes, _ := stub.GetState(bankInvoicesKey)
		erro := json.Unmarshal(invoiceStackAsbytes, &invoiceStack)
		if erro != nil {
			return nil, errors.New("C:couldnt unmarshal invoice stack")
		}
		for invoiceIndex := range invoiceStack {
			if invoiceStack[invoiceIndex].InvoiceId == invoiceId {
				invoiceStack[invoiceIndex].Status = `completed`
			}
		}
		newInvoiceStackAsbytes, _ := json.Marshal(invoiceStack)
		_ = stub.PutState(bankInvoicesKey, newInvoiceStackAsbytes)
	}
	return nil, nil
}
