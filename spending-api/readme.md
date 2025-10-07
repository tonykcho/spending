# Install go-lang migrate
brew install golang-migrate

# Create migrations
migrate create -ext sql -dir migrations {name}
