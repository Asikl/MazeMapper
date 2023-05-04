package Cache

import (
	// "github.com/miekg/dns"
	"fmt"
	"log"
	"os"
)

type CacheStruct struct {
	Domain      string
	cacheFIX    map[Cachekey]Cachevalue
	DomainERROR ERROR
}

type Cachekey struct {
	CacheDomain string
	QType       uint16
	CacheIp     string
}
type Cachevalue struct {
	Flag bool
	NS   []string
	IP   []string
}

type ERROR struct {
	Edge   int //边的个数
	ERROR0 int //没有错误查询
	ERROR1 int //格式错误
	ERROR2 int //Server failure
	ERROR3 int //NXDOMAIN
	ERROR4 int //Not Implemented
	ERROR5 int //Refused
	ERROR6 int //time out
	ERROR7 int //数据包中没有NS记录，无法继续进行下去
	ERROR8 int //hijack\
	ERROR9 int //IPerror

	//ERROR8 int
}

// 先定义一个比较low的全局缓存cache
//var cacheFIX = make(map[Cachekey]Cachevalue, 0)

// func InitCache() map[Cachekey]Cachevalue {
// 	// 先定义一个比较low的全局缓存cache
// 	var cacheFIX = make(map[Cachekey]Cachevalue, 0)

// 	return cacheFIX
// }

//var DomainERROR ERROR

func (c *CacheStruct) ERROR0Init() {
	c.DomainERROR.ERROR0 = 0
}
func (c *CacheStruct) GetERROR0() (num int) {
	return c.DomainERROR.ERROR0
}
func (c *CacheStruct) AddERROR0() {
	c.DomainERROR.ERROR0++
}

func (c *CacheStruct) ERROR1Init() {
	c.DomainERROR.ERROR1 = 0
}
func (c *CacheStruct) GetERROR1() (num int) {
	return c.DomainERROR.ERROR1
}
func (c *CacheStruct) AddERROR1() {
	c.DomainERROR.ERROR1++
}

func (c *CacheStruct) ERROR2Init() {
	c.DomainERROR.ERROR2 = 0
}
func (c *CacheStruct) GetERROR2() (num int) {
	return c.DomainERROR.ERROR2
}
func (c *CacheStruct) AddERROR2() {
	c.DomainERROR.ERROR2++
}

func (c *CacheStruct) ERROR3Init() {
	c.DomainERROR.ERROR3 = 0
}
func (c *CacheStruct) GetERROR3() (num int) {
	return c.DomainERROR.ERROR3
}
func (c *CacheStruct) AddERROR3() {
	c.DomainERROR.ERROR3++
}

func (c *CacheStruct) ERROR4Init() {
	c.DomainERROR.ERROR4 = 0
}
func (c *CacheStruct) GetERROR4() (num int) {
	return c.DomainERROR.ERROR4
}
func (c *CacheStruct) AddERROR4() {
	c.DomainERROR.ERROR4++
}

func (c *CacheStruct) ERROR5Init() {
	c.DomainERROR.ERROR5 = 0
}
func (c *CacheStruct) GetERROR5() (num int) {
	return c.DomainERROR.ERROR5
}
func (c *CacheStruct) AddERROR5() {
	c.DomainERROR.ERROR5++
}

func (c *CacheStruct) ERROR6Init() {
	c.DomainERROR.ERROR6 = 0
}
func (c *CacheStruct) GetERROR6() (num int) {
	return c.DomainERROR.ERROR6
}
func (c *CacheStruct) AddERROR6() {
	c.DomainERROR.ERROR6++
}

func (c *CacheStruct) ERROR7Init() {
	c.DomainERROR.ERROR7 = 0
}
func (c *CacheStruct) GetERROR7() (num int) {
	return c.DomainERROR.ERROR7
}
func (c *CacheStruct) AddERROR7() {
	c.DomainERROR.ERROR7++
}

func (c *CacheStruct) ERROR8Init() {
	c.DomainERROR.ERROR8 = 0
}
func (c *CacheStruct) GetERROR8() (num int) {
	return c.DomainERROR.ERROR8
}
func (c *CacheStruct) AddERROR8() {
	c.DomainERROR.ERROR8++
}

func (c *CacheStruct) ERROR9Init() {
	c.DomainERROR.ERROR9 = 0
}
func (c *CacheStruct) GetERROR9() (num int) {
	return c.DomainERROR.ERROR9
}
func (c *CacheStruct) AddERROR9() {
	c.DomainERROR.ERROR9++
}

func (c *CacheStruct) EdegeInit() {
	c.DomainERROR.Edge = 0
}

func (c *CacheStruct) AddEdge(num int) {
	c.DomainERROR.Edge += num
}
func (c *CacheStruct) GetEdge() (num int) {
	return c.DomainERROR.Edge
}

func Init(domain string, c *CacheStruct) {

	c.cacheFIX = make(map[Cachekey]Cachevalue, 0)
	c.ERROR1Init()
	c.ERROR2Init()
	c.ERROR3Init()
	c.ERROR4Init()
	c.ERROR5Init()
	c.ERROR6Init()
	c.ERROR7Init()
	c.ERROR8Init()
	c.ERROR9Init()
}

func (c *CacheStruct) Add(domain string, server string, Qtype uint16, value Cachevalue) {
	var temp Cachekey
	temp.CacheDomain = domain
	temp.CacheIp = server
	temp.QType = Qtype
	c.cacheFIX[temp] = value
}

func (c *CacheStruct) GetCache(domain string, server string, Qtype uint16) (value Cachevalue) {
	var key Cachekey
	key.CacheDomain = domain
	key.CacheIp = server
	key.QType = Qtype
	return c.cacheFIX[key]
}

func (c *CacheStruct) Has(domain string, server string, Qtype uint16) (flag bool) {
	var temp Cachekey
	temp.CacheDomain = domain
	temp.CacheIp = server
	temp.QType = Qtype
	if _, ok := c.cacheFIX[temp]; ok {
		//fmt.Println("缓存中已有，是重复节点！")
		return true
		//存在
	} else {
		return false
	}
}

func (c *CacheStruct) Dump() {
	fmt.Println(c.cacheFIX)
}

func WriteTimeout(domain string, file os.File) {
	//str := domain + '\t'
	str := fmt.Sprintf("\"%s\" \n", domain)
	if _, err := file.WriteString(str); err != nil {
		log.Fatal(err)
	}

}
