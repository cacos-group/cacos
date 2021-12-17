package model

import "fmt"

const (
	etcdNamespacePrefix = "/cacos/namespaces/"
	etcdNamespaceKey    = "/cacos/namespaces/%s"

	etcdAppidPrefix = "/cacos/services/%s/"
	etcdAppidKey    = "/cacos/services/%s/%s"

	// /cacos/kvs/{namespace}/{service}/
	storeKVPrefix = "/cacos/kvs/%s/%s/"
	// /cacos/kvs/{namespace}/{service}/{k}
	storeKV = "/cacos/kvs/%s/%s/%s"
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

func GenStoreKVPrefix(namespace string, service string) string {
	return fmt.Sprintf(storeKVPrefix, namespace, service)
}

func GenStoreKV(namespace string, service string, k string) string {
	return fmt.Sprintf(storeKV, namespace, service, k)
}
