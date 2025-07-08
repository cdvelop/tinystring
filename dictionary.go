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
	At            LocStr // "at"
	Base          LocStr // "base"
	Boolean       LocStr // "boolean"
	Cannot        LocStr // "cannot"
	Character     LocStr // "character"
	Configuration LocStr // "configuration"
	Decimal       LocStr // "decimal"
	Delimiter     LocStr // "delimiter"
	Digit         LocStr // "digit"
	Empty         LocStr // "empty"
	End           LocStr // "end"
	Float         LocStr // "float"
	For           LocStr // "for"
	Format        LocStr // "format"
	Found         LocStr // "found"
	Handler       LocStr // "handler"
	In            LocStr // "in"
	Integer       LocStr // "integer"
	Invalid       LocStr // "invalid"
	Missing       LocStr // "missing"
	More          LocStr // "more"
	Mutex         LocStr // "mutex"
	Nano          LocStr // "nano"
	Negative      LocStr // "negative"
	NonNumeric    LocStr // "non-numeric"
	Not           LocStr // "not"
	Number        LocStr // "number"
	Numbers       LocStr // "numbers"
	Of            LocStr // "of"
	Options       LocStr // "options"
	Out           LocStr // "out"
	Overflow      LocStr // "overflow"
	Point         LocStr // "point"
	Range         LocStr // "range"
	Required      LocStr // "required"
	Round         LocStr // "round"
	Seconds       LocStr // "seconds"
	Session       LocStr // "session"
	Specifier     LocStr // "specifier"
	String        LocStr // "string"
	Supported     LocStr // "supported"
	Sync          LocStr // "sync"
	Time          LocStr // "time"
	Type          LocStr // "type"
	Unknown       LocStr // "unknown"
	Unsigned      LocStr // "unsigned"
	Value         LocStr // "value"
}{
	LocStr{"allowed", "permitido", "允许", "अनुमति", "مسموح", "permitido", "autorisé", "erlaubt", "разрешено"},
	LocStr{"argument", "argumento", "参数", "तर्क", "وسيط", "argumento", "argument", "Argument", "аргумент"},
	LocStr{"at", "en", "在", "पर", "في", "em", "à", "bei", "в"},
	LocStr{"base", "base", "进制", "आधार", "قاعدة", "base", "base", "Basis", "основание"},
	LocStr{"boolean", "booleano", "布尔", "बूलियन", "منطقي", "booleano", "booléen", "boolescher", "логический"},
	LocStr{"cannot", "no puede", "不能", "नहीं कर सकते", "لا يمكن", "não pode", "ne peut pas", "kann nicht", "не может"},
	LocStr{"character", "caracter", "字符", "वर्ण", "حرف", "caractere", "caractère", "Zeichen", "символ"},
	LocStr{"configuration", "configuración", "配置", "कॉन्फ़िगरेशन", "تكوين", "configuração", "configuration", "Konfiguration", "конфигурация"},
	LocStr{"decimal", "decimal", "十进制", "दशमलव", "عشري", "decimal", "décimal", "Dezimal", "десятичная"},
	LocStr{"delimiter", "delimitador", "分隔符", "सीमांकक", "محدد", "delimitador", "délimiteur", "Trennzeichen", "разделитель"},
	LocStr{"digit", "dígito", "数字", "अंक", "رقم", "dígito", "chiffre", "Ziffer", "цифра"},
	LocStr{"empty", "vacío", "空", "खाली", "فارغ", "vazio", "vide", "leer", "пустой"},
	LocStr{"end", "fin", "结束", "अंत", "نهاية", "fim", "fin", "Ende", "конец"},
	LocStr{"float", "flotante", "浮点", "फ्लोट", "عائم", "flutuante", "flottant", "Gleitkomma", "число с плавающей точкой"},
	LocStr{"for", "para", "为", "के लिए", "لـ", "para", "pour", "für", "для"},
	LocStr{"format", "formato", "格式", "प्रारूप", "تنسيق", "formato", "format", "Fmt", "формат"},
	LocStr{"found", "encontrado", "找到", "मिला", "موجود", "encontrado", "trouvé", "gefunden", "найден"},
	LocStr{"handler", "manejador", "处理程序", "हैंडलर", "معالج", "manipulador", "gestionnaire", "Handler", "обработчик"},
	LocStr{"in", "en", "在", "में", "في", "em", "dans", "in", "в"},
	LocStr{"integer", "entero", "整数", "पूर्णांक", "عدد صحيح", "inteiro", "entier", "ganze Zahl", "целое число"},
	LocStr{"invalid", "inválido", "无效", "अमान्य", "غير صالح", "inválido", "invalide", "ungültig", "недопустимый"},
	LocStr{"missing", "falta", "缺少", "गुम", "مفقود", "ausente", "manquant", "fehlend", "отсутствует"},
	LocStr{"more", "más", "更多", "अधिक", "أكثر", "mais", "plus", "mehr", "больше"},
	LocStr{"mutex", "mutex", "互斥锁", "म्यूटेक्स", "قفل", "mutex", "mutex", "Mutex", "мьютекс"},
	LocStr{"nano", "nano", "纳秒", "नैनो", "نانو", "nano", "nano", "Nano", "нано"},
	LocStr{"negative", "negativo", "负", "नकारात्मक", "سالب", "negativo", "négatif", "negativ", "отрицательный"},
	LocStr{"non-numeric", "no numérico", "非数字", "गैर-संख्यात्मक", "غير رقمي", "não numérico", "non numérique", "nicht numerisch", "нечисловой"},
	LocStr{"not", "no", "不", "नहीं", "ليس", "não", "pas", "nicht", "не"},
	LocStr{"number", "número", "数字", "संख्या", "رقم", "número", "nombre", "Zahl", "число"},
	LocStr{"numbers", "números", "数字", "संख्याएं", "أرقام", "números", "nombres", "Zahlen", "числа"},
	LocStr{"of", "de", "的", "का", "من", "de", "de", "von", "из"},
	LocStr{"options", "opciones", "选项", "विकल्प", "خيارات", "opções", "options", "Optionen", "опции"},
	LocStr{"out", "fuera", "出", "बाहर", "خارج", "fora", "hors", "aus", "вне"},
	LocStr{"overflow", "desbordamiento", "溢出", "ओवरफ्लो", "فيض", "estouro", "débordement", "Überlauf", "переполнение"},
	LocStr{"point", "punto", "点", "बिंदु", "نقطة", "ponto", "point", "Punkt", "точка"},
	LocStr{"range", "rango", "范围", "रेंज", "نطاق", "intervalo", "plage", "Bereich", "диапазон"},
	LocStr{"required", "requerido", "必需", "आवश्यक", "مطلوب", "necessário", "requis", "erforderlich", "обязательный"},
	LocStr{"round", "redondear", "圆", "गोल", "جولة", "arredondar", "arrondir", "runden", "округлить"},
	LocStr{"seconds", "segundos", "秒", "सेकंड", "ثواني", "segundos", "secondes", "Sekunden", "секунды"},
	LocStr{"session", "sesión", "会话", "सत्र", "جلسة", "sessão", "session", "Sitzung", "сессия"},
	LocStr{"specifier", "especificador", "说明符", "निर्दिष्टकर्ता", "محدد", "especificador", "spécificateur", "Spezifizierer", "спецификатор"},
	LocStr{"string", "cadena", "字符串", "स्ट्रिंग", "سلسلة", "string", "chaîne", "Zeichenkette", "строка"},
	LocStr{"supported", "soportado", "支持", "समर्थित", "مدعوم", "suportado", "pris en charge", "unterstützt", "поддерживается"},
	LocStr{"sync", "sincronización", "同步", "सिंक", "مزامنة", "sincronização", "synchronisation", "Synchronisierung", "синхронизация"},
	LocStr{"time", "tiempo", "时间", "समय", "وقت", "tempo", "temps", "Zeit", "время"},
	LocStr{"type", "tipo", "类型", "प्रकार", "نوع", "tipo", "type", "Typ", "тип"},
	LocStr{"unknown", "desconocido", "未知", "अज्ञात", "غير معروف", "desconhecido", "inconnu", "unbekannt", "неизвестный"},
	LocStr{"unsigned", "sin signo", "无符号", "अहस्ताक्षरित", "غير موقع", "sem sinal", "non signé", "vorzeichenlos", "безzнаковый"},
	LocStr{"value", "valor", "值", "मूल्य", "قيمة", "valor", "valeur", "Wert", "значение"},
}
