# Image API

To run, copy `.env-example` to `.env`, and run `docker-compose up`.
To run without docker, point `.env` to the db, and run `make run`.

`Image API.postman_collection.json` contains API calls that can be imported to Postman.

## Stuff that could be improved

- Add tests
- Proper db migration (currently just creates a table when launching the API)
- Ability to configure the API endpoint address/port
