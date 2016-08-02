package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/japanese"
)

var (
	ja    = flag.Bool("ja", false, "language")
	count = flag.Int("count", 30, "count of words")
)

func loremEn(count int) (string, error) {
	if count <= 0 {
		count = 30
	}

	common := []string{"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipisicing", "elit"}
	words := []string{
		"exercitationem", "perferendis", "perspiciatis", "laborum", "eveniet",
		"sunt", "iure", "nam", "nobis", "eum", "cum", "officiis", "excepturi",
		"odio", "consectetur", "quasi", "aut", "quisquam", "vel", "eligendi",
		"itaque", "non", "odit", "tempore", "quaerat", "dignissimos",
		"facilis", "neque", "nihil", "expedita", "vitae", "vero", "ipsum",
		"nisi", "animi", "cumque", "pariatur", "velit", "modi", "natus",
		"iusto", "eaque", "sequi", "illo", "sed", "ex", "et", "voluptatibus",
		"tempora", "veritatis", "ratione", "assumenda", "incidunt", "nostrum",
		"placeat", "aliquid", "fuga", "provident", "praesentium", "rem",
		"necessitatibus", "suscipit", "adipisci", "quidem", "possimus",
		"voluptas", "debitis", "sint", "accusantium", "unde", "sapiente",
		"voluptate", "qui", "aspernatur", "laudantium", "soluta", "amet",
		"quo", "aliquam", "saepe", "culpa", "libero", "ipsa", "dicta",
		"reiciendis", "nesciunt", "doloribus", "autem", "impedit", "minima",
		"maiores", "repudiandae", "ipsam", "obcaecati", "ullam", "enim",
		"totam", "delectus", "ducimus", "quis", "voluptates", "dolores",
		"molestiae", "harum", "dolorem", "quia", "voluptatem", "molestias",
		"magni", "distinctio", "omnis", "illum", "dolorum", "voluptatum", "ea",
		"quas", "quam", "corporis", "quae", "blanditiis", "atque", "deserunt",
		"laboriosam", "earum", "consequuntur", "hic", "cupiditate",
		"quibusdam", "accusamus", "ut", "rerum", "error", "minus", "eius",
		"ab", "ad", "nemo", "fugit", "officia", "at", "in", "id", "quos",
		"reprehenderit", "numquam", "iste", "fugiat", "sit", "inventore",
		"beatae", "repellendus", "magnam", "recusandae", "quod", "explicabo",
		"doloremque", "aperiam", "consequatur", "asperiores", "commodi",
		"optio", "dolor", "labore", "temporibus", "repellat", "veniam",
		"architecto", "est", "esse", "mollitia", "nulla", "a", "similique",
		"eos", "alias", "dolore", "tenetur", "deleniti", "porro", "facere",
		"maxime", "corrupti"}

	ret := []string{}
	sentence := 0
	for i := 0; i < count; i++ {
		if sentence > 0 {
			common = append(common, words...)
		}
		r := rand.Int() % len(common)
		word := common[r%len(common)]
		if sentence == 0 {
			word = string(unicode.ToUpper(rune(word[0]))) + word[1:]
		}
		sentence++
		ret = append(ret, word)
		if (sentence > 5 && rand.Int() < 10000) || i == count-1 {
			endc := ""
			if i == count-1 {
				endc = string("?!..."[rand.Int()%5])
				ret = append(ret, endc)
			} else {
				endc = string("?!,..."[rand.Int()%6])
				ret = append(ret, endc+" ")
			}
			if strings.Index(endc, ",") != -1 {
				sentence = 0
			}
		} else {
			ret = append(ret, " ")
		}
	}
	return strings.Join(ret, ""), nil
}

func loremJa(count int) (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", fmt.Errorf("$HOME not found")
	}
	fname := filepath.Join(home, ".lorem")

	text := ""
	_, err := os.Stat(fname)
	if err != nil {
		resp, err := http.Get("http://www.aozora.gr.jp/cards/000081/files/470_15407.html")
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		f, err := os.Create(fname)
		if err != nil {
			return "", err
		}
		defer f.Close()

		root, err := html.Parse(japanese.ShiftJIS.NewDecoder().Reader(resp.Body))
		if err != nil {
			return "", err
		}

		contents := scrape.FindAllNested(root, scrape.ByClass("main_text"))
		if len(contents) == 0 {
			return "", fmt.Errorf("contens not found")
		}

		scrape.FindAll(contents[0], func(n *html.Node) bool {
			if n.Type == html.TextNode {
				text += strings.TrimSpace(n.Data)
			} else if n.Type == html.ElementNode && strings.ToLower(n.Data) == "br" {
				text += "\n"
			}
			return false
		})

		out := []string{}
		for _, line := range strings.Split(text, "\n") {
			if len(line) > 0 {
				out = append(out, line)
			}
		}
		text = strings.Join(out, "\n")

		err = ioutil.WriteFile(fname, []byte(text), 0644)
		if err != nil {
			return "", err
		}
	} else {
		b, err := ioutil.ReadFile(fname)
		if err != nil {
			return "", err
		}
		text = string(b)
	}

	output := ""
	lines := strings.Split(text, "\n")
	for {
		c := lines[rand.Int()%len(lines)]
		if len([]rune(output+c)) > count {
			break
		}
		output += c
	}
	if len(output) == 0 {
		mlines := []string{}
		for _, line := range lines {
			if len([]rune(line)) <= count {
				mlines = append(mlines, line)
			}
		}
		if len(mlines) > 0 {
			return mlines[rand.Int()%len(mlines)], nil
		}
	}
	return output, nil
}

func main() {
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	var content string
	var err error
	if *ja {
		content, err = loremJa(*count)
	} else {
		content, err = loremEn(*count)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(content)
}
