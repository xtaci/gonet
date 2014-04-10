package cfg

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

var _DEF_CONFIG = os.Getenv("GOPATH") + "/config.ini"

var (
	_map  map[string]string
	_lock sync.Mutex
)

func init() {
	Reload()
}

func Get() map[string]string {
	_lock.Lock()
	defer _lock.Unlock()
	return _map
}

func Reload() {
	var path string
	if path = os.Getenv("GONET_CONFIG"); path == "" {
		path = _DEF_CONFIG
	}

	_lock.Lock()
	log.Println("Loading Config from:", path)
	defer log.Println("Config Loaded.")
	_map = _load_config(path)
	_lock.Unlock()
}

func _load_config(path string) (ret map[string]string) {
	ret = make(map[string]string)
	f, err := os.Open(path)
	if err != nil {
		log.Println(path, err)
		return
	}

	re := regexp.MustCompile(`[\t ]*([0-9A-Za-z_]+)[\t ]*=[\t ]*([^\t\n\f\r# ]+)[\t #]*`)

	// using scanner to read config file
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// expression match
		slice := re.FindStringSubmatch(line)

		if slice != nil {
			ret[slice[1]] = slice[2]
			log.Println(slice[1], "=", slice[2])
		}
	}

	return
}
