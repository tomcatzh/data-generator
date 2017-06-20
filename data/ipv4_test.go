package data

import "testing"

func TestRandomIPv4(t *testing.T) {
	ip, err := newRandomIPv4("test", "192.168.0.1/30")
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}
	if ip.Title() != "test" {
		t.Errorf("Unexcepted title: %v", ip.Title())
	}

	c := ip.Clone()
	for i := 0; i < 20; i++ {
		s, _ := c.Data()
		if s != "192.168.0.0" && s != "192.168.0.1" && s != "192.168.0.2" && s != "192.168.0.3" {
			t.Errorf("Unexcepted data: %v", s)
		}
	}
}
