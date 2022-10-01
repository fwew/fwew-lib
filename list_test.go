package fwew_lib

import (
	"errors"
	"reflect"
	"testing"
)

func TestList(t *testing.T) {
	type args struct {
		args []string
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
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "pos has adj.",
			args: args{
				args: []string{
					"pos",
					"has",
					"adj.",
				},
			},
			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "pos is vim.",
			args: args{
				args: []string{
					"pos",
					"is",
					"vim.",
				},
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
			},
			wantResults: nil,
			wantErr:     InvalidNumber,
		},
		{
			name: "stress = -1 and word like 'aw",
			args: args{
				args: []string{"stress", "=", "-1", "and", "word", "like", "'aw"},
			},
			wantResults: []Word{
				{
					ID:             "12",
					Navi:           "'aw",
					IPA:            "ʔaw",
					InfixLocations: "NULL",
					PartOfSpeech:   "num.",
					Source:         "NULL",
					Stressed:       "1",
					Syllables:      "'aw",
					InfixDots:      "NULL",
					DE:             "eins",
					EN:             "one",
					ET:             "üks",
					FR:             "1 (un)",
					HU:             "egy, 1",
					NL:             "één",
					PL:             "jeden",
					RU:             "один",
					SV:             "en, ett",
					TR:             "bir",
					Affixes: affix{
						Prefix:   nil,
						Infix:    nil,
						Suffix:   nil,
						Lenition: nil,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "stress != -1 and word like 'aw",
			args: args{
				args: []string{"stress", "!=", "-1", "and", "word", "like", "'aw"},
			},
			wantResults: nil,
			wantErr:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResults, err := List(tt.args.args)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// test new negative stress feature
			if tt.name == "stress = -1 and word like 'aw" || tt.name == "stress != -1 and word like 'aw" {
				if !reflect.DeepEqual(gotResults, tt.wantResults) {
					t.Errorf("List() gotResults = %v, want %v", gotResults, tt.wantResults)
				}
			} else if err == nil && len(gotResults) == 0 {
				// for now, only check if something returns
				t.Errorf("List() got empty result, expected something!")
			}
			//if !reflect.DeepEqual(gotResults, tt.wantResults) {
			//	t.Errorf("List() gotResults = %v, want %v", gotResults, tt.wantResults)
			//}
		})
	}
}
