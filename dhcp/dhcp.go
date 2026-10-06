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

	src    []byte
	parsed []Option
}

// Parse 解析选项区。
func Parse(data []byte, c *Counter) ([]Option, error) {
	if len(data) < HeaderLen+MagicLen {
		return nil, errors.New("too-short")
	}
	if c.sameSource(data) {
		return c.parsed, nil
	}
	off := HeaderLen + MagicLen
	out := []Option{}
	for off < len(data) {
		code := data[off]
		if code == PadCode {
			off++
			continue
		}
		if code == EndCode {
			break
		}
		if off+1 >= len(data) {
			return nil, errors.New("bad-option")
		}
		length := int(data[off+1])
		if off+2+length > len(data) {
			return nil, errors.New("bad-option")
		}
		c.Scanned++
		out = append(out, Option{Code: code, Len: length, Data: data[off+2 : off+2+length]})
		off += 2 + length
	}
	c.src = data
	c.parsed = out
	return out, nil
}

// sameSource 判断 data 是否与缓存的解析结果来自同一份报文。
func (c *Counter) sameSource(data []byte) bool {
	return c.src != nil && len(c.src) == len(data) && &c.src[0] == &data[0]
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
