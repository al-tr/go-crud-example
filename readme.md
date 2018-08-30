# CRUD example in Go language

[![Build Status](https://travis-ci.com/al-tr/go-crud-example.svg?branch=master)](https://travis-ci.com/al-tr/go-crud-example) [![Go Report Card](https://goreportcard.com/badge/github.com/al-tr/go-crud-example)](https://goreportcard.com/report/github.com/al-tr/go-crud-example) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/324b9e1f987d4e0dbbb8e930e97c1b26)](https://www.codacy.com/app/al-tr/go-crud-example?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=al-tr/go-crud-example&amp;utm_campaign=Badge_Grade) [![Open Source Helpers](https://www.codetriage.com/al-tr/go-crud-example/badges/users.svg)](https://www.codetriage.com/al-tr/go-crud-example) [![Coverage Status](https://coveralls.io/repos/github/al-tr/go-crud-example/badge.svg?branch=master)](https://coveralls.io/github/al-tr/go-crud-example?branch=master)

### Endpoints:

| Method | Endpoint | Description |
| --- | --- | --- |
| GET | /articles(/) | Gets all articles that don't have isDeleted=true |
| GET | /articles/all | Gets all articles (priority over {id}) |
| GET | /articles/{id} | Gets article by id |
| POST | /articles(/) | Creates an article from the request body, validates all the stuff before creation | 
| PUT | /articles/{id} | Updates an article from the request body, validates all the stuff before update | 
| DELETE | /articles(/) | Cleans the database from bad articles: marks them as isDeleted=true |
| DELETE | /articles/{id} | Marks article as isDeleted=true |

### Challenges

- go report A+
- code coverage 100%
- makefile
- advise if you have any ideas
