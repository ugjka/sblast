package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

func chooseAudioSource(lookup string) (string, error) {
	srcCMD := exec.Command("pactl", "-f", "json", "list", "sources", "short")
	srcData, err := srcCMD.Output()
	if err != nil {
		return "", fmt.Errorf("pactl sources: %v", err)
	}

	var srcJSON Sources
	err = json.Unmarshal(srcData, &srcJSON)
	if err != nil {
		return "", err
	}
	if len(srcJSON) == 0 {
		return "", fmt.Errorf("no audio sources found")
	}
	// append for on-demand loading of sblast sink
	srcJSON = append(srcJSON, struct{ Name string }{sblastMONITOR})
	if lookup != "" {
		for _, v := range srcJSON {
			if v.Name == lookup {
				return lookup, nil
			}
		}
		return "", fmt.Errorf("%s: not found", lookup)
	}

	fmt.Println("Audio sources")
	for i, v := range srcJSON {
		fmt.Printf("%d: %s\n", i, v.Name)
	}

	fmt.Println("----------")
	fmt.Println("Select the audio source:")

	selected := selector(srcJSON)
	return srcJSON[selected].Name, nil
}

type Sources []struct {
	Name string
}
