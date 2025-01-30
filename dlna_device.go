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
