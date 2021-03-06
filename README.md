# Web app

## Development environment setup

### Prerequisites

* Go
* MySQL _client_ on your local machine

    ```sh
    # On a Mac this should suffice
    $ brew install mysql-client
    ...
    echo 'export PATH="/usr/local/opt/mysql-client/bin:$PATH"' >> ~/.zshrc
    source ~/.zshrc
    ```

The app consists of a Go-based web server and a MariaDB (MySQL compatible) database server.

Set a project directory. Not really needed but it saves typing:

```sh
export PROJECT=~/projects/gogo
cd $PROJECT
```

#### Database server

The docker imager for MariaDB can automatically create a database and user for you. Choose some passwords for your `root` and `web` database users on your local environment.

```sh
export MYSQL_ROOT_PASSWORD='your-root-password'
export MYSQL_PASSWORD='your-user-password'
```

Look in the `docker-compose.yml` file and you'll see that the application database will be called `snippetbox` and the application user will be `web`.

Use docker to create a clean DB environment:

```sh
docker-compose up -d
```

And then try connecting to the database via the container. Below it shows `gogo-mariadb_1` but use `docker ps` to see what your container is called:

> REMINDER: The container can take a few seconds to start up. It is, after all
  spinning up an entire MySQL database and a boatload of files in your `./data/`
  directory. If you get errors connecting, first try waiting a while.

```sh
docker exec -it gogo_mariadb_1 bash
mysql -u web -p
******** # the $MYSQL_PASSWORD value
mysql> quit
```

If that worked it's time to load the initial DB schema:

```sh
mysql -h 0.0.0.0 -u web -D snippetbox -p < pkg/models/mysql/schema.sql
```

## App Server

Generate a TLS key pair so secure hosting works on the localhost. This assumes that Go is installed in a typical location. Check your system for that. The `generate_cert.go` utility is part of every Go distribution:

```sh
cd $PROJECT/tls
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --hos
t=localhost
```

## Teardown

If you need a fresh start:

```sh
docker-compose down; docker-compose rm
# and if you want to zap the db files...
rm -rf db/tmp
```

## Iterative development

First get the web server running. To avoid hardcoding user name and password for the database we'll set some environments up. **Do not** check that in to version control!

```sh
cd $PROJECT
export DBUSER='web'
export DBPASS='something-super-secure-like-password123'
export SESSION_SECRET=$(openssl rand -base64 32)
# include optional arg or default to port 4000
# The docker image host should be 0.0.0.0 and defaults to an empty value
$ go run ./com/web -help
$ go run ./cmd/web -secret=$SESSION_SECRET -dbpass=$DBPASS -dbuser=$DBUSER
INFO Etc
...
[Ctrl-C to kill]
```

If the DB connection works you'll see a message telling you. If not, you'll get an ERROR log.

### Database schema

Make changes to the database schema in `pkg/models/mysql/schema.sql` then apply to the DB. From your dev machine:

```sh
mysql -h 0.0.0.0 -u web -D snippetbox -p < pkg/models/mysql/schema.sql
```

### Test data

Keep any test data that you might need for testing in the `pkg/models/mysql/test_data.sql` file and load as follows:

```sh
mysql -h 0.0.0.0 -u web -D snippetbox -p < pkg/models/mysql/test_data.sql
```
