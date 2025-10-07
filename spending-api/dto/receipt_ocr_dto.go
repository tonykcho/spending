package dto

import "time"

type ReceiptOcrDto struct {
	StoreName string
	Date      time.Time
	Items     []ReceiptItemOcrDto
}

type ReceiptItemOcrDto struct {
	Name  string
	Price float64
}
