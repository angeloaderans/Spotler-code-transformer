package main

import (
	"fmt"
	_ "html/template"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

func main() {
	body, err := ioutil.ReadFile("form.html")
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	outcome := string(body)

	replaceListItemStart := regexp.MustCompile(`(?miU)<li |<li`)
	replaceListItemEnd := regexp.MustCompile(`(?miU)</li>`)
	replaceListOpeners := regexp.MustCompile(`(?miU)<ul `)
	replaceListEnders := regexp.MustCompile(`(?miU)</ul>`)

	findClusters := regexp.MustCompile(`(?miU)<div id=([\d\D]*)<\/div>`)
	findRequiredMark := regexp.MustCompile(`(?m)mandatorySign`)
	findLabelId := regexp.MustCompile(`(?mU)<label ([\d\D]+) id="([\d\D]+)"`)
	extractId := regexp.MustCompile(`(?mU)id="([\d\D]+)"`)

	outcome = replaceListItemStart.ReplaceAllString(outcome, "<div ")
	outcome = replaceListItemEnd.ReplaceAllString(outcome, "</div>")
	outcome = replaceListOpeners.ReplaceAllString(outcome, "<div ")
	outcome = replaceListEnders.ReplaceAllString(outcome, "</div>")

	clusters := findClusters.FindAllString(outcome, -1)

	for i := range clusters {
		found := findRequiredMark.FindString(clusters[i])
		idString := []string{}

		if found == "mandatorySign" {
			id := findLabelId.FindString(clusters[i])
			idString = append(idString, `id="`+strings.Replace(extractId.FindStringSubmatch(id)[1], "lbl-", "", -1)+`"`)
			idString = append(idString, `id="`+strings.Replace(extractId.FindStringSubmatch(id)[1], "lbl-", "", -1)+`-dd"`)
			idString = append(idString, `id="`+strings.Replace(extractId.FindStringSubmatch(id)[1], "lbl-", "", -1)+`-mm"`)

			for i2 := range idString {
				outcome = strings.Replace(outcome, idString[i2], idString[i2]+" required", -1)
				outcome = strings.Replace(outcome, idString[i2]+"-dd", idString[i2]+" required", 1)
				outcome = strings.Replace(outcome, idString[i2]+"-mm", idString[i2]+" required", 1)
			}
		}
	}

	fmt.Println(outcome)
}
