package article

import (
	"crud/util"
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/google/uuid"
	"log"
)

const articleBucket string = "articles"

var db *bolt.DB

func InitDatabase() {
	dbt, err := bolt.Open("my.db", 0600, nil)
	db = dbt // thanks, go
	util.Panicerr(err)

	log.Print(db.Info().Data)
	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(articleBucket))
		util.Panicerr(err)

		var id string
		id = uuid.New().String()

		now := util.NowUtc()
		email := "tester@rago.com"
		text := "Random text for the first entry"
		title := "Title?"
		article := Article{Uuid: &id, DatePublished: &now, Publisher: &email, Text: &text, Title: &title}
		bytes, err := json.Marshal(article)
		util.Panicerr(err)

		key, _ := b.Cursor().First()
		if key == nil {
			b.Put([]byte(id), bytes)
			get := b.Get([]byte(id))
			log.Print("Successfully put into database ", id, string(get))
		}

		return nil
	})
}

func getAllArticles() []Article {
	var articles []Article

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(articleBucket))
		bucket.ForEach(func(key, value []byte) error {
			if value == nil {
				return nil
			}

			var tmp *Article
			err := json.Unmarshal(value, &tmp)
			util.Panicerr(err)

			articles = append(articles, *tmp)

			return nil
		})
		return nil
	})
	util.Panicerr(err)
	log.Print("Received from db: ", articles)

	return articles
}

func getAllArticlesNotDeleted() []Article {
	var articles []Article

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(articleBucket))
		bucket.ForEach(func(key, value []byte) error {
			if value == nil {
				return nil
			}

			var tmp *Article
			err := json.Unmarshal(value, &tmp)
			util.Panicerr(err)

			if tmp.IsDeleted == nil || !*tmp.IsDeleted {
				articles = append(articles, *tmp)
			}

			return nil
		})
		return nil
	})
	util.Panicerr(err)
	log.Print("Received from db: ", articles)

	return articles
}

func getArticleById(id string) *Article {
	var article *Article
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(articleBucket))
		articleFromDb := bucket.Get([]byte(id))
		err := json.Unmarshal(articleFromDb, &article)
		util.Panicerr(err)
		return nil
	})
	util.Panicerr(err)
	log.Print("Received from db: ", article)

	return article
}

func updateArticle(article Article) *Article {
	var articleInserted *Article
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(articleBucket))

		articleJson, err := json.Marshal(article)
		util.Panicerr(err)

		bucket.Put([]byte(*article.Uuid), articleJson)
		articleFromDb := bucket.Get([]byte(*article.Uuid))
		er := json.Unmarshal(articleFromDb, &articleInserted)
		util.Panicerr(er)
		return nil
	})
	util.Panicerr(err)
	log.Print("Received from db after adding: ", articleInserted)

	return articleInserted
}

func bulkUpdateArticles(articles []Article) {
	for i := 0; i < len(articles); i++ {
		updateArticle(articles[i])
	}
}

func deleteArticleById(id string) {

}
