package model

type EventType string

const (
	InfoNamespacePut    EventType = "InfoNamespacePut"
	InfoAppidPut        EventType = "InfoAppidPut"
	AppidPut            EventType = "AppidPut"
	KVPut               EventType = "KVPut"
	UserAdd             EventType = "UserAdd"
	RoleAdd             EventType = "RoleAdd"
	UserGrantRole       EventType = "UserGrantRole"
	RoleGrantPermission EventType = "RoleGrantPermission"

	InfoAppidDel EventType = "InfoAppidDel"
	AppidDel     EventType = "AppidDel"
	KVDel        EventType = "KVDel"
	UserDel      EventType = "UserDel"
	RoleDel      EventType = "RoleDel"
)
