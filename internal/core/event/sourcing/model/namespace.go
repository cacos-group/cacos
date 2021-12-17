package model

import (
	"fmt"
	"github.com/cacos-group/cacos/internal/core/event/sourcing/pb"
	"github.com/golang/protobuf/proto"
	clientV3 "go.etcd.io/etcd/client/v3"
)

type HandlerReq struct {
	Namespace string
	Events    []*pb.Event
}

func GeneratorNamespaceReq(namespace string) *HandlerReq {
	return &HandlerReq{
		Namespace: namespace,
		Events:    generatorNamespaceEvents(namespace),
	}
}

func generatorNamespaceEvents(namespace string) []*pb.Event {
	return []*pb.Event{
		{
			Payload: &pb.Event_NamespacePut{NamespacePut: &pb.NamespacePut{
				Key: GenNamespaceKey(namespace),
				Val: "",
			}},
		},
	}
}

func GeneratorServiceReq(namespace string, service string) *HandlerReq {
	return &HandlerReq{
		Namespace: namespace,
		Events:    generatorServiceEvents(namespace, service),
	}
}

func generatorServiceEvents(namespace string, service string) []*pb.Event {
	key := GenAppidKey(namespace, service)

	//kv
	kvPerfix := GenStoreKVPrefix(namespace, service)

	user := fmt.Sprintf("u_%s_%s", namespace, service)
	// todo 生成password 和 加密
	password := "password"

	role := fmt.Sprintf("%s_%s", namespace, service)

	// 只读权限
	permissionType := int32(clientV3.PermRead)

	return []*pb.Event{
		{
			Payload: &pb.Event_ServicePut{ServicePut: &pb.ServicePut{
				Key: key,
				Val: "",
			}},
		},
		{
			Payload: &pb.Event_UserAdd{UserAdd: &pb.UserAdd{
				Key:      kvPerfix,
				User:     user,
				Password: password,
			}},
		},
		{
			Payload: &pb.Event_RoleAdd{RoleAdd: &pb.RoleAdd{
				Role: role,
			}},
		},
		{
			Payload: &pb.Event_UserGrantRole{UserGrantRole: &pb.UserGrantRole{
				User: user,
				Role: role,
			}},
		},
		{
			Payload: &pb.Event_RoleGrantPermission{RoleGrantPermission: &pb.RoleGrantPermission{
				Role: role,
				Key:  kvPerfix,
				Perm: permissionType,
			}},
		},
	}
}

func GeneratorKvReq(namespace string, service string, k string, v string) *HandlerReq {
	return &HandlerReq{
		Namespace: namespace,
		Events:    generatorKvEvents(namespace, service, k, v),
	}
}

func generatorKvEvents(namespace string, service string, k string, v string) []*pb.Event {
	return []*pb.Event{
		{
			Payload: &pb.Event_NamespacePut{NamespacePut: &pb.NamespacePut{
				Key: GenStoreKV(namespace, service, k),
				Val: v,
			}},
		},
	}
}

// 85 41
func MustMarshal(event *pb.Event) []byte {
	bytes, err := proto.Marshal(event)
	if err != nil {
		panic(err)
	}

	return bytes
}

func MustUnmarshal(data []byte, event *pb.Event) {
	err := proto.Unmarshal(data, event)
	if err != nil {
		panic(err)
	}
}

func ConvertToCancelEvent(event *pb.Event) *pb.Event {
	switch event.Payload.(type) {
	case *pb.Event_NamespacePut:
		return &pb.Event{Payload: &pb.Event_NamespaceDel{
			NamespaceDel: &pb.NamespaceDel{
				Key: event.GetNamespacePut().Key,
			},
		}}
	case *pb.Event_ServicePut:
		return &pb.Event{Payload: &pb.Event_ServiceDel{
			ServiceDel: &pb.ServiceDel{
				Key: event.GetServicePut().Key,
			},
		}}
	case *pb.Event_UserAdd:
		return &pb.Event{Payload: &pb.Event_UserDel{
			UserDel: &pb.UserDel{
				User: event.GetUserAdd().User,
			},
		}}
	case *pb.Event_RoleAdd:
		return &pb.Event{Payload: &pb.Event_RoleDel{
			RoleDel: &pb.RoleDel{
				Role: event.GetRoleAdd().Role,
			},
		}}
	default:
		return nil
	}
}
