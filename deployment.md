# Deployment Guide for Go Web Applications

## Pre-Deployment Checklist
- [ ] Local testing complete
- [ ] Database backups configured
- [ ] Environment variables documented
- [ ] Required ports identified
- [ ] Domain DNS configured

## Server Setup
1. Basic Server Configuration
   - [ ] Update system packages
   - [ ] Configure SSH access
   - [ ] Set up UFW firewall
     ```bash
     sudo ufw allow 22/tcp    # SSH
     sudo ufw allow 80/tcp    # HTTP
     sudo ufw allow 443/tcp   # HTTPS
     sudo ufw enable
     ```
   - [ ] Verify firewall status
     ```bash
     sudo ufw status numbered
     ```

2. Web Server Configuration
   - [ ] Install nginx
   - [ ] Configure nginx as reverse proxy
   - [ ] Set up systemd service
   - [ ] Configure logging

3. SSL/TLS Security
   - [ ] Install Certbot
   - [ ] Obtain SSL certificate
   - [ ] Configure auto-renewal
   - [ ] Test SSL configuration

SSL (Secure Sockets Layer) certificates are crucial for:
- Encrypting traffic between users and server
- Providing HTTPS for your domain
- Showing security indicators (padlock) in browsers
- Meeting modern web security standards

### Certificate Options
1. Let's Encrypt (Free, Automated)
   - Most common choice for small-medium sites
   - Auto-renews every 90 days
   - Widely trusted by browsers
   - Easy to set up with Certbot

2. Commercial Certificates
   - Paid options from various providers
   - Extended validation available
   - Longer renewal periods
   - Additional features/warranty

4. Database Configuration
   - [ ] Set up database directory
   - [ ] Configure permissions
   - [ ] Set up backup strategy
   - [ ] Verify backup automation

5. Monitoring and Logging
   - [ ] Configure application logs
   - [ ] Set up log rotation
   - [ ] Implement health checks
   - [ ] Configure monitoring alerts

## Deployment Process
1. Application Preparation
2. File Transfer
3. Service Configuration
4. Testing
5. Monitoring

## Rollback Procedures
1. Backup Verification
2. Rollback Steps
3. Service Restoration
4. Verification Steps

## SSL/TLS Security

### Setting up Let's Encrypt Certificate
1. Install Certbot and nginx plugin:
```bash
sudo apt install certbot python3-certbot-nginx
```

2. Obtain and install initial certificate:
```bash
sudo certbot --nginx -d yourdomain.com
```

3. Configure nginx for www and non-www domains:
```nginx
server {
    listen 80;
    listen [::]:80;
    server_name yourdomain.com www.yourdomain.com;
    return 301 https://yourdomain.com$request_uri;
}

server {
    listen 443 ssl;
    listen [::]:443 ssl;
    server_name www.yourdomain.com;
    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
    return 301 https://yourdomain.com$request_uri;
}

server {
    listen 443 ssl;
    listen [::]:443 ssl;
    server_name yourdomain.com;
    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

4. Expand certificate to cover www subdomain:
```bash
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com
```

### Verification Steps
1. Check certificate status:
```bash
sudo certbot certificates
```

2. Verify auto-renewal:
```bash
sudo certbot renew --dry-run
```

3. Test nginx configuration:
```bash
sudo nginx -t
sudo systemctl reload nginx
```

4. Verify in browser:
- Check for padlock icon
- Verify HTTPS in URL
- Confirm www to non-www redirect
- Test HTTP to HTTPS redirect

### Certificate Maintenance
- Certificates auto-renew every 90 days
- Renewal cron job installed automatically
- Logs located in `/var/log/letsencrypt/`
- Certificates stored in `/etc/letsencrypt/live/`

### Troubleshooting
- Check nginx logs: `sudo nginx -t`
- Check certbot logs: `sudo certbot certificates`
- Renewal issues: `sudo systemctl status certbot.timer`

## Environment Configuration

### Configuration Files
1. Production Environment File (`/etc/recipe-app/env`):
```bash
RECIPE_APP_PORT=8080
RECIPE_APP_ENV=production
RECIPE_APP_DB_PATH=/var/lib/recipe-app/recipes.db
RECIPE_APP_LOG_DIR=/var/log/recipe-app
RECIPE_APP_LOG_LEVEL=info
RECIPE_APP_BASE_URL=https://memeticuniverse.com
RECIPE_APP_LOG_FORMAT=json
RECIPE_APP_READ_TIMEOUT=15s
RECIPE_APP_WRITE_TIMEOUT=15s
```

2. Development Environment File (`.env.example`):
```bash
RECIPE_APP_PORT=8080
RECIPE_APP_ENV=development
RECIPE_APP_DB_PATH=data/recipes.db
RECIPE_APP_LOG_DIR=logs
RECIPE_APP_LOG_LEVEL=debug
RECIPE_APP_BASE_URL=http://localhost:8080
RECIPE_APP_LOG_FORMAT=text
```

### Directory Setup
```bash
# Create service user
sudo useradd -r -s /bin/false recipe-app

# Create required directories
sudo mkdir -p /var/log/recipe-app
sudo mkdir -p /var/lib/recipe-app
sudo mkdir -p /etc/recipe-app

# Set ownership
sudo chown recipe-app:recipe-app /var/log/recipe-app
sudo chown recipe-app:recipe-app /var/lib/recipe-app
sudo chown recipe-app:recipe-app /etc/recipe-app

# Set permissions
sudo chmod 755 /var/log/recipe-app
sudo chmod 755 /var/lib/recipe-app
sudo chmod 755 /etc/recipe-app
```

### Configuration Management
1. **Local Development**
   - Copy `.env.example` to `.env`
   - Modify values for local environment
   - Never commit `.env` to version control

2. **Production Deployment**
   - Create/update `/etc/recipe-app/env`
   - Secure file permissions
   - Reload systemd after changes:
     ```bash
     sudo systemctl daemon-reload
     sudo systemctl restart recipe-app
     ```

3. **Validation**
   - Application validates config on startup
   - Check logs for configuration errors
   - Use health check endpoint to verify service

### Environment Variables Reference
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| RECIPE_APP_PORT | HTTP server port | 8080 | No |
| RECIPE_APP_ENV | Environment name | development | No |
| RECIPE_APP_DB_PATH | Database file location | data/recipes.db | No |
| RECIPE_APP_LOG_DIR | Log directory | logs | No |
| RECIPE_APP_LOG_LEVEL | Log level (debug/info/warn/error) | info | No |
| RECIPE_APP_BASE_URL | Base URL for the application | http://localhost:8080 | No |
| RECIPE_APP_LOG_FORMAT | Log format (json/text) | text | No |
| RECIPE_APP_READ_TIMEOUT | HTTP read timeout | 15s | No |
| RECIPE_APP_WRITE_TIMEOUT | HTTP write timeout | 15s | No |

## Service User Setup
1. Create Service User
   ```bash
   # Create non-login service user
   sudo useradd -r -s /bin/false recipe-app
   ```

2. Directory Permissions
   ```bash
   # Application directories
   sudo mkdir -p /var/www/recipe-app/current
   sudo chown -R recipe-app:recipe-app /var/www/recipe-app/current
   sudo chmod 755 /var/www/recipe-app/current

   # Log directory
   sudo mkdir -p /var/log/recipe-app
   sudo chown -R recipe-app:recipe-app /var/log/recipe-app
   sudo chmod 755 /var/log/recipe-app

   # Database directory
   sudo mkdir -p /var/lib/recipe-app
   sudo chown -R recipe-app:recipe-app /var/lib/recipe-app
   sudo chmod 755 /var/lib/recipe-app

   # Environment configuration
   sudo mkdir -p /etc/recipe-app
   sudo touch /etc/recipe-app/env
   sudo chown root:recipe-app /etc/recipe-app/env
   sudo chmod 640 /etc/recipe-app/env
   ```

3. Security Verification
   ```bash
   # Verify user creation
   id recipe-app

   # Verify directory permissions
   ls -l /var/www/recipe-app/current
   ls -l /var/log/recipe-app
   ls -l /var/lib/recipe-app
   ls -l /etc/recipe-app/env
   ```
