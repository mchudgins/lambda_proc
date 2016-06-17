package main

import (
	"log"
	"strings"

	"github.com/mchudgins/lambda_proc"
)

type WidgetList struct {
	Items []string `json:items`
}

func ProcessWidgets(context *lambda_proc.Context, evt APIEvent) (interface{}, error) {
	switch strings.ToUpper(evt.Method) {
	case "GET":
		return getWidgetList(context, evt)
	case "POST":

	default:
	}

	return nil, nil
}

func getWidgetList(context *lambda_proc.Context, evt APIEvent) (WidgetList, error) {
	l, err := Lists.Fetch("0", "ToDo's")
	if err != nil {
		return WidgetList{}, err
	}

	log.Printf("List: %+v\n", l)

	list := WidgetList{Items: l.Items}

	return list, nil

}
