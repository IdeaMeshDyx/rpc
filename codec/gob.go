package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

/**
* TODO:
* 1. 将 GobCodec 修改为：Protobuf
* 2. 增加 ebpf 读取ip地址的函数
 */

type GobCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *gob.Decoder
	enc  *gob.Encoder
}

/**
* 本段代码用于安全性检测
* 看是否GobCodec实现了Codec， _ 只是一个占位符
 */
var _ Codec = (*GobCodec)(nil)

// NewGobCodec 初始化函数，传入io.ReadWriteCloser对象并创建buf消息队列
func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}

func (c *GobCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}

func (c *GobCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *GobCodec) Write(h *Header, body interface{}) (err error) {
	// 当写入完数据之后，清空缓冲区并将缓冲区中的数据写入底层输出流
	defer func() {
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()
	//
	if err := c.enc.Encode(h); err != nil {
		log.Printf("rpc codec: gob error encoding header: %s\n", err)
		return err
	}
	if err := c.enc.Encode(body); err != nil {
		log.Panicf("rpc codec: gob err encoding body:%s", err)
		return err
	}
	return nil
}

func (c *GobCodec) Close() error {
	return c.conn.Close()
}
