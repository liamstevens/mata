# Basic Usage Examples

## Echo Mode (Testing)

Test mata with echo mode:

```bash
# Start mata in echo mode on port 8080
./bin/mata -source :8080 -echo

# In another terminal, test with netcat
echo "hello world" | nc localhost 8080
```

## Single Target Forwarding

Forward traffic from port 8080 to a single target:

```bash
# Forward traffic to localhost:9000
./bin/mata -source :8080 -targets localhost:9000

# Test with netcat (assuming something is listening on 9000)
echo "test message" | nc localhost 8080
```

## Traffic Duplication

Duplicate traffic to multiple targets:

```bash
# Duplicate traffic to multiple services
./bin/mata -source :8080 -targets localhost:9001,localhost:9002,localhost:9003

# All three targets will receive the same traffic
echo "duplicated message" | nc localhost 8080
```

## Real-world Example

```bash
# Duplicate HTTP traffic from production to staging environments
./bin/mata -source :80 -targets prod.example.com:80,staging1.example.com:80,staging2.example.com:80
```

This will:
1. Accept connections on port 80
2. Forward all traffic to production server
3. Duplicate the same traffic to two staging environments
4. Return responses from the production server to clients