package column

import "testing"

func TestIPv4(t *testing.T) {
	tmpl := map[string]interface{}{}
	tmpl["CIDR"] = "192.168.0.1/30"

	ip, err := newIPv4Factory(columnChangePerRow, tmpl)
	if err != nil {
		t.Errorf("Unexcepted error: %v", err)
	}

	c := ip.Create()
	for i := 0; i < 20; i++ {
		s, _ := c.Data()
		if s != "192.168.0.0" && s != "192.168.0.1" && s != "192.168.0.2" && s != "192.168.0.3" {
			t.Errorf("Unexcepted data: %v", s)
		}
	}
}
