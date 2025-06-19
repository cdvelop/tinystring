package tinystring

// Language enumeration for supported languages
type lang uint8

const (
	EN lang = iota // 0 - English (default)
	ES             // 1 - Spanish
	PT             // 2 - Portuguese
	FR             // 3 - French
	RU             // 4 - Russian
	DE             // 5 - German
	IT             // 6 - Italian
	HI             // 7 - Hindi
	BN             // 8 - Bengali
	ID             // 9 - Indonesian
	AR             // 10 - Arabic
	UR             // 11 - Urdu
	ZH             // 12 - Chinese
)

// OL represents Output Language translations using fixed array for efficiency
type OL [13]string

// get returns translation for specified language with English fallback
func (o OL) get(l lang) string {
	if int(l) < len(o) && o[l] != "" {
		return o[l]
	}
	return o[EN] // Fallback to English
}

// Dictionary structure containing all translatable terms
type dictionary struct {
	// Basic words sorted alphabetically for maximum reusability
	Argument         OL // "argument"
	At               OL // "at"
	Base             OL // "base"
	Boolean          OL // "boolean"
	Cannot           OL // "cannot"
	Empty            OL // "empty"
	End              OL // "end"
	Float            OL // "float"
	For              OL // "for"
	Fmt              OL // "format"
	Found            OL // "found"
	In               OL // "in"
	Integer          OL // "integer"
	Invalid          OL // "invalid"
	Missing          OL // "missing"
	Negative         OL // "negative"
	NegativeUnsigned OL // "negative numbers are not supported for unsigned integers"
	NonNumeric       OL // "non-numeric"
	Not              OL // "not"
	Number           OL // "number"
	Numbers          OL // "numbers"
	Of               OL // "of"
	One              OL // "one"
	Only             OL // "only"
	Overflow         OL // "overflow"
	Range            OL // "range"
	Required         OL // "required"
	Round            OL // "round"
	Specifier        OL // "specifier"
	String           OL // "string"
	Supported        OL // "supported"
	Text             OL // "text"
	Type             OL // "type"
	Unknown          OL // "unknown"
	Unsigned         OL // "unsigned"
	Unsupported      OL // "unsupported"
	Value            OL // "value"
	Wrong            OL // "wrong"
}

// Global dictionary instance - populated with all translations using horizontal format
var D = dictionary{
	// Language order: EN, ES, PT, FR, RU, DE, IT, HI, BN, ID, AR, UR, ZH
	Argument:         OL{"argument", "argumento", "argumento", "argument", "аргумент", "Argument", "argomento", "तर्क", "যুক্তি", "argumen", "وسيط", "دلیل", "参数"},
	At:               OL{"at", "en", "em", "à", "в", "bei", "a", "पर", "এ", "di", "في", "میں", "在"},
	Base:             OL{"base", "base", "base", "base", "основание", "Basis", "base", "आधार", "ভিত্তি", "basis", "قاعدة", "بنیاد", "进制"},
	Boolean:          OL{"boolean", "booleano", "booleano", "booléen", "логический", "boolescher", "booleano", "बूलियन", "বুলিয়ান", "boolean", "منطقي", "بولین", "布尔"},
	Cannot:           OL{"cannot", "no puede", "não pode", "ne peut pas", "не может", "kann nicht", "non può", "नहीं कर सकते", "পারে না", "tidak bisa", "لا يمكن", "نہیں کر سکتے", "不能"},
	Empty:            OL{"empty", "vacío", "vazio", "vide", "пустой", "leer", "vuoto", "खाली", "খালি", "kosong", "فارغ", "خالی", "空"},
	End:              OL{"end", "fin", "fim", "fin", "конец", "Ende", "fine", "अंत", "শেষ", "akhir", "نهاية", "اختتام", "结束"},
	Float:            OL{"float", "flotante", "flutuante", "flottant", "число с плавающей точкой", "Gleitkomma", "virgola mobile", "फ्लोट", "ফ্লোট", "float", "عائم", "فلوٹ", "浮点"},
	For:              OL{"for", "para", "para", "pour", "для", "für", "per", "के लिए", "জন্য", "untuk", "لـ", "کے لیے", "为"},
	Fmt:              OL{"format", "formato", "formato", "format", "формат", "Fmt", "formato", "प्रारूप", "বিন্যাস", "format", "تنسيق", "فارمیٹ", "格式"},
	Found:            OL{"found", "encontrado", "encontrado", "trouvé", "найден", "gefunden", "trovato", "मिला", "পাওয়া", "ditemukan", "موجود", "ملا", "找到"},
	Integer:          OL{"integer", "entero", "inteiro", "entier", "целое число", "ganze Zahl", "intero", "पूर्णांक", "পূর্ণসংখ্যা", "integer", "عدد صحيح", "انٹیجر", "整数"},
	Invalid:          OL{"invalid", "inválido", "inválido", "invalide", "недопустимый", "ungültig", "non valido", "अमान्य", "অবৈধ", "tidak valid", "غير صالح", "غیر درست", "无效"},
	Missing:          OL{"missing", "falta", "ausente", "manquant", "отсутствует", "fehlend", "mancante", "गुम", "অনুপস্থিত", "hilang", "مفقود", "غائب", "缺少"},
	Negative:         OL{"negative", "negativo", "negativo", "négatif", "отрицательный", "negativ", "negativo", "नकारात्मक", "নেগেটিভ", "negatif", "سالب", "منفی", "负"},
	NegativeUnsigned: OL{"negative numbers are not supported for unsigned integers", "números negativos no soportados para enteros sin signo", "números negativos não suportados para inteiros sem sinal", "nombres négatifs non pris en charge pour les entiers non signés", "отрицательные числа не поддерживаются для беззнаковых целых чисел", "negative Zahlen werden für vorzeichenlose Ganzzahlen nicht unterstützt", "numeri negativi non supportati per interi senza segno", "नकारात्मक संख्याएं अहस्ताक्षरित पूर्णांकों के लिए समर्थित नहीं हैं", "স্বাক্ষরহীন পূর্ণসংখ্যার জন্য নেগেটিভ সংখ্যা সমর্থিত নয়", "angka negatif tidak didukung untuk integer tanpa tanda", "الأرقام السالبة غير مدعومة للأعداد الصحيحة غير الموقعة", "منفی نمبرز غیر دستخط انٹیجرز کے لیے معاون نہیں", "无符号整数不支持负数"},
	NonNumeric:       OL{"non-numeric", "no numérico", "não numérico", "non numérique", "нечисловой", "nicht numerisch", "non numerico", "गैर-संख्यात्मक", "অ-সংখ্যাসূচক", "non-numerik", "غير رقمي", "غیر عددی", "非数字"},
	Not:              OL{"not", "no", "não", "pas", "не", "nicht", "non", "नहीं", "না", "tidak", "ليس", "نہیں", "不"},
	Number:           OL{"number", "número", "número", "nombre", "число", "Zahl", "numero", "संख्या", "সংখ্যা", "angka", "رقم", "نمبر", "数字"},
	Numbers:          OL{"numbers", "números", "números", "nombres", "числа", "Zahlen", "numeri", "संख्याएं", "সংখ্যা", "angka", "أرقام", "نمبرز", "数字"},
	Of:               OL{"of", "de", "de", "de", "из", "von", "di", "का", "এর", "dari", "من", "کا", "的"},
	One:              OL{"one", "uno", "um", "un", "один", "eins", "uno", "एक", "একটি", "satu", "واحد", "ایک", "一"},
	Only:             OL{"only", "solo", "apenas", "seulement", "только", "nur", "solo", "केवल", "শুধুমাত্র", "hanya", "فقط", "صرف", "仅"},
	Overflow:         OL{"overflow", "desbordamiento", "estouro", "débordement", "переполнение", "Überlauf", "overflow", "ओवरफ्लो", "ওভারফ্লো", "overflow", "فيض", "اوور فلو", "溢出"},
	Range:            OL{"range", "rango", "intervalo", "plage", "диапазон", "Bereich", "intervallo", "रेंज", "পরিসর", "rentang", "نطاق", "رینج", "范围"},
	Required:         OL{"required", "requerido", "necessário", "requis", "обязательный", "erforderlich", "richiesto", "आवश्यक", "প্রয়োজনীয়", "diperlukan", "مطلوب", "ضروری", "必需"},
	Round:            OL{"round", "redondear", "arredondar", "arrondir", "округлить", "runden", "arrotondare", "गोल", "গোল", "bulatkan", "جولة", "گول", "圆"},
	Specifier:        OL{"specifier", "especificador", "especificador", "spécificateur", "спецификатор", "Spezifizierer", "specificatore", "निर्दिष्टकर्ता", "নির্দিষ্টকারী", "penentu", "محدد", "تعین کنندہ", "说明符"},
	String:           OL{"string", "cadena", "string", "chaîne", "строка", "Zeichenkette", "stringa", "स्ट्रिंग", "স্ট্রিং", "string", "سلسلة", "سٹرنگ", "字符串"},
	Supported:        OL{"supported", "soportado", "suportado", "pris en charge", "поддерживается", "unterstützt", "supportato", "समर्थित", "সমর্থিত", "didukung", "مدعوم", "معاون", "支持"},
	Text:             OL{"text", "texto", "texto", "texte", "текст", "Text", "testo", "पाठ", "পাঠ", "teks", "نص", "متن", "文本"},
	Type:             OL{"type", "tipo", "tipo", "type", "тип", "Typ", "tipo", "प्रकार", "টাইপ", "tipe", "نوع", "قسم", "类型"},
	Unknown:          OL{"unknown", "desconocido", "desconhecido", "inconnu", "неизвестный", "unbekannt", "sconosciuto", "अज्ञात", "অজানা", "tidak diketahui", "غير معروف", "نامعلوم", "未知"},
	Unsigned:         OL{"unsigned", "sin signo", "sem sinal", "non signé", "беззнаковый", "vorzeichenlos", "senza segno", "अहस्ताक्षरित", "স্বাক্ষরহীন", "tidak bertanda", "غير موقع", "غیر دستخط شدہ", "无符号"},
	Unsupported:      OL{"unsupported", "no soportado", "não suportado", "non pris en charge", "не поддерживается", "nicht unterstützt", "non supportato", "असमर्थित", "অসমর্থিত", "tidak didukung", "غير مدعوم", "غیر معاون", "不支持"},
	Value:            OL{"value", "valor", "valor", "valeur", "значение", "Wert", "valore", "मूल्य", "মান", "nilai", "قيمة", "قیمت", "值"},
	Wrong:            OL{"wrong", "incorrecto", "errado", "mauvais", "неправильный", "falsch", "sbagliato", "गलत", "ভুল", "salah", "خطأ", "غلط", "错误"},
}

// Private global configuration
var defLang lang = EN

// OutLang sets the default output language
// OutLang() without parameters auto-detects system language
// OutLang(ES) sets Spanish as default
func OutLang(l ...lang) {
	if len(l) == 0 {
		defLang = getSystemLang()
	} else {
		defLang = l[0]
	}
}
