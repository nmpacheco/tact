package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"strings"

	"github.com/brunotm/tact"
	_ "github.com/brunotm/tact/collector/aix"
	_ "github.com/brunotm/tact/collector/linux"

	//	_ "github.com/brunotm/tact/collector/oracle"
	"github.com/brunotm/tact/log"
	"github.com/brunotm/tact/scheduler"
)

var (
	currentPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	cron           = flag.String("cron", "0 */1 * * * *", "Cron like scheduling expression: 0 */1 * * * *")
	nodeCfgFile    = flag.String("nodecfg", fmt.Sprintf("%s/nodes.json", currentPath), "Nodes json config file")
	collector      = flag.String("c", "", "Collector or group to run")
	logLevel       = flag.String("log", "info", "Log level")
	dataPath       = flag.String("datapath", "./statedb", "Path for state data")
)

func nodesFromFile(path string) (*tact.Nodes, error) {
	var file []byte
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	nodes := &tact.Nodes{}
	err = json.Unmarshal(file, nodes)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func addJobs(cron string, sched *scheduler.Scheduler, nodes *tact.Nodes) error {
	for _, node := range nodes.Nodes {
		log.Info("%v", node)
		if node.NetAddr == "" {
			node.NetAddr = node.HostName
		}

		if node.Collectors == nil {
			return fmt.Errorf("no colector specified")
		}

		for _, collector := range node.Collectors {
			var coll *tact.Collector
			var collGroup []*tact.Collector
			if len(strings.Split(collector, "/")) > 3 {
				coll = tact.Registry.Get(collector)
			} else {
				collGroup = tact.Registry.GetGroup(collector)
			}
			if collGroup != nil {
				for _, c := range collGroup {
					if err := sched.AddJob(cron, c, &node, 290*time.Second); err != nil {
						return err
					}
				}
			}
			if coll != nil {
				if err := sched.AddJob(cron, coll, &node, 290*time.Second); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func main() {
	// grmon.Start()
	flag.Parse()
	var err error

	switch strings.ToLower(*logLevel) {
	case "debug":
		log.SetDebug()
	case "info":
		log.SetInfo()
	case "warn":
		log.SetWarn()
	case "error":
		log.SetError()
	default:
		log.Error("invalid log level")
		os.Exit(1)
	}

	tact.Init(*dataPath)

	nodes, err := nodesFromFile(*nodeCfgFile)
	if err != nil {
		panic(err)
	}

	wchan := make(chan []byte, 10)
	go func() {
		for e := range wchan {
			fmt.Println(string(e))
		}
	}()

	sched := scheduler.New(100, 60*time.Second, wchan)
	if err := addJobs(*cron, sched, nodes); err != nil {
		panic(err)
	}

	sched.Start()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	sched.Stop()

	log.Info("Shutting down")
	tact.Close()
}
