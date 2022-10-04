package util

import (
	"strings"

	"golang.org/x/net/html"
)

func DomTextContent(node *html.Node) string {
	var sb strings.Builder
	var finder func(*html.Node)

	finder = func(n *html.Node) {
		switch n.Type {
		case html.TextNode:
			sb.WriteString(n.Data)
		case html.ElementNode:
			if n.Data == "br" {
				sb.WriteString("\n")
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			finder(child)
		}
	}

	finder(node)
	return sb.String()
}
