package metadata

import "net/url"

type Metadata string

const (
	Key       Metadata = "key"
	Val       Metadata = "val"
	Namespace Metadata = "namespace"
	Appid     Metadata = "appid"
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
