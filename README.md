Welcome to TeamsACS project!

      _____                                    ___     ___     ___   
     |_   _|   ___    __ _    _ __     ___    /   \   / __|   / __|  
       | |    / -_)  / _` |  | '  \   (_-<    | - |  | (__    \__ \  
      _|_|_   \___|  \__,_|  |_|_|_|  /__/_   |_|_|   \___|   |___/  
    _|"""""|_|"""""|_|"""""|_|"""""|_|"""""|_|"""""|_|"""""|_|"""""| 
    "`-0-0-'"`-0-0-'"`-0-0-'"`-0-0-'"`-0-0-'"`-0-0-'"`-0-0-'"`-0-0-' 

# TeamsACS

TeamsACS is committed to providing superior ease of network management for work teams. We use Mikrotik's networking products as the base core.

The core of the system is based on Golang technology, providing superior performance and an easy deployment experience.

![TeamsACS Struct](https://user-images.githubusercontent.com/377938/166147509-c5df9824-52f1-43c3-ae46-842a1cbe9400.png)

## Install and config

TeamsACS uses PostgreSQL as its primary database and uses the Timescaledb extension

- Create Database

    CREATE USER teamsacs WITH PASSWORD 'teamsacs'
    
    CREATE DATABASE teamsacs OWNER teamsacs;
    
    GRANT ALL PRIVILEGES ON DATABASE teamsacs TO teamsacs;

## Docker Deploy

```yml
version: "2"
services:
  pgdb:
    image: timescale/timescaledb:latest-pg14
    container_name: "pgdb"
    ports:
      - "127.0.0.1:5432:5432"
    environment:
      POSTGRES_DB: teamsacs
      POSTGRES_USER: teamsacs
      POSTGRES_PASSWORD: teamsacs
    volumes:
      - /data/pgdata:/var/lib/postgresql/data
    networks:
      teamsacs_network:

  teamsacs:
    depends_on:
      - 'pgdb'
    image: ca17/teamsacs:latest
    container_name: "teamsacs"
    restart: always
    ports:
      - 8000
      - 8106
    expose:
      - 8000
      - 8106
    volumes:
      - /data/teamsacs:/var/teamsacs
    environment:
      - GODEBUG=x509ignoreCN=0
      - TEAMSACS_WEB_DEBUG=1
      - TEAMSACS_SECRET=9b6de5cc-xxxx-xxxx-xxx-0f568ac9da37
    networks:
      teamsacs_network:

networks:
  teamsacs_network:

```
