package main

import (
	"fmt"

	"github.com/klinkonskydev/first_go_llm/internal/tokenizer"
)

func main() {
	bpe := &tokenizer.BPE{}
	bpe.Train("o gato subiu no tenhado e o gato dormiu", 30)
	fmt.Println("vocab: ", bpe.VocabSize())
	ids := bpe.Encode("o gato subiu")
	fmt.Println("ids: ", ids)
	fmt.Println("decoded: ", bpe.Decode(ids))
}
