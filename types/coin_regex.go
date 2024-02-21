
//line coin_regex.rl:1
package types

func MatchDenom(data []byte) bool {

//line coin_regex.rl:5

//line coin_regex.go:10
var _scanner_actions []byte = []byte{
	0, 1, 0, 
}

var _scanner_key_offsets []int16 = []int16{
	0, 0, 4, 11, 18, 25, 32, 39, 
	46, 53, 60, 67, 74, 81, 88, 95, 
	102, 109, 116, 123, 130, 137, 144, 151, 
	158, 165, 172, 179, 186, 193, 200, 207, 
	214, 221, 228, 235, 242, 249, 256, 263, 
	270, 277, 284, 291, 298, 305, 312, 319, 
	326, 333, 340, 347, 354, 361, 368, 375, 
	382, 389, 396, 403, 410, 417, 424, 431, 
	438, 445, 452, 459, 466, 473, 480, 487, 
	494, 501, 508, 515, 522, 529, 536, 543, 
	550, 557, 564, 571, 578, 585, 592, 599, 
	606, 613, 620, 627, 634, 641, 648, 655, 
	662, 669, 676, 683, 690, 697, 704, 711, 
	718, 725, 732, 739, 746, 753, 760, 767, 
	774, 781, 788, 795, 802, 809, 816, 823, 
	830, 837, 844, 851, 858, 865, 872, 879, 
	886, 893, 
}

var _scanner_trans_keys []byte = []byte{
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 
}

var _scanner_single_lengths []byte = []byte{
	0, 0, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 0, 
}

var _scanner_range_lengths []byte = []byte{
	0, 2, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 0, 
}

var _scanner_index_offsets []int16 = []int16{
	0, 0, 3, 8, 13, 18, 23, 28, 
	33, 38, 43, 48, 53, 58, 63, 68, 
	73, 78, 83, 88, 93, 98, 103, 108, 
	113, 118, 123, 128, 133, 138, 143, 148, 
	153, 158, 163, 168, 173, 178, 183, 188, 
	193, 198, 203, 208, 213, 218, 223, 228, 
	233, 238, 243, 248, 253, 258, 263, 268, 
	273, 278, 283, 288, 293, 298, 303, 308, 
	313, 318, 323, 328, 333, 338, 343, 348, 
	353, 358, 363, 368, 373, 378, 383, 388, 
	393, 398, 403, 408, 413, 418, 423, 428, 
	433, 438, 443, 448, 453, 458, 463, 468, 
	473, 478, 483, 488, 493, 498, 503, 508, 
	513, 518, 523, 528, 533, 538, 543, 548, 
	553, 558, 563, 568, 573, 578, 583, 588, 
	593, 598, 603, 608, 613, 618, 623, 628, 
	633, 638, 
}

var _scanner_indicies []byte = []byte{
	0, 0, 1, 2, 2, 2, 2, 1, 
	3, 3, 3, 3, 1, 4, 4, 4, 
	4, 1, 5, 5, 5, 5, 1, 6, 
	6, 6, 6, 1, 7, 7, 7, 7, 
	1, 8, 8, 8, 8, 1, 9, 9, 
	9, 9, 1, 10, 10, 10, 10, 1, 
	11, 11, 11, 11, 1, 12, 12, 12, 
	12, 1, 13, 13, 13, 13, 1, 14, 
	14, 14, 14, 1, 15, 15, 15, 15, 
	1, 16, 16, 16, 16, 1, 17, 17, 
	17, 17, 1, 18, 18, 18, 18, 1, 
	19, 19, 19, 19, 1, 20, 20, 20, 
	20, 1, 21, 21, 21, 21, 1, 22, 
	22, 22, 22, 1, 23, 23, 23, 23, 
	1, 24, 24, 24, 24, 1, 25, 25, 
	25, 25, 1, 26, 26, 26, 26, 1, 
	27, 27, 27, 27, 1, 28, 28, 28, 
	28, 1, 29, 29, 29, 29, 1, 30, 
	30, 30, 30, 1, 31, 31, 31, 31, 
	1, 32, 32, 32, 32, 1, 33, 33, 
	33, 33, 1, 34, 34, 34, 34, 1, 
	35, 35, 35, 35, 1, 36, 36, 36, 
	36, 1, 37, 37, 37, 37, 1, 38, 
	38, 38, 38, 1, 39, 39, 39, 39, 
	1, 40, 40, 40, 40, 1, 41, 41, 
	41, 41, 1, 42, 42, 42, 42, 1, 
	43, 43, 43, 43, 1, 44, 44, 44, 
	44, 1, 45, 45, 45, 45, 1, 46, 
	46, 46, 46, 1, 47, 47, 47, 47, 
	1, 48, 48, 48, 48, 1, 49, 49, 
	49, 49, 1, 50, 50, 50, 50, 1, 
	51, 51, 51, 51, 1, 52, 52, 52, 
	52, 1, 53, 53, 53, 53, 1, 54, 
	54, 54, 54, 1, 55, 55, 55, 55, 
	1, 56, 56, 56, 56, 1, 57, 57, 
	57, 57, 1, 58, 58, 58, 58, 1, 
	59, 59, 59, 59, 1, 60, 60, 60, 
	60, 1, 61, 61, 61, 61, 1, 62, 
	62, 62, 62, 1, 63, 63, 63, 63, 
	1, 64, 64, 64, 64, 1, 65, 65, 
	65, 65, 1, 66, 66, 66, 66, 1, 
	67, 67, 67, 67, 1, 68, 68, 68, 
	68, 1, 69, 69, 69, 69, 1, 70, 
	70, 70, 70, 1, 71, 71, 71, 71, 
	1, 72, 72, 72, 72, 1, 73, 73, 
	73, 73, 1, 74, 74, 74, 74, 1, 
	75, 75, 75, 75, 1, 76, 76, 76, 
	76, 1, 77, 77, 77, 77, 1, 78, 
	78, 78, 78, 1, 79, 79, 79, 79, 
	1, 80, 80, 80, 80, 1, 81, 81, 
	81, 81, 1, 82, 82, 82, 82, 1, 
	83, 83, 83, 83, 1, 84, 84, 84, 
	84, 1, 85, 85, 85, 85, 1, 86, 
	86, 86, 86, 1, 87, 87, 87, 87, 
	1, 88, 88, 88, 88, 1, 89, 89, 
	89, 89, 1, 90, 90, 90, 90, 1, 
	91, 91, 91, 91, 1, 92, 92, 92, 
	92, 1, 93, 93, 93, 93, 1, 94, 
	94, 94, 94, 1, 95, 95, 95, 95, 
	1, 96, 96, 96, 96, 1, 97, 97, 
	97, 97, 1, 98, 98, 98, 98, 1, 
	99, 99, 99, 99, 1, 100, 100, 100, 
	100, 1, 101, 101, 101, 101, 1, 102, 
	102, 102, 102, 1, 103, 103, 103, 103, 
	1, 104, 104, 104, 104, 1, 105, 105, 
	105, 105, 1, 106, 106, 106, 106, 1, 
	107, 107, 107, 107, 1, 108, 108, 108, 
	108, 1, 109, 109, 109, 109, 1, 110, 
	110, 110, 110, 1, 111, 111, 111, 111, 
	1, 112, 112, 112, 112, 1, 113, 113, 
	113, 113, 1, 114, 114, 114, 114, 1, 
	115, 115, 115, 115, 1, 116, 116, 116, 
	116, 1, 117, 117, 117, 117, 1, 118, 
	118, 118, 118, 1, 119, 119, 119, 119, 
	1, 120, 120, 120, 120, 1, 121, 121, 
	121, 121, 1, 122, 122, 122, 122, 1, 
	123, 123, 123, 123, 1, 124, 124, 124, 
	124, 1, 125, 125, 125, 125, 1, 126, 
	126, 126, 126, 1, 127, 127, 127, 127, 
	1, 128, 128, 128, 128, 1, 1, 
}

var _scanner_trans_targs []byte = []byte{
	2, 0, 3, 4, 5, 6, 7, 8, 
	9, 10, 11, 12, 13, 14, 15, 16, 
	17, 18, 19, 20, 21, 22, 23, 24, 
	25, 26, 27, 28, 29, 30, 31, 32, 
	33, 34, 35, 36, 37, 38, 39, 40, 
	41, 42, 43, 44, 45, 46, 47, 48, 
	49, 50, 51, 52, 53, 54, 55, 56, 
	57, 58, 59, 60, 61, 62, 63, 64, 
	65, 66, 67, 68, 69, 70, 71, 72, 
	73, 74, 75, 76, 77, 78, 79, 80, 
	81, 82, 83, 84, 85, 86, 87, 88, 
	89, 90, 91, 92, 93, 94, 95, 96, 
	97, 98, 99, 100, 101, 102, 103, 104, 
	105, 106, 107, 108, 109, 110, 111, 112, 
	113, 114, 115, 116, 117, 118, 119, 120, 
	121, 122, 123, 124, 125, 126, 127, 128, 
	129, 
}

var _scanner_trans_actions []byte = []byte{
	0, 0, 0, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 
}

const scanner_start int = 1
const scanner_first_final int = 4
const scanner_error int = 0

const scanner_en_main int = 1


//line coin_regex.rl:6
    cs, p, pe, eof := 0, 0, len(data), len(data)
    _ = eof
    
//line coin_regex.go:344
	{
	cs = scanner_start
	}

//line coin_regex.go:349
	{
	var _klen int
	var _trans int
	var _acts int
	var _nacts uint
	var _keys int
	if p == pe {
		goto _test_eof
	}
	if cs == 0 {
		goto _out
	}
_resume:
	_keys = int(_scanner_key_offsets[cs])
	_trans = int(_scanner_index_offsets[cs])

	_klen = int(_scanner_single_lengths[cs])
	if _klen > 0 {
		_lower := int(_keys)
		var _mid int
		_upper := int(_keys + _klen - 1)
		for {
			if _upper < _lower {
				break
			}

			_mid = _lower + ((_upper - _lower) >> 1)
			switch {
			case data[p] < _scanner_trans_keys[_mid]:
				_upper = _mid - 1
			case data[p] > _scanner_trans_keys[_mid]:
				_lower = _mid + 1
			default:
				_trans += int(_mid - int(_keys))
				goto _match
			}
		}
		_keys += _klen
		_trans += _klen
	}

	_klen = int(_scanner_range_lengths[cs])
	if _klen > 0 {
		_lower := int(_keys)
		var _mid int
		_upper := int(_keys + (_klen << 1) - 2)
		for {
			if _upper < _lower {
				break
			}

			_mid = _lower + (((_upper - _lower) >> 1) & ^1)
			switch {
			case data[p] < _scanner_trans_keys[_mid]:
				_upper = _mid - 2
			case data[p] > _scanner_trans_keys[_mid + 1]:
				_lower = _mid + 2
			default:
				_trans += int((_mid - int(_keys)) >> 1)
				goto _match
			}
		}
		_trans += _klen
	}

_match:
	_trans = int(_scanner_indicies[_trans])
	cs = int(_scanner_trans_targs[_trans])

	if _scanner_trans_actions[_trans] == 0 {
		goto _again
	}

	_acts = int(_scanner_trans_actions[_trans])
	_nacts = uint(_scanner_actions[_acts]); _acts++
	for ; _nacts > 0; _nacts-- {
		_acts++
		switch _scanner_actions[_acts-1] {
		case 0:
//line coin_regex.rl:16
 return true 
//line coin_regex.go:431
		}
	}

_again:
	if cs == 0 {
		goto _out
	}
	p++
	if p != pe {
		goto _resume
	}
	_test_eof: {}
	_out: {}
	}

//line coin_regex.rl:20

    return false
}

func MatchDecCoin(data []byte) (amountStart, amountEnd, denomEnd int, isValid bool) {
    
//line coin_regex.rl:26
    
//line coin_regex.go:456
var _dec_coin_actions []byte = []byte{
	0, 1, 0, 1, 1, 1, 2, 
}

var _dec_coin_key_offsets []int16 = []int16{
	0, 0, 3, 5, 12, 19, 26, 35, 
	42, 49, 56, 63, 70, 77, 84, 91, 
	98, 105, 112, 119, 126, 133, 140, 147, 
	154, 161, 168, 175, 182, 189, 196, 203, 
	210, 217, 224, 231, 238, 245, 252, 259, 
	266, 273, 280, 287, 294, 301, 308, 315, 
	322, 329, 336, 343, 350, 357, 364, 371, 
	378, 385, 392, 399, 406, 413, 420, 427, 
	434, 441, 448, 455, 462, 469, 476, 483, 
	490, 497, 504, 511, 518, 525, 532, 539, 
	546, 553, 560, 567, 574, 581, 588, 595, 
	602, 609, 616, 623, 630, 637, 644, 651, 
	658, 665, 672, 679, 686, 693, 700, 707, 
	714, 721, 728, 735, 742, 749, 756, 763, 
	770, 777, 784, 791, 798, 805, 812, 819, 
	826, 833, 840, 847, 854, 861, 868, 875, 
	882, 889, 896, 903, 910, 910, 
}

var _dec_coin_trans_keys []byte = []byte{
	46, 48, 57, 48, 57, 32, 9, 13, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 32, 9, 13, 48, 57, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 95, 45, 
	58, 65, 90, 97, 122, 95, 45, 58, 
	65, 90, 97, 122, 95, 45, 58, 65, 
	90, 97, 122, 95, 45, 58, 65, 90, 
	97, 122, 95, 45, 58, 65, 90, 97, 
	122, 95, 45, 58, 65, 90, 97, 122, 
	95, 45, 58, 65, 90, 97, 122, 95, 
	45, 58, 65, 90, 97, 122, 32, 46, 
	9, 13, 48, 57, 65, 90, 97, 122, 
	
}

var _dec_coin_single_lengths []byte = []byte{
	0, 1, 0, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 1, 1, 1, 1, 
	1, 1, 1, 1, 0, 2, 
}

var _dec_coin_range_lengths []byte = []byte{
	0, 1, 1, 3, 3, 3, 4, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 0, 4, 
}

var _dec_coin_index_offsets []int16 = []int16{
	0, 0, 3, 5, 10, 15, 20, 26, 
	31, 36, 41, 46, 51, 56, 61, 66, 
	71, 76, 81, 86, 91, 96, 101, 106, 
	111, 116, 121, 126, 131, 136, 141, 146, 
	151, 156, 161, 166, 171, 176, 181, 186, 
	191, 196, 201, 206, 211, 216, 221, 226, 
	231, 236, 241, 246, 251, 256, 261, 266, 
	271, 276, 281, 286, 291, 296, 301, 306, 
	311, 316, 321, 326, 331, 336, 341, 346, 
	351, 356, 361, 366, 371, 376, 381, 386, 
	391, 396, 401, 406, 411, 416, 421, 426, 
	431, 436, 441, 446, 451, 456, 461, 466, 
	471, 476, 481, 486, 491, 496, 501, 506, 
	511, 516, 521, 526, 531, 536, 541, 546, 
	551, 556, 561, 566, 571, 576, 581, 586, 
	591, 596, 601, 606, 611, 616, 621, 626, 
	631, 636, 641, 646, 651, 652, 
}

var _dec_coin_indicies []byte = []byte{
	0, 2, 1, 3, 1, 4, 4, 5, 
	5, 1, 6, 6, 6, 6, 1, 7, 
	7, 7, 7, 1, 8, 8, 3, 9, 
	9, 1, 10, 10, 10, 10, 1, 11, 
	11, 11, 11, 1, 12, 12, 12, 12, 
	1, 13, 13, 13, 13, 1, 14, 14, 
	14, 14, 1, 15, 15, 15, 15, 1, 
	16, 16, 16, 16, 1, 17, 17, 17, 
	17, 1, 18, 18, 18, 18, 1, 19, 
	19, 19, 19, 1, 20, 20, 20, 20, 
	1, 21, 21, 21, 21, 1, 22, 22, 
	22, 22, 1, 23, 23, 23, 23, 1, 
	24, 24, 24, 24, 1, 25, 25, 25, 
	25, 1, 26, 26, 26, 26, 1, 27, 
	27, 27, 27, 1, 28, 28, 28, 28, 
	1, 29, 29, 29, 29, 1, 30, 30, 
	30, 30, 1, 31, 31, 31, 31, 1, 
	32, 32, 32, 32, 1, 33, 33, 33, 
	33, 1, 34, 34, 34, 34, 1, 35, 
	35, 35, 35, 1, 36, 36, 36, 36, 
	1, 37, 37, 37, 37, 1, 38, 38, 
	38, 38, 1, 39, 39, 39, 39, 1, 
	40, 40, 40, 40, 1, 41, 41, 41, 
	41, 1, 42, 42, 42, 42, 1, 43, 
	43, 43, 43, 1, 44, 44, 44, 44, 
	1, 45, 45, 45, 45, 1, 46, 46, 
	46, 46, 1, 47, 47, 47, 47, 1, 
	48, 48, 48, 48, 1, 49, 49, 49, 
	49, 1, 50, 50, 50, 50, 1, 51, 
	51, 51, 51, 1, 52, 52, 52, 52, 
	1, 53, 53, 53, 53, 1, 54, 54, 
	54, 54, 1, 55, 55, 55, 55, 1, 
	56, 56, 56, 56, 1, 57, 57, 57, 
	57, 1, 58, 58, 58, 58, 1, 59, 
	59, 59, 59, 1, 60, 60, 60, 60, 
	1, 61, 61, 61, 61, 1, 62, 62, 
	62, 62, 1, 63, 63, 63, 63, 1, 
	64, 64, 64, 64, 1, 65, 65, 65, 
	65, 1, 66, 66, 66, 66, 1, 67, 
	67, 67, 67, 1, 68, 68, 68, 68, 
	1, 69, 69, 69, 69, 1, 70, 70, 
	70, 70, 1, 71, 71, 71, 71, 1, 
	72, 72, 72, 72, 1, 73, 73, 73, 
	73, 1, 74, 74, 74, 74, 1, 75, 
	75, 75, 75, 1, 76, 76, 76, 76, 
	1, 77, 77, 77, 77, 1, 78, 78, 
	78, 78, 1, 79, 79, 79, 79, 1, 
	80, 80, 80, 80, 1, 81, 81, 81, 
	81, 1, 82, 82, 82, 82, 1, 83, 
	83, 83, 83, 1, 84, 84, 84, 84, 
	1, 85, 85, 85, 85, 1, 86, 86, 
	86, 86, 1, 87, 87, 87, 87, 1, 
	88, 88, 88, 88, 1, 89, 89, 89, 
	89, 1, 90, 90, 90, 90, 1, 91, 
	91, 91, 91, 1, 92, 92, 92, 92, 
	1, 93, 93, 93, 93, 1, 94, 94, 
	94, 94, 1, 95, 95, 95, 95, 1, 
	96, 96, 96, 96, 1, 97, 97, 97, 
	97, 1, 98, 98, 98, 98, 1, 99, 
	99, 99, 99, 1, 100, 100, 100, 100, 
	1, 101, 101, 101, 101, 1, 102, 102, 
	102, 102, 1, 103, 103, 103, 103, 1, 
	104, 104, 104, 104, 1, 105, 105, 105, 
	105, 1, 106, 106, 106, 106, 1, 107, 
	107, 107, 107, 1, 108, 108, 108, 108, 
	1, 109, 109, 109, 109, 1, 110, 110, 
	110, 110, 1, 111, 111, 111, 111, 1, 
	112, 112, 112, 112, 1, 113, 113, 113, 
	113, 1, 114, 114, 114, 114, 1, 115, 
	115, 115, 115, 1, 116, 116, 116, 116, 
	1, 117, 117, 117, 117, 1, 118, 118, 
	118, 118, 1, 119, 119, 119, 119, 1, 
	120, 120, 120, 120, 1, 121, 121, 121, 
	121, 1, 122, 122, 122, 122, 1, 123, 
	123, 123, 123, 1, 124, 124, 124, 124, 
	1, 125, 125, 125, 125, 1, 126, 126, 
	126, 126, 1, 127, 127, 127, 127, 1, 
	128, 128, 128, 128, 1, 129, 129, 129, 
	129, 1, 130, 130, 130, 130, 1, 131, 
	131, 131, 131, 1, 132, 132, 132, 132, 
	1, 133, 133, 133, 133, 1, 134, 134, 
	134, 134, 1, 1, 8, 135, 8, 136, 
	9, 9, 1, 
}

var _dec_coin_trans_targs []byte = []byte{
	2, 0, 133, 6, 3, 4, 5, 7, 
	3, 4, 8, 9, 10, 11, 12, 13, 
	14, 15, 16, 17, 18, 19, 20, 21, 
	22, 23, 24, 25, 26, 27, 28, 29, 
	30, 31, 32, 33, 34, 35, 36, 37, 
	38, 39, 40, 41, 42, 43, 44, 45, 
	46, 47, 48, 49, 50, 51, 52, 53, 
	54, 55, 56, 57, 58, 59, 60, 61, 
	62, 63, 64, 65, 66, 67, 68, 69, 
	70, 71, 72, 73, 74, 75, 76, 77, 
	78, 79, 80, 81, 82, 83, 84, 85, 
	86, 87, 88, 89, 90, 91, 92, 93, 
	94, 95, 96, 97, 98, 99, 100, 101, 
	102, 103, 104, 105, 106, 107, 108, 109, 
	110, 111, 112, 113, 114, 115, 116, 117, 
	118, 119, 120, 121, 122, 123, 124, 125, 
	126, 127, 128, 129, 130, 131, 132, 2, 
	133, 
}

var _dec_coin_trans_actions []byte = []byte{
	1, 0, 1, 0, 0, 0, 0, 0, 
	3, 3, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 
}

var _dec_coin_eof_actions []byte = []byte{
	0, 0, 0, 0, 0, 0, 3, 5, 
	5, 5, 5, 5, 5, 5, 5, 5, 
	5, 5, 5, 5, 5, 5, 5, 5, 
	5, 5, 5, 5, 5, 5, 5, 5, 
	5, 5, 5, 5, 5, 5, 5, 5, 
	5, 5, 5, 5, 5, 5, 5, 5, 
	5, 5, 5, 5, 5, 5, 5, 5, 
	5, 5, 5, 5, 5, 5, 5, 5, 
	5, 5, 5, 5, 5, 5, 5, 5, 
	5, 5, 5, 5, 5, 5, 5, 5, 
	5, 5, 5, 5, 5, 5, 5, 5, 
	5, 5, 5, 5, 5, 5, 5, 5, 
	5, 5, 5, 5, 5, 5, 5, 5, 
	5, 5, 5, 5, 5, 5, 5, 5, 
	5, 5, 5, 5, 5, 5, 5, 5, 
	5, 5, 5, 5, 5, 5, 5, 5, 
	5, 5, 5, 5, 5, 3, 
}

const dec_coin_start int = 1
const dec_coin_first_final int = 6
const dec_coin_error int = 0

const dec_coin_en_main int = 1


//line coin_regex.rl:27

    // Initialize positions and validity flag
    amountStart, amountEnd, denomEnd = -1, -1, -1
    isValid = false

    // Ragel state variables
    var cs, p, pe, eof int
    p, pe, eof = 0, len(data), len(data)

    
//line coin_regex.go:826
	{
	cs = dec_coin_start
	}

//line coin_regex.go:831
	{
	var _klen int
	var _trans int
	var _acts int
	var _nacts uint
	var _keys int
	if p == pe {
		goto _test_eof
	}
	if cs == 0 {
		goto _out
	}
_resume:
	_keys = int(_dec_coin_key_offsets[cs])
	_trans = int(_dec_coin_index_offsets[cs])

	_klen = int(_dec_coin_single_lengths[cs])
	if _klen > 0 {
		_lower := int(_keys)
		var _mid int
		_upper := int(_keys + _klen - 1)
		for {
			if _upper < _lower {
				break
			}

			_mid = _lower + ((_upper - _lower) >> 1)
			switch {
			case data[p] < _dec_coin_trans_keys[_mid]:
				_upper = _mid - 1
			case data[p] > _dec_coin_trans_keys[_mid]:
				_lower = _mid + 1
			default:
				_trans += int(_mid - int(_keys))
				goto _match
			}
		}
		_keys += _klen
		_trans += _klen
	}

	_klen = int(_dec_coin_range_lengths[cs])
	if _klen > 0 {
		_lower := int(_keys)
		var _mid int
		_upper := int(_keys + (_klen << 1) - 2)
		for {
			if _upper < _lower {
				break
			}

			_mid = _lower + (((_upper - _lower) >> 1) & ^1)
			switch {
			case data[p] < _dec_coin_trans_keys[_mid]:
				_upper = _mid - 2
			case data[p] > _dec_coin_trans_keys[_mid + 1]:
				_lower = _mid + 2
			default:
				_trans += int((_mid - int(_keys)) >> 1)
				goto _match
			}
		}
		_trans += _klen
	}

_match:
	_trans = int(_dec_coin_indicies[_trans])
	cs = int(_dec_coin_trans_targs[_trans])

	if _dec_coin_trans_actions[_trans] == 0 {
		goto _again
	}

	_acts = int(_dec_coin_trans_actions[_trans])
	_nacts = uint(_dec_coin_actions[_acts]); _acts++
	for ; _nacts > 0; _nacts-- {
		_acts++
		switch _dec_coin_actions[_acts-1] {
		case 0:
//line coin_regex.rl:37

            amountStart = p;
        
		case 1:
//line coin_regex.rl:40

            amountEnd = p;
        
//line coin_regex.go:920
		}
	}

_again:
	if cs == 0 {
		goto _out
	}
	p++
	if p != pe {
		goto _resume
	}
	_test_eof: {}
	if p == eof {
		__acts := _dec_coin_eof_actions[cs]
		__nacts := uint(_dec_coin_actions[__acts]); __acts++
		for ; __nacts > 0; __nacts-- {
			__acts++
			switch _dec_coin_actions[__acts-1] {
			case 1:
//line coin_regex.rl:40

            amountEnd = p;
        
			case 2:
//line coin_regex.rl:43

            denomEnd = p-1; // Adjusted to exclude space if present
        
//line coin_regex.go:949
			}
		}
	}

	_out: {}
	}

//line coin_regex.rl:60


    isValid = (cs >= 6);

    // Return the captured positions and validity
    return amountStart, amountEnd, denomEnd, isValid
}
