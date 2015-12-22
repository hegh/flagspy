// Package flagspy provides basic access to command line flags before they
// are made available by flag.Parse().
//
// This access is not as straightforward, as the flag library, and that is
// intentional. If you can get away with using the flag library, use it instead.
//
// Not intended for use after program initialization. Copy any values out that
// you need during init(), and do not call into flagspy outside of init(). The
// main package is expected to call Free(), invalidating all flagspy data.
//
// Flags are parsed from the os.Args slice by these rules:
//  * Flags are parsed from os.Args[1] onward.
//  * Flag parsing stops if an element of os.Args is the string "--".
//  * An element matches a flag name if the element consists of one or two "-"
//    characters immediately followed by the flag name.
//  * If an element beginning with a flag name is immediately followed by an "="
//    then the flag value is the rest of that element after the "=".
//  * Otherwise:
//    * If the following element does not begin with a "-" then the
//      flag value is that next element and flag parsing resumes after it.
//    * Otherwise the flag value is "".
//  * If a value begins with ' or " and ends with the same character, those
//    characters are removed.
//
// Anticipated usage:
//     func init() {
//       logdir, ok := flagspy.Get("logdir")
//       if !ok {
//         logdir = os.Getwd()
//       }
//       // ... set up logging ...
//     }
//
// A concrete example of how flags are parsed:
//     ./exe -flag1 --flag2 value2 -flag3='value3' -- -flag4=value4
//     Get("flag1"): "", true
//     Get("flag2"): "value2", true
//     Get("flag3"): "value3", true
//     Get("flag4"): "", false
package flagspy

import (
	"os"
	"strings"
	"sync"
)

var (
	flags map[string]string
	once  sync.Once
)

func initialize() {
	flags = make(map[string]string)
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--" {
			break
		}

		if strings.HasPrefix(os.Args[i], "-") {
			var f string
			if strings.HasPrefix(os.Args[i], "--") {
				f = os.Args[i][2:]
			} else {
				f = os.Args[i][1:]
			}

			var v string
			if p := strings.Index(f, "="); p != -1 {
				v = f[p+1:]
				f = f[:p]
			} else if i < len(os.Args)-1 && !strings.HasPrefix(os.Args[i+1], "-") {
				i++
				v = os.Args[i]
			}

			if (strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'")) ||
				(strings.HasPrefix(v, "\"") && strings.HasSuffix(v, "\"")) {
				v = v[1 : len(v)-1]
			}

			flags[f] = v
		}
	}
}

// Get the value of the named flag, and whether it was specified on the command line.
func Get(name string) (string, bool) {
	once.Do(initialize)
	v, ok := flags[name]
	return v, ok
}

// Free clears the cached flag values, allowing GC to free the memory.
// DO NOT call Get after calling Free.
func Free() {
	flags = nil // Any future call to Get() will panic on a nil dereference.
}
