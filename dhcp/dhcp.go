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
}

// Parse 解析选项区。
func Parse(data []byte, c *Counter) ([]Option, error) {
	if len(data) < HeaderLen+MagicLen {
		return nil, errors.New("too-short")
	}
	off := HeaderLen + MagicLen
	out := []Option{}
	seen := map[byte]bool{}
	for off+1 < len(data) {
		code := data[off]
		length := int(data[off+1])
		if off+2+length > len(data) {
			return nil, errors.New("bad-option")
		}
		c.Scanned++
		if seen[code] {
			off += 2 + length
			continue
		}
		seen[code] = true
		out = append(out, Option{Code: code, Len: length, Data: data[off+2 : off+2+length]})
		off += 2 + length
	}
	return out, nil
}

// Lookup 按 code 找选项。
func Lookup(data []byte, code byte, c *Counter) ([]Option, error) {
	options, err := Parse(data, c)
	if err != nil {
		return nil, err
	}
	found := []Option{}
	for _, opt := range options {
		if opt.Code == code {
			found = append(found, opt)
		}
	}
	return found, nil
}
