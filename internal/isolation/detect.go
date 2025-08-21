package isolation

import "runtime"

func DetectOS() string {
	switch runtime.GOOS {
	case "linux":
		return "linux"
	case "windows":
		return "windows"
	case "android":
		return "android"
	default:
		return "unknown"
	}
}
