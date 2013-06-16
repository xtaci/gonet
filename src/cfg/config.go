package cfg

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

var _DEF_CONFIG = os.Getenv("GOPATH") + "/config.ini"

var (
	_map        map[string]string
	_lock       sync.RWMutex
	config_file = flag.String("config", _DEF_CONFIG, "specify absolute path for config.ini")
)

func init() {
	Reload()
}

func Get() map[string]string {
	_lock.RLock()
	defer _lock.RUnlock()
	return _map
}

func Reload() {
	path := *config_file
	log.Println("Loading Config.")
	defer log.Println("Config Loaded.")
	_lock.Lock()
	_map = _load_config(path)
	_lock.Unlock()
}

func _load_config(path string) (ret map[string]string) {
	ret = make(map[string]string)
	f, err := os.Open(path)

	if err != nil {
		fmt.Println("error opening file %v\n", err)
		os.Exit(-1)
	}

	re := regexp.MustCompile(`[\t ]*([0-9A-Za-z_]+)[\t ]*=[\t ]*([^\t\n\f\r# ]+)[\t #]*`)

	r := bufio.NewReader(f)

	for {
		line, e := r.ReadString('\n')
		line = strings.TrimSpace(line)

		// empty-line & #comment
		if line == "" {
			if e == nil {
				continue
			} else {
				break
			}
		}

		if []byte(line)[0] == '#' {
			continue
		}

		// maping
		slice := re.FindStringSubmatch(line)

		if slice != nil {
			ret[slice[1]] = slice[2]
		}

		if e != nil {
			break
		}
	}

	return
}
