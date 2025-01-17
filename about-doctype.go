package xmlutils

/* This file is documentation only

https://en.wikipedia.org/wiki/Document_type_declaration

http://www.blooberry.com/indexdot/html/tagpages/d/doctype.htm

Sources say that both FPI and URI must be in double quotes, not single.

The most common keywords are DTD, ELEMENT, ENTITIES, and TEXT.
DTD is used only for DTD files.
ELEMENT is usually for DTD fragments i.e. entity or element declarations.
TEXT is used for XML content (text and tags).

Declarations in an external subset are located in a separate text file.
The external subset may be referenced via a formal public identifier
(FPI) and/or a system identifier (URI). Programs for reading documents
may not be required to read the external subset.

Internal subset ONLY:
<!DOCTYPE rootelm [
     <!-- internal subset --> ]>

PUBLIC only:
<!DOCTYPE rootelm PUBLIC "FPI">
or
<!DOCTYPE rootelm PUBLIC "FPI" [
    <!-- internal subset --> ]>

Both PUBLIC and [implied] SYSTEM:
<!DOCTYPE rootelm PUBLIC "FPI" "URI">
or
<!DOCTYPE rootelm PUBLIC "FPI" "URI" [
    <!-- internal subset --> ]>

SYSTEM only:
<!DOCTYPE rootelm SYSTEM "URI">
or
<!DOCTYPE rootelm SYSTEM "URI" [
    <!-- internal subset --> ]>
*/
