// Copyright 2021-2022 Faye Amacker
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Reference implementation of CircleHash64 is maintained at
// https://github.com/fxamacker/circlehash

// This file is for Go versions >= 1.17.
// Older versions of Go will build circlehash64_oldgo.go.
//go:build go1.17
// +build go1.17

package circlehash

import (
	"unsafe"
)

// Hash64 returns a 64-bit digest of b.
// Digest is compatible with CircleHash64f.
func Hash64(b []byte, seed uint64) uint64 {
	fn := circle64fShortInput
	if len(b) > 64 {
		fn = circle64f
	}
	return uint64(fn(*(*unsafe.Pointer)(unsafe.Pointer(&b)), seed, uint64(len(b))))
}

// Hash64String returns a 64-bit digest of s.
// Digest is compatible with Hash64.
func Hash64String(s string, seed uint64) uint64 {
	fn := circle64fShortInput
	if len(s) > 64 {
		fn = circle64f
	}
	return uint64(fn(*(*unsafe.Pointer)(unsafe.Pointer(&s)), seed, uint64(len(s))))
}

// Hash64Uint64x2 returns a 64-bit digest of a and b.
// Digest is compatible with Hash64 with byte slice of len 16.
func Hash64Uint64x2(a uint64, b uint64, seed uint64) uint64 {
	return circle64fUint64x2(a, b, seed)
}
