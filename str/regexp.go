package str

import (
	"regexp"
	"strings"

	"github.com/Muskchen/toolkits/logger"
)

var IPReg, _ = regexp.Compile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)
var MailReg, _ = regexp.Compile(`\w[-._\w]*@\w[-._\w]*\.\w+`)

func IsMatch(s, pattern string) bool {
	match, err := regexp.Match(pattern, []byte(s))
	if err != nil {
		return false
	}
	return match
}

func IsIdentifier(s string, pattern ...string) bool {
	defpattern := "^[a-zA-Z0-9\\-\\_\\.]+$"
	if len(pattern) > 0 {
		defpattern = pattern[0]
	}
	return IsMatch(s, defpattern)
}

func IsMail(s string) bool {
	return MailReg.MatchString(s)
}

func IsPhone(s string) bool {
	if strings.HasPrefix(s, "+") {
		return IsMatch(s[1:], `^\d{13}$`)
	} else {
		return IsMatch(s, `^{11}$`)
	}
}

func IsIp(s string) bool {
	return IPReg.MatchString(s)
}

func Dangerous(s string) bool {
	if strings.Contains(s, "<") {
		return true
	}

	if strings.Contains(s, ">") {
		return true
	}

	if strings.Contains(s, "&") {
		return true
	}

	if strings.Contains(s, "'") {
		return true
	}

	if strings.Contains(s, "\"") {
		return true
	}

	if strings.Contains(s, "file://") {
		return true
	}

	if strings.Contains(s, "../") {
		return true
	}

	return false
}

func GetPatAndTimeFormat(s string) (string, string) {
	var pattern, timeFormat string
	switch s {
	case "dd/mmm/yyyy:HH:MM:SS":
		pattern = `([012][0-9]|3[01])/[JFMASONDjfmasond][a-zA-Z]{2}/(2[0-9]{3}):([01][0-9]|2[0-4])(:[012345][0-9]){2}`
		timeFormat = "02/Jan/2006:15:04:05"
	case "dd/mmm/yyyy HH:MM:SS":
		pattern = `([012][0-9]|3[01])/[JFMASONDjfmasond][a-zA-Z]{2}/(2[0-9]{3})\s([01][0-9]|2[0-4])(:[012345][0-9]){2}`
		timeFormat = "02/Jan/2006 15:04:05"
	case "yyyy-mm-ddTHH:MM:SS":
		pattern = `(2[0-9]{3})-(0[1-9]|1[012])-([012][0-9]|3[01])T([01][0-9]|2[0-4])(:[012345][0-9]){2}`
		timeFormat = "2006-01-02T15:04:05"
	case "dd-mmm-yyyy HH:MM:SS":
		pattern = `([012][0-9]|3[01])-[JFMASONDjfmasond][a-zA-Z]{2}-(2[0-9]{3})\s([01][0-9]|2[0-4])(:[012345][0-9]){2}`
		timeFormat = "02-Jan-2006 15:04:05"
	case "yyyy-mm-dd HH:MM:SS":
		pattern = `(2[0-9]{3})-(0[1-9]|1[012])-([012][0-9]|3[01])\s([01][0-9]|2[0-4])(:[012345][0-9]){2}`
		timeFormat = "2006-01-02 15:04:05"
	case "yyyy/mm/dd HH:MM:SS":
		pattern = `(2[0-9]{3})/(0[1-9]|1[012])/([012][0-9]|3[01])\s([01][0-9]|2[0-4])(:[012345][0-9]){2}`
		timeFormat = "2006/01/02 15:04:05"
	case "yyyymmdd HH:MM:SS":
		pattern = `(2[0-9]{3})(0[1-9]|1[012])([012][0-9]|3[01])\s([01][0-9]|2[0-4])(:[012345][0-9]){2}`
		timeFormat = "20060102 15:04:05"
	case "mmm dd HH:MM:SS":
		pattern = `[JFMASONDjfmasond][a-zA-Z]{2}\s+([1-9]|[1-2][0-9]|3[01])\s([01][0-9]|2[0-4])(:[012345][0-9]){2}`
		timeFormat = "Jan 2 15:04:05"
	case "mmdd HH:MM:SS":
		pattern = `(0[1-9]|1[012])([012][0-9]|3[01])\s([01][0-9]|2[0-4])(:[012345][0-9]){2}`
		timeFormat = "0102 15:04:05"
	default:
		logger.Errorf("match time pac failed : [timeFormat:%s]", s)
		return "", ""
	}
	return pattern, timeFormat
}
