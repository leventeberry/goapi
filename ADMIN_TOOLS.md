# Admin Tools Guide

This project includes web-based admin tools for managing Redis and PostgreSQL databases.

## Redis Commander

Web-based Redis management interface for viewing and managing Redis cache data.

### Access

- **URL**: http://localhost:8081
- **Username**: `admin`
- **Password**: `admin`

### Features

- Browse Redis keys and values
- View cached user data (`user:id:*`, `user:email:*`)
- View rate limiting data (`ratelimit:*`)
- Edit/delete keys
- Monitor Redis operations
- Execute Redis commands

### Usage

1. Start the services:
   ```bash
   docker-compose up -d
   ```

2. Access Redis Commander:
   ```bash
   # Using make command (opens in browser)
   make docker-open-redis-commander
   
   # Or manually navigate to
   http://localhost:8081
   ```

3. Login with:
   - Username: `admin`
   - Password: `admin`

### Common Tasks

**View cached users:**
- Filter keys by pattern: `user:*`
- Click on a key to view its JSON value
- Keys follow patterns:
  - `user:id:{id}` - User cached by ID
  - `user:email:{email}` - User cached by email
  - `ratelimit:{ip}` - Rate limiting counters

**Monitor cache activity:**
- Watch keys being created/updated in real-time
- View TTL (Time To Live) for each key
- See when cache entries expire

## pgAdmin

Web-based PostgreSQL administration and development platform.

### Access

- **URL**: http://localhost:5050
- **Email**: `admin@goapi.com`
- **Password**: `admin`

### Features

- Database browser and query tool
- SQL editor with syntax highlighting
- Table data viewer and editor
- Query history
- Database statistics and monitoring
- Export/import data

### Usage

1. Start the services:
   ```bash
   docker-compose up -d
   ```

2. Access pgAdmin:
   ```bash
   # Using make command (opens in browser)
   make docker-open-pgadmin
   
   # Or manually navigate to
   http://localhost:5050
   ```

3. Login with:
   - Email: `admin@goapi.com`
   - Password: `admin`

### Setting Up Database Connection

After logging in, you need to add a server connection:

1. Right-click "Servers" → "Register" → "Server"

2. **General Tab:**
   - Name: `GoAPI Database` (or any name)

3. **Connection Tab:**
   - Host name/address: `db` (Docker service name)
   - Port: `5432`
   - Maintenance database: `goapi`
   - Username: `goapi_user`
   - Password: `goapi_password`
   - Check "Save password"

4. Click "Save"

### Common Tasks

**View tables:**
- Navigate: Servers → GoAPI Database → Databases → goapi → Schemas → public → Tables
- Right-click on `users` table → "View/Edit Data" → "All Rows"

**Run SQL queries:**
- Right-click on database → "Query Tool"
- Write SQL queries:
  ```sql
  SELECT * FROM users;
  SELECT * FROM users WHERE email = 'test@example.com';
  ```

**View table structure:**
- Right-click on table → "Properties"
- See columns, constraints, indexes

## Makefile Commands

Quick access commands:

```bash
# Open Redis Commander in browser
make docker-open-redis-commander

# Open pgAdmin in browser
make docker-open-pgadmin

# View Redis Commander logs
make docker-logs-redis-commander

# View pgAdmin logs
make docker-logs-pgadmin
```

## Security Note

⚠️ **Important**: The default credentials (`admin`/`admin`) are for **development only**. 

For production:
1. Change default passwords
2. Use environment variables for credentials
3. Consider removing admin tools or restricting access
4. Use secure passwords and enable authentication

## Troubleshooting

### Redis Commander won't connect

- Ensure Redis container is running: `docker ps | grep redis`
- Check Redis is healthy: `docker-compose ps redis`
- Verify Redis Commander logs: `make docker-logs-redis-commander`

### pgAdmin can't connect to database

- Ensure database container is running: `docker ps | grep db`
- Check database is healthy: `docker-compose ps db`
- Verify connection details:
  - Host: `db` (not `localhost` when connecting from pgAdmin container)
  - Port: `5432`
  - Username: `goapi_user`
  - Password: `goapi_password`

### Port conflicts

If ports 5050 or 8081 are already in use, you can change them in `docker-compose.yml`:

```yaml
redis-commander:
  ports:
    - "8082:8081"  # Change 8082 to your preferred port

pgadmin:
  ports:
    - "5051:80"    # Change 5051 to your preferred port
```

