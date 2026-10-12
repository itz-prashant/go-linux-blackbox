package collector

import (
	"bufio"
	"fmt"
	"os"
)

type rawCPUticks struct {
	user    uint64 // run normal application
	nice    uint64 // low priority process
	system  uint64 // linux kernel syscall
	idle    uint64 // when cpu no use and on rest mode
	iowait  uint64 // cpu wait for hhd/ssd data
	irq     uint64 // handle hardware interrupt
	softirq uint64 // software interrupt
	steal   uint64 // how much cpu steal by hypervisor on VM
}

func (t *rawCPUticks) Total() uint64 {
	total := t.user + t.nice + t.system + t.idle + t.iowait + t.irq + t.softirq + t.steal

	return total
}

func (t *rawCPUticks) Idle() uint64 {
	return t.idle + t.iowait
}

func readRawCPUTicks(statPath string) (*rawCPUticks, error) {
	file, err := os.Open(statPath)

	if err != nil {
		return nil, fmt.Errorf("failed to open stat file: %w", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		err := scanner.Err()
		return nil, err
	}

	line := scanner.Text()

	var ticks rawCPUticks

	_, err = fmt.Sscanf(line, "cpu %d %d %d %d %d %d %d %d", &ticks.user, &ticks.nice, &ticks.system, &ticks.idle, &ticks.iowait, &ticks.irq, &ticks.softirq, &ticks.steal)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cpu line: %w", err)
	}

	return &ticks, nil
}