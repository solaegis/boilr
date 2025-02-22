package template

import (
	"fmt"
	"os"
	"os/user"
	"reflect"
	"strconv"
	"strings"
	"text/template"
	"time"

	"crypto/rand"
	"encoding/base64"

	"github.com/sethvargo/go-password/password"
)

var (
	// custom symbols: removed \ and "
	customSymbols = "~!@#$%^&*()_+`-={}|[]:<>?,./"

	// FuncMap contains the functions exposed to templating engine.
	FuncMap = template.FuncMap{
		// TODO confirmation prompt
		// TODO value prompt
		// TODO encoding utilities (e.g. toBinary)
		// TODO GET, POST utilities
		// TODO Hostname(Also accessible through $HOSTNAME), interface IP addr, etc.
		// TODO add validate for custom regex and expose validate package
		"env":  os.Getenv,
		"time": CurrentTimeInFmt,
		"now": func() time.Time {
			return time.Now()
		},
		"addYear": func(t time.Time) time.Time {
			return t.AddDate(1, 0, 0)
		},
		"addMonth": func(t time.Time) time.Time {
			return t.AddDate(0, 1, 0)
		},
		"addDay": func(t time.Time) time.Time {
			return t.AddDate(0, 0, 1)
		},
		"modifyYear": func(count int, t time.Time) time.Time {
			return t.AddDate(count, 0, 0)
		},
		"modifyMonth": func(count int, t time.Time) time.Time {
			return t.AddDate(0, count, 0)
		},
		"modifyDay": func(count int, t time.Time) time.Time {
			return t.AddDate(0, 0, count)
		},
		"timeToRfc3339": func(t time.Time) string {
			return t.Format(time.RFC3339)
		},
		"timeToDay": func(t time.Time) string {
			return strconv.Itoa(t.Day())
		},
		"timeToHour": func(t time.Time) string {
			return strconv.Itoa(t.Hour())
		},
		"timeToMinute": func(t time.Time) string {
			return strconv.Itoa(t.Minute())
		},
		"timeToMonth": func(t time.Time) string {
			return strconv.Itoa(int(t.Month()))
		},
		"timeToSecond": func(t time.Time) string {
			return strconv.Itoa(t.Second())
		},
		"timeToYear": func(t time.Time) string {
			return strconv.Itoa(t.Year())
		},
		"hostname": func() string { return os.Getenv("HOSTNAME") },
		"username": func() string {
			t, err := user.Current()
			if err != nil {
				return "Unknown"
			}

			return t.Name
		},
		"toBinary": func(s string) string {
			n, err := strconv.Atoi(s)
			if err != nil {
				return s
			}

			return fmt.Sprintf("%b", n)
		},

		"formatFilesize": func(value interface{}) string {
			var size float64

			v := reflect.ValueOf(value)
			switch v.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				size = float64(v.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				size = float64(v.Uint())
			case reflect.Float32, reflect.Float64:
				size = v.Float()
			default:
				return ""
			}

			var KB float64 = 1 << 10
			var MB float64 = 1 << 20
			var GB float64 = 1 << 30
			var TB float64 = 1 << 40
			var PB float64 = 1 << 50

			filesizeFormat := func(filesize float64, suffix string) string {
				return strings.Replace(fmt.Sprintf("%.1f %s", filesize, suffix), ".0", "", -1)
			}

			var result string
			if size < KB {
				result = filesizeFormat(size, "bytes")
			} else if size < MB {
				result = filesizeFormat(size/KB, "KB")
			} else if size < GB {
				result = filesizeFormat(size/MB, "MB")
			} else if size < TB {
				result = filesizeFormat(size/GB, "GB")
			} else if size < PB {
				result = filesizeFormat(size/TB, "TB")
			} else {
				result = filesizeFormat(size/PB, "PB")
			}

			return result
		},
		"password": func(length, numDigits, numSymbols int, noUpper, allowRepeat bool) string {
			generator, err := password.NewGenerator(&password.GeneratorInput{Symbols: customSymbols})

			if err != nil {
				return fmt.Sprintf("failed to generate password generator: %s", err)
			}

			res, err := generator.Generate(length, numDigits, numSymbols, noUpper, allowRepeat)
			if err != nil {
				return fmt.Sprintf("failed to generate password: %s", err)
			}

			return res
		},
		// generate a random base64 string based on random bytes of length n
		"randomBase64": func(length int) string {
			b := make([]byte, length)
			_, err := rand.Read(b)

			if err != nil {
				return fmt.Sprintf("failed to generate randomBase64: %s", err)
			}

			return base64.StdEncoding.EncodeToString(b)
		},

		// String utilities
		"toLower": strings.ToLower,
		"toUpper": strings.ToUpper,
		"toTitle": strings.ToTitle,
		"title":   strings.Title,

		"trimSpace":  strings.TrimSpace,
		"trimPrefix": strings.TrimPrefix,
		"trimSuffix": strings.TrimSuffix,

		"repeat":     strings.Repeat,
		"replace":    strings.Replace,
		"replaceAll": strings.ReplaceAll,

		"kebabCase": func(str string) string {
			return strings.ReplaceAll(str, "_", "-")
		},
		"snakeCase": func(str string) string {
			return strings.ReplaceAll(str, "-", "_")
		},
	}

	// Options contain the default options for the template execution.
	Options = []string{
		// TODO ignore a field if no value is found instead of writing <no value>
		"missingkey=invalid",
	}
)

// CurrentTimeInFmt returns the current time in the given format.
// See time.Time.Format for more details on the format string.
func CurrentTimeInFmt(fmt string) string {
	t := time.Now()

	return t.Format(fmt)
}
