// 2014-12-30 Adam Bryt

package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/adbr/resimg/internal/github.com/nfnt/resize"
)

var (
	dir      = flag.String("dir", "/tmp/small", "katalog docelowy")
	size     = flag.String("size", "medium", "rozmiar obrazka")
	verbose  = flag.Bool("v", false, "verbose")
	helpFlag = flag.Bool("help", false, "help")
)

func usage() {
	fmt.Fprintln(os.Stderr, usageStr)
	os.Exit(1)
}

func help() {
	fmt.Fprintln(os.Stdout, helpStr)
	os.Exit(0)
}

// makeDir tworzy katalog docelowy.
// Jeśli katalog już istnieje to drukuje warning.
func makeDir() error {
	if *verbose {
		log.Printf("utworzenie katalogu '%s'", *dir)
	}

	err := os.Mkdir(*dir, os.ModePerm)
	if os.IsExist(err) {
		log.Printf("warning: katalog '%s' już istnieje", *dir)
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

// resizeFile konwertuje obrazek z pliku file do rozmiaru max w, h
// z zachowaniem proporcji i zapisuje wynikowy (pomniejszony) obrazek do
// katalogu *dir.
// Obsługuje pliki w formatach gif, jpeg i png.
func resizeFile(file string, w, h uint) error {
	if *verbose {
		log.Printf("konwersja pliku '%s'", file)
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	img, typ, err := image.Decode(f)
	if err != nil {
		return fmt.Errorf("%s: %s", err, file)
	}

	img = resize.Thumbnail(w, h, img, resize.Lanczos3)

	nfile := filepath.Join(*dir, filepath.Base(file))
	nf, err := os.Create(nfile)
	if err != nil {
		return err
	}
	defer nf.Close()

	switch typ {
	case "gif":
		err := gif.Encode(nf, img, nil)
		if err != nil {
			return err
		}
	case "jpeg":
		err := jpeg.Encode(nf, img, nil)
		if err != nil {
			return err
		}
	case "png":
		err := png.Encode(nf, img)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("nie znany format obrazka: '%s'", typ)
	}

	return nil
}

// Predefiniowane nazwy rozmiarów obrazków.
// Tych nazw można używać jako wartości opcji -size.
var sizes = map[string]string{
	"small":  "320x240",
	"medium": "640x480",
	"large":  "800x600",
}

// parseSize parsuje rozmiar obrazka zawarty w *size.
// Rozmiar ma format typu '640x480' lub jest predefiniowaną nazwą rozmiaru:
// small, medium, large. Nazwa rozmiaru może być skrócona.
func parseSize() (width, height uint, err error) {
	// zamiana nazwy rozmiaru na rozmiar typu '640x480'
	for k, v := range sizes {
		if strings.HasPrefix(k, *size) {
			*size = v
			break
		}
	}

	a := strings.Split(*size, "x")
	if len(a) != 2 {
		err = fmt.Errorf("nie nie poprawny rozmiar: %q", *size)
		return
	}

	var n uint64
	n, err = strconv.ParseUint(a[0], 10, 0)
	if err != nil {
		return
	}
	width = uint(n)

	n, err = strconv.ParseUint(a[1], 10, 0)
	if err != nil {
		return
	}
	height = uint(n)

	return
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("resimg: ")

	flag.Usage = usage
	flag.Parse()

	if *helpFlag {
		help()
	}

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "brak nazwy pliku")
		usage()
	}

	err := makeDir()
	if err != nil {
		log.Fatal(err)
	}

	w, h, err := parseSize()
	if err != nil {
		log.Fatal(err)
	}

	for _, fname := range flag.Args() {
		err := resizeFile(fname, w, h)
		if err != nil {
			log.Fatal(err)
		}
	}
}
