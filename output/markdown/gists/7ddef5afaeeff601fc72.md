# [REST design: PUT vs POST] #REST #API #PUT #POST #DELETE #GET #methods #http

## rest-design-put-vs-post.md

```markdown
**UPDATE Dec 2020**

HTTP Method | Description
------------|------------
`GET`       | To _read_ a single or collection resource.
`POST`      | To _create_ a resource.
`PUT`       | To _entirely replace_ a resource.
`PATCH`     | To _partially_ update a resource.
`DELETE`    | To _delete_ a resource.

> See more at: [restcookbook.com](http://restcookbook.com/HTTP%20Methods/put-vs-post/#sthash.uH9tVV9x.dpuf), [restapitutorial.com](https://www.restapitutorial.com/lessons/httpmethods.html) and this blog post on [best practices](https://hackernoon.com/restful-api-designing-guidelines-the-best-practices-60e1d954e7c9).

## PUT

Use PUT when you can update a resource completely through a specific resource. For instance, if you know that an article resides at http://example.org/article/1234, you can PUT a new resource representation of this article directly through a PUT on this URL.

## POST

If you do not know the actual resource location, for instance, when you add a new article, but do not have any idea where to store it, you can POST it to an URL, and let the server decide the actual URL.

> Note: there must always be a body param passed to the server when executing a POST or a PUT operation. Meaning you only pass enough information in the URL to identify the resource and not including any data for that resource. For example, this is bad: `/user/<id>/username/<new_username>` but this is ok: `/user/<id>/username` upon which you'll then pass the new username as the body param: `curl --method PUT --header 'Content-Type: application/json' --data '{"username": "foobar"}' https://api.example.com/user/id/123/username`.
```

