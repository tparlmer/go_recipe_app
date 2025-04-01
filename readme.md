# Go Recipe App

This is a simple recipe app built with Go. It allows you to add, get, update, and delete recipes.

### How to run the application:

1. Clone the repository
2. Run `go run cmd/main.go`

We will use BoltDB as our database.

### Overview of current file structure:

recipe-app/cmd/main.go (main application entry point)
recipe-app/internal/ (for application logic)
recipe-app/web/templates/ (for HTML templates)


### Viewing Server Logs

#### View application logs
tail -f /var/log/recipe-app/app.log

#### View error logs
tail -f /var/log/recipe-app/error.log

#### View system service logs
journalctl -u recipe-app -f

#### View nginx access logs
tail -f /var/log/nginx/access.log

#### View nginx error logs
tail -f /var/log/nginx/error.log