package main

import (
        "fmt"
        "math"
        "io"
        "net/http"
	"strings"
        "os"
        "log"
        "encoding/json"
        "io/ioutil"
        "strconv"
)

type Geo struct{
	Tpe string
	Coordinates []float64
}

type Prop struct{
	Name string
	Population int
} 

type Feat struct{
	Tpe string
	Id string
	Geometry Geo
	Properties Prop
}

type NewType struct {
	Tpe string
	Features []Feat
}

func downloadFromUrl(url string, filename string) {
        // TODO: check file existence first with io.IsExist
        output, err := os.Create(filename)
        if err != nil {
                fmt.Println("Error while creating kai.txt")
                return
        }
        defer output.Close()

        client := &http.Client{}
        req, _ := http.NewRequest("GET", url, nil)

        req.SetBasicAuth("api@undertone.com", "aiYoXo8I!")

        response, err := client.Do(req)
        if err != nil {
                fmt.Println("Error while downloading", url, "-", err)
                return
        }
        defer response.Body.Close()

        n, err := io.Copy(output, response.Body)
        if err != nil {
                fmt.Println("Error while downloading", url, "-", err)
                return
        }

        fmt.Println(n, "bytes downloaded.")
}

func topBrands() {
        url := "https://insights-admin.undertone.com/api/custom_report_115"
	filename := "wc_insights.txt"
        downloadFromUrl(url, filename)

        var data [][]string
        file, err := ioutil.ReadFile(filename)
        if err != nil {
                log.Fatal(err)
        }
        err = json.Unmarshal(file, &data)
        if err != nil {
                log.Fatal(err)
        }

        total := 0
        for _, pair := range data {
                i,_ := strconv.Atoi(pair[1])
                total += i
        }

        //open wc.txt for writing
        f, _ := os.Create("wc.txt")
        defer f.Close()

        for _, pair := range data {
                i,_ := strconv.Atoi(pair[1])
                count := int((float64(i) / float64(total)) * 1000)
                var printed string
                words := strings.Split(pair[0], "(")
                stripped := strings.Trim(words[0], " ")
                printed = strings.Replace(stripped, " ", "##", -1)
                for i := 0; i < count; i++ {
                        f.WriteString(printed)
                        f.WriteString(" ")
                }
        }
}

func ctr() {
        url := "https://insights-admin.undertone.com/api/custom_report_116"
	filename := "ctr_insights.txt"
        downloadFromUrl(url, filename)

        var data [][]string
        file, err := ioutil.ReadFile(filename)
        if err != nil {
                log.Fatal(err)
        }
        err = json.Unmarshal(file, &data)
        if err != nil {
                log.Fatal(err)
        }

        totalCtr := 0.0
        for _, pair := range data {
                i,_ := strconv.ParseFloat(pair[1], 32)
                totalCtr += i
        }

        f, _ := os.Create("ctr.json")
        defer f.Close()

        f.WriteString("{\n")
        f.WriteString("\t\"name\" : \"CTR\",\n")
        f.WriteString("\t\"children\": [\n")

        for i, pair := range data {
		if( i!= 0){
			f.WriteString(",\n")
		}
                i,_ := strconv.ParseFloat(pair[1], 32)
                font := math.Floor(i/totalCtr * 10000)
                f.WriteString("\t\t\t{\"name\" : \""+ pair[0] + "\", \"size\": ")
                f.WriteString(strconv.FormatFloat(font,'f',-1,32))
                f.WriteString("}")
        }
        f.WriteString("\n\t]\n")
        f.WriteString("}")
}

func handle_geo_data() {
	url2 := "https://insights-admin.undertone.com/api/custom_report_117"
	downloadFromUrl(url2, "geo_insights.txt")
	var data NewType
    	file, err := ioutil.ReadFile("/var/www/html/hackathon/geo/geo.json")
    	if err != nil {
        	log.Fatal(err)
    	}
    	err = json.Unmarshal(file, &data)
    	if err != nil {
        	log.Fatal(err)
    	}
	
	var states [][]string
	file1, err := ioutil.ReadFile("geo_insights.txt")
        if err != nil {
                log.Fatal(err)
        }
        err = json.Unmarshal(file1, &states)
        if err != nil {
                log.Fatal(err)
        }

	var m map[string]int
	m = make(map[string]int)
	for _, pair := range states {
		m[strings.ToUpper(pair[0])],_ = strconv.Atoi(pair[1])
	}	
	
	data.Tpe = "FeatureCollection"
	for i, _ := range data.Features {
		data.Features[i].Tpe = "Feature"
		data.Features[i].Geometry.Tpe = "Point"
		data.Features[i].Properties.Population = m[data.Features[i].Properties.Name]	
	}
	

	str, _ := json.Marshal(data)
	f, _ := os.Create("geo_out.json")
        defer f.Close()
	s := strings.Replace(string(str), "Tpe", "type", -1)
	s = strings.Replace(s, "Geometry", "geometry", -1)
	s = strings.Replace(s, "Features", "features", -1)
	s = strings.Replace(s, "Coordinates", "coordinates", -1)
	s = strings.Replace(s, "Properties", "properties", -1)
	s = strings.Replace(s, "Population", "population", -1)
	s = strings.Replace(s, "Name", "name", -1)
	s = strings.Replace(s, "Id", "id", -1)
	f.WriteString(s)
}

func main(){
	 topBrands()
	 ctr()
	 handle_geo_data()	
	
}


