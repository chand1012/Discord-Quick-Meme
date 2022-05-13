set dotenv-load := true

default:
    just --list

build:
    go build -o GoDiscordBot . 

clean:
    go clean -cache
    rm GoDiscordBot

run:
    ./GoDiscordBot

build-docker:
    docker build -t discord-quick-meme .

docker-build:
 @just build-docker

docker-run:
    docker run --env-file .env --name quickmeme discord-quick-meme

run-docker:
 @just docker-run

test:
    go test -v ./...

fly:
    fly secrets set --app discord-quick-meme MONGO_CONNECT_STR=$MONGO_CONNECT_STR
    fly secrets set --app discord-quick-meme MONGO_DATABASE=$MONGO_DATABASE
    fly secrets set --app discord-quick-meme REDDIT_ID=$REDDIT_ID
    fly secrets set --app discord-quick-meme REDDIT_SECRET=$REDDIT_SECRET
    fly secrets set --app discord-quick-meme MODE=$MODE
    fly secrets set --app discord-quick-meme ADMINS=$ADMINS
    fly secrets set --app discord-quick-meme DISCORD_TOKEN=$DISCORD_TOKEN
    fly deploy --app discord-quick-meme

shoot-down:
    fly apps destroy discord-quick-meme