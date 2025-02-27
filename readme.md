# Telebackup

## What is Telebackup?
Telebackup is a simple backup tool for Telegram. It allows you to back up your local files to Telegram. It uses tar & gzip to compress the files and then sends them to Telegram chat

## Features
- Supports multiple files/directories
- File upload up to 2GB

## Installation
### Download binary
You can download prebuilt binaries from [releases](https://github.com/aiexz/telebackup/releases/latest) page

### Build from source
```bash
git clone https://github.com/aiexz/telebackup.git
cd telebackup
go build cmd/telebackup/main.go
```

### Run in Docker
```bash
docker run -v /path/to/config.yml:/config.yml -v /path/to/dir:/dr ghcr.io/aiexz/telebackup:master
```



## Usage
1. Create a Telegram bot using [BotFather](https://t.me/botfather) and get bot token
2. Get APP ID & API Hash from [my.telegram.org](https://my.telegram.org) or use provided in example
3. Edit `config.example.yml` and rename it to `config.yml`
4. Run `./telebackup` or `./telebackup --config /path/to/config.yml`

## Configuration
```yaml
appId: 6 # Telegram APP ID
appHash: eb06d4abfb49dc3eeb1aeb98ae0f581e # Telegram API Hash
botToken: 123:AAA # Telegram Bot Token
target: 56789123 # Telegram chat/channel username or chat ID
targets:
    - /tmp/test # List of files/directories to back up
    - /tmp/test2
```

> [!TIP]
> It is recommended to create a group or channel with the bot to not spam your personal messages


## Real-world example
Here is our minecraft server
```yaml
services:
  mc:
    image: itzg/minecraft-server
    ports:
      - "25565:25565"
    environment:
      EULA: "TRUE"
    volumes:
      - ./data:/data
  telebackup:
    container_name: telebackup_mc
    image: ghcr.io/aiexz/telebackup:master
    depends_on:
      - mc
    environment:
      APP_ID: 6
      APP_HASH: 123
      BOT_TOKEN: 123:abc
      TARGET: 123
      TARGETS: |
        /data
    volumes:
      - ./data:/data:ro
  ```
We also add a cron job to backup the server every day
```bash
0 0 * * * docker start telebackup_mc
```

## Roadmap
- [ ] Handle files larger than 2GB
- [x] Support for forums (chats with topics)
- [x] Support for usage without username, just chat ID
- [ ] Encryption/password protection
- [ ] Signing backups

## Contributing
All contributions are welcome. Feel free to open an issue or a pull request

## Awesome libraries used
- [gogram](https://github.com/AmarnathCJD/gogram) - Awesome Telegram API library

## Contact
- Telegram: [@aiexz](https://t.me/aiexz)

## License
[MIT](LICENSE)
