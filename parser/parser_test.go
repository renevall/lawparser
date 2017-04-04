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

func TestFindBasicDataName(t *testing.T) {
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
		// {"Aprobado el 8 de Abril de 1988", "1988-04-08 00:00:00 +0000 UTC"},
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

}
