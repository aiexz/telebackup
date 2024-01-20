# Telebackup

## What is Telebackup?
Telebackup is a simple backup tool for Telegram. It allows you to back up your local files to Telegram. It uses tar & gzip to compress the files and then sends them to Telegram chat

## Features
- Supports multiple files/directories
- File upload up to 2GB

## Installation
> [!NOTE]
> Releases will be available soon

### Build from source
```bash
git clone https://github.com/aiexz/telebackup.git
cd telebackup
go build cmd/telebackup/main.go
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
target: "@aiexz" # Telegram chat/channel username
targets:
    - /tmp/test # List of files/directories to backup
    - /tmp/test2
```

## Roadmap
- [ ] Make automated releases with GitHub Actions
- [ ] Handle files larger than 2GB
- [ ] Support for forums (chats with topics)
- [ ] Support for usage without username, just chat ID
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
