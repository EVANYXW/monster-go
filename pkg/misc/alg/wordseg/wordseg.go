/*******************************************************************************
 * DANIEL'S ALGORITHM IMPLEMENTAIONS
 *
 *  /\  |  _   _  ._ o _|_ |_  ._ _   _
 * /--\ | (_| (_) |  |  |_ | | | | | _>
 *         _|
 *
 * CHINESE WORD SEGMENTATION
 *
 * Features:
 * 1. based on Hidden Markov Model
 * 2. definition the states in B,E,M,S
 * 3. solve the hidden states with viterbi algorithm.
 *
 * http://en.wikipedia.org/wiki/Hidden_Markov_model
 * http://en.wikipedia.org/wiki/Viterbi_algorithm
 *
 ******************************************************************************/
package wordseg

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

// B,E,M,S
const (
	SINGLE = 0
	BEGIN  = 1
	MIDDLE = 2
	END    = 3
)

var (
	// the start probability of state
	// SP = [4]float64{0.23101714745, 0.76898285255, 0.0, 0.0}
	SP = [4]float64{0, 0, 0, 0}
	TP = [4][4]float64{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		/* pre-calculated value
		{0.5138298266610544, 0.48617017333894563, 0.0, 0.0},
		{0.0, 0.0, 0.14817814348183422, 0.8518218565181658},
		{0.0, 0.0, 0.2835512540013088, 0.7164487459986911},
		{0.44551469488355755, 0.554485305116442, 0.0, 0.0},
		*/
	}

	IGNORES = ",. ;。，；"
)

var (
	SPCount = [4]float64{0, 0, 0, 0}
	TPCount = [4][4]float64{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}
)

//------------------------------------------------------ 单个字的观察
type WordEP struct {
	EP [4]float64 // the emission probability in each state
	SC [4]uint32  // the count of appearence in each state
}

type Wordseg struct {
	words map[rune]*WordEP
}

func NewWordseg() *Wordseg {
	ws := &Wordseg{}
	ws.words = make(map[rune]*WordEP)
	return ws
}

func (ws *Wordseg) AddWord(word string, count uint32) {
	runes := []rune(word)

	// return zero
	if len(runes) == 0 {
		return
	}

	// single word
	if len(runes) == 1 {
		if ws.words[runes[0]] == nil {
			ws.words[runes[0]] = &WordEP{}
		}

		ws.words[runes[0]].SC[SINGLE] += count
		TPCount[END][SINGLE] += float64(count / 2)
		TPCount[SINGLE][SINGLE] += float64(count / 2)
		SPCount[SINGLE] += float64(count)
		return
	}

	// more than 1 word
	for k := range runes {
		if ws.words[runes[k]] == nil {
			ws.words[runes[k]] = &WordEP{}
		}

		// begin
		if k == 0 {
			ws.words[runes[k]].SC[BEGIN] += count
			TPCount[END][BEGIN] += float64(count / 2)
			TPCount[SINGLE][BEGIN] += float64(count / 2)
			SPCount[BEGIN] += float64(count)
			continue
		}

		// end
		if k == len(runes)-1 {
			ws.words[runes[k]].SC[END] += count
			if k >= 2 {
				TPCount[MIDDLE][END] += float64(count)
			} else {
				TPCount[BEGIN][END] += float64(count)
			}
			continue
		}

		// otherwise in MIDDLE
		ws.words[runes[k]].SC[MIDDLE] += count
		if k >= 2 {
			TPCount[MIDDLE][MIDDLE] += float64(count)
		} else {
			TPCount[BEGIN][MIDDLE] += float64(count)
		}
	}

	return
}

//----------------------------------------------------- should be called when finished adding words
func (ws *Wordseg) Update() {
	for _, word := range ws.words {
		// total occurence
		sum := float64(word.SC[SINGLE] + word.SC[BEGIN] + word.SC[END] + word.SC[MIDDLE])
		word.EP[SINGLE] = float64(word.SC[SINGLE]) / sum
		word.EP[BEGIN] = float64(word.SC[BEGIN]) / sum
		word.EP[END] = float64(word.SC[END]) / sum
		word.EP[MIDDLE] = float64(word.SC[MIDDLE]) / sum
	}

	for i := 0; i < 4; i++ {
		rowcount := float64(0)
		for j := 0; j < 4; j++ {
			rowcount += TPCount[i][j]
		}

		for j := 0; j < 4; j++ {
			TP[i][j] = TPCount[i][j] / rowcount
		}
	}

	rowcount := float64(0)
	for i := 0; i < 4; i++ {
		rowcount += SPCount[i]
	}

	for i := 0; i < 4; i++ {
		SP[i] = SPCount[i] / rowcount
	}
}

//------------------------------------------------------ words split
func (ws *Wordseg) Split(str string) []string {
	old_runes := []rune(str)
	runes_set := make([][]rune, 0)

	// strip sentence
	runes := make([]rune, 0)
	for _, v := range old_runes {
		if !strings.ContainsRune(IGNORES, v) {
			runes = append(runes, v)
		} else {
			runes_set = append(runes_set, runes)
			runes = make([]rune, 0)
		}
	}
	// final one
	runes_set = append(runes_set, runes)
	// strip unknow chars
	for _, v := range runes_set {
		for k, r := range v {
			if ws.words[r] == nil {
				v[k] = '空'
			}
		}
	}

	result := make([]string, 0)
	for _, v := range runes_set {
		if len(v) > 0 {
			result = append(result, ws._split(v)...)
		}
	}

	return result
}

//---------------------------------------------------------- viterbi algorithm
func (ws *Wordseg) _split(runes []rune) []string {
	length := len(runes)
	// init
	V := make([][]float64, length)
	for k := 0; k < length; k++ {
		V[k] = make([]float64, 4)
	}

	path := make([][]int8, 4)
	for k := 0; k < len(path); k++ {
		path[k] = make([]int8, length)
	}

	newpath := make([][]int8, 4)
	for k := 0; k < len(newpath); k++ {
		newpath[k] = make([]int8, length+1)
	}

	// first observation
	for s := 0; s < 4; s++ {
		V[0][s] = math.Log(SP[s]) + math.Log(ws.words[runes[0]].EP[s])
		path[s][0] = int8(s)
	}

	for k := 1; k < len(runes); k++ {
		word := runes[k]
		wordep := ws.words[word]

		for j := 0; j < 4; j++ {
			prob_max := -math.MaxFloat64
			state_max := int(0)
			ep := wordep.EP[j]

			// relaxation
			for s := 0; s < 4; s++ {
				prob := V[k-1][s] + math.Log(ep) + math.Log(TP[s][j])
				if prob > prob_max {
					prob_max = prob
					state_max = s
				}
			}

			// update maximum
			V[k][j] = prob_max
			copy(newpath[j], path[state_max])
			newpath[j][k] = int8(j)

		}

		for j := 0; j < 4; j++ {
			copy(path[j], newpath[j])
		}
	}

	prob_max := -math.MaxFloat64
	state_max := 0

	for i := 0; i < 4; i++ {
		if V[len(runes)-1][i] > prob_max {
			prob_max = V[len(runes)-1][i]
			state_max = i
		}
	}

	words := make([]string, 0)
	begin := int(0)
	for k, v := range path[state_max] {
		if v == BEGIN {
			begin = k
			continue
		}

		if v == SINGLE {
			words = append(words, string(runes[k]))
			begin = k + 1
			continue
		}

		if v == END {
			words = append(words, string(runes[begin:k+1]))
			begin = k + 1
			continue
		}

	}

	if begin < len(path[state_max]) {
		words = append(words, string(runes[begin:]))
	}

	return words
}

var (
	_default_word_seg *Wordseg
)

func init() {
	ws := NewWordseg()
	ws.LoadFile(os.Getenv("GOPATH") + "/src/misc/alg/wordseg/dict.txt")
	_default_word_seg = ws
}

func (ws *Wordseg) LoadFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("error opening file %v\n", err)
		return
	}

	// using scanner to read config file
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// expression match
		slice := strings.Split(line, " ")

		if slice != nil {
			cnt, _ := strconv.Atoi(slice[1])
			ws.AddWord(slice[0], uint32(cnt))
		}
	}

	ws.Update()
}

//---------------------------------------------------------- default word segmentation
func Split(str string) []string {
	return _default_word_seg.Split(str)
}

//---------------------------------------------------------- return default wordseg
func Default() *Wordseg {
	return _default_word_seg
}
