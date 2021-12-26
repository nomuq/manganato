package manganatoapi

import (
	"fmt"

	"github.com/gocolly/colly"
)

type Chapter struct {
	ID          string `json:"id"`
	MangaID     string `json:"manga_id"`
	ChapterName string `json:"chapter_name"`
	Views       string `json:"views"`
	Uploaded    string `json:"uploaded"`
	Pages       []Page `json:"pages"`
}

// gets every chapter info for a specific manga
func createChapterList(m *Manga, c2 *colly.Collector) {

	c2.OnHTML(".row-content-chapter li.a-h", func(h *colly.HTMLElement) {
		ch := Chapter{}
		ch.getID(h.ChildAttr("a.chapter-name", "href"))
		ch.MangaID = m.ID
		ch.ChapterName = h.ChildText("a.chapter-name")
		ch.Views = h.ChildText("span.chapter-view")
		ch.Uploaded = h.ChildText("span.chapter-time")

		m.Chapters = append(m.Chapters, ch)

	})
}

// helper of ReadMangaChapter
func (ch *Chapter) getChapterByID() {
	ch.Pages = createPages(ch.getChapterURL())
}

func (ch *Chapter) getID(url string) {
	ch.ID = getID(url, "-")
}

func (ch *Chapter) getChapterURL() string {
	return fmt.Sprintf("%s%s/chapter-%s", specificMangaURL, ch.MangaID, ch.ID)
}
