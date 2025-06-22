# Tracey

Simple, user-friendly CAPTCHA that keeps bots and spam out of your apps and websites.

## Configuration

1. Copy the example environment file:

   `cp .env.example .env`

2. Open `.env` and configure the variables.

## Run

`go run cmd/main.go`

## Todo

- [ ] Add rate limiting
- [ ] Add fingerprint collection & validation.
- [ ] Enhance security (obfuscation, encoding, encryption)
- [ ] Add frontend JavaScript SDK
- [ ] Add documentation

## License

This project is licensed under the [MIT License](LICENSE).