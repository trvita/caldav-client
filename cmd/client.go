package main

import (
	"os"

	"github.com/trvita/caldav-client/menu"
)

func main() {
	//menu.StartMenu("http://127.0.0.1:5232", os.Stdin)       //radicale
	menu.StartMenu("http://127.0.0.1:90/dav.php", os.Stdin) // baikal
}
