
package epic

import (
	"github.com/scionproto/scion/go/lib/serrors"
	"github.com/scionproto/scion/go/lib/slayers/path"
	"github.com/scionproto/scion/go/lib/slayers/path/scion"
)

const PathType path.Type = 3

func RegisterPath() {
	path.RegisterPath(path.Metadata{
		Type: PathType,
		Desc: "Epic",
		New: func() path.Path {
			return &EpicPath{ScionRaw: &scion.Raw{}}
		},
	})
}

type EpicPath struct {
	PacketTimestamp uint64
	PHVF            []byte
	LHVF            []byte
	ScionRaw		*scion.Raw
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
	// todo: validate
	return 16 + p.ScionRaw.Len()
}

func (p EpicPath) Type() path.Type {
	return PathType
}
