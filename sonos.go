package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"strings"

	"github.com/huin/goupnp"
	"github.com/huin/goupnp/dcps/av1"
)

func detectSonos(dev *goupnp.MaybeRootDevice) bool {
	var xmldata struct {
		Device struct {
			Manufacturer string `xml:"manufacturer"`
		} `xml:"device"`
	}

	clients, err := av1.NewAVTransport1ClientsByURL(dev.Location)
	if err != nil {
		return false
	}
	for _, client := range clients {
		resp, err := http.Get(client.Location.String())
		if err != nil {
			return false
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return false
		}
		err = xml.Unmarshal(data, &xmldata)
		if err != nil {
			return false
		}
		man := xmldata.Device.Manufacturer
		man = strings.ToLower(man)
		if strings.Contains(man, "sonos") {
			return true
		}
	}
	return false
}
