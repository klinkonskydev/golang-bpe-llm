package tokenizer

import (
	"sort"
	"strings"
)

type merges struct {
	a, b, result string
}

type BPE struct {
	vocab map[string]int
	iVocab []string
	merges []merges
}

func (b *BPE) VocabSize() int {
	return len(b.vocab)
}

func (b *BPE) Train(corpus string, vocabSize int) {
	words := strings.Fields(corpus)

	seqs := make([][]string, len(words))
	for i, word := range words {
		chars := make([]string, 0, len(word)+1)
		for _, ch := range word {
			chars = append(chars, string(ch))
		}

		chars = append(chars, "</w>")
		seqs[i] = chars
	}

	b.vocab = make(map[string]int)
	for _, seq := range seqs {
		for _, word := range seq {
			if _, ok := b.vocab[word]; !ok {
				b.vocab[word] = len(b.vocab)
			}
		}
	}

	for len(b.vocab) < vocabSize {
		counts := countPairs(seqs)
		if len(counts) == 0 {
			break
		}

		pair := bestPairs(counts)
		merged := pair[0] + pair[1]
		b.merges = append(b.merges, merges{pair[0], pair[1], merged})
		b.vocab[merged] = len(b.vocab)

		seqs = applyMergeAll(seqs, pair[0], pair[1], merged)
	}

	b.iVocab = make([]string, len(b.vocab))
	for token, id := range b.vocab {
		b.iVocab[id] = token
	}
}

func countPairs(seq [][]string) map[[2]string]int {
	counts := make(map[[2]string]int)
	for _, id := range seq {
		for i := 0; i+1 < len(id); i++ {
			counts[[2]string{id[i], id[i+1]}]++
		}
	}

	return counts
}

func (b *BPE) Encode(text string) []int {
	var ids []int

	for _, word := range strings.Fields(text) {
		chars := make([]string, 0, len(word) + 1)
		for _, ch := range word {
			chars = append(chars, string(ch))
		}

		chars = append(chars, "</w>")

		for _, ch := range b.merges {
			chars = applyMergeSeq(chars, ch.a, ch.b, ch.result)
		}

		for _, ch := range chars {
			if id, ok := b.vocab[ch]; ok {
				ids = append(ids, id)
			}
		}
	}

	return ids
} 

func (b *BPE) Decode(ids []int) string {
	var sb strings.Builder

	for _, id := range ids {
		if id >= len(b.iVocab) {
			continue
		}

		tok := b.iVocab[id]
		if strings.HasSuffix(tok, "</w>") {
			sb.WriteString(strings.TrimSuffix(tok, "</w>"))
			sb.WriteByte(' ')
		} else {
			sb.WriteString(tok)
		}
	}

	return strings.TrimSpace(sb.String())
} 

func bestPairs(counts map[[2]string]int) [2]string {
	type entry struct {
		pair [2]string
		count int
	}

	entries := make([]entry, 0, len(counts))
	for pair, count := range counts {
		entries = append(entries, entry{pair, count})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].count == entries[j].count {
			return entries[i].pair[0] < entries[j].pair[0]
		}

		return entries[i].pair[0] < entries[j].pair[0]
	})

	return entries[0].pair
}

func applyMergeAll(seq [][]string, a, b, merged string) [][]string {
	out := make([][]string, len(seq))
	for i, seq := range seq {
		out[i] = applyMergeSeq(seq, a, b, merged)
	}

	return out
}

func applyMergeSeq(seq []string, a, b, merged string) []string {
	out := make([]string, 0, len(seq))
	i := 0

	for i < len(seq) {
		if i+1 < len(seq) && seq[i] == a && seq[i+1] == b {
			out = append(out, merged)
			i += 2
		} else {
			out = append(out, seq[i])
			i++
		}
	}

	return out
}
