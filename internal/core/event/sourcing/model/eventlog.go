package model

type EventType string

const (
	AppidPut            EventType = "AppidPut"
	KVPut               EventType = "KVPut"
	UserAdd             EventType = "UserAdd"
	RoleAdd             EventType = "RoleAdd"
	UserGrantRole       EventType = "UserGrantRole"
	RoleGrantPermission EventType = "RoleGrantPermission"
)
