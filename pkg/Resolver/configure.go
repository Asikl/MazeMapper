package Resolver

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	dnsTimeout time.Duration = 5 * time.Second
)

// Root zone ipv4/6 servers
var root46servers = []string{
	"198.41.0.4",          //a发包发不过去
	"199.9.14.201",        //b
	"192.33.4.12",         //c
	"199.7.91.13",         //d
	"192.203.230.10",      //e
	"192.5.5.241",         //f
	"192.112.36.4",        //g
	"198.97.190.53",       //h
	"192.36.148.17",       //i
	"192.58.128.30",       //j
	"193.0.14.129",        //k
	"199.7.83.42",         //l
	"202.12.27.33",        //m
	"2001:503:ba3e::2:30", //a
	"2001:500:200::b",     //b
	"2001:500:2::c",       //c
	"2001:500:2d::d",      //d
	"2001:500:a8::e",      //e
	"2001:500:2f::f",      //f
	"2001:500:12::d0d",    //g
	"2001:500:1::53",      //h
	"2001:7fe::53",        //i
	"2001:503:c27::2:30",  //j
	"2001:7fd::1",         //k
	"2001:500:9f::42",     //l
	"2001:dc3::35"}        //m

// SetTimeOut set read write dial timeout
func (d *Dig) SetTimeOut(t time.Duration) {
	d.ReadTimeout = t
	d.WriteTimeout = t
	d.DialTimeout = t
}

// SetDNS 设置查询的dns server
func (d *Dig) SetDNS(host string) error {
	var ip string
	port := "53"
	switch strings.Count(host, ":") {
	case 0: //ipv4 no port
		ip = host
		d.RemoteAddr = fmt.Sprintf("[%s]:%v", ip, port)
	case 1: //ipv4 has port
		var err error
		ip, port, err = net.SplitHostPort(host)
		if err != nil {
			return err
		}
		d.RemoteAddr = fmt.Sprintf("[%s]:%v", ip, port)
	default: //ipv6
		// if net.ParseIP(host).To16() != nil {
		// 	//ParseIP parses s as an IP address
		// 	//To16 converts the IP address ip to a 16-byte representation. If ip is not an IP address (it is the wrong length), To16 returns nil.
		// 	ip = host
		// 	d.RemoteAddr = fmt.Sprintf("[%s]:%v", ip, port)
		// } else {
		// 	ip = host[:strings.LastIndex(host, ":")]
		// 	port = host[strings.LastIndex(host, ":")+1:]
		// }
		ip = host
		d.RemoteAddr = fmt.Sprintf("[%s]:%v", ip, port)
	}
	return nil
}

func (d *Dig) protocol() string {
	if d.Protocol != "" {
		return d.Protocol
	}
	return "udp"
}

func (d *Dig) dialTimeout() time.Duration {
	if d.DialTimeout != 0 {
		return d.DialTimeout
	}
	return dnsTimeout
}

func (d *Dig) readTimeout() time.Duration {
	if d.ReadTimeout != 0 {
		return d.ReadTimeout
	}
	return dnsTimeout
}

func (d *Dig) writeTimeout() time.Duration {
	if d.WriteTimeout != 0 {
		return d.WriteTimeout
	}
	return dnsTimeout
}

// 可以设置一下发多少数据包，我们一般默认是1
func (d *Dig) SetRetry(k int) int {
	return k
}

func (d *Dig) conn(ctx context.Context) (net.Conn, error) {
	remoteaddr := d.RemoteAddr
	di := net.Dialer{Timeout: d.dialTimeout()}           //func (d *Dialer) DialContext(ctx context.Context, network, address string) (Conn, error)
	return di.DialContext(ctx, d.protocol(), remoteaddr) //我们一般d.protocol()是udp
}
