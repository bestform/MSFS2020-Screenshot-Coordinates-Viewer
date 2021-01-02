package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := flag.String("p", "8080", "port to serve on")
	directory := flag.String("d", "", "the directory of static file to host")
	flag.Parse()

	if *directory == "" {
		println("Please specify the screenshot directory with -d")
		os.Exit(1)
	}

	http.HandleFunc("/", handleIndex(*directory))

	// resources files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./resources"))))

	// image files
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir(*directory))))

	log.Printf("Serving %s on HTTP port: %s\n", "./resources", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

type Marker struct {
	Lat      string
	Lon      string
	Filename string
}

func handleIndex(path string) func(writer http.ResponseWriter, request *http.Request) {

	return func(writer http.ResponseWriter, request *http.Request) {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}

		var markers []Marker

		for _, f := range files {
			if f.IsDir() || (!strings.HasSuffix(f.Name(), ".png") && !strings.HasSuffix(f.Name(), ".jpg") && !strings.HasSuffix(f.Name(), ".jpeg")){
				continue
			}

			parts := strings.Split(f.Name(), ".")
			geoFileName := strings.Join(parts[0:len(parts) - 1], ".") + ".geo"

			geoData, err := ioutil.ReadFile(path + string(os.PathSeparator) + geoFileName)
			if err != nil {
				continue
			}

			geoDataStr := string(geoData)
			latLon := strings.Split(geoDataStr, ",")

			markers = append(markers, Marker{
				Lat:      latLon[0],
				Lon:      latLon[1],
				Filename: f.Name(),
			})
		}


		tmpl := template.Must(template.ParseFiles("templates/index.html"))

		err = tmpl.Execute(writer, markers)
		if err != nil {
			println(err)
		}

	}

}
