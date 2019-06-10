# API

## Dependency management
Just use `go mod`. `go get -u` to upgrade.

## Schema management
Using `golang-migrate/migrate` to build timestamped sql migrations and to follow some sort of convention.
When running migrations docker-compose is assumed to be running. << IMPORTANT >>

Enter container shell and run additional commands.
`docker-compose exec api bash`

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

## PSQL
Find a GUI for this later. Currently using just the CLI. The service container for postgres
must be running in order to connect.
```
docker-compose exec db bash
psql -U recipey -d recipey_dev

```
If your host has a compatible psql client then this should work as well.
`psql -h localhost -U recipey -d recipey_dev`

## Development
Start the app by running `docker-compose up` and when ready to test code changes recompile the app
by running `docker-compose restart api`. As part of the `docker-compose.yml` file it will rebuild
the app before running the binary again.
