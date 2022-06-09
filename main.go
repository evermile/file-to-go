package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var imports = []string{
	"import (",
	"\t\"bytes\"",
	"\t\"compress/gzip\"",
	"\t\"encoding/base64\"",
	"\t\"strings\"",
	")",
}

var decodeFunc = []string{
	"func GetFullSwagger() string {",
	"	zipped, _ := base64.StdEncoding.DecodeString(strings.Join(openapiSpec, \"\"))",
	"	zr, _ := gzip.NewReader(bytes.NewReader(zipped))",
	"	var buf bytes.Buffer",
	"	_, _ = buf.ReadFrom(zr)",
	"	return buf.String()",
	"}",
}

func main() {
	apiFile := flag.String("file", "bla.json", "The file to read")
	outFile := flag.String("outfile", "bla.go", "The output file to write")
	pkg := flag.String("package", "main", "The package for the output class")

	flag.Parse()
	encoded := getZippedContentBase64(apiFile)
	if file, err := os.Create(*outFile); err != nil {
		panic(err)
	} else {
		writer := bufio.NewWriter(file)
		writer.WriteString(fmt.Sprintf("package %s \n\n", *pkg))
		for _, l := range imports {
			writer.WriteString(l + "\n")
		}

		writer.WriteString("var openapiSpec = []string{\n")
		writer.WriteString("\t\"")
		for i, r := range encoded {
			writer.WriteRune(r)
			if (i+1)%100 == 0 {
				writer.WriteString("\", \n\t\"")
			}
		}
		writer.WriteString("\", \n")
		writer.WriteString("}\n")

		for _, l := range decodeFunc {
			writer.WriteString(l + "\n")
		}
		writer.WriteString("\n")

		writer.Flush()
	}
}

func getZippedContentBase64(apiFile *string) string {
	var content []byte
	if file, err := os.Open(*apiFile); err != nil {
		panic(err)
	} else {
		reader := bufio.NewReader(file)
		if content, err = ioutil.ReadAll(reader); err != nil {
			panic(err)
		}
	}
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	defer func(zw *gzip.Writer) {
		_ = zw.Close()
	}(zw)

	if _, err := zw.Write(content); err != nil {
		panic(err)
	}
	_ = zw.Flush()
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}
