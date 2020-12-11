package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/egeoz/iterminal"
)

var commandLineOptions struct {
	hideCPU    bool
	hideRAM    bool
	hideDisk   bool
	hideUptime bool
	hideKernel bool
	hideOS     bool
	hideShell  bool
	hideDE     bool
	hideGPU    bool

	showHelp bool
	showVer  bool
}

func help() {
	fmt.Println("Gofetch is a command line system information tool.\n\n",
		"Optional Flags:\n",
		"	-h\tPrint this help message\n",
		"	-v\tPrint version info\n",
		"	-cpu\tHide CPU info\n",
		"	-ram\tHide RAM info\n",
		"	-disk\tHide disk info\n",
		"	-gpu\tHide GPU info\n",
		"	-kernel\tHide kernel version\n",
		"	-shell\tHide shell name\n",
		"	-uptime\tHide uptime\n",
		"	-de\tHide desktop enviroment info\n",
	)
}

func ver() {
	fmt.Println("Gofetch 1.0")
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func genInfo(hideCPU bool, hideRAM bool, hideDisk bool, hideGPU bool, hideKernel bool, hideShell bool, hideUptime bool, hideDE bool) {
	var data, desc [9]string
	i := 0

	f, err := os.Open("/etc/hostname")
	var temp []string

	checkError(err)
	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)

	for s.Scan() {
		temp = append(temp, s.Text())
	}

	f.Close()
	data[i] = os.Getenv("USER") + "@" + temp[0]
	i++

	if !hideKernel {
		out, _ := exec.Command("uname", "-r").Output()

		data[i] = string(out)[:len(string(out))-1]
		desc[i] = "Kernel:"
		i++
	}

	if !hideCPU {
		f, err := os.Open("/proc/cpuinfo")
		checkError(err)

		s := bufio.NewScanner(f)
		line := 0

		for s.Scan() {
			line++
			if line == 5 {
				data[i] = s.Text()[13:]
			}
		}
		f.Close()

		desc[i] = "CPU:"
		i++
	}

	if !hideGPU {
		var temp []string
		var l int
		out, _ := exec.Command("lspci", "-mm").Output()
		temp = strings.Split(string(out), "\"")

		for i, str := range temp {
			if str == "VGA compatible controller" {
				l = i
			}
		}

		data[i] = temp[l+4]
		desc[i] = "GPU:"
		i++
	}

	if !hideDisk {
		var temp []string
		var l int
		out, _ := exec.Command("df", "-h").Output()
		temp = strings.Fields(string(out))

		for i, str := range temp {
			if str == "/" {
				l = i
			}
		}

		data[i] = temp[l-2] + " / " + temp[l-3]
		desc[i] = "Disk:"
		i++
	}

	if !hideRAM {
		var MemUsed, MemTotal, Shmem, MemFree, Buffers, Cached, SReclaimable int
		var temp []string

		f, err := os.Open("/proc/meminfo")
		checkError(err)

		s := bufio.NewScanner(f)
		line := 0

		for s.Scan() {
			line++

			if line == 1 {
				temp = strings.Fields(s.Text())
				MemTotal, _ = strconv.Atoi(temp[1])
			} else if line == 2 {
				temp = strings.Fields(s.Text())
				MemFree, _ = strconv.Atoi(temp[1])
			} else if line == 4 {
				temp = strings.Fields(s.Text())
				Buffers, _ = strconv.Atoi(temp[1])
			} else if line == 5 {
				temp = strings.Fields(s.Text())
				Cached, _ = strconv.Atoi(temp[1])
			} else if line == 21 {
				temp = strings.Fields(s.Text())
				Shmem, _ = strconv.Atoi(temp[1])
			} else if line == 24 {
				temp = strings.Fields(s.Text())
				SReclaimable, _ = strconv.Atoi(temp[1])
			}
		}
		f.Close()
		MemUsed = (MemTotal + Shmem - MemFree - Buffers - Cached - SReclaimable) / 1024
		strUsed := strconv.Itoa(MemUsed)
		MemTotal /= 1024
		strTotal := strconv.Itoa(MemTotal)

		data[i] = strUsed + " MB / " + strTotal + " MB"
		desc[i] = "RAM:"
		i++
	}

	if !hideShell {
		data[i] = os.Getenv("SHELL")
		desc[i] = "Shell:"
		i++
	}

	if !hideDE {
		data[i] = os.Getenv("XDG_CURRENT_DESKTOP")
		desc[i] = "DE:"
		i++
	}

	if !hideUptime {
		var temp []string
		var fsec float64
		var day, hour, min, sec int

		f, err := os.Open("/proc/uptime")

		checkError(err)
		s := bufio.NewScanner(f)
		s.Split(bufio.ScanLines)

		for s.Scan() {
			temp = append(temp, s.Text())
		}

		f.Close()

		fsec, _ = strconv.ParseFloat(strings.Fields(temp[0])[0], 64)
		sec = int(fsec)

		day = sec / 60 / 60 / 24
		hour = sec / 60 / 60 % 24
		min = sec / 60 % 60

		data[i] = strconv.Itoa(int(day)) + " days, " + strconv.Itoa(int(hour)) + " hours, " + strconv.Itoa(int(min)) + " minutes"
		desc[i] = "Uptime:"
		i++
	}

	fmt.Printf("\t\t%s%s%s%s\n", iterminal.MakeColor(iterminal.ForegroundColor.LightBlue, ""), iterminal.MakeFont(iterminal.Font.Bold), data[0], iterminal.MakeFont(iterminal.Font.Default))
	fmt.Printf("      ___\t%s%s%s%s%s\t%s\n", iterminal.MakeFont(iterminal.Font.Bold), iterminal.MakeColor(iterminal.ForegroundColor.Yellow, ""), desc[1], iterminal.MakeColor(iterminal.ForegroundColor.Default, ""), iterminal.MakeFont(iterminal.Font.Default), data[1])
	fmt.Printf("     (.. |\t%s%s%s%s%s\t%s\n", iterminal.MakeFont(iterminal.Font.Bold), iterminal.MakeColor(iterminal.ForegroundColor.Yellow, ""), desc[2], iterminal.MakeColor(iterminal.ForegroundColor.Default, ""), iterminal.MakeFont(iterminal.Font.Default), data[2])
	fmt.Printf("     (%s<>%s |\t%s%s%s%s%s\t%s\n", iterminal.MakeColor(iterminal.ForegroundColor.Yellow, ""), iterminal.MakeColor(iterminal.ForegroundColor.Default, ""), iterminal.MakeFont(iterminal.Font.Bold), iterminal.MakeColor(iterminal.ForegroundColor.Yellow, ""), desc[3], iterminal.MakeColor(iterminal.ForegroundColor.Default, ""), iterminal.MakeFont(iterminal.Font.Default), data[3])
	fmt.Printf("    / __  \\\t%s%s%s%s%s\t%s\n", iterminal.MakeFont(iterminal.Font.Bold), iterminal.MakeColor(iterminal.ForegroundColor.Yellow, ""), desc[4], iterminal.MakeColor(iterminal.ForegroundColor.Default, ""), iterminal.MakeFont(iterminal.Font.Default), data[4])
	fmt.Printf("   ( /  \\ /|\t%s%s%s%s%s\t%s\n", iterminal.MakeFont(iterminal.Font.Bold), iterminal.MakeColor(iterminal.ForegroundColor.Yellow, ""), desc[5], iterminal.MakeColor(iterminal.ForegroundColor.Default, ""), iterminal.MakeFont(iterminal.Font.Default), data[5])
	fmt.Printf("  _/\\ __)/_)\t%s%s%s%s%s\t%s\n", iterminal.MakeFont(iterminal.Font.Bold), iterminal.MakeColor(iterminal.ForegroundColor.Yellow, ""), desc[6], iterminal.MakeColor(iterminal.ForegroundColor.Default, ""), iterminal.MakeFont(iterminal.Font.Default), data[6])
	fmt.Printf(" %s\\|/%s-___%s\\|/%s\t%s%s%s%s%s\t%s\n", iterminal.MakeColor(iterminal.ForegroundColor.Yellow, ""), iterminal.MakeColor(iterminal.ForegroundColor.Default, ""), iterminal.MakeColor(iterminal.ForegroundColor.Yellow, ""), iterminal.MakeColor(iterminal.ForegroundColor.Default, ""), iterminal.MakeFont(iterminal.Font.Bold), iterminal.MakeColor(iterminal.ForegroundColor.Yellow, ""), desc[7], iterminal.MakeColor(iterminal.ForegroundColor.Default, ""), iterminal.MakeFont(iterminal.Font.Default), data[7])
	fmt.Printf("\t\t%s%s%s%s%s\t%s\n", iterminal.MakeFont(iterminal.Font.Bold), iterminal.MakeColor(iterminal.ForegroundColor.Yellow, ""), desc[8], iterminal.MakeColor(iterminal.ForegroundColor.Default, ""), iterminal.MakeFont(iterminal.Font.Default), data[8])
	fmt.Println()
}

func main() {
	flag.BoolVar(&commandLineOptions.showHelp, "h", false, "Print this help message")
	flag.BoolVar(&commandLineOptions.showVer, "v", false, "Print version info")
	flag.BoolVar(&commandLineOptions.hideCPU, "cpu", false, "Hide CPU info")
	flag.BoolVar(&commandLineOptions.hideRAM, "ram", false, "Hide RAM info")
	flag.BoolVar(&commandLineOptions.hideDisk, "disk", false, "Hide disk info")
	flag.BoolVar(&commandLineOptions.hideGPU, "gpu", false, "Hide GPU info")
	flag.BoolVar(&commandLineOptions.hideKernel, "kernel", false, "Hide kernel version")
	flag.BoolVar(&commandLineOptions.hideShell, "shell", false, "Hide shell name")
	flag.BoolVar(&commandLineOptions.hideUptime, "uptime", false, "Hide uptime")
	flag.BoolVar(&commandLineOptions.hideDE, "de", false, "Hide desktop enviroment info")
	flag.Parse()

	switch true {
	case commandLineOptions.showHelp:
		help()

	case commandLineOptions.showVer:
		ver()
	default:
		genInfo(commandLineOptions.hideCPU, commandLineOptions.hideRAM, commandLineOptions.hideDisk, commandLineOptions.hideGPU, commandLineOptions.hideKernel, commandLineOptions.hideShell, commandLineOptions.hideUptime, commandLineOptions.hideDE)
	}

}
