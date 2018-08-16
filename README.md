# A simple Go URL Shortener

This is a simple URL shortening API that I built without any previous experience in Go. Built with ♥ in a solid day of hacking.

## Features
- centralized JSON configuration file (accessible as middleware by request handlers)
- same short code served for the same links submitted
- SQLite integration

## Endpoints
### `/`
#### GET
Display the number of rows in the table (i.e. the number of link to short code combinations)

#### POST
Submit a new URL to be shortened by the API. This will take a JSON object in the request body with the format:
```json
{
  "url": "http://dmuhs.com"
}
```
and yield a response
```json
{
    "LongURL": "http://dmuhs.com/",
    "ShortURL": "http://my-shortener.com:8000/ByH81bL0"
}
```

### `/{code}`
#### GET
Delivers a 302 permanent redirect response towards the URL associated with the short code in the database.


## Configuration
Configuration is handled through a JSON file in `url-shortener/config.json`.
```json
{
  "SQLitePath": "./urls.db",
  "Port": 8000,
  "ShortCodeLength": 8,
  "Domain": "http://my-shortener.com",
  "Public": true
}
```

Here the `Public` key toggles whether the service is accessible from hosts other than `localhost`. The `shortCodeLength` property defines the length of the generated short codes, which are randomly chosen from the charset `a-zA-Z0-9`. The `Domain` property decides which domain is used by the shortener's response, e.g. `http://my-shortener.com:8000/ByH81bL0` for the above example configuration. To be able to directly follow links in a local deployment, it can also be set to `http://localhost`.
