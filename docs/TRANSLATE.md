# 🌍 TinyString: Multilingual Message System

**TinyString** is a lightweight, dependency-free multilingual dictionary for generating composable error and validation messages. It supports 9 major languages:

**Supported Languages:**

- 🇺🇸 EN (English, default)
- 🇪🇸 ES (Spanish)
- 🇨🇳 ZH (Chinese)
- 🇮🇳 HI (Hindi)
- 🇸🇦 AR (Arabic)
- 🇧🇷 PT (Portuguese)
- 🇫🇷 FR (French)
- 🇩🇪 DE (German)
- 🇷🇺 RU (Russian)

---

## 🚀 Features

- ✅ 9 Languages with 35+ essential terms
- 🧱 Composable error messages from dictionary words
- 🌐 Auto-detects system/browser language
- 🛠️ Language override (global or inline)
- 🧩 Custom dictionaries for domain-specific terms
- 🔒 Zero external dependencies
- ⚙️ Compatibility: Go + TinyGo (WASM ready)


---

## 🌍 Basic Usage

```go
// Set global language to Spanish (using lang constant), returns "ES"
code := OutLang(ES) // returns "ES"
code = OutLang()    // auto-detects and returns code (e.g. "EN")
// If an error occurs or the language is not recognized, "EN" is always returned by default

// Usage examples:

// return strings
// Force to Spanish (ES) only for this response, not globally.
// Useful for personalized user replies.
msg := Translate(ES, D.Format, D.Invalid).String()
// → "formato inválido"

// Capitalize translation (first letter of each word uppercase)
msgCap := Translate(ES, D.Format, D.Invalid).Capitalize().String()
// → "Formato Inválido"

// Force French
err = Err(FR, D.Empty, D.String)
// → "vide chaîne" (forced French)


// Use global language (e.g. Spanish) for error messages
// return error
err := Err(D.Format, D.Invalid)
// → "formato inválido"

err = Err(D.Number, D.Negative, D.Not, D.Supported)
// → "número negativo no soportado"

err = Err(D.Cannot, D.Round, D.Value, D.NonNumeric)
// → "no se puede redondear valor no numérico"
```


---

## ⚡ Memory Management

`Translate` returns a pooled `*Conv` object for high performance.

- **Automatic Release**: Calling `.String()` or `.Apply()` automatically returns the object to the pool.
- **Manual Release**: If you use `.Bytes()` or keep the object, you **MUST** call `.PutConv()` manually.

```go
// ✅ Automatic release (Recommended)
msg := Translate(D.Format).String()

// ⚠️ Manual release required
c := Translate(D.Format)
bytes := c.Bytes()
// ... use bytes ...
c.PutConv() // Don't forget this!
```

### 🚀 Zero-Allocation Performance

For hot paths requiring zero allocations, pass **pointers** to `LocStr`:

```go
// Standard usage (1 alloc/op)
msg := Translate(D.Format, D.Invalid).String()

// Zero-allocation usage (0 allocs/op)
msg := Translate(&D.Format, &D.Invalid).String()
```

**Benchmark Results:**
- `Translate(D.Format)`: 1 alloc/op, 144 B/op
- `Translate(&D.Format)`: **0 allocs/op**, 0 B/op

This optimization is useful when allocation-free operation is critical.

---



## 🌐 Minimal HTTP API Example

```go
import (
    "encoding/json"
    "net/http"
    . "github.com/tinywasm/fmt"
)

func handler(w http.ResponseWriter, r *http.Request) {
    lang := r.URL.Query().Get("lang") // e.g. ?lang=ES
    resp := map[string]string{
        "error": Translate(lang, D.Format, D.Invalid).String(),
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}
```

## 🧩 Custom Dictionary

Define domain-specific words:

```go
type MyDict struct {
    User  LocStr
    Email LocStr
}

var MD = MyDict{
    User:  LocStr{"user", "usuario", "usuário", "utilisateur", "пользователь", "Benutzer", "utente", "उपयोगकर्ता", "用户"},
    Email: LocStr{"email", "correo", "email", "email", "البريد الإلكتروني", "Courriel", "Эл. адрес", "电邮", "ईमेल"},
}

// Usage with custom dictionary
err := Err("es",D.Format, MD.Email, MD.User, D.Invalid)
// → "formato correo usuario inválido"
```

---

## ✅ Validation Example

```go
validate := func(input string) error {
    if input == "" {
        return Err(D.Empty, D.String, D.Not, D.Supported)
    }
    if _, err := Convert(input).Int(); err != nil {
        return Err(D.Invalid, D.Number, D.Format)
    }
    return nil
}
```

---

## 🔍 Dictionary Reference

See [`dictionary.go`](../dictionary.go) for built-in words.
Combine `D.` (default terms) and custom dictionaries for flexible messaging.


