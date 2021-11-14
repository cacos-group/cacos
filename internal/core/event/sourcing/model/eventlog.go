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

	AppidDel             EventType = "AppidDel"
	UserDel              EventType = "UserDel"
	RoleDel              EventType = "RoleDel"
	UserRevokeRole       EventType = "UserRevokeRole"
	RoleRevokePermission EventType = "RoleRevokePermission"
)
