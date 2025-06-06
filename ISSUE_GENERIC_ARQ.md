Eres un asistente especializado en Go y TinyGo/WebAssembly. A continuación tienes un repositorio completo en GitHub: github.com/cdvelop/tinystring. Tu objetivo es refactorizar todo el código de ese paquete siguiendo estas instrucciones:

1. **Mantener la función pública Convert(v any) *conv**:
   - La función `Convert` debe seguir existiendo con exactamente ese nombre y firma pública.  
   - Internamente, debe usar un constructor que hemos de crear (por ejemplo, `newConv(...)`) que reciba “opciones” (functional options).

2. **Reemplazar toda la lógica de conversión de tipos numéricos por genéricos**:
   - Actualmente hay métodos como `handleIntTypes`, `handleUintTypes`, `handleFloatTypes`, y sus correspondientes versiones “ForAny2s” con switches que distinguen `int`, `int8`, etc.  
   - Debes eliminar esos switches repetitivos y crear **3 funciones genéricas**:
     ```go
     type anyInt interface { ~int | ~int8 | ~int16 | ~int32 | ~int64 }
     type anyUint interface { ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 }
     type anyFloat interface { ~float32 | ~float64 }
 
     func (c *conv) genInt[T anyInt](v T) { … }
     func (c *conv) genUint[T anyUint](v T) { … }
     func (c *conv) genFloat[T anyFloat](v T) { … }
 
     func (c *conv) genAny2sInt[T anyInt](v T) { … }
     func (c *conv) genAny2sUint[T anyUint](v T) { … }
     func (c *conv) genAny2sFloat[T anyFloat](v T) { … }
     ```


3. **Crear un constructor `newConv(opts ...convOption) *conv` con functional options**:
   - El constructor inicializa un struct `conv` con valores por defecto (por ejemplo, `separator: "_"`).  
   - Todas las ramas de `convInit` deben convertirse en una opción:
     - Por ejemplo, `withValue(v any) convOption` que encapsula la lógica de decidir “si es string, si es anyInt, anyUint, anyFloat, bool, []string, *string, errorType, o fallback a string”.  
     - Otra opción: `withSeparator(sep string) convOption` para cambiar el separador.  
     - Otra opción: `withRoundDown(b bool) convOption`.  
     - Y así con todas las configuraciones que antes estaban embebidas en `convInit` o en métodos dispersos.  
   - Después, la función pública `Convert(v any)` debe simplemente llamar a `newConv(withValue(v))`.

4. **Convertir métodos de transformación en métodos encadenables (builder-like)**:
   - En lugar de `func (t *conv) tmap(mappings []charMapping) *conv`, mantenlo como un método que devuelve `*conv` para que pueda encadenarse:
     ```go
     func (c *conv) tMap(m []charMapping) *conv { … return c }
     ```
   - Lo mismo para cualquier método que transforme palabras (`transformWord`) o que afecte internamente a `conv`.

5. **Privacidad de tipos e interfaces**:
   - Todos los tipos internos (por ejemplo, `vTpe`, `charMapping`, `wordTransform`, y cualquier otro helper) deben quedar **no exportados** (nombre en minúscula).  
   - Solo exporta lo mínimo necesario: 
     ```go
     // Convert es la única función pública para inicializar la cadena.
     func Convert(v any) *conv

     // Si quieres exponer la opción de encadenar transformaciones, deja públicos los métodos del receiver *conv
     // Ejemplo: func (c *conv) tMap(m []charMapping) *conv
     // Pero deja charMapping sin exportar; en su lugar define un tipo público si deseas que el usuario cree mapeos.
     ```
   - Revisa cada tipo y función: si no debe formar parte de la API pública, ponlo en minúscula.

6. **Mantener la lógica interna (split, getString, addRne2Buf, etc.) igual de eficiente**:
   - No elimines ni modifiques esa lógica, solo reorganízala dentro del nuevo constructor y métodos genéricos.
   - Asegúrate de que la funcionalidad sea equivalente: las conversiones de entero a string, floats, booleans, slices de strings, etc., sigan comportándose igual.

7. **Actualizar documentación y comentarios**:
   - Ajusta los comentarios de GoDoc para que reflejen el nuevo patrón de uso.  
   - Ejemplo de uso en GoDoc: 
     ```go
     // Convert convierte cualquier valor a un conv, permitiendo luego encadenar transformaciones:
     //   result := Convert(123).tMap(misMapas).String()
     ```

8. **Instrucciones de salida**:
   - Devuélveme el diff completo o los archivos refactorizados como si fuese un “pull request” con todos los cambios en `github.com/cdvelop/tinystring`.  
   - Mantén la estructura de carpetas intacta; solo modifica los archivos `.go`.  
   - Asegúrate de que el nuevo código compile con Go 1.18+ y con TinyGo ≥0.27 para WebAssembly.

--- 

Opcionalmente, puedes incluir esta línea al principio del prompt para que el modelo clone el repositorio antes de refactorizar (dependiendo de tu entorno de ejecución del LLM):


```go 
package tinystring

// ——— Definiciones de tipos y constantes como antes ———
type vTpe uint8
const (
    typeStr vTpe = iota
    typeInt
    typeUint
    typeFloat
    typeBool
    typeStrSlice
    typeStrPtr
    typeErr
)

type charMapping struct{ from, to rune }
type wordTransform int
const (
    toLower wordTransform = iota
    toUpper
)

// ——— Interfases genéricas ———
type anyInt interface{ ~int | ~int8 | ~int16 | ~int32 | ~int64 }
type anyUint interface{ ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 }
type anyFloat interface{ ~float32 | ~float64 }

// ——— Estructura central ———
type conv struct {
    stringVal      string
    intVal         int64
    uintVal        uint64
    floatVal       float64
    boolVal        bool
    stringSliceVal []string
    stringPtrVal   *string
    vTpe           vTpe
    roundDown      bool
    separator      string
    tmpStr         string
    lastConvType   vTpe
    err            errorType
}

// ——— Funciones genéricas de asignación ———
func (c *conv) genHandleInt[T anyInt](v T) {
    c.intVal = int64(v)
    c.vTpe = typeInt
}

func (c *conv) genHandleUint[T anyUint](v T) {
    c.uintVal = uint64(v)
    c.vTpe = typeUint
}

func (c *conv) genHandleFloat[T anyFloat](v T) {
    c.floatVal = float64(v)
    c.vTpe = typeFloat
}

// ——— Funciones genéricas para any2s ———
func (c *conv) genAny2sInt[T anyInt](v T) {
    c.intVal = int64(v)
    c.fmtInt(10)
}
func (c *conv) genAny2sUint[T anyUint](v T) {
    c.uintVal = uint64(v)
    c.fmtUint(10)
}
func (c *conv) genAny2sFloat[T anyFloat](v T) {
    c.floatVal = float64(v)
    c.f2s()
}

// ——— Opciones para construir/editar conv ———
type ConvOption func(*conv)

// Inicializar con valor “any”
func WithValue(v any) ConvOption {
    return func(c *conv) {
        if v == nil {
            c.stringVal, c.vTpe = "", typeStr
            return
        }
        switch val := v.(type) {
        case string:
            c.stringVal, c.vTpe = val, typeStr
        case anyInt:
            c.genHandleInt(val)
        case anyUint:
            c.genHandleUint(val)
        case anyFloat:
            c.genHandleFloat(val)
        case bool:
            c.boolVal, c.vTpe = val, typeBool
        case []string:
            c.stringSliceVal, c.vTpe = val, typeStrSlice
        case *string:
            c.stringPtrVal = val
            c.stringVal, c.vTpe = *val, typeStrPtr
        case errorType:
            c.err, c.vTpe = val, typeErr
        default:
            // Fallback a string
            c.vTpe = typeStr
            c.any2s(val)
        }
    }
}

func WithSeparator(sep string) ConvOption {
    return func(c *conv) {
        c.separator = sep
    }
}

func WithRoundDown(rd bool) ConvOption {
    return func(c *conv) {
        c.roundDown = rd
    }
}

// (… más opciones según necesites: WithTMap, WithTransformWord, etc.)

// ——— Constructor principal ———
func NewConv(opts ...ConvOption) *conv {
    c := &conv{
        separator: "_", // valor por defecto
    }
    for _, opt := range opts {
        opt(c)
    }
    return c
}

// Ejemplo de un método encadenable para aplicar mapas de caracteres
func (c *conv) WithTMap(mappings []charMapping) *conv {
    str := c.getString()
    buf := make([]byte, 0, len(str)*2)
    hc := false
    for _, r := range str {
        mapped := false
        for _, m := range mappings {
            if r == m.from {
                buf = addRne2Buf(buf, m.to)
                hc = true; mapped = true
                break
            }
        }
        if !mapped {
            buf = addRne2Buf(buf, r)
        }
    }
    if !hc {
        return c
    }
    c.setString(string(buf))
    return c
}
// (… el resto de métodos—split, transformWord, getString, addRne2Buf, etc.—igual que antes)
```


