package model

type Entry struct {
	TableName string
	Events    []Event
}

type EventSourcingName int

const (
	AddNamespace EventSourcingName = iota
	AddAppid
	AddKV
)
