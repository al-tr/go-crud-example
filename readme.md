# Crud example

Handlers map:

| Method | Endpoint | Description |
| --- | --- | --- |
| GET | /articles(/) | Get all articles that don't have isDeleted=true |
| GET | /articles/{id} | Get article by id |
| PUT | /articles(/) | Creates (no id in body / id in body and id not found) or updates (id in body and id found) an article from the post body, validates all the stuff before update | 
| DELETE | /articles(/) | Cleans the database from bad articles: marks them as isDeleted=true |
| DELETE | /articles/{id} | Marks article as isDeleted=true |
