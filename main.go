package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
	"strings"
)

var links []string
var url = "https://dungeon.su"

func main() {
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{"https://dungeon.su/items/"},
		ParseFunc: linksParse,
	}).Start()

	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: links,
		ParseFunc: mainParser,
		Exporters: []export.Exporter{&export.JSON{FileName: "items.json"}},
		ConcurrentRequests: 10,
	}).Start()
/*	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{"https://dungeon.su/spells/"},
		ParseFunc: linksParse,
	}).Start()

	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: links,
		ParseFunc: mainParser,
		Exporters: []export.Exporter{&export.JSON{FileName: "spells.json"}},
		ConcurrentRequests: 10,
	}).Start()


	links = []string{}

	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{"https://dungeon.su/bestiary/"},
		ParseFunc: linksParse,
	}).Start()

	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: links,
		ParseFunc: mainParser,
		Exporters: []export.Exporter{&export.JSON{FileName: "beasts.json"}},
		ConcurrentRequests: 10,
	}).Start()
*/
}

func linksParse(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find("ul.list-of-items").First().Find("a").Each(func(i int, selection *goquery.Selection) {
		link, _ := selection.Attr("href")
		links = append(links, url+link)
	})
}

func imgParse(g *geziyor.Geziyor, r *client.Response) []string {
	var imgs []string
	r.HTMLDoc.Find("section.gallery").Find("img").Each(func(i int, selection *goquery.Selection) {
		link, _ := selection.Attr("src")
		imgs = append(imgs, url+link)
	})
	return imgs
}


func mainParser(g *geziyor.Geziyor, r *client.Response) {
	values := map[string]string{}

	outData := map[string]interface{}{
		"Название":   r.HTMLDoc.Find("a.item-link").Text(),
	}

	r.HTMLDoc.Find("ul.params").ChildrenFiltered("li").Slice(1, goquery.ToEnd).Each(func(i int, selection *goquery.Selection) {
		if _, found := selection.Attr("class"); !found {
			tempValues := strings.Split(selection.Text(), ": ")
			values[tempValues[0]] = tempValues[1]
		}

	})
	r.HTMLDoc.Find("div.stat").Each(func(i int, selection *goquery.Selection) {
		attr, _ := selection.Attr("title")
		values[attr] = selection.Text()[6:]

	})

	r.HTMLDoc.Find("li.subsection").Each(func(i int, selection *goquery.Selection) {
		title := selection.Find("h3.subsection-title").Text()
		if len(title) == 0 {
			return
		}

		var texts []string

		selection.Find("p").Each(func(i int, selection *goquery.Selection) {
			texts = append(texts, selection.Text())
		})
		outData[title] = texts
	})

	imgs := imgParse(g, r)
	if len(imgs) > 0 {
		outData["Изображения"] = imgs
	}
	for k, v := range values {
		outData[k] = v
	}

	g.Exports <- outData
}

