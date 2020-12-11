# go-rest-api
REST api in golang (again)

## Stack
```bash
Golang > 1.15
PostgreSQL
Redis
Ngingx
Docker/Docker-compose
```

## Plan

Approximated plan for the realization of the project

### 1. Authentification part
For the Authentification the tokenized system was selected, JWT system and Redis cache will be used. 

### For Authentification part the next packages are used: 
```bash
github.com/dgrijalva/jwt-go     // For tokenization
github.com/go-redis/redis/v8    // For Redis 
github.com/twinj/uuid           // For uniqueIDs
os                              // For environ
```

### 2. Modules and DB // Will be soon
### 3. Roles and Methods // Will be soon
```

## END