package spotdl

import (
    "os"
    "os/exec"
	"net/http"
    "io"
    "io/ioutil"
    "fmt"
    "regexp"
    "errors"
    "encoding/json"
    "strconv"
    "bufio"
    "strings"
    "path"

    "github.com/JoshuaDoes/spotigo"
    // "github.com/dhowden/tag"
    id3 "github.com/bogem/id3v2"
    // id3go "github.com/mikkyang/id3-go"
    // "github.com/landaire/go-taglib"
)

// Config ...
type Config struct {
    SpotigoHost string
    SpotigoPass string
    DefaultMusicDir string
}

var (
    Client *spotigo.Client
    MusicDir string
)

func LoadConfig(filepath string) (*Config, error) {
    file, err := os.Open(filepath)
    if err != nil {
        return nil, err
    }
    bytes, err := ioutil.ReadAll(file)
    if err != nil {
        return nil, err
    }
    config := &Config{}
    json.Unmarshal(bytes, &config)
    return config, nil
}

func CreateConfig(dir string) (error) {
    err := os.MkdirAll(dir, 0777)
    if err != nil {
        return err
    }
    file, err := os.Create(path.Join(dir, "config.json"))
    if err != nil {
        return err
    }
    writer := bufio.NewWriter(file)

    fmt.Fprintf(writer, "{\n\"SpotigoHost\": \"\",\n\"SpotigoPass\": \"\",\n\"DefaultMusicDir\": \"\"\n}")

    writer.Flush()
    return nil
}

func ReadList(filepath string) ([]string, error) {
    var list []string

    file, err := os.Open(filepath)
    if err != nil {
        return nil, err
    }

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        list = append(list, scanner.Text())
        // fmt.Println(scanner.Text())
    }

    for i := 0; i < len(list); i++ {
        fmt.Println(list)
    }

    return list, nil
}

func CreateList(filepath string, url string, urltype int) (error) {
    os.Remove(filepath)
    file, err := os.Create(filepath)
    if err != nil {
        return err
    }
    writer := bufio.NewWriter(file)

    if urltype == 1 {
        album, err := Client.GetAlbumInfo(url)
        if err != nil {
            return err
        }
        for _, track := range album.Discs[0].Tracks {
            fmt.Fprintln(writer, GetURL(track.URI, 3))
        }
    } else if urltype == 2 {
        playlist, err := Client.GetPlaylist(url)
        if err != nil {
            return err
        }
        for _, track := range playlist.Contents.Items {
            fmt.Fprintln(writer, GetURL(track.TrackURI, 3))
        }
    } else if urltype == 3 {
        fmt.Fprintln(writer, url)
    } else {
        return errors.New("Invalid value for urltype")
    }

    writer.Flush()

    return nil
}

func AppendToList(filepath string, url string, urltype int) (error) {
    oldFilePath := filepath + ".old"
    err := os.Rename(filepath, oldFilePath)
    if err != nil {
        return err
    }
    file, err := os.Create(filepath)
    if err != nil {
        return err
    }
    oldFile, err := os.OpenFile(oldFilePath, os.O_RDONLY, 0777)
    if err != nil {
        return err
    }

    scanner := bufio.NewScanner(oldFile)
    // scanner.Scan() // this allows the preservation of old contents of the file
    for scanner.Scan() {
        fmt.Fprintln(file, scanner.Text())
    }
    // writer := bufio.NewWriter(file)

    if urltype == 1 {
        album, err := Client.GetAlbumInfo(url)
        if err != nil {
            return err
        }
        for _, track := range album.Discs[0].Tracks {
            fmt.Fprintln(file, GetURL(track.URI, 3))
        }
    } else if urltype == 2 {
        playlist, err := Client.GetPlaylist(url)
        if err != nil {
            return err
        }
        for _, track := range playlist.Contents.Items {
            fmt.Fprintln(file, GetURL(track.TrackURI, 3))
        }
    } else if urltype == 3 {
        fmt.Fprintln(file, url)
    } else {
        return errors.New("Invalid value for urltype")
    }

    file.Close()
    oldFile.Close()
    os.Remove(oldFilePath)

    return nil
}

func GetURL(uri string, uritype int) string {
    if uritype == 3 {
        id := strings.Replace(uri, "spotify:track:", "", -1)
        return "https://open.spotify.com/track/" + id
    } else if uritype == 1 {
        id := strings.Replace(uri, "spotify:album:", "", -1)
        return "https://open.spotify.com/album/" + id
    }

    return ""
}

func DownloadList(filepath string) (error) {
    file, err := os.Open(filepath)
    if err != nil {
        return err
    }

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        err := GetTrack(scanner.Text())
        if err != nil {
            return err
        }
    }
    return nil
}

func GetTrack(url string) (error) {
    var (
        trackURL string
        trackName string
        trackMP3 string
        trackOGG string
    )

    trackURL = url
        
	track, err := Client.GetTrackInfo(trackURL)
	if err != nil {
		return(err)
    }

    info, err := GetTrackInfo(trackURL)
    if err != nil {
        return(err)
    }

    trackNumber := ""

    if info.TrackNumber < 10 {
        trackNumber = "0" + strconv.Itoa(info.TrackNumber)
    } else {
        trackNumber = strconv.Itoa(info.TrackNumber)
    }

    trackDir := path.Join(MusicDir, track.Artist, info.Album.Name)

    err = os.Mkdir(MusicDir, 0777)
    err = os.MkdirAll(trackDir, 0777)
    if err != nil {
        return err
    }

    trackName = trackNumber + " " + track.Title
    trackMP3 = path.Join(trackDir, trackName + ".mp3")
    trackOGG = path.Join(trackDir, trackName + ".ogg")

    os.Remove(trackMP3)

	fmt.Println("Downloading " + track.Title)
	DownloadFile(trackOGG, track.StreamURL)
    fmt.Println("Done!")
    err = exec.Command("ffmpeg", "-i", trackOGG, trackMP3).Run()
    if err != nil {
        return(err)
    }
    DownloadFile("cover.jpeg", track.ArtURL)

    err = SetMetadata(trackMP3, *info)
    if err != nil {
        return(err)
    }

    os.Remove(trackOGG)
    os.Remove("cover.jpeg")

    return nil
}

func DownloadFile(filepath string, url string) error {

    // Create the file
    out, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer out.Close()

    // Get the data
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Write the body to file
    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return err
    }

    return nil
}

func GetMetadata(filepath string) (*id3.Tag, error) {
    data, err := id3.Open(filepath, id3.Options{Parse: true})
    if err != nil {
        return nil, err
    }
    return data, nil
}

func SetMetadata(filepath string, track spotigo.SpotigoTrackInfo) (error) {
    file, err := id3.Open(filepath, id3.Options{Parse: true})
    if err != nil {
        return err
    }

    cover, err := ioutil.ReadFile("cover.jpeg")
    if err != nil {
        return err
    }

    albumGid := &spotigo.SpotigoGid{Gid: track.Album.Gid}
    albumID, err := albumGid.GetID(Client.Host, Client.Pass)
    if err != nil {
        return err
    }

    album, err := Client.GetAlbumInfo(GetURL("spotify:album:" + albumID, 1))
    if err != nil {
        return err
    }

    file.SetDefaultEncoding(id3.EncodingUTF8)
    file.SetVersion(4)
    file.SetTitle(track.Name)
    file.SetArtist(track.Artist[0].Name)
    file.SetAlbum(track.Album.Name)
    file.SetYear(strconv.Itoa(track.Album.Date.Year))

    albumCover := id3.PictureFrame {
        Encoding: id3.EncodingUTF8,
        MimeType: "image/jpeg",
        PictureType: id3.PTFrontCover,
        Description: "Cover",
        Picture: cover,
    }

    file.AddAttachedPicture(albumCover)
    file.AddTextFrame("TRCK",
                      id3.EncodingUTF8,
                      strconv.Itoa(track.TrackNumber) + "/" + strconv.Itoa(len(album.Discs[0].Tracks)))

    err = file.Save()
    if err != nil {
        return err
    }

    return nil
}

func GetTrackInfo(url string) (*spotigo.SpotigoTrackInfo, error) {
    regex := regexp.MustCompile("^(https:\\/\\/open.spotify.com\\/track\\/|spotify:track:)([a-zA-Z0-9]+)(.*)$")
	trackID := regex.FindStringSubmatch(url)
	if len(trackID) <= 0 {
		return nil, errors.New("error finding track ID")
	}

	trackJSON, err := http.Get(fmt.Sprintf("http://%s/track/%s?pass=%s", Client.Host, trackID[len(trackID)-2], Client.Pass))
	if err != nil {
		return nil, err
    }

	trackInfo := &spotigo.SpotigoTrackInfo{}
	err = unmarshal(trackJSON, trackInfo)
	if err != nil {
		return nil, err
	}
	if trackInfo.Name == "" {
		return nil, errors.New("error getting track info")
    }
    return trackInfo, nil
}

func unmarshal(body *http.Response, target interface{}) error {
	defer body.Body.Close()
	return json.NewDecoder(body.Body).Decode(target)
}