package config_file

type Config struct {
	Blocks []*Block
}

// Block representa um bloco de configuração (ex: http, server, location)
type Block struct {
	Name     string
	Args     []string
	Children []Entry
	Parent   *Block
}

// Entry pode ser um Block ou uma Directive
type Entry interface {
	IsEntry()
}

// Directive representa uma diretiva simples (ex: listen 80;)
type Directive struct {
	Name string
	Args []string
}

func (Directive) IsEntry() {}
func (*Block) IsEntry()    {}

func (c *Config) FindBlocksByName(name string) []*Block {
	var result []*Block
	for _, b := range c.Blocks {
		result = append(result, b.FindBlocksByName(name)...)
	}
	return result
}

func (b *Block) FindBlocksByName(name string) []*Block {
	var result []*Block
	if b.Name == name {
		result = append(result, b)
	}
	for _, child := range b.Children {
		if childBlock, ok := child.(*Block); ok {
			result = append(result, childBlock.FindBlocksByName(name)...)
		}
	}
	return result
}

// GetDirectives retorna todas as diretivas com o nome dado dentro do bloco (filhos diretos)
func (b *Block) GetDirectives(name string) []Directive {
	var result []Directive
	for _, child := range b.Children {
		if d, ok := child.(Directive); ok && d.Name == name {
			result = append(result, d)
		}
	}
	return result
}

// GetFirstDirective retorna a primeira diretiva com o nome dado dentro do bloco
func (b *Block) GetFirstDirective(name string) (Directive, bool) {
	for _, child := range b.Children {
		if d, ok := child.(Directive); ok && d.Name == name {
			return d, true
		}
	}
	return Directive{}, false
}
