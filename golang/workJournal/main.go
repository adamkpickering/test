package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

type JSONStringDuration time.Duration

func (d JSONStringDuration) MarshalJSON() ([]byte, error) {
	asDuration := time.Duration(d)
	output := fmt.Sprintf("%q", asDuration)
	return []byte(output), nil
}

type Day struct {
	ToDo  []string
	Done  []string
	Tasks []Task
}

var toDoLineRegex *regexp.Regexp
var doneLineRegex *regexp.Regexp
var taskLineRegex *regexp.Regexp
var taskRegex *regexp.Regexp

func init() {
	toDoLineRegex = regexp.MustCompile("\nTo Do\n")
	doneLineRegex = regexp.MustCompile("\nDone\n")
	taskLineRegex = regexp.MustCompile("\n[0-9]{2}:[0-9]{2} .*?\n")
	taskRegex = regexp.MustCompile("(?s)([0-9]{2}:[0-9]{2}) ([a-zA-Z0-9,\\_\\-]+) (.*?)\n(.*)$")
}

// Constructs a Task and reads relevant fields from a Task title line
// into that Task. If the title line does not contain info on
// a given field of the Task, that field is left as its zero value.
func partialTaskFromTitleLine(line string) (Task, error) {
	task := Task{}

	parts := strings.SplitN(line, " ", 2)
	if len(parts) != 2 {
		return Task{}, fmt.Errorf("failed to split line %q", line)
	}

	rawStartTime := parts[0]
	parsedTime, err := time.Parse("15:04", rawStartTime)
	if err != nil {
		return Task{}, fmt.Errorf("failed to parse time %q: %w", rawStartTime, err)
	}
	task.StartTime = parsedTime

	task.Title = parts[1]

	return task, nil
}

func printDay(day Day) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(day)
	if err != nil {
		panic(err)
	}
}

func printSummary(day Day) {
	for _, task := range day.Tasks {
		fmt.Printf("%s\t\t%s\n", time.Duration(task.Duration), task.Title)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("must pass exactly one file to parse as arg")
		os.Exit(1)
	}
	fileName := os.Args[1]
	rawContents, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	contents := string(rawContents)
	day, err := parseDay(contents)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// printDay(day)
	printSummary(day)
	// fmt.Printf("%#v\n", day)
}

func parseDay(contents string) (Day, error) {
	toDoIndices := toDoLineRegex.FindStringIndex(contents)
	if toDoIndices == nil {
		return Day{}, errors.New("no match for To Do regex")
	}
	doneIndices := doneLineRegex.FindStringIndex(contents)
	if doneIndices == nil {
		return Day{}, errors.New("no match for Done regex")
	}
	taskLineIndexPairs := taskLineRegex.FindAllStringIndex(contents, -1)
	if len(taskLineIndexPairs) == 0 {
		return Day{}, errors.New("no match for task line regex")
	}

	day := Day{}
	toDoContents := contents[toDoIndices[1]:doneIndices[0]]
	day.ToDo = parseDashList(toDoContents)
	doneContents := contents[doneIndices[1]:taskLineIndexPairs[0][0]]
	day.Done = parseDashList(doneContents)

	for i := range taskLineIndexPairs {
		var taskContents string
		if i+1 == len(taskLineIndexPairs) {
			// the last pair is a special case
			taskContents = contents[taskLineIndexPairs[i][0]:]
		} else {
			taskContents = contents[taskLineIndexPairs[i][0]:taskLineIndexPairs[i+1][0]]
		}
		task := Task{}
		if err := task.UnmarshalText([]byte(taskContents)); err != nil {
			return Day{}, fmt.Errorf("failed to parse task: %w", err)
		}
		day.Tasks = append(day.Tasks, task)
	}

	// Set the duration of each task. Skip the last task, allowing
	// its duration to remain set to 0.
	for i := 0; i < len(day.Tasks)-1; i++ {
		duration := day.Tasks[i+1].StartTime.Sub(day.Tasks[i].StartTime)
		day.Tasks[i].Duration = JSONStringDuration(duration)
	}

	return day, nil
}

func parseDashList(dashListText string) []string {
	lines := strings.Split(strings.TrimSpace(dashListText), "\n")
	dashList := make([]string, 0, len(lines))
	for _, line := range lines {
		dashListEntry := strings.TrimPrefix(line, "- ")
		dashList = append(dashList, dashListEntry)
	}
	return dashList
}

type Task struct {
	Title     string
	Duration  JSONStringDuration
	StartTime time.Time
	Tags      []string
	Body      string
}

func (task *Task) UnmarshalText(text []byte) error {
	submatches := taskRegex.FindSubmatch(text)
	if submatches == nil {
		return fmt.Errorf("no match for task regex for text %q", bytes.TrimSpace(text)[:40])
	}

	rawStartTime := string(submatches[1])
	parsedTime, err := time.Parse("15:04", rawStartTime)
	if err != nil {
		return fmt.Errorf("failed to parse time %q: %w", rawStartTime, err)
	}
	task.StartTime = parsedTime

	task.Tags = strings.Split(string(submatches[2]), ",")
	task.Title = string(submatches[3])
	task.Body = strings.TrimSpace(string(submatches[4]))

	return nil
}
