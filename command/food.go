package command

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
)

// Menu returns food in available restaurants
func Menu() Command {
	return NewCommand("menu", func(args ...string) ([]byte, error) {
		resturant := []string{"angel", "hisa", "menza"}

		// check amount of arguments, if restaurant is missing return list
		if len(args) < 2 {
			return []byte(strings.Join(resturant, " ")), nil
		}

		out := []string{}
		switch args[1] {
		case "angel":
			out = angel()
		case "hisa":
			out = dobrahisa()
		case "menza":
			out = menza()
		default:
			out = append(out, "Don't know "+args[1]+"!")
		}
		return []byte(strings.Join(out, "\n")), nil
	})
}

func angel() []string {
	r := []string{}

	// url for Dobra Hisa
	url := "http://www.kaval-group.si/ANGEL,,ponudba/kosila"

	// load the URL
	res, err := http.Get(url)
	if err != nil {
		r = append(r, err.Error())
	}
	defer res.Body.Close()

	// convert windows-1250 HTML to utf-8 encoded HTML
	utfBody, err := iconv.NewReader(res.Body, "windows-1250", "utf-8")
	if err != nil {
		r = append(r, err.Error())
	}

	// get document
	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		r = append(r, err.Error())
		return r
	}

	// add title
	r = append(r, "Ponudba v Piceriji Angel:")

	// get right class selector from day of week
	dow := int(time.Now().Weekday()) - 1
	cls := ".show-" + strconv.Itoa(dow)

	// get div for first day
	doc.Find(cls).Each(func(i int, d *goquery.Selection) {
		d.Find("p").Each(func(i int, p *goquery.Selection) {
			html, err := p.Html()
			if err != nil {
				r = append(r, err.Error())
			}

			// manipulate html for nicer output
			html = strings.Replace(html, "<strong>", "*", -1)
			html = strings.Replace(html, "</strong>", "*", -1)
			html = strings.Replace(html, "* *", "", -1)
			html = strings.Replace(html, "* *", "", -1)
			html = strings.Replace(html, "<br/>", "", -1)
			html = strings.TrimSpace(html)

			// append menu to output
			if len(html) > 0 {
				r = append(r, html)
			}
		})
	})

	return r
}

func dobrahisa() []string {
	r := []string{}

	// url for Dobra Hisa
	url := "https://api.malcajt.com/getApiData.php?action=embed&id=2030&show=1001"

	// get document
	doc, err := goquery.NewDocument(url)
	if err != nil {
		r = append(r, err.Error())
		return r
	}

	// find first day
	doc.Find("a").Each(func(i int, a *goquery.Selection) {
		h, ok := a.Attr("href")
		if ok && h == "#day0" {
			r = append(r, "Ponudba v Dobri Hisi: *"+a.Text()+"*")
		}
	})

	// get div for day0
	doc.Find("#day0").Each(func(i int, d *goquery.Selection) {
		html, err := d.Html()
		if err != nil {
			r = append(r, err.Error())
		}

		// manipulate html for nicer output
		html = strings.Replace(html, "—", "", -1)
		html = strings.Replace(html, "<br/></i></b>", "</i></b><br/>", -1)
		html = strings.Replace(html, "<b><i>", "*", -1)
		html = strings.Replace(html, "</i></b>", "*", -1)

		// split by lines
		s := strings.Split(html, "<br/>")

		// append menu to output
		r = append(r, s...)
	})

	return r
}

func menza() []string {
	r := []string{}

	// url for Menza
	url := "https://www.studentska-prehrana.si/sl/restaurant/Details/2710"

	// get document
	doc, err := goquery.NewDocument(url)
	if err != nil {
		r = append(r, err.Error())
		return r
	}

	// add title
	r = append(r, "Ponudba v Menzi:")

	doc.Find("#menu-list").Find(".shadow-wrapper").Each(func(i int, div *goquery.Selection) {
		s := div.Find("h5").Find("strong").Text()

		s = strings.ToLower(s)
		split := strings.Split(s, " ")
		// capitalize first word. for some reason first split element is empty string
		split[1] = strings.Title(split[1])
		s = strings.Join(split, " ")
		r = append(r, s)
	})

	r = append(r, "© 2016 Študentska organizacija Slovenije – Vse pravice pridržane")

	return r
}
