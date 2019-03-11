package main

import (
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var (
	singleFileSizeMax int64
	totalFiles        int
	bpfFilter         string
	name              string
	ni                string
	nn                string
	packetSize        int
)

func main() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		fmt.Println(err)
		return
	}

	flag.Int64Var(&singleFileSizeMax, "Z", 50, "use to specify single File size Max MB")
	flag.IntVar(&totalFiles, "J", 5, "use to specify total files")
	flag.StringVar(&bpfFilter, "BPF", "", "port list")
	flag.StringVar(&name, "Name", "Test", "port list")
	flag.StringVar(&ni, "NI", "", "network card ip")
	flag.StringVar(&nn, "NN", "", "network card name")
	flag.IntVar(&packetSize, "ps", 1024, "packet size unit is byte")
	flag.Parse()

	fmt.Println("flag parse finished.")

	finalChannel := make(chan gopacket.Packet, 5000)
	go trueWrite(finalChannel)

	if ni != "" {
		//this is ip mode
		networkInterfaceViaIp := findNetworkInterfaceViaIp(ni, devices)
		go goCaptureDevice(networkInterfaceViaIp, finalChannel)
	} else if nn != "any" {
		go goCaptureDevice(nn, finalChannel)
	} else {
		for _, device := range devices {
			go goCaptureDevice(device.Name, finalChannel)
		}
	}

	ch := make(chan int)
	_ = <-ch
	os.Exit(0)
}

func goCaptureDevice(deviceName string, finalChannel chan gopacket.Packet) {
	handle, err := pcap.OpenLive(deviceName, int32(packetSize), false, 0)
	if bpfFilter != "" {
		handle.SetBPFFilter(bpfFilter)
	}
	if err != nil {
		fmt.Printf("Error opening device %s: %v", deviceName, err)
		os.Exit(1)
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		finalChannel <- packet
	}
}

func trueWrite(packets chan gopacket.Packet) {
	var isStart = true
	var nowFileIndex = 0
	var currentFile *os.File
	var currentWriter *pcapgo.Writer
	for packet := range packets {
		if isStart {
			currentFile, currentWriter = getFileWriter(getPcapWriteName(nowFileIndex))
			currentWriter.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
			isStart = false
			continue
		}
		//fmt.Println(packet)
		currentWriter.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		size := getCurrentFileSize(nowFileIndex)
		if size > singleFileSizeMax*1000*1000 {
			closeWriter(currentFile, currentWriter)
			judgeAndRemovePcap()
			nowFileIndex++
			isStart = true
		}
	}
}

func closeWriter(file *os.File, writer *pcapgo.Writer) {
	var aux_value = file.Name()
	file.Close()
	zipFile(aux_value)
}

func zipFile(fileName string) {
	ZipFiles(fileName+".zip", fileName)
	os.Remove(fileName)
}

func getCurrentFileSize(index int) int64 {
	var aux_value = getPcapWriteName(index)
	info, _ := os.Stat(aux_value)
	return info.Size()
}

func getPcapWriteName(index int) string {
	itoa := strconv.Itoa(index)
	var midZero = ""
	for i := 0; i < 4-len(itoa); i++ {
		midZero += "0"
	}
	return name + midZero + itoa + ".pcap"
}

func judgeAndRemovePcap() {
	var currentPcapNum = getPcapCountCurrentDir()
	if currentPcapNum >= totalFiles {
		deleteOldestPcap()
	}
}

func deleteOldestPcap() {
	var oldFile os.FileInfo
	files, _ := ioutil.ReadDir("./")
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".pcap.zip") {
			if oldFile == nil {
				oldFile = file
			} else {
				if file.ModTime().Before(oldFile.ModTime()) {
					oldFile = file
				}
			}
			file.ModTime()
		}
	}
	err := os.Remove(oldFile.Name())
	if err != nil {
		fmt.Println("delet file error, error is %v", err)
	}
}

func getPcapCountCurrentDir() int {
	var count = 0
	files, _ := ioutil.ReadDir("./")
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".pcap.zip") {
			count++
		}
	}
	return count
}

func getFileWriter(name string) (*os.File, *pcapgo.Writer) {
	f, _ := os.Create(name)
	w := pcapgo.NewWriter(f)

	w.WriteFileHeader(0, layers.LinkTypeEthernet)
	return f, w
}

func findNetworkInterfaceViaIp(ip string, devices []pcap.Interface) string {
	fmt.Println("begin to traversing Devices found:")
	for _, device := range devices {
		for _, address := range device.Addresses {
			fmt.Println("- IP address: ", address.IP)
			if address.IP.String() == ip {
				return device.Name
			}
		}
	}
	return ""
}
