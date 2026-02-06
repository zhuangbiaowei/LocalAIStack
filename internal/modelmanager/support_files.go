package modelmanager

import "strings"

var requiredModelFiles = []string{
	"config.json",
	"params.json",
	"tokenizer.json",
	"tokenizer.model",
	"tokenizer_config.json",
	"special_tokens_map.json",
	"vocab.json",
	"merges.txt",
	"generation_config.json",
	"preprocessor_config.json",
	"image_processor.json",
	"feature_extractor.json",
}

func IsRequiredModelFile(base string) bool {
	base = strings.ToLower(strings.TrimSpace(base))
	for _, name := range requiredModelFiles {
		if base == name {
			return true
		}
	}
	return false
}

func RequiredModelFiles() []string {
	out := make([]string, len(requiredModelFiles))
	copy(out, requiredModelFiles)
	return out
}
