package business

import (
	"encoding/hex"
	"fmt"
	"pdf_raw_printing/internal/libs/ionreader"
	"pdf_raw_printing/internal/libs/wion"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSerialize(t *testing.T) {
	tests := []struct {
		name     string
		v        any
		expected string
	}{
		{
			name: "eidbucket_24",
			v: Eidbucket{
				Block: 24,
				Contains: []ContainsElement{{
					Eid:         "i5",
					SectionName: "c0",
				},
				},
			},
			expected: "e00100eaeea18204e2de9c04da211801b5be94de9201b9e68204d682693501aee68204d6826330",
		},
		{
			name: "cA-spm",
			v: SectionPositionIdMap{
				Contains: []ValueMap{
					{
						ID:        1,
						Reference: "tB",
					},
					{
						ID:        2,
						Reference: "iE",
					},
					{
						ID:        3,
						Reference: "iF",
					},
					{
						ID:        4,
						Reference: "tD",
					},
				},
				SectionName: "cA",
			},
			expected: "e00100eaeeba8204e1deb501aee68204d682634101b5bea8b92101e68204d6827442b92102e68204d6826945b92103e68204d6826946b92104e68204d6827444",
		},
		{
			name:     "book_metadata",
			expected: "e00100eaee02dc8203eade02d603ebbe02d1ded403ef8e956b696e646c655f7469746c655f6d657461646174610282beb7dea503ec87626f6f6b5f696402b38e97644e64724b68566a523236745f635a5f7548596b4f4130de8e03ec857469746c6502b383747474deed03ef8e9a6b696e646c655f6361706162696c6974795f6d657461646174610282becbde9703ec8e8f796a5f66697865645f6c61796f757402b32101de9c03ec8e9467726170686963616c5f686967686c696768747302b32101de9203ec8b796a5f74657874626f6f6b02b32101deb503ef8e956b696e646c655f65626f6f6b5f6d657461646174610282be98de9603ec8973656c656374696f6e02b387656e61626c6564ded303ef8e956b696e646c655f61756469745f6d657461646174610282beb6de9403ec8c66696c655f63726561746f7202b3824b43de9e03ec8e8f63726561746f725f76657273696f6e02b388312e39332e302e30",
		},
		{
			name:     "book_navigation",
			expected: "e00100eaeea7820385bea2dea001b272015f0388be97ee95820387de9001eb71d401efe68204d6826e4101f7b0",
		},
		{
			name:     "c0",
			expected: "e00100eaeefa820284def501aee68204d6826330018dbee8eea78204e0dea204d6e68204d682743101b0e68204d6826c3201abc372020e019c720143019f72010eeebd8204e0deb804d6e68204d6827433b8d902b3216402b272013a01b0e68204d6826c32c2d902b3216402b272013a01abc372020d019c720145019f72010e",
		},
		{
			name:     "c0-ad",
			expected: "e00100eaeeaa8204d5dea504d6e98204d68563302d61640282be95de9303ec8d706167655f726f746174696f6e02b320",
		},
		{
			name:     "c0-spm",
			expected: "e00100eaeeba8204e1deb501aee68204d682633001b5bea8b92101e68204d6827431b92102e68204d6826934b92103e68204d6826935b92104e68204d6827433",
		},
		{
			name:     "cA",
			expected: "e00100eaeefa820284def501aee68204d6826341018dbee8eea78204e0dea204d6e68204d682744201b0e68204d6826c4301abc372020e019c720143019f72010eeebd8204e0deb804d6e68204d6827444b8d902b3216402b272013a01b0e68204d6826c43c2d902b3216402b272013a01abc372020d019c720145019f72010e",
		},
		{
			name:     "cA-ad",
			expected: "e00100eaeeaa8204d5dea504d6e98204d68563412d61640282be95de9303ec8d706167655f726f746174696f6e02b320",
		},
		{
			name:     "d6",
			expected: "e00100eaee01aa8204d5de01a404d6e68204d68264360282be0196de9203ec847479706502b3887265736f75726365de9b03ec8e8f7265736f757263655f73747265616d02b3857273726338de8f03ec8473697a6502b3853132383434de9d03ec8d6d6f6469666965645f74696d6502b38a31373235343339303135deb303ec886c6f636174696f6e02b38ea42f55736572732f726f6e616e6a616d65742f446f776e6c6f6164732f746573742e706466",
		},
		{
			name:     "d7",
			expected: "e00100eaeeb78204d5deb204d6e68204d68264370282bea5dea303ec8e95617578446174615f7265736f757263655f6c69737402b3b7e68204d6826436",
		},
		{
			name:     "document_data",
			expected: "e00100eaeebb82049adeb688211201c072017804c57201b904d5d904e5e68204d682643701a9be99de9701b272015f01aabe8ee68204d6826330e68204d6826341",
		},
		{
			name:     "e9",
			expected: "e00100eaeecd8201a4dec8b04001a1720235b140b2483fe000000000000004b42001a585727372633804d5e68204d682643603a6484082a0000000000003a748408a50000000000001afe68204d6826539af40",
		},
		{
			name:     "eG",
			expected: "e00100eaeece8201a4dec9b04001a1720235b140b2483fe000000000000004b4210101a585727372633804d5e68204d682643603a6484082a0000000000003a748408a50000000000001afe68204d6826547af40",
		},
		{
			name:     "eidbucket_13",
			expected: "e00100eaeea18204e2de9c04da210d01b5be94de9201b9e68204d682633001aee68204d6826330",
		},
		{
			name:     "eidbucket_23",
			expected: "e00100eaeea18204e2de9c04da211701b5be94de9201b9e68204d682693401aee68204d6826330",
		},

		{
			name:     "eidbucket_30",
			expected: "e00100eaeea18204e2de9c04da211e01b5be94de9201b9e68204d682634101aee68204d6826341",
		},
		{
			name:     "eidbucket_31",
			expected: "e00100eaeea18204e2de9c04da211f01b5be94de9201b9e68204d682743101aee68204d6826330",
		},
		{
			name:     "eidbucket_33",
			expected: "e00100eaeea18204e2de9c04da212101b5be94de9201b9e68204d682743301aee68204d6826330",
		},
		{
			name:     "eidbucket_40",
			expected: "e00100eaeea18204e2de9c04da212801b5be94de9201b9e68204d682694501aee68204d6826341",
		},
		{
			name:     "eidbucket_41",
			expected: "e00100eaeea18204e2de9c04da212901b5be94de9201b9e68204d682694601aee68204d6826341",
		},
		{
			name:     "eidbucket_48",
			expected: "e00100eaeea18204e2de9c04da213001b5be94de9201b9e68204d682744201aee68204d6826341",
		},
		{
			name:     "eidbucket_50",
			expected: "e00100eaeea18204e2de9c04da213201b5be94de9201b9e68204d682744401aee68204d6826341",
		},
		{
			name:     "i4",
			expected: "e00100eaeeb58204e0deb004d6e68204d6826934c222e8d0c3230148e803db7201d8019c720146018c720140019f72010e0192b7e68204d6826935",
		},
		{
			name:     "i5",
			expected: "e00100eaeeb28204e0dead04d6e68204d6826935b8d902b3216402b272013ab9d902b3216402b272013a019f72010f01afe68204d6826539",
		},
		{
			name:     "iE",
			expected: "e00100eaeeb58204e0deb004d6e68204d6826945c222e8d0c3230148e803db7201d8019c720146018c720140019f72010e0192b7e68204d6826946",
		},
		{
			name:     "iF",
			expected: "e00100eaeeb28204e0dead04d6e68204d6826946b8d902b3216402b272013ab9d902b3216402b272013a019f72010f01afe68204d6826547",
		},
		{
			name:     "l2",
			expected: "e00100eaee98820283de9301b0e68204d6826c320192b7e68204d6826934",
		},
		{
			name:     "lC",
			expected: "e00100eaee98820283de9301b0e68204d6826c430192b7e68204d6826945",
		},
		{
			name:     "max_id",
			expected: "e00100ea220342",
		},
		{
			name:     "metadata",
			expected: "e00100eaeea2820282de9d01a9be99de9701b272015f01aabe8ee68204d6826330e68204d6826341",
		},
		{
			name:     "rsrc8",
			expected: "7265732f7273726338",
		},
		{
			name:     "yj",
			expected: "e00100eaeea58204e3dea001b5be9cdd01aee68204d682633001902104dd01aee68204d682634101902104",
		},
	}

	for _, ts := range tests {
		if ts.v == nil {
			continue
		}

		t.Run(ts.name, func(t *testing.T) {
			hexenc1, err := hex.DecodeString(ts.expected)
			if err != nil {
				panic(err)
			}

			expectedString, err := ionreader.IonToString(hexenc1)
			if err != nil {
				panic(err)
			}
			fmt.Println(expectedString)

			hashString, err := wion.MarshalString(ts.v)
			require.NoError(t, err)
			hash24, err := wion.Marshal(ts.v)

			require.NoError(t, err)
			hexenc := hex.EncodeToString(hash24)

			require.NoError(t, ionreader.ReadDouble(hexenc1, hash24))

			require.Equal(t, expectedString, hashString)
			require.Equal(t, ts.expected, hexenc)
		})
	}
}
