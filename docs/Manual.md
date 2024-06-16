# Manual

`lac` supports Yaml and JSON files, but this manual only uses JSON.

## Requests

Create request:

```javascript
// createPet.json

{
  "method": "POST",
  "path": "https://api.example.com/v1/pet",
  "headers": {
    "Content-Type": "application/json"
  },
  "body": {
    "name": "Princess",
    "type": "dog",
    "breed": "bulldog"
  }
}
```

Send request:

```sh
lac createPet.json
```

## Variables and .env

`lac` can use paased variables and preload `.env`

```javascript
// variables.json

{
  "petId": "id-1"
}
```

```sh
# .env

API_KEY=mysecretkey
```

To use variables in request, write it surrounded by `${` and `}`.

```javascript
// getPet.json

{
  "method": "GET",
  "path": "https://api.example.com/v1/pet/${petId}",
  "headers": {
    "x-api-key": "${API_KEY}"
  }
}
```

Then pass variables to `lac`

```sh
lac getPet.json --vars variables.json
```

## Headers

`lac` can use passed headers and headers can contain variables.

```javascript
// headers.json

{
  "Content-Type": "application/json",
  "x-api-key": "${API_KEY}"
}
```

Pass headers and variables to `lac`

```sh
lac getPet.json --vars variables.json --headers headers.json
```
