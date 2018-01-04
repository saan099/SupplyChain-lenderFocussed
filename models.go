package main

type product struct {
	ProductName string  `json:"productName"`
	Quantity    int     `json:"quantity"`
	Rate        float64 `json:"rate"`
	Value       float64 `json:"value"`
}

type purchaseOrder struct {
	OrderId      string    `json:"orderId"`
	Date         string    `json:"date"`
	Time         string    `json:"time"`
	Products     []product `json:"products"`
	Buyer        string    `json:"buyer"`
	Supplier     string    `json:"supplier"`
	TotalValue   float64   `json:"totalValue"`
	CreditPeriod int       `json:"creditPeriod"`
	Status       string    `json:"status"`
}

type invoice struct {
	InvoiceId      string          `json:"invoiceId"`
	PurchaseOrders []purchaseOrder `json:"purchaseOrders"`
	Date           string          `json:"date"`
	Time           string          `json:"time"`
	Buyer          string          `json:"buyer"`
	Supplier       string          `json:"supplier"`
	Subtotal       float64         `json:"subtotal"`
	Taxes          float64         `json:"taxes"`
	Total          float64         `json:"total"`
	Status         string          `json:"status"`
}

type buyer struct {
	BuyerId        string                `json:"buyerId"`
	BuyerName      string                `json:"buyerName"`
	PurchaseOrders []purchaseOrder       `json:"purchaseOrders"`
	Invoices       []invoice             `json:"invoices"`
	DInvoices      []disbursementInvoice `json:"dInvoices"`
	EmployeeName   string                `json:"employeeName"`
	EmailId        string                `json:"emailId"`
}

type supplier struct {
	SupplierId     string                `json:"supplierId"`
	SupplierName   string                `json:"supplierName"`
	PurchaseOrders []purchaseOrder       `json:"purchaseOrders"`
	Invoices       []invoice             `json:"invoices"`
	DInvoices      []disbursementInvoice `json:"dInvoices"`
	EmployeeName   string                `json:"employeeName"`
	EmailId        string                `json:"emailId"`
	Offers         []offer               `json:"offers"`
}

type bank struct {
	BankId       string                `json:"bankId"`
	BankName     string                `json:"bankName"`
	EmployeeName string                `json:"employeeName"`
	EmailId      string                `json:"emailId"`
	DInvoices    []disbursementInvoice `json:"dInvoices"`
	Offers       []offer               `json:"offers"`
}

type disbursementInvoice struct {
	Details         invoice `json:"details"`
	DisRate         float64 `json:"disRate"`
	DisAmount       float64 `json:"disAmount"`
	Days            int     `json:"days"`
	DisbursedAmount float64 `json:"disbursedAmount"`
	PenaltyPerDay   float64 `json:"penaltyPerDay"`
	Recourse        string  `json:"recourse"`
	Bank            string  `json:"bank"`
	Date            string  `json:"date"`
	Time            string  `json:"time"`
}

type offer struct {
	Id      string              `json:"id"`
	Details disbursementInvoice `json:"details"`
	Status  string              `json:"status"`
}
