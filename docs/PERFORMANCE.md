# Performance Guide

## Performance Metrics

SwiftCal is designed to handle high-performance calendar and email operations with the following targets:

### Response Times

- **API Endpoints**: < 200ms average response time
- **Database Queries**: < 50ms average query time
- **Email Processing**: < 5 seconds per email
- **Event Creation**: < 100ms

### Throughput

- **Concurrent Users**: 1000+ simultaneous users
- **API Requests**: 10,000+ requests per minute
- **Email Processing**: 100+ emails per minute
- **Database Operations**: 50,000+ operations per minute

## Performance Optimization

### Database Optimization

#### Indexing Strategy

```sql
-- Users table
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Events table
CREATE INDEX idx_events_user_id ON events(user_id);
CREATE INDEX idx_events_start_time ON events(start_time);
CREATE INDEX idx_events_end_time ON events(end_time);
CREATE INDEX idx_events_user_start ON events(user_id, start_time);

-- Emails table
CREATE INDEX idx_emails_user_id ON emails(user_id);
CREATE INDEX idx_emails_processed_at ON emails(processed_at);
```

#### Query Optimization

- Use prepared statements
- Implement connection pooling
- Optimize complex queries with EXPLAIN
- Use appropriate data types

### Caching Strategy

#### Redis Caching

```go
// Cache frequently accessed data
type CacheService struct {
    redis *redis.Client
}

func (c *CacheService) GetUserEvents(userID int) ([]Event, error) {
    key := fmt.Sprintf("user_events:%d", userID)

    // Try cache first
    if cached, err := c.redis.Get(key).Result(); err == nil {
        var events []Event
        json.Unmarshal([]byte(cached), &events)
        return events, nil
    }

    // Fetch from database
    events := fetchEventsFromDB(userID)

    // Cache for 5 minutes
    c.redis.Set(key, events, 5*time.Minute)
    return events, nil
}
```

#### Application-Level Caching

- Cache user sessions
- Cache frequently accessed configuration
- Cache parsed email templates

### API Performance

#### Response Optimization

```go
// Use pagination for large datasets
func (h *EventHandler) GetEvents(c *gin.Context) {
    page := c.DefaultQuery("page", "1")
    limit := c.DefaultQuery("limit", "20")

    events, total := h.service.GetEventsPaginated(page, limit)

    c.JSON(200, gin.H{
        "events": events,
        "pagination": gin.H{
            "page": page,
            "limit": limit,
            "total": total,
        },
    })
}
```

#### Compression

- Enable gzip compression for API responses
- Compress static assets
- Use efficient JSON serialization

### Background Processing

#### Email Processing Queue

```go
// Use worker pools for email processing
type EmailProcessor struct {
    workers int
    queue   chan Email
}

func (ep *EmailProcessor) Start() {
    for i := 0; i < ep.workers; i++ {
        go ep.worker()
    }
}

func (ep *EmailProcessor) worker() {
    for email := range ep.queue {
        ep.processEmail(email)
    }
}
```

#### Scheduled Tasks

- Use cron jobs for maintenance tasks
- Implement retry mechanisms
- Monitor job execution times

## Monitoring and Profiling

### Application Metrics

```go
// Use Prometheus metrics
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )

    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
        },
        []string{"method", "endpoint"},
    )
)
```

### Database Monitoring

- Monitor slow queries
- Track connection pool usage
- Monitor index usage
- Set up alerts for performance degradation

### Profiling Tools

```bash
# CPU profiling
go tool pprof -http=:8080 cpu.prof

# Memory profiling
go tool pprof -http=:8080 mem.prof

# Goroutine profiling
go tool pprof -http=:8080 goroutine.prof
```

## Load Testing

### Test Scenarios

```bash
# Run load tests
./scripts/load_test.sh

# Test specific endpoints
ab -n 1000 -c 10 http://localhost:8080/api/events

# Test database performance
pgbench -c 10 -t 1000 swiftcal
```

### Performance Benchmarks

- API endpoint response times
- Database query performance
- Memory usage under load
- CPU utilization patterns

## Scaling Strategies

### Horizontal Scaling

- Use load balancers
- Implement stateless design
- Use shared databases
- Deploy multiple instances

### Vertical Scaling

- Increase server resources
- Optimize database configuration
- Use faster storage (SSD)
- Increase connection limits

### Microservices Architecture

- Split into smaller services
- Use message queues
- Implement service discovery
- Use API gateways

## Performance Checklist

### Development

- [ ] Use efficient algorithms
- [ ] Implement proper indexing
- [ ] Add caching where appropriate
- [ ] Optimize database queries
- [ ] Use connection pooling
- [ ] Implement pagination

### Deployment

- [ ] Enable compression
- [ ] Configure caching headers
- [ ] Set up monitoring
- [ ] Configure load balancing
- [ ] Optimize database settings
- [ ] Set up alerting

### Maintenance

- [ ] Monitor performance metrics
- [ ] Analyze slow queries
- [ ] Update dependencies
- [ ] Review and optimize indexes
- [ ] Clean up old data
- [ ] Update performance baselines
