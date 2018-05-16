package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
	yaml "gopkg.in/yaml.v2"

	_ "github.com/victorgama/howe/widgets/handlers/banner"
	_ "github.com/victorgama/howe/widgets/handlers/blank"
	_ "github.com/victorgama/howe/widgets/handlers/disks"
	_ "github.com/victorgama/howe/widgets/handlers/docker"
	_ "github.com/victorgama/howe/widgets/handlers/load"
	_ "github.com/victorgama/howe/widgets/handlers/print"
	_ "github.com/victorgama/howe/widgets/handlers/systemd-services"
	_ "github.com/victorgama/howe/widgets/handlers/updates"
	_ "github.com/victorgama/howe/widgets/handlers/uptime"

	"github.com/victorgama/howe/config"
	"github.com/victorgama/howe/widgets"
)

const configPath = "/etc/howe/config.yml"

func main() {
	config := config.Root{}
	color.NoColor = false

	_, err := os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("howe: %s not found. Please refer to the documentation.", configPath)
		} else {
			log.Fatalf("howe: error: %s", err)
		}
	}

	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("howe: error: %s", err)
	}

	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		log.Fatalf("howe: error: %v", err)
	}

	wg := sync.WaitGroup{}
	result := []chan interface{}{}
	for i, w := range config.Items {
		rawName, ok := w["type"]
		if !ok {
			log.Fatalf("howe: error: widget %d is missing a type attribute", i+1)
		}

		name, ok := rawName.(string)
		if !ok {
			log.Fatalf("howe: error: widget %d is has an invalid type attribute", i+1)
		}

		handler, ok := widgets.Handlers[name]
		if !ok {
			log.Fatalf("howe: error: widget %d uses unknown type %s", i+1, name)
		}

		wg.Add(1)
		output := make(chan interface{}, 1)
		result = append(result, output)
		go handler(w, output, &wg)
	}

	wg.Wait()

	for _, ch := range result {
		res := <-ch
		if err, ok := res.(error); ok {
			log.Fatalf("howe: error: %s", err)
		} else {
			fmt.Println(strings.Trim(res.(string), "\n"))
		}
	}
}
