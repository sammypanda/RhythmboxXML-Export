package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/user"
	"regexp"
	"strings"
)

// The structs for playlist.xml

type Location struct {
	XMLName xml.Name `xml:"location"`
	Path    string   `xml:",chardata"`
}

type Playlist struct {
	// XMLName xml.Name `xml:"playlist"`
	Name      string     `xml:"name,attr"`
	Locations []Location `xml:"location"`
}

type RhythmdbPlaylists struct {
	XMLName   xml.Name   `xml:"rhythmdb-playlists"`
	Playlists []Playlist `xml:"playlist"`
}

// End of the structs for playlist.xml

func main() {
	user, err := user.Current() // get user details

	if err != nil {
		panic(err)
	}

	path := "/home/" + user.Username + "/.local/share/rhythmbox/playlists.xml" // default rhythmbox path TODO: put in a config file
	rbPlaylists, _ := os.ReadFile(path)                                        // get go to read the file
	playlist := &RhythmdbPlaylists{}                                           // assign the pattern from the structs to playlist var

	manipulate, err := regexp.Compile(`file\:\/\/`) // used for manipulatedPath (removes "file://")

	if err != nil {
		panic(err)
	}

	xml.Unmarshal([]byte(rbPlaylists), playlist) // deserialise the xml to the structs

	for _, list := range playlist.Playlists {
		fmt.Println(list.Name) // output the playlist name
		for _, location := range list.Locations {
			// Simple playlists:
			manipulatedPath := strings.Replace(manipulate.ReplaceAllString(location.Path, ""), "%20", " ", -1) // remove "file://" and replace "%20" with " "

			fmt.Println(manipulatedPath) // temporary output of the path to put to the new file
		}
		fmt.Println("\n ") // put space between playlists
	}
}
