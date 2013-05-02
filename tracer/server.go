package tracer

import (
	"flag"
	"fmt"
	"github.com/deboshire/exp/math/vector"
	"go/build"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"path"
	"runtime"
	"sort"
	"sync"
)

var (
	data       = make(map[string][]float64)
	mu         = &sync.Mutex{}
	tracerPort = flag.Int("tracer-port", 1234, "Port to use for tracer")

	templateDir string
	staticDir   string
)

type webTracer struct {
	prefix string
}

func NewWebTracer() Tracer {
	result := &webTracer{}
	go result.start()
	return result
}

type IndexPageData struct {
	Names []string
}

func executeTemplate(name string, data interface{}, w io.Writer) {
	templateName := path.Join(templateDir, name)

	t := template.Must(template.ParseFiles(templateName))
	err := t.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

type VarPageData struct {
	Name   string
	Values []float64
}

func handleVar(w http.ResponseWriter, r *http.Request) {
	var pageData VarPageData
	pageData.Name = r.URL.Query().Get("n")

	{
		mu.Lock()
		defer mu.Unlock()

		pageData.Values = data[pageData.Name]
	}

	executeTemplate("v.html", &pageData, w)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	var pageData IndexPageData
	{
		mu.Lock()
		defer mu.Unlock()

		// This is unbelievable that in XXI century you have to manually copy keys out.
		names := make([]string, len(data))
		i := 0
		for k, _ := range data {
			names[i] = k
			i++
		}
		sort.Strings(names)
		pageData.Names = names
	}

	executeTemplate("index.html", &pageData, w)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		handleIndex(w, r)
	case "/v":
		handleVar(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (t *webTracer) TraceFloat64(label string, value float64) {
	mu.Lock()
	defer mu.Unlock()

	key := t.prefix + "." + label
	data[key] = append(data[key], value)
}

func (t *webTracer) TraceF64(label string, value vector.F64) {
}

func (t *webTracer) TraceInt(label string, value int) {
}

func (t *webTracer) Sub(label string) Tracer {
	if len(t.prefix) > 0 {
		label = "." + label
	}
	return &webTracer{prefix: t.prefix + label}
}

func (t *webTracer) Algorithm(name string) Tracer {
	return t.Sub(fmt.Sprintf("%s.%x", name, rand.Int()))
}

func (t *webTracer) Iter(i int64) Tracer {
	return t
}

func (t *webTracer) LastIter(i int64) Tracer {
	return t.Iter(i)
}

func (t *webTracer) start() {
	http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir(staticDir))))
	http.HandleFunc("/", handleRequest)
	fmt.Printf("Tracer started on http://localhost:%d/\n", *tracerPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *tracerPort), nil); err != nil {
		panic(err)
	}
}

func initDirs() {
	ctx := build.Default
	// todo(mike): can you reference current package somehow?
	pkg, err := ctx.Import("github.com/deboshire/exp/tracer", "", build.FindOnly)
	if err != nil {
		panic(err)
	}
	dir := pkg.Dir
	templateDir = path.Join(dir, "t")
	staticDir = path.Join(dir, "s")
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	initDirs()
}
