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
OutLang(ES) // Set global language to Spanish (using lang constant)

// Usage examples:

// return strings
msg := T(ES, D.Format, D.Invalid)
// → "formato inválido"

// return error
err := Err(D.Format, D.Invalid)
// → "formato inválido"

err = Err(D.Number, D.Negative, D.Not, D.Supported)
// → "número negativo no soportado"

// Force French
err = Err(FR, D.Empty, D.String)
// → "vide chaîne" (forced French)

OutLang() // Auto-detect system/browser language
err = Err(D.Cannot, D.Round, D.NonNumeric, D.Value)
```

---


## 🌐 Minimal HTTP API Example

```go
import (
    "encoding/json"
    "net/http"
    . "github.com/cdvelop/tinystring"
)

func handler(w http.ResponseWriter, r *http.Request) {
    lang := r.URL.Query().Get("lang") // e.g. ?lang=ES
    resp := map[string]string{
        "error": T(lang, D.Format, D.Invalid),
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

err := Err(D.Format, MD.Email, D.Invalid)
// → "formato correo inválido"
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


