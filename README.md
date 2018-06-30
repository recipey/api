# API

## Dependency management
Using `dep` to manage packages. Ran `dep init` to set up project with lockfile and vendor
directory. When you want to add new packages run `dep ensure -add git_url_1 git_url_2 ... n` to
install however many packages you want 1 to many. To get an overview of how dependency looks
like run `dep status` to see packages and several attributes like constraint and versioning.

Update dependencies running `dep ensure -update`.

## Schema management
Using `golang-migrate/migrate` to build timestamped sql migrations and to follow some sort of convention.
When running migrations docker-compose is assumed to be running. << IMPORTANT >>

Either docker-compose run to perform the command in the container or bash into the container
to run command.

`docker-compose run api migrate-create <name_of_migration>`

OR

`docker-compose exec api bash` to enter container shell and run additional commands.

### Create migration
`migrate-create <name_of_migration>`

### Run migration
Run up migration
`migrate-up <n>` where n is optional integer > 0

Run down migration
`migrate-down <n>` where n is optional integer > 0

### Bad migration
If you're migration failed then the last migration will have its version kept in the
schema_migrations table marking it as dirty. To fix this you must rewrite your migration to
have no errors and then run `migrate-force <migration_version>` to undirty the version. Then run
`migrate-down 1` to undo the migration. You can now run the migration up command.

## Development
Start the app by running `docker-compose up` and when ready to test code changes recompile the app
by running `docker-compose restart api`. As part of the `docker-compose.yml` file it will rebuild
the app before running the binary again.

## PSQL
For now been bashing into postgres container and running psql client to run some queries and quick
view of the database and table schemas. To do this run `docker-compose exec db bash` then connect
to psql with `psql -U postgres` as the postgres user. Then connect to the dev database with
`\c recipey_dev` when you are in the psql shell.
