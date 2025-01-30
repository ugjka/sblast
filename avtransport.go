package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/huin/goupnp"
	"github.com/huin/goupnp/dcps/av1"
)

type avsetup struct {
	device    *goupnp.MaybeRootDevice
	stream    stream
	logoURI   string
	streamURI string
}

type avtransport interface {
	SetAVTransportURI(InstanceID uint32, CurrentURI string, CurrentURIMetaData string) (err error)
	Play(InstanceID uint32, Speed string) (err error)
	Stop(InstanceID uint32) (err error)
}

func detectAVtransport(dev *goupnp.MaybeRootDevice) string {
	transport := dev.Root.Device.FindService(av1.URN_AVTransport_1)
	if len(transport) > 0 {
		return av1.URN_AVTransport_1
	}
	transport = dev.Root.Device.FindService(av1.URN_AVTransport_2)
	if len(transport) > 0 {
		return av1.URN_AVTransport_2
	}
	return ""
}

func AVSetAndPlay(av avsetup) error {
	urn := detectAVtransport(av.device)
	var client avtransport

	switch {
	case urn == av1.URN_AVTransport_1:
		clients, err := av1.NewAVTransport1ClientsByURL(av.device.Location)
		if err != nil {
			return err
		}
		client = avtransport(clients[0])
	case urn == av1.URN_AVTransport_2:
		clients, err := av1.NewAVTransport2ClientsByURL(av.device.Location)
		if err != nil {
			return err
		}
		client = avtransport(clients[0])
	default:
		return fmt.Errorf("no avtransport found")
	}

	var err error
	try := func(metadata string) error {
		err = client.SetAVTransportURI(0, av.streamURI, metadata)
		if err != nil {
			return fmt.Errorf("set uri: %v", err)
		}
		time.Sleep(time.Second)
		err = client.Play(0, "1")
		if err != nil {
			return fmt.Errorf("play: %v", err)
		}
		return nil
	}

	metadata := fmt.Sprintf(
		didlTemplate,
		av.logoURI,
		av.stream.mime,
		av.stream.contentfeat,
		av.stream.bitdepth,
		av.stream.samplerate,
		av.stream.channels,
		av.streamURI,
	)
	metadata = strings.ReplaceAll(metadata, "\n", " ")
	metadata = strings.ReplaceAll(metadata, "> <", "><")

	err = try(metadata)
	if err == nil {
		return nil
	}
	log.Println(err)
	log.Println("trying without metadata")
	return try("")
}

func AVStop(device *goupnp.MaybeRootDevice) {
	urn := detectAVtransport(device)
	var client avtransport

	switch {
	case urn == av1.URN_AVTransport_1:
		clients, err := av1.NewAVTransport1ClientsByURL(device.Location)
		if err != nil {
			return
		}
		client = avtransport(clients[0])
	case urn == av1.URN_AVTransport_2:
		clients, err := av1.NewAVTransport2ClientsByURL(device.Location)
		if err != nil {
			return
		}
		client = avtransport(clients[0])
	default:
		return
	}

	client.Stop(0)
}

const didlTemplate = `<DIDL-Lite
xmlns="urn:schemas-upnp-org:metadata-1-0/DIDL-Lite/"
xmlns:upnp="urn:schemas-upnp-org:metadata-1-0/upnp/"
xmlns:dc="http://purl.org/dc/elements/1.1/"
xmlns:dlna="urn:schemas-dlna-org:metadata-1-0/"
xmlns:sec="http://www.sec.co.kr/"
xmlns:pv="http://www.pv.com/pvns/">
<item id="0" parentID="-1" restricted="1">
<upnp:class>object.item.audioItem.musicTrack</upnp:class>
<dc:title>Audio Cast</dc:title>
<dc:creator>sblast</dc:creator>
<upnp:artist>sblast</upnp:artist>
<upnp:albumArtURI>%s</upnp:albumArtURI>
<res protocolInfo="http-get:*:%s:%s"
bitsPerSample="%d"
sampleFrequency="%d"
nrAudioChannels="%d">%s</res>
</item>
</DIDL-Lite>`
