# Redis Integration Testing Guide

This guide provides specific test scenarios for validating Redis integration.

## Quick Test Commands

### Test Redis Connection

```bash
# Check if Redis is running
docker exec goapi_redis redis-cli ping
# Should return: PONG

# Or if Redis is local
redis-cli ping
```

### Test User Caching

```bash
# 1. Get user (first time - cache miss)
curl -H "Authorization: Bearer TOKEN" http://localhost:8080/users/1

# 2. Check Redis for cached user
docker exec goapi_redis redis-cli GET "user:id:1"

# 3. Get user again (should be faster - cache hit)
curl -H "Authorization: Bearer TOKEN" http://localhost:8080/users/1
```

### Test Cache Invalidation

```bash
# 1. Cache a user
curl -H "Authorization: Bearer TOKEN" http://localhost:8080/users/1

# 2. Update the user
curl -X PUT -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"first_name":"Updated"}' \
  http://localhost:8080/users/1

# 3. Check Redis - should have updated data
docker exec goapi_redis redis-cli GET "user:id:1"
```

### Test Rate Limiting with Redis

```bash
# Make 70 requests quickly
for i in {1..70}; do
  curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/
done

# Check Redis rate limit counter
docker exec goapi_redis redis-cli GET "ratelimit:127.0.0.1"
```

### Test Distributed Rate Limiting

```bash
# Terminal 1: Start API on port 8080
PORT=8080 go run main.go

# Terminal 2: Start API on port 8081
PORT=8081 go run main.go

# Terminal 3: Make requests to both instances
for i in {1..30}; do
  curl http://localhost:8080/ &
  curl http://localhost:8081/ &
done

# Check shared rate limit
docker exec goapi_redis redis-cli GET "ratelimit:127.0.0.1"
```

## Manual Testing Scenarios

### Scenario 1: Test with Redis Disabled

1. Set `REDIS_ENABLED=false` in `.env`
2. Start application: `make run`
3. Verify logs show: "Redis is disabled"
4. Make API requests - should work normally
5. Verify no Redis connection attempts

**Expected:** Application works without Redis, uses no-op cache

### Scenario 2: Test with Redis Enabled

1. Start Redis: `docker run -d -p 6379:6379 redis:7-alpine`
2. Set `REDIS_ENABLED=true` in `.env`
3. Start application: `make run`
4. Verify logs show: "Redis connection established"
5. Make API requests
6. Check Redis for cached data

**Expected:** Application uses Redis, caches data, faster responses

### Scenario 3: Test Cache Performance

1. Enable Redis and start application
2. Time first request: `time curl http://localhost:8080/users/1`
3. Time second request: `time curl http://localhost:8080/users/1`
4. Compare response times

**Expected:** Second request should be significantly faster (cache hit)

### Scenario 4: Test Graceful Degradation

1. Start application with Redis enabled
2. Stop Redis: `docker stop <redis-container>`
3. Make API requests
4. Verify application continues to work
5. Check logs for cache errors (should not crash)

**Expected:** Application continues working, falls back to no-op cache

## Validation Checklist

Use this checklist to verify Redis integration:

### Phase 1: Basic Setup
- [ ] Redis service added to docker-compose.yml
- [ ] Redis environment variables configured
- [ ] Application connects to Redis when enabled
- [ ] Application works without Redis (no-op cache)

### Phase 2: User Caching
- [ ] User data cached on first read
- [ ] Cache hit on subsequent reads (faster)
- [ ] Both ID and email keys cached
- [ ] Cache TTL is 15 minutes
- [ ] Cache invalidated on update
- [ ] Cache invalidated on delete
- [ ] Old email key deleted on email change

### Phase 3: Rate Limiting
- [ ] Rate limiting uses Redis when enabled
- [ ] Rate limiting uses in-memory when Redis disabled
- [ ] Distributed rate limiting works across instances
- [ ] Rate limit falls back to in-memory if Redis fails

### Phase 4: Error Handling
- [ ] Application starts if Redis unavailable
- [ ] Cache errors don't crash application
- [ ] Rate limiting works if Redis fails
- [ ] Graceful degradation to no-op cache

### Phase 5: Performance
- [ ] Cached requests are faster
- [ ] Database load reduced with caching
- [ ] No performance degradation without Redis

## Troubleshooting

### Redis Connection Issues

```bash
# Check if Redis is running
docker ps | grep redis

# Check Redis logs
docker logs goapi_redis

# Test Redis connection manually
docker exec goapi_redis redis-cli ping
```

### Cache Not Working

```bash
# Check if Redis is enabled
grep REDIS_ENABLED .env

# Check application logs for Redis connection
# Should see "Redis connection established"

# Check Redis keys
docker exec goapi_redis redis-cli KEYS "*"
```

### Rate Limiting Issues

```bash
# Check rate limit keys in Redis
docker exec goapi_redis redis-cli KEYS "ratelimit:*"

# Check rate limit count
docker exec goapi_redis redis-cli GET "ratelimit:127.0.0.1"

# Reset rate limit (for testing)
docker exec goapi_redis redis-cli DEL "ratelimit:127.0.0.1"
```

### Performance Issues

```bash
# Monitor Redis performance
docker exec goapi_redis redis-cli INFO stats

# Check cache hit rate
docker exec goapi_redis redis-cli INFO stats | grep keyspace_hits
```

## Expected Test Results

When all Redis tests pass:

```
✓ Application starts with Redis disabled
✓ Application starts with Redis enabled
✓ User caching works (cache hits/misses)
✓ Cache invalidation on updates
✓ Cache invalidation on deletes
✓ Cache invalidation on email changes
✓ Distributed rate limiting works
✓ Graceful degradation when Redis unavailable
✓ Performance improvement with caching
```

