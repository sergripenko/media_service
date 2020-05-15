# media_service
## 1: Prepare project env
#### Create app database:
####  Install PostgreSQL:
```bash
sudo apt-get update

sudo apt-get install postgresql
```
#### Or follow the installation instructions for PostgreSQL from the official website:
https://www.postgresql.org/download
#### Check installation:
```bash
sudo -u postgres psql
```
```bash
CREATE DATABASE media_service;
```
#### Or exit from postgres command line and create db with linux command:
```bash
\q

sudo -u postgres createdb media_service
```
#### Enter to the database:
```bash
sudo -u postgres psql -d media_service
```
#### Create user of the db from postgres command line:
```bash
CREATE USER admin WITH PASSWORD '123';

\q
```
#### Go to the postgres command line and do some customization:
```bash
sudo -u postgres psql

ALTER ROLE admin SET client_encoding TO 'utf8';

ALTER ROLE admin SET default_transaction_isolation TO 'read committed';

ALTER ROLE admin SET timezone TO 'UTC';

GRANT ALL PRIVILEGES ON DATABASE media_service TO admin;

ALTER USER admin CREATEDB;

\q
```
## 2: Getting Started
#### Inside media_service/db directory create "dbconf.yml" file with such content:
```bash
development:
    driver: postgres
    open: user=admin dbname=media_service password=123 host=127.0.0.1 port=5432 sslmode=disable
```
#### Install Goose migration tool:
https://bitbucket.org/liamstask/goose/src/master/
```bash
go get bitbucket.org/liamstask/goose/cmd/goose
```

#### Inside the project root directory migrate the database:
```bash
goose up
```
#### To delete tables, repeat below command for each table:
```bash
goose down
```
#### To install all dependencies run:
```bash
dep ensure -v
```

#### To run app:
```bash
bee run -downdoc=true -gendoc=true
```

#### Swagger docs:
```bash
http://127.0.0.1:8080/swagger/
```
