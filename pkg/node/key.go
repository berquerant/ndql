package node

import "slices"

const (
	KeyPath    = "path"     // file or directory path.
	KeySize    = "size"     // file size in bytes.
	KeyIsDir   = "is_dir"   // if true, it's a directory.
	KeyModTime = "mod_time" // last modified time.
	KeyMode    = "mode"     // file or directory mode.
)

func BuiltinKeys() []string {
	return []string{
		KeyPath,
		KeySize,
		KeyIsDir,
		KeyModTime,
		KeyMode,
	}
}

func IsBuiltinKey(key string) bool { return slices.Contains(BuiltinKeys(), key) }
