package yamlmanipulation

import "strings"

type Document string

func DocumentsSeperate(fileContent string) []Document {
	docRaw := strings.Split(fileContent, "---")

	documents := make([]Document, 0, len(docRaw))
	for _, docContent := range docRaw {

		docContentWithoutWhitespace := strings.ReplaceAll(docContent, " ", "")
		docContentWithoutWhitespace = strings.ReplaceAll(docContentWithoutWhitespace, "\n", "")
		docContentWithoutWhitespace = strings.ReplaceAll(docContentWithoutWhitespace, "\r", "")
		docContentWithoutWhitespace = strings.ReplaceAll(docContentWithoutWhitespace, "\t", "")

		if len(docContentWithoutWhitespace) > 0 {
			documents = append(documents, Document(strings.Trim(docContent, "\n ")))
		}
	}

	if len(documents) == 0 {
		documents = append(documents, Document(docRaw[0]))
	}

	return documents
}

func DocumentsJoin(documents []Document) string {
	xs := make([]string, 0, len(documents))

	for _, v := range documents {
		xs = append(xs, string(v))
	}

	joined := strings.Join(xs, "\n---\n")

	addTrailingNewline := true
	if len(documents) == 1 {
		docContent := string(documents[0])
		docContentWithoutWhitespace := strings.ReplaceAll(docContent, " ", "")
		docContentWithoutWhitespace = strings.ReplaceAll(docContentWithoutWhitespace, "\n", "")
		docContentWithoutWhitespace = strings.ReplaceAll(docContentWithoutWhitespace, "\r", "")
		docContentWithoutWhitespace = strings.ReplaceAll(docContentWithoutWhitespace, "\t", "")

		if len(docContentWithoutWhitespace) == 0 {
			addTrailingNewline = false
		}
	}

	if addTrailingNewline {
		joined += "\n"
	}

	return joined
}
