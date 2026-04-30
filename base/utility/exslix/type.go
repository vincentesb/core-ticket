package exslix

import "github.com/xuri/excelize/v2"

type Filter struct {
	Label string
	Value string
	Style *excelize.Style
}

type Header struct {
	Attribute   string
	Label       string
	Position    int
	Style       *excelize.Style
	ColumnStyle *excelize.Style
}

type Value[T ~string | ~int | ~float64] struct {
	Attribute string
	Value     T
	Style     *excelize.Style
}

type Sheet[T any] struct {
	SheetName string
	Title     string
	Company   string
	Filters   []Filter
	Headers   []Header
	Values    []T
}
