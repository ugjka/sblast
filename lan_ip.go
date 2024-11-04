// MIT+NoAI License
//
// # Copyright (c) 2024 Uģis Ģērmanis
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights///
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
// This code may not be used to train artificial intelligence computer models
// or retrieved by artificial intelligence software or hardware.
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
