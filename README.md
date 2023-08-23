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
