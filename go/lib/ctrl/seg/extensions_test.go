// Copyright 2020 Anapaya Systems
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package seg

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/scionproto/scion/go/lib/ctrl/seg/extensions/digest"
	"github.com/scionproto/scion/go/lib/ctrl/seg/unsigned_extensions/epic_detached"
)

func TestDecodeEncodeEpicDigest(t *testing.T) {
	hop := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	peers := make([][]byte, 0, 5)
	for i := 0; i < 5; i++ {
		peer := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
		peers = append(peers, peer)

		ed := &epic_detached.EpicDetached{
			AuthHopEntry:    hop,
			AuthPeerEntries: peers,
		}
		hash, _ := digest.CalcEpicDigest(ed, true)
		dig := &digest.DigestExtension{
			Epic: hash,
		}

		ext := Extensions{
			Digests: dig,
		}
		ext2 := ExtensionsFromPB(ExtensionsToPB(ext))
		assert.Equal(t, ext, ext2)
	}
}