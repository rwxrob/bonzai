//go:build aix || js || nacl || plan9 || windows || android || solaris

package term

func WinSizeUpdate() {
	WinSize = WinSizeStruct{80, 24, 100, 100}
}
