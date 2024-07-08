package main

import (
	"os"

	"github.com/trvita/caldav-client/ui"
)

func main() {
	//ui.StartMenu("http://127.0.0.1:5232", os.Stdin)       //radicale
	ui.StartMenu("http://127.0.0.1:90/dav.php", os.Stdin) // baikal
}
