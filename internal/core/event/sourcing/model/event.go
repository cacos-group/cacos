package model

import (
	"net/url"
)

type Event struct {
	EventType EventType
	Args      url.Values
}

func NewAppidPutEvent(key string, val string) Event {
	uv := url.Values{}
	uv.Set("key", key)
	uv.Set("val", val)

	return Event{
		EventType: AppidPut,
		Args:      uv,
	}
}

func NewKVPutEvent(key string, val string) Event {
	uv := url.Values{}
	uv.Set("key", key)
	uv.Set("val", val)

	return Event{
		EventType: KVPut,
		Args:      uv,
	}
}

func NewUserAddEvent(key string, username string, password string) Event {
	uv := url.Values{}
	uv.Set("key", key)
	uv.Set("user", username)
	uv.Set("password", password)

	return Event{
		EventType: UserAdd,
		Args:      uv,
	}
}

func NewRoleAddEvent(key string, role string) Event {
	uv := url.Values{}
	uv.Set("key", key)
	uv.Set("role", role)

	return Event{
		EventType: RoleAdd,
		Args:      uv,
	}
}

func NewUserGrantRoleEvent(key string, user string, role string) Event {
	uv := url.Values{}
	uv.Set("key", key)
	uv.Set("user", user)
	uv.Set("role", role)

	return Event{
		EventType: UserGrantRole,
		Args:      uv,
	}
}

func NewRoleGrantPermissionEvent(role string, key string, perm string) Event {
	uv := url.Values{}
	uv.Set("role", role)
	uv.Set("key", key)
	uv.Set("range", "\\0")
	uv.Set("perm", perm)

	return Event{
		EventType: RoleGrantPermission,
		Args:      uv,
	}
}
