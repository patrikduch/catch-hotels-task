# 🧪 catch-hotels-task

## 🌍 CLI aplikace v Go pro monitoring dostupnosti webových stránek

Tato aplikace monitoruje dostupnost a odezvu webových stránek pomocí periodických HTTP GET požadavků.

### ✨ Funkce

- Každých 5 sekund odesílá požadavek na každou URL.
- Požadavky na různé URL běží paralelně.
- Na stejnou URL se nový požadavek odešle až po dokončení předchozího.
- Výsledky se zobrazují v přehledné tabulce v terminálu.
- Ukončení pomocí `CTRL+C` zachová poslední stav tabulky na obrazovce.


### ✅ Požadavky

- Go 1.20 nebo vyšší
- Internetové připojení

### 🧹 Instalace závislostí

```bash
go mod tidy
```

---

## 🌐 CLI Application in Go for Website Availability Monitoring

This CLI app monitors the availability and response metrics of given websites using periodic HTTP GET requests.

### ✨ Features

- Sends a GET request every 5 seconds to each provided URL.
- Processes all URLs in parallel.
- Only one request at a time per URL (sequential per URL, parallel across URLs).
- Terminal output is refreshed with a live statistics table.
- Exits gracefully with `CTRL+C`, preserving the last table on screen.



### ✅ Requirements

- Go 1.20 or higher
- Internet access

### 🧹 Install dependencies

```bash
go mod tidy
```


## 🏃‍♂️ Spuštění / Running

```bash
# Kompilace a spuštění / Compile and run
go run main.go [URL1] [URL2] [URL3]
```

### 🧪 Example Runs

```bash
go run main.go https://example.com https://seznam.cz
```

```bash
go run main.go https://google.com https://facebook.com https://youtube.com https://instagram.com https://wikipedia.org https://amazon.com https://apple.com https://microsoft.com https://linkedin.com https://reddit.com https://github.com https://stackoverflow.com https://paypal.com https://cnn.com https://bbc.com https://nytimes.com https://weather.com https://seznam.cz
```

---

## 🧪 Testování / Testing

Aplikace podporuje **end-to-end testy** pomocí [`httpmock`](https://github.com/jarcoal/httpmock). Testy lze spustit pomocí:

```bash
go test ./...
```

> This project avoids global state and is designed for parallel-safe testing.

---

## ✍️ Autor / Author

**Bc. Patrik Duch**


