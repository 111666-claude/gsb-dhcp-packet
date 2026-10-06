// Command dhcpdump 是选项解析的场景入口。
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"example.com/dhcp-packet/dhcp"
)

func emit(v any) {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func packet(options []byte) []byte {
	data := make([]byte, dhcp.HeaderLen+dhcp.MagicLen)
	copy(data[236:], []byte{0x63, 0x82, 0x53, 0x63})
	return append(data, options...)
}

func main() {
	sample := flag.String("sample", "", "pad|end|repeat|work")
	flag.Parse()

	c := &dhcp.Counter{}
	switch *sample {
	case "pad":
		got, err := dhcp.Parse(packet([]byte{0, 0, 53, 1, 2}), c)
		if err != nil {
			emit(map[string]string{"error": err.Error()})
			return
		}
		emit(map[string]any{"codes": codes(got)})
	case "end":
		got, err := dhcp.Parse(packet([]byte{53, 1, 2, 255, 0, 12, 1, 'x'}), c)
		if err != nil {
			emit(map[string]string{"error": err.Error()})
			return
		}
		emit(map[string]any{"codes": codes(got)})
	case "repeat":
		got, err := dhcp.Parse(packet([]byte{53, 1, 2, 53, 1, 3}), c)
		if err != nil {
			emit(map[string]string{"error": err.Error()})
			return
		}
		emit(map[string]int{"codes": len(got)})
	case "work":
		options := []byte{}
		for i := 0; i < 2000; i++ {
			options = append(options, 53, 1, byte(i))
		}
		data := packet(options)
		total := 0
		for i := 0; i < 2000; i++ {
			found, err := dhcp.Lookup(data, 53, c)
			if err != nil {
				panic(err)
			}
			total += len(found)
		}
		emit(map[string]int{"options": 2000, "scanned": c.Scanned})
	default:
		fmt.Fprintln(os.Stderr, "未知场景："+*sample)
		os.Exit(2)
	}
}

func codes(options []dhcp.Option) []int {
	out := []int{}
	for _, opt := range options {
		out = append(out, int(opt.Code))
	}
	return out
}
