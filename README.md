apimaker
========

apimaker is a service to scrape websites and provide the scraped data back to
you as easy to parse JSON.

**Still a work in progress**

Requirements
------------

- [Selenium](http://www.seleniumhq.org/download/)

Installation
------------

    $ go install https://github.com/kanemathers/apimaker

Example Usage
-------------

Start Selenium:

    $ java -jar selenium-server-standalone-*.jar

Start apimaker. See ``apimaker --help`` for more options:

    $ apimaker

Start a python HTTP server to host the demo page:

    $ cd ./hacking
    $ python -m SimpleHTTPServer

### Adding a Job

    $ cd hacking
    $ curl -X POST -d @fields.json http://127.0.0.1:8080/jobs
    {"id":"5e772c1f-96e0-11e4-b6fa-8c705a804d80"}

### Retrieve a Job

    $ curl http://127.0.0.1:8080/jobs/5e772c1f-96e0-11e4-b6fa-8c705a804d80 2>/dev/null | python -m json.tool
    {
        "collections": {
            "news": {
                "group": ".news-post",
                "selectors": [
                    {
                        "name": "title",
                        "selector": ".title"
                    },
                    {
                        "name": "content",
                        "selector": ".content"
                    }
                ]
            },
            "page_title": {
                "group": null,
                "selectors": [
                    {
                        "name": "page_title",
                        "selector": "body > h1"
                    }
                ]
            }
        },
        "id": "5e772c1f-96e0-11e4-b6fa-8c705a804d80",
        "interval": 3600000000000,
        "last_scraped": "2015-01-08T10:45:15.022743202+08:00",
        "scraped_data": {
            "news": [
                [
                    {
                        "title": "First News Post"
                    },
                    {
                        "content": "The content of the first news post"
                    }
                ],
                [
                    {
                        "title": "And the Second Post!"
                    },
                    {
                        "content": "The second news post has boring content. Like the first."
                    }
                ]
            ],
            "page_title": [
                [
                    {
                        "page_title": "Demo Page"
                    }
                ]
            ]
        },
        "url": "http://127.0.0.1:8000/demo.html"
    }

### Removing a Job

    $ curl -X DELETE http://127.0.0.1:8080/jobs/5e772c1f-96e0-11e4-b6fa-8c705a804d80

