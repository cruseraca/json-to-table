package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"reflect"

	"github.com/cruseraca/json-to-table/models"
	"github.com/labstack/echo/v4"
	"gitlab.com/c0b/go-ordered-json"
)

type checkJson struct {
}

type CheckJsonHandler interface {
	CheckJson(c echo.Context) error
	GenerateTable(c echo.Context) error
	GenerateTableOrder(c echo.Context) error
}

func NewCheckJsonHandler() CheckJsonHandler {
	return &checkJson{}
}

func (h *checkJson) CheckJson(c echo.Context) error {

	//validate json
	var bodyJson map[string]any

	err := c.Bind(&bodyJson)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid JSON type",
		})
	}

	//count length
	numOfLength := len(bodyJson)

	//count depth of json
	depth := maxDepth(bodyJson, 0)

	//get listed name - value pairs
	fields := getNameValuePairs(bodyJson, "", []models.Field{})

	responseData := models.CheckJsonResponse{
		NumberOfFields: numOfLength,
		MaximumDepth:   depth,
		ListOfFields:   fields,
	}

	responseApp := models.ResponseCheckJson{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "JSON is valid",
		Data:            responseData,
	}

	return c.JSON(http.StatusOK, responseApp)
}

func (h *checkJson) GenerateTable(c echo.Context) error {

	//validate json
	var bodyJson map[string]any

	err := c.Bind(&bodyJson)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid JSON type",
		})
	}

	//define FuncMap
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}

	//parse template
	tmpl, err := template.New("templateTable.html").Funcs(funcMap).ParseFiles("views/templateTable.html")
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, models.Response{
			ResponseCode:    http.StatusInternalServerError,
			ResponseMessage: "template not found",
		})
	}

	//get listed name - value pairs
	fields := getNameValuePairs(bodyJson, "", []models.Field{})

	responseData := models.CheckJsonResponse{
		ListOfFields: fields,
	}

	err = tmpl.Execute(c.Response(), responseData.ListOfFields)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, models.Response{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "error parsing template",
		})
	}
	return nil
}

func (h *checkJson) GenerateTableOrder(c echo.Context) error {

	//validate json
	bodyJson := ordered.NewOrderedMap()

	err := c.Bind(&bodyJson)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid JSON type",
		})
	}

	//define FuncMap
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}

	//parse template
	tmpl, err := template.New("templateTable.html").Funcs(funcMap).ParseFiles("views/templateTable.html")
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, models.Response{
			ResponseCode:    http.StatusInternalServerError,
			ResponseMessage: "template not found",
		})
	}

	//get listed name - value pairs
	fields := getNameValuePairOrder(bodyJson, "", []models.Field{})

	responseData := models.CheckJsonResponse{
		ListOfFields: fields,
	}

	err = tmpl.Execute(c.Response(), responseData.ListOfFields)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, models.Response{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "error parsing template",
		})
	}
	
	return nil
}

func getNameValuePairOrder(bodyJson *ordered.OrderedMap, prefix string, fields []models.Field) ([]models.Field) {
	kvPairs := bodyJson.EntriesIter()
	for {
		pair, ok := kvPairs()
		if !ok {
			break
		}
		switch v := pair.Value.(type) {
		case *ordered.OrderedMap:
			fields = getNameValuePairOrder(v, pair.Key + " >> ", fields)
		default:
			fields = append(fields, models.Field{Name: prefix + pair.Key, Value: v})
		}
	}
	return fields
}

func maxDepth(data map[string]any, currentDepth int) int {
	max := currentDepth
	for _, value := range data {
		if nestedMap, ok := value.(map[string]interface{}); ok {
			depth := maxDepth(nestedMap, currentDepth+1)
			if depth > max {
				max = depth
			}
		}
	}
	return max
}

func getNameValuePairs(data map[string]any, prefix string, fields []models.Field) []models.Field {
	for key, value := range data {
		fieldName := prefix + key
		switch v := value.(type) {
		case map[string]interface{}:
			fields = append(fields, models.Field{Name: fieldName, Type: reflect.TypeOf(v).String()})
			fields = getNameValuePairs(v, fieldName+">>", fields)
		case []interface{}:
			//check length
			if len(v) > 0 {
				if _, ok := v[0].(map[string]interface{}); ok {
					fields = append(fields, models.Field{Name: fieldName, Type: "array of " + reflect.TypeOf(v[0]).String()})
					//check the most inner fields
					indexMost := searchMaximumField(v)
					arrayOfObject := v[indexMost].(map[string]interface{})
					fields = getNameValuePairs(arrayOfObject, fieldName+">>", fields)
					fields = append(fields, models.Field{Name: fieldName, Type: "array of " + reflect.TypeOf(v[0]).String(), Value: fmt.Sprintf("%v", v)})
				} else {
					fields = append(fields, models.Field{Name: fieldName, Type: "array of " + reflect.TypeOf(v[0]).String(), Value: fmt.Sprintf("%v", v)})
				}
			}
		default:
			fields = append(fields, models.Field{Name: fieldName, Type: reflect.TypeOf(v).String(), Value: v})
		}
	}
	return fields
}

func searchMaximumField(array []any) int {
	max := 0
	indexMost := 0
	for index, value := range array {
		if nestedMap, ok := value.(map[string]interface{}); ok {
			if len(nestedMap) <= max {
				continue
			} else {
				max = index
				indexMost = index
			}
		}
	}
	return indexMost
}
