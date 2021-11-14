package model

import "fmt"

const (
	etcdNamespacePrefix = "/info/namespace/"
	etcdNamespaceKey    = "/info/namespace/%s"

	etcdAppidPrefix = "/info/appid/%s/"
	etcdAppidKey    = "/info/appid/%s/%s"
)

func GenNamespacePrefix() string {
	return fmt.Sprintf(etcdNamespacePrefix)
}

func GenNamespaceKey(namespace string) string {
	return fmt.Sprintf(etcdNamespaceKey, namespace)
}

func GenAppidPrefix(namespace string) string {
	return fmt.Sprintf(etcdAppidPrefix, namespace)
}

func GenAppidKey(namespace string, appid string) string {
	return fmt.Sprintf(etcdAppidKey, namespace, appid)
}
