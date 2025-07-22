package helper

import (
	"strings"

	"golang.org/x/net/html"
)

var allowedTags = map[string]bool{
	"b": true, "strong": true,
	"i": true, "em": true,
	"u": true, "ins": true,
	"s": true, "strike": true, "del": true,
	//"tg-spoiler": true, "span": true, // <span class="tg-spoiler">
	"code": true, "pre": true,
	"a": true,
}

func SanitizeTelegramHTML(input string) string {
	// Telegram не поддерживает <br> и <p>, но поддерживает \n
	// Поэтому заменяем их на \n
	input = strings.ReplaceAll(input, "<br>", "\n")
	input = strings.ReplaceAll(input, "<br/>", "\n")
	input = strings.ReplaceAll(input, "<p>", "")
	input = strings.ReplaceAll(input, "</p>", "\n")

	node, err := html.Parse(strings.NewReader(input))
	if err != nil {
		return input // если HTML поломан — просто возвращаем исходный
	}
	var b strings.Builder
	var render func(*html.Node)
	render = func(n *html.Node) {
		switch n.Type {
		case html.ElementNode:
			if allowedTags[n.Data] {
				b.WriteString("<" + n.Data)
				for _, attr := range n.Attr {
					// Telegram разрешает только href у <a> и class="tg-spoiler" у <span>
					if (n.Data == "a" && attr.Key == "href") ||
						(n.Data == "span" && attr.Key == "class" && attr.Val == "tg-spoiler") {
						b.WriteString(" " + attr.Key + "=\"" + attr.Val + "\"")
					}
				}
				b.WriteString(">")
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					render(c)
				}
				b.WriteString("</" + n.Data + ">")
			} else {
				// Просто рендерим детей, пропуская сам тег
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					render(c)
				}
			}
		case html.TextNode:
			b.WriteString(n.Data)
		default:
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				render(c)
			}
		}
	}
	render(node)
	return b.String()
}
