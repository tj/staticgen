package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Post model.
type Post struct {
	ID     int
	Title  string
	Teaser string
	Body   string
}

var posts []Post

func init() {
	for i := 0; i < 10000; i++ {
		posts = append(posts, Post{
			ID:     i,
			Title:  fmt.Sprintf("Blog post #%d", i),
			Teaser: "Some content here.",
			Body:   "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse venenatis sem sit amet felis imperdiet porta. Cras eget cursus ante. Suspendisse lobortis aliquam felis feugiat vehicula. Phasellus id mattis lacus. Morbi pellentesque nibh ut turpis finibus porttitor. Vivamus vestibulum dui vel mi sagittis dictum. Ut faucibus commodo magna vel auctor. Donec ut egestas nunc, at ullamcorper mi. Maecenas in fringilla mauris. Mauris non enim eu diam elementum ornare quis at est. Morbi et semper ex, vitae cursus lectus. Nulla auctor, nibh vel tempus ultrices, erat metus auctor ligula, vel tincidunt quam quam ac ligula. Nullam in elit dui. In tellus quam, suscipit eu ultrices sit amet, volutpat eu tellus. Integer pretium et dolor eu faucibus. Proin aliquam blandit rhoncus.",
		})
	}
}

func main() {
	http.HandleFunc("/", list)
	http.Handle("/posts/", http.StripPrefix("/posts/", http.HandlerFunc(post)))
	fmt.Printf("Server starting on :3000\n")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func list(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	for _, p := range posts {
		fmt.Fprintf(w, `
			<article>
				<h2>%s</h2>
				<p>%s</p>
				<a href="/posts/%d">View article</a>
			</article>
		`, p.Title, p.Teaser, p.ID)
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.URL.Path, 10, 64)
	p := posts[id]
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<h1>%s</h1>
		<p>%s</p>
	`, p.Title, p.Body)
}
