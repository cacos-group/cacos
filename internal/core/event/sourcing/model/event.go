package model

import (
	"github.com/cacos-group/cacos/internal/core/event/sourcing/qparams"
)

type Event struct {
	EventType EventType
	Params    qparams.Metadatas
}

type Cancel struct {
}

func NewInfoAppidPutEvent(namespace string, appid string) Event {
	params := qparams.Metadatas{}
	params.Set(qparams.Key, GenAppidKey(namespace, appid))
	params.Set(qparams.Val, appid)

	return Event{
		EventType: InfoAppidPut,
		Params:    params,
	}
}

func NewAppidPutEvent(key string, val string) Event {
	params := qparams.Metadatas{}
	params.Set(qparams.Key, key)
	params.Set(qparams.Val, val)

	return Event{
		EventType: AppidPut,
		Params:    params,
	}
}

func NewKVPutEvent(key string, val string) Event {
	params := qparams.Metadatas{}
	params.Set(qparams.Key, key)
	params.Set(qparams.Val, val)

	return Event{
		EventType: KVPut,
		Params:    params,
	}
}

func NewUserAddEvent(key string, user string, password string) Event {
	params := qparams.Metadatas{}
	params.Set(qparams.Key, key)
	params.Set(qparams.User, user)
	params.Set(qparams.Password, password)

	return Event{
		EventType: UserAdd,
		Params:    params,
	}
}

func NewRoleAddEvent(key string, role string) Event {
	params := qparams.Metadatas{}
	params.Set(qparams.Key, key)
	params.Set(qparams.Role, role)

	return Event{
		EventType: RoleAdd,
		Params:    params,
	}
}

func NewUserGrantRoleEvent(key string, user string, role string) Event {
	params := qparams.Metadatas{}
	params.Set(qparams.Key, key)
	params.Set(qparams.User, user)
	params.Set(qparams.Role, role)

	return Event{
		EventType: UserGrantRole,
		Params:    params,
	}
}

func NewRoleGrantPermissionEvent(role string, key string, perm string) Event {
	params := qparams.Metadatas{}
	params.Set(qparams.Role, role)
	params.Set(qparams.Key, key)
	params.Set(qparams.Perm, perm)

	return Event{
		EventType: RoleGrantPermission,
		Params:    params,
	}
}
