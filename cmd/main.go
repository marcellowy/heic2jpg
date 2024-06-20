package main

import (
	"flag"
	"gopkg.in/gographics/imagick.v3/imagick"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	sourceDir  = ""
	distDir    = ""
	sourceFile = ""
	distFile   = ""
)

func init() {
	flag.StringVar(&sourceDir, "s", "", "-s /var/path/to/your/source/dir")
	flag.StringVar(&distDir, "d", "", "-d /var/path/to/your/output/dir")

	flag.StringVar(&sourceFile, "i", "", "-i /var/path/to/your/input/filename")
	flag.StringVar(&distFile, "o", "", "-o /var/path/to/your/output/filename")
	flag.Parse()
}

func main() {

	if sourceDir == "" && sourceFile == "" {
		flag.PrintDefaults()
		return
	}

	if sourceFile != "" {
		if sourceFile != "" && distFile == "" {
			log.Println("need -o param")
			return
		}
		if err := convert(sourceFile, distFile); err != nil {
			log.Println(sourceFile, "==>", err)
			return
		}
		log.Println(sourceFile, "==>", distFile)
		log.Println("success")
		return
	}

	if sourceDir != "" && distDir == "" {
		log.Println("need -d param")
		return
	}

	dirEntry, err := os.ReadDir(sourceDir)
	if err != nil {
		log.Println(err)
		return
	}

	imagick.Initialize()
	defer imagick.Terminate()

	for _, v := range dirEntry {
		if v.IsDir() {
			continue
		}
		filename := sourceDir + "/" + v.Name()
		ext := path.Ext(filename)
		_, name := filepath.Split(filename)
		distName := distDir + "/" + name + ".jpg"
		if strings.ToLower(ext) == ".heic" {
			if err = convert(filename, distName); err != nil {
				log.Println(filename, "==>", err)
				break
			} else {
				log.Println(filename, "==>", distName)
			}
		}
	}
	log.Println("success")
}

func convert(s, d string) (err error) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	if err = mw.ReadImage(s); err != nil {
		log.Println("read image err: ", err.Error())
		return
	}

	if err = mw.WriteImages(d, false); err != nil {
		log.Println("write image err: " + err.Error())
		return
	}

	return nil
}
