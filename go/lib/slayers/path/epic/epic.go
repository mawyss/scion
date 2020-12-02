
package epic

import (
	"github.com/scionproto/scion/go/lib/serrors"
	"github.com/scionproto/scion/go/lib/slayers/path"
	"github.com/scionproto/scion/go/lib/slayers/path/scion"
)

const PathType path.Type = 3

type EpicPath struct {
	PacketTimestamp uint64
	PHVF            []byte
	LHVF            []byte
	ScionPath		*scion.Raw
}

func (p EpicPath) SerializeTo(b []byte) error {
	return serrors.New("todo")
}

func (p EpicPath) DecodeFromBytes(b []byte) error {
	return serrors.New("todo")
}

func (p EpicPath) Reverse() (path.Path, error) {
	return nil, serrors.New("todo")
}

func (p EpicPath) Len() int {
	return 0
}

func (p EpicPath) Type() path.Type {
	return 0
}
