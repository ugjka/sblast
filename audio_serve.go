package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"slices"
	"strings"
	"sync"

	"github.com/davecgh/go-spew/spew"
)

type stream struct {
	sink         string
	mime         string
	format       string
	bitrate      int
	chunk        int
	printheaders bool
	contentfeat  dlnaContentFeatures
	bitdepth     int
	samplerate   int
	channels     int
	nochunked    bool
	be           bool
}

func (s stream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.printheaders {
		spew.Fdump(os.Stderr, r.Proto)
		spew.Fdump(os.Stderr, r.RemoteAddr)
		spew.Fdump(os.Stderr, r.URL)
		spew.Fdump(os.Stderr, r.Method)
		spew.Fdump(os.Stderr, r.Header)
	}
	// Set some headers
	w.Header().Add("Cache-Control", "No-Cache, No-Store")
	w.Header().Add("Pragma", "No-Cache")
	w.Header().Add("Expires", "0")
	w.Header().Add("User-Agent", "sblast-DLNA UPnP/1.0 DLNADOC/1.50")
	// handle devices like Samsung TVs
	if r.Header.Get("GetContentFeatures.DLNA.ORG") == "1" {
		w.Header().Set("ContentFeatures.DLNA.ORG", s.contentfeat.String())
	}

	var yearSeconds = 365 * 24 * 60 * 60
	if r.Header.Get("Getmediainfo.sec") == "1" {
		w.Header().Set("MediaInfo.sec", fmt.Sprintf("SEC_Duration=%d", yearSeconds*1000))
	}
	w.Header().Add("Content-Type", s.mime)

	flusher, ok := w.(http.Flusher)
	chunked := ok && r.Proto == "HTTP/1.1" && !s.nochunked

	if !chunked {
		size := yearSeconds * (s.bitrate / 8) * 1000
		if s.bitrate == 0 {
			size = s.samplerate * s.bitdepth * s.channels * yearSeconds
		}
		w.Header().Add(
			"Content-Length",
			fmt.Sprint(size),
		)
	}

	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	endianess := "le"
	if s.be {
		endianess = "be"
	}
	parecCMD := exec.Command(
		"parec",
		"--device="+s.sink,
		"--client-name=sblast-rec",
		"--rate="+fmt.Sprint(s.samplerate),
		"--channels="+fmt.Sprint(s.channels),
		"--format="+fmt.Sprintf("s%d%s", s.bitdepth, endianess),
		"--raw",
	)

	var raw bool
	// wav can't have big endian
	var pcm = fmt.Sprintf("pcm_s%dle", s.bitdepth)
	if s.format == "lpcm" || s.format == "wav" {
		raw = true
	}
	if s.format == "lpcm" {
		// lpcm can have big endian
		s.format = fmt.Sprintf("s%d%s", s.bitdepth, endianess)
		pcm = fmt.Sprintf("pcm_s%d%s", s.bitdepth, endianess)
	}

	ffargs := []string{
		"-f", fmt.Sprintf("s%d%s", s.bitdepth, endianess),
		"-ac", fmt.Sprint(s.channels),
		"-ar", fmt.Sprint(s.samplerate),
		"-i", "-",
		"-f", s.format, "-",
	}
	if s.bitrate != 0 {
		ffargs = slices.Insert(
			ffargs,
			len(ffargs)-3,
			"-b:a", fmt.Sprintf("%dk", s.bitrate),
		)
	}
	if raw {
		ffargs = slices.Insert(
			ffargs,
			len(ffargs)-1,
			"-c:a", pcm,
		)
	}
	//spew.Dump(strings.Join(ffargs, " "))
	ffmpegCMD := exec.Command("ffmpeg", ffargs...)

	if *logsblast {
		fmt.Fprintln(os.Stderr, strings.Join(parecCMD.Args, " "))
		parecCMD.Stderr = os.Stderr
		fmt.Fprintln(os.Stderr, strings.Join(ffmpegCMD.Args, " "))
		ffmpegCMD.Stderr = os.Stderr
	}

	parecReader, parecWriter := io.Pipe()
	parecCMD.Stdout = parecWriter
	ffmpegCMD.Stdin = parecReader

	ffmpegReader, ffmpegWriter := io.Pipe()
	ffmpegCMD.Stdout = ffmpegWriter

	var wg sync.WaitGroup
	//defer fmt.Println("done")
	defer wg.Wait()

	err := parecCMD.Start()
	if err != nil {
		log.Printf("parec failed: %v", err)
		return
	}
	wg.Add(1)
	go func() {
		err := parecCMD.Wait()
		if err != nil && !strings.Contains(err.Error(), "signal") {
			log.Println("parec:", err)
		}
		wg.Done()
	}()

	err = ffmpegCMD.Start()
	if err != nil {
		log.Printf("ffmpeg failed: %v", err)
		return
	}
	wg.Add(1)
	go func() {
		err := ffmpegCMD.Wait()
		if err != nil && !strings.Contains(err.Error(), "signal") {
			log.Println("ffmpeg:", err)
		}
		ffmpegWriter.Close()
		wg.Done()
	}()
	if chunked {
		var (
			err error
			n   int
		)
		buf := make([]byte, (s.bitrate/8)*1000*s.chunk)
		if s.bitrate == 0 {
			buf = make([]byte, s.samplerate*s.bitdepth*s.channels*s.chunk)
		}
		for {
			n, err = ffmpegReader.Read(buf)
			if err != nil {
				break
			}
			_, err = w.Write(buf[:n])
			if err != nil {
				break
			}
			flusher.Flush()
		}
	} else {
		io.Copy(w, ffmpegReader)
	}

	if parecCMD.Process != nil {
		parecCMD.Process.Kill()
	}
	if ffmpegCMD.Process != nil {
		ffmpegCMD.Process.Kill()
	}
	parecReader.Close()
	parecWriter.Close()
	ffmpegReader.Close()
	ffmpegWriter.Close()
}
