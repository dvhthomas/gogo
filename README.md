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

### Set up

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

```sh
docker exec -it gogo_mariadb_1 bash
mysql -u web -p
******** # the $MYSQL_PASSWORD value
mysql> quit
```

Load the initial DB schema:

```sh
mysql -h 0.0.0.0 -u web -D snippetbox -p < pkg/models/mysql/schema.sql
```

## Teardown

If you need a fresh start:

```sh
docker-compose down; docker-compose rm
# and if you want to zap the db files...
rm -rf db/tmp
```

## Iterative development

### Database schema

Make changes to the database schema in `pkg/models/mysql/schema.sql` then apply to the DB. From your dev machine:

```sh
mysql -h 0.0.0.0 -u web -D snippetbox -p < pkg/models/mysql/schema.sql
```

### Test data

Keep any test data that you might need for testing in the `pkg/models/mysql/test_data.sql` file and load as follows:

```sh
mysql -h 0.0.0.0 -u web -D snippetbox -p < pkg/models/mysql/schema.sql
```

## Appendix - For posterity

### Maria DB setup

This is all one-time setup for a local dev environment.

Following the [MariaDB](https://hub.docker.com/_/mariadb) docker page and [this handy guide](https://towardsdatascience.com/connect-to-mysql-running-in-docker-container-from-a-local-machine-6d996c574e55).

You have to pass a fully qualified (not relative) path for the data dir. If you're already in the proejct root directory then `$PWD/data` should suffice:

```sh
docker run --name demo-db -p 3306:3306 -v $PWD/data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=secret -d mariadb:latest
```

Now connect to the container-based MariaDB and make some one-time config changes.

```sh
docker exec -it demo-db bash
$ mysql -u root -p
******* # 'secret'
```

Now try from your local machine:

```sh
mysql -h 0.0.0.0 -u root -p
******** # 'secret'
mysql> quit;
```

## Create the database and tables




## Appendix - DB steps you might need

### Enable external connections

By default MariaDB only lets local connections work. In this context, 'local' means from within the container itself. Let's enable any connection for the test db, including from the local deve machine. Use the [instructions from the MariaDB docker page](https://mariadb.com/kb/en/installing-and-using-mariadb-via-docker/#connecting-to-mariadb-from-outside-the-container) for this:

We're going to need an editor:

```sh
apt-get update
apt-install vim -y
```

```sql
mysql> update mysql.user set host = '%' where user='root';
Query OK, 1 row affected (0.02 sec)
```

Obviously don't forget to update this with new users as needed (this is just for root).

### Connect to the DB from local machine

First get the IP address of the container:

```sh
export DBIP=$(docker inspect --format '{{ .NetworkSettings.IPAddress }}' demo-db)
$ echo $DBIP
172.17.0.2
```

Now install a mysql client and connect. On Mac using homebrew there are a couple of one-time steps to make sure the mysql client is correctly configured:

$ mysql -h $DBIP -u root -p
******
mysql > quit;
```
