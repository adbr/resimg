// 2014-12-30 Adam Bryt

// Program resimg zmienia rozmiar obrazków. Każdy obrazek (argument
// polecenia) jest konwertowany do rozmiaru określonego opcją -s
// (domyślnie 'medium' czyli '640x480') i kopiowany do katalogu
// docelowego (domyślnie '/tmp/small'). Katalog docelowy jest tworzony
// jeśli nie istnieje. Konwersja jest robiona poleceniem 'convert' z
// pakietu ImageMagick. Program służy do zmniejszania zdjęć przed
// wysłaniem emailem.
//
// Sposób użycia:
//
//	resimg [flags] file ...
//
// Opcje:
//
//	-s=medium
//		rozmiar obrazka (small, medium, large, lub w postaci '300x200')
//		(small=320x240,	medium=640x480, large=800x600)
//	-d=/tmp/small
//		katalog docelowy
//	-v=false
//		verbose: informuje co robi
//
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

var (
	dir     = flag.String("d", "/tmp/small", "katalog docelowy")
	size    = flag.String("s", "medium", "rozmiar obrazka (small, medium, large, lub w postaci '300x200')")
	verbose = flag.Bool("v", false, "verbose: informuje co robi")
)

const usageStr = `usage: resimg [flags] file ...
	-d="/tmp/small"
		katalog docelowy
	-s="medium"
		rozmiar obrazka (small, medium, large, lub w postaci '300x200')
	-v=false
		verbose: informuje co robi`

func usage() {
	fmt.Fprintln(os.Stderr, usageStr)
	os.Exit(1)
}

func mkdir(dir string) error {
	if *verbose {
		log.Printf("utworzenie katalogu '%s'", dir)
	}
	err := os.Mkdir(dir, os.ModePerm)
	if os.IsExist(err) {
		if *verbose {
			log.Printf("  katalog '%s' już istnieje", dir)
		}
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func resize(file string) error {
	if *verbose {
		log.Printf("konwersja pliku '%s'", file)
	}
	newfile := *dir + "/" + path.Base(file)
	cmd := exec.Command("convert", file, "-resize", *size, newfile)
	if *verbose {
		c := strings.Join(cmd.Args, " ")
		log.Printf("  %s", c)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		c := strings.Join(cmd.Args, " ")
		s := fmt.Sprintf("command '%s' failed:\n%s", c, string(out))
		return errors.New(s)
	}
	return nil
}

// Predefiniowane nazwy rozmiarów obrazków.
var sizes = map[string]string{
	"small":  "320x240",
	"medium": "640x480",
	"large":  "800x600",
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("resimg: ")

	flag.Usage = usage
	flag.Parse()

	// ustawia rozmiar na podstawie predefiniowanej nazwy
	s, ok := sizes[*size]
	if ok {
		*size = s
	}

	if flag.NArg() == 0 {
		usage()
	}

	err := mkdir(*dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, fname := range flag.Args() {
		err := resize(fname)
		if err != nil {
			log.Fatal(err)
		}
	}
}
