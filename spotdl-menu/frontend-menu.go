package main

import (
	"../spotdl"
	"github.com/JoshuaDoes/spotigo"
	// gt "github.com/buger/goterm"
	"runtime"
	"fmt"
	"os"
	"os/exec"
	"bufio"
	"strconv"
	"strings"
)

var (
    client *spotigo.Client
    musicDir string
    config *spotdl.Config
)

// Menu ...
type Menu struct {
	header string
	menuitems []MenuItem
	parent *Menu
}

// MenuItem ...
type MenuItem struct {
	id int
	name string
	function string
}

func (menu *Menu) displayMenu() (string) {
	clear()
	fmt.Println(menu.header)
	fmt.Println("")
	for _, item := range menu.menuitems {
		fmt.Printf("%d. %s\n", item.id, item.name)
	}
	fmt.Println("")
	in := inputInt(": ")
	if in < 1 || in > len(menu.menuitems) {
		menu.displayMenu()
	}
	return menu.menuitems[in-1].function
}

func newMenu(header string, parent *Menu) (*Menu) {
	menu := &Menu {
		header: header,
		parent: parent,
	}

	// if parent == nil {
	// 	menu.menuitems = append(menu.menuitems, MenuItem {
	// 		name: "Exit",
	// 		function: "exit",
	// 	})
	// } else {
	// 	menu.menuitems = append(menu.menuitems, MenuItem {
	// 		name: "Back",
	// 		function: "back",
	// 	})
	// }

	return menu
}

func (menu *Menu) addMenuItem(name string, function string) {
	var id int
	if len(menu.menuitems) == 0 {
		id = 1
	} else {
		lastMenu := menu.menuitems[len(menu.menuitems)-1]
		id = lastMenu.id + 1
	}
	menu.menuitems = append(menu.menuitems, MenuItem {
		id: id,
		name: name,
		function: function,
	})
}

func main() {
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

    musicDir = config.DefaultMusicDir

    spotdl.MusicDir = musicDir
	spotdl.Client = client
	
	testMenu := newMenu("Test Menu", nil)
	testMenu.addMenuItem("Test Menu Item", "test")
	testMenu.addMenuItem("Exit", "exit")
	option := testMenu.displayMenu()
	switch option {
	case "test":
		fmt.Println("I work!")
		os.Exit(0)
	case "exit":
		os.Exit(0)
	default:
		panic("Whoops")
	}

}

func input(message string) (string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	text, _ := reader.ReadString('\n')
	return text
}

func inputInt(message string) (int) {
	text := input(message)
	i, err := strconv.Atoi(strings.Replace(text, "\r\n", "", -1)); if err != nil {
		// fmt.Println(err)
		return -1
	}
	return i
}

func clear() {
	ros := runtime.GOOS
	if ros == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		out, _ := exec.Command("clear").Output()
  		os.Stdout.Write(out)
	}
}