package article

import (
	"crud/auth"
	"crud/util"
	"encoding/json"
	"errors"
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

func getArticles(includeDeleted bool) ([]Article, error) {
	var articles []Article

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(articleBucket))
		bucket.ForEach(func(key, value []byte) error {
			if value == nil {
				return nil
			}

			var tmp *Article
			err := json.Unmarshal(value, &tmp)
			if err != nil {
				return err
			}

			if includeDeleted {
				// all articles
				articles = append(articles, *tmp)
			} else {
				//excluding articles marked as deleted
				if tmp.IsDeleted == nil || !*tmp.IsDeleted {
					articles = append(articles, *tmp)
				}
			}

			return nil
		})
		return nil
	})
	if err != nil {
		return articles, err
	}
	log.Print("Received from db: ", articles)

	return articles, nil
}

func getAllArticles() ([]Article, error) {
	return getArticles(true)
}

func getAllArticlesNotDeleted() ([]Article, error) {
	return getArticles(false)
}

func getArticleById(id string) (*Article, error) {
	var article *Article
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(articleBucket))
		articleFromDb := bucket.Get([]byte(id))
		if articleFromDb == nil {
			return errors.New("No article found with " + id)
		}
		err := json.Unmarshal(articleFromDb, &article)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	log.Print("Received from db: ", article)

	return article, nil
}

func putArticle(article Article) (*Article, error) {
	var articleInserted *Article
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(articleBucket))

		articleJson, err := json.Marshal(article)
		if err != nil {
			return err
		}

		bucket.Put([]byte(*article.Uuid), articleJson)
		articleFromDb := bucket.Get([]byte(*article.Uuid))
		err = json.Unmarshal(articleFromDb, &articleInserted)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return &article, err
	}
	log.Print("Received from db after adding: ", articleInserted)

	return articleInserted, nil
}

func bulkPutArticles(articles []Article) ([]Article, []error) {
	var articlesResponses []Article
	var err []error
	for i := 0; i < len(articles); i++ {
		article, e := putArticle(articles[i])
		articlesResponses = append(articlesResponses, *article)
		err = append(err, e)
	}
	return articlesResponses, err
}

func deleteArticleById(id string, user auth.User) {

}
