package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	None = iota
	PDetected
	PSDetected

	P = "pulse"
	S = "space"
)

var (
	T             float64
	BitDump       bool
	HexDump       bool
	DisplayChange bool
	SkipFrameHex  string
	FrameLength   int
)

func init() {
	flag.Float64Var(&T, "T", 445, "t length")
	flag.BoolVar(&HexDump, "h", true, "hex dump")
	flag.BoolVar(&BitDump, "b", false, "bit dump")
	flag.BoolVar(&DisplayChange, "c", false, "display changes")
	flag.StringVar(&SkipFrameHex, "skip", "", "skip frame hex")
	flag.IntVar(&FrameLength, "l", -1, "valid frame length")
	flag.Parse()
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var (
		dc        int
		datum     byte
		prevFrame []byte
		frame     []byte
		pf        = false
		lf        = None
	)

	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " ")
		if len(s) == 2 {
			state := s[0]
			nsec, err := strconv.Atoi(s[1])
			if err != nil {
				log.Fatal(err)
			}

			l := int(math.Round(float64(nsec) / T))

			switch {
			case state == P && l == 8:
				lf = PDetected
			case state == S && l == 4 && lf == PDetected:
				lf = PSDetected
				dc = 0
				pf = false
				datum = 0
				frame = nil
			}

			if lf == PSDetected {
				switch {
				case state == P && l == 1:
					dc++
					pf = true
				case state == S && pf:
					pf = false

					var bit byte = 0
					switch l {
					case 1:
						bit = 0
					case 3:
						bit = 1
					default:
						if len(frame) > 0 {
							if !IsSkipped(frame) {
								printBytes(frame, prevFrame)
								prevFrame = frame
							}
							frame = nil
						}
					}

					datum = datum<<1 | bit

					if dc == 20 || dc == 24 {
						frame = append(frame, reverseBit4(datum))
						datum = 0
					} else if dc%8 == 0 {
						frame = append(frame, reverseBit8(datum))
						datum = 0
					}
				}
			}
		}
	}
}

func reverseBit8(x byte) byte {
	x = ((x & 0x55) << 1) | ((x & 0xAA) >> 1)
	x = ((x & 0x33) << 2) | ((x & 0xCC) >> 2)
	return (x&0x0F)<<4 | x>>4
}

func reverseBit4(x byte) byte {
	x = ((x & 0x5) << 1) | ((x & 0xA) >> 1)
	x = ((x & 0x3) << 2) | ((x & 0xC) >> 2)
	return x
}

func printBytes(bytes []byte, prevBytes []byte) {
	if HexDump {
		if len(prevBytes) == len(bytes) && DisplayChange {
			for i, v := range bytes {
				if v != prevBytes[i] {
					fmt.Print(color.YellowString("%02X", v))
				} else {
					fmt.Printf("%02X", v)
				}
				fmt.Print(" ")
			}
			fmt.Println()
		} else {
			for _, v := range bytes {
				fmt.Printf("%02X ", v)
			}
			fmt.Println()
		}
	}
	if BitDump {
		if len(prevBytes) == len(bytes) && DisplayChange {
			for i, v := range bytes {
				if v != prevBytes[i] {
					a := fmt.Sprintf("%08b", v)
					b := fmt.Sprintf("%08b", prevBytes[i])
					for k := range a {
						if a[k] != b[k] {
							fmt.Print(color.YellowString("%c", a[k]))
						} else {
							fmt.Printf("%c", a[k])
						}
					}
				} else {
					fmt.Printf("%08b", v)
				}
				fmt.Print(" ")
			}
			fmt.Println()
		} else {
			for _, v := range bytes {
				fmt.Printf("%08b ", v)
			}
			fmt.Println()
		}
	}
}

func IsSkipped(bytes []byte) bool {
	return (FrameLength > 0 && len(bytes) != FrameLength) || fmt.Sprintf("%X", bytes) == SkipFrameHex
}
