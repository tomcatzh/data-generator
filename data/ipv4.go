package data

import (
	"math/rand"
	"net"
	"time"
)

var mask = [...]int{
	1,
	256,
	256 * 256,
	256 * 256 * 256,
}

type randomIPv4 struct {
	column
	ip       net.IP
	contains int
	rand     *rand.Rand
}

func (i *randomIPv4) Clone() columnData {
	return &randomIPv4{
		column:   i.column,
		ip:       i.ip,
		contains: i.contains,
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (i *randomIPv4) Data() (string, error) {
	ip := make(net.IP, len(i.ip))
	copy(ip, i.ip)
	add := i.rand.Intn(i.contains)

	for j := len(ip) - 1; j >= 0; j-- {
		ip[j] += byte((add / mask[len(ip)-j-1]) % 256)
	}

	return ip.String(), nil
}

func pow(x, y int) (result int) {
	result = 1

	for i := 0; i < y; i++ {
		result *= x
	}

	return
}

func newRandomIPv4(title string, cidr string) (*randomIPv4, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	ones, bits := ipnet.Mask.Size()

	return &randomIPv4{
		column: column{
			title: title,
		},
		ip:       ip.Mask(ipnet.Mask),
		contains: pow(2, bits-ones),
	}, nil
}
