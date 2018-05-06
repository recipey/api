# API

## Dependency management
Using `dep` to manage packages. Ran `dep init` to set up project with lockfile and vendor
directory. When you want to add new packages run `dep ensure -add git_url_1 git_url_2 ... n` to
install however many packages you want 1 to many. To get an overview of how dependency looks
like run `dep status` to see packages and several attributes like constraint and versioning.

Update dependencies running `dep ensure -update`.

## Schema management
Using mattes/migrate to build timestamped sql migrations and to follow some sort of convention.
When running migrations assumed to be in the api container. Before running the commands below
enter the api container bash with `docker-compose exec api bash`.

### Create migration
`bin/migrate create -ext sql -dir migrations <name_of_migration>`

### Run migration
Run up migration
`bin/migrate -path migrations/ -database postgres://recipey:recipey@db:5432/recipey_dev?sslmode=disable up`

Run down migration
`bin/migrate -path migrations/ -database postgres://recipey:recipey@db:5432/recipey_dev?sslmode=disable down`

## Development
Start the app by running `docker-compose up` and when ready to test code changes recompile the app
by running `docker-compose restart api`. As part of the `docker-compose.yml` file it will rebuild
the app before running the binary again.
