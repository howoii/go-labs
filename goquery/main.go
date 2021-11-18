package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/lipgloss"
)

const (
	lineSize = 64
)

var (
	query = flag.String("query", "", "input text to translate")
	url   = "http://dict.youdao.com/search?q="

	title    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#B8860B"))
	subTitle = lipgloss.NewStyle().Bold(false).Foreground(lipgloss.Color("#B8860B"))
	li       = lipgloss.NewStyle().Foreground(lipgloss.Color("#797979"))
	head     = lipgloss.NewStyle().Foreground(lipgloss.Color("#54A0CF"))
)

func printLine(line string, style lipgloss.Style) {
	var count int
	var buf bytes.Buffer
	for _, v := range line {
		if strings.ContainsRune("()", v) {
			count += 2
		} else {
			count++
		}
		buf.WriteRune(v)
		if count%lineSize == 0 {
			fmt.Println(style.Render(buf.String()))
			buf.Reset()
		}
	}
	if buf.Len() > 0 {
		fmt.Println(style.Render(buf.String()))
	}
}

func textClean(text string) string {
	var buf bytes.Buffer
	for _, v := range text {
		if !unicode.IsSpace(v) {
			buf.WriteRune(v)
		}
	}
	return buf.String()
}

func styleText(text string, style lipgloss.Style) string {
	return style.Render(text)
}

func phrsList(s *goquery.Selection) {
	printLine(s.Find(".wordbook-js .keyword").Text(), title)
	s.Find(".trans-container li").Each(func(i int, ss *goquery.Selection) {
		printLine(textClean(ss.Text()), li)
	})
	fmt.Println()
}

func webTrans(s *goquery.Selection) {
	printLine(s.Find("#webPhrase .title").Text(), subTitle)
	s.Find("#webPhrase .wordGroup").Each(func(i int, ss *goquery.Selection) {
		contentTitle := ss.Find(".contentTitle a").Text()
		contentBody := textClean(ss.Nodes[0].LastChild.Data)
		fmt.Printf("%s %s\n", styleText(contentTitle, head), styleText(contentBody, li))
	})
}

func main() {
	flag.Parse()
	if *query == "" {
		log.Fatalln("query text is empty")
	}
	res, err := http.Get(url + *query)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Fatalln("http get failed")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	phrsList(doc.Find("#phrsListTab"))
	webTrans(doc.Find("#webTrans"))
}
