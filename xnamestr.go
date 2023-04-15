package xmlutils

import (
	// "io"
	S "strings"
)

// Echo implements Markupper.
func (N XName) Echo() string {
	// if N.Space == "" {
	// 	return N.Local
	// }
	// Assert colon at the end of `N.Space`
	if N.Space != "" && !S.HasSuffix(N.Space, ":") {
		// panic("Missing colon on NS")
		return N.Space + ":" + N.Local
	}
	return N.Space + N.Local
}

/* OBS print stuff

// EchoTo implements Markupper.
func (N XName) EchoTo(w io.Writer) {
	w.Write([]byte(N.Echo()))
}

// DumpTo implements Markupper.
func (N XName) DumpTo(w io.Writer) {
	w.Write([]byte(N.String()))
}

*/
