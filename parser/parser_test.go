package parser

import "testing"

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
