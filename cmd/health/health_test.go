package health

import (
	"net"
	"reflect"
	"testing"
)

func TestInc(t *testing.T) {
	ip := net.ParseIP("127.0.2.255")
	want := net.ParseIP("127.0.3.0")
	inc(ip)
	if !want.Equal(ip) {
		t.Fatalf("want %q, got %q", want, ip)
	}
}

func TestHosts(t *testing.T) {
	ts := []struct {
		cidr    string
		want    []string
		wanterr bool
	}{
		{
			cidr:    "192.168.178.0/30",
			want:    []string{"192.168.178.1", "192.168.178.2", "192.168.178.3"},
			wanterr: false,
		},
		{
			cidr:    "127.0.0.1/31",
			want:    []string{"127.0.0.1"},
			wanterr: false,
		},
		{
			cidr:    "10.0.0.255/32",
			want:    []string{},
			wanterr: false,
		},
		{
			cidr:    "not a cidr",
			want:    nil,
			wanterr: true,
		},
	}

	for _, tt := range ts {
		got, goterr := hosts(tt.cidr)
		if tt.wanterr && goterr == nil {
			t.Fatal("want error, got no error")
		}
		if !tt.wanterr && goterr != nil {
			t.Fatal("want no error, got error")
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("want %q, got %q", tt.want, got)
		}
	}
}
