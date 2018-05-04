# API

## Dependency management
Using `dep` to manage packages. Ran `dep init` to set up project with lockfile and vendor
directory. When you want to add new packages run `dep ensure -add git_url_1 git_url_2 ... n` to
install however many packages you want 1 to many. To get an overview of how dependency looks
like run `dep status` to see packages and several attributes like constraint and versioning.

Update dependencies running `dep ensure -update`.

## Running the app
We are using Docker to isolate runtime. TBD... docker-compose coming
