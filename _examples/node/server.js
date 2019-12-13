
const express = require('express')
const app = express()

app.set('view engine', 'pug')
app.use(express.static('public'))

// generate fake posts, normally these would come
// from disk, or a database.
const posts = []

for (let i = 0; i < 1000; i++) {
  posts.push({
    id: i,
    title: `Blog Post #${i}`,
    teaser: 'Some content here.',
    body: `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse venenatis sem sit amet felis imperdiet porta. Cras eget cursus ante. Suspendisse lobortis aliquam felis feugiat vehicula. Phasellus id mattis lacus. Morbi pellentesque nibh ut turpis finibus porttitor. Vivamus vestibulum dui vel mi sagittis dictum. Ut faucibus commodo magna vel auctor. Donec ut egestas nunc, at ullamcorper mi. Maecenas in fringilla mauris. Mauris non enim eu diam elementum ornare quis at est. Morbi et semper ex, vitae cursus lectus. Nulla auctor, nibh vel tempus ultrices, erat metus auctor ligula, vel tincidunt quam quam ac ligula. Nullam in elit dui. In tellus quam, suscipit eu ultrices sit amet, volutpat eu tellus. Integer pretium et dolor eu faucibus. Proin aliquam blandit rhoncus.`,
  })
}

app.get('/posts/:id', (req, res) => {
  const post = posts[req.params.id]
  res.render('post', { post })
})

app.get('/', (req, res) => {
  res.render('index', { posts })
})

console.log('Server starting on :3000')
app.listen(3000)