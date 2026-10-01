package tokenizer

type Tokenizer interface {
	Train(corpus string, vocabSize int)
	Enconde(text string) []int
	Decode(id int) string
	VocabSize() int
}
