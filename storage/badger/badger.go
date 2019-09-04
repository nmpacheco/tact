package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/brunotm/tact/storage/badgerdb"
)

var (
	currentPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	dataPath       = flag.String("datapath", "../../daemon/tactd/statedb", "Path for state data")
	prefix         = flag.String("prefix", "delta", "Storage prefix")
)

func main() {
	flag.Parse()
	fmt.Printf("Args:\ndatapath = %s\nprefix = %s\n", *dataPath, *prefix)

	s, err := badgerdb.Open(*dataPath, true)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	tx := s.NewTxn(false)
	entries, err := tx.GetTree([]byte(*prefix))
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		fmt.Printf("key: %s\nvalue: %s\n", entry.Key, entry.Value)
	}
}
