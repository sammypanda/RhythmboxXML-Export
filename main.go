package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/user"
)

type Location struct {
	XMLName  xml.Name `xml:"location"`
	Location string   `xml:",chardata"`
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

func main() {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	path := "/home/" + user.Username + "/.local/share/rhythmbox/playlists.xml"
	rbPlaylists, _ := os.ReadFile(path)
	playlist := &RhythmdbPlaylists{}

	xml.Unmarshal([]byte(rbPlaylists), playlist)

	fmt.Println("test: " + user.Username)
	// fmt.Println(playlist.Playlists[7].Locations[0].Location)

	for _, list := range playlist.Playlists {
		fmt.Println(list.Name)
		for _, location := range list.Locations {
			fmt.Println(location)
		}
		fmt.Println("\n ")
	}
}
