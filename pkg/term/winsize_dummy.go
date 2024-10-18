// Copyright 2022 Robert S. Muhlestein
// SPDX-License-Identifier: Apache-2.0

//go:build aix || js || nacl || plan9 || windows || android || solaris

package term

func WinSizeUpdate() {
	WinSize = WinSizeStruct{24, 80, 100, 100}
}
