package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	geo "github.com/oschwald/geoip2-golang"
)

type month struct {
	Hits         hits          `json:"hits"`
	RegionalHits regionalStats `json:"regionalhits"`
}

type regionalStats struct {
	NorthAmerica hits `json:"northamerica"`
	SouthAmerica hits `json:"southamerica"`
	Europe       hits `json:"europe"`
	Asia         hits `json:"asia"`
	Africa       hits `json:"africa"`
	Oceania      hits `json:"oceania"`
}

type hits struct {
	Bandwidth float64 `json:"bandwidth"`
	Qps       float64 `json:"qps"`
	TotalHits float64 `json:"totalhits"`
	TotalGet  float64 `json:"totalget"`
	TotalPost float64 `json:"totalpost"`
}

var maxmind *geo.Reader

func main() {
	var (
		err error
	)
	stats := make(map[string]month)
	next := make(chan bool, 1)
	next <- true
	maxmind, err = geo.Open("GeoIP2-City.mmdb")
	if err != nil {
		panic(err)
	}

	// 1. build a slice with a list of files to inspect
	// 2. open each file, gunzip if needed
	// 3. read each line and
	//	3.1 parse the date
	//	3.2 add to bandwidth stat to current month
	//	3.3 add current request to req/s rate
	//	3.4 geolocate src ip region
	// 4. print stats
	pattern := os.Args[1]
	files, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Fprintf(os.Stderr, "entering %s\n", file)
		fd, err := os.Open(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open %s: %v\n", file, err)
		}
		defer fd.Close()
		magic := make([]byte, 2)
		n, err := fd.Read(magic)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read magic number from %s: %v\n", file, err)
		}
		if n != 2 {
			fmt.Fprintf(os.Stderr, "read wrong number of bytes for magic number from %s: %v\n", file, err)
		}
		// rewind
		_, err = fd.Seek(0, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to seek(0,0) on %s: %v\n", file, err)
		}
		if magic[0] == 0x1f && magic[1] == 0x8b {
			// this is gzip, gunzip it
			gzip, err := gzip.NewReader(fd)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to gunzip %s: %v\n", file, err)
			}
			defer gzip.Close()
			<-next
			go processFile(gzip, stats, next)
		} else {
			<-next
			go processFile(fd, stats, next)
		}
	}
	<-next
	fmt.Printf("month; bandwidth; qps; totalhits; totalget; totalpost; " +
		"northamerica_bandwidth; northamerica_qps; northamerica_totalhits; northamerica_totalget; northamerica_totalpost; " +
		"southamerica_bandwidth; southamerica_qps; southamerica_totalhits; southamerica_totalget; southamerica_totalpost; " +
		"europe_bandwidth; europe_qps; europe_totalhits; europe_totalget; europe_totalpost; " +
		"asia_bandwidth; asia_qps; asia_totalhits; asia_totalget; asia_totalpost; " +
		"africa_bandwidth; africa_qps; africa_totalhits; africa_totalget; africa_totalpost; " +
		"oceania_bandwidth; oceania_qps; oceania_totalhits; oceania_totalget; oceania_totalpost;\n")
	for cmonth := range stats {
		s := stats[cmonth]
		fmt.Printf("%s; %.0f; %.3f; %.0f; %.0f; %.0f; %.0f; %.3f; %.0f; %.0f; %.0f; %.0f; %.3f; %.0f; %.0f; %.0f; %.0f; %.3f; %.0f; %.0f; %.0f; %.0f; "+
			"%.3f; %.0f; %.0f; %.0f; %.0f; %.3f; %.0f; %.0f; %.0f; %.0f; %.3f; %.0f; %.0f; %.0f;\n", cmonth, s.Hits.Bandwidth, s.Hits.Qps,
			s.Hits.TotalHits, s.Hits.TotalGet, s.Hits.TotalPost, s.RegionalHits.NorthAmerica.Bandwidth,
			s.RegionalHits.NorthAmerica.Qps, s.RegionalHits.NorthAmerica.TotalHits, s.RegionalHits.NorthAmerica.TotalGet,
			s.RegionalHits.NorthAmerica.TotalPost, s.RegionalHits.SouthAmerica.Bandwidth, s.RegionalHits.SouthAmerica.Qps,
			s.RegionalHits.SouthAmerica.TotalHits, s.RegionalHits.SouthAmerica.TotalGet, s.RegionalHits.SouthAmerica.TotalPost,
			s.RegionalHits.Europe.Bandwidth, s.RegionalHits.Europe.Qps, s.RegionalHits.Europe.TotalHits,
			s.RegionalHits.Europe.TotalGet, s.RegionalHits.Europe.TotalPost, s.RegionalHits.Asia.Bandwidth,
			s.RegionalHits.Asia.Qps, s.RegionalHits.Asia.TotalHits, s.RegionalHits.Asia.TotalGet, s.RegionalHits.Asia.TotalPost,
			s.RegionalHits.Africa.Bandwidth, s.RegionalHits.Africa.Qps, s.RegionalHits.Africa.TotalHits,
			s.RegionalHits.Africa.TotalGet, s.RegionalHits.Africa.TotalPost, s.RegionalHits.Oceania.Bandwidth,
			s.RegionalHits.Oceania.Qps, s.RegionalHits.Oceania.TotalHits, s.RegionalHits.Oceania.TotalGet,
			s.RegionalHits.Oceania.TotalPost)
	}
	jsondata, _ := json.Marshal(stats)
	fmt.Printf("\n%s\n", jsondata)

}

// example logline:
// 195.1.2.3 www.mozilla.org - [24/Sep/2015:01:00:28 -0700] "GET /images/feed-icon-14x14.png HTTP/1.1" 301 201 "http://somereferer.ch/" "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.99 Safari/537.36" "-"
var logre = regexp.MustCompile(`(.+)\s.+\s-\s\[(.+)\]\s"(GET|POST|PUT|HEAD|OPTIONS|DELETE).+HTTP/1\.[0-1]"\s[0-9]{3}\s([0-9]{1,20})\s".+`)

const logDateFormat string = "2/Jan/2006:15:04:05 -0700"

const secondsInMonth float64 = 2628000

func processFile(r io.Reader, stats map[string]month, next chan bool) {
	defer func() {
		next <- true
	}()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fields := logre.FindStringSubmatch(scanner.Text())
		if fields == nil {
			continue
		}
		if len(fields) < 5 {
			fmt.Fprintf(os.Stderr, "found wrong number of fields on line %s\n", scanner.Text())
			continue
		}
		var ip, date, method, respsize string = fields[1], fields[2], fields[3], fields[4]
		if ip == "" || date == "" || method == "" || respsize == "" {
			fmt.Fprintf(os.Stderr, "found empty fields '%s', '%s', '%s', '%s'\n",
				ip, date, method, respsize)
			continue
		}
		bandwidth, err := strconv.ParseFloat(respsize, 64)
		t, err := time.Parse(logDateFormat, date)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse time %s\n", date)
			continue
		}
		curmonth := t.UTC().Format("2006-01")
		if _, ok := stats[curmonth]; !ok {
			stats[curmonth] = month{}
		}
		statmonth := stats[curmonth]
		statmonth.Hits.Bandwidth += bandwidth
		statmonth.Hits.TotalHits += 1
		statmonth.Hits.Qps = statmonth.Hits.TotalHits / secondsInMonth
		switch method {
		case "GET":
			statmonth.Hits.TotalGet += 1
		case "POST":
			statmonth.Hits.TotalPost += 1
		}
		switch getContinent(ip) {
		case "North America":
			statmonth.RegionalHits.NorthAmerica.Bandwidth += bandwidth
			statmonth.RegionalHits.NorthAmerica.TotalHits += 1
			statmonth.RegionalHits.NorthAmerica.Qps = statmonth.RegionalHits.NorthAmerica.TotalHits / secondsInMonth
			switch method {
			case "GET":
				statmonth.RegionalHits.NorthAmerica.TotalGet += 1
			case "POST":
				statmonth.RegionalHits.NorthAmerica.TotalPost += 1
			}
		case "South America":
			statmonth.RegionalHits.SouthAmerica.Bandwidth += bandwidth
			statmonth.RegionalHits.SouthAmerica.TotalHits += 1
			statmonth.RegionalHits.SouthAmerica.Qps = statmonth.RegionalHits.SouthAmerica.TotalHits / secondsInMonth
			switch method {
			case "GET":
				statmonth.RegionalHits.SouthAmerica.TotalGet += 1
			case "POST":
				statmonth.RegionalHits.SouthAmerica.TotalPost += 1
			}
		case "Europe":
			statmonth.RegionalHits.Europe.Bandwidth += bandwidth
			statmonth.RegionalHits.Europe.TotalHits += 1
			statmonth.RegionalHits.Europe.Qps = statmonth.RegionalHits.Europe.TotalHits / secondsInMonth
			switch method {
			case "GET":
				statmonth.RegionalHits.Europe.TotalGet += 1
			case "POST":
				statmonth.RegionalHits.Europe.TotalPost += 1
			}
		case "Asia":
			statmonth.RegionalHits.Asia.Bandwidth += bandwidth
			statmonth.RegionalHits.Asia.TotalHits += 1
			statmonth.RegionalHits.Asia.Qps = statmonth.RegionalHits.Asia.TotalHits / secondsInMonth
			switch method {
			case "GET":
				statmonth.RegionalHits.Asia.TotalGet += 1
			case "POST":
				statmonth.RegionalHits.Asia.TotalPost += 1
			}
		case "Oceania":
			statmonth.RegionalHits.Oceania.Bandwidth += bandwidth
			statmonth.RegionalHits.Oceania.TotalHits += 1
			statmonth.RegionalHits.Oceania.Qps = statmonth.RegionalHits.Oceania.TotalHits / secondsInMonth
			switch method {
			case "GET":
				statmonth.RegionalHits.Oceania.TotalGet += 1
			case "POST":
				statmonth.RegionalHits.Oceania.TotalPost += 1
			}
		case "Africa":
			statmonth.RegionalHits.Africa.Bandwidth += bandwidth
			statmonth.RegionalHits.Africa.TotalHits += 1
			statmonth.RegionalHits.Africa.Qps = statmonth.RegionalHits.Africa.TotalHits / secondsInMonth
			switch method {
			case "GET":
				statmonth.RegionalHits.Africa.TotalGet += 1
			case "POST":
				statmonth.RegionalHits.Africa.TotalPost += 1
			}
		}
		stats[curmonth] = statmonth
	}
}

func getContinent(ip string) string {
	record, err := maxmind.City(net.ParseIP(ip))
	if err != nil {
		return ""
	}
	return record.Continent.Names["en"]
}
