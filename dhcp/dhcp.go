// Package dhcp 解析 DHCPv4 选项区。
package dhcp

import "errors"

// 报文字节布局常量。
const (
	HeaderLen = 236
	MagicLen  = 4
	EndCode   = 255
	PadCode   = 0
)

// Option 是一条选项。
type Option struct {
	Code byte   `json:"code"`
	Len  int    `json:"len"`
	Data []byte `json:"-"`
}

// Counter 记录本轮检查过的选项条数。
type Counter struct {
	Scanned int
	cache   map[string]parsedPacket
}

type parsedPacket struct {
	options []Option
	byCode  map[byte][]Option
}

// Parse 解析选项区。
func Parse(data []byte, c *Counter) ([]Option, error) {
	out, _, err := parseOptions(data, c)
	return out, err
}

// Lookup 按 code 找选项。
func Lookup(data []byte, code byte, c *Counter) ([]Option, error) {
	_, byCode, err := parseOptions(data, c)
	if err != nil {
		return nil, err
	}
	return byCode[code], nil
}

func parseOptions(data []byte, c *Counter) ([]Option, map[byte][]Option, error) {
	if len(data) < HeaderLen+MagicLen {
		return nil, nil, errors.New("too-short")
	}
	if c != nil {
		if cached, ok := c.cache[string(data)]; ok {
			return cached.options, cached.byCode, nil
		}
	}

	off := HeaderLen + MagicLen
	out := []Option{}
	byCode := map[byte][]Option{}
	for off < len(data) {
		code := data[off]
		off++
		if code == EndCode {
			break
		}
		if code == PadCode {
			continue
		}
		if off >= len(data) {
			return nil, nil, errors.New("bad-option")
		}
		length := int(data[off])
		off++
		if off+length > len(data) {
			return nil, nil, errors.New("bad-option")
		}
		if c != nil {
			c.Scanned++
		}
		opt := Option{Code: code, Len: length, Data: data[off : off+length]}
		out = append(out, opt)
		byCode[code] = append(byCode[code], opt)
		off += length
	}

	if c != nil {
		if c.cache == nil {
			c.cache = map[string]parsedPacket{}
		}
		c.cache[string(data)] = parsedPacket{options: out, byCode: byCode}
	}
	return out, byCode, nil
}
