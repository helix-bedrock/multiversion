# multiversion

Older Minecraft: Bedrock Edition versions for Helix: one adapter per version
that converts packets, blocks, items, recipes and chunks to and from the latest
protocol.

**Status:** in development

Part of [Helix](https://github.com/helix-bedrock), an open-source Minecraft: Bedrock
Edition server engine in Go. Each adapter implements
[`version.Protocol`](https://github.com/helix-bedrock/protocol/tree/main/version)
from `helix-bedrock/protocol`; pass it to `conn.Config.Protocol`.

## Protocols

| Protocol ID | Minecraft version | Adapter | Support | Tested |
|------------:|-------------------|---------|:-------:|:------:|
| 2193 | 1.26.50 | Native gophertunnel | ✅ | ✅ |
| 2169 | 1.26.45 | `v1_26_45` | ✅ | ✅ |
| 2168 | 1.26.40-1.26.44 | `v1_26_44` | ✅ | ✅ |
| 1001 | 1.26.30-1.26.34, 1.26.36 | `v1_26_30` | ✅ | ✅ |
| 975 | 1.26.20, 1.26.21, 1.26.23 | `v1_26_20` | ✅ | ✅ |
| 944 | 1.26.10-1.26.14 | `v1_26_10` | ✅ | ✅ |
| 924 | 1.26.0-1.26.3 | `v1_26_0` | ✅ | ✅ |
| 898 | 1.21.130-1.21.132 | `v1_21_130` | ✅ | ✅ |
| 844 | 1.21.110-1.21.114 | `v1_21_110` | ✅ | ✅ |
| 827 | 1.21.100-1.21.102 | `v1_21_100` | ✅ | ✅ |
| 766 | 1.21.50-1.21.51 | `v1_21_50` | ✅ | ✅ |
| 748 | 1.21.40, 1.21.41, 1.21.43, 1.21.44 | `v1_21_40` | ✅ | ✅ |
| 486 | 1.18.10-1.18.12 | `v1_18_10` | ✅ | ✅ |
| 475 | 1.18.0-1.18.2 | `v1_18_0` | ✅ | ✅ |
| 419 | 1.16.100 | `v1_16_100` | ✅ | ✅ |

Unlisted releases and preview builds are not implied.

## Usage

Registry-aware adapters need the server's block and item registries:

```go
protocols, err := (multiversion.Config{MinimumProtocol: 748}).ProtocolsWithRegistries(blocks, items)
```

`MinimumProtocol` is an inclusive floor: `748` allows 1.21.40 and newer, `0`
(the default) allows every supported version.

Not ready yet:

- Protocols 486, 475 and 419 (1.18 and 1.16) need RakNet protocol 10 in
  Helix raknet, and a server handshake that starts with Login.
- Chunk and palette conversion needs hooks in the Helix world, which does not
  exist yet. Upstream used a Dragonfly fork for this.

## Credits

Imported from [go-multiversion](https://github.com/shawtymarco/go-multiversion)
by **shawtymarco** at commit
[`2e8c8c3aad13`](https://github.com/shawtymarco/go-multiversion/commit/2e8c8c3aad13). It is MIT
licensed (copyright Sandertv); the original license is in
[`licenses/go-multiversion-LICENSE`](licenses/go-multiversion-LICENSE). Helix
maintains this code from here on. See [CREDITS.md](CREDITS.md).

go-multiversion credits:
- [Sandertv/gophertunnel](https://github.com/Sandertv/gophertunnel)
- [df-mc/dragonfly](https://github.com/df-mc/dragonfly)
- [df-mc/worldupgrader](https://github.com/df-mc/worldupgrader)
- [Mojang/bedrock-protocol-docs](https://github.com/Mojang/bedrock-protocol-docs)
- [EndstoneMC/bedrock-server-data](https://github.com/EndstoneMC/bedrock-server-data)
- [TedacMC/tedac](https://github.com/TedacMC/tedac) for protocol-419 registry cross-checks
- Mascot adapted with OpenAI ImageGen from the [Go gopher](https://go.dev/blog/gopher) by Renee French, under [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/).

## License

MIT. See [LICENSE](LICENSE).
