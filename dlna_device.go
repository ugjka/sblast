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

	"github.com/huin/goupnp"
	"github.com/huin/goupnp/dcps/av1"
)

func chooseUPNPDevice(lookup string) (*goupnp.MaybeRootDevice, error) {
	if lookup == "" {
		fmt.Println("Loading...")
	}

	roots, err := goupnp.DiscoverDevices(av1.URN_AVTransport_1)

	if lookup == "" {
		fmt.Print("\033[1A\033[K")
		fmt.Println("----------")
	}

	if err != nil {
		return nil, fmt.Errorf("discover: %v", err)
	}
	if lookup != "" {
		for _, v := range roots {
			if v.Root != nil {
				if v.Root.Device.FriendlyName == lookup {
					return &v, nil
				}
			}
		}
		return nil, fmt.Errorf("%s: not found", lookup)
	}

	if len(roots) == 0 {
		return nil, fmt.Errorf("no dlna devices on the network found")
	}
	fmt.Println("DLNA receivers")

	for i, v := range roots {
		if v.Root != nil {
			fmt.Printf("%d: %s\n", i, v.Root.Device.FriendlyName)
		}
	}

	fmt.Println("----------")
	fmt.Println("Select the DLNA device:")

	selected := selector(roots)
	return &roots[selected], nil
}
