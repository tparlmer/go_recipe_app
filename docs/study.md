# Go Concepts to Review

## Language Features
- [ ] Struct tags (e.g. `json:"id"`)
  - [Official Go Struct Tags](https://golang.org/ref/spec#Struct_types)
  - [Practical Go: Struct Tags](https://www.digitalocean.com/community/tutorials/how-to-use-struct-tags-in-go)
- [ ] Array vs Slice syntax
  - [Go Blog: Slices](https://blog.golang.org/slices-intro)
- [ ] Interface implementation
  - [Go Interface Types](https://golang.org/doc/effective_go#interfaces)
- [ ] Error handling patterns
  - [Error Handling in Go](https://go.dev/blog/error-handling-and-go)
- [ ] Pointer receivers vs value receivers
  - [Go Methods](https://golang.org/doc/effective_go#methods)

## Web Development
- [ ] Form handling in Go
  - [Go Web Examples: Forms](https://gowebexamples.com/forms/)
- [ ] Template syntax and patterns
  - [Go Templates Documentation](https://pkg.go.dev/text/template)
- [ ] HTTP request/response lifecycle
  - [Go HTTP Handler](https://golang.org/doc/articles/wiki/#tmp_3)
- [ ] Middleware patterns
  - [Go Web Examples: Middleware](https://gowebexamples.com/basic-middleware/)
- [ ] Javascript DOM manipulation in templates
  - [MDN: DOM Manipulation](https://developer.mozilla.org/en-US/docs/Learn/JavaScript/Client-side_web_APIs/Manipulating_documents)
- [ ] Template security considerations
  - [Go Template Security](https://pkg.go.dev/html/template#hdr-Security)
- [ ] Form data parsing in Go
  - [Go Form Parsing](https://pkg.go.dev/net/http#Request.ParseForm)

## Project Patterns
- [ ] Project structure best practices
  - [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [ ] Dependency injection
  - [Go Dependency Injection](https://www.alexedwards.net/blog/organising-database-access)
- [ ] Error handling strategies
  - [Go Error Handling Best Practices](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)
- [ ] Testing patterns
  - [Go Testing Package](https://pkg.go.dev/testing)

## Debugging Skills
- [ ] Go logging best practices
  - [Go Logger Package](https://pkg.go.dev/log)
  - [Effective Logging in Go](https://www.digitalocean.com/community/tutorials/how-to-use-the-logger-package-in-go)
- [ ] Browser Developer Tools
  - [Chrome DevTools Network](https://developer.chrome.com/docs/devtools/network/)
- [ ] HTTP request/response debugging
  - [Go HTTP Debugging Guide](https://pkg.go.dev/net/http/httputil#DumpRequest)
- [ ] Template debugging
  - [Go Template Debugging](https://pkg.go.dev/text/template#hdr-Functions)

## Form Handling Lessons Learned
- [ ] PUT Request Form Data Handling
  - Key Issues:
    1. Browser handling of PUT requests differs from POST
    2. Form data encoding must be explicitly set
    3. Content-Type headers are crucial
  - Solution Components:
    - Form attributes: `enctype="application/x-www-form-urlencoded"`
    - Fetch headers: `'Content-Type': 'application/x-www-form-urlencoded'`
    - Data conversion: `new URLSearchParams(formData).toString()`
  - References:
    - [MDN: Using FormData](https://developer.mozilla.org/en-US/docs/Web/API/FormData/Using_FormData_Objects)
    - [Go Form Parsing](https://pkg.go.dev/net/http#Request.ParseForm)

## Current Questions
- How do browsers handle different HTTP methods with form data?
- What are the best practices for handling PUT/DELETE requests in Go web apps?
- When should we use different form encoding types?

## Resources
- [Go Template Documentation](https://pkg.go.dev/html/template)
- [MDN JavaScript in HTML](https://developer.mozilla.org/en-US/docs/Learn/HTML/Howto/Use_JavaScript_within_a_webpage)
- [MDN DOM Manipulation](https://developer.mozilla.org/en-US/docs/Learn/JavaScript/Client-side_web_APIs/Manipulating_documents)
- [Go Form Processing](https://pkg.go.dev/net/http#Request.ParseForm)

## Database Concepts
- [ ] BoltDB Fundamentals
  - Key Concepts:
    - Key-value storage
    - B+ tree implementation
    - ACID transactions
    - Single-writer, multiple-reader design
  - References:
    - [BoltDB GitHub](https://github.com/boltdb/bolt)
    - [BoltDB Docs](https://pkg.go.dev/go.etcd.io/bbolt)

- [ ] Transaction Management
  - [ ] Read transactions
  - [ ] Write transactions
  - [ ] Transaction patterns in Go
  - References:
    - [BoltDB Transactions](https://pkg.go.dev/go.etcd.io/bbolt#DB.Begin)
    - [Go Database Patterns](https://www.alexedwards.net/blog/organising-database-access)

- [ ] Data Serialization
  - [ ] JSON encoding/decoding
  - [ ] Binary encoding options
  - [ ] Error handling in serialization
  - References:
    - [Go JSON Package](https://pkg.go.dev/encoding/json)
    - [Go Binary Encoding](https://pkg.go.dev/encoding/gob)

## Deployment Concepts
- [ ] Server Configuration
  - Key Concepts:
    - Process management (systemd vs supervisor)
    - Environment variables vs config files
    - File permissions and security
    - Path management in production
  - References:
    - [Systemd Service Files](https://www.digitalocean.com/community/tutorials/how-to-use-systemctl-to-manage-systemd-services-and-units)
    - [Linux File Permissions](https://www.digitalocean.com/community/tutorials/linux-permissions-basics-and-how-to-use-umask-on-a-vps)

- [ ] Database Management
  - Key Concepts:
    - BoltDB file location strategies
    - Backup procedures
    - Permission settings
    - Recovery plans
  - References:
    - [Linux Directory Structure](https://www.digitalocean.com/community/tutorials/linux-directory-structure-explained)
    - [Backup Strategies](https://www.digitalocean.com/community/tutorials/how-to-choose-an-effective-backup-strategy-for-your-vps)

- [ ] Web Server Configuration
  - Key Concepts:
    - Nginx as reverse proxy
    - SSL/TLS setup
    - Static file serving
    - Security headers
  - References:
    - [Nginx Reverse Proxy](https://www.digitalocean.com/community/tutorials/how-to-configure-nginx-as-a-reverse-proxy-on-ubuntu-22-04)
    - [Let's Encrypt SSL](https://www.digitalocean.com/community/tutorials/how-to-secure-nginx-with-let-s-encrypt-on-ubuntu-22-04)

- [ ] Production Logging
  - Key Concepts:
    - Log rotation
    - Error reporting
    - Performance monitoring
    - Access logs
  - References:
    - [Go Production Logging](https://www.digitalocean.com/community/tutorials/how-to-use-logrotate-to-manage-log-files-on-ubuntu-20-04)
    - [Nginx Access Logs](https://www.digitalocean.com/community/tutorials/how-to-configure-logging-and-log-rotation-in-nginx-on-an-ubuntu-vps)

- [ ] Deployment Process
  - Key Concepts:
    - Build process
    - File transfer methods
    - Zero-downtime updates
    - Rollback procedures
  - References:
    - [Go Deployment Best Practices](https://www.digitalocean.com/community/tutorials/how-to-deploy-a-go-web-application-using-nginx-on-ubuntu-18-04)
    - [SCP File Transfer](https://www.digitalocean.com/community/tutorials/how-to-use-scp-to-securely-transfer-files-between-servers)

## Linux System Concepts
- [ ] Directory Structure
  - Key Concepts:
    - Root filesystem (`/`): The top level of the entire filesystem
    - Key directories:
      - `/bin/`: System binaries
      - `/etc/`: System configuration
      - `/home/`: Regular user home directories
      - `/root/`: Root user's home directory
      - `/var/`: Variable data (logs, websites, databases)
      - `/usr/`: User programs and libraries
      - `/tmp/`: Temporary files
    - Web applications:
      - `/var/www/`: Standard location for website files
      - `/var/log/`: System and application logs
      - `/var/lib/`: Database files
  - References:
    - [Linux Directory Structure Explained](https://www.digitalocean.com/community/tutorials/linux-directory-structure-explained)
    - [Linux Filesystem Hierarchy Standard](https://refspecs.linuxfoundation.org/FHS_3.0/fhs/index.html)

- [ ] User and Permission Management
  - Key Concepts:
    - System users vs regular users
    - Service users (non-login users)
    - File ownership and permissions
    - Security best practices
  - Implementation:
    - Creating service users (`useradd -r`)
    - Setting ownership (`chown`)
    - Permission modes (`chmod`)
    - Directory structure permissions
  - References:
    - [Linux Users and Groups](https://www.digitalocean.com/community/tutorials/initial-server-setup-with-ubuntu-22-04)
    - [Service User Security](https://www.digitalocean.com/community/tutorials/how-to-create-a-sudo-user-on-ubuntu-quickstart)

- [ ] User ID Generation
  - Key Concepts:
    - Use UUIDs for unique user identification
    - UUIDs are universally unique and prevent collisions
  - Implementation:
    - Use `github.com/google/uuid` package
    - Generate UUIDs with `uuid.New().String()`
  - References:
    - [UUID Package Documentation](https://pkg.go.dev/github.com/google/uuid)

- [ ] Password Security
  - Key Concepts:
    - Hash passwords before storing them
    - Use a strong hashing algorithm like bcrypt
    - Validate passwords by comparing hashes
  - Implementation:
    - Use `golang.org/x/crypto/bcrypt` package
    - Hash passwords with `bcrypt.GenerateFromPassword`
    - Validate with `bcrypt.CompareHashAndPassword`
  - References:
    - [Bcrypt Package Documentation](https://pkg.go.dev/golang.org/x/crypto/bcrypt)

## Database Management
- [ ] Backup Strategies
  - Key Concepts:
    - Daily automated backups using cron
    - Backup rotation and retention policies
    - Timestamped backup files
    - Backup location (`/var/backups/`)
    - Permission management for backup files
  - Implementation:
    - Backup script in `/usr/local/bin/`
    - Cron job configuration
    - Maximum backup retention setting
    - Automated cleanup of old backups
  - References:
    - [Cron Job Tutorial](https://www.digitalocean.com/community/tutorials/how-to-use-cron-to-automate-tasks-ubuntu-1804)
    - [Linux Backup Best Practices](https://www.digitalocean.com/community/tutorials/how-to-backup-your-linux-system)

- [ ] Monitoring and Health Checks
  - Key Concepts:
    - Health check endpoints
    - Uptime monitoring
    - Service status tracking
    - Email alerts
  - Implementation:
    - UptimeRobot configuration
    - Health check HTML file
    - Nginx health endpoint
    - Systemd service monitoring
  - Tools:
    - UptimeRobot
    - Systemd service files
    - Nginx configuration
    - Email notifications
  - References:
    - [Server Monitoring Basics](https://www.digitalocean.com/community/tutorials/an-introduction-to-service-monitoring-with-systemd)
    - [Nginx Health Checks](https://www.nginx.com/blog/smart-health-checks-nginx-plus/)

## Security Fundamentals
- [ ] Input Validation & Sanitization
  - Key Concepts:
    - XSS prevention
    - SQL injection prevention (even with BoltDB)
    - Form data validation
    - Content security policies
  - References:
    - [OWASP Input Validation Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Input_Validation_Cheat_Sheet.html)
    - [Go Security Practices](https://golang.org/security/best-practices)

- [ ] HTTP Security Headers
  - Key Concepts:
    - Content-Security-Policy
    - X-Frame-Options
    - X-XSS-Protection
    - HSTS
  - References:
    - [OWASP Secure Headers](https://owasp.org/www-project-secure-headers/)
    - [Security Headers Explained](https://securityheaders.com/)

## Error Handling & Monitoring
- [ ] Structured Logging
  - Key Concepts:
    - Log levels and when to use them:
      - DEBUG: Detailed information for debugging
      - INFO: General operational events
      - WARN: Warning messages for potentially harmful situations
      - ERROR: Error events that might still allow the app to continue
      - FATAL: Very severe errors that will terminate the app
    - Structured log format:
      ```go
      // Instead of:
      log.Printf("Error getting recipe: %v", err)
      
      // Use structured logging:
      logger.Error("failed to get recipe",
          "recipeID", id,
          "error", err,
          "path", r.URL.Path,
          "method", r.Method,
      )
      ```
    - Context fields to include:
      - Request path
      - HTTP method
      - User ID (when auth is added)
      - Timing information
      - Error details
      - Resource IDs
  - References:
    - [Go slog Package](https://pkg.go.dev/log/slog)
    - [Go Structured Logging Best Practices](https://www.digitalocean.com/community/tutorials/how-to-use-structured-logging-in-go)

- [ ] Error Reporting
  - Key Concepts:
    - Error classification
    - Stack trace capture
    - Error aggregation
    - Alert thresholds
  - Tools:
    - Sentry
    - ELK Stack
    - Prometheus/Grafana
  - References:
    - [Go Error Handling Patterns](https://go.dev/blog/error-handling-and-go)
    - [Monitoring Go Applications](https://prometheus.io/docs/guides/go-application/)

## Environment Configuration
- [ ] Configuration Management
  - Key Concepts:
    - Environment variables
    - Configuration files
    - Secrets management
    - Configuration validation
  - Implementation:
    - Using `os.Getenv()`
    - `.env` files
    - Configuration structs
    - Validation checks
  - References:
    - [12-Factor App Config](https://12factor.net/config)
    - [Go Configuration Patterns](https://golang.org/doc/effective_go#constants) 