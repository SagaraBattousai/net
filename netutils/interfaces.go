
package netutils

import (
  "fmt"
  "net"
  "strings"
  "log"
)

const (
  WIFI string = "WiFi"
  LAN string = "Local Area Connection* 11"
  LAN_BACKUP string = "Local Area Connection* 1"
)

type UnsupportedError string

func (err UnsupportedError) Error() string {
  return fmt.Sprintf("%s is currently unsupported!", string(err))
}

func DisplayInterfaces() {
  ifaces, _ := net.Interfaces()
  for _, i := range ifaces {
    addr, _ := i.Addrs()
    for _, a := range addr {
      fmt.Printf("%s: %s,  %s, %s\n\n", i.Name, a.(*net.IPNet).IP.To16().String(), a.Network(), a.String())
    }
  }
}

func IsIPv6(ip net.IP) bool {
  return strings.Contains(ip.String(), "::")
}

func GetCommunicationAddress() string {
  i, err := net.InterfaceByName(WIFI)
  if err != nil {
    i, err = net.InterfaceByName(LAN)
    if err != nil {
      i, err = net.InterfaceByName(LAN_BACKUP)
      if err != nil {
        log.Fatal(UnsupportedError("Network interface"))
      }
    }
  }
  addrs, err := i.Addrs()
  if err != nil {
    log.Fatal(err)
  }

  var ipAddress string

  //Usually only two and usually the second one is ipv6
  for i := len(addrs) - 1; i >= 0; i-- {
    ipNet, okay := addrs[i].(*net.IPNet)
    if !okay {
      log.Fatal("Address type")
    }
    ip := ipNet.IP
    if IsIPv6(ip) {
      return "[" + ip.String() + "]"
    }

    if ipAddress == "" {
      ipAddress = ip.String()
    }
  }

  return ipAddress
}
