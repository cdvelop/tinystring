package tinystring

// Slice of mappings to remove accents and diacritics
var accentMappings = []charMapping{
	{'á', 'a'}, {'à', 'a'}, {'ã', 'a'}, {'â', 'a'}, {'ä', 'a'},
	{'é', 'e'}, {'è', 'e'}, {'ê', 'e'}, {'ë', 'e'},
	{'í', 'i'}, {'ì', 'i'}, {'î', 'i'}, {'ï', 'i'},
	{'ó', 'o'}, {'ò', 'o'}, {'õ', 'o'}, {'ô', 'o'}, {'ö', 'o'},
	{'ú', 'u'}, {'ù', 'u'}, {'û', 'u'}, {'ü', 'u'},
	{'ý', 'y'}, {'ÿ', 'y'},
	{'ñ', 'n'},
	// Upper case
	{'Á', 'A'}, {'À', 'A'}, {'Ã', 'A'}, {'Â', 'A'}, {'Ä', 'A'},
	{'É', 'E'}, {'È', 'E'}, {'Ê', 'E'}, {'Ë', 'E'},
	{'Í', 'I'}, {'Ì', 'I'}, {'Î', 'I'}, {'Ï', 'I'},
	{'Ó', 'O'}, {'Ò', 'O'}, {'Õ', 'O'}, {'Ô', 'O'}, {'Ö', 'O'},
	{'Ú', 'U'}, {'Ù', 'U'}, {'Û', 'U'}, {'Ü', 'U'},
	{'Ý', 'Y'},
	{'Ñ', 'N'},
}

// Mappings to convert upper case to lower case
var lowerMappings = []charMapping{
	{'A', 'a'}, {'B', 'b'}, {'C', 'c'}, {'D', 'd'}, {'E', 'e'},
	{'F', 'f'}, {'G', 'g'}, {'H', 'h'}, {'I', 'i'}, {'J', 'j'},
	{'K', 'k'}, {'L', 'l'}, {'M', 'm'}, {'N', 'n'}, {'O', 'o'},
	{'P', 'p'}, {'Q', 'q'}, {'R', 'r'}, {'S', 's'}, {'T', 't'},
	{'U', 'u'}, {'V', 'v'}, {'W', 'w'}, {'X', 'x'}, {'Y', 'y'},
	{'Z', 'z'},
}

// Mappings to convert lower case to upper case
var upperMappings = []charMapping{
	{'a', 'A'}, {'b', 'B'}, {'c', 'C'}, {'d', 'D'}, {'e', 'E'},
	{'f', 'F'}, {'g', 'G'}, {'h', 'H'}, {'i', 'I'}, {'j', 'J'},
	{'k', 'K'}, {'l', 'L'}, {'m', 'M'}, {'n', 'N'}, {'o', 'O'},
	{'p', 'P'}, {'q', 'Q'}, {'r', 'R'}, {'s', 'S'}, {'t', 'T'},
	{'u', 'U'}, {'v', 'V'}, {'w', 'W'}, {'x', 'X'}, {'y', 'Y'},
	{'z', 'Z'},
}
