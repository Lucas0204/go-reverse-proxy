package parser

import (
	"fmt"
	"go-reverse-proxy/config-file"
	"os"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type RawConfig struct {
	Blocks []*RawBlock `@@*`
}

type RawBlock struct {
	Name     string      `@Value`
	Args     []string    `@Value*`
	Open     string      `"{"`
	Children []*RawEntry `@@*`
	Close    string      `"}"`
}

type RawEntry struct {
	Block     *RawBlock     `@@`
	Directive *RawDirective `| @@`
}

type RawDirective struct {
	Name string   `@Value`
	Args []string `@Value*`
	End  string   `";"`
}

var nginxLexer = lexer.MustSimple([]lexer.SimpleRule{
	{"Whitespace", `[ \t\n\r]+`},
	{"Punct", `[{};]`},
	{"Value", `[^\s;{}]+`},
})

func ParseConfig(configFilepath string) *config_file.Config {
	parser := participle.MustBuild[RawConfig](
		participle.Lexer(nginxLexer),
		participle.Elide("Whitespace"),
		participle.UseLookahead(4),
	)

	file, err := os.Open(configFilepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	raw, err := parser.Parse(configFilepath, file)
	if err != nil {
		panic(err)
	}

	config := raw.ToConfig()

	if os.Getenv("DEBUG") == "true" {
		executeDebugExamples(config)
	}

	return config
}

func (rc *RawConfig) ToConfig() *config_file.Config {
	cfg := &config_file.Config{}
	for _, rb := range rc.Blocks {
		cfg.Blocks = append(cfg.Blocks, rb.toBlock(nil))
	}
	return cfg
}

func (rb *RawBlock) toBlock(parent *config_file.Block) *config_file.Block {
	b := &config_file.Block{
		Name:   rb.Name,
		Args:   rb.Args,
		Parent: parent,
	}
	for _, child := range rb.Children {
		if child.Block != nil {
			b.Children = append(b.Children, child.Block.toBlock(b))
		} else if child.Directive != nil {
			b.Children = append(b.Children, config_file.Directive{
				Name: child.Directive.Name,
				Args: child.Directive.Args,
			})
		}
	}
	return b
}

func executeDebugExamples(config *config_file.Config) {
	httpBlocks := config.FindBlocksByName("http")
	if len(httpBlocks) > 0 {
		httpBlock := httpBlocks[0]
		servers := httpBlock.FindBlocksByName("server")
		for _, server := range servers {
			if listen, ok := server.GetFirstDirective("listen"); ok {
				fmt.Printf("Server listen: %v\n", listen.Args)
			}
			locations := server.FindBlocksByName("location")
			for _, loc := range locations {
				fmt.Printf("Location path: %v\n", loc.Args)
				if proxyPass, ok := loc.GetFirstDirective("proxy_pass"); ok {
					fmt.Printf("  proxy_pass: %v\n", proxyPass.Args)
				}
			}
		}
	}
	printConfig(config, 0)
}

func printConfig(c *config_file.Config, indent int) {
	for _, block := range c.Blocks {
		PrintBlock(block, indent)
	}
}

func PrintBlock(b *config_file.Block, indent int) {
	prefix := strings.Repeat("  ", indent)
	fmt.Printf("%s%s %v {\n", prefix, b.Name, b.Args)
	for _, child := range b.Children {
		if d, ok := child.(config_file.Directive); ok {
			fmt.Printf("%s  %s %v;\n", prefix, d.Name, d.Args)
		} else if block, ok := child.(*config_file.Block); ok {
			PrintBlock(block, indent+1)
		}
	}
	fmt.Printf("%s}\n", prefix)
}
