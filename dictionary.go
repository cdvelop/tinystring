package tinystring

// Global dictionary instance - populated with all translations using horizontal format
// Language order: EN, ES, ZH, HI, AR, PT, FR, DE, RU
//
// By using an anonymous struct, we define and initialize the dictionary in a single step,
// avoiding the need for a separate 'type dictionary struct' declaration.
// The usage API (e.g., D.Argument) remains unchanged.
var D = struct {
	// Basic words sorted alphabetically for maximum reusability
	Allowed    LocStr // "allowed"
	Argument   LocStr // "argument"
	At         LocStr // "at"
	Base       LocStr // "base"
	Boolean    LocStr // "boolean"
	Cannot     LocStr // "cannot"
	Character  LocStr // "character"
	Decimal    LocStr // "decimal"
	Delimiter  LocStr // "delimiter"
	Digit      LocStr // "digit"
	Empty      LocStr // "empty"
	End        LocStr // "end"
	Float      LocStr // "float"
	For        LocStr // "for"
	Format     LocStr // "format"
	Found      LocStr // "found"
	In         LocStr // "in"
	Integer    LocStr // "integer"
	Invalid    LocStr // "invalid"
	Missing    LocStr // "missing"
	Negative   LocStr // "negative"
	NonNumeric LocStr // "non-numeric"
	Not        LocStr // "not"
	Number     LocStr // "number"
	Numbers    LocStr // "numbers"
	Of         LocStr // "of"
	Out        LocStr // "out"
	Overflow   LocStr // "overflow"
	Range      LocStr // "range"
	Required   LocStr // "required"
	Round      LocStr // "round"
	Specifier  LocStr // "specifier"
	String     LocStr // "string"
	Supported  LocStr // "supported"
	Type       LocStr // "type"
	Unknown    LocStr // "unknown"
	Unsigned   LocStr // "unsigned"
	Value      LocStr // "value"
}{
	Allowed:    LocStr{"allowed", "permitido", "允许", "अनुमति", "مسموح", "permitido", "autorisé", "erlaubt", "разрешено"},
	Argument:   LocStr{"argument", "argumento", "参数", "तर्क", "وسيط", "argumento", "argument", "Argument", "аргумент"},
	At:         LocStr{"at", "en", "在", "पर", "في", "em", "à", "bei", "в"},
	Base:       LocStr{"base", "base", "进制", "आधार", "قاعدة", "base", "base", "Basis", "основание"},
	Boolean:    LocStr{"boolean", "booleano", "布尔", "बूलियन", "منطقي", "booleano", "booléen", "boolescher", "логический"},
	Cannot:     LocStr{"cannot", "no puede", "不能", "नहीं कर सकते", "لا يمكن", "não pode", "ne peut pas", "kann nicht", "не может"},
	Character:  LocStr{"character", "caracter", "字符", "वर्ण", "حرف", "caractere", "caractère", "Zeichen", "символ"},
	Decimal:    LocStr{"decimal", "decimal", "十进制", "दशमलव", "عشري", "decimal", "décimal", "Dezimal", "десятичная"},
	Delimiter:  LocStr{"delimiter", "delimitador", "分隔符", "सीमांकक", "محدد", "delimitador", "délimiteur", "Trennzeichen", "разделитель"},
	Digit:      LocStr{"digit", "dígito", "数字", "अंक", "رقم", "dígito", "chiffre", "Ziffer", "цифра"},
	Empty:      LocStr{"empty", "vacío", "空", "खाली", "فارغ", "vazio", "vide", "leer", "пустой"},
	End:        LocStr{"end", "fin", "结束", "अंत", "نهاية", "fim", "fin", "Ende", "конец"},
	Float:      LocStr{"float", "flotante", "浮点", "फ्लोट", "عائم", "flutuante", "flottant", "Gleitkomma", "число с плавающей точкой"},
	For:        LocStr{"for", "para", "为", "के लिए", "لـ", "para", "pour", "für", "для"},
	Format:     LocStr{"format", "formato", "格式", "प्रारूप", "تنسيق", "formato", "format", "Fmt", "формат"},
	Found:      LocStr{"found", "encontrado", "找到", "मिला", "موجود", "encontrado", "trouvé", "gefunden", "найден"},
	Integer:    LocStr{"integer", "entero", "整数", "पूर्णांक", "عدد صحيح", "inteiro", "entier", "ganze Zahl", "целое число"},
	Invalid:    LocStr{"invalid", "inválido", "无效", "अमान्य", "غير صالح", "inválido", "invalide", "ungültig", "недопустимый"},
	Missing:    LocStr{"missing", "falta", "缺少", "गुम", "مفقود", "ausente", "manquant", "fehlend", "отсутствует"},
	Negative:   LocStr{"negative", "negativo", "负", "नकारात्मक", "سالب", "negativo", "négatif", "negativ", "отрицательный"},
	NonNumeric: LocStr{"non-numeric", "no numérico", "非数字", "गैर-संख्यात्मक", "غير رقمي", "não numérico", "non numérique", "nicht numerisch", "нечисловой"},
	Not:        LocStr{"not", "no", "不", "नहीं", "ليس", "não", "pas", "nicht", "не"},
	Number:     LocStr{"number", "número", "数字", "संख्या", "رقم", "número", "nombre", "Zahl", "число"},
	Numbers:    LocStr{"numbers", "números", "数字", "संख्याएं", "أرقام", "números", "nombres", "Zahlen", "числа"},
	Of:         LocStr{"of", "de", "的", "का", "من", "de", "de", "von", "из"},
	Out:        LocStr{"out", "fuera", "出", "बाहर", "خارج", "fora", "hors", "aus", "вне"},
	Overflow:   LocStr{"overflow", "desbordamiento", "溢出", "ओवरफ्लो", "فيض", "estouro", "débordement", "Überlauf", "переполнение"},
	Range:      LocStr{"range", "rango", "范围", "रेंज", "نطاق", "intervalo", "plage", "Bereich", "диапазон"},
	Required:   LocStr{"required", "requerido", "必需", "आवश्यक", "مطلوب", "necessário", "requis", "erforderlich", "обязательный"},
	Round:      LocStr{"round", "redondear", "圆", "गोल", "جولة", "arredondar", "arrondir", "runden", "округлить"},
	Specifier:  LocStr{"specifier", "especificador", "说明符", "निर्दिष्टकर्ता", "محدد", "especificador", "spécificateur", "Spezifizierer", "спецификатор"},
	String:     LocStr{"string", "cadena", "字符串", "स्ट्रिंग", "سلسلة", "string", "chaîne", "Zeichenkette", "строка"},
	Supported:  LocStr{"supported", "soportado", "支持", "समर्थित", "مدعوم", "suportado", "pris en charge", "unterstützt", "поддерживается"},
	Type:       LocStr{"type", "tipo", "类型", "प्रकार", "نوع", "tipo", "type", "Typ", "тип"},
	Unknown:    LocStr{"unknown", "desconocido", "未知", "अज्ञात", "غير معروف", "desconhecido", "inconnu", "unbekannt", "неизвестный"},
	Unsigned:   LocStr{"unsigned", "sin signo", "无符号", "अहस्ताक्षरित", "غير موقع", "sem sinal", "non signé", "vorzeichenlos", "безzнаковый"},
	Value:      LocStr{"value", "valor", "值", "मूल्य", "قيمة", "valor", "valeur", "Wert", "значение"},
}
