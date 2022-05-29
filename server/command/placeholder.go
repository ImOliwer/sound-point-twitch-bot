package command

import "fmt"

type PlaceholderFunc func(*Context) any

type PlaceholderIndex struct {
	Start   int
	End     int
	Content string
}

//////////////////////
//      SHARED      //
//////////////////////

var points_placeholders = map[string]PlaceholderFunc{
	"target": func(ctx *Context) any {
		return ctx.Arguments[0]
	},
	"points": func(ctx *Context) any {
		value, ok := ctx.Temp["response-points"]
		if !ok {
			return 0
		}
		return value
	},
}

//////////////////////
//    PROCESSING    //
//////////////////////

func (r *Context) process_placeholders(in string, placeholders map[string]PlaceholderFunc) string {
	population := make([]PlaceholderIndex, 0)
	firstCharIndex := -1

	for index, char := range in {
		switch char {
		case '{':
			firstCharIndex = index
		case '}':
			if firstCharIndex == -1 {
				continue
			}
			population = append(population, PlaceholderIndex{
				Start:   firstCharIndex,
				End:     index,
				Content: in[firstCharIndex+1 : index],
			})
			firstCharIndex = -1
		default:
			continue
		}
	}

	for _, placeholder := range population {
		function := placeholders[placeholder.Content]
		if function == nil {
			continue
		}
		in = fmt.Sprintf("%v%v%v", in[:placeholder.Start], function(r), in[placeholder.End+1:])
	}

	return in
}
