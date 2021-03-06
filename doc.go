// 2015-04-23 Adam Bryt

/*
Program resimg zmienia rozmiar obrazków. Każdy obrazek (argument
polecenia) jest konwertowany do rozmiaru określonego opcją -size
(domyślnie 'medium' czyli '640x480') i kopiowany do katalogu
docelowego (domyślnie '/tmp/small'). Katalog docelowy jest tworzony
jeśli nie istnieje. Obsługuje formaty GIF, JPEG i PNG.  Program
może służyć do zmniejszania zdjęć przed wysłaniem emailem.

Sposób użycia:
	resimg [opcje] file ...

Opcje:
	-dir="/tmp/small"
		Katalog docelowy
	-size="medium"
		Rozmiar obrazków: 'small', 'medium', 'large', lub w formacie
		typu '300x200'.
		Zdefiniowane nazwy rozmiarów: small=320x240, medium=640x480,
		large=800x600.
		Można używać skróconych nazw (np. -size m).
	-v=false
		Verbose: informuje co robi
	-help=false
		Wyświetla help
*/
package main
