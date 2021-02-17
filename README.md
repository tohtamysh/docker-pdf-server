# PDF conversion server

POST request

```sh
127.0.0.1:8000
```

You can pass both a link to the page and the body of the page.

Request params
| Param | Description |
| ------ | ------ |
| url | page adress (eg. https://google.com) |
| html | string html page |
| landscape | pdf orientation (eg. true or null) |
