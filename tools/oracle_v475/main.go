package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type output struct {
	Server  map[uint32]string `json:"server"`
	Client  map[uint32]string `json:"client"`
	Skipped map[string]string `json:"skipped,omitempty"`
}

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1) //nolint:forbidigo // a CLI must report failure with its exit code
	}
}

func run() error {
	path := flag.String("out", "", "output JSON path")
	flag.Parse()
	if *path == "" && flag.NArg() == 3 {
		return roundTrip(flag.Arg(0), flag.Arg(1), flag.Arg(2))
	}
	if *path == "" {
		flag.Usage()
		return errors.New("-out is required")
	}
	result := output{Server: map[uint32]string{}, Client: map[uint32]string{}, Skipped: map[string]string{}}
	encodePool("server", packet.NewPool(), result.Server, result.Skipped)
	encodePool("client", packet.NewPool(), result.Client, result.Skipped)
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if err := os.MkdirAll(filepath.Dir(*path), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(*path, data, 0o644); err != nil {
		return err
	}
	return nil
}

func roundTrip(direction, rawID, encoded string) error {
	id, err := strconv.ParseUint(rawID, 10, 32)
	if err != nil {
		return err
	}
	data, err := hex.DecodeString(encoded)
	if err != nil {
		return err
	}
	pool := packet.NewPool()
	if direction == "client" {
		pool = packet.NewPool()
	}
	constructor, ok := pool[uint32(id)]
	if !ok {
		return fmt.Errorf("packet %d is missing", id)
	}
	pk := constructor()
	pk.Unmarshal(protocol.NewReader(zeroSafeReader{Reader: bytes.NewReader(data)}, -1))
	var buffer bytes.Buffer
	pk.Marshal(protocol.NewWriter(&buffer, -1))
	_, err = fmt.Fprint(os.Stdout, hex.EncodeToString(buffer.Bytes()))
	return err
}

type zeroSafeReader struct{ *bytes.Reader }

func (r zeroSafeReader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	return r.Reader.Read(p)
}

func encodePool(direction string, pool packet.Pool, encoded map[uint32]string, skipped map[string]string) {
	ids := make([]int, 0, len(pool))
	for id := range pool {
		ids = append(ids, int(id))
	}
	sort.Ints(ids)
	for _, rawID := range ids {
		id := uint32(rawID)
		func() {
			defer func() {
				if recovered := recover(); recovered != nil {
					skipped[fmt.Sprintf("%s/%d", direction, id)] = fmt.Sprint(recovered)
				}
			}()
			var buffer bytes.Buffer
			pool[id]().Marshal(protocol.NewWriter(&buffer, -1))
			encoded[id] = hex.EncodeToString(buffer.Bytes())
		}()
	}
}
