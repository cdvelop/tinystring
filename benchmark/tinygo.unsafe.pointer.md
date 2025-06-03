TinyGo 0.37.0, al compilar a WebAssembly con interoperabilidad con JavaScript usando syscall/js, soporta el uso de `unsafe.Pointer`, `unsafe.String()` y `unsafe.SliceData()` de la librería estándar. 

# Soporte en TinyGo 0.37.0

TinyGo 0.37.0 (Go 1.24) **sí soporta** `unsafe.Pointer`, `unsafe.String` y `unsafe.SliceData`. De hecho, estas funciones (añadidas en Go 1.20) fueron incorporadas al compilador TinyGo en enero de 2023. La documentación oficial de TinyGo indica explícitamente que estas funciones están *“fully supported”* (es decir, “totalmente compatibles”). Se recomienda usarlas para convertir entre punteros y slices/strings sin copiar, en lugar de usar los antiguos `reflect.StringHeader` o `reflect.SliceHeader`.

# Restricciones y errores conocidos

No hay reportes de **errores de compilación** específicos en el target WASM/js al usar estas funciones. TinyGo compila el código WASM estándar con `syscall/js` sin rechazar `unsafe.Pointer`, `unsafe.String` o `unsafe.SliceData`. Un detalle importante es que TinyGo define internamente los campos `Len`/`Cap` de un slice como `uintptr` en lugar de `int`, por lo que cualquier código legacy que manipule directamente los *SliceHeader* puede fallar. Por ello la guía advierte **no usar** esos *Header* y sí las funciones `unsafe.SliceData`/`unsafe.String` como se indicó arriba. En la interoperabilidad WASM–JavaScript se debe respetar la semántica de memoria de WebAssembly (por ejemplo, pasar offsets o usar `syscall/js.CopyBytesToJS`), pero esto es independiente del soporte de TinyGo para las funciones `unsafe`.

# Workarounds documentados

La “solución” oficial consiste en usar las funciones `unsafe.Slice`, `unsafe.SliceData` y `unsafe.String` para hacer conversiones sin copiar datos. Por ejemplo, TinyGo muestra cómo convertir un `string` a `[]byte` con `unsafe.Slice(unsafe.StringData(s), len(s))`. No se requieren hacks adicionales: basta usar estas funciones de `unsafe` como indica la documentación. En resumen, TinyGo 0.37.0 admite plenamente `unsafe.Pointer`, `unsafe.String()` y `unsafe.SliceData()` en proyectos WASM/JS; solo hay que tener en cuenta la recomendación de TinyGo de evitar los reflect.Headers y preferir las funciones `unsafe` mencionadas.

**Fuentes:** Documentación oficial y notas de TinyGo sobre compatibilidad y el cambio Go 1.20; changelog de TinyGo (enero 2023) incorporando `unsafe.SliceData` y `unsafe.String`.
