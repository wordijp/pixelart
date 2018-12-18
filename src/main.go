package main

import (
	"fmt"
	"log"
	"strings"

	"pixela_art/src/lib/file"
	"pixela_art/src/lib/svg"
)

func main() {
	buf, err := file.ReadAll("./vim-pixela.svg")
	if err != nil {
		log.Fatal(err)
		return
	}

	data, err := svg.ParsePixelaSvg(strings.NewReader(string(buf)))
	if err != nil {
		log.Fatal(err)
		return
	}

	// TEST: 読み込めてる？
	for _, elem := range data.Elems {
		fmt.Printf("{date: %s, count: %3d, color: %s}\n",
			elem.Date.GetString(),
			elem.Count,
			elem.Rgb.ToColorCode())
	}
}
