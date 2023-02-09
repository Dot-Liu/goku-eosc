package http_dubbo

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo/impl"
	"dubbo.apache.org/dubbo-go/v3/protocol/invocation"
	"fmt"
	"github.com/apache/dubbo-go-hessian2"
	"net"
	"reflect"
	"time"
)

// HttpToDubbo
// addr:dubbo://192.168.198.160:20000
func HttpToDubbo(addr string, serviceName, methodName string, typesList []string, valuesList []hessian.Object) (interface{}, error) {
	arguments := make([]interface{}, 3)
	parameterValues := make([]reflect.Value, 3)

	arguments[0] = methodName
	arguments[1] = typesList
	arguments[2] = valuesList

	parameterValues[0] = reflect.ValueOf(arguments[0])
	parameterValues[1] = reflect.ValueOf(arguments[1])
	parameterValues[2] = reflect.ValueOf(arguments[2])

	invoc := invocation.NewRPCInvocationWithOptions(invocation.WithMethodName("$invoke"),
		invocation.WithArguments(arguments),
		invocation.WithParameterValues(parameterValues))

	url, err := common.NewURL(addr,
		common.WithProtocol(dubbo.DUBBO), common.WithParamsValue(constant.SerializationKey, constant.ProtobufSerialization),
		common.WithParamsValue(constant.GenericFilterKey, "true"),
		common.WithParamsValue(constant.TimeoutKey, "5s"),
		common.WithParamsValue(constant.InterfaceKey, serviceName),
		common.WithParamsValue(constant.ReferenceFilterKey, "generic,filter"),
		//dubboAttachment must contains group and version info
		common.WithPath(serviceName),
	)
	if err != nil {
		return nil, err
	}
	dubboProtocol := dubbo.NewDubboProtocol()
	invoker := dubboProtocol.Refer(url)
	var resp interface{}
	invoc.SetReply(&resp)

	result := invoker.Invoke(context.Background(), invoc)
	if result.Error() != nil {
		return nil, result.Error()
	}

	return result.Result(), nil
}

func TcpToDubbo() {

	conn, err := net.Dial("tcp", "127.0.0.1:3333")
	if err != nil {
		panic(err)
	}
	methodName := "UpdateUser"

	attachments := make(map[string]interface{})
	attachments["jwt"] = "fdsf1ds23f1sdf5ds64fds123"

	var params []interface{}
	params = append(params, methodName)

	//todo types长度和values长度一一对应
	types := make([]string, 0)
	types = append(types, "object")

	values := make([]hessian.Object, 0)

	mm := make(map[string]interface{})
	mm["id"] = 10
	mm["name"] = "nihao"
	mm["age"] = "18"

	values = append(values, mm)

	params = append(params, types)
	params = append(params, values)

	//todo codec.EncodeRequest限制了序列化接口只能用Hessian2和proto
	//codec := &dubbo.DubboCodec{}
	//rpcInvocation := invocation.NewRPCInvocation(methodName, params, attachments)
	//rpcInvocation.SetAttachment(constant.InterfaceKey, "cn.zzy.api.UserService")
	//var invocationProtocol protocol.Invocation = rpcInvocation
	//
	//request := remoting.NewRequest("")
	//request.SerialID = constant.SHessian2
	//
	//request.Data = &invocationProtocol
	//
	//buffer, err := codec.EncodeRequest(request)
	//if err != nil {
	//	panic(err)
	//}

	//todo 可扩展序列化接口   需要实现impl.Serializer()接口  然后设置dubboPackage.SetSerializer()
	dubboPackage := impl.NewDubboPackage(nil)
	dubboPackage.Service = impl.Service{
		Path:      "/api.UserService",
		Interface: "api.UserService",
		Method:    "$invoke",       //todo 固定写死
		Timeout:   time.Second * 3, //request Timeout
	}

	dubboPackage.Header = impl.DubboHeader{
		SerialID: constant.SHessian2,
		Type:     impl.PackageRequest,
		ID:       10, //request.ID uuid
	}
	body := impl.NewRequestPayload(params, attachments)
	dubboPackage.Body = body

	buffer1, err := dubboPackage.Marshal()
	if err != nil {
		panic(err)
	}
	bytes1 := buffer1.Bytes()
	fmt.Println(string(bytes1))

	go read(conn)

	_, err = conn.Write(bytes1)

	select {}

}

func read(conn net.Conn) {
	//
	fmt.Println(conn.LocalAddr())
	buf := make([]byte, 512)
	fmt.Println("开始读")
	t := time.Now()
	n, err := conn.Read(buf)

	fmt.Println(n)
	fmt.Println(string(buf))
	fmt.Println("读取时间", time.Now().Sub(t).String())
	fmt.Println("结束读")
	fmt.Println(err)

}
