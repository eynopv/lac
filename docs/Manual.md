# Manual

`gorcli` supports Yaml and JSON files, but this manual only uses JSON.

## Requests

Create request:
```json
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
gorcli run createPet.json
```

## Variables and .env

`gorcli` can use paased variables and preload `.env`


```json
// variables.json

{
    "petId": "id-1"
}

```

```sh
// .env

API_KEY=mysecretkey
```

To use variables in request, write it surrounded by `${` and `}`.

```json
// getPet.json

{
    "method": "GET",
    "path": "https://api.example.com/v1/pet/${petId}",
    "headers": {
        "x-api-key": "${API_KEY}"
    }
}
```

Then pass variables to `gorcli`

```sh
gorcli run getPet.json --vars variables.json
```

## Headers

`gorcli` can use passed headers and headers can contain variables.

```json
{
    "Content-Type": "application/json",
    "x-api-key": "${API_KEY}"
}
```

Pass headers and variables to `gorcli`

```sh
gorcli run getPet.json --vars variables.json --headers headers.json
```

## Tests

```json
// createAndGetPet.json

[
  {
    "id": "createPet",
    "uses": "createPet.json"
    "expect": {
      "status": 201,
      "timeLessThan": 3000
    }
  },
  {
    "uses": "getPet.json"
    "with": {
      "petId": "${createPet.id}"
    },
    "expect": {
      "status": 200,
      "timeLessThan": 1000
    }
  }
]
```
