package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"log"
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

const usage = `Usage of RhythmboxXML-Export

	-to --trackonly only the track instead of the full path to song file

	-h --help prints help information

`

var isTrackOnly bool

func init() {
	// handling cli options/args
	flag.BoolVar(&isTrackOnly, "trackonly", false, "")
	flag.BoolVar(&isTrackOnly, "to", false, "")
}

func main() {
	flag.Usage = func() { fmt.Print(usage) } // import for substituting in a better help page
	flag.Parse()                             // handle the flags
	user, err := user.Current()              // get user details

	if err != nil {
		panic(err)
	}

	rbPath := "/home/" + user.Username + "/.local/share/rhythmbox"       // default rhythmbox path, TODO: put in a config file
	playlistPath := "/home/" + user.Username + "/Documents/rb-playlists" // default playlist path, TODO: put in a config file

	// Create directory if doesn't exist

	if _, err := os.Stat(playlistPath); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(playlistPath, os.ModePerm) // make the directory if it doesn't exist
		if err != nil {
			log.Println(err)
		}
	}

	rbPlaylists, _ := os.ReadFile(rbPath + "/playlists.xml") // get go to read the file
	playlist := &RhythmdbPlaylists{}                         // assign the pattern from the structs to playlist var

	manipulate, err := regexp.Compile(`file\:\/\/`) // used for manipulatedPath (removes "file://")

	if err != nil {
		panic(err)
	}

	xml.Unmarshal([]byte(rbPlaylists), playlist) // deserialise the xml to the structs

	for _, list := range playlist.Playlists {
		fmt.Println(list.Name) // output the playlist name
		if len(list.Locations) != 0 {

			// Simple playlists:

			f, _ := os.Create(playlistPath + "/" + list.Name + ".m3u")

			for _, location := range list.Locations {

				manipulatedPath := strings.Replace(manipulate.ReplaceAllString(location.Path, ""), "%20", " ", -1) // remove "file://" and replace "%20" with " "

				if isTrackOnly {
					dropPath, _ := regexp.Compile(`[^_]*\/`)
					manipulatedPath = dropPath.ReplaceAllString(manipulatedPath, "")
				}

				fmt.Println(manipulatedPath) // output path visually TODO: tie to verbose option?

				f.WriteString(manipulatedPath + "\n") // put to the new file
			}

			defer f.Close()
		} else {
			fmt.Println("empty (no tasks) OR auto-playlist (unsupported)")

			// TODO: if list has <conjunction/> tag then it's auto-playlist, otherwise it's empty
		}
		fmt.Println("\n ") // put space between playlists
	}
}
