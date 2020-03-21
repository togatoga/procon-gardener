package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/skratchdot/open-golang/open"
	"github.com/thoas/go-funk"
	cli "github.com/urfave/cli/v2"
)

const APP_NAME = "procon-gardener"
const ATCODER_API_SUBMISSION_URL = "https://kenkoooo.com/atcoder/atcoder-api/results?user="

type AtCoderSubmission struct {
	ID            int     `json:"id"`
	EpochSecond   int     `json:"epoch_second"`
	ProblemID     string  `json:"problem_id"`
	ContestID     string  `json:"contest_id"`
	UserID        string  `json:"user_id"`
	Language      string  `json:"language"`
	Point         float64 `json:"point"`
	Length        int     `json:"length"`
	Result        string  `json:"result"`
	ExecutionTime int     `json:"execution_time"`
}

func isDirExist(path string) bool {
	info, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
func isFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

type Service struct {
	RepositoryPath string `json:"repository_path"`
	UserID         string `json:"user_id"`
}
type Config struct {
	Atcoder Service `json:"atcoder"`
}

func languageToFileName(language string) string {
	//e.g C++14 (GCC 5.4.1)
	//C++14
	language = strings.Split(language, "(")[0]
	//remove extra last whitespace
	language = language[:len(language)-1]
	if strings.HasPrefix(language, "C++") {
		return "Main.cpp"
	}
	if strings.HasPrefix(language, "Bash") {
		return "Main.sh"
	}

	//C (GCC 5.4.1)
	//C (Clang 3.8.0)
	if language == "C" {
		return "Main.c"
	}

	if language == "C#" {
		return "Main.cs"
	}

	if language == "Clojure" {
		return "Main.clj"
	}

	if strings.HasPrefix(language, "Common Lisp") {
		return "Main.lisp"
	}

	//D (DMD64 v2.070.1)
	if language == "D" {
		return "Main.d"
	}

	if language == "Fortran" {
		return "Main.f08"
	}

	if language == "Go" {
		return "Main.go"
	}

	if language == "Haskell" {
		return "Main.hs"
	}

	if language == "JavaScript" {
		return "Main.js"
	}
	if language == "Java" {
		return "Main.java"
	}
	if language == "OCaml" {
		return "Main.ml"
	}

	if language == "Pascal" {
		return "Main.pas"
	}

	if language == "Perl" {
		return "Main.pl"
	}

	if language == "PHP" {
		return "Main.php"
	}

	if strings.HasPrefix(language, "Python") {
		return "Main.py"
	}

	if language == "Ruby" {
		return "Main.rb"
	}

	if language == "Scala" {
		return "Main.scala"
	}

	if language == "Scheme" {
		return "Main.scm"
	}

	if language == "Main.txt" {
		return "Main.txt"
	}

	if language == "Visual Basic" {
		return "Main.vb"
	}

	if language == "Objective-C" {
		return "Main.m"
	}

	if language == "Swift" {
		return "Main.swift"
	}

	if language == "Rust" {
		return "Main.rs"
	}

	if language == "Sed" {
		return "Main.sed"
	}

	if language == "Awk" {
		return "Main.awk"
	}

	if language == "Brainfuck" {
		return "Main.bf"
	}

	if language == "Standard ML" {
		return "Main.sml"
	}

	if strings.HasPrefix(language, "PyPy") {
		return "Main.py"
	}

	if language == "Crystal" {
		return "Main.cr"
	}

	if language == "F#" {
		return "Main.fs"
	}

	if language == "Unlambda" {
		return "Main.unl"
	}

	if language == "Lua" {
		return "Main.lua"
	}

	if language == "LuaJIT" {
		return "Main.lua"
	}

	if language == "MoonScript" {
		return "Main.moon"
	}

	if language == "Ceylon" {
		return "Main.ceylon"
	}

	if language == "Julia" {
		return "Main.jl"
	}

	if language == "Octave" {
		return "Main.m"
	}

	if language == "Nim" {
		return "Main.nim"
	}

	if language == "TypeScript" {
		return "Main.ts"
	}

	if language == "Perl6" {
		return "Main.p6"
	}

	if language == "Kotlin" {
		return "Main.kt"
	}

	if language == "COBOL" {
		return "Main.cob"
	}

	log.Printf("Unknown ... %s", language)
	return "Main.txt"
}

func init() {

	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	configDir := filepath.Join(home, "."+APP_NAME)
	if !isDirExist(configDir) {
		err = os.MkdirAll(configDir, 0700)
		if err != nil {
			panic(err)
		}
	}

	configFile := filepath.Join(configDir, "config.json")
	if !isFileExist(configFile) {
		//initial config
		atcoder := Service{RepositoryPath: "", UserID: ""}

		config := Config{Atcoder: atcoder}

		jsonBytes, err := json.MarshalIndent(config, "", "\t")
		if err != nil {
			panic(err)
		}
		json := string(jsonBytes)
		file, err := os.Create(filepath.Join(configDir, "config.json"))
		if err != nil {
			panic(err)
		}
		defer file.Close()
		file.WriteString(json)
	}
}

func load_config() Config {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	configDir := filepath.Join(home, "."+APP_NAME)
	configFile := filepath.Join(configDir, "config.json")
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	var config Config
	if err = json.Unmarshal(bytes, &config); err != nil {

		panic(err)
	}
	return config
}

func archiveFile(code, fileName, path string) error {
	if err := os.MkdirAll(path, 0700); err != nil {
		return err
	}
	filePath := filepath.Join(path, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(code)
	return nil
}

func archive() {
	config := load_config()
	resp, err := http.Get(ATCODER_API_SUBMISSION_URL + config.Atcoder.UserID)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var ss []AtCoderSubmission
	err = json.Unmarshal(bytes, &ss)
	if err != nil {
		panic(err)
	}

	//only ac
	ss = funk.Filter(ss, func(s AtCoderSubmission) bool {
		return s.Result == "AC"
	}).([]AtCoderSubmission)

	//rev sort by EpochSecond
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].EpochSecond > ss[j].EpochSecond
	})

	//filter latest submission for each problem
	v := map[string]struct{}{}
	ss = funk.Filter(ss, func(s AtCoderSubmission) bool {
		_, ok := v[s.ContestID+"_"+s.ProblemID]
		if ok {
			return false
		}
		v[s.ContestID+"_"+s.ProblemID] = struct{}{}
		return true
	}).([]AtCoderSubmission)
	startTime := time.Now()
	funk.ForEach(ss, func(s AtCoderSubmission) {
		url := fmt.Sprintf("https://atcoder.jp/contests/%s/submissions/%s", s.ContestID, strconv.Itoa(s.ID))

		log.Printf("Requesting... %s", url)
		elapsedTime := time.Now().Sub(startTime)
		if elapsedTime.Milliseconds() < 800 {
			sleepTime := time.Duration(800 - elapsedTime.Milliseconds())
			time.Sleep(time.Millisecond * sleepTime)
		}
		resp, err := http.Get(url)
		startTime = time.Now()
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatal(err)
			return
		}
		language := s.Language
		contestID := s.ContestID
		problemID := s.ProblemID

		doc.Find(".linenums").Each(func(i int, s *goquery.Selection) {
			code := s.Text()
			if code == "" {
				log.Print("Empty string...")
				return
			}
			fileName := languageToFileName(language)
			archiveDirPath := filepath.Join(config.Atcoder.RepositoryPath, contestID, problemID)

			if err = archiveFile(code, fileName, archiveDirPath); err != nil {
				log.Println("Fail to archive the code", filepath.Join(archiveDirPath, fileName))
				return
			}
			log.Println("Success!! archive the code", filepath.Join(archiveDirPath, fileName))
			return
		})
	})
}
func validateConfig(config Config) bool {
	//TODO check path
	return false
}
func edit() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	configFile := filepath.Join(home, "."+APP_NAME, "config.json")
	editor := os.Getenv("EDITOR")
	if editor != "" {
		c := exec.Command(editor, configFile)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Run()
	} else {
		open.Run(configFile)
	}

}

func main() {

	app := cli.App{Name: "procon-gardener", Usage: "archive your AC submissions",
		Commands: []*cli.Command{
			{
				Name:    "archive",
				Aliases: []string{"a"},
				Usage:   "archive your AC submissions",
				Action: func(c *cli.Context) error {
					archive()
					return nil
				},
			},
			{
				Name:    "edit",
				Aliases: []string{"e"},
				Usage:   "edit your config file",
				Action: func(c *cli.Context) error {

					edit()
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
