package tinystring

// Global dictionary instance - populated with all translations using horizontal format
// Language order: EN, ES, ZH, HI, AR, PT, FR, DE, RU
//
// By using an anonymous struct, we define and initialize the dictionary in a single step,
// avoiding the need for a separate 'type dictionary struct' declaration.
// The usage API (e.g., D.Argument) remains unchanged.
var D = struct {
	// A
	All        LocStr // "all"
	Allowed    LocStr // "allowed"
	Arrow      LocStr // "arrow"
	Argument   LocStr // "argument"
	Assign     LocStr // "assign"
	Assignable LocStr // "assignable"

	// B
	BackingUp LocStr // "backing up"
	Be        LocStr // "be"
	Begin     LocStr // "begin"
	Binary    LocStr // "binary"

	// C
	Call          LocStr // "call"
	Can           LocStr // "can"
	Cannot        LocStr // "cannot"
	Cancel        LocStr // "cancel"
	Changed       LocStr // "changed"
	Character     LocStr // "character"
	Chars         LocStr // "chars"
	Checker       LocStr // "checker"
	Coding        LocStr // "coding"
	Compilation   LocStr // "compilation"
	Configuration LocStr // "configuration"
	Connection    LocStr // "connection"
	Content       LocStr // "content"
	Create        LocStr // "create"

	// D
	Date       LocStr // "date"
	Debugging  LocStr // "debugging"
	Decimal    LocStr // "decimal"
	Delimiter  LocStr // "delimiter"
	Dictionary LocStr // "dictionary"
	Digit      LocStr // "digit"
	Down       LocStr // "down"

	// E
	Edit    LocStr // "edit"
	Element LocStr // "element"
	Email   LocStr // "email"
	Empty   LocStr // "empty"
	End     LocStr // "end"
	Example LocStr // "example"
	Exceeds LocStr // "exceeds"
	Execute LocStr // "execute"

	// F
	Failed LocStr // "failed"
	Female LocStr // "female
	Field  LocStr // "field"
	Fields LocStr // "fields"
	Files  LocStr // "files"
	Format LocStr // "format"
	Found  LocStr // "found"

	// H
	Handler LocStr // "handler"
	Hyphen  LocStr // "hyphen"
	Hour    LocStr // "hour"

	// I
	Icons        LocStr // "icons"
	Insert       LocStr // "insert"
	Left         LocStr // "left"
	Implemented  LocStr // "implemented"
	In           LocStr // "in"
	Index        LocStr // "index"
	Information  LocStr // "information"
	Input        LocStr // "input"
	Install      LocStr // "install"
	Installation LocStr // "installation"
	Invalid      LocStr // "invalid"

	// K
	Keyboard LocStr // "keyboard"

	// L
	Language LocStr // "language"
	Letters  LocStr // "letters"
	Line     LocStr // "line"

	// M
	Male     LocStr // "male"
	Maximum  LocStr // "maximum"
	Method   LocStr // "method"
	Missing  LocStr // "missing"
	Mismatch LocStr // "mismatch"
	Mode     LocStr // "mode"
	Modes    LocStr // "modes"
	More     LocStr // "more"
	Move     LocStr // "move"
	Must     LocStr // "must"

	// N
	Negative   LocStr // "negative"
	Nil        LocStr // "null"
	NonNumeric LocStr // "non-numeric"
	Not        LocStr // "not"
	NotOfType  LocStr // "not of type"
	Number     LocStr // "number"
	Numbers    LocStr // "numbers"

	// O
	Of       LocStr // "of"
	Options  LocStr // "options"
	Out      LocStr // "out"
	Overflow LocStr // "overflow"

	// P
	Page       LocStr // "page"
	Pointer    LocStr // "pointer"
	Point      LocStr // "point"
	Preparing  LocStr // "preparing"
	Production LocStr // "production"
	Provided   LocStr // "provided"

	// Q
	Quit LocStr // "quit"

	// R
	Range    LocStr // "range"
	Read     LocStr // "read"
	Required LocStr // "required"
	Right    LocStr // "right"
	Round    LocStr // "round"

	// S
	Seconds   LocStr // "seconds"
	Session   LocStr // "session"
	Slice     LocStr // "slice"
	Space     LocStr // "space"
	Status    LocStr // "status"
	String    LocStr // "string"
	Shortcuts LocStr // "shortcuts"
	Supported LocStr // "supported"
	Switch    LocStr // "switch"
	Switching LocStr // "switching"
	Sync      LocStr // "sync"
	System    LocStr // "system"

	// Translate
	Tab     LocStr // "tab"
	Test    LocStr // "test"
	Testing LocStr // "testing"
	Text    LocStr // "text"
	Time    LocStr // "time"
	To      LocStr // "to"
	Type    LocStr // "type"

	// U
	Unexported LocStr // "unexported"
	Unknown    LocStr // "unknown"
	Unsigned   LocStr // "unsigned"
	Up         LocStr // "up"
	Use        LocStr // "use"

	// V
	Valid      LocStr // "valid"
	Validating LocStr // "validating"
	Value      LocStr // "value"
	Visible    LocStr // "visible"

	// W
	Warning LocStr // "warning"
	With    LocStr // "with"

	// Z
	Zero LocStr // "zero"
}{
	// A
	LocStr{"All", "Todo", "所有", "सभी", "كل", "Todo", "Tout", "Alle", "Все"},
	LocStr{"Allowed", "Permitido", "允许", "अनुमति", "مسموح", "Permitido", "Autorisé", "Erlaubt", "Разрешено"},
	LocStr{"Arrow", "Flecha", "箭头", "तीर", "سهم", "Seta", "Flèche", "Pfeil", "Стрелка"},
	LocStr{"Argument", "Argumento", "参数", "तर्क", "وسيط", "Argumento", "Argument", "Argument", "Аргумент"},
	LocStr{"Assign", "Asignar", "分配", "असाइन", "تعيين", "Atribuir", "Assigner", "Zuweisen", "Присвоить"},
	LocStr{"Assignable", "Asignable", "可分配", "असाइन करने योग्य", "قابل للتعيين", "Atributável", "Assignable", "Zuweisbar", "Присваиваемый"},

	// B
	LocStr{"Backing up", "Respaldando", "备份", "बैकअप", "نسخ احتياطي", "Fazendo backup", "Sauvegarde", "Sicherung", "Резервное копирование"},
	LocStr{"be", "ser", "是", "होना", "كون", "ser", "être", "sein", "быть"},
	LocStr{"Begin", "Comenzar", "开始", "शुरू", "ابدأ", "Começar", "Commencer", "Beginnen", "Начать"},
	LocStr{"Binary", "Binario", "二进制", "二进制", "ثنائي", "Binário", "Binaire", "Binär", "Двоичный"},

	// C
	LocStr{"Call", "Llamar", "调用", "कॉल", "استدعاء", "Chamar", "Appeler", "Aufrufen", "Вызвать"},
	LocStr{"Can", "Puede", "可以", "कर सकते हैं", "يمكن", "Pode", "Peut", "Kann", "Может"},
	LocStr{"Cannot", "No puede", "不能", "नहीं कर सकते", "لا يمكن", "Não pode", "Ne peut pas", "Kann nicht", "Не может"},
	LocStr{"Cancel", "Cancelar", "取消", "रद्द करें", "إلغاء", "Cancelar", "Annuler", "Abbrechen", "Отменить"},
	LocStr{"Changed", "Cambiado", "更改", "परिवर्तित", "تم التغيير", "Alterado", "Changé", "Geändert", "Изменено"},
	LocStr{"Character", "Caracter", "字符", "वर्ण", "حرف", "Caractere", "Caractère", "Zeichen", "Символ"},
	LocStr{"Chars", "Caracteres", "字符", "अक्षर", "حروف", "Caracteres", "Caractères", "Zeichen", "Символы"},
	LocStr{"Checker", "Verificador", "检查器", "चेककर", "مدقق", "Verificador", "Vérificateur", "Prüfer", "Проверяющий"},
	LocStr{"Coding", "Codificación", "编码", "कोडिंग", "ترميز", "Codificação", "Codage", "Kodierung", "Кодирование"},
	LocStr{"Compilation", "Compilación", "编译", "संकलन", "تجميع", "Compilação", "Compilation", "Kompilierung", "Компиляция"},
	LocStr{"Configuration", "Configuración", "配置", "कॉन्फ़िगरेशन", "تكوين", "Configuração", "Configuration", "Konfiguration", "Конфигурация"},
	LocStr{"Connection", "Conexión", "连接", "संपर्क", "اتصال", "Conexão", "Connexion", "Verbindung", "Соединение"},
	LocStr{"Content", "Contenido", "内容", "सामग्री", "محتوى", "Conteúdo", "Contenu", "Inhalt", "Содержимое"},
	LocStr{"Create", "Crear", "创建", "बनाएँ", "إنشاء", "Criar", "Créer", "Erstellen", "Создать"},

	// D
	LocStr{"Date", "Fecha", "日期", "तारीख", "تاريخ", "Data", "Date", "Datum", "Дата"},
	LocStr{"Debugging", "Depuración", "调试", "डिबगिंग", "تصحيح", "Depuração", "Débogage", "Debuggen", "Отладка"},
	LocStr{"Decimal", "Decimal", "十进制", "दशमलव", "عشري", "Decimal", "Décimal", "Dezimal", "Десятичная"},
	LocStr{"Delimiter", "Delimitador", "分隔符", "सीमांकक", "محدد", "Delimitador", "Délimiteur", "Trennzeichen", "Разделитель"},
	LocStr{"Dictionary", "Diccionario", "字典", "शब्दकोश", "قاموس", "Dicionário", "Dictionnaire", "Wörterbuch", "Словарь"},
	LocStr{"Digit", "Dígito", "数字", "अंक", "رقم", "Dígito", "Chiffre", "Ziffer", "Цифра"},
	LocStr{"Down", "Abajo", "下", "नीचे", "أسفل", "Baixo", "Bas", "Unten", "Вниз"},

	// E
	LocStr{"Edit", "Editar", "编辑", "संपादित करें", "تحرير", "Editar", "Éditer", "Bearbeiten", "Редактировать"},
	LocStr{"Element", "Elemento", "元素", "एलिमेंट", "عنصر", "Elemento", "Élément", "Element", "Элемент"},
	LocStr{"Email", "Correo electrónico", "电子邮件", "ईमेल", "البريد الإلكتروني", "Email", "Email", "E-Mail", "Электронная почта"},
	LocStr{"Empty", "Vacío", "空", "खाली", "فارغ", "Vazio", "Vide", "Leer", "Пустой"},
	LocStr{"End", "Fin", "结束", "अंत", "نهاية", "Fim", "Fin", "Ende", "Конец"},
	LocStr{"Example", "Ejemplo", "例子", "उदाहरण", "مثال", "Exemplo", "Exemple", "Beispiel", "Пример"},
	LocStr{"Exceeds", "Excede", "超过", "अधिक", "يتجاوز", "Excede", "Dépasse", "Überschreitet", "Превышает"},
	LocStr{"Execute", "Ejecutar", "执行", "निष्पादित करें", "تنفيذ", "Executar", "Exécuter", "Ausführen", "Выполнить"},

	// F
	LocStr{"Failed", "Falló", "失败", "असफल", "فشل", "Falhou", "Échoué", "Fehlgeschlagen", "Не удалось"},
	LocStr{"Female", "Femenino", "女性", "महिला", "أنثى", "Feminino", "Féminin", "Weiblich", "Женский"},
	LocStr{"Field", "Campo", "字段", "फील्ड", "حقل", "Campo", "Champ", "Feld", "Поле"},
	LocStr{"Fields", "Campos", "字段集", "फील्ड्स", "حقول", "Campos", "Champs", "Felder", "Поля"},
	LocStr{"Files", "Archivos", "文件", "फ़ाइलें", "ملفات", "Arquivos", "Fichiers", "Dateien", "Файлы"},
	LocStr{"Format", "Formato", "格式", "प्रारूप", "تنسيق", "Formato", "Format", "Fmt", "Формат"},
	LocStr{"Found", "Encontrado", "找到", "मिला", "موجود", "Encontrado", "Trouvé", "Gefunden", "Найден"},

	// H
	LocStr{"Handler", "Manejador", "处理程序", "हैंडलर", "معالج", "Manipulador", "Gestionnaire", "Handler", "Обработчик"},
	LocStr{"Hyphen", "Guion", "连字符", "हाइफ़न", "شرطة", "Hífen", "Trait d'union", "Bindestrich", "Дефис"},
	LocStr{"Hour", "Hora", "小时", "घंटा", "ساعة", "Hora", "Heure", "Stunde", "Час"},

	// I
	LocStr{"Icons", "Iconos", "图标", "आइकन", "أيقونات", "Ícones", "Icônes", "Symbole", "Значки"},
	LocStr{"Insert", "Insertar", "插入", "सम्मिलित करें", "إدراج", "Inserir", "Insérer", "Einfügen", "Вставить"},
	LocStr{"Left", "Izquierda", "左", "बाएं", "يسار", "Esquerda", "Gauche", "Links", "Слева"},
	LocStr{"Implemented", "Implementado", "已实现", "कार्यान्वित", "مُنفذ", "Implementado", "Implémenté", "Implementiert", "Реализовано"},
	LocStr{"In", "En", "在", "में", "في", "Em", "Dans", "In", "В"},
	LocStr{"Index", "Índice", "索引", "इंडेक्स", "فهرس", "Índice", "Index", "Index", "Индекс"},
	LocStr{"Information", "Información", "信息", "सूचना", "معلومات", "Informação", "Information", "Information", "Информация"},
	LocStr{"Input", "Entrada", "输入", "इनपुट", "إدخال", "Entrada", "Entrée", "Eingabe", "Ввод"},
	LocStr{"Install", "Instalar", "安装", "इंस्टॉल", "تثبيت", "Instalar", "Installer", "Installieren", "Установить"},
	LocStr{"Installation", "Instalación", "安装", "स्थापना", "التثبيت", "Instalação", "Installation", "Installation", "Установка"},
	LocStr{"Invalid", "Inválido", "无效", "अमान्य", "غير صالح", "Inválido", "Invalide", "Ungültig", "Недопустимый"},

	// K
	LocStr{"Keyboard", "Teclado", "键盘", "कीबोर्ड", "لوحة المفاتيح", "Teclado", "Clavier", "Tastatur", "Клавиатура"},

	// L
	LocStr{"Language", "Idioma", "语言", "भाषा", "لغة", "Idioma", "Langue", "Sprache", "Язык"},
	LocStr{"Letters", "Letras", "字母", "अक्षर", "حروف", "Letras", "Lettres", "Buchstaben", "Буквы"},
	LocStr{"Line", "Línea", "行", "लाइन", "خط", "Linha", "Ligne", "Zeile", "Строка"},

	// M
	LocStr{"Male", "Masculino", "男性", "पुरुष", "ذكر", "Masculino", "Masculin", "Männlich", "Мужской"},
	LocStr{"Maximum", "Máximo", "最大", "अधिकतम", "الحد الأقصى", "Máximo", "Maximum", "Maximum", "Максимум"},
	LocStr{"Method", "Método", "方法", "विधि", "طريقة", "Método", "Méthode", "Methode", "Метод"},
	LocStr{"Missing", "Falta", "缺少", "गुम", "مفقود", "Ausente", "Manquant", "Fehlend", "Отсутствует"},
	LocStr{"Mismatch", "Desajuste", "不匹配", "बेमेल", "عدم تطابق", "Incompatibilidade", "Incompatibilité", "Nichtübereinstimmung", "Несоответствие"},
	LocStr{"Mode", "Modo", "模式", "मोड", "وضع", "Modo", "Mode", "Modus", "Режим"},
	LocStr{"Modes", "Modos", "模式", "मोड", "أوضاع", "Modos", "Modes", "Modi", "Режимы"},
	LocStr{"More", "Más", "更多", "अधिक", "أكثر", "Mais", "Plus", "Mehr", "Больше"},
	LocStr{"Move", "Mover", "移动", "स्थानांतरित करें", "نقل", "Mover", "Déplacer", "Verschieben", "Переместить"},
	LocStr{"Must", "Debe", "必须", "चाहिए", "يجب", "Deve", "Doit", "Muss", "Должен"},

	// N
	LocStr{"Negative", "Negativo", "负", "नकारात्मक", "سالب", "Negativo", "Négatif", "Negativ", "Отрицательный"},
	LocStr{"Nil", "Nulo", "空值", "शून्य मान", "قيمة فارغة", "Nulo", "Nul", "Nullwert", "Нуль"},
	LocStr{"Non-numeric", "No numérico", "非数字", "गैर-संख्यात्मक", "غير رقمي", "Não numérico", "Non numérique", "Nicht numerisch", "Нечисловой"},
	LocStr{"Not", "No", "不", "नहीं", "ليس", "Não", "Pas", "Nicht", "Не"},
	LocStr{"Not of type", "No es del tipo", "不是类型", "प्रकार नहीं है", "ليس من النوع", "Não é do tipo", "N'est pas du type", "Ist nicht vom Typ", "Не того типа"},
	LocStr{"Number", "Número", "数字", "संख्या", "رقم", "Número", "Nombre", "Zahl", "Число"},
	LocStr{"Numbers", "Números", "数字", "संख्याएं", "أرقام", "Números", "Nombres", "Zahlen", "Числа"},

	// O
	LocStr{"of", "de", "的", "का", "من", "de", "de", "von", "из"},
	LocStr{"Options", "Opciones", "选项", "विकल्प", "خيارات", "Opções", "Options", "Optionen", "Опции"},
	LocStr{"Out", "Fuera", "出", "बाहर", "خارج", "Fora", "Hors", "Aus", "Вне"},
	LocStr{"Overflow", "Desbordamiento", "溢出", "ओवरफ्लो", "فيض", "Estouro", "Débordement", "Überlauf", "Переполнение"},

	// P
	LocStr{"Page", "Página", "页面", "पृष्ठ", "صفحة", "Página", "Page", "Seite", "Страница"},
	LocStr{"Pointer", "Puntero", "指针", "पॉइंटर", "مؤشر", "Ponteiro", "Pointeur", "Zeiger", "Указатель"},
	LocStr{"Point", "Punto", "点", "बिंदु", "نقطة", "Ponto", "Point", "Punkt", "Точка"},
	LocStr{"Preparing", "Preparando", "准备", "तैयारी", "تحضير", "Preparando", "Préparation", "Vorbereitung", "Подготовка"},
	LocStr{"Production", "Producción", "生产", "उत्पादन", "إنتاج", "Produção", "Production", "Produktion", "Производство"},
	LocStr{"Provided", "Proporcionado", "提供", "प्रदान किया गया", "مقدم", "Fornecido", "Fourni", "Bereitgestellt", "Предоставлено"},

	// Q
	LocStr{"Quit", "Salir", "退出", "छोड़ें", "إنهاء", "Sair", "Quitter", "Beenden", "Выйти"},

	// R
	LocStr{"Range", "Rango", "范围", "रेंज", "نطاق", "Intervalo", "Plage", "Bereich", "Диапазон"},
	LocStr{"Read", "Leer", "读取", "पढ़ें", "قراءة", "Ler", "Lire", "Lesen", "Читать"},
	LocStr{"Required", "Requerido", "必需", "आवश्यक", "مطلوب", "Necessário", "Requis", "Erforderlich", "Обязательный"},
	LocStr{"Right", "Derecha", "右", "दाएं", "يمين", "Direita", "Droite", "Rechts", "Справа"},
	LocStr{"Round", "Redondear", "圆", "गोल", "جولة", "Arredondar", "Arrondir", "Runden", "Округлить"},

	// S
	LocStr{"Seconds", "Segundos", "秒", "सेकंड", "ثواني", "Segundos", "Secondes", "Sekunden", "Секунды"},
	LocStr{"Session", "Sesión", "会话", "सत्र", "جلسة", "Sessão", "Session", "Sitzung", "Сессия"},
	LocStr{"Slice", "Segmento", "切片", "स्लाइस", "شريحة", "Fatia", "Tranche", "Scheibe", "Срез"},
	LocStr{"Space", "Espacio", "空间", "जगह", "فضاء", "Espaço", "Espace", "Raum", "Пространство"},
	LocStr{"Status", "Estado", "状态", "स्थिति", "حالة", "Status", "Statut", "Status", "Статус"},
	LocStr{"String", "Cadena", "字符串", "स्ट्रिंग", "سلسلة", "String", "Chaîne", "Zeichenkette", "Строка"},
	LocStr{"Shortcuts", "Atajos", "快捷键", "शॉर्टकट्स", "اختصارات", "Atalhos", "Raccourcis", "Kurzbefehle", "Ярлыки"},
	LocStr{"Supported", "Soportado", "支持", "समर्थित", "مدعوم", "Suportado", "Pris en charge", "Unterstützt", "Поддерживается"},
	LocStr{"Switch", "Cambiar", "切换", "स्विच", "تبديل", "Mudar", "Changer", "Wechseln", "Переключить"},
	LocStr{"Switching", "Cambiando", "切换中", "स्विच कर रहा है", "تبديل", "Mudando", "Changement", "Wechseln", "Переключение"},
	LocStr{"Sync", "Sincronización", "同步", "सिंक", "مزامنة", "Sincronização", "Synchronisation", "Synchronisierung", "Синхронизация"},
	LocStr{"System", "Sistema", "系统", "सिस्टम", "نظام", "Sistema", "Système", "System", "Система"},

	// T
	LocStr{"Tab", "Pestaña", "标签页", "टैब", "علامة تبويب", "Aba", "Onglet", "Registerkarte", "Вкладка"},
	LocStr{"Test", "Prueba", "测试", "परीक्षण", "اختبار", "Teste", "Test", "Test", "Тест"},
	LocStr{"Testing", "Probando", "测试中", "परीक्षण", "اختبارات", "Testando", "Test", "Testen", "Тестирование"},
	LocStr{"Text", "Texto", "文本", "पाठ", "نص", "Texto", "Texte", "Text", "Текст"},
	LocStr{"Time", "Tiempo", "时间", "समय", "وقت", "Tempo", "Temps", "Zeit", "Время"},
	LocStr{"to", "a", "到", "को", "إلى", "para", "à", "zu", "к"},
	LocStr{"Type", "Tipo", "类型", "प्रकार", "نوع", "Tipo", "Type", "Typ", "Тип"},

	// U
	LocStr{"Unexported", "No Exportado", "未导出", "गैर-निर्यातित", "غير مصدر", "Não Exportado", "Non Exporté", "Nicht Exportiert", "Неэкспортированный"},
	LocStr{"Unknown", "Desconocido", "未知", "अज्ञात", "غير معروف", "Desconhecido", "Inconnu", "Unbekannt", "Неизвестный"},
	LocStr{"Unsigned", "Sin Signo", "无符号", "अहस्ताक्षरित", "غير موقع", "Sem Sinal", "Non Signé", "Vorzeichenlos", "Беззнаковый"},
	LocStr{"Up", "Arriba", "上", "ऊपर", "أعلى", "Cima", "Haut", "Oben", "Вверх"},
	LocStr{"Use", "Usar", "使用", "उपयोग", "استخدام", "Usar", "Utiliser", "Verwenden", "Использовать"},

	// V
	LocStr{"Valid", "Válido", "有效", "वैध", "صحيح", "Válido", "Valide", "Gültig", "Действительный"},
	LocStr{"Validating", "Validando", "验证中", "सत्यापन हो रहा है", "التحقق من", "Validando", "Validation", "Validieren", "Проверка"},
	LocStr{"Value", "Valor", "值", "मूल्य", "قيمة", "Valor", "Valeur", "Wert", "Значение"},
	LocStr{"Visible", "Visible", "可见", "दृश्य", "مرئي", "Visível", "Visible", "Sichtbar", "Видимый"},

	// W
	LocStr{"Warning", "Advertencia", "警告", "चेतावनी", "تحذير", "Aviso", "Avertissement", "Warnung", "Предупреждение"},
	LocStr{"With", "Con", "与", "के साथ", "مع", "Com", "Avec", "Mit", "С"},

	// Z
	LocStr{"Zero", "Cero", "零", "शून्य", "صفر", "Zero", "Zéro", "Null", "Ноль"},
}
