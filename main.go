package main

import (
	"log"
	"mian/pkg/invoice"
	"mian/pkg/invoiceheader"
	"mian/pkg/invoiceitem"
	"mian/storage"
)

func main() {
	storage.NewPostgresDB()

	storageHeader := storage.NewPsqlInvoiceHeader(storage.Pool())
	storageItems := storage.NewPsqlInvoiceItem(storage.Pool())
	storageInvoice := storage.NewPsqlInvoice(
		storage.Pool(),
		storageHeader,
		storageItems,
	)

	m := &invoice.Model{
		Header: &invoiceheader.Model{
			Client: "Juan Galindo",
		},
		Items: invoiceitem.Models{
			&invoiceitem.Model{ProductID: 4},
			//&invoiceitem.Model{ProductID: 35},
		},
	}

	serviceInvoice := invoice.NewService(storageInvoice)
	if err := serviceInvoice.Create(m); err != nil {
		log.Fatalf("invoice.Create: %v", err)
	}
}
