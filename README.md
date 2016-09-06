# README #

# Law Parser

The following repo attemps to show a golang file parser. Following
parser basics, the code looks for tags of interest, and joing them
in a Linked List structure, to be used later to extract the content
of interest. If you look for a much mature parser, take a look at:
(https://github.com/dotabuff/sange)[https://github.com/dotabuff/sange] 

This does not intend to showcase the best posible route to parse a
text file, but my attemp to learn golang basics.

I started learning golang because I became obsessed with performance
and scalability.

# Requeriments

As of Sept 5, 2016 I'm using:
+ Golang 1.6
+ Migrating from Sqlite3 to Postgres 9.5
+ Docker to run my postgres Instance
+ Migrating to Goose to better handle db versions

# Reminders

The following are commands I should remember to start my dev enviroment

To Start Docker Container

    docker-compose up -d --build

To stop

    docker-compose down


