package dhcp

import "testing"

func build(options []byte) []byte {
	data := make([]byte, HeaderLen+MagicLen)
	copy(data[236:], []byte{0x63, 0x82, 0x53, 0x63})
	return append(data, options...)
}

func TestSingleOption(t *testing.T) {
	got, err := Parse(build([]byte{53, 1, 2}), &Counter{})
	if err != nil {
		t.Fatalf("解析失败：%v", err)
	}
	if len(got) != 1 || got[0].Code != 53 || got[0].Len != 1 {
		t.Fatalf("应解出一条 code=53 的选项：%+v", got)
	}
}

func TestTwoOptions(t *testing.T) {
	got, err := Parse(build([]byte{53, 1, 2, 12, 3, 'a', 'b', 'c'}), &Counter{})
	if err != nil {
		t.Fatalf("解析失败：%v", err)
	}
	if len(got) != 2 {
		t.Fatalf("应解出两条选项：%+v", got)
	}
}

func TestTooShort(t *testing.T) {
	if _, err := Parse([]byte{1, 2, 3}, &Counter{}); err == nil {
		t.Fatal("太短的报文应报错")
	}
}

func TestTruncatedOption(t *testing.T) {
	if _, err := Parse(build([]byte{53, 9, 1}), &Counter{}); err == nil {
		t.Fatal("选项长度越界应报错")
	}
}

func TestLookupByCode(t *testing.T) {
	got, err := Lookup(build([]byte{53, 1, 2, 12, 1, 'x'}), 53, &Counter{})
	if err != nil {
		t.Fatalf("查找失败：%v", err)
	}
	if len(got) != 1 || got[0].Code != 53 {
		t.Fatalf("应找到一条 code=53：%+v", got)
	}
}
