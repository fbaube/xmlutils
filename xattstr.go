package xmlutils

import (
// "io"
)

// Echo implements Markupper (and inserts a leading space).
func (A XAtt) Echo() string {
	return " " + XName(A.Name).Echo() + "=\"" + A.Value + "\""
}

// Echo implements Markupper (and inserts spaces).
func (AL XAtts) Echo() string {
	var s string
	for _, A := range AL {
		s += " " + XName(A.Name).Echo() + "=\"" + A.Value + "\""
	}
	return s
}

/* OBS print stuff

// EchoTo implements Markupper.
func (A XAtt) EchoTo(w io.Writer) {
	w.Write([]byte(A.Echo()))
}

// EchoTo implements Markupper.
func (AL XAtts) EchoTo(w io.Writer) {
	w.Write([]byte(AL.Echo()))
}

// String implements Markupper.
func (A XAtt) String() string {
	return A.Echo()
}

// String implements Markupper.
func (AL XAtts) String() string {
	return AL.Echo()
}

// DumpTo implements Markupper.
func (A XAtt) DumpTo(w io.Writer) {
	w.Write([]byte(A.String()))
}

// DumpTo implements Markupper.
func (AL XAtts) DumpTo(w io.Writer) {
	w.Write([]byte(AL.String()))
}

*/
