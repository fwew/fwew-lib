package fwew_lib

import (
	"errors"
	"testing"
)

func TestList(t *testing.T) {
	type args struct {
		args     []string
		langCode string
	}
	tests := []struct {
		name        string
		args        args
		wantResults []Word
		wantErr     error
	}{
		// TODO: Add test cases.
		{
			name: "pos starts v",
			args: args{
				args: []string{
					"pos",
					"starts",
					"v",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "pos has svin.",
			args: args{
				args: []string{
					"pos",
					"has",
					"svin.",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "pos is v.",
			args: args{
				args: []string{
					"pos",
					"is",
					"v.",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "pos not-starts v",
			args: args{
				args: []string{
					"pos",
					"not-starts",
					"v",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "pos not-ends m.",
			args: args{
				args: []string{
					"pos",
					"not-ends",
					"m.",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "pos not-has svin.",
			args: args{
				args: []string{
					"pos",
					"not-has",
					"svin.",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "pos not-is v.",
			args: args{
				args: []string{
					"pos",
					"not-is",
					"v.",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "pos not-like *",
			args: args{
				args: []string{
					"pos",
					"not-like",
					"n",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "word starts ft",
			args: args{
				args: []string{
					"word",
					"starts",
					"ft",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "word ends ang",
			args: args{
				args: []string{
					"word",
					"ends",
					"ang",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "word has ts",
			args: args{
				args: []string{
					"word",
					"has",
					"ts",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "word like *",
			args: args{
				args: []string{
					"word",
					"like",
					"'u",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "word not-starts ft",
			args: args{
				args: []string{
					"word",
					"not-starts",
					"ft",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "word not-ends ang",
			args: args{
				args: []string{
					"word",
					"not-ends",
					"ang",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "word not-has ts",
			args: args{
				args: []string{
					"word",
					"not-has",
					"ts",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "word not-like *",
			args: args{
				args: []string{
					"word",
					"not-like",
					"'u",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "words first 20",
			args: args{
				args: []string{
					"words",
					"first",
					"20",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "words last 30",
			args: args{
				args: []string{
					"words",
					"last",
					"30",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "syllables > 1",
			args: args{
				args: []string{
					"syllables",
					">",
					"1",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "syllables = 2",
			args: args{
				args: []string{
					"syllables",
					"=",
					"2",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "syllables <= 3",
			args: args{
				args: []string{
					"syllables",
					"<=",
					"3",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "pos ends j.",
			args: args{
				args: []string{
					"pos",
					"ends",
					"j.",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "pos like adj.",
			args: args{
				args: []string{
					"pos",
					"like",
					"adj.",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "words last g",
			args: args{
				args: []string{
					"words",
					"last",
					"g",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     InvalidNumber,
		},
		{
			name: "syllables = g",
			args: args{
				args: []string{
					"syllables",
					"=",
					"g",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     InvalidNumber,
		},
		{
			name: "syllables < 2",
			args: args{
				args: []string{
					"syllables",
					"<",
					"2",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "syllables >= 4",
			args: args{
				args: []string{
					"syllables",
					">=",
					"4",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "syllables != 1",
			args: args{
				args: []string{
					"syllables",
					"!=",
					"1",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "syllables != 1 and words last 5",
			args: args{
				args: []string{
					"syllables",
					"!=",
					"1",
					"and",
					"words",
					"last",
					"5",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "syllables != 1 and words",
			args: args{
				args: []string{
					"syllables",
					"!=",
					"1",
					"and",
					"words",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "stress = 1",
			args: args{
				args: []string{
					"stress",
					"=",
					"1",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "stress >= 1",
			args: args{
				args: []string{
					"stress",
					">=",
					"1",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "stress <= 1",
			args: args{
				args: []string{
					"stress",
					"<=",
					"1",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "stress < 2",
			args: args{
				args: []string{
					"stress",
					"<",
					"2",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "stress > 1",
			args: args{
				args: []string{
					"stress",
					">",
					"1",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "stress != 2",
			args: args{
				args: []string{
					"stress",
					"!=",
					"2",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "stress = g",
			args: args{
				args: []string{
					"stress",
					"=",
					"g",
				},
				langCode: "en",
			},
			wantResults: nil,
			wantErr:     InvalidNumber,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResults, err := List(tt.args.args, tt.args.langCode)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// for know, only check if something returns
			if err == nil && len(gotResults) == 0 {
				t.Errorf("List() got empty result, expected something!")
			}
			//if !reflect.DeepEqual(gotResults, tt.wantResults) {
			//	t.Errorf("List() gotResults = %v, want %v", gotResults, tt.wantResults)
			//}
		})
	}
}
