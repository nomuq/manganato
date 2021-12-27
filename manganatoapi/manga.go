package manganatoapi

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gocolly/colly"
)

type Manga struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Image        string    `json:"image"`
	Alternatives string    `json:"alternatives"`
	Author       Author    `json:"author"`
	Status       string    `json:"status"`
	Updated      string    `json:"updated"`
	Views        string    `json:"views"`
	Rating       string    `json:"rating"`
	Description  string    `json:"description"`
	Genres       []Genre   `json:"genres"`
	Chapters     []Chapter `json:"chapters"`
}

func getMangaList(name string) []Manga {
	url := searchMangaURL + name
	return getMangaListHelper(url)
}

// helper of SearchLatestUpdatedManga,
func getLatestUpdatedManga() []Manga {
	mgs := []Manga{}

	c.OnHTML(".content-homepage-item", func(h *colly.HTMLElement) {
		m := Manga{}
		m.getID(h.ChildAttr(".content-homepage-item-right h3 a", "href"))
		m.Name = h.ChildText(".content-homepage-item-right h3 a")
		m.Author.Name = h.ChildText(".content-homepage-item-right .item-author")
		m.Image = h.ChildAttr(".img-loading", "src")

		// h.child(".content-homepage-item-right .item-chapter")
		// get chapters

		h.ForEach(".content-homepage-item-right .item-chapter", func(_ int, el *colly.HTMLElement) {
			ch := Chapter{}
			ch.getID(el.ChildAttr("a", "href"))
			ch.ChapterName = el.ChildText("a")
			ch.Views = el.ChildText("span")
			ch.Uploaded = el.ChildText("span")

			ch.MangaID = m.ID
			ch.Uploaded = el.ChildText("i")

			m.Chapters = append(m.Chapters, ch)
		})

		// m.Chapters = append(m.Chapters, Chapter{
		// 	// Name: h.ChildText(".content-homepage-item-right .item-chapter"),
		// })

		mgs = append(mgs, m)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", e)
	})

	c.Visit(manganatoURLWithHTTPS)
	c.Wait()

	return mgs
}

func (m *Manga) getMangaByID() {

	c.OnHTML(".story-info-right", func(h *colly.HTMLElement) {
		name := h.ChildText("h1")
		m.Name = name
	})

	c.OnHTML(".variations-tableInfo", func(h *colly.HTMLElement) {
		alternatives := h.ChildText("tr:nth-child(1) .table-value")
		status := h.ChildText("tr:nth-child(3) .table-value")

		m.Alternatives = alternatives
		m.Status = status
	})

	c.OnHTML(".story-info-right-extent", func(h *colly.HTMLElement) {
		updated := h.ChildText("p:nth-child(1) .stre-value")
		views := h.ChildText("p:nth-child(2) .stre-value")

		m.getMangaRating(h.ChildText("em#rate_row_cmd"))
		m.Updated = updated
		m.Views = views
	})

	c.OnHTML(".panel-story-info-description", func(h *colly.HTMLElement) {
		m.getMangaDescription(h.Text)
	})

	createGenreList(m, c)
	createChapterList(m, c)
	createAuthor(m, c)

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", e)
	})

	c.Visit(specificMangaURL + m.ID)
	c.Wait()
}

// formats manga description
func (m *Manga) getMangaDescription(desc string) {
	pref := "Description :\n"
	desc = strings.Trim(desc, "\n")
	m.Description = strings.TrimPrefix(desc, pref)
}

// formats manga rating
func (m *Manga) getMangaRating(rating string) {
	tmp := (strings.Fields(rating))[3:]
	m.Rating = strings.Join(tmp, " ")
}

func (m *Manga) getID(url string) {
	m.ID = getID(url, "-")
}

func (m *Manga) compareManga(tmp *Manga) bool {
	return reflect.DeepEqual(m, tmp)
}
