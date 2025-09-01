package utils

import "regexp"

func GetImageLink(version, path string) string {
	imageRefRegex := regexp.MustCompile(`^([a-zA-Z0-9.-/]*)?[a-zA-Z0-9._-]+:[a-zA-Z0-9._-]+$`)
	if imageRefRegex.MatchString(version) {
		return version
	}
	return path + ":" + version
}
