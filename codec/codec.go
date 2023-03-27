package codec

import (
	"io"
)

/**
* TODO:
* 1. 将 GobCodec 修改为：Protobuf
* 2. 增加 ebpf 读取ip地址的函数: ReadeBpf(interface{}) error : 作用是读取这个rpc调用的进程所在linux内核端口和进程IP地址
 */

type Header struct {
	// 方法名称，标识为：“Service.Method”
	ServiceMethod string
	// 序列号，客户端设置，用于标识当前这个数据包序列，表示不同的请求
	Seq uint64
	// 错误类型
	Error string
}

type Codec interface {
	//
	io.Closer

	ReadHeader(*Header) error

	ReadBody(interface{}) error

	Write(*Header, interface{}) error
}

type NewCodecFunc func(io.ReadWriteCloser) Codec

type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json"
)

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
