# Go URL Shortener

A simple URL shortener built with **Go + Redis + Base62 + Feistel encryption**.

It converts a long URL into a short URL and redirects the short URL back to the original URL.

## How it works

```text
Long URL
   │
   ▼
Generate ID
   │
   ▼
Feistel Obfuscation
   │
   ▼
Base62 Encoding
   │
   ▼
Short Code
   │
   ▼
http://localhost:8080/r/x7Kp91
```

When someone opens the short URL:

```text
Short Code
   │
   ▼
Base62 Decode
   │
   ▼
Feistel Reverse
   │
   ▼
Original ID
   │
   ▼
Redis
   │
   ▼
Long URL
   │
   ▼
Redirect
```

## Why Base62?

Base62 uses:

```text
a-z + A-Z + 0-9
```

So numbers can become shorter codes.

```text
1       → b
62      → ba
1000    → qi
```

## Why Feistel?

Redis IDs are sequential:

```text
1 → 2 → 3 → 4 → 5
```

If we directly convert them to Base62, people can easily guess the next URL.

Feistel uses a secret key to **obfuscate the ID**:

```text
ID: 1
 ↓
Feistel + Secret
 ↓
Random-looking number
 ↓
Base62
 ↓
Short Code
```

The same secret can reverse the process.


## Run

Install dependencies:

```bash
go mod tidy
```
Connect and start redis and also add shortner key

Start Redis, then:

```bash
go run .
```

Server:

```text
http://localhost:8080
```

## Example

Shorten:

```text
GET /shorten?url=https://google.com
```

The server generates something like:

```text
http://localhost:8080/r/x7Kp91
```

Opening:

```text
/r/x7Kp91
```

redirects to:

```text
https://google.com
```

## Main Technologies

- **Go** — Backend server
- **Redis** — Stores URL mappings and generates IDs
- **Base62** — Converts numbers into short codes
- **Feistel + HMAC** — Makes sequential IDs harder to guess
- **godotenv** — Loads `.env` configuration