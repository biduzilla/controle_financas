package utils

import (
	"encoding/json"
	"maps"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Envelope map[string]any

func MinifySQL(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func WriteJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')
	maps.Copy(w.Header(), headers)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func GetTypeName(v any) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return strings.ToLower(t.Name())
}

func ValidateTelefone(telefone string) bool {
	telefone = regexp.MustCompile(`[^\d]`).ReplaceAllString(telefone, "")

	match, _ := regexp.MatchString(`^\d{2}9\d{8}$`, telefone)
	return match
}

func ValidateDate(data string) bool {
	pattern := `^(0[1-9]|[12][0-9]|3[01])/(0[1-9]|1[012])/(19|20)\d\d$`
	matched, _ := regexp.MatchString(pattern, data)
	if !matched {
		return false
	}

	parts := strings.Split(data, "/")
	if len(parts) != 3 {
		return false
	}

	day, _ := strconv.Atoi(parts[0])
	month, _ := strconv.Atoi(parts[1])
	year, _ := strconv.Atoi(parts[2])

	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	return date.Year() == year && int(date.Month()) == month && date.Day() == day
}

func ValidateDateISO(data string) bool {
	pattern := `^(19|20)\d\d-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])$`
	matched, _ := regexp.MatchString(pattern, data)
	if !matched {
		return false
	}

	parts := strings.Split(data, "-")
	if len(parts) != 3 {
		return false
	}

	year, _ := strconv.Atoi(parts[0])
	month, _ := strconv.Atoi(parts[1])
	day, _ := strconv.Atoi(parts[2])

	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	return date.Year() == year && int(date.Month()) == month && date.Day() == day
}
