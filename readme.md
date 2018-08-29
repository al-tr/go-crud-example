# Crud example 

[![Build Status](https://travis-ci.com/al-tr/playing-with-go.svg?branch=master)](https://travis-ci.com/al-tr/playing-with-go) [![Go Report Card](https://goreportcard.com/badge/github.com/al-tr/playing-with-go)](https://goreportcard.com/report/github.com/al-tr/playing-with-go) [![Open Source Helpers](https://www.codetriage.com/al-tr/playing-with-go/badges/users.svg)](https://www.codetriage.com/al-tr/playing-with-go)

Handlers map:

| Method | Endpoint | Description |
| --- | --- | --- |
| GET | /articles(/) | Gets all articles that don't have isDeleted=true |
| GET | /articles/all | Gets all articles (priority over {id}) |
| GET | /articles/{id} | Gets article by id |
| PUT | /articles(/) | Creates (no id in body / id in body and id not found) or updates (id in body and id found) an article from the post body, validates all the stuff before update | 
| DELETE | /articles(/) | Cleans the database from bad articles: marks them as isDeleted=true |
| DELETE | /articles/{id} | Marks article as isDeleted=true |
