package tinystring

// Global dictionary instance - populated with all translations using horizontal format
// Language order: EN, ES, ZH, HI, AR, PT, FR, DE, RU
//
// By using an anonymous struct, we define and initialize the dictionary in a single step,
// avoiding the need for a separate 'type dictionary struct' declaration.
// The usage API (e.g., D.Argument) remains unchanged.
var D = struct {
	Allowed       LocStr // "allowed"
	Argument      LocStr // "argument"
	Assign        LocStr // "assign"
	Assignable    LocStr // "assignable"
	Be            LocStr // "be"
	Binary        LocStr // "binary"
	Call          LocStr // "call"
	Cannot        LocStr // "cannot"
	Character     LocStr // "character"
	Coding        LocStr // "coding"
	Compilation   LocStr // "compilation"
	Configuration LocStr // "configuration"
	Debugging     LocStr // "debugging"
	Decimal       LocStr // "decimal"
	Delimiter     LocStr // "delimiter"
	Digit         LocStr // "digit"
	Element       LocStr // "element"
	Empty         LocStr // "empty"
	End           LocStr // "end"
	Exceeds       LocStr // "exceeds"
	Failed        LocStr // "failed"
	Field         LocStr // "field"
	Fields        LocStr // "fields"
	Format        LocStr // "format"
	Found         LocStr // "found"
	Handler       LocStr // "handler"
	Implemented   LocStr // "implemented"
	In            LocStr // "in"
	Index         LocStr // "index"
	Input         LocStr // "input"
	Install       LocStr // "install"
	Installation  LocStr // "installation"
	Invalid       LocStr // "invalid"
	Maximum       LocStr // "maximum"
	Method        LocStr // "method"
	Missing       LocStr // "missing"
	Mismatch      LocStr // "mismatch"
	Mode          LocStr // "mode"
	Modes         LocStr // "modes"
	More          LocStr // "more"
	Must          LocStr // "must"
	Negative      LocStr // "negative"
	Nil           LocStr // "null"
	NonNumeric    LocStr // "non-numeric"
	Not           LocStr // "not"
	NotOfType     LocStr // "not of type"
	Number        LocStr // "number"
	Numbers       LocStr // "numbers"
	Of            LocStr // "of"
	Options       LocStr // "options"
	Out           LocStr // "out"
	Overflow      LocStr // "overflow"
	Pointer       LocStr // "pointer"
	Point         LocStr // "point"
	Production    LocStr // "production"
	Provided      LocStr // "provided"
	Range         LocStr // "range"
	Required      LocStr // "required"
	Round         LocStr // "round"
	Seconds       LocStr // "seconds"
	Session       LocStr // "session"
	Slice         LocStr // "slice"
	String        LocStr // "string"
	Supported     LocStr // "supported"
	Switching     LocStr // "switching"
	Sync          LocStr // "sync"
	Time          LocStr // "time"
	To            LocStr // "to"
	Type          LocStr // "type"
	Unexported    LocStr // "unexported"
	Unknown       LocStr // "unknown"
	Unsigned      LocStr // "unsigned"
	Use           LocStr // "use"
	Valid         LocStr // "valid"
	Value         LocStr // "value"
	Zero          LocStr // "zero"
}{
	LocStr{"allowed", "permitido", "允许", "अनुमति", "مسموح", "permitido", "autorisé", "erlaubt", "разрешено"},
	LocStr{"argument", "argumento", "参数", "तर्क", "وسيط", "argumento", "argument", "Argument", "аргумент"},
	LocStr{"assign", "asignar", "分配", "असाइन", "تعيين", "atribuir", "assigner", "zuweisen", "присвоить"},
	LocStr{"assignable", "asignable", "可分配", "असाइन करने योग्य", "قابل للتعيين", "atributável", "assignable", "zuweisbar", "присваиваемый"},
	LocStr{"be", "ser", "是", "होना", "كون", "ser", "être", "sein", "быть"},
	LocStr{"binary", "binario", "二进制", "二进制", "ثنائي", "binário", "binaire", "binär", "двоичный"},
	LocStr{"call", "llamar", "调用", "कॉल", "استدعاء", "chamar", "appeler", "aufrufen", "вызвать"},
	LocStr{"cannot", "no puede", "不能", "नहीं कर सकते", "لا يمكن", "não pode", "ne peut pas", "kann nicht", "не может"},
	LocStr{"character", "caracter", "字符", "वर्ण", "حرف", "caractere", "caractère", "Zeichen", "символ"},
	LocStr{"coding", "codificación", "编码", "कोडिंग", "ترميز", "codificação", "codage", "kodierung", "кодирование"},
	LocStr{"compilation", "compilación", "编译", "संकलन", "تجميع", "compilação", "compilation", "kompilierung", "компиляция"},
	LocStr{"configuration", "configuración", "配置", "कॉन्फ़िगरेशन", "تكوين", "configuração", "configuration", "Konfiguration", "конфигурация"},
	LocStr{"debugging", "depuración", "调试", "डिबगिंग", "تصحيح", "depuração", "débogage", "debuggen", "отладка"},
	LocStr{"decimal", "decimal", "十进制", "दशमलव", "عشري", "decimal", "décimal", "Dezimal", "десятичная"},
	LocStr{"delimiter", "delimitador", "分隔符", "सीमांकक", "محدد", "delimitador", "délimiteur", "Trennzeichen", "разделитель"},
	LocStr{"digit", "dígito", "数字", "अंक", "رقم", "dígito", "chiffre", "Ziffer", "цифра"},
	LocStr{"element", "elemento", "元素", "एलिमेंट", "عنصر", "elemento", "élément", "Element", "элемент"},
	LocStr{"empty", "vacío", "空", "खाली", "فارغ", "vazio", "vide", "leer", "пустой"},
	LocStr{"end", "fin", "结束", "अंत", "نهاية", "fim", "fin", "Ende", "конец"},
	LocStr{"exceeds", "excede", "超过", "अधिक", "يتجاوز", "excede", "dépasse", "überschreitet", "превышает"},
	LocStr{"failed", "falló", "失败", "असफल", "فشل", "falhou", "échoué", "fehlgeschlagen", "не удалось"},
	LocStr{"field", "campo", "字段", "फील्ड", "حقل", "campo", "champ", "Feld", "поле"},
	LocStr{"fields", "campos", "字段集", "फील्ड्स", "حقول", "campos", "champs", "Felder", "поля"},
	LocStr{"format", "formato", "格式", "प्रारूप", "تنسيق", "formato", "format", "Fmt", "формат"},
	LocStr{"found", "encontrado", "找到", "मिला", "موجود", "encontrado", "trouvé", "gefunden", "найден"},
	LocStr{"handler", "manejador", "处理程序", "हैंडलर", "معالج", "manipulador", "gestionnaire", "Handler", "обработчик"},
	LocStr{"implemented", "implementado", "已实现", "कार्यान्वित", "مُنفذ", "implementado", "implémenté", "implementiert", "реализовано"},
	LocStr{"in", "en", "在", "में", "في", "em", "dans", "in", "в"},
	LocStr{"index", "índice", "索引", "इंडेक्स", "فهرس", "índice", "index", "Index", "индекс"},
	LocStr{"input", "entrada", "输入", "इनपुट", "إدخال", "entrada", "entrée", "eingabe", "ввод"},
	LocStr{"install", "instalar", "安装", "इंस्टॉल", "تثبيت", "instalar", "installer", "installieren", "установить"},
	LocStr{"installation", "instalación", "安装", "स्थापना", "التثبيت", "instalação", "installation", "installation", "установка"},
	LocStr{"invalid", "inválido", "无效", "अमान्य", "غير صالح", "inválido", "invalide", "ungültig", "недопустимый"},
	LocStr{"maximum", "máximo", "最大", "अधिकतम", "الحد الأقصى", "máximo", "maximum", "Maximum", "максимум"},
	LocStr{"method", "método", "方法", "विधि", "طريقة", "método", "méthode", "Methode", "метод"},
	LocStr{"missing", "falta", "缺少", "गुम", "مفقود", "ausente", "manquant", "fehlend", "отсутствует"},
	LocStr{"mismatch", "desajuste", "不匹配", "बेमेल", "عدم تطابق", "incompatibilidade", "incompatibilité", "Nichtübereinstimmung", "несоответствие"},
	LocStr{"mode", "modo", "模式", "मोड", "وضع", "modo", "mode", "Modus", "режим"},
	LocStr{"modes", "modos", "模式", "मोड", "أوضاع", "modos", "modes", "Modi", "режимы"},
	LocStr{"more", "más", "更多", "अधिक", "أكثر", "mais", "plus", "mehr", "больше"},
	LocStr{"must", "debe", "必须", "चाहिए", "يجب", "deve", "doit", "muss", "должен"},
	LocStr{"negative", "negativo", "负", "नकारात्मक", "سالب", "negativo", "négatif", "negativ", "отрицательный"},
	LocStr{"nil", "nulo", "空值", "शून्य मान", "قيمة فارغة", "nulo", "nul", "Nullwert", "нуль"},
	LocStr{"non-numeric", "no numérico", "非数字", "गैर-संख्यात्मक", "غير رقمي", "não numérico", "non numérique", "nicht numerisch", "нечисловой"},
	LocStr{"not", "no", "不", "नहीं", "ليس", "não", "pas", "nicht", "не"},
	LocStr{"not of type", "no es del tipo", "不是类型", "प्रकार नहीं है", "ليس من النوع", "não é do tipo", "n'est pas du type", "ist nicht vom Typ", "не того типа"},
	LocStr{"number", "número", "数字", "संख्या", "رقم", "número", "nombre", "Zahl", "число"},
	LocStr{"numbers", "números", "数字", "संख्याएं", "أرقام", "números", "nombres", "Zahlen", "числа"},
	LocStr{"of", "de", "的", "का", "من", "de", "de", "von", "из"},
	LocStr{"options", "opciones", "选项", "विकल्प", "خيارات", "opções", "options", "Optionen", "опции"},
	LocStr{"out", "fuera", "出", "बाहर", "خارج", "fora", "hors", "aus", "вне"},
	LocStr{"overflow", "desbordamiento", "溢出", "ओवरफ्लो", "فيض", "estouro", "débordement", "Überlauf", "переполнение"},
	LocStr{"pointer", "puntero", "指针", "पॉइंटर", "مؤشر", "ponteiro", "pointeur", "Zeiger", "указатель"},
	LocStr{"point", "punto", "点", "बिंदु", "نقطة", "ponto", "point", "Punkt", "точка"},
	LocStr{"production", "producción", "生产", "उत्पादन", "إنتاج", "produção", "production", "produktion", "производство"},
	LocStr{"provided", "proporcionado", "提供", "प्रदान किया गया", "مقدم", "fornecido", "fourni", "bereitgestellt", "предоставлено"},
	LocStr{"range", "rango", "范围", "रेंज", "نطاق", "intervalo", "plage", "Bereich", "диапазон"},
	LocStr{"required", "requerido", "必需", "आवश्यक", "مطلوب", "necessário", "requis", "erforderlich", "обязательный"},
	LocStr{"round", "redondear", "圆", "गोल", "جولة", "arredondar", "arrondir", "runden", "округлить"},
	LocStr{"seconds", "segundos", "秒", "सेकंड", "ثواني", "segundos", "secondes", "Sekunden", "секунды"},
	LocStr{"session", "sesión", "会话", "सत्र", "جلسة", "sessão", "session", "Sitzung", "сессия"},
	LocStr{"slice", "segmento", "切片", "स्लाइस", "شريحة", "fatia", "tranche", "Scheibe", "срез"},
	LocStr{"string", "cadena", "字符串", "स्ट्रिंग", "سلسلة", "string", "chaîne", "Zeichenkette", "строка"},
	LocStr{"supported", "soportado", "支持", "समर्थित", "مدعوم", "suportado", "pris en charge", "unterstützt", "поддерживается"},
	LocStr{"switching", "cambiando", "切换", "स्विच", "تبديل", "mudando", "changement", "wechseln", "переключение"},
	LocStr{"sync", "sincronización", "同步", "सिंक", "مزامنة", "sincronização", "synchronisation", "Synchronisierung", "синхронизация"},
	LocStr{"time", "tiempo", "时间", "समय", "وقت", "tempo", "temps", "Zeit", "время"},
	LocStr{"to", "para", "到", "को", "إلى", "para", "à", "zu", "к"},
	LocStr{"type", "tipo", "类型", "प्रकार", "نوع", "tipo", "type", "Typ", "тип"},
	LocStr{"unexported", "no exportado", "未导出", "गैर-निर्यातित", "غير مصدر", "não exportado", "non exporté", "nicht exportiert", "неэкспортированный"},
	LocStr{"unknown", "desconocido", "未知", "अज्ञात", "غير معروف", "desconhecido", "inconnu", "unbekannt", "неизвестный"},
	LocStr{"unsigned", "sin signo", "无符号", "अहस्ताक्षरित", "غير موقع", "sem sinal", "non signé", "vorzeichenlos", "безzнаковый"},
	LocStr{"use", "usar", "使用", "उपयोग", "استخدام", "usar", "utiliser", "verwenden", "использовать"},
	LocStr{"valid", "válido", "有效", "वैध", "صحيح", "válido", "valide", "gültig", "действительный"},
	LocStr{"value", "valor", "值", "मूल्य", "قيمة", "valor", "valeur", "Wert", "значение"},
	LocStr{"zero", "cero", "零", "शून्य", "صفر", "zero", "zéro", "Null", "ноль"},
}
