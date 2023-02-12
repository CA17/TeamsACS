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

TeamsACS uses Timescaledb as the primary data store, supporting time-series data and providing support for system observability

- TR069-based device management

TeamsACS implements the TR069 protocol server, supports Mikrotik as a TR069 protocol client for secure access, 
and supports Mikrotik TR069 private features, supports configuration reading and modification, 
and supports uploading Mikrotik script files for execution.

- Device Backup Management

The TR069 protocol enables Mikrotik devices to be backed up regularly and then downloaded to the TeamsACS file system for unified management, 
and supports backup restoration via TR069.


### Northbound Interface

- Provide a unified API for various third-party management systems, based on the HTTPS Json protocol.
- Provide basic equipment information and status data query API, and data maintenance API.
- Provide various policy management APIs, such as firewall rules, routing tables, etc.

# Quick Start

TeamsACS uses PostgreSQL as its primary database and uses the Timescaledb extension

- Create Database

    CREATE USER teamsacs WITH PASSWORD 'teamsacs'
    
    CREATE DATABASE teamsacs OWNER teamsacs;
    
    GRANT ALL PRIVILEGES ON DATABASE teamsacs TO teamsacs;


- Install TeamsACS 

The following installation method will download and build the latest teamsacs version

```
go install github.com/ca17/teamsacs@latest

teamsacs -install

```

> If you want to download the compiled binaries directly, you can visit [Github Release](https://github.com/CA17/TeamsACS/releases)

- Config TeamsACS

Modifying configuration file [/etc/teamsacs.yml](https://github.com/CA17/TeamsACS/wiki/Configuration)

Start the service with the following commands

    systemctl enable teamsacs
    systemctl start teamsacs


# Docker Deploy

```yml
version: "3"
services:
  pgdb:
    image: timescale/timescaledb:latest-pg14
    container_name: "pgdb"
    ports:
      - "127.0.0.1:15432:5432"
    environment:
      POSTGRES_DB: teamsacs
      POSTGRES_USER: teamsacs
      POSTGRES_PASSWORD: teamsacs
    volumes:
      - pgdb-volume:/var/lib/postgresql/data
    networks:
      teamsacs_network:

  teamsacs:
    depends_on:
      - 'pgdb'
    image: ca17/teamsacs:latest
    container_name: "teamsacs"
    restart: always
    ports:
      - "2979:2979"
      - "2989:2989"
      - "2999:2999"
    volumes:
      - teamsacs-volume:/var/teamsacs
    environment:
      - GODEBUG=x509ignoreCN=0
      - TEAMSACS_DB_HOST=pgdb
      - TEAMSACS_DB_PORT=5432
      - TEAMSACS_DB_NAME=teamsacs
      - TEAMSACS_DB_USER=teamsacs
      - TEAMSACS_DB_PWD=teamsacs
      - TEAMSACS_WEB_DEBUG=1
    networks:
      teamsacs_network:

networks:
  teamsacs_network:

volumes:
  pgdb-volume:
  teamsacs-volume:

```

## Access web console

Open the browser and enter the URL: `http://your-ip:2979` or `https://your-ip:2989` 

The default username and password are `admin/teamsacs`

## Links

- [TeamsACS Documentation Wiki](https://github.com/ca17/teamsacs/wiki)
- [Mikrotik Tr069 best practices](https://wiki.mikrotik.com/wiki/Tr069-best-practices)
- [Mikrotik Tr069 client](https://wiki.mikrotik.com/wiki/Manual:TR069-client)
- [Mikrotik TR069 client supported parameter reference](https://wiki.mikrotik.com/tr069ref/current.html)


# Contribute

We welcome contributions of any kind, including but not limited to issues, pull requests, documentation, examples, etc.

