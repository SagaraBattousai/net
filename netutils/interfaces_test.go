
package netutils

import (
  "testing"
  "net"
)

func TestIsIPv6(t *testing.T) {
  cases := []struct {
          ip net.IP
          expected bool
  }{
          {net.IPv4bcast,                 false},
          {net.IPv4allsys,                false},
          {net.IPv4allrouter,             false},
          {net.IPv4zero,                  false},
          {net.IPv6zero,                   true},
          {net.IPv6unspecified,            true},
          {net.IPv6loopback,               true},
          {net.IPv6interfacelocalallnodes, true},
          {net.IPv6linklocalallnodes,      true},
          {net.IPv6linklocalallrouters,    true},
  }
  for _, c := range cases {
    got := IsIPv6(c.ip)
    if got != c.expected {
      t.Errorf("IsIPv6(%q) should be %t not %t", c.ip, c.expected, got)
    }
  }
}
