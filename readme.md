# Go Watch That Site

Time to experiment with Golang! I'm waiting on deals to buy a new mountain bike and I want a simple program to scrape a site then compare it to the last scrape. If there's a change, I want to know about it.

## Flags

All flags can be set via environment variable or switch. If both are supplied, the switch will be used.

| Required | Environment         | Switch        | Purpose                               |
| -------- | ------------------- | ------------- | ------------------------------------- |
|    x     | GWTS_EMAIL_TO       | emailTo       | Set where the email will be sent to   |
|    x     | GWTS_EMAIL_FROM     | emailFrom     | Set where the email will be sent from |
|    x     | GWTS_EMAIL_PASSWORD | emailPassword | GMail password for `emailFrom` SMTP   |

### Examples

#### Go

```bash
./go-watch-that-site -emailTo=to@email.test -emailFrom=from@email.test -emailPassword=superComplexPassword
```

#### Docker

```bash
docker build -t go-watch-that-site .
docker run -e "GWTS_EMAIL_TO=to@email.test" -e "GWTS_EMAIL_FROM=from@email.test" -e "GWTS_EMAIL_PASSWORD=superComplexPassword" go-watch-that-site
```
