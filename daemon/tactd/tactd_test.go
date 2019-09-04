package main

import (
	"testing"
	"time"

	"github.com/brunotm/tact"
	_ "github.com/brunotm/tact/collector/aix"
	_ "github.com/brunotm/tact/collector/linux"
	"github.com/brunotm/tact/scheduler"
)

func TestNodesFromFile(t *testing.T) {
	nodes := &tact.Nodes{}
	nodes, err := nodesFromFile("./nodes1.json")
	if err != nil {
		t.Fatalf("%v", err)
	}
	for key, node := range nodes.Nodes {
		t.Logf("Nodes %s %v", key, node)
		t.Logf("Nodes %v", node.SSHKey)
		t.Logf("Nodes %v", node.LogFiles)
	}
}

func TestAddJobs(t *testing.T) {
	wchan := make(chan []byte)
	sched := scheduler.New(100, 60*time.Second, wchan)
	nodes := &tact.Nodes{}
	nodes, err := nodesFromFile("./nodes1.json")
	if err != nil {
		t.Fatalf("%v", err)
	}

	if err := addJobs("0/30 * * * * *", sched, nodes); err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("Listing crons")
	for _, cron := range sched.GetCrons() {
		t.Logf("%v", cron)
	}
}
