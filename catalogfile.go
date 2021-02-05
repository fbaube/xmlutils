package xmlmodels

import (
	"encoding/xml"
	"fmt"
	"os"

	// "path"

	FP "path/filepath"
	S "strings"

	SU "github.com/fbaube/stringutils"
)

// XmlCatalogFile represents a parsed XML catalog file, at the top level.
type XmlCatalogFile struct {
	XMLName xml.Name `xml:"catalog"̀`
	// "public" or "system"
	Prefer                string                    `xml:"prefer,attr"`
	XmlPublicIDsubrecords []PIDSIDcatalogFileRecord `xml:"public"`
	// We do this so we can peel off the directory path
	AbsFilePath string
}

func (p *XmlCatalogFile) GetByPublicID(s string) *PIDSIDcatalogFileRecord {
	if s == "" {
		return nil
	}
	for _, entry := range p.XmlPublicIDsubrecords {
		if s == string(entry.XmlPublicID) {
			return &entry
		}
	}
	return nil
}

// NewXmlCatalogFile is a convenience function that reads in the
// file and then processes the file contents. It is not clear what the
// constraints on the path are (but a relative path should work okay).
func NewXmlCatalogFile(fpath string) (pXC *XmlCatalogFile, err error) {
	if fpath == "" {
		return nil, nil
	}
	var raw []byte
	var e error
	raw, e = os.ReadFile(fpath)
	if e != nil {
		println("==> Can't read XML catalog file:", fpath, ", reason:", e)
		return nil, fmt.Errorf("gparse.NewXmlCatalog.ReadFile<%s>: %w", fpath, e)
	}

	var pCPR *ParserResults_xml
	pCPR, e = GenerateParserResults_xml(string(raw))
	if e != nil {
		return nil, fmt.Errorf("gparse.xml.parseResults: %w", e)
	}
	var catRoot xml.StartElement      // *gtoken.GToken
	var pubEntries []xml.StartElement // []*gtoken.GToken
	catRoot = GetFirstStartElmByTag(pCPR.NodeSlice, "catalog")
	pubEntries = GetAllStartElmsByTag(pCPR.NodeSlice, "public")
	if catRoot.Name.Local == "" {
		panic("No <catalog> root elm")
	}
	pXC = new(XmlCatalogFile)
	pXC.XMLName = catRoot.Name // xml.Name(gktnRoot.GName)
	pXC.Prefer = GetAttVal(catRoot, "prefer")
	pXC.XmlPublicIDsubrecords = make([]PIDSIDcatalogFileRecord, 0)

	for _, GT := range pubEntries {
		// println("  CAT-ENTRY:", GT.Echo()) // entry.GAttList.Echo())
		pID, e := NewSIDPIDcatalogRecordfromStartElm(GT)
		// NOTE Gotta fix the filepath
		// // ## pID.AbsFilePath = // FU.AbsFilePath(
		// // ## 	FU.AbsWRT(string(pID.AbsFilePath), FP.Dir(string(fpath))) // )
		if e != nil {
			panic(e)
		}
		if pID == nil {
			fmt.Printf("Got NIL from: %+v \n", GT)
		}
		pXC.XmlPublicIDsubrecords = append(pXC.XmlPublicIDsubrecords, *pID)
	}

	// ==============================

	// NOTE The following code is UGLY and needs to be FIXED.
	// fileDir := pXC.AbsFilePath.DirPath()
	fileDir, _ := FP.Split(pXC.AbsFilePath)
	// return AbsFilePath(dp)
	println("XML catalog fileDir:", fileDir)
	for _, entry := range pXC.XmlPublicIDsubrecords {
		println("  Entry's AbsFilePath:" /* FIXME:60 MU.Tilded*/, (entry.AbsFilePath))
		entry.AbsFilePath = fileDir + "/" + string(entry.AbsFilePath)
	}
	ok := pXC.Validate()
	if !ok {
		panic("BAD CAT")
	}
	// println("==> Processed XML catalog at:", pXC.FileFullName.String())
	// println("TODO: insert file path for catalog file")
	return pXC, nil
}

func NewSIDPIDcatalogRecordfromStartElm(se xml.StartElement) (pID *PIDSIDcatalogFileRecord, err error) {
	if se.Name.Local == "" {
		return nil, nil
	}
	fmt.Printf("L.174 GT: %+v \n", se)
	NS := se.Name.Space
	if NS != "" && NS != NS_OASIS_XML_CATALOG {
		panic("XML catalog entry has bad NS: " + NS)
	}
	println("L.179 Space:", se.Name.Space, "/ Local:", se.Name.Local)
	attPid := GetAttVal(se, "publicId")
	attUri := GetAttVal(se, "uri")
	if attPid == "" && attUri == "" {
		println("Empty GToken for Public ID!")
		return nil, nil
	}
	println("L.186 attPid is:", attPid)
	println("L.187 attUri is:", attUri)

	// -//OASIS//DTD LIGHTWEIGHT DITA Topic//EN
	var ss []string
	ss = S.Split(attPid, "/")
	// fmt.Printf("(DD) (%d) %#v \n", len(ss), ss)
	ss = SU.DeleteEmptyStrings(ss)
	// {"-", "OASIS", "DTD LIGHTWEIGHT DITA Topic", "EN"}
	// fmt.Printf("(DD:PIDss) (%d) %v \n", len(ss), ss)
	if len(ss) != 4 || ss[0] != "-" || ss[3] != "EN" {
		return nil, fmt.Errorf("Malformed Public ID<%s>", attPid)
	}
	pID = new(PIDSIDcatalogFileRecord)
	// NOTE DANGER This is probably relative not absolute,
	// and has to be fixed by the caller
	pID.XmlPublicID = XmlPublicID(attPid)
	pID.XmlSystemID = XmlSystemID(attUri)
	pID.AbsFilePath = attUri // FU.AbsFilePath(attUri)
	pID.Organization = ss[1]
	pID.IsOasis = ("OASIS" == pID.Organization)
	pID.PublicTextClass, pID.PublicTextDesc = SU.SplitOffFirstWord(ss[2])
	// ilog.Printf("PubID|%s|%s|%s|\n", p.Organization, p.PTClass, p.PTDesc)
	// fmt.Printf("(DD:pPID) PubID<%s|%s|%s> AFP<%s>\n",
	//  	pID.Organization, pID.PublicTextClass,
	//		pID.PublicTextDesc, pID.AbsFilePath)
	return pID, nil
}

// Validate validates an XML catalog. It checks that the listed files exist
// and that the IDs (as strings that are not parsed yet) are well-formed.
// It assumes that the catalog has already been loaded from an XML catalog
// file on-disk. The return value is false if _any_ entry fails to load,
// but also each entry has its own error field.
func (p *XmlCatalogFile) Validate() (retval bool) {
	retval = true
	for i, pEntry := range p.XmlPublicIDsubrecords {
		if "" == pEntry.XmlPublicID {
			println("OOPS:", pEntry.String())
			panic(fmt.Sprintf("Missing Public ID in catalog entry[%d]: %s",
				i, p.AbsFilePath)) // Parts.String()))
		}
		var abspath, dirpart string
		dirpart, _ = FP.Split(p.AbsFilePath)
		abspath = dirpart + "/" + string(pEntry.XmlSystemID)

		var e error
		_, e = os.ReadFile(abspath)
		if e != nil {
			fmt.Printf("==> Catalog<%s>: Bad System ID / URI <%s> for Public ID <%s> (%s) \n",
				p.AbsFilePath, pEntry.XmlSystemID, pEntry.XmlPublicID, e.Error())
			retval = false
			continue
		}
		// NOTE The loop variable "entry" is by value, not reference !
		// entry.FilePath = FU.FilePath(pIF.FileFullName.String())
		p.XmlPublicIDsubrecords[i].AbsFilePath = abspath

		// Now do some fancy parsing of the Public ID
		var s = string(pEntry.XmlPublicID)
		if SU.IsXmlQuoted(s) {
			s = SU.MustXmlUnquote(s)
		}
		// -//OASIS//DTD LIGHTWEIGHT DITA Topic//EN
		var ss []string
		ss = S.Split(s, "/")
		// fmt.Printf("(DD) (%d) %#v \n", len(ss), ss)
		ss = SU.DeleteEmptyStrings(ss)
		// {"-", "OASIS", "DTD LIGHTWEIGHT DITA Topic", "EN"}
		// fmt.Printf("(DD) (%d) %#v \n", len(ss), ss)
		if len(ss) != 4 || ss[0] != "-" || ss[3] != "EN" {
			retval = false
			pEntry.Err = fmt.Errorf("Malformed Public ID: %s", s)
			continue
		}
		pEntry.Organization = ss[1]
		pEntry.IsOasis = ("OASIS" == pEntry.Organization)
		pEntry.PublicTextClass, pEntry.PublicTextDesc = SU.SplitOffFirstWord(ss[2])
		// ilog.Printf("PubID|%s|%s|%s|\n", p.Organization, p.PTClass, p.PTDesc)
		// fmt.Printf("(DD:pPID) PubID<%s|%s|%s>\n", p.Organization, p.PTClass, p.PTDesc)
	}
	return true
}
