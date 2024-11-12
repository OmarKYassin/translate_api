# Translate API

This is a simple API that receives a request containing JSON structured call log and translates the non English words into English.

---

## Quick start

To run this application make sure you have docker and docker compose on your machine then run:

``` bash
git clone https://github.com/OmarKYassin/translate_api.git
cd translate_api
docker-compose up
```

---

## Endpoint

### Path

``` http
POST "http://localhost:3000/translate"
```

## Body

```json
[
  { "speaker": "John", "time": "00:00:04", "sentence": "Hello Maria." },
  { "speaker": "Maria", "time": "00:00:09", "sentence": "صباح الخير" }
]
```

## Expected response

```json
[
  { "speaker": "John", "time": "00:00:04", "sentence": "Hello Maria." },
  { "speaker": "Maria", "time": "00:00:09", "sentence": "Good Morning" }
]
```

---
