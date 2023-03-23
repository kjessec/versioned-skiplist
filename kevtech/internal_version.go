package kevtech

import "unsafe"

var (
	sizeIv = int(unsafe.Sizeof(internalVersionValue{}))
)

type internalVersionValue struct {
	version uint64
	value   unsafe.Pointer
}

func findNearestVersion(arr []unsafe.Pointer, targetVersion uint64) unsafe.Pointer {
	var l, r = uint64(0), uint64(len(arr) - 1)

	for l <= r {
		m := l + (r-l)/2

		if (*internalVersionValue)(arr[m]).version == targetVersion {
			return arr[m]
		} else if (*internalVersionValue)(arr[m]).version < targetVersion {
			l = m + 1
		} else {
			r = m - 1
		}
	}

	return unsafe.Pointer(nil) // x is not found in arr
}
