package main

import (
    "1BRC/src/generate"
    "github.com/alecthomas/kingpin/v2"
    "fmt"
    "os"
    "bufio"
)

var (
    debug = kingpin.Flag("debug", "Enable debug mode").Short('d').Bool()
    rows = kingpin.Flag("rows", "Number of Rows to output").Default("30").Int()
    
    citySet = kingpin.Flag("station-set", "Select charset for station name").Default("basic").String()
    useStdout = kingpin.Flag("use-stdout", "Use stdout instead of output file. Overrides `output-file`").Default("false").Bool()
    outputFile = kingpin.Arg("output-file", "Output file path").Default("./input.txt").String()
)

func main() {
    kingpin.Parse()
    if *debug == true{
        fmt.Printf("Debug mode is on\n")
        fmt.Printf("Rows is %d\n", *rows)
        fmt.Printf("Output file path is %s\n", *outputFile)
    }
    if *useStdout == true {
        generate.GenerateReal(*rows, os.Stdout)
    } else {
        w, err := os.OpenFile(*outputFile, os.O_WRONLY | os.O_CREATE, 0666)
        if err != nil {
            panic(err)
        }
        defer w.Close()
        bufw := bufio.NewWriterSize(w, 4096 * 1024)
        defer bufw.Flush()
        gen := generate.ParseGen(*citySet)
        generate.Generate(*rows, max(min(10000, *rows / 12), 30), gen, bufw)
    }
}
