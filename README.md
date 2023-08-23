# Open Bar
Open Bar is a project that combines Open Source, Open Data, and Open Hardware to allow anyone to build their own fully automated
cocktail machine. The Open Bar source and the Cocktails database are the first two pieces of this project to be made available publicly.
There will be additional components published in the future including a client application, and different robot bartender hardware files
such as the EAGLE files for PCB creation, and different models for robot bartenders and their components which can be used for creating your
own robot bartender.

![Robot Bartender](https://raw.githubusercontent.com/cocktailrobots/openbar-server/main/rb.jpeg)

# openbar-server
openbar-server provides two REST APIs written in Go. 

* The Cocktails API allows you to add, edit, and delete cocktails, ingredients, and recipes. The
* The OpenBar API provides administrative functions like calibrating the pumps, telling the system about what ingredients are in each bottle,
and of course making drinks.

The server is backed by a [Dolt](https://www.doltdb.com/) database, and uses [DoltHub](https://www.dolthub.com/)
to manage the database of cocktail recipes, and the API provides the means of downloading new recipes from a fork of
the [cocktails database](https://www.dolthub.com/repositories/openbar/cocktails) on [DoltHub](https://www.dolthub.com/).
Since you manage your own fork of the database, you can add your own recipes, and create pull requests on [Dolthub](https://www.dolthub.com/)
to share your recipes with the world.

# Setup
The setup instructions for installing this onto an actual cocktail robot will be different from those used to setup the project for development. Below are the
directions for setting up this project for local development.

## Prerequisites
* [Go](https://golang.org/doc/install) - openbar-server is written in Go, so you'll need to install Go to build and run the project.
* [Dolt](https://www.doltdb.com) - openbar-server uses Dolt as its database, so you'll need to install Dolt to run the project.
* [Dolthub Account](https://www.dolthub.com) - openbar-server uses Dolthub to manage the cocktail recipes database, so you'll need to create an account on Dolthub so you can fork the data.

## Install The Open Bar Server
Once you've cloned the repository, install openbar-server by opening a terminal prompt and running 

```bash
go install .
```

in the `cmd/openbar-server` directory of the project.

## Fork and Clone the Cocktails DB
The cocktails database is hosted on [Dolthub](https://www.dolthub.com/repositories/openbar/cocktails), and you'll want 
to fork it to your own account so you can push changes to your cocktails database.  This will also allow you to create
pull requests to share your recipes with the world.

Once you've forked the database, you'll want to clone it to your local machine.  You should create a new directory to hold
your datasbases. From your new directory can clone your fork by running the following

```bash
dolt clone <owner>/cocktails
```

Where `<owner>` is your Dolthub username, org organization name that owns the fork.

## Create a server config for dolt sql-server

In the same folder you cloned the cocktails database to you will need to create a `config.yaml` file. This file configures
the dolt server and should look something like:

```yaml
log_level: "debug"
user:
  name: "openbar"
listener:
  host: "0.0.0.0"
  port: 3306
  max_connections: 5
```

For details on the configuration options see the [dolt documentation](https://docs.dolthub.com/cli-reference/cli#dolt-sql-server).

## Start dolt sql-server

Now that you have cloned your database fork, and created a `config.yaml` file, you can start the dolt sql-server by running

```bash
dolt sql-server --config config.yaml
```

from the directory where you cloned your database fork and created the config.yaml file.

## Create a debugconfig.yaml file for openbar-server

openbar-server uses a configuration file to tell it how to connect to the databases it uses.  You can create a debugconfig.yaml
file that looks like:

```yaml
hardware:
  debug:
    num-pumps: 8
    out-file: "/Users/brian/openbar.out"

db:
  host: "127.0.0.1"
  port: 3306
  user: openbar
  pass:

cocktails-api:
  host: "0.0.0.0"
  port: 8675

openbar-api:
  host: "0.0.0.0"
  port: 3099
```

Put it wherever you want to run openbar-server from.

## Creating the OpenBarDB
The second database you'll need is the OpenBarDB. This is the database that openbar-server uses to store its configuration,
calibration, and some other data related to the robot.  You can create this database by running the following commands.

```bash
openbar-server -migration-dir=<path to openbar-server source>/schema/openbardb debugconfig.yaml
```

## Running openbar-server

Now that you have set everything up you can run openbar-server by running

```bash
openbar-server debugconfig.yaml
```

And you can test it by running:

```bash
curl http://localhost:8675/cocktails
```

you should see output that looks like:

```json
[
  {
    "name":"americano",
    "display_name":"Americano",
    "description":"The Americano is an IBA official cocktail composed of Campari, sweet vermouth, and club soda."
  },
  {
    "name":"boulevardier",
    "display_name":"Boulevardier",
    "description":"A Boulevardier is a bourbon-based cocktail made with Campari, sweet vermouth, and bourbon whiskey.",
  },
  {
    "name":"gin_and_tonic",
    "display_name":"Gin and Tonic",
    "description":"A gin and tonic is a highball cocktail made with gin and tonic water poured over a large amount of ice."
  },
  {
    "name":"manhattan",
    "display_name":"Manhattan",
    "description":"A Manhattan is a cocktail made with whiskey, sweet vermouth, and bitters.",
  },
  {
    "name":"negroni",
    "display_name":"Negroni",
    "description":"A Negroni is an Italian cocktail, made of one part gin, one part vermouth rosso (red, semi-sweet) and one part Campari, garnished with orange peel. It is considered an apéritif.",
  {
    "name":"paper_plane",
    "display_name":"Paper Plane",
    "description":"The Paper Plane is a modern variation on the Last Word composed of Bourbon, Aperol, Amaro Nonino and Lemon Juice."
  }
]
```