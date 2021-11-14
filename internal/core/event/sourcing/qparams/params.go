package qparams

import "net/url"

type Metadata string

const (
	Key      Metadata = "key"
	Val      Metadata = "val"
	User     Metadata = "user"
	Password Metadata = "password"
	Role     Metadata = "role"
	Perm     Metadata = "perm"
)

type Metadatas url.Values

func (m Metadatas) Set(key Metadata, val string) {
	url.Values(m).Set(string(key), val)
}

func (m Metadatas) Get(key Metadata) string {
	return url.Values(m).Get(string(key))
}

func (m Metadatas) Encode() string {
	return url.Values(m).Encode()
}
