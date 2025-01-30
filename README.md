# sblast

<img src="logo.png" width=256px alt="sblast logo" title="sblast logo">

## Cast your Linux audio to DLNA receivers

You need `pactl`, `parec` and `ffmpeg` executables/dependencies on your system to run sblast.

If you have all that then you can launch `sblast` and it looks like this when you run it:

```
[user@user sblast]$ ./sblast 
----------
DLNA receivers
0: Kitchen
1: Phone
2: Bedroom
3: Livingroom TV
----------
Select the DLNA device:
[1]
----------
Audio sources
0: alsa_output.pci-0000_00_1b.0.analog-stereo.monitor
1: alsa_input.pci-0000_00_1b.0.analog-stereo
2: bluez_output.D8_AA_59_95_96_B7.1.monitor
3: sblast.monitor
----------
Select the audio source:
[2]
----------
Your LAN ip addresses
0: 192.168.1.14
1: 192.168.122.1
2: 2a04:ec00:b9ab:555:3c50:e6e8:8ea:211f
3: 2a04:ec00:b9ab:555:806d:800b:1138:8b1b
4: fe80::f4c2:c827:a865:35e5
----------
Select the lan IP address for the stream:
[0]
----------
2023/07/08 23:53:07 starting the stream on port 9000 (configure your firewall if necessary)
2023/07/10 23:53:07 stream URI: http://192.168.1.14:9000/stream.mp3
2023/07/08 23:53:07 setting av1transport URI and playing
```

There are also `-debug` and `-headers` flags if you want to inspect your DLNA device. Also, `-log` to inspect what parec and ffmpeg are doing.

### Non-interactive usage and extra flags

```
[ugjka@ugjka sblast]$ sblast -h
Usage of sblast:
  -bitrate int
        audio format bitrate (default 320)
  -bits int
        audio bitdepth (default 16)
  -channels int
        audio channels (default 2)
  -chunk int
        chunk size in seconds (default 1)
  -debug
        print debug info
  -device string
        dlna device's friendly name
  -dummy
        only serve content
  -format string
        stream audio format (default "mp3")
  -headers
        print request headers
  -ip string
        host ip address
  -log
        log parec and ffmpeg stderr
  -mime string
        stream mime type (default "audio/mpeg")
  -nochunked
        disable chunked tranfer endcoding
  -port int
        stream port (default 9000)
  -rate int
        audio sample rate (default 44100)
  -source string
        audio source (pactl list sources short | cut -f2)
  -useaac
        use aac audio
  -useflac
        use flac audio
  -uselpcm
        use lpcm audio
  -uselpcmle
        use lpcm little-endian audio
  -usewav
        use wav audio
  -version
        show sblast version
```

## Tips and tricks

* If you choose `sblast.monitor` as a source, you can send apps' audio to it (in pavucotrol or whatever applet you use) without streaming entire the desktop audio

<img src="img.sblast.monitor.png" width=300px alt="sblast.monitor example" title="sblast.monitor example">

* If none of the built-in codecs presets satisfy you, you can specify your own with `-mime` and `-format`. For example: `-mime audio/ac3 -format ac3`, `-mime audio/opus -format opus`, `-mime "audio/x-caf" -format caf` or `-mime "audio/mpeg" -format mp2`

* You can change audio features with `-rate`, `-bits` and `-channels`, e.g. `sblast -rate 48000 -bits 24 -channels 1`

## Building

You need the `go` and `go-tools` toolchain, also `git`

then execute:

```
git clone https://github.com/ugjka/sblast
cd sblast
go build
```

now you can run sblast with:
```
[user@user sblast]$ ./sblast
```

## Firewall

besides opening tcp port 9000 for the audio stream, you may also need to allow in pretty rad udp port range for ssdp discovery to work

ufw example, a subnet would be 192.168.1.0/24 or something like that

```bash
sudo ufw insert 1 allow from <subnet> to any port 32768:61000 proto udp
sudo ufw insert 1 allow from <subnet> to any port 9000 proto tcp
```

## Bins

Prebuilt Linux binaries are available on the releases [page](https://github.com/ugjka/sblast/releases)

## Archlinux

`paru -S sblast` or any other aur helper of your choice

## Why not use pulseaudio-dlna?

This is for pipewire-pulse users.

## Caveats

* You need to allow port 9000 from LAN for the DLNA receiver to be able to access the HTTP stream, you can change it with `-port` flag
* sblast monitor sink may not be visible in the pulse control applet unless you enable virtual streams

## Trivia

What on earth is "x-rincon-mp3radio"

## License

```
MIT LICENSE
Copyright (c) 2024 Uģis Ģērmanis

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights/
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
