# dhcp-packet

解析 DHCPv4 报文的选项区。只用标准库。

```
go test ./...
go vet ./...
go run ./cmd/dhcpdump --sample=pad
go run ./cmd/dhcpdump --sample=end
go run ./cmd/dhcpdump --sample=repeat
go run ./cmd/dhcpdump --sample=work
```

## 口径

- **布局**：前 236 字节是固定头，接着 4 字节魔术字，之后才是选项区。
- **选项**：每项是 `code`（1 字节）+ `len`（1 字节）+ 数据；数据长度越界就是坏报文。
- **填充**：`code` 为 0 是单字节填充，没有长度字段，跳过不算选项。
- **结束**：遇到 `code` 255 就结束，后面的字节不参与解析。
- **重复**：同一个 `code` 出现多次要各留一条，顺序与报文一致。

## 不变量

- 选项列表与报文里的顺序一致，重复项不合并。
- 结束标记之后的字节不产生选项。
- `scanned` 不随选项数乘查询次数放大：两千条选项查两千次的 `scanned` 不超过 9000。

## 输出契约

每个场景打印一行 JSON；`--sample=work` 只打印 `{"options": N, "scanned": N}`。
