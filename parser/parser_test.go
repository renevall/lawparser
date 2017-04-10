package parser

import (
	"sync"
	"testing"
)

func TestPrepareTags(t *testing.T) {
	tag1 := Tags{
		Tag{"Name", "LEY DE|CÓDIGO"},
		Tag{"Number", "No\\."},
	}
	tag2 := Tags{
		Tag{"Aproved", "Aprobada"},
		Tag{"Diary", "Publicada"},
		Tag{"Arto", "Art\\.\\s\\d+"},
	}
	tag3 := Tags{
		Tag{"Diary", "Publicada"},
	}

	tests := []struct {
		in  Tags
		out string
	}{
		{tag1, "LEY DE|CÓDIGO"},
		{tag2, "Aprobada"},
		{tag3, "Publicada"},
	}

	for _, tt := range tests {
		actual := prepareTags(tt.in)
		if actual[0].exp.String() != tt.out {
			t.Errorf("Tag(%v): expected %v, actual %d", tt.in, tt.out, actual[0].exp.String())
			break
		}
	}

}

func TestFindBasicData(t *testing.T) {
	// fmt.Println("TestFindBasicDataName started")

	input := make(chan string, 1)
	defer close(input)
	done := make(chan struct{})
	defer close(done)
	wg := new(sync.WaitGroup)

	nameTests := []struct {
		in  string
		out string
	}{
		{"LEY DE CONCERTACIÓN TRIBUTARIA", "LEY DE CONCERTACIÓN TRIBUTARIA"},
		{"CÓDIGO TRIBUTARIO", "CÓDIGO TRIBUTARIO"},
		{"Código Tributario", ""},
		{"Some Text CÓDIGO TRIBUTARIO", ""},
		{"Some Text LEY DE CONCERTACIÓN TRIBUTARIA", ""},
		{"CÓDIGO PROCESAL CIVIL DE LA REPÚBLICA DE NICARAGUA", "CÓDIGO PROCESAL CIVIL DE LA REPÚBLICA DE NICARAGUA"},
		{"CODIGO DEL TRABAJO", "CODIGO DEL TRABAJO"},
	}

	for _, tt := range nameTests {
		law := NewLaw()
		FindBasicData(law, done, input, wg)
		input <- tt.in
		input <- "Art. 1"

		<-done
		if law.Name != tt.out {
			t.Errorf("Name testing %q, expected %q, actual %q", tt.in, tt.out, law.Name)
		}
	}

	numberTests := []struct {
		in  string
		out int
	}{
		{"LEY N°. 902", 902},
		{"LEY N°. 902 25", 902},
		{"LEY No. 822", 822},
		{"LEY No. 185, Aprobada el 5 de Septiembre de 1996", 185},
	}

	for _, tt := range numberTests {
		law := NewLaw()
		FindBasicData(law, done, input, wg)
		input <- tt.in
		input <- "Art. 1"

		<-done
		if law.Number != tt.out {
			t.Errorf("Name testing %q, expected %d, actual %d", tt.in, tt.out, law.Number)
		}
	}

	approvalTests := []struct {
		in  string
		out string
	}{
		{"Aprobada el 30 de Noviembre del 2012", "2012-11-30 00:00:00 +0000 UTC"},
		{"LEY No. 185, Aprobada el 5 de Septiembre de 1996", "1996-09-05 00:00:00 +0000 UTC"},
		{"Aprobado el 8 de Abril de 1988", "1988-04-08 00:00:00 +0000 UTC"},
	}

	for _, tt := range approvalTests {
		law := NewLaw()
		FindBasicData(law, done, input, wg)
		input <- tt.in
		input <- "Art. 1"

		<-done
		if law.ApprovalDate.String() != tt.out {
			t.Errorf("Name testing %q, expected %q, actual %q", tt.in, tt.out, law.ApprovalDate)
		}
	}

	diaryTests := []struct {
		in  string
		out string
	}{
		{"Publicada en La Gaceta No. 241 del 17 de Diciembre del 2012", "241"},
		{"Publicada en La Gaceta No. 205 del 30 de Octubre de 1996", "205"},
		{"Publicado en La Gaceta No. 121 del 27 de Junio del 2000.", "121"},
	}

	for _, tt := range diaryTests {
		law := NewLaw()
		FindBasicData(law, done, input, wg)
		input <- tt.in
		input <- "Art. 1"

		<-done
		if law.Journal != tt.out {
			t.Errorf("Name testing %q, expected %q, actual %d", tt.in, tt.out, law.Journal)
		}
	}

	introData := []string{
		"CODIGO DEL TRABAJO ( CON SUS REFORMAS, ADICIONES E",
		"INTERPRETACIÓN AUTENTICA)",
		"LEY No. 185, Aprobada el 5 de Septiembre de 1996",
		"Publicada en La Gaceta No. 205 del 30 de Octubre de 1996",
		"EL PRESIDENTE DE LA REPUBLICA DE NICARAGUA",
		"Hace saber al Pueblo Nicaragüense que",
		"LA ASAMBLEA NACIONAL DE LA REPUBLICA DE NICARAGUA",
		"En uso de sus facultades",
		"HA DICTADO",
		"El siguiente",
		"CODIGO DEL TRABAJO",
		"LIBRO PRIMERO",
		"DERECHO SUSTANTIVO",
		"TITULO PRELIMINAR",
		"PRINCIPIOS FUNDAMENTALES",
		"I",
		"El trabajo es un derecho, una responsabilidad social y goza de la especial",
		"protección del Estado. El Estado procurará la ocupación plena y productiva de",
		"todos los nicaragüenses.",
		"TITULO",
		"DISPOSICIONES GENERALES",
		"CAPITULO I",
		"OBJETO Y AMBITO DE APLICACION",
		"Artículo 1.- El presente código regula las relaciones de trabajo estableciendo",
	}

	inIntro := make(chan string)
	done2 := make(chan struct{})
	law2 := NewLaw()
	wg2 := new(sync.WaitGroup)
	FindBasicData(law2, done2, inIntro, wg2)
	var l int

	for _, row := range introData {
		l = l + len(row) + 1
		inIntro <- row
	}
	l--
	<-done2
	if len(law2.Intro) != l {
		t.Errorf("Testing law intro, expected %d length, actual %d", l, len(law2.Intro))

	}

}

func TestFindCTags(t *testing.T) {
	in := make(chan string)
	index := make(chan []int)
	defer close(index)
	tags := make(chan foundTag)

	wg := new(sync.WaitGroup)

	keywordTests := []struct {
		in    string
		out   string
		found string ``
	}{
		{"Artículo 22", "Arto", ""},
		{"Arto. 22", "Arto", ""},
		{"Articulo", "", ""},
		{"Articulo 11", "Arto", ""},
		{"Este es no es Artículo", "", ""},
		{"Este es no es Articulo", "", ""},
		{"Este es no es articulo", "", ""},
		{"Artículo 1. Objeto.", "Arto", ""},
		{"Art. 1", "Arto", ""},
		{"Art. 2 Principios tributarios.", "Arto", ""},
		{"Art. 12 Vínculos económicos de las rentas del trabajo de fuente", "Arto", ""},
		{"Arto. 12", "Arto", ""},
		{"Articulo 12", "Arto", ""},
		{"Artículo 12", "Arto", ""},
		{"LIBRO I", "Libro", ""},
		{"LIBRO V", "Libro", ""},
		{"LIBRO IV", "Libro", ""},
		{"LIBRO PRIMERO", "Libro", ""},
		{"LIBRO SEGUNDO", "Libro", ""},
		{"LIBRO OCTAVO", "Libro", ""},
		{"TÍTULO I", "Titulo", ""},
		{"TITULO II", "Titulo", ""},
		{"TÍTULO IV", "Titulo", ""},
		{"TÍTULO ÚNICO", "Titulo", ""},
		{"TÍTULO PRELIMINAR", "Titulo", ""},
		{"Este es no es título", "", ""},
		{"titulo II", "", ""},
		{"Título II", "", ""},
		{"Capítulo I", "Capitulo", ""},
		{"Capitulo II", "Capitulo", ""},
		{"Capítulo IV", "Capitulo", ""},
		{"Capítulo Único", "", ""},
		{"Este es no es capítulo", "", ""},
		{"Este es no es Capitulo", "", ""},
	}
	notFound := keywordTests

	tags, index = FindCTags(in, wg)

	go func() {
		for _, tt := range keywordTests {
			in <- tt.in
		}
		close(in)
	}()
	<-index

	for tag := range tags {
		found := tag
		for i, keyword := range notFound {
			if keyword.in == found.text {
				if keyword.out == found.tagname {
					notFound = append(notFound[:i], notFound[i+1:]...)
				} else {
					notFound[i].found = found.tagname
				}
				break
			}
		}
	}
	// fmt.Println(len(notFound))
	for _, fail := range notFound {
		if fail.out != "" {
			t.Errorf("keyword testing %q, expected %q, actual %q", fail.in, fail.out, fail.found)

		}
	}
}
