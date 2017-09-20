package column

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/tomcatzh/data-generator/misc"
)

var mask = [...]int{
	1,
	256,
	256 * 256,
	256 * 256 * 256,
}

type ipv4Factory struct {
	ip       net.IP
	contains int
}

func newIPv4Factory(columnMethod int, c map[string]interface{}) (*ipv4Factory, error) {
	cidr, ok := c["CIDR"].(string)
	if !ok || cidr == "" {
		return nil, fmt.Errorf("column does not have IPv4 CIDR")
	}

	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	ones, bits := ipnet.Mask.Size()

	return &ipv4Factory{
		ip:       ip.Mask(ipnet.Mask),
		contains: misc.Pow(2, bits-ones),
	}, nil
}

type ipv4 struct {
	ipv4Factory
	rand *rand.Rand
}

func (i *ipv4Factory) Create() Column {
	return &ipv4{
		ipv4Factory: *i,
		rand:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (i *ipv4) Data() (string, error) {
	ip := make(net.IP, len(i.ip))
	copy(ip, i.ip)
	add := i.rand.Intn(i.contains)

	for j := len(ip) - 1; j >= 0; j-- {
		ip[j] += byte((add / mask[len(ip)-j-1]) % 256)
	}

	return ip.String(), nil
}
