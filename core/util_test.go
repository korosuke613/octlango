package core

import (
	"reflect"
	"testing"
)

var reverseSampleLanguageSizes = func() []LanguageSize {
	return []LanguageSize{
		SampleLanguageSizes[2],
		SampleLanguageSizes[1],
		SampleLanguageSizes[0],
	}
}()

func Test_reverseSlice(t *testing.T) {
	type args struct {
		s []LanguageSize
	}
	tests := []struct {
		name string
		args args
		want []LanguageSize
	}{
		{
			name: "順番が逆になる",
			args: args{s: SampleLanguageSizes},
			want: reverseSampleLanguageSizes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reverseSlice(tt.args.s)
			got := tt.args.s
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRepositoriesContributedTo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
