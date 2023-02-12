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

# Install and config

TeamsACS uses PostgreSQL as its primary database and uses the Timescaledb extension

- Create Database

    CREATE USER teamsacs WITH PASSWORD 'teamsacs'
    
    CREATE DATABASE teamsacs OWNER teamsacs;
    
    GRANT ALL PRIVILEGES ON DATABASE teamsacs TO teamsacs;

# Docker Deploy

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
      - "2979:2979"
      - "2989:2989"
      - "2999:2999"
    volumes:
      - /data/teamsacs:/var/teamsacs
    environment:
      - GODEBUG=x509ignoreCN=0
      - TEAMSACS_DB_HOST=pgdb
      - TEAMSACS_DB_PORT=5432
      - TEAMSACS_DB_USER=teamsacs
      - TEAMSACS_DB_PWD=teamsacs
      - TEAMSACS_WEB_DEBUG=1
    networks:
      teamsacs_network:

networks:
  teamsacs_network:

```

## Links

- [TeamsACS Documentation Wiki](https://github.com/ca17/teamsacs/wiki)
- [Mikrotik Tr069 best practices](https://wiki.mikrotik.com/wiki/Tr069-best-practices)
- [Mikrotik Tr069 client](https://wiki.mikrotik.com/wiki/Manual:TR069-client)
- [Mikrotik TR069 client supported parameter reference](https://wiki.mikrotik.com/tr069ref/current.html)


# Contribute

We welcome contributions of any kind, including but not limited to issues, pull requests, documentation, examples, etc.

