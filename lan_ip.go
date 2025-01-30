package main

import (
	"fmt"
	"net"
)

func chooseStreamIP(lookup string) (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	ips := make([]net.IP, 0)
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok &&
			!ipnet.IP.IsLoopback() &&
			(ipnet.IP.To4() != nil || ipnet.IP.To16() != nil) {
			ips = append(ips, ipnet.IP)
		}
	}

	if len(ips) == 0 {
		return nil, fmt.Errorf("no usable lan ip addresses found")
	}
	if lookup != "" {
		lookupIp := net.ParseIP(lookup)
		if lookupIp == nil {
			return nil, fmt.Errorf("%s: not found", lookup)
		}
		for _, ip := range ips {
			if ip.Equal(lookupIp) {
				return lookupIp, nil
			}
		}
		return nil, fmt.Errorf("%s: not found", lookup)
	}
	fmt.Println("Your LAN ip addresses")
	for i, ip := range ips {
		fmt.Printf("%d: %s\n", i, ip)
	}

	fmt.Println("----------")
	fmt.Println("Select the lan IP address for the stream:")

	selected := selector(ips)
	return ips[selected], nil
}

func findInterface(ip net.IP) (string, error) {
	infs, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, inf := range infs {
		addrs, err := inf.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if addr.(*net.IPNet).IP.Equal(ip) {
				return inf.Name, nil
			}
		}
	}
	return "", fmt.Errorf("no interface found for ip")
}
