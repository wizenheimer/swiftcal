# Troubleshooting Guide

This guide helps you resolve common issues with SwiftCal.

## Common Issues

### Database Connection Issues

#### Error: "connection refused"

**Symptoms**: Application fails to start with database connection errors.

**Solutions**:

1. Check if PostgreSQL is running:

   ```bash
   sudo systemctl status postgresql
   # or
   brew services list | grep postgresql
   ```

2. Verify database configuration in `.env`:

   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_NAME=swiftcal
   DB_USER=postgres
   DB_PASSWORD=your_password
   ```

3. Test database connection:
   ```bash
   psql -h localhost -U postgres -d swiftcal
   ```

#### Error: "authentication failed"

**Solutions**:

1. Check user credentials
2. Verify pg_hba.conf configuration
3. Reset password if needed:
   ```sql
   ALTER USER postgres PASSWORD 'new_password';
   ```

### Email Service Issues

#### Error: "SMTP authentication failed"

**Symptoms**: Email sending fails with authentication errors.

**Solutions**:

1. Verify SMTP credentials in `.env`:

   ```env
   SMTP_HOST=smtp.gmail.com
   SMTP_PORT=587
   SMTP_USERNAME=your_email@gmail.com
   SMTP_PASSWORD=your_app_password
   ```

2. For Gmail, use App Password instead of regular password
3. Check if 2FA is enabled and generate app password

#### Error: "connection timeout"

**Solutions**:

1. Check firewall settings
2. Verify SMTP server is accessible
3. Try different SMTP port (587 or 465)
4. Check network connectivity

### OpenAI Integration Issues

#### Error: "invalid API key"

**Solutions**:

1. Verify OpenAI API key in `.env`:

   ```env
   OPENAI_API_KEY=sk-your-actual-api-key
   ```

2. Check API key format and validity
3. Verify account has sufficient credits
4. Check API rate limits

#### Error: "rate limit exceeded"

**Solutions**:

1. Implement exponential backoff
2. Reduce request frequency
3. Upgrade OpenAI plan if needed
4. Cache responses where possible

### Application Startup Issues

#### Error: "port already in use"

**Solutions**:

1. Check if another process is using the port:

   ```bash
   lsof -i :8080
   ```

2. Kill the process or change port:

   ```env
   SERVER_PORT=8081
   ```

3. Use different port in configuration

#### Error: "permission denied"

**Solutions**:

1. Check file permissions:

   ```bash
   chmod +x swiftcal
   ```

2. Run with appropriate user permissions
3. Check directory write permissions

## Performance Issues

### Slow API Responses

**Diagnosis**:

1. Check database query performance
2. Monitor CPU and memory usage
3. Review application logs

**Solutions**:

1. Add database indexes
2. Implement caching
3. Optimize queries
4. Scale application resources

### High Memory Usage

**Diagnosis**:

1. Monitor memory usage with `top` or `htop`
2. Check for memory leaks
3. Review goroutine count

**Solutions**:

1. Implement proper resource cleanup
2. Use connection pooling
3. Limit concurrent operations
4. Monitor and kill long-running goroutines

## Log Analysis

### Understanding Logs

SwiftCal uses structured logging with different levels:

```go
// Debug level - detailed information
logger.Debug("Processing email", "email_id", email.ID)

// Info level - general information
logger.Info("User registered", "user_id", user.ID)

// Warn level - potential issues
logger.Warn("Database connection slow", "duration", duration)

// Error level - errors that need attention
logger.Error("Failed to process email", "error", err)
```

### Common Log Patterns

#### Database Connection Issues

```
ERROR: failed to connect to database: connection refused
WARN: database connection slow: 2.5s
```

#### Email Processing Issues

```
ERROR: failed to send email: SMTP authentication failed
WARN: email processing delayed: rate limit exceeded
```

#### API Issues

```
ERROR: invalid request: missing required field 'email'
WARN: API rate limit approaching: 95% used
```

## Debug Mode

Enable debug mode for detailed logging:

```env
LOG_LEVEL=debug
ENVIRONMENT=development
```

### Debug Commands

```bash
# Check application status
curl http://localhost:8080/health

# Test database connection
go run cmd/server/main.go --test-db

# Validate configuration
go run cmd/server/main.go --validate-config
```

## Recovery Procedures

### Database Recovery

1. **Backup Restoration**:

   ```bash
   pg_restore -h localhost -U postgres -d swiftcal backup.dump
   ```

2. **Schema Reset**:
   ```bash
   psql -h localhost -U postgres -d swiftcal -f scripts/init.sql
   ```

### Application Recovery

1. **Graceful Restart**:

   ```bash
   # Send SIGTERM to application
   kill -TERM <pid>

   # Wait for graceful shutdown
   # Start application again
   make run
   ```

2. **Force Restart**:

   ```bash
   # Kill process if needed
   kill -9 <pid>

   # Start fresh
   make run
   ```

## Monitoring and Alerts

### Health Checks

Implement health check endpoints:

```go
func healthCheck(c *gin.Context) {
    // Check database
    if err := db.Ping(); err != nil {
        c.JSON(503, gin.H{"status": "unhealthy", "database": "down"})
        return
    }

    // Check external services
    if err := checkOpenAI(); err != nil {
        c.JSON(503, gin.H{"status": "degraded", "openai": "down"})
        return
    }

    c.JSON(200, gin.H{"status": "healthy"})
}
```

### Alerting

Set up alerts for:

- Application downtime
- High error rates
- Performance degradation
- Resource exhaustion

## Getting Help

### Before Asking for Help

1. Check this troubleshooting guide
2. Review application logs
3. Verify configuration
4. Test with minimal setup
5. Search existing issues

### Providing Information

When reporting issues, include:

- Error messages and logs
- Environment details
- Steps to reproduce
- Expected vs actual behavior
- Configuration files (without secrets)

### Support Channels

- GitHub Issues: For bug reports and feature requests
- Documentation: Check docs/ directory
- Community: Join discussions
- Email: For security issues
