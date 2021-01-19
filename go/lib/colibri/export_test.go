// Copyright 2020 ETH Zurich
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package colibri

import (
	"github.com/scionproto/scion/go/lib/slayers"
	"github.com/scionproto/scion/go/lib/slayers/path/colibri"
)

func PrepareMacInputStatic(s *slayers.SCION, inf *colibri.InfoField,
	hop *colibri.HopField) ([]byte, error) {

	return prepareMacInputStatic(s, inf, hop)
}

func PrepareMacInputSigma(s *slayers.SCION, inf *colibri.InfoField,
	hop *colibri.HopField) ([]byte, error) {

	return prepareMacInputSigma(s, inf, hop)
}

func PrepareMacInputPacket(s *slayers.SCION, packetTimestamp uint64) ([]byte, error) {

	return prepareMacInputPacket(s, packetTimestamp)
}
