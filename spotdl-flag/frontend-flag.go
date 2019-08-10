package main

import (
    "../spotdl"
    "flag"
    "os"
    "fmt"
    "github.com/JoshuaDoes/spotigo"
)

var (
    client *spotigo.Client
    musicDir string
    config *spotdl.Config
)

func main() {
    listPtr := flag.Bool("l", false, "Whether to download or just append tracks to list")
    listFilePtr := flag.String("f", "list.txt", "File to save download list to")
    newlistPtr := flag.Bool("n", false, "Create new list instead of appending to existing")
    helpPtr := flag.Bool("h", false, "Display this message")
    trackPtr := flag.String("t", "", "Track URL to download")
    albumPtr := flag.String("a", "", "Album URL to download")
    playlistPtr := flag.String("p", "", "Playlist URL to download")
    musicDirPtr := flag.String("m", "", "Directory to save music to")

    // flag.BoolVar(listPtr, "l", false, "--list")
    // flag.StringVar(listFilePtr, "f", "list.txt", "--file")
    // flag.BoolVar(newlistPtr, "n", false, "--newlist")
    // flag.BoolVar(helpPtr, "h", false, "--help")
    // flag.StringVar(trackPtr, "t", "", "--track")
    // flag.StringVar(albumPtr, "a", "", "--album")
    // flag.StringVar(playlistPtr, "p", "", "--playlist")
    // flag.StringVar(musicDirPtr, "m", "", "--musicdir")

    flag.Parse()

    if *helpPtr == true {
        flag.PrintDefaults()
        os.Exit(0)
    }
    
    config, err := spotdl.LoadConfig("config.json")
    if err != nil {
        fmt.Printf("err loading config: %v", err)
        os.Exit(1)
    }

    client = &spotigo.Client{
        Host: config.SpotigoHost,
        Pass: config.SpotigoPass,
    }

    if config.DefaultMusicDir == "" {
        fmt.Println("Please edit the config.json file before running again.")
        os.Exit(0)
    }

    if *musicDirPtr != "" {
        musicDir = *musicDirPtr
    } else {
        musicDir = config.DefaultMusicDir
    }

    spotdl.MusicDir = musicDir
    spotdl.Client = client

    if *trackPtr != "" {
        if *newlistPtr == true {
            if *listFilePtr != "" {
                spotdl.CreateList(*listFilePtr, *trackPtr, 3) // 3 is single track
            } else {
                spotdl.CreateList("list.txt", *trackPtr, 3)
            }
        } else {
            if *listFilePtr != "" {
                spotdl.AppendToList(*listFilePtr, *trackPtr, 3)
            } else {
                spotdl.AppendToList("list.txt", *trackPtr, 3)
            }
        }
    } else if *albumPtr != "" {
        if *newlistPtr == true {
            if *listFilePtr != "" {
                spotdl.CreateList(*listFilePtr, *albumPtr, 1) // 1 is album 
            } else {
                spotdl.CreateList("list.txt", *albumPtr, 1)
            }
        } else {
            if *listFilePtr != "" {
                spotdl.AppendToList(*listFilePtr, *albumPtr, 1)
            } else {
                spotdl.AppendToList("list.txt", *albumPtr, 1)
            }
        }
    } else if *playlistPtr != "" {
        if *newlistPtr == true {
            if *listFilePtr != "" {
                spotdl.CreateList(*listFilePtr, *playlistPtr, 2) // 1 is album 
            } else {
                spotdl.CreateList("list.txt", *playlistPtr, 2)
            }
        } else {
            if *listFilePtr != "" {
                spotdl.AppendToList(*listFilePtr, *playlistPtr, 2)
            } else {
                spotdl.AppendToList("list.txt", *playlistPtr, 2)
            }
        }
    }

    if *listPtr == true {
        os.Exit(0)
    } else {
        if *listFilePtr != "" {
            spotdl.DownloadList(*listFilePtr)
        } else {
            spotdl.DownloadList("list.txt")
        }
    }
    // spotdl.DownloadList("list.txt")
    // err := spotdl.AppendToList("list.txt", "https://open.spotify.com/album/1CEODgTmTwLyabvwd7HBty", 1)
    // if err != nil {
    //     panic(err)
    // }
    // getTrack("https://open.spotify.com/track/7b0yOQHH8uSsMt0RsfDfCw")
}