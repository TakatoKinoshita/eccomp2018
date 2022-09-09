package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/term"
)

type Output struct {
	Obj [][]float64
	Con [][]float64
	Err []string
}

func main() {
	var err error
	env := os.Getenv("EVAL_MODULE")

	//read input
	var input []byte
	if input, err = readInput(); err != nil {
		msgs := []string{"System Error on Reading Input (readInput)", err.Error()}
		outputError(msgs)
		os.Exit(1)
	}

	//read schema
	var schema []byte
	if schema, err = readSchema(fmt.Sprintf("./schema/%s.json", env)); err != nil {
		msgs := []string{"System Error on Reading Input Schema (readSchema)", err.Error()}
		outputError(msgs)
		os.Exit(1)
	}

	//validation
	schemaLoader := gojsonschema.NewBytesLoader(schema)
	documentLoader := gojsonschema.NewBytesLoader(input)
	var result *gojsonschema.Result
	if result, err = gojsonschema.Validate(schemaLoader, documentLoader); err != nil {
		msgs := []string{"System Error on Validating", err.Error()}
		outputError(msgs)
		os.Exit(1)
	}
	if !result.Valid() {
		msgs := make([]string, len(result.Errors())+1)
		msgs[0] = "Input is invalid"
		for i, e := range result.Errors() {
			msgs[i+1] = e.String()
		}
		outputError(msgs)
		os.Exit(10)
	}

	//convert json string to golang data
	var inputMat [][]float64
	if inputMat = convertToFloat(input); inputMat == nil {
		msgs := []string{"System Error on Converting Input (convertToFloat)"}
		outputError(msgs)
		os.Exit(1)
	}

	//save input data as pop_vars_eval.txt
	if err = saveCSV(inputMat, "pop_vars_eval.txt", '\t'); err != nil {
		msgs := []string{"System Error on Writing Input Data (saveCSV)", err.Error()}
		outputError(msgs)
		os.Exit(1)
	}

	//execute evaluation module
	if err = exec.Command("./"+env, ".").Run(); err != nil {
		msgs := []string{"System Error on Running Evaluation Module", err.Error()}
		outputError(msgs)
		os.Exit(20)
	}

	//read data from evaluation module
	var objMat, conMat [][]float64
	if objMat, err = readCSV("pop_objs_eval.txt", '\t'); err != nil {
		msgs := []string{"System Error on Reading Output Data (readCSV)", err.Error()}
		outputError(msgs)
		os.Exit(1)
	}
	if conMat, err = readCSV("pop_cons_eval.txt", '\t'); err != nil {
		msgs := []string{"System Error on Reading Output Data (readCSV)", err.Error()}
		outputError(msgs)
		os.Exit(1)
	}

	//output
	output := Output{objMat, conMat, nil}
	outputStr, _ := json.Marshal(output)
	fmt.Println(string(outputStr))
}

func readInput() ([]byte, error) {
	//debug case
	if term.IsTerminal(int(syscall.Stdin)) {
		s := bufio.NewScanner(os.Stdin)
		s.Scan()
		input := s.Bytes()
		return input, nil
	}

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}

	return input, nil
}

func readSchema(filepath string) ([]byte, error) {
	fp, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	schema, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, err
	}

	return schema, nil
}

func outputError(msgs []string) {
	out := Output{nil, nil, msgs}
	jsonStr, _ := json.Marshal(out)
	fmt.Println(string(jsonStr))
}

func convertToFloat(jsonStr []byte) [][]float64 {
	var vec []float64
	if err := json.Unmarshal(jsonStr, &vec); err == nil {
		mat := make([][]float64, 1)
		mat[0] = vec
		return mat
	}

	var mat [][]float64
	if err := json.Unmarshal(jsonStr, &mat); err == nil {
		return mat
	}

	return nil
}

func floatToString(floatMat [][]float64) [][]string {
	strMat := make([][]string, len(floatMat))

	for i, row := range floatMat {
		strMat[i] = make([]string, len(row))
		for j, num := range row {
			strMat[i][j] = strconv.FormatFloat(num, 'f', -1, 64)
		}
	}

	return strMat
}

func stringToFloat(strMat [][]string) ([][]float64, error) {
	floatMat := make([][]float64, len(strMat))

	for i, row := range strMat {
		floatMat[i] = make([]float64, len(row))
		for j, str := range row {
			var err error
			if floatMat[i][j], err = strconv.ParseFloat(str, 64); err != nil {
				return nil, err
			}
		}
	}

	return floatMat, nil
}

func saveCSV(data [][]float64, filepath string, delim rune) error {
	fp, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer fp.Close()

	strData := floatToString(data)
	w := csv.NewWriter(fp)
	w.Comma = delim
	w.WriteAll(strData)
	if err := w.Error(); err != nil {
		return err
	}

	if err := fp.Sync(); err != nil {
		return err
	}
	return nil
}

func readCSV(filepath string, delim rune) ([][]float64, error) {
	fp, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	var strData [][]string
	r := csv.NewReader(fp)
	r.Comma = delim
	if strData, err = r.ReadAll(); err != nil {
		return nil, err
	}

	var floatData [][]float64
	if floatData, err = stringToFloat(strData); err != nil {
		return nil, err
	}

	return floatData, nil
}
