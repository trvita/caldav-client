package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"
	"golang.org/x/term"
)

// const url = "http://localhost:5232"
const url = "http://localhost:5232"

func main() {
	StartMenu()
}
