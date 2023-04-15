package Resolver

import (
	"context"
	"fmt"
	"time"

	"github.com/miekg/dns"
)

// 封装一下，方便调用 返回query msg
func NewMsg(Type uint16, domain string) *dns.Msg {
	return newMsg(Type, domain)
}

func newMsg(Type uint16, domain string) *dns.Msg {
	domain = dns.Fqdn(domain)
	msg := new(dns.Msg)
	msg.Id = dns.Id() //随机生成16bit的整数
	//msg.Id = 4096
	msg.RecursionDesired = false
	msg.Question = make([]dns.Question, 1)
	msg.Question[0] = dns.Question{
		Name:   domain,
		Qtype:  Type,
		Qclass: dns.ClassINET,
	}
	msg.SetEdns0(4096, false) // 设置 UDP 数据包最大长度为 4096 字节   true和false表示是否支持DNSSEC
	return msg
}

// Exchange 发送msg 接收响应
func (d *Dig) Exchange(m *dns.Msg) (msg1 *dns.Msg, err1 error, num int) {
	var msg *dns.Msg
	var err error
	for i := 0; i < d.SetRetry(3); i++ {
		msg, err, num = d.exchange(context.TODO(), m) //TODO返回一个空的context，todo 通常用在并不知道传递什么 context的情形
		if err == nil {
			return msg, err, num
		}
	}
	return msg, err, num
}

func (d *Dig) exchange(ctx context.Context, m *dns.Msg) (msg1 *dns.Msg, err1 error, num int) {
	var err error
	// var res dns.Msg
	// res.SetEdns0(4096, false)
	c := new(dns.Conn) //The new built-in function allocates memory. The first argument is a type, not a value,
	//and the value returned is a pointer to a newly allocated zero value of that type.
	c.UDPSize = 4096 //这句话必不可少，要不会崩溃，缓冲区溢出
	c.Conn, err = d.conn(ctx)

	if err != nil {
		fmt.Println("连接Conn error!", c.Conn, ctx)
		return nil, err, 63
	}
	defer c.Close()
	// SetWriteDeadline sets the deadline for future Write calls
	c.SetWriteDeadline(time.Now().Add(d.writeTimeout()))
	err = c.WriteMsg(m) //WriteMsg通过连接co发送消息。
	if err != nil {
		return nil, err, 61
	}
	// SetReadDeadline sets the deadline for future Read calls
	c.SetReadDeadline(time.Now().Add(d.readTimeout()))

	res, err := c.ReadMsg() //ReadMsg通过连接co读取消息。
	if err != nil {
		return nil, err, 62
	}
	if res.Id != m.Id { //there is a mismatch with the message's ID.
		return res, dns.ErrId, 9
	}
	return res, nil, 100
}

// GetMsg 返回msg响应体
func (d *Dig) GetMsg(Type uint16, domain string) (msg1 *dns.Msg, err1 error, num int) {
	m := newMsg(Type, domain)
	return d.Exchange(m)
}
