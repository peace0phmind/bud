// Copyright (C) 2022 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

//go:build windows || dragonfly || illumos || solaris || openbsd
// +build windows dragonfly illumos solaris openbsd

package fs

import (
	"github.com/peace0phmind/bud/lib/protocol"
)

func (f *BasicFilesystem) GetXattr(path string, xattrFilter XattrFilter) ([]protocol.Xattr, error) {
	return nil, ErrXattrsNotSupported
}

func (f *BasicFilesystem) SetXattr(path string, xattrs []protocol.Xattr, xattrFilter XattrFilter) error {
	return ErrXattrsNotSupported
}
